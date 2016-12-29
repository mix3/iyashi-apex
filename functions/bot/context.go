package main

import (
	"fmt"
	"strings"

	"github.com/nlopes/slack"
)

var postMessageParameters = slack.PostMessageParameters{
	AsUser:      true,
	UnfurlLinks: true,
	UnfurlMedia: true,
}

type Context struct {
	Iyashi      *Iyashi
	Token       string
	TeamDomain  string
	ChannelName string
	UserID      string
	Timestamp   string
	UserName    string
	Text        string
	TriggerWord string
}

func (c Context) DM(text string) {
	c.Iyashi.api.PostMessage(
		fmt.Sprintf("@%s", c.UserName),
		fmt.Sprintf("<@%s> %s", c.UserID, text),
		postMessageParameters,
	)
}

func (c Context) Reply(text string) {
	c.Iyashi.api.PostMessage(
		c.ChannelName,
		fmt.Sprintf("<@%s> %s %s", c.UserID, text, c.Permalink()),
		postMessageParameters,
	)
}

func (c Context) ReplyWithoutPermalink(text string) {
	c.Iyashi.api.PostMessage(
		c.ChannelName,
		fmt.Sprintf("<@%s> %s", c.UserID, text),
		postMessageParameters,
	)
}

func (c Context) Permalink() string {
	tss := strings.Split(c.Timestamp, ".")
	if len(tss) != 2 {
		return ""
	}
	return fmt.Sprintf(
		"https://%s.slack.com/archives/%s/p%s%s",
		c.TeamDomain,
		c.ChannelName,
		tss[0],
		tss[1],
	)
}
