# goModeS - Introduction

goModeS is a Go library for decoding Mode S messages. You can use it to decode live data in your command line or import it into your application and use however you see fit.

## Installation and Usage

1. Install the package

```
go get github.com/pragmatic-zac/goModeS
```

2. Import the package

```
import "github.com/pragmatic-zac/goModeS"
```

3. Use it

```
// example of decoding category
message = "8D406B902015A678D4D220AA4BDA"

category, err := decode.Category(message)
if err != nil {
    fmt.Printf("Error: %v", err)
}

fmt.Printf("Category: %d\n", category)  
```

## Command line instructions

Coming soon.

## Work in progress

This package is an active work in progress! Currently, ADS-B messages are supported. 

This application also supports connection to a networked RTL-SDR receiver to decode and display messages in the command line. Eventually I would like to add the ability to connect directly to the RTL-SDR receiver itself.

## Attributions

This package was inspired by antirez's popular dump1090 program. 

[The 1090MHz Riddle](https://mode-s.org/decode/index.html) as a reference made this code possible!