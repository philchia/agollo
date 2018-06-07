package agollo

import (
	"fmt"

	"github.com/philchia/agollo/internal/mockserver"
)

// func TestMain(m *testing.M) {
// 	setup()
// 	defer teardown()
// 	m.Run()
// }

func setup() {
	go func() {
		fmt.Println("start mock server")
		if err := mockserver.Run(); err != nil {
			fmt.Println(err)
		}
	}()
}

func teardown() {
	mockserver.Close()
}

// func TestAgolloStart(t *testing.T) {
// 	if err := StartWithConfFile("./testdata/app.properties"); err != nil {
// 		t.FailNow()
// 	}

// 	defer Stop()

// 	mockserver.Set("application", "key", "value")

// 	updates := WatchUpdate()

// 	select {
// 	case event := <-updates:
// 		_ = event
// 	case <-time.After(time.Millisecond * 3000):
// 	}

// 	val := GetStringValue("key", "defaultValue")
// 	_ = val
// }
