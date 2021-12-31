package main

import (
	"fmt"

	"github.com/HotPotatoC/snowflake"
)

func main() {
	machineID := uint64(1)
	processID := uint64(24)
	sf := snowflake.New2(machineID, processID)

	id := sf.NextID()
	fmt.Println(id)
	// 1292062458947571712

	// or
	id = snowflake.New2(machineID, processID).NextID()
	fmt.Println(id)
	// 1292062458947571712
}