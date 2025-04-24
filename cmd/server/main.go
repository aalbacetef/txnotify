package main

import (
	"context"
	"flag"
	"log"
)

func main() {
	addr := ":8080"
	rpcEndpoint := "https://eth.nodeconnect.org"
	pollInterval := "5s"

	flag.StringVar(&addr, "addr", addr, "server address")
	flag.StringVar(&rpcEndpoint, "rpc", rpcEndpoint, "RPC endpoint")
	flag.StringVar(&pollInterval, "interval", pollInterval, "poll interval")
	flag.Parse()

	if rpcEndpoint == "" || pollInterval == "" {
		flag.Usage()
		return
	}

	server, err := NewServer(addr, rpcEndpoint, pollInterval)
	if err != nil {
		log.Fatalf("server init error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := server.Start(ctx); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
