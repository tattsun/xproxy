package main

import (
	"log"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "xproxy"
	app.Usage = "yet another pac file killer"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "profiler",
			Usage: "profiler will be started on port 1081",
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.Bool("profiler") {
			go func() {
				r := http.NewServeMux()
				r.HandleFunc("/debug/pprof/", pprof.Index)
				r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
				r.HandleFunc("/debug/pprof/profile", pprof.Profile)
				r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
				r.HandleFunc("/debug/pprof/trace", pprof.Trace)
				http.ListenAndServe("127.0.0.1:1081", r)
			}()
		}

		configFile, err := os.Open("./config.yml")
		if err != nil {
			return err
		}

		config, err := ParseConfig(configFile)
		if err != nil {
			return err
		}

		server, err := NewServer(config.Host, config.Port, config)
		if err != nil {
			return err
		}

		return server.Start()
	}

	log.Fatal(app.Run(os.Args))
}
