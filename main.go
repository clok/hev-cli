package main

import (
	"encoding/json"
	"fmt"
	"github.com/OrlovEvgeny/go-mcache"
	"github.com/clok/cdocs"
	"github.com/clok/kemba"
	"github.com/umahmood/haversine"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

var (
	version string
	URL     = "https://heb-ecom-covid-vaccine.hebdigital-prd.com/vaccine_locations.json"
	k       = kemba.New("hev")
	client  = &http.Client{Timeout: 10 * time.Second}
)

func getJSON(url string, target interface{}) error {
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		if err = r.Body.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	return json.NewDecoder(r.Body).Decode(target)
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

type location struct {
	Zip                  string
	URL                  string
	SlotDetails          []slotDetail
	OpenTimeSlots        int
	OpenAppointmentSlots int
	Name                 string
	Longitude            float64
	Latitude             float64
	City                 string
}

type slotDetail struct {
	OpenTimeslots        int
	OpenAppointmentSlots int
	Manufacturer         string
}

func main() {
	k.Println("executing")

	im, err := cdocs.InstallManpageCommand(&cdocs.InstallManpageCommandInput{
		AppName: "hev",
	})
	if err != nil {
		log.Fatal(err)
	}

	app := cli.NewApp()
	app.Name = "hev"
	app.Version = version
	app.Usage = "scan H-E-B vaccine availabliby and open a browser when there is one available within a radius of miles"
	app.Commands = []*cli.Command{
		{
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
				kl := k.Extend("watcher")
				radius := c.Int("miles")
				delay := c.Int("delay")
				suppressTTL := c.Int("suppress-ttl")
				kl.Printf("radius: %d miles -- delay: %d seconds -- ttl: %s", radius, delay, suppressTTL)

				origin := haversine.Coord{Lat: c.Float64("latitude"), Lon: c.Float64("longitude")}

				MCache := mcache.New()

				for {
					var data map[string][]*location

					err := getJSON(URL, &data)
					if err != nil {
						return err
					}
					// kl.Log(data)

					var locations []*location
					kll := kl.Extend("location")
					var totalSlots int
					for _, loc := range data["locations"] {
						if loc.URL != "" {
							kll.Log("VACCINES", loc)
							locations = append(locations, loc)
							totalSlots += loc.OpenAppointmentSlots
						}
					}

					kl.Printf("found %d locations with %d total slots", len(locations), totalSlots)

					if totalSlots > 0 {
						for _, loc := range locations {
							_, ok := MCache.Get(loc.Name)

							destination := haversine.Coord{Lat: loc.Latitude, Lon: loc.Longitude}
							mi, _ := haversine.Distance(origin, destination)
							if mi <= float64(radius) {
								fmt.Print("\a")
								fmt.Printf("FOUND %s [%d %s] : %.0f miles", loc.Name, loc.OpenAppointmentSlots, loc.SlotDetails[0].Manufacturer, mi)
								if !ok {
									err := MCache.Set(loc.Name, true, time.Duration(suppressTTL)*time.Minute)
									if err != nil {
										log.Fatal(err)
									}

									fmt.Printf("\nCLICK TO OPEN: %s\n", loc.URL)
									openBrowser(loc.URL)
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
		},
		im,
		{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "Print version info",
			Action: func(c *cli.Context) error {
				fmt.Printf("%s %s (%s/%s)\n", "hev", version, runtime.GOOS, runtime.GOARCH)
				return nil
			},
		},
	}

	if os.Getenv("DOCS_MD") != "" {
		docs, err := cdocs.ToMarkdown(app)
		if err != nil {
			panic(err)
		}
		fmt.Println(docs)
		return
	}

	if os.Getenv("DOCS_MAN") != "" {
		docs, err := cdocs.ToMan(app)
		if err != nil {
			panic(err)
		}
		fmt.Println(docs)
		return
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
