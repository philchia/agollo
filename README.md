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
* Realtime change notification
* Unmarshal

## Required

**go 1.9** or later

## Installation

```sh
    go get -u github.com/philchia/agollo
```

## Usage

### 1, Start

#### Start use default app.properties config file

```golang
    agollo.Start()
```

#### Start use given config file path

```golang
    agollo.StartWithConfFile(name)
```

#### Start with given conf struct
```golang
    agollo.StartWithConf(yourConf)
```

### 2, Subscribe updates

#### WatchUpdate

```golang
    events := agollo.WatchUpdate()
    changeEvent := <-events
    bytes, _ := json.Marshal(changeEvent)
    fmt.Println("event:", string(bytes))
```

#### Or watch any change with OnConfigChange
```golang
	agollo.OnConfigChange(func() {
        // to do
	})
```

### 3, Get config using multi-method

#### Get apollo values

```golang
    // default namespace: application
    agollo.GetStringValue(Key, defaultValue)
    // user specify namespace
    agollo.GetStringValueWithNameSpace(namespace, key, defaultValue)
```

#### Get namespace file contents

```golang
    agollo.GetNameSpaceContent(namespace, defaultValue)
```

#### Get all keys

```golang
    agollo.GetAllKeys(namespace)
```

#### Unmarshal

There is a config in apollo like this:
![](https://github.com/xujintao/agollo/blob/master/apollo.png)


So our meta-config should like:
```json
{
    "appId": "001",
    "cluster": "default",
    "namespaceNames": ["application","dnspod1","dnspod2.yaml","db"],
    "ip": "localhost:8080"
}
```

At last, we make a structure to get all the config
```golang
package main

import (
	"fmt"
	"log"

	"github.com/philchia/agollo"
)

type config struct {
    // dns 配置
    DNS1 struct {
        ID     string `mapstructure:"id"`
        Token  string `mapstructure:"token"`
        Domain string `mapstructure:"domain"`
    } `mapstructure:"dnspod1"`
    DNS2 struct {
        ID     int    `mapstructure:"id"`
        Token  string `mapstructure:"token"`
        Domain string `mapstructure:"domain"`
    } `mapstructure:"dnspod2.yaml"`

    // DB
    DB struct {
        DSN     string `mapstructure:"dsn"`
        MaxConn string `mapstructure:"max_conn"`
    } `mapstructure:"db"`
}

func main(){
    agollo.Start()

    // first
    var c config
    agollo.Unmarshal(&c)
    fmt.Printf("%v", c)

    // watch
	agollo.OnConfigChange(func() {
		var c config
		agollo.Unmarshal(&c)
		fmt.Println(c)
	})
}
```
