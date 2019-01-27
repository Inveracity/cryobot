package discord

import (
    "github.com/bwmarrin/discordgo"
    "log"
    "os"
    "time"
)

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

// Discordfeed connects to discord and receives twitter messages to post
func Discordfeed(twitter chan string, closeDiscord chan string) {
    token := os.Getenv("CRYO_TOKEN")
    discord, err := discordgo.New("Bot " + token)
    if err != nil {
        log.Println("error creating Discord session, ", err)
        return
    }

    // Register message handler
    discord.AddHandler(messageCreate)

    // Open socket
    log.Println("Opening discord connection")
    err = discord.Open()
    if err != nil {
        log.Println("Could not open socket to discord, ", err)
        return
    }

    // Listen for messages from twitter
    go func() {
        for {
            select {
            case tweet := <-twitter:
                discord.ChannelMessageSend("538683678708989952", tweet) // sends tweet to discord channel "#twitter"
            default:
                time.Sleep(1 * time.Second)
            }

            select {
            case <-closeDiscord: // If anything is received on this channel, disconnect from discord
                log.Println("Closing discord connection")
                discord.Close()
                close(twitter)
                close(closeDiscord)
                return
            default:
                time.Sleep(1 * time.Second)
            }
        }
        return
    }()
}
