package bot

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/LlamaNite/llamaimage"
	"github.com/haashemi/AcountLandBot/generator"
	"github.com/haashemi/tgo"
	"github.com/haashemi/tgo/routers/message"
	printer "golang.org/x/text/message"
)

const SetPriceHelp = `✅| <b>Set %s Price command</b>

— <b>Arguments:</b>
—— <i>Price/1,000 V-Bucks</i> as <code>(number)</code>

— <b>Example:</b>
—— <code>/%s %d</code>

🔥| <i>Developed by haashemi.dev</i>`

const SetPriceText = `✅| <b>%s Prices Updated!</b>

— V-Bucks — Price
💵| <code>200V — %d%s</code>
💵| <code>500V — %d%s</code>
💵| <code>800V — %d%s</code>
💵| <code>1,200V — %d%s</code>
💵| <code>1,500V — %d%s</code>
💵| <code>1,600V — %d%s</code>
💵| <code>2,000V — %d%s</code>

🔥| <i>Developed by haashemi.dev</i>`

const ItemshopText = `⏳| I'm cooking. wait a few seconds...

📑| %d Tabs
#️⃣| Hash: <code>%s</code>

🔥| <i>Developed by haashemi.dev</i>`

const ItemshopTabText = `✅| Tab %d/%d

⏳| Generated in: <code>%v</code>

🔥| <i>Developed by haashemi.dev</i>`

func (c *Client) SetPrimaryPrice(ctx *message.Context) {
	caption := strings.SplitN(ctx.String(), " ", 2)
	if len(caption) != 2 {
		ctx.Send(&tgo.SendMessage{Text: fmt.Sprintf(SetPriceHelp, "Primary", "setpp", c.config.Itemshop.PrimaryPrice)})
		return
	}

	newPrice, err := strconv.Atoi(caption[1])
	if err != nil {
		handleError(err, ctx)
	}

	c.config.Itemshop.PrimaryPrice = newPrice
	if err = c.config.Save(); err != nil {
		handleError(err, ctx)
	}

	ctx.Send(&tgo.SendMessage{Text: printer.NewPrinter(printer.MatchLanguage("en")).Sprintf(
		SetPriceText, "Primary",
		200*newPrice, c.config.Itemshop.PrimaryCurrency,
		500*newPrice, c.config.Itemshop.PrimaryCurrency,
		800*newPrice, c.config.Itemshop.PrimaryCurrency,
		1200*newPrice, c.config.Itemshop.PrimaryCurrency,
		1500*newPrice, c.config.Itemshop.PrimaryCurrency,
		1600*newPrice, c.config.Itemshop.PrimaryCurrency,
		2000*newPrice, c.config.Itemshop.PrimaryCurrency,
	)})
}

func (c *Client) SetSecondaryPrice(ctx *message.Context) {
	caption := strings.SplitN(ctx.String(), " ", 2)
	if len(caption) != 2 {
		ctx.Send(&tgo.SendMessage{Text: fmt.Sprintf(SetPriceHelp, "Secondary", "setsp", c.config.Itemshop.SecondaryPrice)})
		return
	}

	newPrice, err := strconv.Atoi(caption[1])
	if err != nil {
		handleError(err, ctx)
	}

	c.config.Itemshop.SecondaryPrice = newPrice
	if err = c.config.Save(); err != nil {
		handleError(err, ctx)
	}

	ctx.Send(&tgo.SendMessage{Text: printer.NewPrinter(printer.MatchLanguage("en")).Sprintf(
		SetPriceText, "Secondary",
		200*newPrice, c.config.Itemshop.SecondaryCurrency,
		500*newPrice, c.config.Itemshop.SecondaryCurrency,
		800*newPrice, c.config.Itemshop.SecondaryCurrency,
		1200*newPrice, c.config.Itemshop.SecondaryCurrency,
		1500*newPrice, c.config.Itemshop.SecondaryCurrency,
		1600*newPrice, c.config.Itemshop.SecondaryCurrency,
		2000*newPrice, c.config.Itemshop.SecondaryCurrency,
	)})
}

func (c *Client) Itemshop(ctx *message.Context) {
	// get a copy of the current data
	c.shopMut.RLock()
	items := c.shopItems
	itemsHash := c.shopHash
	c.shopMut.RUnlock()

	// don't confuse the end-user by saying generating zero images =)
	if len(items) == 0 {
		ctx.Send(&tgo.SendMessage{Text: "⚠️| ItemShop is not ready, yet."})
		return
	}

	shopTabs := generator.SplitSlice(items, 12)

	// just saying that we are cooking
	ctx.Send(&tgo.SendMessage{Text: fmt.Sprintf(ItemshopText, len(shopTabs), hex.EncodeToString(itemsHash))})

	// now cook the images!
	for idx, tab := range shopTabs {
		ctx.Bot.SendChatAction(&tgo.SendChatAction{ChatId: tgo.ID(ctx.Chat.Id), Action: "upload_document"})

		start := time.Now()
		img, err := c.generator.GenerateItemshop(tab)
		if err != nil {
			handleError(err, ctx)
			return // if only a tab fails, it's a useless response to the user, so stop.
		}
		end := time.Since(start)

		buf := bytes.NewBuffer([]byte{})
		if err = llamaimage.SaveToStream(img, buf); err != nil {
			handleError(err, ctx)
			return // if only a tab fails, it's a useless response to the user, so stop.
		}

		ctx.Send(&tgo.SendDocument{
			Document: tgo.FileFromReader(fmt.Sprintf("itemshop-%d.png", start.Unix()), buf),
			Caption:  fmt.Sprintf(ItemshopTabText, idx+1, len(shopTabs), end),
		})
	}
}
