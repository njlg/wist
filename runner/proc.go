package runner

import (
	"os/exec"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"wist/log"
	"wist/metrics"
)

func getProc() string {
	path := "/proc/net/wireless"

	if useFakeData {
		path = "/Users/nlevingreenhaw/go/src/wist/files/proc/net/wireless"
	}

	out, err := exec.Command("cat", path).Output()

	if err != nil {
		log.Error("Could not read file", zap.Error(err))
		return ""
	}

	return string(out)
}

func processProc(t time.Time) {
	log.Debug("running processProc")

	str := getProc()

	if len(str) < 1 {
		return
	}

	// log.Debug("str", zap.Int("len", len(str)), zap.String("st", str))

	lines := strings.Split(str, "\n")

	if len(lines) < 1 {
		return
	}

	// log.Debug("lines", zap.Int("len", len(lines)), zap.String("st", lines[2]))

	vals := strings.Fields(lines[2])

	if len(vals) < 1 {
		return
	}

	// log.Debug("vals", zap.Int("len", len(vals)), zap.String("st", vals[2]))

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
