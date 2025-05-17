//go:build js && wasm

package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"syscall/js"
	"time"

	"github.com/aalbacetef/txnotify"
	"github.com/aalbacetef/txnotify/ethereum"
)

const (
	RetOK = iota
	RetErr
)

type State struct {
	watcher     *txnotify.Watcher
	rpcEndpoint string
	cancel      context.CancelFunc
	address     string
}

type Settings struct {
	RPCEndpoint string `json:"rpcEndpoint"`
}

func main() {
	done := make(chan struct{}, 1)
	state := State{}

	cfg := txnotify.Config{PollInterval: time.Second * 15}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	state.cancel = cancel

	onSettingsUpdated := func(this js.Value, args []js.Value) any {
		n := len(args)
		if n != 2 {
			fmt.Println("expected 2 arguments, got ", n)
			return RetErr
		}

		data, err := deserialize[Settings](args[0], args[1])
		if err != nil {
			fmt.Println("error deserializing: ", err)
			return RetErr
		}

		state.rpcEndpoint = data.RPCEndpoint

		return RetOK
	}

	onSubscribe := func(this js.Value, args []js.Value) any {
		if len(args) < 2 {
			return RetErr
		}

		address, err := readStr(args[0], args[1])
		if err != nil {
			fmt.Println("error in onSubscribe: ", err)
			return RetErr
		}

		state.address = address

		return RetOK
	}

	onStarted := func(this js.Value, args []js.Value) any {
		if state.rpcEndpoint == "" {
			fmt.Println("please set RPC endpoint")
			return RetErr
		}

		if state.address == "" {
			fmt.Println("please set address")
			return RetErr
		}

		watcher, err := txnotify.NewWatcher(state.rpcEndpoint, cfg, mockNotifier{})
		if err != nil {
			fmt.Println("error: ", err)
			return RetErr
		}

		if err := watcher.Subscribe(state.address); err != nil {
			fmt.Println("could not subscribe: ", err)
			return RetErr
		}

		state.watcher = watcher

		go func() {
			if err := watcher.Listen(ctx); err != nil {
				fmt.Println("listen ended: ", err)
			}
		}()

		return RetOK
	}

	js.Global().Set("WASM_subscribe", js.FuncOf(onSubscribe))
	js.Global().Set("WASM_updateSettings", js.FuncOf(onSettingsUpdated))
	js.Global().Set("WASM_start", js.FuncOf(onStarted))

	<-done
}

type mockNotifier struct{}

func (mockNotifier) Notify(address string, txList []ethereum.Transaction) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	for _, tx := range txList {
		jsArray, err := serialize[ethereum.Transaction](tx)
		if err != nil {
			logger.Error("could not serialize data", "error", err)
			continue
		}

		js.Global().Get("WASM_listenNotification").Invoke(jsArray)
	}
}
