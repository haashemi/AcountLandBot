package bot

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/LlamaNite/llamaimage"
	"github.com/davegardnerisme/deephash"
	"github.com/haashemi/AcountLandBot/fortnite"
	"github.com/haashemi/AcountLandBot/generator"
	"github.com/haashemi/tgo"
	"github.com/samber/lo"
)

const AnnounceShopText = `✅| New ItemShop!

#️⃣| Hash: <code>%s</code>
⏳| Generated in: <code>%v</code>

🔥| <i>Developed by haashemi.dev</i>`

func (c *Client) trackItemShop() {
	ticker := time.NewTicker(time.Second * 15)

	c.checkForUpdates()
	for range ticker.C {
		isUpdated := c.checkForUpdates()
		if isUpdated {
			c.announceItemshop()
		}
	}
}

func (c *Client) checkForUpdates() (isUpdated bool) {
	data, err := fortnite.GetItemshop()
	if err != nil {
		c.log.Error("Failed to fetch the itemshop > %s", err.Error())
		return false
	}

	shopItems := filterItems(append(data.Data.Featured.Entries, data.Data.Daily.Entries...))
	shopHash := deephash.Hash(shopItems)

	c.shopMut.RLock()
	isUpdated = !bytes.Equal(c.shopHash, shopHash)
	c.shopMut.RUnlock()

	if isUpdated {
		c.shopMut.Lock()
		c.shopHash = shopHash
		c.shopItems = shopItems
		c.shopMut.Unlock()

		c.log.Info("Shop data updated > %s", hex.EncodeToString(shopHash))
	}

	return isUpdated
}

func (c *Client) announceItemshop() {
	// get a copy of the current data
	c.shopMut.RLock()
	items := c.shopItems
	itemsHash := c.shopHash
	c.shopMut.RUnlock()

	if len(items) == 0 {
		return
	}

	shopTabs := generator.SplitSlice(items, 12)
	shopImages := make([]*tgo.InputFile, len(shopTabs))

	genStart := time.Now()
	for _, tab := range shopTabs {
		img, err := c.generator.GenerateItemshop(tab)
		if err != nil {
			c.log.Error("announcer: Failed to generate shop > %v", err)
			return
		}

		buf := bytes.NewBuffer([]byte{})
		if err = llamaimage.SaveToStream(img, buf); err != nil {
			c.log.Error("announcer: Failed to encode the image > %v", err)
			return
		}

		shopImages = append(shopImages, tgo.FileFromReader(fmt.Sprintf("itemshop-%d.png", time.Now().Unix()), buf))
	}
	genEnd := time.Since(genStart)

	c.bot.SendMediaGroup(&tgo.SendMediaGroup{
		ChatId: tgo.ID(c.config.Itemshop.Channel),
		Media: lo.Map(shopImages, func(item *tgo.InputFile, index int) tgo.InputMedia {
			return &tgo.InputMediaDocument{
				Type:    "document",
				Media:   item,
				Caption: lo.Ternary(index+1 == len(shopImages), fmt.Sprintf(AnnounceShopText, itemsHash, genEnd), ""),
			}
		}),
	})
}

// filterItems filters the shop items with bundles and outfits only.
func filterItems(items []fortnite.ItemShopItem) (newItems []fortnite.ItemShopItem) {
	for _, item := range items {
		if item.Bundle != nil {
			newItems = append(newItems, item)
		} else if len(item.Items) > 0 {
			if item.Items[0].Type.Value == "outfit" {
				newItems = append(newItems, item)
			}
		}
	}

	return newItems
}
