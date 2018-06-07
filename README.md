# agollo is a golang client for apollo

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
