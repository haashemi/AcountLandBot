package bot

import (
	"accountland/global"
	"gogram"
	"gogram/filters"
	"strconv"
	"strings"

	"golang.org/x/text/message"
)

func registerCommands(client *gogram.Client, usrename string) {
	client.OnNewMessage(filters.And(filters.IsPrivate(), filters.Or(isSuperUser(), isAdmin()), filters.Command("/", usrename, "setPrice")), onSetPrice)
	client.OnNewMessage(filters.And(filters.IsPrivate(), filters.Or(isSuperUser(), isAdmin()), filters.Command("/", usrename, "setPriceIllegal")), onSetPriceIllegal)
	client.OnNewMessage(filters.And(filters.IsPrivate(), filters.Or(isSuperUser(), isAdmin()), filters.Command("/", usrename, "forcePost")), onForcePost)
}

func onSetPrice(msg *gogram.Message, c *gogram.Client) {
	caption := strings.SplitN(msg.GetCaption(), " ", 2)
	if len(caption) != 2 {
		msg.Send(c.HTML(
			"✅| <b>Set Price command</b>\n\n" +
				"— <b>Arguments:</b>\n" +
				"—— <i>Price/V-Bucks</i> <code>(number)</code>\n\n" +
				"— <b>Example:</b>\n" +
				"—— <code>/setPrice 120</code>\n\n" +
				"🔥| <i>Developed by Seyyed</i>",
		))
		return
	}

	newPrice, err := strconv.Atoi(caption[1])
	if err != nil {
		onError(err, msg, c)
	}

	global.Config.Itemshop.PriceLegal = newPrice
	if err = global.Config.UpdateConfig(); err != nil {
		onError(err, msg, c)
	}

	msg.Send(c.HTML(message.NewPrinter(message.MatchLanguage("en")).Sprintf(
		"✅| <b>PRICES UPDATED</b>\n\n"+
			"— V-Bucks — Tomans\n"+
			"💵| <code>200V — %dT</code>\n"+
			"💵| <code>500V — %dT</code>\n"+
			"💵| <code>800V — %dT</code>\n"+
			"💵| <code>1,200V — %dT</code>\n"+
			"💵| <code>1,500V — %dT</code>\n"+
			"💵| <code>1,600V — %dT</code>\n"+
			"💵| <code>2,000V — %dT</code>\n\n"+
			"🔥| <i>Developed by Seyyed</i>",
		200*newPrice,
		500*newPrice,
		800*newPrice,
		1200*newPrice,
		1500*newPrice,
		1600*newPrice,
		2000*newPrice,
	)))
}

func onSetPriceIllegal(msg *gogram.Message, c *gogram.Client) {
	caption := strings.SplitN(msg.GetCaption(), " ", 2)
	if len(caption) != 2 {
		msg.Send(c.HTML(
			"✅| <b>Set Nime Legal Price command</b>\n\n" +
				"— <b>Arguments:</b>\n" +
				"—— <i>Price/V-Bucks</i> <code>(number)</code>\n\n" +
				"— <b>Example:</b>\n" +
				"—— <code>/setPrice 120</code>\n\n" +
				"🔥| <i>Developed by Seyyed</i>",
		))
		return
	}

	newPrice, err := strconv.Atoi(caption[1])
	if err != nil {
		onError(err, msg, c)
	}

	global.Config.Itemshop.PriceIllegal = newPrice
	if err = global.Config.UpdateConfig(); err != nil {
		onError(err, msg, c)
	}

	msg.Send(c.HTML(message.NewPrinter(message.MatchLanguage("en")).Sprintf(
		"✅| <b>Nime Legal PRICES UPDATED</b>\n\n"+
			"— V-Bucks — Tomans\n"+
			"💵| <code>200V — %dT</code>\n"+
			"💵| <code>500V — %dT</code>\n"+
			"💵| <code>800V — %dT</code>\n"+
			"💵| <code>1,200V — %dT</code>\n"+
			"💵| <code>1,500V — %dT</code>\n"+
			"💵| <code>1,600V — %dT</code>\n"+
			"💵| <code>2,000V — %dT</code>\n\n"+
			"🔥| <i>Developed by Seyyed</i>",
		200*newPrice,
		500*newPrice,
		800*newPrice,
		1200*newPrice,
		1500*newPrice,
		1600*newPrice,
		2000*newPrice,
	)))
}

func onForcePost(msg *gogram.Message, c *gogram.Client) {
	caption := strings.SplitN(msg.GetCaption(), " ", 2)
	if len(caption) != 2 {
		msg.Send(c.HTML(
			"✅| <b>Force Post command</b>\n\n" +
				"— <b>Arguments:</b>\n" +
				"—— <i>Target Chat</i>: <code>here / channel</code>\n\n" +
				"— <b>Example:</b>\n" +
				"—— <code>/forcePost here</code>\n" +
				"—— <code>/forcePost channel</code>\n\n" +
				"🔥| <i>Developed by Seyyed</i>",
		))
		return
	}

	switch strings.ToLower(caption[1]) {
	case "here":
		tracker.CheckForUpdates(true, msg.SenderID())
	case "channel":
		tracker.CheckForUpdates(true, global.Config.Itemshop.Channel)
	}
}
