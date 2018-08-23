package runner

import (
	"strconv"
	"strings"
	"time"
	"io/ioutil"

	"go.uber.org/zap"

	"wist/log"
	"wist/metrics"
)

func getWireless() string {
	path := "/proc/net/wireless"

	if useFakeData {
		path = "./files/proc/net/wireless"
	}

	out, err := ioutil.ReadFile(path)

	if err != nil {
		log.Error("Could not read file", zap.Error(err))
		return ""
	}

	return string(out)
}

func processWireless(t time.Time) {
	log.Debug("running processWireless")

	str := getWireless()

	if len(str) < 1 {
		return
	}

	lines := strings.Split(str, "\n")

	if len(lines) < 1 {
		return
	}

	vals := strings.Fields(lines[2])

	if len(vals) < 1 {
		return
	}

	if f, err := strconv.ParseFloat(vals[2], 32); err == nil {
		metrics.Add("link.quality", f)
	}
	if f, err := strconv.ParseFloat(vals[3], 32); err == nil {
		metrics.Add("link.level", f)
	}
	if f, err := strconv.ParseFloat(vals[4], 32); err == nil {
		metrics.Add("link.noise", f)
	}
	if f, err := strconv.ParseFloat(vals[5], 32); err == nil {
		metrics.Add("discard_packets.nwid", f)
	}
	if f, err := strconv.ParseFloat(vals[6], 32); err == nil {
		metrics.Add("discard_packets.crypt", f)
	}
	if f, err := strconv.ParseFloat(vals[7], 32); err == nil {
		metrics.Add("discard_packets.frag", f)
	}
	if f, err := strconv.ParseFloat(vals[8], 32); err == nil {
		metrics.Add("discard_packets.retry", f)
	}
	if f, err := strconv.ParseFloat(vals[9], 32); err == nil {
		metrics.Add("discard_packets.misc", f)
	}
}

func getStat() string {
	path := "/proc/stat"

	if useFakeData {
		path = "./files/proc/stat"
	}

	out, err := ioutil.ReadFile(path)

	if err != nil {
		log.Error("Could not read file", zap.Error(err))
		return ""
	}

	return string(out)
}

func processStat(t time.Time) {
	log.Debug("running processStat")

	str := getStat()

	if len(str) < 1 {
		return
	}

	lines := strings.Split(str, "\n")

	if len(lines) < 1 {
		return
	}

	vals := strings.Fields(lines[0])

	if len(vals) < 1 {
		return
	}

	one, _ := strconv.ParseFloat(vals[1], 32)
	two, _ := strconv.ParseFloat(vals[3], 32)
	thr, _ := strconv.ParseFloat(vals[4], 32)

	cpu := (one + two) * 100 / (one + two + thr)

	metrics.Add("cpu", cpu)
}

func getMemInfo() string {
	path := "/proc/meminfo"

	if useFakeData {
		path = "./files/proc/meminfo"
	}

	out, err := ioutil.ReadFile(path)

	if err != nil {
		log.Error("Could not read file", zap.Error(err))
		return ""
	}

	return string(out)
}

func processMemInfo(t time.Time) {
	log.Debug("running processMemInfo")

	str := getMemInfo()

	if len(str) < 1 {
		return
	}

	lines := strings.Split(str, "\n")

	if len(lines) < 1 {
		return
	}

	vals := strings.Fields(lines[0])

	if len(vals) < 1 {
		return
	}

	total, _ := strconv.ParseFloat(vals[1], 32)

	vals = strings.Fields(lines[1])

	if len(vals) < 1 {
		return
	}

	available, _ := strconv.ParseFloat(vals[1], 32)

	mem := total / available

	metrics.Add("mem", mem)
}
