package runner

import (
	"os/exec"
	"strings"
	"time"
	"regexp"

	"go.uber.org/zap"

	"wist/log"
	"strconv"
	"wist/metrics"
)

type Config struct {
	SSID string
	Frequency string
}

var IWConfig Config

func getIWConfig() string {
	cmd := exec.Command("iwconfig", "wlan1")


	if useFakeData {
		cmd = exec.Command("cat", "/Users/nlevingreenhaw/go/src/wist/files/iwconfig")
	}

	cut := exec.Command("cut", "-c", "11-")

	cmdOut, err := cmd.StdoutPipe()

	if err != nil {
		log.Error("cmd StdoutPipe failed", zap.Error(err))
		return ""
	}

	if err := cmd.Start(); err != nil {
		log.Error("cmd start failed", zap.Error(err))
		return ""
	}

	defer cmd.Process.Kill()

	cut.Stdin = cmdOut

	out, err := cut.Output()

	if err != nil {
		log.Error("cut no stdout", zap.Error(err))
		return ""
	}

	return string(out)
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
