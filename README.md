# agollo is a golang client for apollo 🚀 [![CircleCI](https://circleci.com/gh/philchia/agollo/tree/master.svg?style=svg)](https://circleci.com/gh/philchia/agollo/tree/master)

[![Go Report Card](https://goreportcard.com/badge/github.com/philchia/agollo)](https://goreportcard.com/report/github.com/philchia/agollo)
[![golang](https://img.shields.io/badge/Language-Go-green.svg?style=flat)](https://golang.org)
[![GoDoc](https://godoc.org/github.com/philchia/zen?status.svg)](https://godoc.org/github.com/philchia/agollo)


## Feature

* Multiple namespace support
* Fail tolerant
* Zero dependency

## Dependency

required **go 1.9** or later

## Installation

```sh
    go get -u github.com/philchia/agollo
```

## Usage

### Start use default app.properties config file

```golang
    agollo.Start()
```

### Start use given config file path

```golang
    agollo.StartWithConfFile(name)
```

### Subscribe to updates

```golang
    events := agollo.WatchUpdate()
    changeEvent := <-event
    bytes, _ := json.Marshal(changeEvent)
    fmt.Println("event:", string(bytes))
```

### Get apollo values

```golang
    agollo.GetStringValue(Key, defaultValue)
    agollo.GetStringValueWithNameSapce(namespace, key, defaultValue)
```

## License

agollo is released under MIT lecense