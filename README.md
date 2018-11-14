# agollo is a golang client for apollo 🚀 [![CircleCI](https://circleci.com/gh/philchia/agollo/tree/master.svg?style=svg)](https://circleci.com/gh/philchia/agollo/tree/master)

[![Go Report Card](https://goreportcard.com/badge/github.com/philchia/agollo)](https://goreportcard.com/report/github.com/philchia/agollo)
[![codebeat badge](https://codebeat.co/badges/e31b4a09-f531-4b74-a86a-775f46436539)](https://codebeat.co/projects/github-com-philchia-agollo-master)
[![Coverage Status](https://coveralls.io/repos/github/philchia/agollo/badge.svg?branch=master)](https://coveralls.io/github/philchia/agollo?branch=master)
[![golang](https://img.shields.io/badge/Language-Go-green.svg?style=flat)](https://golang.org)
[![GoDoc](https://godoc.org/github.com/philchia/zen?status.svg)](https://godoc.org/github.com/philchia/agollo)
![GitHub release](https://img.shields.io/github/release/philchia/agollo.svg)

## Simple chinese

[简体中文](./README_CN.md)

## Feature

* Multiple namespace support
* Fail tolerant
* Zero dependency
* Realtime change notification

## Required

**go 1.9** or later

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
    changeEvent := <-events
    bytes, _ := json.Marshal(changeEvent)
    fmt.Println("event:", string(bytes))
```

### Get apollo values

```golang
    agollo.GetStringValue(Key, defaultValue)
    agollo.GetStringValueWithNameSpace(namespace, key, defaultValue)
```

### Get namespace file contents

```golang
    agollo.GetNameSpaceContent(namespace, defaultValue)
```

### Get all keys

```golang
    agollo.GetAllKeys(namespace)
```

## License

agollo is released under MIT license
