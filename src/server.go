package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	Time "time"

	"github.com/kyokomi/goslash/goslash"
	"github.com/kyokomi/goslash/plugins"
	"github.com/kyokomi/goslash/plugins/akari"
	"github.com/kyokomi/goslash/plugins/echo"
	"github.com/kyokomi/goslash/plugins/lgtm"
	"github.com/kyokomi/goslash/plugins/suddendeath"
	"github.com/kyokomi/goslash/plugins/time"

	"github.com/unrolled/render"
)

// option flags
var (
	port    uint
	timeout Time.Duration
	isDev   bool
)

func main() {
	flag.UintVar(&port, "port", 8080, "server port")
	flag.BoolVar(&isDev, "dev", false, "is development mode")
	flag.DurationVar(&timeout, "timeout", 10*Time.Second, "http client timeout")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n  %s [OPTIONS]\n", os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	renderer := render.New(render.Options{})

	http.HandleFunc("/v1/cmd", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)

		req, err := goslash.ParseFormSlashCommandRequest(r)
		if err != nil {
			renderer.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		client := &http.Client{
			Timeout: timeout * Time.Second,
		}
		slashPlugins := map[string]plugins.Plugin{
			"echo":  echo.New(),
			"time":  time.New(),
			"突然":    suddendeath.New(),
			"LGTM":  lgtm.New(client),
			"akari": akari.New(),
		}

		slashCmd := plugins.New(client, slashPlugins)

		if isDev {
			// development
			cmd, _ := req.CmdArgs()
			p, ok := slashPlugins[cmd]
			if !ok {
				renderer.JSON(w, http.StatusNotFound, "cmd not found")
				return
			}

			msg := p.Do(req)
			var jsonData bytes.Buffer
			if err := json.NewEncoder(&jsonData).Encode(&msg); err != nil {
				renderer.JSON(w, http.StatusInternalServerError, err.Error())
				return
			}
			renderer.JSON(w, http.StatusOK, jsonData.String())
		} else {
			// production
			renderer.Text(w, http.StatusOK, slashCmd.Execute(req))
		}
	})

	log.Println("Start Server. Port:", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
