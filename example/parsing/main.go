package main

import (
	"fmt"

	"github.com/HotPotatoC/snowflake"
)

func main() {
	parsed := snowflake.Parse(1292053924173320192)

	fmt.Printf("Timestamp: %d\n", parsed.Timestamp)      // 1640942460724
	fmt.Printf("Sequence: %d\n", parsed.Sequence)        // 0
	fmt.Printf("Machine ID: %d\n", parsed.Discriminator) // 1
}
