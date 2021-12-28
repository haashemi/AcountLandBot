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
}

func onSetPrice(msg *gogram.Message, c *gogram.Client) {
	caption := strings.SplitN(msg.GetCaption(), " ", 2)
	if len(caption) != 2 {
		return
	}

	newPrice, err := strconv.Atoi(caption[1])
	if err != nil {
		onError(err, msg, c)
	}

	global.Config.Itemshop.Price = newPrice
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
