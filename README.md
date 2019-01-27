# Cryobot

A bot for the CRYOSPHERE Discord server

It takes the twitter feed from [cryosphere twitter](https://twitter.com/cryosphereband) and posts the tweets to [discord](https://discord.gg/ejtXgkj)

## Requirements

|software|version   |
|:-------|:---------|
| Golang | 1.10+    |
| Docker | 18.09.1+ |

## Instructions

```
go get github.com/mitchellh/gox
go get github.com/Inveracity/cryobot
cd %GOPATH/src/github.com/Inveracity/cryobot
gox -osarch="linux/amd64"
docker build -t <username>/cryobot:latest .
docker run cryobot
```

## Links

[Twitter](https://twitter.com/cryosphereband)

[Instagram](https://www.instagram.com/cryosphereband/)

[Bandcamp](https://cryosphere.bandcamp.com)

[Discordapp](https://discord.gg/ejtXgkj)

[Youtube](https://www.youtube.com/channel/UCOwnbdRqpukvpcQqmrc6cIQ)

[Facebook](https://www.facebook.com/cryosphereband/)
