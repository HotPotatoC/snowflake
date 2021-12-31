package main

import (
	"fmt"

	"github.com/HotPotatoC/snowflake"
)

func main() {
	machineID := uint64(1)
	sf := snowflake.New(machineID)

	id := sf.NextID()
	fmt.Println(id)

	// or
	id = snowflake.New(machineID).NextID()
	fmt.Println(id)
}