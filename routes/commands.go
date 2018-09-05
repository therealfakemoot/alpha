package main

import (
	"github.com/Necroforger/dgrouter/exrouter"
	trash "github.com/therealfakemoot/trash-talk"
)

func All() *exrouter.Route {
	r := exrouter.New()

	r.Default = r.On("help", func(ctx *exrouter.Context) {
		var text = ""
		for _, v := range r.Routes {
			text += v.Name + " : \t" + v.Description + "\n"
		}
		ctx.Reply("```" + text + "```")
	}).Desc("prints this help menu")

	r.On("mock", func(ctx *exrouter.Context) {
		msgMap := conf.Get("msgMap").(map[string]*dgo.Message)
		if len(ctx.Msg.Mentions) == 0 {
			ctx.Reply("Who do you want me to make fun of, dumbass?")
			return
		}
		if len(msgMap) == 0 {
			ctx.Reply("Nobody's said anything yet, idiot.")
			return
		}

		target := ctx.Msg.Mentions[0].ID
		targetMsg := msgMap[target].Content
		ctx.Reply(trash.Mock(targetMsg))

	}).Desc("Makes fun of the mentioned user's last message")
}
