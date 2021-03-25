package commands

import (
	"encoding/json"
	"fmt"
	"github.com/clok/hev-cli/helpers"
	"github.com/clok/hev-cli/types"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"time"
)

var (
	CommandPublish = &cli.Command{
		Name:  "publish",
		Usage: "start the publisher",
		UsageText: `
Poll the HEB Vaccine location API, pushing locations with vaccines available to a REDIS pub/sub channel

NOTE: The REDIS_HOST_URL environment variable is require. Example: redis://redis.yolo.co:6379
`,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "delay",
				Aliases: []string{"d"},
				Usage:   "number of seconds to wait between polling",
				Value:   5,
			},
		},
		Action: func(c *cli.Context) error {
			kl := k.Extend("publish")
			delay := c.Int("delay")
			kl.Printf("delay: %d seconds", delay)

			redisHostURL := os.Getenv("REDIS_HOST_URL")
			if redisHostURL == "" {
				return fmt.Errorf("missing REDIS_HOST_URL env variable")
			}
			kl.Println(redisHostURL)
			creds := helpers.ParseRedisURL(redisHostURL)
			kl.Log(creds)
			pool := helpers.NewRedisPool(fmt.Sprintf("%s:%s", creds.Host, creds.Port), creds.User, creds.Password)

			fmt.Printf("[%s] starting\n", time.Now().UTC())
			go helpers.Heartbeat()

			for {
				locations, err := helpers.GetHEBData()
				if err != nil {
					return err
				}

				if len(locations) > 0 {
					// publish
					kl.Println("publishing")
					for i, location := range locations {
						go func(l *types.Location, c int) {
							conn := pool.Get()
							defer func() {
								if err = conn.Close(); err != nil {
									log.Fatal(err)
								}
							}()
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
							if _, err := conn.Do("PUBLISH", "locations", loc); err != nil {
								log.Fatal(err)
								return
							}
						}(location, i)
					}
				} else {
					kl.Println("No open slots found")
				}

				kl.Printf("-> Checking again in %d seconds ...", delay)
				time.Sleep(time.Duration(delay) * time.Second)
			}
		},
	}
)
