% hev 8
# NAME
hev - scan H-E-B vaccine availability and open a browser when there is one available within a radius of miles
# SYNOPSIS
hev


# COMMAND TREE

- [watch, w](#watch-w)
- [publish](#publish)
- [websocket](#websocket)
- [install-manpage](#install-manpage)
- [version, v](#version-v)

**Usage**:
```
hev [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# COMMANDS

## watch, w

start the watcher

**--delay, -d**="": number of seconds to wait between polling (default: 5)

**--latitude, --lat**="": origin latitude (default: 30.345123)

**--longitude, --lon**="": origin longitude (default: -97.967556)

**--miles, -m**="": radius in miles from location (default: 30)

**--suppress-ttl**="": number of minutes to suppress alerting link to previously seen open slots (default: 5)

## publish

start the publisher

```
Poll the HEB Vaccine location API, pushing locations with vaccines available to a REDIS pub/sub channel

NOTE: The REDIS_HOST_URL environment variable is require. Example: redis://redis.yolo.co:6379
```

**--delay, -d**="": number of seconds to wait between polling (default: 5)

## websocket

start the websocket server

>Poll the HEB Vaccine location API, pushing locations with vaccines available to a websocket

**--addr**="": host address to bind to (default: localhost:8337)

**--delay, -d**="": number of seconds to wait between polling (default: 5)

## install-manpage

Generate and install man page

>NOTE: Windows is not supported

## version, v

Print version info

