// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"
		wsrpc "github.com/Myriad-Dreamin/go-ves/grpc/ws-ves-rpc"
 uiptypes "github.com/Myriad-Dreamin/go-uip/types"
)

var addr = flag.String("port", ":23452", "http service address")

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
	hub *Hub
}

// NewServer return a pointer of Server
func NewServer(addr string) (srv *Server) {
	srv = &Server{Server:new(http.Server)}
	srv.hub = newHub()
	srv.Handler = http.NewServeMux()
	srv.Addr = addr
	return
}

// Start the service of centered ves
func (srv *Server) Start() error {
	go srv.hub.run()
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
	client := &Client{hub: srv.hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
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
	var msg wsrpc.RequestComing
	msg.NsbHost = nsbip
	msg.GrpcHost = grpcHost
	msg.SessionId = iscAddress
	qwq, err := wsrpc.GetDefaultSerializer().Serial(wsrpc.CodeRequestComing, &msg)
	if err != nil {
		return err
	}
	srv.hub.unicast <- &uniMessage{chainID, address, qwq.Bytes()}
	wsrpc.GetDefaultSerializer().Put(qwq)
	return nil
}

func main() {
	flag.Parse()
	if err := NewServer(*addr).Start(); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
