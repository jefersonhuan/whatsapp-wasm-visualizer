package main

import (
	"bytes"
	"fmt"
	"github.com/jefersonhuan/whatsapp-vizualizer-wasm/main/parser"
	"syscall/js"
	"time"
)

func handleFileBytes(this js.Value, args []js.Value) interface{} {
	start := time.Now()
	var buf []byte
	array := args[0]

	buf = make([]byte, array.Get("byteLength").Int())
	js.CopyBytesToGo(buf, array)

	result, nMessages := parser.Parse(bytes.NewReader(buf))
	data := parser.Convert(result)

	js.Global().Call("loadChart", data, nMessages)
	fmt.Println("Parsing took", time.Now().Sub(start))
	return js.Undefined()
}

func main() {
	c := make(chan bool)
	fmt.Println("WhatsApp Chat parser has been initialized")

	js.Global().Set("parseChat", js.FuncOf(handleFileBytes))
	<-c
}
