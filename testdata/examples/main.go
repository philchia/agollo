package main

import (
	"fmt"
	"time"

	"github.com/philchia/agollo/v4"
)

func main() {
	_ = agollo.Start(&agollo.Conf{
		AppID:          "SampleApp",
		Cluster:        "",
		NameSpaceNames: []string{"test-json.json"},
		MetaAddr:       "http://106.54.227.205:8080",
	}, agollo.SkipLocalCache())
	agollo.OnUpdate(func(event *agollo.ChangeEvent) {
		fmt.Println("change:", event)
	})
	appContent := agollo.GetContent()
	fmt.Println("application content is", appContent)

	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			fmt.Println("+1s")
		default:
		}
	}
}
