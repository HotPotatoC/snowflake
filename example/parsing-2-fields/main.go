package main

import (
	"fmt"

	"github.com/HotPotatoC/snowflake"
)

func main() {
	parsed := snowflake.Parse2(1292062458947571712)

	fmt.Printf("Timestamp: %d\n", parsed.Timestamp)       // 1640944495572
	fmt.Printf("Sequence: %d\n", parsed.Sequence)         // 0
	fmt.Printf("Machine ID: %d\n", parsed.Discriminator1) // 1
	fmt.Printf("Process ID: %d\n", parsed.Discriminator2) // 24
}
