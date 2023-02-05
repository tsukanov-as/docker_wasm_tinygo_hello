package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello there")

	var ns uint64
	clock_time_get(0, 1000, &ns) // WASI direct call example

	fmt.Println("Time from WASI: ", time.Unix(0, int64(ns)))
}

//go:wasm-module wasi_snapshot_preview1
//export clock_time_get
func clock_time_get(clockid uint32, precision uint64, time *uint64) (errno uint16)
