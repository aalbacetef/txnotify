package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/aalbacetef/txnotify"
	"github.com/aalbacetef/txnotify/ethereum"
)

type mockNotifier struct{}

func (mockNotifier) Notify(address string, txList []ethereum.Transaction) {
	for _, tx := range txList {
		fmt.Printf("%s) got tx: %s\n", address, tx.Hash)
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

	watcher, err := txnotify.NewWatcher(rpcEndpoint, interval, mockNotifier{})
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
