package agollo

import (
	"fmt"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	if err := StartWithConfFile("./testdata/app.properties"); err != nil {
		t.Error(err)
	}

	defer Stop()

	changes := WatchUpdate()

	go func() {
		for {
			time.Sleep(time.Second)
			timeout := GetStringValue("timeout", "default")
			fmt.Println("timeout is:", timeout)
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
