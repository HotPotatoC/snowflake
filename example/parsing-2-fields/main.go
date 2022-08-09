package main

import (
	"fmt"

	"github.com/HotPotatoC/snowflake"
)

func main() {
	parsed := snowflake.Parse2(1292065108376162304)

	fmt.Printf("Timestamp: %d\n", parsed.Timestamp)       // 1640945127245
	fmt.Printf("Sequence: %d\n", parsed.Sequence)         // 0
	fmt.Printf("Machine ID: %d\n", parsed.Field1) // 1
	fmt.Printf("Process ID: %d\n", parsed.Field2) // 24
}
