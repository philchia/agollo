# agollo is a golang client for apollo 🚀 

[![Go Report Card](https://goreportcard.com/badge/github.com/philchia/agollo)](https://goreportcard.com/report/github.com/philchia/agollo)
[![codebeat badge](https://codebeat.co/badges/e31b4a09-f531-4b74-a86a-775f46436539)](https://codebeat.co/projects/github-com-philchia-agollo-master)
[![Coverage Status](https://coveralls.io/repos/github/philchia/agollo/badge.svg?branch=v4)](https://coveralls.io/github/philchia/agollo?branch=v4)
[![golang](https://img.shields.io/badge/Language-Go-green.svg?style=flat)](https://golang.org)
[![GoDoc](https://godoc.org/github.com/philchia/agollo?status.svg)](https://pkg.go.dev/github.com/philchia/agollo/v4)
![GitHub release](https://img.shields.io/github/release/philchia/agollo.svg)

**v1:**[![goproxy.cn](https://goproxy.cn/stats/github.com/philchia/agollo/badges/download-count.svg)](https://goproxy.cn)
**v3:**[![goproxy.cn](https://goproxy.cn/stats/github.com/philchia/agollo/v3/badges/download-count.svg)](https://goproxy.cn)
**v4:**[![goproxy.cn](https://goproxy.cn/stats/github.com/philchia/agollo/v4/badges/download-count.svg)](https://goproxy.cn)


## Feature

* Multiple namespace support
* Fail tolerant
* Zero dependency
* Realtime change notification
* API to get contents of namespace

## Required

**go 1.11** or later

## Installation

```sh
go get -u github.com/philchia/agollo/v4
```

## Usage

### Import agollo

```golang
import "github.com/philchia/agollo/v4"
```

### In order to use agollo, issue a client or use the built-in default client

#### to use the default global client

for namespaces with the format of properties, you need to specific the full name 

```golang
agollo.Start(&agollo.Conf{
    AppID:          "your app id",
    Cluster:        "default",
    NameSpaceNames: []string{"application.properties"},
    MetaAddr:       "your apollo meta addr",
})
```

#### or to issue a new client to embedded into your program

```golang
apollo := agollo.New(&agollo.Conf{
                            AppID:          "your app id",
                            Cluster:        "default",
                            NameSpaceNames: []string{"application.properties"},
                            MetaAddr:       "your apollo meta addr",
                        })
apollo.Start()
```

### Set config update callback

```golang
agollo.OnUpdate(func(event *ChangeEvent) {
    // do your business logic to handle config update
})
```

### Get apollo values

```golang
// get values in the application.properties default namespace
val := agollo.GetString(Key)
// or indicate a namespace
other := agollo.GetString(key, agollo.WithNamespace("other namespace"))
```

### Get namespace file contents
**to get namespace of with a format of properties, you need to specific the full name of the namespace, e.g. namespace.properties in both options and configs**
```golang
namespaceContent := agollo.GetContent(agollo.WithNamespace("application.properties"))
```

### Get all keys

```golang
allKyes := agollo.GetAllKeys(namespace)
```

### Subscribe to new namespaces

```golang
agollo.SubscribeToNamespaces("newNamespace1", "newNamespace2")
```

## License

agollo is released under MIT license
