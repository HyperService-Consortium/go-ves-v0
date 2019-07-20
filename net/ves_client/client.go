package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"sync"
	"time"
	"unsafe"

	signaturer "github.com/Myriad-Dreamin/go-uip/signaturer"
	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc/uip-rpc"
	wsrpc "github.com/Myriad-Dreamin/go-ves/grpc/ws-ves-rpc"
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"

	filedb "github.com/Myriad-Dreamin/go-ves/database/filedb"
)

type ECCKey struct {
	PrivateKey []byte `json:"private_key"`
	ChainID    uint64 `json:"chain_id"`
}

type ECCKeyAlias struct {
	PrivateKey string `json:"private_key"`
	ChainID    uint64 `json:"chain_id"`
	Alias      string `json:"alias"`
}

type EthAccount struct {
	Address    string `json:"address"`
	ChainID    uint64 `json:"chain_id"`
	PassPhrase string `json:"pass_phrase"`
}

type EthAccountAlias struct {
	EthAccount
	Alias string `json:"alias"`
}

type ECCKeys struct {
	Keys  []*ECCKey
	Alias map[string]ECCKey
}

type EthAccounts struct {
	Accs  []*EthAccount
	Alias map[string]EthAccount
}

// VesClient is the web socket client interactive with veses
type VesClient struct {
	rwMutex sync.RWMutex

	name []byte

	cb chan *bytes.Buffer

	shortSendMessage          *wsrpc.Message
	shortReplyMessage         *wsrpc.Message
	clientHello               *wsrpc.ClientHello
	clientHelloReply          *wsrpc.ClientHelloReply
	requestComingRequest      *wsrpc.RequestComingRequest
	requestComingReply        *wsrpc.RequestComingReply
	requestGrpcServiceRequest *wsrpc.RequestGrpcServiceRequest
	requestGrpcServiceReply   *wsrpc.RequestGrpcServiceReply
	requestNsbServiceRequest  *wsrpc.RequestNsbServiceRequest
	requestNsbServiceReply    *wsrpc.RequestNsbServiceReply
	userRegisterRequest       *wsrpc.UserRegisterRequest
	userRegisterReply         *wsrpc.UserRegisterReply
	sessionListRequest        *wsrpc.SessionListRequest
	sessionListReply          *wsrpc.SessionListReply
	transactionListRequest    *wsrpc.TransactionListRequest
	transactionListReply      *wsrpc.TransactionListReply

	sessionStart *uiprpc.SessionStartRequest

	conn *websocket.Conn
	fdb  *filedb.FileDB

	keys *ECCKeys
	accs *EthAccounts

	nsbip  string
	grpcip string
}

func NewVesClient(dbpath string) (*VesClient, error) {
	if fdb, err := filedb.NewFileDB(dbpath); err != nil {
		return nil, err
	} else {
		return (&VesClient{
			fdb: fdb,
			cb:  make(chan *bytes.Buffer, 1),
		}).predo()
	}
}

func (vc *VesClient) predo() (*VesClient, error) {
	filedb.Register(&ECCKeys{})
	filedb.Register(&EthAccounts{})
	ev, err := vc.fdb.ReadWithPath("keys")
	if err != nil {
		return nil, err
	}
	err = ev.Decode(vc.keys)
	if err != nil {
		return nil, err
	}
	err = ev.Settle()
	if err != nil {
		return nil, err
	}

	ev, err = vc.fdb.ReadWithPath("accs")
	if err != nil {
		return nil, err
	}
	err = ev.Decode(vc.accs)
	if err != nil {
		return nil, err
	}
	err = ev.Settle()
	if err != nil {
		return nil, err
	}

	return vc, nil
}

func (vc *VesClient) setName(b []byte) {
	vc.rwMutex.Lock()
	defer vc.rwMutex.Unlock()
	vc.name = make([]byte, len(b))
	copy(vc.name, b)
}

func (vc *VesClient) getName() []byte {
	vc.rwMutex.RLock()
	defer vc.rwMutex.RUnlock()
	return vc.name
}

func (vc *VesClient) updateFileObj(name string, obj interface{}) error {
	ev, err := vc.fdb.WriteWithPath(name)
	if err != nil {
		return err
	}
	err = ev.Encode(obj)
	if err != nil {
		err2 := ev.Settle()
		return errors.New(err.Error() + "\n" + err2.Error())
	}
	err = ev.Settle()
	if err != nil {
		return err
	}
	return nil
}

func (vc *VesClient) updateKeys() error {
	return vc.updateFileObj("keys", vc.keys)
}

func (vc *VesClient) updateAccs() error {
	return vc.updateFileObj("accs", vc.accs)
}

func (vc *VesClient) getClientHello() *wsrpc.ClientHello {
	if vc.clientHello == nil {
		vc.clientHello = new(wsrpc.ClientHello)
	}
	return vc.clientHello
}

func (vc *VesClient) getClientHelloReply() *wsrpc.ClientHelloReply {
	if vc.clientHelloReply == nil {
		vc.clientHelloReply = new(wsrpc.ClientHelloReply)
	}
	return vc.clientHelloReply
}

func (vc *VesClient) getShortSendMessage() *wsrpc.Message {
	if vc.shortSendMessage == nil {
		vc.shortSendMessage = new(wsrpc.Message)
	}
	return vc.shortSendMessage
}

func (vc *VesClient) getShortReplyMessage() *wsrpc.Message {
	if vc.shortReplyMessage == nil {
		vc.shortReplyMessage = new(wsrpc.Message)
	}
	return vc.shortReplyMessage
}

func (vc *VesClient) getUserRegisterRequest() *wsrpc.UserRegisterRequest {
	if vc.userRegisterRequest == nil {
		vc.userRegisterRequest = new(wsrpc.UserRegisterRequest)
	}
	return vc.userRegisterRequest
}

func (vc *VesClient) getUserRegisterReply() *wsrpc.UserRegisterReply {
	if vc.userRegisterReply == nil {
		vc.userRegisterReply = new(wsrpc.UserRegisterReply)
	}
	return vc.userRegisterReply
}

func (vc *VesClient) getrequestComingRequest() *wsrpc.RequestComingRequest {
	if vc.requestComingRequest == nil {
		vc.requestComingRequest = new(wsrpc.RequestComingRequest)
	}
	return vc.requestComingRequest
}

func (vc *VesClient) getrequestComingReply() *wsrpc.RequestComingReply {
	if vc.requestComingReply == nil {
		vc.requestComingReply = new(wsrpc.RequestComingReply)
	}
	return vc.requestComingReply
}

func (vc *VesClient) getrequestGrpcServiceRequest() *wsrpc.RequestGrpcServiceRequest {
	if vc.requestGrpcServiceRequest == nil {
		vc.requestGrpcServiceRequest = new(wsrpc.RequestGrpcServiceRequest)
	}
	return vc.requestGrpcServiceRequest
}

func (vc *VesClient) getrequestGrpcServiceReply() *wsrpc.RequestGrpcServiceReply {
	if vc.requestGrpcServiceReply == nil {
		vc.requestGrpcServiceReply = new(wsrpc.RequestGrpcServiceReply)
	}
	return vc.requestGrpcServiceReply
}

func (vc *VesClient) getrequestNsbServiceRequest() *wsrpc.RequestNsbServiceRequest {
	if vc.requestNsbServiceRequest == nil {
		vc.requestNsbServiceRequest = new(wsrpc.RequestNsbServiceRequest)
	}
	return vc.requestNsbServiceRequest
}

func (vc *VesClient) getrequestNsbServiceReply() *wsrpc.RequestNsbServiceReply {
	if vc.requestNsbServiceReply == nil {
		vc.requestNsbServiceReply = new(wsrpc.RequestNsbServiceReply)
	}
	return vc.requestNsbServiceReply
}

func (vc *VesClient) getuserRegisterRequest() *wsrpc.UserRegisterRequest {
	if vc.userRegisterRequest == nil {
		vc.userRegisterRequest = new(wsrpc.UserRegisterRequest)
	}
	return vc.userRegisterRequest
}

func (vc *VesClient) getuserRegisterReply() *wsrpc.UserRegisterReply {
	if vc.userRegisterReply == nil {
		vc.userRegisterReply = new(wsrpc.UserRegisterReply)
	}
	return vc.userRegisterReply
}

func (vc *VesClient) getsessionListRequest() *wsrpc.SessionListRequest {
	if vc.sessionListRequest == nil {
		vc.sessionListRequest = new(wsrpc.SessionListRequest)
	}
	return vc.sessionListRequest
}

func (vc *VesClient) getsessionListReply() *wsrpc.SessionListReply {
	if vc.sessionListReply == nil {
		vc.sessionListReply = new(wsrpc.SessionListReply)
	}
	return vc.sessionListReply
}

func (vc *VesClient) gettransactionListRequest() *wsrpc.TransactionListRequest {
	if vc.transactionListRequest == nil {
		vc.transactionListRequest = new(wsrpc.TransactionListRequest)
	}
	return vc.transactionListRequest
}

func (vc *VesClient) gettransactionListReply() *wsrpc.TransactionListReply {
	if vc.transactionListReply == nil {
		vc.transactionListReply = new(wsrpc.TransactionListReply)
	}
	return vc.transactionListReply
}

func (vc *VesClient) getSessionStart() *uiprpc.SessionStartRequest {
	if vc.sessionStart == nil {
		vc.sessionStart = new(uiprpc.SessionStartRequest)
	}
	return vc.sessionStart
}

func (vc *VesClient) postMessage(code wsrpc.MessageType, msg proto.Message) error {
	buf, err := wsrpc.GetDefaultSerializer().Serial(code, msg)
	if err != nil {
		fmt.Println(err)
		return err
	}
	vc.conn.WriteMessage(websocket.BinaryMessage, buf.Bytes())
	wsrpc.GetDefaultSerializer().Put(buf)
	return nil
}

func (vc *VesClient) write() {
	var (
		reader                                         = bufio.NewReader(os.Stdin)
		cmdBytes, toBytes, filePath, alias, fileBuffer []byte
		buf                                            *bytes.Buffer
	)
	fileBuffer = make([]byte, 65536)
	for {
		strBytes, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println(err)
			return
		}

		buf = bytes.NewBuffer(bytes.TrimSpace(strBytes))

		cmdBytes, err = buf.ReadBytes(' ')
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return
		}

		switch string(bytes.TrimSpace(cmdBytes)) {
		case "set-name":
			vc.name, err = buf.ReadBytes(' ')
			if err != nil && err != io.EOF {
				fmt.Println(err)
				continue
			}
			vc.name = bytes.TrimSpace(vc.name)
			if err = vc.sayClientHello(vc.name); err != nil {
				fmt.Println(err)
				continue
			}

		case "send-to":
			toBytes, err = buf.ReadBytes(' ')
			if err != nil && err != io.EOF {
				fmt.Println(err)
				continue
			}

			if err = vc.sendMessage(
				bytes.TrimSpace(toBytes),
				bytes.TrimSpace(buf.Bytes()),
			); err != nil {
				fmt.Println(err)
				continue
			}
		case "register-key":
			filePath, err = buf.ReadBytes(' ')
			if err != nil && err != io.EOF {
				fmt.Println(err)
				continue
			}

			if err = vc.registerKey(
				bytes.TrimSpace(filePath),
				fileBuffer,
			); err != nil {
				fmt.Println(err)
				continue
			}
		case "register-eth":
			filePath, err = buf.ReadBytes(' ')
			if err != nil && err != io.EOF {
				fmt.Println(err)
				continue
			}

			if err = vc.configEth(
				bytes.TrimSpace(filePath),
				fileBuffer,
			); err != nil {
				fmt.Println(err)
				continue
			}
		case "send-eth-alias-to-ves":
			alias, err = buf.ReadBytes(' ')
			if err != nil && err != io.EOF {
				fmt.Println(err)
				continue
			}

			if err = vc.sendEthAlias(
				bytes.TrimSpace(alias),
			); err != nil {
				fmt.Println(err)
				continue
			}
		case "send-alias-to-ves":
			alias, err = buf.ReadBytes(' ')
			if err != nil && err != io.EOF {
				fmt.Println(err)
				continue
			}
			if err = vc.sendAlias(
				bytes.TrimSpace(alias),
			); err != nil {
				fmt.Println(err)
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
	if err != nil {
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
	if err != nil {
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
		userRegister.Account = &wsrpc.Account{Address: b, ChainId: acc.ChainID}
		userRegister.UserName = *(*string)(unsafe.Pointer(&vc.name))
		err := vc.postMessage(wsrpc.CodeUserRegisterRequest, userRegister)
		if err != nil {
			return err
		}
		for {
			select {
			case msgBuf := <-vc.cb:
				var messageID uint16
				binary.Read(msgBuf, binary.BigEndian, &messageID)
				if messageID != wsrpc.CodeUserRegisterReply {
					continue
				}
				var s = vc.getUserRegisterReply()
				err = proto.Unmarshal(msgBuf.Bytes(), s)
				if err != nil {
					// ignoring
					// todo: add hidden log
					continue
				}
				//todo: checkCharacteristicFlag
				if !s.GetOk() {
					return errors.New("register user failed")
				}
				return nil
			case <-time.After(time.Second * 5):
				return errors.New("timeout")
			}
		}
	}
	return errors.New("not found")
}

func (vc *VesClient) sendAlias(alias []byte) error {
	if key, ok := vc.keys.Alias[*(*string)(unsafe.Pointer(&alias))]; ok {
		userRegister := vc.getUserRegisterRequest()

		signer := signaturer.NewTendermintNSBSigner(key.PrivateKey)

		userRegister.Account = &wsrpc.Account{Address: signer.GetPublicKey(), ChainId: key.ChainID}
		userRegister.UserName = *(*string)(unsafe.Pointer(&vc.name))

		return vc.postMessage(wsrpc.CodeUserRegisterRequest, userRegister)
	}
	return errors.New("not found")
}

func (vc *VesClient) sayClientHello(name []byte) error {
	clientHello := vc.getClientHello()
	clientHello.Name = name

	return vc.postMessage(wsrpc.CodeClientHelloRequest, clientHello)
}

func (vc *VesClient) sendMessage(to, msg []byte) error {
	shortSendMessage := vc.getShortSendMessage()
	shortSendMessage.From = vc.name
	shortSendMessage.To = to
	shortSendMessage.Contents = string(msg)

	return vc.postMessage(wsrpc.CodeMessageRequest, shortSendMessage)
}

func (vc *VesClient) read() {
	for {
		_, message, err := vc.conn.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}

		var buf = bytes.NewBuffer(message)
		var messageID uint16
		binary.Read(buf, binary.BigEndian, &messageID)
		switch messageID {
		case wsrpc.CodeMessageReply:

			var s = vc.getShortReplyMessage()
			err = proto.Unmarshal(buf.Bytes(), s)
			if err != nil {
				// ignoring
				// todo: add hidden log
				continue
			}

			if bytes.Equal(s.To, vc.getName()) {
				fmt.Printf("%v is saying: %v\n", string(s.From), s.Contents)
			}
		case wsrpc.CodeClientHelloReply:
			var s = vc.getClientHelloReply()
			err = proto.Unmarshal(buf.Bytes(), s)
			if err != nil {
				// ignoring
				// todo: add hidden log
				continue
			}

			vc.grpcip, err = decodeIP(s.GetGrpcHost())
			if err != nil {
				// ignoring
				// todo: add hidden log
			}
			vc.nsbip, err = decodeIP(s.GetNsbHost())
			if err != nil {
				// ignoring
				// todo: add hidden log
			}
		case wsrpc.CodeRequestComingRequest:
			var s = vc.getrequestComingRequest()
			err = proto.Unmarshal(buf.Bytes(), s)
			if err != nil {
				// ignoring
				// todo: add hidden log
				continue
			}

		case wsrpc.CodeUserRegisterReply:
			buf.Reset()
			vc.cb <- buf

		default:
			// fmt.Println("aborting message", string(message))
			// abort
		}

	}
}

func decodeIP(ip []byte) (string, error) {
	if len(ip) == 6 {
		return fmt.Sprintf("%v.%v.%v.%v:%v", ip[0], ip[1], ip[2], ip[3], (uint16(ip[4])<<8)|uint16(ip[5])), nil
	} else if len(ip) == 18 {
		return fmt.Sprintf("[%v]:%v", net.IP(ip[0:16]), (uint16(ip[16])<<8)|uint16(ip[17])), nil
	} else {
		return "", errors.New("invalid length")
	}
}

func main() {
	var (
		dialer        *websocket.Dialer
		addr          = flag.String("addr", "localhost:23452", "http service address")
		u             = url.URL{Scheme: "vc", Host: *addr, Path: "/"}
		vcClient, err = NewVesClient("./data")
	)
	if err != nil {
		log.Println(err)
		return
	}

	vcClient.conn, _, err = dialer.Dial(u.String(), nil)
	if err != nil {
		log.Println(err)
		return
	}

	go vcClient.write()
	go vcClient.read()
}
