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
		ctx := Context{
			Iyashi:      iyashi,
			TeamDomain:  r.FormValue("team_domain"),
			ChannelName: r.FormValue("channel_name"),
			Timestamp:   r.FormValue("timestamp"),
			UserID:      r.FormValue("user_id"),
			UserName:    r.FormValue("user_name"),
			Text:        r.FormValue("text"),
			TriggerWord: r.FormValue("trigger_word"),
		}

		var err error
		if _, ok := iyashi.joinChannelMap[ctx.ChannelName]; ok {
			err = next(ctx, w, r)
		}

		w.Header().Set("Content-Type", "text/plain")
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
