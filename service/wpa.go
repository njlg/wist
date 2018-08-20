package service

import (
	"io/ioutil"
	"wist/log"
	"go.uber.org/zap"
	"strings"
)

type Network struct {
	Name string
	Password string
	Priority int
}

var useFakeData = true

func readFile() string {
	path := "/etc/wpa_supplicant/wpa_supplicant-wlan1.conf"

	if useFakeData {
		path = "/Users/nlevingreenhaw/go/src/wist/files/etc/wpa_supplicant/wpa_supplicant-wlan1.conf"
	}

	out, err := ioutil.ReadFile(path)

	if err != nil {
		log.Error("Could not read file", zap.Error(err))
		return ""
	}

	return string(out)
}

func parseConfig() {
	content := readFile()

	if len(content) < 1 {
		return
	}

	lines := strings.Split(content, "\n")

	for _, line := range lines {
		if line == "network={" {

		}
	}

}