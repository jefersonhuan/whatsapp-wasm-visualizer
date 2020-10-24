package main

import (
	"bytes"
	"fmt"
	"github.com/jefersonhuan/whatsapp-vizualizer-wasm/main/parser"
	"syscall/js"
	"time"
)

var c chan bool

func init() {
	c = make(chan bool)
}

func handleFileBytes(this js.Value, args []js.Value) interface{} {
	start := time.Now()
	var buf []byte
	array := args[0]

	buf = make([]byte, array.Get("byteLength").Int())
	js.CopyBytesToGo(buf, array)

	data := parser.LoadChat(bytes.NewReader(buf)).Parse()
	dst := js.Global().Get("Uint32Array").New(len(data))
	for i, d := range data {
		dst.SetIndex(i, d)
	}
	fmt.Println("Parsing took", time.Now().Sub(start))
	return dst
}

func main() {
	fmt.Println("WhatsApp Chat parser has been initialized")

	js.Global().Set("parseChat", js.FuncOf(handleFileBytes))
	<-c
}
