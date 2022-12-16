package bot

import (
	"accountland/global"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/LlamaNite/llamalog"
	"github.com/er-azh/gogram"
	"github.com/er-azh/gogram/bindings/td"
	"github.com/er-azh/gogram/filters"
)

var log = llamalog.NewLogger("AccountLand")

func Run() {
	checkErr := func(err error) {
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
	}

	client := gogram.New()

	checkErr(client.AuthorizeBot(
		context.Background(),
		global.Config.TelegramBot.APIID,
		global.Config.TelegramBot.APIHash,
		global.Config.TelegramBot.BotToken,
		global.Config.TelegramBot.CachePath,
	))
	info, err := client.RawClient().GetMe(context.Background())
	checkErr(err)
	info.Username = "@" + info.Username

	tracker = rawTracker{client: client}
	go tracker.Start()
	registerCommands(client, info.Username)

	// Get channel to make sure about the config file
	client.RawClient().GetChat(context.Background(), global.Config.Itemshop.Channel)
	log.Info("Chat %d fetched", global.Config.Itemshop.Channel)

	// Start polling with ignoring updates.
	time.Sleep(time.Second)
	log.Info("Bot is online as %s", info.Username)
	client.StartPolling()
}

func isSuperUser() *filters.FuncFilter {
	return filters.NewFuncFilter(func(update td.Update) bool {
		if data, ok := update.(*td.UpdateNewMessage); ok {
			if sender, ok := data.Message.SenderID.(*td.MessageSenderUser); ok {
				for _, ownerID := range global.Config.TelegramBot.SuperUsers {
					if sender.UserID == ownerID {
						return true
					}
				}
			}
		}
		return false
	})
}

func isAdmin() *filters.FuncFilter {
	return filters.NewFuncFilter(func(update td.Update) bool {
		if data, ok := update.(*td.UpdateNewMessage); ok {
			if sender, ok := data.Message.SenderID.(*td.MessageSenderUser); ok {
				for _, ownerID := range global.Config.TelegramBot.Admins {
					if sender.UserID == ownerID {
						return true
					}
				}
			}
		}
		return false
	})
}

func onError(err error, ctx gogram.MessageContext) {
	if err == nil {
		return
	}

	ctx.Send(gogram.Text(fmt.Sprintf(
		"⚠️| <b>ERROR OCCURED</b>\n\n"+
			"— <code>%s</code>\n\n"+
			"✅| Goodluck btw",
		err.Error(),
	), gogram.TextHTML))
}
