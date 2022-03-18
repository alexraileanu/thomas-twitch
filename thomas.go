package main

import (
    "encoding/json"
    "fmt"
    "github.com/bwmarrin/discordgo"
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strings"
)

const (
    TwitchMessageIdHeader        = "Twitch-Eventsub-Message-Id"
    TwitchMessageTimestampHeader = "Twitch-Eventsub-Message-Timestamp"
    TwitchMessageSignatureHeader = "Twitch-Eventsub-Message-Signature"
    TwitchMessageTypeHeader      = "Twitch-Eventsub-Message-Type"

    TwitchStreamOnlineType = "stream.online"

    MessageTypeVerification = "webhook_callback_verification"
    MessageTypeNotification = "notification"
)

type Notification struct {
    Challenge    string
    Subscription struct {
        Id        string
        Status    string
        Type      string
        Version   string
        Condition struct {
            BroadcasterUserId string
        }
        Transport struct {
            Method   string
            Callback string
        }
    }
    Event struct {
        BroadcasterUserName string `json:"broadcaster_user_name"`
    }
}

func main() {
    _ = godotenv.Load()

    r := mux.NewRouter()

    r.HandleFunc("/twitch/alert", func(writer http.ResponseWriter, request *http.Request) {
        // there needs to be a copy of the request body as it gets consumed the first time it is read
        // meaning that when attempting to unmarshal the json body below it will fail
        bodyCopy := request.Body
        reqBody, _ := ioutil.ReadAll(bodyCopy)
        secret := os.Getenv("SUB_SECRET")
        msg := getHmacMsg(request, reqBody)
        hash := strings.SplitN(request.Header.Get(TwitchMessageSignatureHeader), "=", 2)

        isValid := verifyHmac([]byte(msg), []byte(secret), hash[1])
        // only handle the incoming request if the hmac verification is valid. otherwise we don't care
        if isValid {
            r := new(Notification)
            err := json.Unmarshal(reqBody, r)
            if err != nil {
                writer.WriteHeader(http.StatusBadRequest)
                fmt.Fprintf(writer, "error")
                return
            }

            msgType := request.Header.Get(TwitchMessageTypeHeader)
            // we only care about 2 kinds of incoming message types:
            // - verification messages (only emitted once after subscribing to an event)
            // - notification messages (emitted every time the event (stream.online) happens)
            switch msgType {
            case MessageTypeVerification:
                // as per the twitch documentation, if the type is of webhook_callback_verification
                // the response needs to be the challenge field passed in the request in string format
                // note: no extra headers passed or anything as it will result in a failed verification
                writer.WriteHeader(http.StatusOK)
                fmt.Fprintf(writer, r.Challenge)
                return
            case MessageTypeNotification:
                if r.Subscription.Type == TwitchStreamOnlineType {
                    // if the message is of type stream.online, a discord bot is created which then sends a basic
                    // message to an anouncement channel defined in the .env file
                    message := fmt.Sprintf("%[1]s is LIVE at https://twitch.tv/%[1]s", r.Event.BroadcasterUserName)

                    thomas, _ := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
                    thomas.ChannelMessageSend(os.Getenv("DISCORD_CHANNEL_ID"), message)
                }
            }
        }
    })

    log.Fatal(http.ListenAndServe(":8080", r))
}
