package main

import (
    "cryobot/discord"
    "cryobot/twitter"
    "cryobot/config"
    "log"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    cfg := config.ReadConfig()
    tweetsChannel := make(chan string)                       // Tweets passing from twitter to discord
    closeDiscord := make(chan string)                        // Signal to stop discord
    closeTwitter := make(chan string)                        // Signal to stop twitter
    go discord.Discordfeed(tweetsChannel, closeDiscord)      // Connect to Discord
    go twitter.Twitterfeed(tweetsChannel, closeTwitter, cfg) // Connect to Twitter

    sc := make(chan os.Signal, 1) // debug feature, ctrl+x to exit
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc
    closeTwitter <- "close"
    closeDiscord <- "close"
    log.Println("Program Terminated")
}
