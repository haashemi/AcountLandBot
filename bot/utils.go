package bot

import (
	"fmt"

	"github.com/haashemi/tgo"
	"github.com/haashemi/tgo/routers/message"
)

func handleError(err error, ctx *message.Context) {
	if err == nil {
		return
	}

	ctx.Send(&tgo.SendMessage{
		Text: fmt.Sprintf(
			"⚠️| <b>ERROR OCCURRED</b>\n\n"+
				"— <code>%s</code>\n\n"+
				"✅| Good luck btw",
			err.Error(),
		),
	})
}

func whitelist(whitelist []int64) message.Middleware {
	return func(ctx *message.Context) (ok bool) {
		for _, id := range whitelist {
			if id == ctx.From.Id {
				return true
			}
		}

		ctx.Send(&tgo.SendMessage{Text: "⚠️| Insufficient Permissions."})

		return false
	}
}
