package runner

import (
	"strings"
	"time"
	"io/ioutil"

	"go.uber.org/zap"

	"wist/log"
)

type Client struct {
	Mac string
	IP string
	Name string
}

var Clients []Client

func getDNSLease() string {
	path := "/var/lib/misc/dnsmasq.leases"

	if useFakeData {
		path = "./files/var/lib/misc/dnsmasq.leases"
	}

	out, err := ioutil.ReadFile(path)

	if err != nil {
		log.Error("Could not read file", zap.Error(err))
		return ""
	}

	return string(out)
}

func processDNS(t time.Time) {
	log.Debug("running processDNS")

	str := getDNSLease()

	if len(str) < 1 {
		return
	}

	lines := strings.Split(str, "\n")

	if len(lines) < 1 {
		return
	}

	Clients = make([]Client, len(lines))

	for idx, line := range lines {
		vals := strings.Fields(line)

		if len(vals) < 3 {
			break
		}

		Clients[idx] = Client{
			Mac: vals[1],
			IP: vals[2],
			Name: vals[3],
		}
	}
}
