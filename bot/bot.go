package bot

import (
	"accountland/global"
	"fmt"
	"gogram"
	"gogram/bindings/td"
	"gogram/filters"
	"os"
	"time"

	"github.com/LlamaNite/llamalog"
)

var log = llamalog.NewLogger("AccountLand")

func Run() {
	checkErr := func(err error) {
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
	}

	client := gogram.New(
		global.Config.TelegramBot.APIID,
		global.Config.TelegramBot.APIHash,
		global.Config.TelegramBot.CachePath,
	)
	checkErr(client.AuthenticateBot(global.Config.TelegramBot.BotToken))
	info, err := client.GetRawClient().GetMe()
	checkErr(err)
	info.Username = "@" + info.Username

	tracker = rawTracker{client: client}
	go tracker.Start()
	registerCommands(client, info.Username)

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

func onError(err error, msg *gogram.Message, c *gogram.Client) {
	if err == nil {
		return
	}

	msg.Send(c.HTML(fmt.Sprintf(
		"⚠️| <b>ERROR OCCURED</b>\n\n"+
			"— <code>%s</code>\n\n"+
			"✅| Goodluck btw",
		err.Error(),
	)))
}
