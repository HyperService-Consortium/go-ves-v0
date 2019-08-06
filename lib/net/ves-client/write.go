package vesclient

import (
	"bufio"
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
	"unsafe"

	"google.golang.org/grpc"

	log "github.com/Myriad-Dreamin/go-ves/lib/log"

	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc/uiprpc"
	uipbase "github.com/Myriad-Dreamin/go-ves/grpc/uiprpc-base"
	wsrpc "github.com/Myriad-Dreamin/go-ves/grpc/wsrpc"

	signaturer "github.com/Myriad-Dreamin/go-uip/signaturer"

	helper "github.com/Myriad-Dreamin/go-ves/lib/net/help-func"
)

func (vc *VesClient) write() {
	var (
		reader                             = bufio.NewReader(os.Stdin)
		cmdBytes, toBytes, filePath, alias []byte
		fileBuffer                         = make([]byte, 65536)
		buf                                *bytes.Buffer
	)
	for {
		strBytes, _, err := reader.ReadLine()
		if err != nil {
			log.Println(err)
			return
		}

		buf = bytes.NewBuffer(bytes.TrimSpace(strBytes))

		cmdBytes, err = buf.ReadBytes(' ')
		if err != nil && err != io.EOF {
			log.Println(err)
			return
		}

		switch string(bytes.TrimSpace(cmdBytes)) {
		case "set-name":
			vc.name, err = buf.ReadBytes(' ')
			if err != nil && err != io.EOF {
				log.Println(err)
				continue
			}
			vc.name = bytes.TrimSpace(vc.name)
			if err = vc.sayClientHello(vc.name); err != nil {
				log.Println(err)
				continue
			}

		case "send-to":
			toBytes, err = buf.ReadBytes(' ')
			if err != nil && err != io.EOF {
				log.Println(err)
				continue
			}
			if err = vc.sendMessage(
				bytes.TrimSpace(toBytes),
				bytes.TrimSpace(buf.Bytes()),
			); err != nil {
				log.Println(err)
				continue
			}
		case "register-key":
			filePath, err = buf.ReadBytes(' ')
			if err != nil && err != io.EOF {
				log.Println(err)
				continue
			}

			if err = vc.registerKey(
				bytes.TrimSpace(filePath),
				fileBuffer,
			); err != nil {
				log.Println(err)
				continue
			}
		case "register-eth":
			filePath, err = buf.ReadBytes(' ')
			if err != nil && err != io.EOF {
				log.Println(err)
				continue
			}

			if err = vc.configEth(
				bytes.TrimSpace(filePath),
				fileBuffer,
			); err != nil {
				log.Println(err)
				continue
			}
		case "send-eth-alias-to-ves":
			alias, err = buf.ReadBytes(' ')
			if err != nil && err != io.EOF {
				log.Println(err)
				continue
			}

			if err = vc.sendEthAlias(
				bytes.TrimSpace(alias),
			); err != nil {
				log.Println(err)
				continue
			}
		case "send-alias-to-ves":
			alias, err = buf.ReadBytes(' ')
			if err != nil && err != io.EOF {
				log.Println(err)
				continue
			}
			if err = vc.sendAlias(
				bytes.TrimSpace(alias),
			); err != nil {
				log.Println(err)
				continue
			}
		case "keys":
			fmt.Println("privatekeys -> publickeys:")
			for alias, key := range vc.keys.Alias {
				fmt.Println(
					"alias:", alias,
					"public key:", hex.EncodeToString(signaturer.NewTendermintNSBSigner(key.PrivateKey).GetPublicKey()),
					"chain id:", key.ChainID,
				)
			}
			fmt.Println("ethAccounts:")
			for alias, acc := range vc.accs.Alias {
				fmt.Println(
					"alias:", alias,
					"public address:", acc.Address,
					"chain id:", acc.ChainID,
				)
			}
		case "send-op-intents":
			filePath, err = buf.ReadBytes(' ')
			if err != nil && err != io.EOF {
				log.Println(err)
				continue
			}

			if err = vc.sendOpIntents(
				bytes.TrimSpace(filePath),
				fileBuffer,
			); err != nil {
				log.Println(err)
				continue
			}
		}

	}
}

func (vc *VesClient) registerKey(filePath, fileBuffer []byte) error {
	file, err := os.Open(string(filePath))
	if err != nil {
		return err
	}
	var n int
	n, err = io.ReadFull(file, fileBuffer)
	file.Close()
	if err != nil && err != io.ErrUnexpectedEOF {
		return err
	}
	var ks = make([]*ECCKeyAlias, 0)
	err = json.Unmarshal(fileBuffer[0:n], &ks)
	if err != nil {
		return err
	}
	var flag bool
	for _, kk := range ks {
		flag = false
		// todo: check

		b, err := hex.DecodeString(kk.PrivateKey)
		if err != nil {
			return err
		}

		k := ECCKey{PrivateKey: b, ChainID: kk.ChainID}
		for _, key := range vc.keys.Keys {
			if key.ChainID == k.ChainID && bytes.Equal(key.PrivateKey, k.PrivateKey) {
				log.Println("this key is already in the storage, private key:", hex.EncodeToString(k.PrivateKey[0:8]))
				flag = true
				break
			}
		}
		if flag {
			continue
		}
		vc.keys.Keys = append(vc.keys.Keys, &k)
		if len(kk.Alias) != 0 {
			vc.keys.Alias[kk.Alias] = k
		}
		log.Println("imported: private key:", hex.EncodeToString(k.PrivateKey[0:8]), ", chain_id: ", k.ChainID)
	}

	return nil
}

func (vc *VesClient) configEth(filePath, fileBuffer []byte) error {
	file, err := os.Open(string(filePath))
	if err != nil {
		return err
	}

	var n int
	n, err = io.ReadFull(file, fileBuffer)
	file.Close()
	if err != nil && err != io.ErrUnexpectedEOF {
		return err
	}
	var as = make([]*EthAccountAlias, 0)
	err = json.Unmarshal(fileBuffer[0:n], &as)
	if err != nil {
		return err
	}
	var flag bool
	for _, a := range as {
		flag = false
		for _, acc := range vc.accs.Accs {
			if acc.ChainID == a.ChainID && acc.Address == a.Address {

				for alias, acc2 := range vc.accs.Alias {
					if acc2.ChainID == a.ChainID && acc2.Address == a.Address {
						delete(vc.accs.Alias, alias)
					}
				}
				if len(a.Alias) != 0 {
					vc.accs.Alias[a.Alias] = a.EthAccount
				}

				if acc.PassPhrase != a.PassPhrase {
					acc.PassPhrase = a.PassPhrase
					break
				}

				log.Println("this account is already in the storage, public address:", a.Address[0:8])
				flag = true
				break
			}
		}
		if flag {
			continue
		}
		vc.accs.Accs = append(vc.accs.Accs, &a.EthAccount)
		if len(a.Alias) != 0 {
			vc.accs.Alias[a.Alias] = a.EthAccount
		}
		log.Println("imported: public address:", a.Address[0:8], ", chain_id: ", a.ChainID)
	}
	return nil
}

func (vc *VesClient) sendEthAlias(alias []byte) error {
	if acc, ok := vc.accs.Alias[*(*string)(unsafe.Pointer(&alias))]; ok {
		userRegister := vc.getUserRegisterRequest()
		b, _ := hex.DecodeString(acc.Address)
		userRegister.Account = &uipbase.Account{Address: b, ChainId: acc.ChainID}
		userRegister.UserName = *(*string)(unsafe.Pointer(&vc.name))
		err := vc.postMessage(wsrpc.CodeUserRegisterRequest, userRegister)
		if err != nil {
			return err
		}
		return nil
		// for {
		// 	select {
		// 	case msgBuf := <-vc.cb:
		// 		var messageID uint16
		// 		binary.Read(msgBuf, binary.BigEndian, &messageID)
		// 		if messageID != wsrpc.CodeUserRegisterReply {
		// 			continue
		// 		}
		// 		var s = vc.getUserRegisterReply()
		// 		err = proto.Unmarshal(msgBuf.Bytes(), s)
		// 		if err != nil {
		// 			// ignoring
		// 			// todo: add hidden log
		// 			continue
		// 		}
		// 		//todo: checkCharacteristicFlag
		// 		if !s.GetOk() {
		// 			return errors.New("register user failed")
		// 		}
		// 		return nil
		// 	case <-time.After(time.Second * 5):
		// 		return errors.New("timeout")
		// 	}
		// }
	}
	return errNotFound
}

func (vc *VesClient) sendAlias(alias []byte) error {
	if key, ok := vc.keys.Alias[*(*string)(unsafe.Pointer(&alias))]; ok {
		userRegister := vc.getUserRegisterRequest()

		signer := signaturer.NewTendermintNSBSigner(key.PrivateKey)
		if signer == nil {
			return errIlegalPrivateKey
		}
		userRegister.Account = &uipbase.Account{Address: signer.GetPublicKey(), ChainId: key.ChainID}
		userRegister.UserName = *(*string)(unsafe.Pointer(&vc.name))

		return vc.postMessage(wsrpc.CodeUserRegisterRequest, userRegister)
	}
	return errNotFound
}

func (vc *VesClient) sayClientHello(name []byte) error {
	clientHello := vc.getClientHello()
	clientHello.Name = name

	return vc.postMessage(wsrpc.CodeClientHelloRequest, clientHello)
}

type opIntents struct {
	Intents      []json.RawMessage `json:"Op-intents"`
	Dependencies []json.RawMessage `json:"dependencies"`
}

func convRaw(rs []json.RawMessage) (ret [][]byte) {
	for _, rawMessage := range rs {
		ret = append(ret, []byte(rawMessage))
	}
	return ret
}

func (vc *VesClient) sendOpIntents(filePath, fileBuffer []byte) error {

	file, err := os.Open(string(filePath))
	if err != nil {
		return err
	}

	var n int
	n, err = io.ReadFull(file, fileBuffer)
	file.Close()
	if err != nil && err != io.ErrUnexpectedEOF {
		return err
	}
	// Set up a connection to the server.
	conn, err := grpc.Dial(mAddress, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("did not connect: %v", err)
	}
	defer conn.Close()
	c := uiprpc.NewVESClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	var intents opIntents
	err = json.Unmarshal(fileBuffer[:n], &intents)
	if err != nil {
		return fmt.Errorf("Unmarshal failed: %v", err)
	}
	fmt.Println(intents)
	r, err := c.SessionStart(
		ctx,
		&uiprpc.SessionStartRequest{
			Opintents: &uipbase.OpIntents{
				Dependencies: convRaw(intents.Dependencies),
				Contents:     convRaw(intents.Intents),
			},
		})
	if err != nil {
		return fmt.Errorf("could not greet: %v", err)
	}
	fmt.Printf("Session Start: %v, %v\n", r.GetOk(), hex.EncodeToString(r.GetSessionId()))
	return nil
}

func (vc *VesClient) getRawTransaction(sessionID, host []byte) (
	[]byte, uint64, *uipbase.Account, *uipbase.Account, error,
) {
	mhost, err := helper.DecodeIP(host)
	if err != nil {
		return nil, 0, nil, nil, fmt.Errorf("could not decode ip: %v", err)
	}
	conn, err := grpc.Dial(mhost, grpc.WithInsecure())
	if err != nil {
		return nil, 0, nil, nil, fmt.Errorf("did not connect: %v", err)
	}
	defer conn.Close()
	c := uiprpc.NewVESClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	r, err := c.SessionRequireRawTransact(
		ctx,
		&uiprpc.SessionRequireRawTransactRequest{
			SessionId: sessionID,
		},
	)
	if err != nil {
		return nil, 0, nil, nil, fmt.Errorf("could not greet: %v", err)
	}
	return r.GetRawTransaction(), r.GetTid(), r.GetSrc(), r.GetDst(), nil
}
