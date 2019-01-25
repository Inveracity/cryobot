package main

import (
    "github.com/bwmarrin/discordgo"
    "log"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    log.Println("cryobot starting")

    token := os.Getenv("CRYO_TOKEN")

    discord, err := discordgo.New("Bot " + token)

    if err != nil {
        log.Println("error creating Discord session,", err)
        return
    }

    // Register message handler
    discord.AddHandler(messageCreate)

    // Open socket
    err = discord.Open()
    if err != nil {
        log.Println("Could not open socket to discord")
        log.Println(err)

        return
    }

    // Wait until interrupted
    log.Println("cryobot running, ctrl+x to exit")
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc

    // Close connection to discord
    discord.Close()

}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
    // Ignore messages created by the bot
    if message.Author.ID == session.State.User.ID {
        return
    }

    // If message is "Ping" reply "Pong!"
    if message.Content == "ping" {
        session.ChannelMessageSend(message.ChannelID, "Pong!")
    }

    // Vice versa
    if message.Content == "pong" {
        session.ChannelMessageSend(message.ChannelID, "Ping!")
    }
}
