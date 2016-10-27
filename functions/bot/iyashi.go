package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/apex/go-apex"
	"github.com/fujiwara/ridge"
	"github.com/nlopes/slack"
)

type Iyashi struct {
	mux                *http.ServeMux
	api                *slack.Client
	slackBotToken      string
	slackOutgoingToken string
	flickrApiToken     string
	tumblrApiToken     string
	host               string
	port               string
	joinChannelMap     map[string]struct{}
	AuthTest           *slack.AuthTestResponse
	dispatchMap        map[string]Command
}

func NewIyashi() (*Iyashi, error) {
	var (
		slackBotToken     = os.Getenv("SLACK_BOT_TOKEN")
		slackOutgoinToken = os.Getenv("SLACK_OUTGOING_TOKEN")
		flickrApiToken    = os.Getenv("FLICKR_API_TOKEN")
		tumblrApiToken    = os.Getenv("TUMBLR_API_TOKEN")
		api               = slack.New(slackBotToken)
		host              = os.Getenv("HOST")
		port              = os.Getenv("PORT")
	)

	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "8080"
	}

	iyashi := &Iyashi{
		mux:                http.NewServeMux(),
		api:                api,
		slackBotToken:      slackBotToken,
		slackOutgoingToken: slackOutgoinToken,
		flickrApiToken:     flickrApiToken,
		tumblrApiToken:     tumblrApiToken,
		host:               host,
		port:               port,
		joinChannelMap:     map[string]struct{}{},
		dispatchMap:        map[string]Command{},
	}

	authTest, err := api.AuthTest()
	if err != nil {
		return nil, err
	}
	iyashi.AuthTest = authTest

	channelsRes, err := api.GetChannels(true)
	if err != nil {
		return nil, err
	}
	for _, channel := range channelsRes {
		for _, member := range channel.Members {
			if member == authTest.UserID {
				iyashi.joinChannelMap[channel.Name] = struct{}{}
			}
		}
	}

	var (
		iyashiCommand = newIyashiCommand()
		moeCommand    = newTumblrCommand(tumblrApiToken, "honobonoarc", "萌え")
		zoiCommand    = newTumblrCommand(tumblrApiToken, "ganbaruzoi", "ぞい")
		tawawaCommand = newTumblrCommand(tumblrApiToken, "tawawa-of-monday", "たわわ")
		helpCommand   = newHelpCommand()
	)
	iyashi.dispatchMap = map[string]Command{
		"癒やし":  iyashiCommand,
		"萌え":   moeCommand,
		"ぞい":   zoiCommand,
		"たわわ":  tawawaCommand,
		"help": helpCommand,
	}

	iyashi.mux.Handle("/", wrapHandle(iyashi, recoverHandle(handleRoot)))

	return iyashi, nil
}

func (i *Iyashi) Run() {
	if os.Getenv("APEX") == "" {
		addr := fmt.Sprintf("%s:%s", i.host, i.port)
		log.Println("starting up with local httpd: ", addr)
		log.Fatal(http.ListenAndServe(addr, i.mux))
	}
	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		r, err := ridge.NewRequest(event)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		r.ParseForm()
		w := ridge.NewResponseWriter()
		i.mux.ServeHTTP(w, r)
		return w.Response(), nil
	})
}
