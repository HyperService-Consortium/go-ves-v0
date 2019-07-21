// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package centered_ves

import (
	"context"
	"fmt"
	"net"
	"net/http"

	uiptypes "github.com/Myriad-Dreamin/go-uip/types"
	uiprpc "github.com/Myriad-Dreamin/go-ves/grpc/uip-rpc"
	wsrpc "github.com/Myriad-Dreamin/go-ves/grpc/ws-ves-rpc"
	log "github.com/Myriad-Dreamin/go-ves/log"
	"github.com/Myriad-Dreamin/go-ves/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	*http.Server
	hub     *Hub
	vesdb   types.VESDB
	rpcport string
}

// NewServer return a pointer of Server
func NewServer(rpcport, addr string, db types.VESDB) (srv *Server) {
	srv = &Server{Server: new(http.Server)}
	srv.hub = newHub()
	srv.hub.server = srv
	srv.vesdb = db
	srv.Handler = http.NewServeMux()
	srv.Addr = addr
	srv.rpcport = rpcport
	return
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
	client := &Client{hub: srv.hub, helloed: make(chan bool, 1), conn: conn, send: make(chan []byte, 256)}
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

// RequestComing do the service of retransmitting message of new session event
func (srv *Server) RequestComing(accounts []uiptypes.Account, iscAddress, grpcHost []byte) (err error) {
	for _, acc := range accounts {
		if err = srv.requestComing(acc.GetChainId(), acc.GetAddress(), iscAddress, grpcHost); err != nil {
			return
		}
	}
	return nil
}

func (srv *Server) requestComing(chainID uint64, address, iscAddress, grpcHost []byte) error {
	var msg wsrpc.RequestComingRequest
	msg.NsbHost = nsbip
	msg.GrpcHost = grpcHost
	msg.SessionId = iscAddress
	qwq, err := wsrpc.GetDefaultSerializer().Serial(wsrpc.CodeRequestComingRequest, &msg)
	if err != nil {
		return err
	}
	srv.hub.unicast <- &uniMessage{chainID, address, qwq.Bytes()}
	wsrpc.GetDefaultSerializer().Put(qwq)
	return nil
}
