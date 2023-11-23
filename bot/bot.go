package bot

import (
	"sync"

	"github.com/LlamaNite/llamalog"
	"github.com/haashemi/AcountLandBot/config"
	"github.com/haashemi/AcountLandBot/fortnite"
	"github.com/haashemi/AcountLandBot/generator"
	"github.com/haashemi/tgo"
	"github.com/haashemi/tgo/filters"
	"github.com/haashemi/tgo/routers/message"
)

type Client struct {
	bot       *tgo.Bot
	log       *llamalog.Logger
	config    *config.Config
	generator *generator.Generator

	shopMut   sync.RWMutex
	shopHash  []byte
	shopItems []fortnite.ItemShopItem
}

func Start(config *config.Config, generator *generator.Generator) {
	log := llamalog.NewLogger("AKB")

	bot := tgo.NewBot(config.Bot.Token, tgo.Options{DefaultParseMode: tgo.ParseModeHTML})
	info, err := bot.GetMe()
	if err != nil {
		log.Fatal("Failed to fetch the bot info > %v", err)
	}

	client := &Client{bot: bot, log: log, config: config, generator: generator}
	ms := message.NewRouter(whitelist(config.Bot.Admins))
	ms.Handle(filters.Command("setpp", info.Username), client.SetPrimaryPrice)
	ms.Handle(filters.Command("setsp", info.Username), client.SetSecondaryPrice)
	ms.Handle(filters.Command("itemshop", info.Username), client.Itemshop)
	bot.AddRouter(ms)

	_, err = bot.SetMyCommands(&tgo.SetMyCommands{Commands: []*tgo.BotCommand{
		{Command: "itemshop", Description: "Receive the current itemshop"},
		{Command: "setpp", Description: "Set primary price"},
		{Command: "setsp", Description: "Set secondary price"},
	}})
	if err != nil {
		log.Fatal("Failed to set bot commands > %v", err)
	}

	go client.trackItemShop()

	log.Info("Polling started as @%s", info.Username)
	log.Fatal("%v", bot.StartPolling(30))
}
