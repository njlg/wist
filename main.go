package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/urfave/cli"

	"wist/info"
	"wist/service"
)

func main() {
	app := cli.NewApp()

	app.Name = "Wist"
	app.Usage = "Wifi Status"
	app.Version = info.GetVersion() + " (" + runtime.Version() + ") : " + info.GetSHA()

	app.Flags = []cli.Flag{
		// cli.StringFlag{Name: "config", Usage: "Load configuration from `FILE`"},
		cli.UintFlag{Name: "port,p", Usage: "Port to bind HTTP server to. Default 8080.", Value: 8080},
		cli.StringFlag{Name: "templatedir,t", Usage: "Use templates from `DIRECTORY`"},
		cli.BoolFlag{Name: "fakedata,f", Usage: "Use fake data"},
	}

	app.Action = service.Run

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to run app")
		os.Exit(1)
	}
}
