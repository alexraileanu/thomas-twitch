package main

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "net/http"
)

// these 2 functions verify the webhook incoming from twitch
// as per their documentation at
// https://dev.twitch.tv/docs/eventsub/handling-webhook-events#verifying-the-event-message

// getHmacMsg composes the message to be verified against the provided hash based on
// Message-Id and Message-Timestamp headers and the body of the request
func getHmacMsg(request *http.Request, msgBody []byte) string {
    msgId := request.Header.Get(TwitchMessageIdHeader)
    msgTimestamp := request.Header.Get(TwitchMessageTimestampHeader)

    return msgId + msgTimestamp + string(msgBody)
}

func verifyHmac(msg, key []byte, hash string) bool {
    sig, err := hex.DecodeString(hash)
    if err != nil {
        return false
    }
    mac := hmac.New(sha256.New, key)
    mac.Write(msg)

    return hmac.Equal(sig, mac.Sum(nil))
}
