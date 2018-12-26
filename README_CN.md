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
* 实时更新通知
* 支持Umarshal

## 依赖

**go 1.9** 或更新

## 安装

```sh
    go get -u github.com/philchia/agollo
```

## 使用

### 1, 启动

#### 使用 app.properties 配置文件启动

```golang
    agollo.Start()
```

#### 使用自定义配置启动

```golang
    agollo.StartWithConfFile(name)
```

#### 使用自定义结构启动
```golang
    agollo.StartWithConf(yourConf)
```

### 2, 热更新

#### 监听配置更新(回头把用户代码用匿名函数包装起来注册到WatchUpdate)

```golang
    events := agollo.WatchUpdate()
    changeEvent := <-event
    bytes, _ := json.Marshal(changeEvent)
    fmt.Println("event:", string(bytes))
```

### 3, 多种方式获取配置

#### 获取properties配置

```golang
    // default namespace: application
    agollo.GetStringValue(Key, defaultValue)

    // user specify namespace
    agollo.GetStringValueWithNameSapce(namespace, key, defaultValue)
```

#### 获取文件内容

```golang
    agollo.GetNameSpaceContent(namespace, defaultValue)
```

#### 获取配置中所有的键

```golang
    agollo.GetAllKeys(namespace)
```

#### 用Unmarshal获取配置

假设配置中心是这样配置的:
![](https://github.com/xujintao/agollo/blob/master/apollo.png)


那么，我们的元配置(app.properties)应该这样写：
```json
{
    "appId": "001",
    "cluster": "default",
    "namespaceNames": ["application","dnspod1","dnspod2.yaml","db"],
    "ip": "localhost:8080"
}
```

然后像这样定义一个struct去获取所有的配置：
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

    // 第一次读取
    var c config
    agollo.Unmarshal(&c)
    fmt.Printf("%v", c)

    // 热更新
	agollo.OnConfigChange(func() {
		var c config
		agollo.Unmarshal(&c)
		fmt.Println(c)
	})
}
```
