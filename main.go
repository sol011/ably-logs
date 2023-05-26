package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ably/ably-go/ably"
)

type slackMsg struct {
	Text string `json:"text"`
}

func subHandler(c *http.Client, ch *ably.RealtimeChannel, slackUrl string, logChan string) {
	_, err := ch.SubscribeAll(context.Background(), func(msg *ably.Message) {
		msgJson, err := json.Marshal(msg)
		if err != nil {
			log.Printf("%#v", msg)
			return
		}

		msgStr := string(msgJson)
		fmt.Printf("%s, \n", msgStr)

		sMsg := slackMsg{Text: msgStr}
		sMsgJson, err := json.Marshal(sMsg)
		if err != nil {
			log.Printf("could not marshal json %s", err)
			return
		}
		req, err := http.NewRequest(http.MethodPost, slackUrl, bytes.NewReader(sMsgJson))
		if err != nil {
			log.Printf("could not create request to send data to slack %s", err)
			return
		}
		req.Header.Add("Content-type", "application/json")
		_, err = c.Do(req)
		if err != nil {
			log.Printf("could not send data to slack %s", err)
			return
		}
	})
	if err != nil {
		log.Printf("could not subscribe to ably channel %s", logChan)
	}
}

func main() {
	ablyKeyPreprod := os.Getenv("ABLY_KEY_PREPROD")
	ablyKeyProd := os.Getenv("ABLY_KEY_PROD")

	slackPreprodUrl := os.Getenv("SLACK_URL_PREPROD")
	slackProdUrl := os.Getenv("SLACK_URL_PROD")

	logChan := "[meta]log"
	c := &http.Client{}

	var forever chan struct{}

	ablyClientPreprod, err := ably.NewRealtime(ably.WithKey(ablyKeyPreprod))
	if err != nil {
		log.Printf("could not setup ably rest client %s", err)
		return
	}
	preprodCh := ablyClientPreprod.Channels.Get(logChan)
	subHandler(c, preprodCh, slackPreprodUrl, logChan)

	ablyClientProd, err := ably.NewRealtime(ably.WithKey(ablyKeyProd))
	if err != nil {
		log.Printf("could not setup ably rest client %s", err)
		return
	}
	prodCh := ablyClientProd.Channels.Get(logChan)
	subHandler(c, prodCh, slackProdUrl, logChan)

	<-forever
}
