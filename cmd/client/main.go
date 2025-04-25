package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/aalbacetef/txnotify/ethereum"
)

type SubscriptionRequest struct {
	Address string `json:"address"`
}

type Notification struct {
	Address string                 `json:"address"`
	Txs     []ethereum.Transaction `json:"transactions"` //nolint:tagliatelle
}

func main() {
	serverAddr := "ws://localhost:8080/ws"
	addresses := ""
	timeout := "5m"

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
		log.Printf("error parsing timeout: %v", err)
		return
	}

	if err := run(serverAddr, splitAddresses(addresses), timeoutDuration); err != nil {
		log.Printf("error: %v", err)
		return
	}
}

func run(serverAddr string, addresses []string, timeoutDuration time.Duration) error { //nolint:gocognit
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, serverAddr, nil)
	if err != nil {
		return fmt.Errorf("dial error: %w", err)
	}
	defer conn.Close()

	var mu sync.Mutex

	seenTxs := subscribe(conn, addresses)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			_, msg, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("read error: %v", err)
				}

				return fmt.Errorf("conn.ReadMessage: %w", err)
			}

			var notif Notification
			if err := json.Unmarshal(msg, &notif); err != nil {
				log.Printf("unmarshal error: %v", err)
				continue
			}

			// this had to be added to drop duplicate server notifications.
			mu.Lock()
			if _, exists := seenTxs[notif.Address]; !exists {
				seenTxs[notif.Address] = make(map[string]struct{})
			}

			for _, tx := range notif.Txs {
				if _, seen := seenTxs[notif.Address][tx.Hash]; !seen {
					fmt.Printf("%s) got tx: %s\n", notif.Address, tx.Hash)
					seenTxs[notif.Address][tx.Hash] = struct{}{}
				}
			}
			mu.Unlock()
		}
	}
}

func subscribe(conn *websocket.Conn, addresses []string) map[string]map[string]struct{} {
	seenTxs := make(map[string]map[string]struct{})

	for _, addr := range addresses {
		seenTxs[addr] = make(map[string]struct{})

		req := SubscriptionRequest{Address: addr}
		buf := &bytes.Buffer{}

		if err := json.NewEncoder(buf).Encode(req); err != nil {
			log.Printf("marshal error for address %s: %v", addr, err)
			continue
		}

		if err := conn.WriteMessage(websocket.TextMessage, buf.Bytes()); err != nil {
			log.Printf("write error for address %s: %v", addr, err)
			continue
		}
	}

	return seenTxs
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
