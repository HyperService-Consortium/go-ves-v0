// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package centered_ves

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/Myriad-Dreamin/minimum-lib/logger"
	"net"
	"net/http"
	"strings"

	grpc "google.golang.org/grpc"
	reflection "google.golang.org/grpc/reflection"

	log "github.com/HyperService-Consortium/go-ves/lib/log"

	uiptypes "github.com/HyperService-Consortium/go-uip/uiptypes"
	types "github.com/HyperService-Consortium/go-ves/types"

	uiprpc "github.com/HyperService-Consortium/go-ves/grpc/uiprpc"
	uipbase "github.com/HyperService-Consortium/go-ves/grpc/uiprpc-base"
	wsrpc "github.com/HyperService-Consortium/go-ves/grpc/wsrpc"
)

// func serveHome(w http.ResponseWriter, r *http.Request) {
// 	log.Println(r.URL)
// 	if r.URL.Path != "/" {
// 		http.Error(w, "Not found", http.StatusNotFound)
// 		return
// 	}
// 	if r.Method != "GET" {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	http.ServeFile(w, r, "home.html")
// }

// Server is a client manager, named centered ves
// it is not in the standard of uip
type Server struct {
	logger logger.Logger
	*http.Server
	hub     *Hub
	vesdb   types.VESDB
	rpcport string
	nsbip []byte
}

type NSBHostOption string

type ServerOptions struct {
	logger logger.Logger
	nsbHost NSBHostOption
}

func defaultServerOptions() ServerOptions {
	return ServerOptions{
		logger: logger.NewStdLogger(),
		nsbHost: "127.0.0.1:26657",
	}
}

func parseOptions(rOptions []interface{}) ServerOptions {
	var options = defaultServerOptions()
	for i := range rOptions {
		switch option := rOptions[i].(type) {
		case logger.Logger:
			options.logger = option
		case NSBHostOption:
			options.nsbHost = option
		}
	}
	return options
}


// NewServer return a pointer of Server
func NewServer(rpcport, addr string, db types.VESDB, rOptions ...interface{}) (srv *Server, err error) {
	options := parseOptions(rOptions)
	srv = &Server{Server: new(http.Server)}
	srv.nsbip, err = HostFromString(options.nsbHost)
	srv.hub = newHub()
	srv.hub.server = srv
	srv.vesdb = db
	srv.Handler = http.NewServeMux()
	srv.Addr = addr
	srv.rpcport = rpcport
	srv.logger = options.logger
	return
}

func HostFromString(option NSBHostOption) ([]byte, error) {
	r := strings.TrimPrefix(string(strings.TrimPrefix(string(option), "https://")), "http://")
	addr, err := net.ResolveTCPAddr("", r)
	if err != nil {
		return nil, err
	}
	return append(addr.IP.To4(), byte(addr.Port >> 8), byte(addr.Port & 0xff)), nil
}

func (srv *Server) ListenAndServeRpc(port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Println(fmt.Errorf("failed to listen: %v", err))
	}
	s := grpc.NewServer()
	uiprpc.RegisterCenteredVESServer(s, srv)
	reflection.Register(s)

	fmt.Printf("prepare to serve rpc on %v\n", port)
	if err := s.Serve(lis); err != nil {
		log.Println(fmt.Errorf("failed to serve: %v", err))
	}
	return
}

// Start the service of centered ves
func (srv *Server) Start() error {
	go srv.hub.run()
	go srv.ListenAndServeRpc(srv.rpcport)
	srv.Handler.(*http.ServeMux).HandleFunc("/", srv.serveWs)
	return srv.ListenAndServe()
}

// serveWs handles websocket requests from the peer.
func (srv *Server) serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("new ws: %v\n", r.RemoteAddr)
	client := &Client{hub: srv.hub, helloed: make(chan bool, 1), conn: conn, send: make(chan *writeMessageTask, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

func (srv *Server) InternalRequestComing(
	ctx context.Context,
	in *uiprpc.InternalRequestComingRequest,
) (*uiprpc.InternalRequestComingReply, error) {
	if err := srv.RequestComing(func() (accs []uiptypes.Account) {
		for _, acc := range in.GetAccounts() {
			accs = append(accs, acc)
		}
		return accs
	}(), in.GetSessionId(), in.GetHost()); err != nil {
		return nil, err
	}
	return &uiprpc.InternalRequestComingReply{
		Ok: true,
	}, nil
}

func (srv *Server) InternalAttestationSending(
	ctx context.Context,
	in *uiprpc.InternalRequestComingRequest,
) (*uiprpc.InternalRequestComingReply, error) {
	if err := srv.AttestationSending(func() (accs []uiptypes.Account) {
		for _, acc := range in.GetAccounts() {
			accs = append(accs, acc)
		}
		return accs
	}(), in.GetSessionId(), in.GetHost()); err != nil {
		return nil, err
	}
	return &uiprpc.InternalRequestComingReply{
		Ok: true,
	}, nil
}

// RequestComing do the service of retransmitting message of new session event
func (srv *Server) RequestComing(accounts []uiptypes.Account, iscAddress, grpcHost []byte) (err error) {
	// fmt.Println("rpc...", accounts)
	for _, acc := range accounts {
		// fmt.Println("hex", acc.GetChainId(), hex.EncodeToString(acc.GetAddress()))
		log.Println("sending session request", acc.GetChainId(), hex.EncodeToString(acc.GetAddress()))
		if err = srv.requestComing(acc, iscAddress, grpcHost); err != nil {
			return
		}
	}
	return nil
}

// AttestationSending do the service of retransmitting attestation
func (srv *Server) AttestationSending(accounts []uiptypes.Account, iscAddress, grpcHost []byte) (err error) {
	// fmt.Println("rpc...", accounts)
	for _, acc := range accounts {
		log.Println("sending attestation request", acc.GetChainId(), hex.EncodeToString(acc.GetAddress()))
		if err = srv.attestationSending(acc, iscAddress, grpcHost); err != nil {
			return
		}
	}
	return nil
}

func (srv *Server) requestComing(acc uiptypes.Account, iscAddress, grpcHost []byte) error {
	var msg wsrpc.RequestComingRequest
	msg.NsbHost = srv.nsbip
	msg.GrpcHost = grpcHost
	msg.SessionId = iscAddress
	msg.Account = &uipbase.Account{
		Address: acc.GetAddress(),
		ChainId: acc.GetChainId(),
	}
	packet, err := wsrpc.GetDefaultSerializer().Serial(wsrpc.CodeRequestComingRequest, &msg)
	if err != nil {
		return err
	}
	srv.hub.unicast <- &uniMessage{acc.GetChainId(), acc.GetAddress(), packet.Bytes(), func() {
		wsrpc.GetDefaultSerializer().Put(packet)
	}}
	return nil
}

func (srv *Server) attestationSending(acc uiptypes.Account, iscAddress, grpcHost []byte) error {
	var msg wsrpc.RequestComingRequest
	msg.NsbHost = srv.nsbip
	msg.GrpcHost = grpcHost
	msg.SessionId = iscAddress
	msg.Account = &uipbase.Account{
		Address: acc.GetAddress(),
		ChainId: acc.GetChainId(),
	}

	// log.Infof("attestating network gate", )

	packet, err := wsrpc.GetDefaultSerializer().Serial(wsrpc.CodeAttestationSendingRequest, &msg)
	if err != nil {
		return err
	}
	srv.hub.unicast <- &uniMessage{acc.GetChainId(), acc.GetAddress(), packet.Bytes(), func() {
		wsrpc.GetDefaultSerializer().Put(packet)
	}}
	return nil
}

func (srv *Server) InternalCloseSession(
	ctx context.Context,
	in *uiprpc.InternalCloseSessionRequest,
) (*uiprpc.InternalCloseSessionReply, error) {
	if err := srv.CloseSession(func() (accs []uiptypes.Account) {
		for _, acc := range in.GetAccounts() {
			accs = append(accs, acc)
		}
		return accs
	}(), in.GetSessionId(), in.GetGrpcHost(), in.GetNsbHost()); err != nil {
		return nil, err
	}
	return &uiprpc.InternalCloseSessionReply{
		Ok: true,
	}, nil
}

// CloseSession do the service of retransmitting attestation
func (srv *Server) CloseSession(accounts []uiptypes.Account, iscAddress, grpcHost, nsbHost []byte) (err error) {
	// fmt.Println("rpc...", accounts)
	for _, acc := range accounts {
		log.Println("sending close session", acc.GetChainId(), hex.EncodeToString(acc.GetAddress()))
		if err = srv.closeSession(acc, iscAddress, grpcHost, nsbHost); err != nil {
			return
		}
	}
	return nil
}

func (srv *Server) closeSession(acc uiptypes.Account, iscAddress, grpcHost, nsbHost []byte) error {
	var msg wsrpc.CloseSessionRequest
	msg.NsbHost = nsbHost
	msg.GrpcHost = grpcHost
	msg.SessionId = iscAddress

	// log.Infof("attestating network gate", )

	packet, err := wsrpc.GetDefaultSerializer().Serial(wsrpc.CodeCloseSessionRequest, &msg)
	if err != nil {
		return err
	}
	srv.hub.unicast <- &uniMessage{acc.GetChainId(), acc.GetAddress(), packet.Bytes(), func() {
		wsrpc.GetDefaultSerializer().Put(packet)
	}}
	return nil
}
