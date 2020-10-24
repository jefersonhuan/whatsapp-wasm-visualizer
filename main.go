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
	result := parser.LoadChat(bytes.NewReader(buf)).Parse()
	fmt.Println(result)

	data := parser.Convert(result)
	fmt.Println(data)

	js.Global().Call("loadChart", data)
	fmt.Println("Parsing took", time.Now().Sub(start))
	return js.Undefined()
}

func main() {
	fmt.Println("WhatsApp Chat parser has been initialized")

	js.Global().Set("parseChat", js.FuncOf(handleFileBytes))
	<-c
}
