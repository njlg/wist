package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"

	// "github.com/pkg/errors"
	"github.com/gobuffalo/packr"
	"github.com/gorilla/websocket"
	"github.com/urfave/cli"
	"go.uber.org/zap"

	"wist/log"
	"wist/metrics"
	"wist/runner"
)

var upgrader = websocket.Upgrader{}
var box = packr.NewBox("../templates")
var public = packr.NewBox("../public")

// page struct is for html response
type page struct {
	Host   string

	Signal string
	Level string
	Noise string
	BitRate string

	PacketNWID string
	PacketCrypt string
	PacketFrag string
	PacketRetry string
	PacketMisc string

	Clients []runner.Client

	IWConfig runner.Config
}

// data struct is for websocket response
type data struct {
	Signal []float64 `json:"signal,omitempty"`
	Level []float64 `json:"level,omitempty"`
	Noise []float64 `json:"noise,omitempty"`

	NWID float64 `json:"nwid,omitempty"`
	Crypt float64 `json:"crypt,omitempty"`
	Frag float64 `json:"frag,omitempty"`
	Retry float64 `json:"retry,omitempty"`
	Misc float64 `json:"misc,omitempty"`

	BitRate []float64 `json:"bitrate,omitempty"`
}

// Run ...
func Run(ctx *cli.Context) error {
	log.Info("Service starting", zap.String("version", ctx.App.Version))

	// prepare http bind address
	port := ctx.String("port")
	addr := strings.Builder{}
	addr.WriteString("localhost:")
	addr.WriteString(port)

	server := &http.Server{
		Addr: addr.String(),
	}

	useFakeData := ctx.Bool("fakedata")

	s := box.String("index.html")
	tmpl, _ := template.New("name").Parse(s)

	m := http.NewServeMux()
	m.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info("HTTP Request", zap.String("path", r.URL.Path))

		if public.Has(r.URL.Path) && r.URL.Path != "/" {
			res := public.String(r.URL.Path)
			w.Write([]byte(res))
			return
		}

		// signal info
		signalQualityMetric, _ := metrics.GetAll("link.quality")
		signalLevelMetric, _ := metrics.GetAll("link.level")
		signalNoiseMetric, _ := metrics.GetAll("link.noise")

		bitRate, _ := metrics.GetAll("bitrate")

		// packet info
		packetNWID, _ := metrics.Get("discard_packets.nwid")
		packetCrypt, _ := metrics.Get("discard_packets.crypt")
		packetFrag, _ := metrics.Get("discard_packets.frag")
		packetRetry, _ := metrics.Get("discard_packets.retry")
		packetMisc, _ := metrics.Get("discard_packets.misc")


		data := page{
			Host:   r.Host,
			Signal: makeJSArray(signalQualityMetric),
			Level: makeJSArray(signalLevelMetric),
			Noise: makeJSArray(signalNoiseMetric),
			BitRate: makeJSArray(bitRate),

			PacketNWID: fmt.Sprintf("%f", packetNWID),
			PacketCrypt: fmt.Sprintf("%f", packetCrypt),
			PacketFrag: fmt.Sprintf("%f", packetFrag),
			PacketRetry: fmt.Sprintf("%f", packetRetry),
			PacketMisc: fmt.Sprintf("%f", packetMisc),

			Clients: runner.Clients,
			IWConfig: runner.IWConfig,
		}
		tmpl.Execute(w, data)
	}))
	m.Handle("/data", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info("Websocket Request")

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error("websocket upgrade", zap.Error(err))
			return
		}

		defer c.Close()

		for _ = range time.Tick(time.Second * 5) {
			// signal info
			linkQuality, _ := metrics.GetAll("link.quality")
			linkLevel, _ := metrics.GetAll("link.level")
			linkNoise, _ := metrics.GetAll("link.noise")

			// bit rate
			bitRate, _ := metrics.GetAll("bitrate")

			// packet info
			packetNWID, _ := metrics.Get("discard_packets.nwid")
			packetCrypt, _ := metrics.Get("discard_packets.crypt")
			packetFrag, _ := metrics.Get("discard_packets.frag")
			packetRetry, _ := metrics.Get("discard_packets.retry")
			packetMisc, _ := metrics.Get("discard_packets.misc")

			d := data{
				Signal: linkQuality,
				Level: linkLevel,
				Noise: linkNoise,
				NWID: packetNWID,
				Crypt: packetCrypt,
				Frag: packetFrag,
				Retry: packetRetry,
				Misc: packetMisc,
				BitRate: bitRate,
			}

			message, err := json.Marshal(d)
			if err != nil {
				log.Error("json marshal", zap.Error(err))
				break
			}

			log.Debug("websocket sending", zap.String("data", string(message)))

			err = c.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Error("websocket write", zap.Error(err))
				break
			}
		}
	}))
	server.Handler = m

	log.Info("Listening on", zap.String("addr", addr.String()))

	runner.Start(useFakeData)
	server.ListenAndServe()

	return nil
}

func makeJSArray(vals []float64) string {
	slice := make([]string, len(vals))

	for i, val := range vals {
		slice[i] = fmt.Sprintf("%f", val)
	}

	list := strings.Join(slice, ", ")

	return fmt.Sprintf("[%s]", list)
}
