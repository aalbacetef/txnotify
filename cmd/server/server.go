package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/aalbacetef/txnotify"
	"github.com/aalbacetef/txnotify/ethereum"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Server struct {
	mu           sync.Mutex
	conns        map[*websocket.Conn]map[string]struct{}
	watcher      *txnotify.Watcher
	addr         string
	rpcEndpoint  string
	pollInterval time.Duration
}

type SubscriptionRequest struct {
	Address string `json:"address"`
}

type Notification struct {
	Address string                 `json:"address"`
	Txs     []ethereum.Transaction `json:"transactions"`
}

func NewServer(addr, rpcEndpoint, pollInterval string) (*Server, error) {
	interval, err := time.ParseDuration(pollInterval)
	if err != nil {
		return nil, err
	}
	return &Server{
		conns:        make(map[*websocket.Conn]map[string]struct{}),
		addr:         addr,
		rpcEndpoint:  rpcEndpoint,
		pollInterval: interval,
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	notifier := &WebsocketNotifier{server: s}
	watcher, err := txnotify.NewWatcher(s.rpcEndpoint, s.pollInterval, notifier)
	if err != nil {
		return err
	}
	s.watcher = watcher
	defer s.watcher.Close()

	go func() {
		if err := s.watcher.Listen(ctx); err != nil {
			log.Printf("watcher error: %v", err)
		}
	}()

	http.HandleFunc("/ws", s.handleWebSocket)
	return http.ListenAndServe(s.addr, nil)
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade error: %v", err)
		return
	}
	defer conn.Close()

	s.mu.Lock()
	s.conns[conn] = make(map[string]struct{})
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.conns, conn)
		s.mu.Unlock()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("read error: %v", err)
			}
			return
		}

		var req SubscriptionRequest
		if err := json.Unmarshal(msg, &req); err != nil {
			log.Printf("unmarshal error: %v", err)
			continue
		}

		if req.Address == "" {
			continue
		}

		if err := s.watcher.Subscribe(req.Address); err != nil {
			log.Printf("subscribe error: %v", err)
			continue
		}

		s.mu.Lock()
		s.conns[conn][req.Address] = struct{}{}
		s.mu.Unlock()
	}
}
