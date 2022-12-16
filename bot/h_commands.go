package bot

import (
	"accountland/global"
	"strconv"
	"strings"

	"github.com/er-azh/gogram"
	"github.com/er-azh/gogram/filters"
	"golang.org/x/text/message"
)

func registerCommands(client *gogram.Client, username string) {
	client.OnNewMessage(filters.And(filters.IsPrivate(), filters.Or(isSuperUser(), isAdmin()), filters.Command("/", username, "setPrice")), onSetPrice)
	client.OnNewMessage(filters.And(filters.IsPrivate(), filters.Or(isSuperUser(), isAdmin()), filters.Command("/", username, "setPriceIllegal")), onSetPriceIllegal)
	client.OnNewMessage(filters.And(filters.IsPrivate(), filters.Command("/", username, "forcePost")), onForcePost)
}

func onSetPrice(ctx gogram.MessageContext) {
	caption := strings.SplitN(ctx.GetCaption(), " ", 2)
	if len(caption) != 2 {
		ctx.Send(gogram.Text(
			"✅| <b>Set Price command</b>\n\n"+
				"— <b>Arguments:</b>\n"+
				"—— <i>Price/V-Bucks</i> <code>(number)</code>\n\n"+
				"— <b>Example:</b>\n"+
				"—— <code>/setPrice 120</code>\n\n"+
				"🔥| <i>Developed by Seyyed</i>",
			gogram.TextHTML,
		))
		return
	}

	newPrice, err := strconv.Atoi(caption[1])
	if err != nil {
		onError(err, ctx)
	}

	global.Config.Itemshop.PriceLegal = newPrice
	if err = global.Config.UpdateConfig(); err != nil {
		onError(err, ctx)
	}

	ctx.Send(gogram.Text(message.NewPrinter(message.MatchLanguage("en")).Sprintf(
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
	), gogram.TextHTML))
}

func onSetPriceIllegal(ctx gogram.MessageContext) {
	caption := strings.SplitN(ctx.GetCaption(), " ", 2)
	if len(caption) != 2 {
		ctx.Send(gogram.Text(
			"✅| <b>Set Half Legal Price command</b>\n\n"+
				"— <b>Arguments:</b>\n"+
				"—— <i>Price/V-Bucks</i> <code>(number)</code>\n\n"+
				"— <b>Example:</b>\n"+
				"—— <code>/setPrice 120</code>\n\n"+
				"🔥| <i>Developed by Seyyed</i>",
			gogram.TextHTML))
		return
	}

	newPrice, err := strconv.Atoi(caption[1])
	if err != nil {
		onError(err, ctx)
	}

	global.Config.Itemshop.PriceIllegal = newPrice
	if err = global.Config.UpdateConfig(); err != nil {
		onError(err, ctx)
	}

	ctx.Send(gogram.Text(message.NewPrinter(message.MatchLanguage("en")).Sprintf(
		"✅| <b>Half Legal PRICES UPDATED</b>\n\n"+
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
	), gogram.TextHTML))
}

func onForcePost(ctx gogram.MessageContext) {
	caption := strings.SplitN(ctx.GetCaption(), " ", 2)
	if len(caption) != 2 {
		ctx.Send(gogram.Text(
			"✅| <b>Force Post command</b>\n\n"+
				"— <b>Arguments:</b>\n"+
				"—— <i>Target Chat</i>: <code>here / channel</code>\n\n"+
				"— <b>Example:</b>\n"+
				"—— <code>/forcePost here</code>\n"+
				"—— <code>/forcePost channel</code>\n\n"+
				"🔥| <i>Developed by Seyyed</i>",
			gogram.TextHTML))
		return
	}

	switch strings.ToLower(caption[1]) {
	case "here":
		tracker.CheckForUpdates(true, ctx.SenderID())
	case "channel":
		tracker.CheckForUpdates(true, global.Config.Itemshop.Channel)
	}
}
