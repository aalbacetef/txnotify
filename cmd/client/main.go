package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aalbacetef/txnotify/ethereum"
	"github.com/gorilla/websocket"
)

type SubscriptionRequest struct {
	Address string `json:"address"`
}

type Notification struct {
	Address string                 `json:"address"`
	Txs     []ethereum.Transaction `json:"transactions"`
}

func main() {
	serverAddr := "ws://localhost:8080/ws"
	addresses := ""
	timeout := "30s"

	flag.StringVar(&serverAddr, "server", serverAddr, "websocket server address")
	flag.StringVar(&addresses, "addresses", addresses, "comma-separated list of addresses to subscribe")
	flag.StringVar(&timeout, "timeout", timeout, "connection timeout")
	flag.Parse()

	if serverAddr == "" || addresses == "" {
		flag.Usage()
		return
	}

	timeoutDuration, err := time.ParseDuration(timeout)
	if err != nil {
		log.Fatalf("error parsing timeout: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, serverAddr, nil)
	if err != nil {
		log.Fatalf("dial error: %v", err)
	}
	defer conn.Close()

	addrList := splitAddresses(addresses)
	for _, addr := range addrList {
		req := SubscriptionRequest{Address: addr}
		data, err := json.Marshal(req)
		if err != nil {
			log.Printf("marshal error for address %s: %v", addr, err)
			continue
		}
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("write error for address %s: %v", addr, err)
			continue
		}
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, msg, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("read error: %v", err)
				}
				return
			}

			var notif Notification
			if err := json.Unmarshal(msg, &notif); err != nil {
				log.Printf("unmarshal error: %v", err)
				continue
			}

			for _, tx := range notif.Txs {
				fmt.Printf("%s) got tx: %s\n", notif.Address, tx.Hash)
			}
		}
	}
}

func splitAddresses(addresses string) []string {
	var result []string
	for _, addr := range strings.Split(addresses, ",") {
		addr = strings.TrimSpace(addr)
		if addr != "" {
			result = append(result, addr)
		}
	}
	return result
}
