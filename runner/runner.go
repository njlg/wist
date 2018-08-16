// Package runner ...
package runner

import (
	"time"
)

var useFakeData = false

// Start ...
func Start(fakeData bool) {
	useFakeData = fakeData

	go doEvery(time.Second, processProc)
	go doEvery(time.Second * 30, processDNS)
	go doEvery(time.Second * 5, processIWConfig)
}

func doEvery(d time.Duration, f func(time.Time)) {
	f(time.Now())

	for x := range time.Tick(d) {
		f(x)
	}
}
