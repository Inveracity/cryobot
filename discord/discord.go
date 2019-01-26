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

// Discordfeed connects to discord and receives twitter messages from channel "c"
func Discordfeed(c chan string, done chan string) {
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

    // Listen for messages from twitter
    go func() {
        for {
            select {
            case message := <-c:
                discord.ChannelMessageSend("538683678708989952", message) // #twitter channel
            default:
                time.Sleep(1 * time.Second)
            }

            select {
            case donesignal := <-done:
                log.Println(donesignal)
                log.Println("Closing discord connection")
                discord.Close() // Close discord connection if the "close" message has been received
                break
            default:
                time.Sleep(1 * time.Second)
            }
        }
    }()
}
