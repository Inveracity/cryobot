package main

import (
    "flag"
    "github.com/bwmarrin/discordgo"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/coreos/pkg/flagutil"
    "github.com/dghubble/go-twitter/twitter"
    "github.com/dghubble/oauth1"
)

func main() {
    c := make(chan string)
    done := make(chan string)
    go discordfeed(c, done)
    go twitterfeed(c, done)

    time.Sleep(2 * time.Second)
    log.Println("testing a signal")
    c <- "passing a message to the channel"

    log.Println("ctrl+x to exit")
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc
    done <- "close"
    log.Println("closed")

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

func discordfeed(c chan string, done chan string) {
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
                log.Println("discord loop has no messages from twitter")
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

            time.Sleep(1 * time.Second)
        }
    }()
}

func twitterfeed(c chan string, done chan string) {
    flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
    consumerKey := flags.String("consumer-key", "", "Twitter Consumer Key")
    consumerSecret := flags.String("consumer-secret", "", "Twitter Consumer Secret")
    accessToken := flags.String("access-token", "", "Twitter Access Token")
    accessSecret := flags.String("access-secret", "", "Twitter Access Secret")
    flags.Parse(os.Args[1:])
    flagutil.SetFlagsFromEnv(flags, "TWITTER")

    if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" {
        log.Fatal("Consumer key/secret and Access token/secret required")
    }

    config := oauth1.NewConfig(*consumerKey, *consumerSecret)
    token := oauth1.NewToken(*accessToken, *accessSecret)
    // OAuth1 http.Client will automatically authorize Requests
    httpClient := config.Client(oauth1.NoContext, token)

    // Twitter Client
    client := twitter.NewClient(httpClient)

    // Convenience Demux demultiplexed stream messages
    demux := twitter.NewSwitchDemux()
    demux.Tweet = func(tweet *twitter.Tweet) {
        message := "https://twitter.com/" + tweet.User.ScreenName + "/status/" + tweet.IDStr
        log.Println(message)
        c <- message
    }
    demux.DM = func(dm *twitter.DirectMessage) {
        log.Println(dm.SenderID)
    }

    log.Println("Twitterfeed starting")

    // FILTER
    filterParams := &twitter.StreamFilterParams{
        //Track:         []string{"testing"},
        Follow:        []string{"1050383420960980998"}, // cryosphereband
        StallWarnings: twitter.Bool(true),
    }
    stream, err := client.Streams.Filter(filterParams)
    if err != nil {
        log.Fatal(err)
    }

    // Receive messages until stopped or stream quits
    go demux.HandleChan(stream.Messages)

    go func() {
        for {
            select {
            case donesignal := <-done:
                log.Println(donesignal)
                log.Println("Closing twitter connection")
                stream.Stop() // Close twitter connection if the "close" message has been received
                break
            default:
                time.Sleep(1 * time.Second)
            }
        }
    }()
}
