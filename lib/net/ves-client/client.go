package vesclient

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/Myriad-Dreamin/minimum-lib/logger"
	"io"
	"net/url"
	"sync"

	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"

	signaturetype "github.com/HyperService-Consortium/go-uip/const/signature_type"
	"github.com/HyperService-Consortium/go-uip/signaturer"
	uiptypes "github.com/HyperService-Consortium/go-uip/types"

	"github.com/HyperService-Consortium/go-ves/grpc/uiprpc"
	uipbase "github.com/HyperService-Consortium/go-ves/grpc/uiprpc-base"
	"github.com/HyperService-Consortium/go-ves/grpc/wsrpc"

	ethbni "github.com/HyperService-Consortium/go-ves/lib/bni/eth"
	nsbbni "github.com/HyperService-Consortium/go-ves/lib/bni/ten"
	"github.com/HyperService-Consortium/go-ves/lib/database/filedb"
	nsbclient "github.com/HyperService-Consortium/go-ves/lib/net/nsb-client"
)

// VesClient is the web socket client interactive with veses
type VesClient struct {
	rwMutex sync.RWMutex
	logger  logger.Logger

	name   []byte
	signer uiptypes.Signer
	keys   *ECCKeys
	accs   *EthAccounts

	conn      *websocket.Conn
	nsbClient *nsbclient.NSBClient
	waitOpt   *uiptypes.WaitOption

	cb   chan *bytes.Buffer
	quit chan bool

	fdb *filedb.FileDB

	nsbip  string
	grpcip string

	rawMessage                *wsrpc.RawMessage
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
	sessionFinishedRequest    *wsrpc.SessionFinishedRequest
	sessionFinishedReply      *wsrpc.SessionFinishedReply
	// sessionRequireTransactRequest *wsrpc.SessionRequireTransactRequest
	// sessionRequireTransactReply   *wsrpc.SessionRequireTransactReply
	attestationReceiveRequestSend    *wsrpc.AttestationReceiveRequest
	attestationReceiveRequestReceive *wsrpc.AttestationReceiveRequest
	attestationReceiveReply          *wsrpc.AttestationReceiveReply

	sessionStart *uiprpc.SessionStartRequest
}

type CVesHostOption string
type VesName []byte

type ServerOptions struct {
	logger  logger.Logger
	waitOpt *uiptypes.WaitOption
	addr    string
	vesName []byte
}

var globalLogger = logger.NewStdLogger()

func defaultServerOptions() ServerOptions {
	return ServerOptions{
		logger:  globalLogger,
		waitOpt: uiptypes.NewWaitOption(),
		addr:    "127.0.0.1:23452",
	}
}

func parseOptions(rOptions []interface{}) ServerOptions {
	var options = defaultServerOptions()
	for i := range rOptions {
		switch option := rOptions[i].(type) {
		case logger.Logger:
			options.logger = option
		case *uiptypes.WaitOption:
			options.waitOpt = option
		case CVesHostOption:
			options.addr = string(option)
		case VesName:
			options.vesName = option
		}
	}
	return options
}

// NewVesClient return a pointer of VesClinet
func NewVesClient(rOptions ...interface{}) (vc *VesClient, err error) {
	options := parseOptions(rOptions)
	vc = &VesClient{
		cb:        make(chan *bytes.Buffer, 1),
		quit:      make(chan bool, 1),
		nsbClient: nsbclient.NewNSBClient(host),
		logger:    options.logger,
		waitOpt:   options.waitOpt,
		name:      options.vesName,
	}

	vc.conn, _, err = new(websocket.Dialer).Dial((&url.URL{Scheme: "ws", Host: options.addr, Path: "/"}).String(), nil)
	return
}

func (vc *VesClient) Boot() (err error) {
	if err = vc.load(dataPrefix + "/" + string(vc.name)); err != nil {
		vc.logger.Error("load config from filepath error", "path", dataPrefix + "/" + string(vc.name), "error", err)
		return
	}
	phandler.register(vc.save)

	go vc.read()
	if err = vc.SayClientHello(vc.name); err != nil {
		vc.logger.Error("say client hello error", "error", err)
		return
	}
	return
}

func (vc *VesClient) load(dbpath string) error {
	var err, err2 error
	filedb.Register(&ECCKeys{})
	filedb.Register(&EthAccounts{})
	if vc.fdb, err = filedb.NewFileDB(dbpath); err != nil {
		return err
	}
	var ev *filedb.ReadEvent
	ev, err = vc.fdb.ReadWithPath("/keys.dat")
	if err != nil {
		goto bad_load_keys
	}
	vc.keys = new(ECCKeys)
	vc.keys.Alias = make(map[string]ECCKey)
	err = ev.Decode(vc.keys)
	if err != nil {
		goto bad_load_keys
	}
	err = ev.Settle()
	if err != nil {
		goto bad_load_keys
	}
bad_load_keys:
	if err == io.EOF {
		err = nil
	}

	ev, err2 = vc.fdb.ReadWithPath("/accs.dat")
	if err2 != nil {
		goto bad_load_accs
	}
	vc.accs = new(EthAccounts)
	vc.accs.Alias = make(map[string]EthAccount)
	err2 = ev.Decode(&vc.accs)
	if err2 != nil {
		goto bad_load_accs
	}
	err2 = ev.Settle()
	if err2 != nil {
		goto bad_load_accs
	}

	return err
bad_load_accs:
	if err2 == io.EOF {
		err2 = err
	} else {
		if err != nil {
			err2 = fmt.Errorf("error loading keys: %v, error loading accs: %v", err, err2)
		}
	}
	return err2
}

func (vc *VesClient) save() {
	var err error
	if err = vc.updateKeys(); err != nil {
		vc.logger.Error("update keys failed", "error", err)
	}
	if err = vc.updateAccs(); err != nil {
		vc.logger.Error("update accounts failed", "error", err)
	}
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
	return vc.updateFileObj("/keys.dat", vc.keys)
}

func (vc *VesClient) updateAccs() error {
	return vc.updateFileObj("/accs.dat", vc.accs)
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

func (vc *VesClient) getRawMessage() *wsrpc.RawMessage {
	if vc.rawMessage == nil {
		vc.rawMessage = new(wsrpc.RawMessage)
	}
	return vc.rawMessage
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
	// if vc.requestComingRequest == nil {
	// 	vc.requestComingRequest = new(wsrpc.RequestComingRequest)
	// }
	// return vc.requestComingRequest
	return new(wsrpc.RequestComingRequest)
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

func (vc *VesClient) getSessionFinishedRequest() *wsrpc.SessionFinishedRequest {
	if vc.sessionFinishedRequest == nil {
		vc.sessionFinishedRequest = new(wsrpc.SessionFinishedRequest)
	}
	return vc.sessionFinishedRequest
}

func (vc *VesClient) getSessionFinishedReply() *wsrpc.SessionFinishedReply {
	if vc.sessionFinishedReply == nil {
		vc.sessionFinishedReply = new(wsrpc.SessionFinishedReply)
	}
	return vc.sessionFinishedReply
}

// func (vc *VesClient) getSessionRequireTransactRequest() *wsrpc.SessionRequireTransactRequest {
// 	if vc.sessionRequireTransactRequest == nil {
// 		vc.sessionRequireTransactRequest = new(wsrpc.SessionRequireTransactRequest)
// 	}
// 	return vc.sessionRequireTransactRequest
// }
//
// func (vc *VesClient) getSessionRequireTransactReply() *wsrpc.SessionRequireTransactReply {
// 	if vc.sessionRequireTransactReply == nil {
// 		vc.sessionRequireTransactReply = new(wsrpc.SessionRequireTransactReply)
// 	}
// 	return vc.sessionRequireTransactReply
// }

func (vc *VesClient) getSendAttestationReceiveRequest() *wsrpc.AttestationReceiveRequest {
	// if vc.attestationReceiveRequestSend == nil {
	// 	vc.attestationReceiveRequestSend = new(wsrpc.AttestationReceiveRequest)
	// }
	// return vc.attestationReceiveRequestSend
	return new(wsrpc.AttestationReceiveRequest)
}

func (vc *VesClient) getReceiveAttestationReceiveRequest() *wsrpc.AttestationReceiveRequest {
	// if vc.attestationReceiveRequestReceive == nil {
	// 	vc.attestationReceiveRequestReceive = new(wsrpc.AttestationReceiveRequest)
	// }
	// return vc.attestationReceiveRequestReceive
	return new(wsrpc.AttestationReceiveRequest)
}

func (vc *VesClient) getAttestationReceiveReply() *wsrpc.AttestationReceiveReply {
	if vc.attestationReceiveReply == nil {
		vc.attestationReceiveReply = new(wsrpc.AttestationReceiveReply)
	}
	return vc.attestationReceiveReply
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

func (vc *VesClient) postRawMessage(code wsrpc.MessageType, dst *uipbase.Account, msg proto.Message) error {

	buf, err := wsrpc.GetDefaultSerializer().Serial(code, msg)
	/// fmt.Println(buf.Bytes())
	if err != nil {
		fmt.Println(err)
		return err
	}
	var s = vc.getRawMessage()
	s.To, err = proto.Marshal(dst)
	if err != nil {
		fmt.Println(err)
		return err
	}
	s.From = vc.name
	s.Contents = make([]byte, buf.Len())
	copy(s.Contents, buf.Bytes())
	// fmt.Println(s.Contents)
	wsrpc.GetDefaultSerializer().Put(buf)
	return vc.postMessage(wsrpc.CodeRawProto, s)
}

func (vc *VesClient) SendMessage(to, msg []byte) error {
	shortSendMessage := vc.getShortSendMessage()
	shortSendMessage.From = vc.name
	shortSendMessage.To = to
	shortSendMessage.Contents = string(msg)

	// fmt.Println(to, msg)

	return vc.postMessage(wsrpc.CodeMessageRequest, shortSendMessage)
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

func (vc *VesClient) getNSBSigner() (uiptypes.Signer, error) {
	if vc.signer != nil {
		return vc.signer, nil
	}

	if key, ok := vc.keys.Alias["ten1"]; ok {
		vc.signer = signaturer.NewTendermintNSBSigner(key.PrivateKey)
		if vc.signer == nil {
			return nil, errIlegalPrivateKey
		}
	} else {
		return nil, errTen1NotFound
	}
	return vc.signer, nil
}

func (vc *VesClient) getRouter(chainID uint64) uiptypes.Router {
	switch chainID {
	case 1, 2:
		return &ethbni.BN{}
	case 3, 4:
		return &nsbbni.BN{}
	default:
		return nil
	}
}

func (vc *VesClient) getBlockStorage(chainID uint64) uiptypes.Storage {
	switch chainID {
	case 1, 2:
		return &ethbni.BN{}
	case 3, 4:
		return &nsbbni.BN{}
	default:
		return nil
	}
}

func (signer *EthAccount) GetPublicKey() []byte {
	b, _ := hex.DecodeString(signer.Address)
	return b
}

func (signer *EthAccount) Sign(b []byte) uiptypes.Signature {
	// todo: sign b
	return signaturer.FromRaw(b, signaturetype.Secp256k1)
}

func (vc *VesClient) getRespSigner(acc *uipbase.Account) (uiptypes.Signer, error) {
	if vc.signer != nil {
		return vc.signer, nil
	}
	cid := acc.GetChainId()
	switch cid {
	case 1, 2:
		sadd := hex.EncodeToString(acc.GetAddress())
		for _, acc := range vc.accs.Alias {
			if acc.ChainID == cid && acc.Address == sadd {
				return &acc, nil
			}
		}
	case 3, 4:
		for _, key := range vc.keys.Alias {
			if key.ChainID != cid {
				continue
			}

			signer := signaturer.NewTendermintNSBSigner(key.PrivateKey)
			if bytes.Equal(signer.GetPublicKey(), acc.GetAddress()) {
				return signer, nil
			}
		}
	}

	return nil, errNotFound
}

func (vc *VesClient) ListKeys() {
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
}
