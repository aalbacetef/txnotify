//go:build js && wasm

package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

func readBuf(buf, _n js.Value) ([]byte, error) {
	fmt.Println("readBuf")
	n := _n.Int()

	dst := make([]byte, buf.Length())
	if read := js.CopyBytesToGo(dst, buf); read != n {
		return nil, fmt.Errorf("error reading: got %d bytes, want %d", read, n)
	}

	return dst, nil
}

func readStr(buf, n js.Value) (string, error) {
	dst, err := readBuf(buf, n)
	if err != nil {
		return "", err
	}

	return string(dst), nil
}

func deserialize[T any](buf, n js.Value) (T, error) {
	var data T

	b, err := readBuf(buf, n)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(b, &data)

	return data, err
}

func serialize[T any](data T) (js.Value, error) {
	encoded, err := json.Marshal(data)
	if err != nil {
		return js.Value{}, fmt.Errorf("error encoding: %w", err)
	}

	n := len(encoded)
	jsArray := js.Global().Get("Uint8Array").New(n)
	js.CopyBytesToJS(jsArray, encoded)

	return jsArray, nil
}
