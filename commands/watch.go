package commands

import (
	"fmt"
	"github.com/OrlovEvgeny/go-mcache"
	"github.com/clok/hev-cli/helpers"
	"github.com/clok/hev-cli/types"
	"github.com/clok/kemba"
	"github.com/umahmood/haversine"
	"github.com/urfave/cli/v2"
	"log"
	"time"
)

var (
	k       = kemba.New("hev:commands")
	special = map[string]types.Location{
		"COVID Vaccination Clinic at Orange Grove ISD": {
			Latitude:  27.95631880379483,
			Longitude: -97.94116047299124,
		},
		"First United Methodist Church": {
			Latitude:  30.04516797336692,
			Longitude: -99.15011125112252,
		},
		"RX04-COMPOUND FRATT & RITTIMAN": {
			Latitude:  29.48710649868841,
			Longitude: -98.3921795306433,
		},
		"COVID Vaccination Clinic at Aransas Pass Civic Center": {
			Latitude:  27.91092331218623,
			Longitude: -97.15047760843393,
		},
		"H-E-B Hosted J&J COVID Vaccination Clinic at Waco Convention Center": {
			Latitude:  31.559946188014,
			Longitude: -97.12905183071285,
		},
	}
)

var (
	CommandWatch = &cli.Command{
		Name:    "watch",
		Aliases: []string{"w"},
		Usage:   "start the watcher",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "miles",
				Aliases: []string{"m"},
				Usage:   "radius in miles from location",
				Value:   30,
			},
			&cli.IntFlag{
				Name:    "delay",
				Aliases: []string{"d"},
				Usage:   "number of seconds to wait between polling",
				Value:   5,
			},
			&cli.Float64Flag{
				Name:    "latitude",
				Aliases: []string{"lat"},
				Usage:   "origin latitude",
				Value:   30.345122515031083,
			},
			&cli.Float64Flag{
				Name:    "longitude",
				Aliases: []string{"lon"},
				Usage:   "origin longitude",
				Value:   -97.96755574412511,
			},
			&cli.IntFlag{
				Name:  "suppress-ttl",
				Usage: "number of minutes to suppress alerting link to previously seen open slots",
				Value: 5,
			},
		},
		Action: func(c *cli.Context) error {
			kl := k.Extend("watch")
			radius := c.Int("miles")
			delay := c.Int("delay")
			suppressTTL := c.Int("suppress-ttl")
			kl.Printf("radius: %d miles -- delay: %d seconds -- ttl: %s", radius, delay, suppressTTL)

			origin := haversine.Coord{Lat: c.Float64("latitude"), Lon: c.Float64("longitude")}

			MCache := mcache.New()

			for {
				locations, err := helpers.GetHEBData()
				if err != nil {
					return err
				}

				if len(locations) > 0 {
					for _, loc := range locations {
						_, ok := MCache.Get(loc.Name)

						var destination haversine.Coord
						var mi float64

						if loc.Latitude != 0 {
							destination = haversine.Coord{Lat: loc.Latitude, Lon: loc.Longitude}
						} else {
							if tmp, ok := special[loc.Name]; ok {
								destination = haversine.Coord{Lat: tmp.Latitude, Lon: tmp.Longitude}
							} else {
								fmt.Printf("No lat/long found for %s", loc.Name)
								break
							}
						}
						mi, _ = haversine.Distance(origin, destination)
						if mi <= float64(radius) {
							fmt.Printf("FOUND %s [%d %s] : %.0f miles", loc.Name, loc.OpenAppointmentSlots, loc.SlotDetails[0].Manufacturer, mi)
							if !ok {
								fmt.Print("\a")
								err := MCache.Set(loc.Name, true, time.Duration(suppressTTL)*time.Minute)
								if err != nil {
									log.Fatal(err)
								}

								fmt.Printf("\nCLICK TO OPEN: %s\n", loc.URL)
								helpers.OpenBrowser(loc.URL)
							} else {
								fmt.Print(" - Already alerting. Silencing for now.")
							}
							fmt.Println("")
						} else {
							fmt.Printf("! %s [%d %s] : %.0f miles\n", loc.Name, loc.OpenAppointmentSlots, loc.SlotDetails[0].Manufacturer, mi)
						}
					}
				} else {
					fmt.Println("No open slots found")
				}

				fmt.Printf("-> Checking again in %d seconds ...", delay)
				time.Sleep(time.Duration(delay) * time.Second)
				fmt.Println("")
			}
		},
	}
)
