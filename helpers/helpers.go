package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/clok/hev-cli/types"
	"github.com/clok/kemba"
	"github.com/gomodule/redigo/redis"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"time"
)

var (
	URL           = "https://heb-ecom-covid-vaccine.hebdigital-prd.com/vaccine_locations.json"
	redisURLRegex = regexp.MustCompile(`^rediss?:\/\/((?P<user>\w+)?:?(?P<password>\w+)?@)?(?P<host>[\w\-\.]+):(?P<port>\d+)`)
	client        = &http.Client{Timeout: 10 * time.Second}
	k             = kemba.New("hev:helpers")
	kgd           = k.Extend("getHEBData")
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

func OpenBrowser(url string) {
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

func GetHEBData() (locations []*types.Location, err error) {
	var data map[string][]*types.Location

	err = getJSON(URL, &data)
	if err != nil {
		return nil, err
	}

	var totalSlots int
	for _, loc := range data["locations"] {
		if loc.URL != "" {
			kgd.Log("VACCINES", loc)
			locations = append(locations, loc)
			totalSlots += loc.OpenAppointmentSlots
		}
	}
	kgd.Printf("found %d locations with %d total slots", len(locations), totalSlots)

	return locations, err
}

func Heartbeat() {
	for {
		timer := time.After(time.Second * 60)
		<-timer
		fmt.Printf("[%s] heartbeat\n", time.Now().UTC())
	}
}

func ParseRedisURL(url string) *types.RedisCredentials {
	match := redisURLRegex.FindStringSubmatch(url)
	k.Log(match)
	return &types.RedisCredentials{
		URL:      url,
		User:     match[2],
		Password: match[3],
		Host:     match[4],
		Port:     match[5],
	}
}

func NewRedisPool(addr string, username string, password string) *redis.Pool {
	mode := os.Getenv("ENV")
	if mode == "production" {
		return &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				return redis.Dial(
					"tcp",
					addr,
					redis.DialUsername(username),
					redis.DialPassword(password),
					redis.DialTLSSkipVerify(true),
					redis.DialUseTLS(true),
				)
			},
		}
	}

	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				addr,
			)
		},
	}
}
