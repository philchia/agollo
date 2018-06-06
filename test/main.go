package main

import (
	"fmt"
	"log"
	"time"

	"github.com/philchia/agollo"
)

func main() {
	if err := agollo.StartWithConfFile("../testdata/app.properties"); err != nil {
		log.Fatal(err)
	}

	defer agollo.Stop()

	changes := agollo.WatchUpdate()

	time.AfterFunc(time.Second*5, func() { agollo.Stop() })

	go func() {
		for {
			time.Sleep(time.Second)
			timeout := agollo.GetStringValue("timeout", "default")
			fmt.Println("timeout is:", timeout)
			fmt.Println("null value is", agollo.GetStringValue("null", "null"))
			fmt.Println("Client.json is", agollo.GetStringValueWithNameSapce("Client.json", "content", "null"))
		}
	}()
	for change := range changes {
		if change != nil {
			fmt.Println(change.Namespace)
			for k, change := range change.Changes {
				fmt.Println("change", k, "type", change.ChangeType, "val", change.NewValue)
			}
		}
	}
}
