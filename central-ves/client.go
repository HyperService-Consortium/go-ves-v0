// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package centered_ves

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"

	log "github.com/HyperService-Consortium/go-ves/lib/log"

	types "github.com/HyperService-Consortium/go-ves/types"

	uipbase "github.com/HyperService-Consortium/go-ves/grpc/uiprpc-base"
	wsrpc "github.com/HyperService-Consortium/go-ves/grpc/wsrpc"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 4 * 1024
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}

	// nsb ip
	nsbip = []byte{47, 251, 2, 73, uint8(26657 >> 8), uint8(26657 & 0xff)}

	// grpc ips
	grpcips = [][]byte{
		[]byte{127, 0, 0, 1, uint8(23351 >> 8), uint8(23351 & 0xff)},
	}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 1024,
	WriteBufferSize: 1024 * 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type writeMessageTask struct {
	b  []byte
	cb func()
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// owned user
	user types.User

	// Buffered channel of outbound messages.
	send chan *writeMessageTask

	// client hello sended
	helloed chan bool
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		tag := md5.Sum(message)
		c.hub.server.logger.Info("reading message", "tag", hex.EncodeToString(tag[:]))

		var buf = bytes.NewBuffer(message)
		var messageID uint16
		binary.Read(buf, binary.BigEndian, &messageID)
		switch wsrpc.MessageType(messageID) {
		case wsrpc.CodeMessageRequest:
			var s wsrpc.Message
			err = proto.Unmarshal(buf.Bytes(), &s)
			if err != nil {
				log.Println("err:", err)
				continue
			}
			log.Infoln("message request", s.GetContents(), string(s.GetFrom()), s.GetFrom(), "->", string(s.GetTo()), s.GetTo())
			// fmt.Println(s.GetContents(), string(s.GetFrom()), s.GetFrom(), "->", string(s.GetTo()), s.GetTo())
			var qwq, err = wsrpc.GetDefaultSerializer().Serial(wsrpc.CodeMessageReply, &s)

			if err != nil {
				log.Println("err:", qwq)
				continue
			}
			c.hub.broadcast <- &broMessage{qwq.Bytes(), func() {
				wsrpc.GetDefaultSerializer().Put(qwq)
			}}
		case wsrpc.CodeRawProto:

			var s wsrpc.RawMessage
			err = proto.Unmarshal(buf.Bytes(), &s)
			if err != nil {
				log.Println("err:", err)
				continue
			}
			var a uipbase.Account
			err = proto.Unmarshal(s.GetTo(), &a)
			if err != nil {
				log.Println("err:", err)
				continue
			}
			tag := md5.Sum(s.GetContents())
			log.Infoln("raw proto", hex.EncodeToString(tag[:]), string(s.GetFrom()), s.GetFrom(), "->", string(s.GetTo()), s.GetTo())
			// fmt.Println(string(s.GetContents()), string(s.GetFrom()), s.GetFrom(), "->", string(a.GetAddress()), s.GetTo())

			c.hub.unicast <- &uniMessage{a.GetChainId(), a.GetAddress(), s.GetContents(), func() {}}
		case wsrpc.CodeClientHelloRequest:
			var s wsrpc.ClientHello
			err = proto.Unmarshal(buf.Bytes(), &s)
			if err != nil {
				log.Println("err:", err)
			}

			c.user, err = c.hub.server.vesdb.FindUser(string(s.GetName()))
			// fmt.Println(c.user, err)
			if err != nil {
				log.Println(err)
				return
			}

			var t wsrpc.ClientHelloReply
			t.GrpcHost = grpcips[0]
			t.NsbHost = c.hub.server.nsbip

			qwq, err := wsrpc.GetDefaultSerializer().Serial(wsrpc.CodeClientHelloReply, &t)
			if err != nil {
				log.Println("err:", err)
				continue
			}
			c.helloed <- true
			c.hub.unicast <- &uniMessage{placeHolderChain, s.GetName(), qwq.Bytes(), func() {
				wsrpc.GetDefaultSerializer().Put(qwq)
			}}

		case wsrpc.CodeUserRegisterRequest:
			var s wsrpc.UserRegisterRequest
			err = proto.Unmarshal(buf.Bytes(), &s)
			if err != nil {
				log.Println("err:", err)
			}

			// fmt.Println("hexx registering", hex.EncodeToString(s.GetAccount().GetAddress()))
			err = c.hub.server.vesdb.InsertAccount(s.GetUserName(), s.GetAccount())

			if err != nil {
				log.Println("err:", err)
				continue
			}
		default:
			fmt.Println("aborting message", string(message))
			// abort
		}

		// c.hub.broadcast <- &broMessage{bytes.TrimSpace(bytes.Replace(message, newline, space, -1)), func() {}}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				// if message.cb != nil {
				// 	message.cb()
				// }
				return
			}

			w, err := c.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				if message.cb != nil {
					message.cb()
				}
				return
			}
			w.Write(message.b)

			// Add queued chat messages to the current websocket message.
			// n := len(c.send)
			// for i := 0; i < n; i++ {
			// 	w.Write(newline)
			// 	w.Write(<-c.send)
			// }

			if message.cb != nil {
				message.cb()
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
