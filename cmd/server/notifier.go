package main

import (
	"encoding/json"
	"log"

	"github.com/aalbacetef/txnotify/ethereum"
	"github.com/gorilla/websocket"
)

type WebsocketNotifier struct {
	server *Server
}

func (n *WebsocketNotifier) Notify(address string, txList []ethereum.Transaction) {
	n.server.mu.Lock()
	defer n.server.mu.Unlock()

	notification := Notification{Address: address, Txs: txList}
	data, err := json.Marshal(notification)
	if err != nil {
		log.Printf("marshal error: %v", err)
		return
	}

	for conn, subs := range n.server.conns {
		if _, ok := subs[address]; ok {
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("write error: %v", err)
			}
		}
	}
}
