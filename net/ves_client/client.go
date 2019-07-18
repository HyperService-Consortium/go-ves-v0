package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"sync"

	wsrpc "github.com/Myriad-Dreamin/go-ves/grpc/ws-ves-rpc"
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
)

var (
	addr = flag.String("addr", "localhost:23452", "http service address")
	u    = url.URL{Scheme: "ws", Host: *addr, Path: "/"}
)

type WSClient struct {
	rwMutex sync.RWMutex
	msg     wsrpc.Message
	conn    *websocket.Conn
}

func main() {
	var err error
	var dialer *websocket.Dialer
	var ws_client = new(WSClient)

	ws_client.conn, _, err = dialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	go ws_client.write()

	for {
		_, message, err := ws_client.conn.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}

		var buf = bytes.NewBuffer(message)
		var message_id uint16
		binary.Read(buf, binary.BigEndian, &message_id)
		switch message_id {
		case wsrpc.CodeMessageReply:

			var s wsrpc.Message
			proto.Unmarshal(buf.Bytes(), &s)

			// fmt.Println(s.GetContents(), string(s.GetFrom()), "->", string(s.GetTo()))

			if bytes.Equal(s.To, ws_client.getName()) {
				fmt.Printf("%v is saying: %v\n", string(s.From), s.Contents)
			}
		default:
			// fmt.Println("aborting message", string(message))
			// abort
		}

	}
}

func (ws *WSClient) setName(b []byte) {
	ws.rwMutex.Lock()
	defer ws.rwMutex.Unlock()
	ws.msg.From = make([]byte, len(b))
	copy(ws.msg.From, b)
}

func (ws *WSClient) getName() []byte {
	ws.rwMutex.RLock()
	defer ws.rwMutex.RUnlock()
	return ws.msg.From
}

func (ws *WSClient) write() {
	reader := bufio.NewReader(os.Stdin)
	var b []byte
	for {
		strBytes, _, err := reader.ReadLine()
		if err != nil {
			log.Fatal(err)
			return
		}

		var buf = bytes.NewBuffer(bytes.TrimSpace(strBytes))
		b, err = buf.ReadBytes(' ')

		if err != nil && err != io.EOF {
			log.Fatal(err)
			return
		}
		switch string(b) {
		case "set-name ":
			ws.setName(buf.Bytes())
			fmt.Println("from =", string(ws.getName()), ws.getName())
			qwq, err := wsrpc.GetDefaultSerializer().Serial(wsrpc.SetNameRequest, &ws.msg)
			if err != nil {
				log.Fatal(err)
				return
			}
			ws.conn.WriteMessage(websocket.BinaryMessage, qwq.Bytes())
			qwq.Reset()
			wsrpc.GetDefaultSerializer().Put(qwq)
		case "send-to ":
			ws.msg.To, err = buf.ReadBytes(' ')
			if err != nil {
				log.Fatal(err)
				return
			}
			ws.msg.To = append(make([]byte, 0, len(ws.msg.To)-1), ws.msg.To[:len(ws.msg.To)-1]...)
			ws.msg.Contents = string(buf.Bytes())
			qwq, err := wsrpc.GetDefaultSerializer().Serial(wsrpc.CodeMessageRequest, &ws.msg)
			if err != nil {
				log.Fatal(err)
				return
			}
			// fmt.Println(qwq.Len())
			ws.conn.WriteMessage(websocket.BinaryMessage, qwq.Bytes())
			// fmt.Println("sending", string(ws.msg.Contents), "from", ws.getName(), string(ws.getName()), "to", ws.msg.To, string(ws.msg.To))
			qwq.Reset()
			wsrpc.GetDefaultSerializer().Put(qwq)
		}

		if err != nil {
			log.Fatal(err)
			return
		}

	}
}
