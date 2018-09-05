package routes

import (
	"github.com/Necroforger/dgrouter/exrouter"
	dgo "github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
	trash "github.com/therealfakemoot/trash-talk"
)

// Mock is a HOF that returns an exrouter HandlerFunc.
func Mock(conf *viper.Viper) func(*exrouter.Context) {
	return func(ctx *exrouter.Context) {
		msgMap := conf.Get("msgMap").(map[string]*dgo.Message)
		if len(ctx.Msg.Mentions) == 0 {
			ctx.Reply("Who do you want me to make fun of, dumbass?")
			return
		}

		target := ctx.Msg.Mentions[0].ID
		targetMsg, ok := msgMap[target]
		if !ok {
			ctx.Reply("Try again, chucklefuck.")
		}

		ctx.Reply(trash.Mock(targetMsg.Content))

	}
}
