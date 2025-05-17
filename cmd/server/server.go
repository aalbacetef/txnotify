package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/aalbacetef/txnotify"
	"github.com/aalbacetef/txnotify/ethereum"
)

const (
	defaultBufSize      = 1024
	defaultReadTimeout  = 10 * time.Second
	defaultWriteTimeout = 10 * time.Second
	// NOTE: this should probably be lower with clients sending keep alive messages.
	defaultIdleTimeout       = 5 * time.Minute
	defaultReadHeaderTimeout = 5 * time.Second
)

type Server struct {
	mu           sync.Mutex
	conns        map[*websocket.Conn]map[string]struct{}
	watcher      *txnotify.Watcher
	addr         string
	rpcEndpoint  string
	pollInterval time.Duration
	upgrader     websocket.Upgrader
}

type SubscriptionRequest struct {
	Address string `json:"address"`
}

type Notification struct {
	Address string                 `json:"address"`
	Txs     []ethereum.Transaction `json:"transactions"` //nolint:tagliatelle
}

func NewServer(addr, rpcEndpoint, pollInterval string) (*Server, error) {
	interval, err := time.ParseDuration(pollInterval)
	if err != nil {
		return nil, fmt.Errorf("could not parse duration: %w", err)
	}

	return &Server{
		conns:        make(map[*websocket.Conn]map[string]struct{}),
		addr:         addr,
		rpcEndpoint:  rpcEndpoint,
		pollInterval: interval,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  defaultBufSize,
			WriteBufferSize: defaultBufSize,
			CheckOrigin:     func(_ *http.Request) bool { return true },
		},
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	notifier := &WebsocketNotifier{server: s}
	cfg := txnotify.Config{PollInterval: s.pollInterval}

	watcher, err := txnotify.NewWatcher(s.rpcEndpoint, cfg, notifier)
	if err != nil {
		return fmt.Errorf("NewWatcher: %w", err)
	}
	s.watcher = watcher
	defer s.watcher.Close()

	go func() {
		if err := s.watcher.Listen(ctx); err != nil {
			log.Printf("watcher error: %v", err)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.handleWebSocket)

	srv := &http.Server{
		Addr:              s.addr,
		Handler:           mux,
		ReadTimeout:       defaultReadTimeout,
		WriteTimeout:      defaultWriteTimeout,
		IdleTimeout:       defaultIdleTimeout,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("ListenAndServe: %w", err)
	}

	return nil
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
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
