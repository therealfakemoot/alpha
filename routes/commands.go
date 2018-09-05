package routes

import (
	"github.com/Necroforger/dgrouter/exrouter"
)

// New creates an exrouter with a help message.
func New() *exrouter.Route {
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
