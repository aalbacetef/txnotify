//go:build js && wasm

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"syscall/js"
	"time"

	"github.com/aalbacetef/txnotify"
	"github.com/aalbacetef/txnotify/ethereum"
)

type ReturnCode int

const (
	RetOK ReturnCode = iota
	RetErr
)

func main() {
	done := make(chan struct{}, 1)

	cfg := txnotify.Config{PollInterval: time.Second * 15}
	rpcEndpoint := "https://eth.llamarpc.com"

	watcher, err := txnotify.NewWatcher(rpcEndpoint, cfg, mockNotifier{})
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer watcher.Close()

	onSettingsUpdated := func(this js.Value, args []js.Value) any {
		fmt.Println("onSettingsUpdated: ", args)
		return RetOK
	}

	onSubscribe := func(this js.Value, args []js.Value) any {
		fmt.Println("onSubscribe: ", args)
		if len(args) < 2 {
			return 1
		}

		address, err := readStr(args[0], args[1])
		if err != nil {
			fmt.Println("onSubscribe: ", err)
			return 1
		}

		fmt.Println("[GO] address: ", address)

		if err := watcher.Subscribe(address); err != nil {
			fmt.Println("could not subscribe: ", err)
			return 1
		}

		return 0
	}

	onStarted := func(this js.Value, args []js.Value) any {
		fmt.Println("onStarted: ", args)
		return RetOK
	}

	js.Global().Set("WASM_subscribe", js.FuncOf(onSubscribe))
	js.Global().Set("WASM_settingsUpdated", js.FuncOf(onSettingsUpdated))
	js.Global().Set("WASM_start", js.FuncOf(onStarted))

	if err := watcher.Listen(ctx); err != nil {
		fmt.Println("listen ended: ", err)
	}

	<-done
}

func readStr(buf, _n js.Value) (string, error) {
	n := _n.Int()
	dst := make([]byte, n)
	if read := js.CopyBytesToGo(dst, buf); read != n {
		return "", fmt.Errorf("error reading: got %d bytes, want %d", read, n)
	}

	return string(dst), nil
}

type mockNotifier struct{}

func (mockNotifier) Notify(address string, txList []ethereum.Transaction) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	for _, tx := range txList {

		data, err := json.Marshal(tx)
		if err != nil {
			logger.Error("error encoding", "error", err)
			continue
		}

		n := len(data)
		jsArray := js.Global().Get("Uint8Array").New(n)
		js.CopyBytesToJS(jsArray, data)
		js.Global().Get("WASM_listenNotification").Invoke(jsArray)
	}
}
