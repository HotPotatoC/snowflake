# Snowflake

[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/gomods/athens) [![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/HotPotatoC/go/snowflake) [![GoReportCard example](https://goreportcard.com/badge/github.com/HotPotatoC/snowflake)](https://goreportcard.com/report/github.com/HotPotatoC/snowflake) [![GitHub release](https://img.shields.io/github/release/HotPotatoC/snowflake.svg)](https://GitHub.com/HotPotatoC/snowflake/releases/) [![GitHub license](https://badgen.net/github/license/HotPotatoC/snowflake)](https://github.com/HotPotatoC/snowflake/blob/master/LICENSE) [![codecov](https://codecov.io/gh/HotPotatoC/snowflake/branch/master/graph/badge.svg?token=0BZ6BDOO7O)](https://codecov.io/gh/HotPotatoC/snowflake)

Dead simple and fast [Twitter's snowflake](https://blog.twitter.com/engineering/en_us/a/2010/announcing-snowflake) id generator in Go.

## Installation

```bash
go get github.com/HotPotatoC/snowflake
```

## Usage

1. Generating a snowflake id

```go
machineID := uint64(1)
sf := snowflake.New(machineID)

id := sf.NextID()
fmt.Println(id)
// 1292053924173320192

// or

id = snowflake.New(machineID).NextID()
fmt.Println(id)
// 1292053924173320192
```

2. Parsing a snowflake id

```go
parsed := snowflake.Parse(1292053924173320192)

fmt.Printf("Timestamp: %d\n", parsed.Timestamp)      // 1640942460724
fmt.Printf("Sequence: %d\n", parsed.Sequence)        // 0
fmt.Printf("Machine ID: %d\n", parsed.Discriminator) // 1
```

3. Generating a snowflake ID with 2 discriminator fields

```go
machineID := uint64(1)
processID := uint64(24)
sf := snowflake.New2(machineID, processID)

id := sf.NextID()
fmt.Println(id)
// 1292065108376162304

// or

id = snowflake.New2(machineID, processID).NextID()
fmt.Println(id)
// 1292065108376162304
```

4. Parsing a snowflake id with 2 discriminator fields

```go
parsed := snowflake.Parse2(1292065108376162304)

fmt.Printf("Timestamp: %d\n", parsed.Timestamp)       // 1640944495572
fmt.Printf("Sequence: %d\n", parsed.Sequence)         // 0
fmt.Printf("Machine ID: %d\n", parsed.Discriminator1) // 1
fmt.Printf("Process ID: %d\n", parsed.Discriminator2) // 24
```

## Support

<a href="https://www.buymeacoffee.com/hotpotato" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: 41px !important;width: 174px !important;box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;-webkit-box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;" ></a>
