package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello there")

	var ns uint64
	clock_time_get(0, 1000, &ns) // direct wasi call example (see: https://github.com/WasmEdge/WasmEdge/blob/4f1059b4aa934c4de3c3edff4055c28e6f0b9311/lib/host/wasi/wasimodule.cpp)

	fmt.Println("Time from WASI: ", time.Unix(0, int64(ns)))

	var buf [10]byte
	random_get(&buf[0], int32(len(buf)))

	fmt.Println("Random from WASI: ", buf)
}

//go:wasm-module wasi_snapshot_preview1
//export clock_time_get
func clock_time_get(clockid uint32, precision uint64, time *uint64) (errno uint16)

//go:wasm-module wasi_snapshot_preview1
//export random_get
func random_get(buf *byte, len int32) (errno uint16)
