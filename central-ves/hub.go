// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package centered_ves

import (
	"crypto/md5"
	"encoding/hex"
	"time"
	"unsafe"

	"github.com/gorilla/websocket"
)

const (
	localChain       = uint64((127 << 24) + 1)
	placeHolderChain = uint64((127 << 24) + 2)
)

type uniMessage struct {
	chainID uint64
	aim     []byte
	message []byte
	cb      func()
}

type broMessage struct {
	message []byte
	cb      func()
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
	broadcast chan *broMessage

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
		broadcast:      make(chan *broMessage),
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
			select {
			case <-client.helloed:
				// do nothing
			case <-time.After(5 * time.Second):
				message := websocket.FormatCloseMessage(
					websocket.ClosePolicyViolation,
					"client hello please",
				)
				err := client.conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(2))
				if err != nil {
					h.server.logger.Error("write close message error", "address", client.conn.RemoteAddr())
				}
				err = client.conn.Close()
				if err != nil {
					h.server.logger.Error("close error", "address", client.conn.RemoteAddr())
				}
				return
			}

			h.clients[client] = true
			for _, address := range client.user.GetAccounts() {
				var a = address.GetAddress()
				h.reverseClients[clientKey{
					address.GetChainId(),
					*(*string)(unsafe.Pointer(&a)),
				}] = client
				// fmt.Println("maps", address.GetChainId(), hex.EncodeToString(address.GetAddress()), "->", client.user.GetName())
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
			tag := md5.Sum(message.message)
			h.server.logger.Info("message broadcasting", "tag", hex.EncodeToString(tag[:]))
			for client := range h.clients {
				// fmt.Println("msg...", client.user, message)
				select {
				case client.send <- &writeMessageTask{message.message, message.cb}:
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
			message.cb()
		case message := <-h.unicast:
			tag := md5.Sum(message.message)
			if client, ok := h.reverseClients[clientKey{
				message.chainID,
				*(*string)(unsafe.Pointer(&message.aim)),
			}]; ok {
				h.server.logger.Info("message unicasting", "tag", hex.EncodeToString(tag[:]), "chain id", message.chainID, "address", hex.EncodeToString(message.aim))
				select {
				case client.send <- &writeMessageTask{message.message, message.cb}:
				default:
					h.server.logger.Info("remove no response client", "tag", hex.EncodeToString(tag[:]), "chain id", message.chainID, "address", hex.EncodeToString(message.aim))
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
				h.server.logger.Info("debugging unknown aim", "tag", hex.EncodeToString(tag[:]), "chain id", message.chainID, "address", hex.EncodeToString(message.aim))
			}
		}
	}
}
