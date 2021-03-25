package commands

import (
	"encoding/json"
	"fmt"
	"github.com/clok/hev-cli/helpers"
	"github.com/clok/hev-cli/types"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
	"time"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "commands/home.html")
}

func scrapeHEB(delay int, hub *helpers.Hub) {
	kl := k.Extend("scrapeHEB")
	for {
		locations, err := helpers.GetHEBData()
		if err != nil {
			log.Fatal(err)
		}

		if len(locations) > 0 {
			// publish
			kl.Println("found")
			for i, location := range locations {
				go func(l *types.Location, c int) {
					kl.Printf("%d\tpublishing: %s", c, l.Name)
					if l.Latitude == 0 {
						if tmp, ok := special[l.Name]; ok {
							l.Latitude = tmp.Latitude
							l.Longitude = tmp.Longitude
						} else {
							kl.Printf("\tNo lat/long found for %s", l.Name)
						}
					}

					loc, err := json.Marshal(l)
					if err != nil {
						log.Fatal(err)
						return
					}
					hub.Broadcast(loc)
				}(location, i)
			}
		} else {
			kl.Println("No open slots found")
		}

		kl.Printf("-> Checking again in %d seconds ...", delay)
		time.Sleep(time.Duration(delay) * time.Second)
	}
}

var (
	CommandWebsocket = &cli.Command{
		Name:  "websocket",
		Usage: "start the websocket server",
		UsageText: `
Poll the HEB Vaccine location API, pushing locations with vaccines available to a websocket
`,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "delay",
				Aliases: []string{"d"},
				Usage:   "number of seconds to wait between polling",
				Value:   5,
			},
			&cli.StringFlag{
				Name:        "addr",
				EnvVars: []string{"WEBSOCKET_HOST"},
				Value: "localhost:8337",
				Usage: "host address to bind to",
			},
		},
		Action: func(c *cli.Context) error {
			kl := k.Extend("websocket")
			delay := c.Int("delay")
			kl.Printf("delay: %d seconds", delay)

			fmt.Printf("[%s] starting\n", time.Now().UTC())
			go helpers.Heartbeat()

			hub := helpers.NewHub()
			go hub.Run()

			go scrapeHEB(delay, hub)

			http.HandleFunc("/", serveHome)
			http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
				helpers.ServeWs(hub, w, r)
			})
			err := http.ListenAndServe(c.String("addr"), nil)
			if err != nil {
				log.Fatal("ListenAndServe: ", err)
			}
			return nil
		},
	}
)
