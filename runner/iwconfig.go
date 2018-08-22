package runner

import (
	"os/exec"
	"strings"
	"time"
	"regexp"

	"wist/log"
	"strconv"
	"wist/metrics"
	"io"
	"bytes"
)

type Config struct {
	SSID string
	Frequency string
}

var IWConfig Config

func getIWConfig() string {
	cmd := exec.Command("iwconfig", "wlan1")
	cut := exec.Command("cut", "-c", "11-")

	if useFakeData {
		cmd = exec.Command("cat", "/Users/nlevingreenhaw/go/src/wist/files/iwconfig")
	}

	r, w := io.Pipe()
	cmd.Stdout = w
	cut.Stdin = r

	var output bytes.Buffer
	cut.Stdout = &output

	cmd.Start()
	cut.Start()
	cmd.Wait()
	w.Close()
	cut.Wait()

	return output.String()
}

func processIWConfig(t time.Time) {
	log.Debug("running processIWConfig")

	str := getIWConfig()

	if len(str) < 1 {
		return
	}

	lines := strings.Split(str, "\n")

	if len(lines) < 1 {
		return
	}

	// ssid
	r, _ := regexp.Compile(`ESSID:"(.+)"`)

	match := r.FindStringSubmatch(lines[0])

	if len(match) < 1 {
		return
	}

	ssid := match[1]

	// frequence
	r, _ = regexp.Compile(`Frequency:([0-9.]+ GHz)`)
	match = r.FindStringSubmatch(lines[1])

	if len(match) < 1 {
		return
	}

	freq := match[1]

	IWConfig = Config{
		SSID: ssid,
		Frequency: freq,
	}

	// bit rate
	r, _ = regexp.Compile(`Bit Rate=([0-9.]+) Mb/s`)

	match = r.FindStringSubmatch(lines[2])

	if len(match) < 1 {
		return
	}

	if f, err := strconv.ParseFloat(match[1], 32); err == nil {
		metrics.Add("bitrate", f)
	}
}
