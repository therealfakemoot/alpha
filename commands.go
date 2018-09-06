package main

import (
	"github.com/Necroforger/dgrouter/exrouter"
	dgo "github.com/bwmarrin/discordgo"
	trash "github.com/therealfakemoot/trash-talk"
)

// Build creates an exrouter with a help message.
func NewRoute() *exrouter.Route {
	r := exrouter.New()

	r.Default = r.On("help", func(ctx *exrouter.Context) {
		var text = ""
		for _, v := range r.Routes {
			text += v.Name + " : \t" + v.Description + "\n"
		}
		ctx.Reply("```" + text + "```")
	}).Desc("prints this help menu")
	return r
}

// Mock is a HOF that returns an exrouter HandlerFunc.
func Mock(conf Conf) func(*exrouter.Context) {
	return func(ctx *exrouter.Context) {
		msgMap := conf.State["msgMap"].(map[string]*dgo.Message)
		if len(ctx.Msg.Mentions) == 0 {
			ctx.Reply("Who do you want me to make fun of, dumbass?")
			return
		}

		target := ctx.Msg.Mentions[0].ID
		targetMsg, ok := msgMap[target]
		if !ok {
			ctx.Reply("Try again, chucklefuck.")
			return
		}

		ctx.Reply(trash.Mock(targetMsg.Content))

	}
}
