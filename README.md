# hev-cli

> H-E-V: Here Everyone's Vaccinated

Please see [the docs for details on the commands.](./docs/hev.md)

- [Docs](./docs/hev.md)
- [About](#about)
- [What it does](#what-it-does)
- [Installation](#installation)
  - [Linux & Mac OS](#linux--mac-os)
  - [Windows](#windows)
  - [docker](#docker)
- [Usage](#usage)
- [How do I find my latitude and longitude?](#how-do-i-find-my-latitude-and-longitude)
- [Important Links](#important-links)
- [License](#license)

## About

This tool is intended to help those who qualify ([SEE RULES](https://vaccine.heb.com/scheduler)) for the COVID-19 vaccine at H-E-B find an appointment.

Please goto the [COVID-19 vaccines at H-E-B Pharmacy Scheduler](https://vaccine.heb.com/scheduler) to read up on the qualification rules.

## What it does

1. `hev` processes the vaccine schedule API on a regular cadence (every 5 seconds by default).
1. It will check the availability posted for ALL H-E-B locations in the list to determine if there are shots available.
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

### [docker](https://www.docker.com/)

The compiled docker images are maintained on [GitHub Container Registry (ghcr.io)](https://github.com/users/clok/packages/container/package/hev-cli).
We maintain the following tags:

- `edge`: Image that is build from the current `HEAD` of the main line branch.
- `latest`: Image that is built from the [latest released version](https://github.com/clok/hev-cli/releases/releases)
- `x.y.z` (versions): Images that are build from the tagged versions within Github.

```bash
docker pull ghcr.io/clok/hev-cli
docker run -v "$PWD":/workdir ghcr.io/clok/hev-cli --version
```

## Usage

Run the tool with a 50-mile radius, and a 2-second refresh rate. (The lat/long here is fake)

```
$ hev watch --miles 50 --delay 2 --lat 12.345 --long -12.345
```

Help output

```
NAME:
   hev - scan H-E-B vaccine availability and open a browser when there is one available within a radius of miles

USAGE:
   hev [global options] command [command options] [arguments...]

VERSION:
   v0.3.0

COMMANDS:
   watch, w         start the watcher
   publish          start the publisher
   websocket        start the websocket server
   install-manpage  Generate and install man page
   version, v       Print version info
   help, h          Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
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
