package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/aalbacetef/txnotify"
	"github.com/aalbacetef/txnotify/ethereum"
)

type mockNotifier struct{}

func (mockNotifier) Notify(address string, txList []ethereum.Transaction) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	for _, tx := range txList {
		logger.Info(
			"notification: got tx",
			"address", address,
			"hash", tx.Hash,
		)
	}
}

func main() {
	address := ""
	pollInterval := "5s"
	rpcEndpoint := "https://eth.nodeconnect.org"

	flag.StringVar(&address, "address", address, "address to subscribe to")
	flag.StringVar(&pollInterval, "interval", pollInterval, "poll interval")
	flag.StringVar(&rpcEndpoint, "rpc", rpcEndpoint, "RPC endpoint")

	flag.Parse()

	if address == "" || pollInterval == "" || rpcEndpoint == "" {
		flag.Usage()
		return
	}

	interval, err := time.ParseDuration(pollInterval)
	if err != nil {
		fmt.Println("error parsing duration: ", err)
		return
	}

	cfg := txnotify.Config{PollInterval: interval}

	watcher, err := txnotify.NewWatcher(rpcEndpoint, cfg, mockNotifier{})
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer watcher.Close()

	if err := watcher.Subscribe(address); err != nil {
		fmt.Println("subscribe error: ", err)
		return
	}

	if err := watcher.Listen(ctx); err != nil {
		fmt.Println("error: ", err)
		return
	}
}
