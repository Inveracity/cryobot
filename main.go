package main

import (
    "cryobot/discord"
    "cryobot/twitter"
    "log"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    c := make(chan string)
    done := make(chan string)
    go discord.Discordfeed(c, done)
    go twitter.Twitterfeed(c, done)

    log.Println("ctrl+x to exit")
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc
    done <- "close"
    log.Println("closed")
}
