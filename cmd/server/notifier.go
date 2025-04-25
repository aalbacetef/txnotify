package main

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"

	"github.com/aalbacetef/txnotify/ethereum"
)

type WebsocketNotifier struct {
	server *Server
}

func (n *WebsocketNotifier) Notify(address string, txList []ethereum.Transaction) {
	n.server.mu.Lock()
	defer n.server.mu.Unlock()

	notification := Notification{Address: address, Txs: txList}

	buf := &bytes.Buffer{}

	if err := json.NewEncoder(buf).Encode(notification); err != nil {
		log.Printf("marshal error: %v", err)
		return
	}

	data := buf.Bytes()

	for conn, subs := range n.server.conns {
		if _, ok := subs[address]; ok {
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("write error: %v", err)
			}
		}
	}
}
