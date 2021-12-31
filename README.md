# Snowflake

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

## Performance

> Benched on Windows 10 - WSL 2 Ubuntu, Intel(R) Core(TM) i7-7700HQ CPU @ 2.80GHz and 12GB of memory

```bash
goos: linux
goarch: amd64
cpu: Intel(R) Core(TM) i7-7700HQ CPU @ 2.80GHz
BenchmarkNewID
BenchmarkNewID/github.com/HotPotatoC/snowflake
BenchmarkNewID/github.com/HotPotatoC/snowflake-8          14012617         87.52 ns/op        0 B/op        0 allocs/op
BenchmarkNewID/github.com/bwmarrin/snowflake
BenchmarkNewID/github.com/bwmarrin/snowflake-8             4918552        244.3 ns/op        0 B/op        0 allocs/op
BenchmarkNewID/github.com/godruoyi/go-snowflake
BenchmarkNewID/github.com/godruoyi/go-snowflake-8          4916791        245.8 ns/op        0 B/op        0 allocs/op
PASS
ok   command-line-arguments 4.230s
```

Disclaimer: Benchmark results may be faster than other implementations. But I do not guarantee that this library is the safest snowflake id generator.

## Support

<a href="https://www.buymeacoffee.com/hotpotato" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: 41px !important;width: 174px !important;box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;-webkit-box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;" ></a>
