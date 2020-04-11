# agollo 是携程 apollo 配置中心的 golang 客户端 🚀 [![CircleCI](https://circleci.com/gh/philchia/agollo/tree/master.svg?style=svg)](https://circleci.com/gh/philchia/agollo/tree/master)

[![Go Report Card](https://goreportcard.com/badge/github.com/philchia/agollo)](https://goreportcard.com/report/github.com/philchia/agollo)
[![codebeat badge](https://codebeat.co/badges/e31b4a09-f531-4b74-a86a-775f46436539)](https://codebeat.co/projects/github-com-philchia-agollo-master)
[![Coverage Status](https://coveralls.io/repos/github/philchia/agollo/badge.svg?branch=master)](https://coveralls.io/github/philchia/agollo?branch=master)
[![golang](https://img.shields.io/badge/Language-Go-green.svg?style=flat)](https://golang.org)
[![GoDoc](https://godoc.org/github.com/philchia/zen?status.svg)](https://godoc.org/github.com/philchia/agollo)
![GitHub release](https://img.shields.io/github/release/philchia/agollo.svg)

## 功能

* 多 namespace 支持
* 容错，本地缓存
* 零依赖
* 实时更新通知

## 依赖

**go 1.9** 或更新

## 安装

```sh
    go get -u github.com/philchia/agollo/v3
```

## 使用

### Import agollo
```golang
    import github.com/philchia/agollo/v3
```

### 使用 app.properties 配置文件启动

```golang
    agollo.Start()
```

### 使用自定义配置启动

```golang
    agollo.StartWithConfFile(name)
```

### 监听配置更新

```golang
    events := agollo.WatchUpdate()
    changeEvent := <-events
    bytes, _ := json.Marshal(changeEvent)
    fmt.Println("event:", string(bytes))
```

### 获取配置

```golang
    agollo.GetStringValue(Key, defaultValue)
    agollo.GetStringValueWithnamespace(namespace, key, defaultValue)
```

### 获取文件内容

```golang
    agollo.GetNameSpaceContent(namespace, defaultValue)
```

### 获取配置中所有的键

```golang
    agollo.GetAllKeys(namespace)
```

### 订阅namespace的配置

```golang
    agollo.SubscribeToNamespaces("newNamespace1", "newNamespace2")
```

## 许可

agollo 使用 MIT 许可
