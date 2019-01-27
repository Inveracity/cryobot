package twitter

import (
    "flag"
    "log"
    "os"
    "time"

    "github.com/coreos/pkg/flagutil"
    "github.com/dghubble/go-twitter/twitter"
    "github.com/dghubble/oauth1"
)

// Twitterfeed fetches tweets and passes them back through channel "discord"
func Twitterfeed(discord chan string, closeTwitter chan string) {

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
    httpClient := config.Client(oauth1.NoContext, token) // OAuth1 http.Client will automatically authorize Requests
    client := twitter.NewClient(httpClient)

    // Convenience Demux demultiplexed stream messages
    demux := twitter.NewSwitchDemux()
    demux.Tweet = func(tweet *twitter.Tweet) {
        message := "https://twitter.com/" + tweet.User.ScreenName + "/status/" + tweet.IDStr // Get link to tweet
        discord <- message                                                                   // Pass tweet back to discord
    }
    demux.DM = func(dm *twitter.DirectMessage) {
        log.Println(dm.SenderID)
    }

    log.Println("Opening twitter connection")

    // FILTER
    filterParams := &twitter.StreamFilterParams{
        Track:         []string{"cryosphereband"},
        Follow:        []string{"1050383420960980998"}, // cryosphereband userid
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
            case <-closeTwitter:
                log.Println("Closing twitter connection")
                stream.Stop() // Close twitter connection if the "close" message has been received
                close(closeTwitter)
                return
            default:
                time.Sleep(1 * time.Second)
            }
        }
        return
    }()
}
