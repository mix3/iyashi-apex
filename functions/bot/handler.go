package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Handler func(ctx Context, w http.ResponseWriter, r *http.Request) error

func wrapHandle(iyashi *Iyashi, next Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		ctx := Context{
			Iyashi:      iyashi,
			Token:       r.PostFormValue("token"),
			TeamDomain:  r.PostFormValue("team_domain"),
			ChannelName: r.PostFormValue("channel_name"),
			Timestamp:   r.PostFormValue("timestamp"),
			UserID:      r.PostFormValue("user_id"),
			UserName:    r.PostFormValue("user_name"),
			Text:        r.PostFormValue("text"),
			TriggerWord: r.PostFormValue("trigger_word"),
		}

		// token check
		if iyashi.slackOutgoingToken != ctx.Token {
			fmt.Fprintln(w, "token error")
			return
		}

		// メンションされたチャンネルに所属している
		// かつ、自分からのメンションではない
		var err error
		_, ok := iyashi.joinChannelMap[ctx.ChannelName]
		if ok && iyashi.AuthTest.User != ctx.UserName {
			err = next(ctx, w, r)
		}

		if err != nil {
			ctx.Reply(fmt.Sprintf("%v", err))
			fmt.Fprintln(w, "ng")
		} else {
			fmt.Fprintln(w, "ok")
		}
	}
	return http.HandlerFunc(f)
}

func recoverHandle(next Handler) Handler {
	f := func(ctx Context, w http.ResponseWriter, r *http.Request) (err error) {
		defer func() {
			if e := recover(); e != nil {
				err = fmt.Errorf("panic:%v", e)
			}
		}()
		err = next(ctx, w, r)
		return err
	}
	return f
}

func handleRoot(ctx Context, w http.ResponseWriter, r *http.Request) error {
	trimText := strings.Trim(strings.TrimLeft(ctx.Text, ctx.TriggerWord), " ")
	if trimText == "" {
		log.Println("empty command")
		ctx.Reply("???")
		return nil
	}

	words := strings.Split(trimText, " ")

	cmd, ok := ctx.Iyashi.dispatchMap[words[0]]
	if !ok {
		log.Println("unknown command:", words[0])
		ctx.Reply("???")
		return nil
	}

	if err := cmd.Func(
		ctx, words[0], words[1:],
	); err != nil {
		return err
	}

	return nil
}
