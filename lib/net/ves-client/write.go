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

	"github.com/HyperService-Consortium/go-ves/grpc/uiprpc"
	uipbase "github.com/HyperService-Consortium/go-ves/grpc/uiprpc-base"
	"github.com/HyperService-Consortium/go-ves/grpc/wsrpc"

	"github.com/HyperService-Consortium/go-uip/signaturer"

	helper "github.com/HyperService-Consortium/go-ves/lib/net/help-func"
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
			vc.logger.Error("error found", "error", err)
			return
		}

		buf = bytes.NewBuffer(bytes.TrimSpace(strBytes))

		cmdBytes, err = buf.ReadBytes(' ')
		if err != nil && err != io.EOF {
			vc.logger.Error("error found", "error", err)
			return
		}

		if err = vc.SendEthAlias(
			[]byte("eth1"),
		); err != nil {
			vc.logger.Error("error found", "error", err)
			continue
		} else {
			fmt.Println("success")
		}
		switch string(bytes.TrimSpace(cmdBytes)) {
		case "set-name":
			vc.name, err = buf.ReadBytes(' ')
			if err != nil && err != io.EOF {
				vc.logger.Error("error found", "error", err)
				continue
			}
			vc.name = bytes.TrimSpace(vc.name)
			if err = vc.SayClientHello(vc.name); err != nil {
				vc.logger.Error("error found", "error", err)
				continue
			}

		case "send-to":
			toBytes, err = buf.ReadBytes(' ')
			if err != nil && err != io.EOF {
				vc.logger.Error("error found", "error", err)
				continue
			}
			if err = vc.SendMessage(
				bytes.TrimSpace(toBytes),
				bytes.TrimSpace(buf.Bytes()),
			); err != nil {
				vc.logger.Error("error found", "error", err)
				continue
			}
		case "register-key":
			filePath, err = buf.ReadBytes(' ')
			if err != nil && err != io.EOF {
				vc.logger.Error("error found", "error", err)
				continue
			}

			if err = vc.ConfigKey(
				string(bytes.TrimSpace(filePath)),
				fileBuffer,
			); err != nil {
				vc.logger.Error("error found", "error", err)
				continue
			}
		case "register-eth":
			filePath, err = buf.ReadBytes(' ')
			if err != nil && err != io.EOF {
				vc.logger.Error("error found", "error", err)
				continue
			}

			if err = vc.ConfigEth(
				string(bytes.TrimSpace(filePath)),
				fileBuffer,
			); err != nil {
				vc.logger.Error("error found", "error", err)
				continue
			}
		case "send-eth-alias-to-ves":
			alias, err = buf.ReadBytes(' ')
			fmt.Println("....", string(alias))
			if err != nil && err != io.EOF {
				vc.logger.Error("error found", "error", err)
				continue
			}
			if err = vc.SendEthAlias(
				bytes.TrimSpace(alias),
			); err != nil {
				vc.logger.Error("error found", "error", err)
				continue
			}
		case "send-alias-to-ves":
			alias, err = buf.ReadBytes(' ')
			if err != nil && err != io.EOF {
				vc.logger.Error("error found", "error", err)
				continue
			}
			if err = vc.SendAlias(
				bytes.TrimSpace(alias),
			); err != nil {
				vc.logger.Error("error found", "error", err)
				continue
			}
		case "keys":
			vc.ListKeys()
		case "send-op-intents":
			filePath, err = buf.ReadBytes(' ')
			if err != nil && err != io.EOF {
				vc.logger.Error("error found", "error", err)
				continue
			}

			if err = vc.SendOpIntents(
				string(bytes.TrimSpace(filePath)),
				fileBuffer,
			); err != nil {
				vc.logger.Error("error found", "error", err)
				continue
			}
		}

	}
}

func (vc *VesClient) ConfigKey(filePath string, fileBuffer []byte) error {
	file, err := os.Open(filePath)
	if err != nil {
		vc.logger.Error("open file error", "error", err)
		return err
	}
	var n int
	n, err = io.ReadFull(file, fileBuffer)
	file.Close()
	if err != nil && err != io.ErrUnexpectedEOF {
		vc.logger.Error("read error", "error", err)
		return err
	}
	var ks = make([]*ECCKeyAlias, 0)
	err = json.Unmarshal(fileBuffer[0:n], &ks)
	if err != nil {
		vc.logger.Error("unmarshal error", "error", err)
		return err
	}
	var flag bool
	for _, kk := range ks {
		flag = false
		// todo: check

		b, err := hex.DecodeString(kk.PrivateKey)
		if err != nil {
			vc.logger.Error("decode private key error", "error", err)
			return err
		}

		k := ECCKey{PrivateKey: b, ChainID: kk.ChainID}
		for _, key := range vc.keys.Keys {
			if key.ChainID == k.ChainID && bytes.Equal(key.PrivateKey, k.PrivateKey) {
				vc.logger.Info("this key is already in the storage, private key", "address", hex.EncodeToString(k.PrivateKey[0:8]))
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
		vc.logger.Info("imported: private key", "address", hex.EncodeToString(k.PrivateKey[0:8]), ", chain_id", k.ChainID)
	}

	return nil
}

func (vc *VesClient) ConfigEth(filePath string, fileBuffer []byte) error {
	file, err := os.Open(filePath)
	if err != nil {
		vc.logger.Error("open file error", "error", err)
		return err
	}

	var n int
	n, err = io.ReadFull(file, fileBuffer)
	file.Close()
	if err != nil && err != io.ErrUnexpectedEOF {
		vc.logger.Error("read error", "error", err)
		return err
	}
	var as = make([]*EthAccountAlias, 0)
	err = json.Unmarshal(fileBuffer[0:n], &as)
	if err != nil {
		vc.logger.Error("unmarshal error", "error", err)
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

				vc.logger.Info("this account is already in the storage, public address", "address", a.Address[0:8])
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
		vc.logger.Info("imported: public address", "address", a.Address[0:8], ", chain_id", a.ChainID)
	}
	return nil
}

func (vc *VesClient) SendEthAlias(alias []byte) error {
	if acc, ok := vc.accs.Alias[*(*string)(unsafe.Pointer(&alias))]; ok {
		vc.logger.Error("registering", "alias", string(alias))
		userRegister := vc.getUserRegisterRequest()
		b, _ := hex.DecodeString(acc.Address)
		userRegister.Account = &uipbase.Account{Address: b, ChainId: acc.ChainID}
		userRegister.UserName = *(*string)(unsafe.Pointer(&vc.name))
		err := vc.postMessage(wsrpc.CodeUserRegisterRequest, userRegister)
		if err != nil {
			vc.logger.Error("register user error", "alias", string(alias), "error", err)
			return err
		}
		vc.logger.Error("register success", "alias", string(alias), "addr", hex.EncodeToString(b))
		return nil
	}
	vc.logger.Error("find error", "alias", string(alias), "error", errNotFound)
	return errNotFound
}


func (vc *VesClient) SendAlias(alias []byte) error {
	if key, ok := vc.keys.Alias[*(*string)(unsafe.Pointer(&alias))]; ok {
		userRegister := vc.getUserRegisterRequest()

		signer, err := signaturer.NewTendermintNSBSigner(key.PrivateKey)
		if err != nil {
			return err
		}
		if signer == nil {
			vc.logger.Error("init signer error", "alias", key.PrivateKey, "error", errIlegalPrivateKey)
			return errIlegalPrivateKey
		}
		userRegister.Account = &uipbase.Account{Address: signer.GetPublicKey(), ChainId: key.ChainID}
		userRegister.UserName = *(*string)(unsafe.Pointer(&vc.name))
		err = vc.postMessage(wsrpc.CodeUserRegisterRequest, userRegister)
		if err != nil {
			vc.logger.Error("register user error", "alias", string(alias), "error", err)
			return err
		}
		return nil
	}
	vc.logger.Error("find error", "alias", string(alias), "error", errNotFound)
	return errNotFound
}

type opIntents struct {
	Intents      []json.RawMessage `json:"op-intents"`
	Dependencies []json.RawMessage `json:"dependencies"`
}

func convRaw(rs []json.RawMessage) (ret [][]byte) {
	for _, rawMessage := range rs {
		ret = append(ret, []byte(rawMessage))
	}
	return ret
}

func (vc *VesClient) readOpIntents(filePath string, fileBuffer []byte)  (opIntents, error) {
	var intents opIntents
	file, err := os.Open(filePath)
	if err != nil {
		vc.logger.Error("open file error", "error", err)
		return intents, err
	}

	var n int
	n, err = io.ReadFull(file, fileBuffer)
	_ = file.Close()
	if err != nil && err != io.ErrUnexpectedEOF {
		vc.logger.Error("read error", "error", err)
		return intents, err
	}

	err = json.Unmarshal(fileBuffer[:n], &intents)
	if err != nil {
		vc.logger.Error("unmarshal error", "error", err)
		return intents, fmt.Errorf("Unmarshal failed: %v", err)
	}

	vc.logger.Info("read op intents from file successfully", "intent count", len(intents.Intents), "dependencies count", len(intents.Dependencies))
	return intents, nil
}

func (vc *VesClient) SendOpIntents(filePath string, fileBuffer []byte) error {
	intents, err := vc.readOpIntents(filePath, fileBuffer)
	if err != nil {
		return err
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(mAddress, grpc.WithInsecure())
	if err != nil {
		vc.logger.Error("did not connect", "error", err)
		return err
	}
	defer conn.Close()
	c := uiprpc.NewVESClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	r, err := c.SessionStart(
		ctx,
		&uiprpc.SessionStartRequest{
			Opintents: &uipbase.OpIntents{
				Dependencies: convRaw(intents.Dependencies),
				Contents:     convRaw(intents.Intents),
			},
		})
	if err != nil {
		vc.logger.Error("could not greet", "error", err)
		return err
	}
	vc.logger.Info("Session Start, %v", "ok", r.GetOk(), "session id", hex.EncodeToString(r.GetSessionId()))
	return nil
}

func (vc *VesClient) GetRawTransaction(sessionID, host []byte) (
	*uiprpc.SessionRequireRawTransactReply, error,
) {
	mhost, err := helper.DecodeIP(host)
	if err != nil {
		vc.logger.Error("could not decode ip", "error", err)
		return nil, err
	}
	conn, err := grpc.Dial(mhost, grpc.WithInsecure())
	if err != nil {
		vc.logger.Error("did not connect", "error", err)
		return nil, err
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
		vc.logger.Error("could not get raw transaction", "error", err)
		return nil, err
	}
	return r, nil
}
