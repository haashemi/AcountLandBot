package bot

import (
	"accountland/bot/llamahttp"
	"accountland/global"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gogram"
	"image"
	"image/color"
	"image/draw"
	"io"
	"math"
	"sync"
	"time"

	"github.com/LlamaNite/llamaimage"
	"github.com/davegardnerisme/deephash"
	"golang.org/x/text/message"
)

var tracker rawTracker

type rawTracker struct {
	mut      sync.Mutex
	client   *gogram.Client
	shopHash []byte
}

type ItemShopData struct {
	Data struct {
		Featured struct {
			Entries []ItemShopItem `json:"entries"`
		} `json:"featured"`
		Daily struct {
			Name    string         `json:"name"`
			Entries []ItemShopItem `json:"entries"`
		} `json:"daily"`
	} `json:"data"`
}

type ItemShopItem struct {
	FinalPrice int `json:"finalPrice"`
	Bundle     *struct {
		Name  string `json:"name"`
		Image string `json:"image"`
	} `json:"bundle,omitempty"`
	NewDisplayAsset struct {
		MaterialInstances []struct {
			Images struct {
				OfferImage string `json:"OfferImage"`
				Background string `json:"Background"`
			} `json:"images"`
		} `json:"materialInstances"`
	} `json:"newDisplayAsset"`
	Items []struct {
		Name        string      `json:"name"`
		ShopHistory []time.Time `json:"shopHistory"`
		Rarity      struct {
			Value string `json:"value"`
		} `json:"rarity"`
		Images struct {
			SmallIcon string `json:"smallIcon"`
			Icon      string `json:"icon"`
			Featured  string `json:"featured"`
		} `json:"images"`
		Type struct {
			Value string `json:"value"`
		} `json:"type"`
	} `json:"items"`
}

func (main *rawTracker) Start() {
	log.Info("Tracker >> Started")
	for range time.NewTicker(time.Second * 15).C {
		log.Info("tracking itemshop...")
		main.CheckForUpdates()
	}
}

func (main *rawTracker) CheckForUpdates() {
	main.mut.Lock()
	defer main.mut.Unlock()

	data := ItemShopData{}
	statusCode, rawData, err := llamahttp.Get("https://fortnite-api.com/v2/shop/br/combined", nil, nil)
	if err != nil || statusCode != 200 {
		log.Error("Tracker >> S: %d >> Err: %v", statusCode, err)
	} else if err = json.Unmarshal(rawData, &data); err != nil {
		log.Error("Tracker >> Err: %s", err.Error())
	} else if newHash := deephash.Hash(data); bytes.Equal(newHash, main.shopHash) {
		return
	} else {
		main.shopHash = newHash
	}

	items := append(data.Data.Featured.Entries, data.Data.Daily.Entries...)
	itemShopTabs := splitSlice(items, 12)
	main.PostUpdateAlert(len(itemShopTabs))

	for index, tabItems := range itemShopTabs {
		start := time.Now()
		shopImage := main.GenerateImage(tabItems)
		generateTime := time.Since(start)

		buf := new(bytes.Buffer)
		if err = llamaimage.SaveToStream(shopImage, buf); err != nil {
			log.Error(err.Error())
			continue
		}
		main.PostUpdate(buf, index+1, len(itemShopTabs), generateTime)
	}
}

func (main *rawTracker) GenerateImage(items []ItemShopItem) image.Image {
	mainImage := llamaimage.NewImage(1080, 1920)

	if bg, err := global.GetImage("background"); err != nil {
		llamaimage.FillGradient(mainImage, global.Colors.Background.Start, global.Colors.Background.End, llamaimage.GradientOrientationVertical)
	} else {
		llamaimage.Paste(mainImage, bg, 0, 0)
	}

	var wg sync.WaitGroup
	wg.Add(len(items))

	for index, item := range items {
		go func(index int, item ItemShopItem) {
			main.GenerateIcon(mainImage, item, 75+(index%3)*320, 215+(index/3)*390)
			wg.Done()
		}(index, item)
	}

	wg.Wait()

	return mainImage
}

func (main *rawTracker) GenerateIcon(mainImage draw.Image, item ItemShopItem, startX, startY int) {
	// IconRes:  290 - 370
	// Base X-Y:  75 - 215

	// TEMP
	cA, err := llamaimage.HexToRGBA("#0091fa")
	if err == nil {
		cB, err := llamaimage.HexToRGBA("#014dbd")
		if err == nil {
			fillGradient(
				mainImage, cA, cB,
				startX, startX+290,
				startY, startY+370,
			)
		}
	}

	var name, iconURL string
	if len(item.Items) > 0 {
		name = item.Items[0].Name

		if iconURL = item.Items[0].Images.Featured; iconURL == "" {
			if iconURL = item.Items[0].Images.Icon; iconURL == "" {
				iconURL = item.Items[0].Images.SmallIcon
			}
		}
	}

	if item.Bundle != nil {
		name = item.Bundle.Name
		iconURL = item.Bundle.Image
	}

	instances := item.NewDisplayAsset.MaterialInstances
	if len(instances) > 0 {
		iconURL = instances[0].Images.Background
	}

	icon, err := llamahttp.GetImage(iconURL)
	if err == nil {
		icon = llamaimage.Resize(icon, 280+6, 280+6, llamaimage.ResizeFit)
		llamaimage.Paste(mainImage, icon, startX+2, startY+2)
	}

	// NAME
	if len(name) > 0 {
		fontFace, textWidth := global.Fonts.Burbank.FitTextWidth(name, 33, 290-20)
		llamaimage.Write(mainImage, name, color.Black, fontFace, startX+((290-textWidth)/2), startY+297)
	}

	fontFace := global.Fonts.Burbank.NewFace(20)
	printPrice := message.NewPrinter(message.MatchLanguage("en")).Sprintf

	// MainPrice
	if icon, err := global.GetImage("icon-vbucks"); err == nil {
		icon = llamaimage.Resize(icon, 30, 30, llamaimage.ResizeFit)
		llamaimage.Paste(mainImage, icon, startX+7, startY+335)
	}
	priceVBucks := printPrice("%d", item.FinalPrice)
	llamaimage.Write(mainImage, priceVBucks, color.White, fontFace, startX+44, startY+342)

	// RialsPrice
	priceTomans := printPrice("%d T", item.FinalPrice*global.Config.Itemshop.Price)
	margin := global.Fonts.Burbank.GetWidth(20, priceTomans) + 10
	llamaimage.Write(mainImage, priceTomans, color.White, fontFace, startX+290-margin, startY+342)
}

func (main *rawTracker) PostUpdateAlert(itemShopTabs int) {
	post, err := main.client.SendMessage(
		global.Config.Itemshop.Channel,
		main.client.HTML(fmt.Sprintf(
			"✅| <b>ITEMSHOP UPATED!</b>\n\n"+
				"📊| <i>TABS COUNT:</i> <code>%d</code>\n"+
				"🧬| <i>HASH:</i> <code>%s</code>\n\n"+
				"🔥| <i>developed by Seyyed</i>",
			itemShopTabs, hex.EncodeToString(main.shopHash),
		)),
	)

	if err == nil && post != nil {
		defer main.client.GetRawClient().PinChatMessage(
			global.Config.Itemshop.Channel,
			post.GetID(), false, false,
		)
	}
}

func (main *rawTracker) PostUpdate(reader io.Reader, tabCount, allTabsCount int, generateTime time.Duration) {
	doc := gogram.FileFromReaderWithPattern(reader, "AccountLand-*.png")
	main.client.SendMessage(
		global.Config.Itemshop.Channel,
		gogram.Document(
			doc, nil, gogram.Text(fmt.Sprintf(
				"> Tab %d/%d\nGenerated in: %v",
				tabCount, allTabsCount, generateTime,
			)), false,
		),
	)
	doc.Dispose()
}

func splitSlice(items []ItemShopItem, tabSize int) [][]ItemShopItem {
	newItems := []ItemShopItem{}
	for _, item := range items {
		if item.Bundle != nil {
			newItems = append(newItems, item)
		} else if len(item.Items) > 0 {
			if item.Items[0].Type.Value == "outfit" {
				newItems = append(newItems, item)
			}
		}
	}

	tabs := [][]ItemShopItem{}

	for i := 0; i < int(math.Ceil(float64(len(newItems))/float64(tabSize))); i++ {
		tab := []ItemShopItem{}
		for index, item := range newItems[i*tabSize:] {
			if index == tabSize {
				break
			}
			tab = append(tab, item)
		}
		tabs = append(tabs, tab)
	}

	return tabs
}

func fillGradient(img draw.Image, startColor, endColor color.RGBA, startX, endX, startY, endY int) {
	row, column := endX-startX, endY-startY

	SR, SG, SB, SA := float64(startColor.R), float64(startColor.G), float64(startColor.B), float64(startColor.A)
	LR, LG, LB, LA := float64(endColor.R), float64(endColor.G), float64(endColor.B), float64(endColor.A)

	difference_R := (LR - SR) / float64(column)
	difference_G := (LG - SG) / float64(column)
	difference_B := (LB - SB) / float64(column)
	difference_A := (LA - SA) / float64(column)

	R, G, B, A := SR, SG, SB, SA

	for columnP := 0; columnP < column; columnP++ {
		for rowP := 0; rowP < row; rowP++ {
			if rowP > 1 && rowP < row-2 { // 2px from left and right
				if columnP < 2 { // 2px from top
					img.Set(startX+rowP, startY+columnP, color.RGBA{uint8(R), uint8(G), uint8(B), uint8(A)})
				} else if columnP > column-83 && columnP < column-38 {
					img.Set(startX+rowP, startY+columnP, color.White) //RGBA{144, 173, 216, 255}
				} else if columnP > column-38 && columnP < column-3 {
					img.Set(startX+rowP, startY+columnP, color.Black)
				} else if columnP > column-3 { // 2px from bottom
					img.Set(startX+rowP, startY+columnP, color.RGBA{uint8(R), uint8(G), uint8(B), uint8(A)})
				}

			} else {
				img.Set(startX+rowP, startY+columnP, color.RGBA{uint8(R), uint8(G), uint8(B), uint8(A)})
			}
		}

		R += difference_R
		G += difference_G
		B += difference_B
		A += difference_A
	}
}
