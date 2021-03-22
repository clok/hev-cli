# hev-cli

> H-E-V: Here Everyone's Vaccinated

Please see [the docs for details on the commands.](./docs/hev.md)

- [Docs](./docs/hev.md)
- [About](#about)
- [What it does](#what-it-does)
- [Installation](#installation)
  - [Linux & Mac OS](#linux--mac-os)
  - [Windows](#windows)
- [Usage](#usage)
- [How do I find my latitude and longitude?](#how-do-i-find-my-latitude-and-longitude)
- [Important Links](#important-links)
- [License](#license)

## About

This tool is intended to help those who qualify ([SEE RULES](https://vaccine.heb.com/scheduler)) for the COVID-19 vaccine at H-E-B find an appointment.

Please goto the [COVID-19 vaccines at H-E-B Pharmacy Scheduler](https://vaccine.heb.com/scheduler) to read up on the qualification rules.

## What it does

1. `hev` processes the vaccine schedule API on a regular cadence (every 5 seconds by default).
1. It will check the availability posted for ALL H-EB locations in the list to determine if there are shots available.
1. It will then determine if the location is within a radius of miles based on your lat/long (user provided)
1. If the available shot is within your radius, it will alert you and open a browser directly to the appointment sign up page.

> PLEASE NOTE: Slots are taken up FAST. Once the browser opens, it is very possible that the slot has been taken. Keep trying. Stay persistent. You will get one.

## Installation

### Linux & Mac OS
```
$ curl https://i.jpillora.com/clok/hev-cli! | sed s/PROG=\"hev-cli\"/PROG=\"hev\"/ | bash
```

### Windows

1. Goto [Releases](https://github.com/clok/hev-cli/releases)
1. Download `hev_<version>_windows_amd64.zip`
1. Unzip (extract all)
1. Run the app in a Terminal

## Usage

Run the tool with a 50-mile radius, and a 2-second refresh rate. (The lat/long here is fake)

```
$ hev watch --miles 50 --delay 2 --lat 12.345 --long -12.345
```

Help output

```
$ hev watch --help
NAME:
   hev watch - start the watcher

USAGE:
   hev watch [command options] [arguments...]

OPTIONS:
   --miles value, -m value         radius in miles from location (default: 30)
   --delay value, -d value         number of seconds to wait between polling (default: 5)
   --latitude value, --lat value   origin latitude (default: 30.345122515031083)
   --longitude value, --lon value  origin longitude (default: -97.96755574412511)
   --suppress-ttl value            number of minutes to suppress alerting link to previously seen open slots (default: 5)
   --help, -h                      show help (default: false)
```

## How do I find my latitude and longitude?

Goto [https://map.google.com](https://map.google.com) and enter your address. Right-Click on the Pin marker that shows up. Click the first line, the lat/long. That will copy the value.

![image](https://user-images.githubusercontent.com/1429775/111990513-539ced00-8ae1-11eb-9bcd-c3999933adc1.png)

## Important Links

- [CDC COVID-19 Vaccine Resource Center](https://www.cdc.gov/vaccines/covid-19/index.html)
- [H-E-B Vaccine Scheduler](https://vaccine.heb.com/scheduler)
- [DSHS Texas Guidelines](https://www.dshs.texas.gov/coronavirus/)

## License

[MIT @ clok](LICENSE)

> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
> IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
> FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
> AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
> LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
> OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
> SOFTWARE.
