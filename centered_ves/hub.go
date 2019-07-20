// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"unsafe"
)

const (
	localChain       = uint64((127 << 24) + 1)
	placeHolderChain = uint64((127 << 24) + 2)
)

type uniMessage struct {
	chainID uint64
	aim     []byte
	message []byte
}

type clientKey struct {
	chainID uint64
	address string
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Registered clients.
	reverseClients map[clientKey]*Client

	// Inbound messages from the clients.
	broadcast chan []byte

	// messages to single clients
	unicast chan *uniMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	server *Server
}

func newHub() *Hub {
	return &Hub{
		broadcast:      make(chan []byte),
		unicast:        make(chan *uniMessage),
		reverseClients: make(map[clientKey]*Client),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		clients:        make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			for _, address := range client.user.GetAccounts() {
				var a = address.GetAddress()
				h.reverseClients[clientKey{
					address.GetChainId(),
					*(*string)(unsafe.Pointer(&a)),
				}] = client
			}
			var a = client.user.GetName()
			h.reverseClients[clientKey{
				placeHolderChain,
				*(*string)(unsafe.Pointer(&a)),
			}] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				for _, address := range client.user.GetAccounts() {
					var a = address.GetAddress()
					delete(h.reverseClients, clientKey{
						address.GetChainId(),
						*(*string)(unsafe.Pointer(&a)),
					})
				}
				var a = client.user.GetName()
				delete(h.reverseClients, clientKey{
					placeHolderChain,
					*(*string)(unsafe.Pointer(&a)),
				})
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
					for _, address := range client.user.GetAccounts() {
						var a = address.GetAddress()
						delete(h.reverseClients, clientKey{
							address.GetChainId(),
							*(*string)(unsafe.Pointer(&a)),
						})
					}
					var a = client.user.GetName()
					delete(h.reverseClients, clientKey{
						placeHolderChain,
						*(*string)(unsafe.Pointer(&a)),
					})
				}
			}
		case message := <-h.unicast:
			if client, ok := h.reverseClients[clientKey{
				message.chainID,
				*(*string)(unsafe.Pointer(&message.aim)),
			}]; ok {
				select {
				case client.send <- message.message:
				default:
					close(client.send)
					delete(h.clients, client)
					for _, address := range client.user.GetAccounts() {
						var a = address.GetAddress()
						delete(h.reverseClients, clientKey{
							address.GetChainId(),
							*(*string)(unsafe.Pointer(&a)),
						})
					}
					var a = client.user.GetName()
					delete(h.reverseClients, clientKey{
						placeHolderChain,
						*(*string)(unsafe.Pointer(&a)),
					})
				}
			} else {
				log.Println("debugging unknown aim", string(message.aim), message.aim)
			}

		}
	}
}
