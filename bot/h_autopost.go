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
	"net/http"
	"sync"
	"time"

	"github.com/LlamaNite/llamaimage"
	"github.com/abdullahdiaa/garabic"
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
	Status int `json:"status"`
	Data   struct {
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
	Bundle     any `json:"bundle,omitempty"`
	Items      []struct {
		Type struct {
			Value string `json:"value"`
		} `json:"type"`
	} `json:"items"`
	NewDisplayAsset struct {
		MaterialInstances []struct {
			ID     string `json:"id"`
			Images struct {
				Background string `json:"Background"`
			} `json:"images"`
		} `json:"materialInstances"`
	} `json:"newDisplayAsset"`
}

func (main *rawTracker) Start() {
	log.Info("Tracker >> Started")
	main.CheckForUpdates(false, global.Config.Itemshop.Channel)
	for range time.NewTicker(time.Second * 15).C {
		log.Info("tracking itemshop...")
		main.CheckForUpdates(false, global.Config.Itemshop.Channel)
	}
}

func (main *rawTracker) CheckForUpdates(forcePost bool, target int64) {
	main.mut.Lock()
	defer main.mut.Unlock()

	resp, err := llamahttp.Do(http.MethodGet, "https://fortnite-api.com/v2/shop/br/combined", llamahttp.Options{})
	if err != nil {
		log.Error("Tracker >> Err: %s", err.Error())
		return
	}
	defer resp.Body.Close()

	data := ItemShopData{}
	if resp.StatusCode != http.StatusOK {
		log.Error("Tracker >> StatusCode: %d", resp.StatusCode)
		return
	} else if json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Error("Tracker >> DecoderErr: %s", err.Error())
	} else if newHash := deephash.Hash(data); bytes.Equal(newHash, main.shopHash) && !forcePost {
		return
	} else {
		oldHash := main.shopHash
		main.shopHash = newHash

		if bytes.Equal(oldHash, []byte{}) && !forcePost {
			return
		}
	}

	items := append(data.Data.Featured.Entries, data.Data.Daily.Entries...)
	itemShopTabs := splitSlice(items, 12)
	main.PostUpdateAlert(len(itemShopTabs), target)

	for index, tabItems := range itemShopTabs {
		start := time.Now()
		shopImage := main.GenerateImage(tabItems)
		generateTime := time.Since(start)

		buf := new(bytes.Buffer)
		if err = llamaimage.SaveToStream(shopImage, buf); err != nil {
			log.Error(err.Error())
			continue
		}
		main.PostUpdate(buf, index+1, len(itemShopTabs), generateTime, target)
	}
}

func (main *rawTracker) GenerateImage(items []ItemShopItem) image.Image {
	mainImage := llamaimage.NewImage(1080, 1920)

	if bg, err := global.GetImage("background"); err != nil {
		llamaimage.FillGradient(mainImage, global.Colors.DefaultBackground.Start, global.Colors.DefaultBackground.End, llamaimage.GradientOrientationVertical)
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
	if instances := item.NewDisplayAsset.MaterialInstances; len(instances) > 0 {
		iconURL := instances[0].Images.Background
		img, err := llamahttp.GetImage(instances[0].ID, iconURL)
		if err != nil {
			log.Error("failed to fetch {%s} icon > %s", instances[0].ID, err.Error())
		} else {
			img = llamaimage.Resize(img, 286, 286, llamaimage.ResizeFit)
			llamaimage.Paste(mainImage, img, startX+2, startY+2)
		}
	}

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

	priceFont := global.Fonts.Burbank.NewFace(27)
	priceTypeFont := global.Fonts.KalamehBold.NewFace(23)
	printPrice := message.NewPrinter(message.MatchLanguage("en")).Sprintf

	// RialsPrice Legal
	priceTomans := printPrice("%d T", item.FinalPrice*global.Config.Itemshop.PriceIllegal)
	llamaimage.Write(mainImage, priceTomans, color.Black, priceFont, startX+10, startY+295)

	priceType := garabic.Shape("نيمه قانوني")
	margin := global.Fonts.KalamehBold.GetWidth(23, priceType) + 10
	llamaimage.Write(mainImage, priceType, color.Black, priceTypeFont, startX+287-margin, startY+286)

	// RialsPrice Illegal
	priceTomans = printPrice("%d T", item.FinalPrice*global.Config.Itemshop.PriceLegal)
	llamaimage.Write(mainImage, priceTomans, color.White, priceFont, startX+10, startY+333)

	priceType = garabic.Shape("قانوني")
	margin = global.Fonts.KalamehBold.GetWidth(23, priceType) + 10
	llamaimage.Write(mainImage, priceType, color.White, priceTypeFont, startX+287-margin, startY+324)
}

func (main *rawTracker) PostUpdateAlert(itemShopTabs int, target int64) {
	post, err := main.client.SendMessage(
		target,
		main.client.HTML(fmt.Sprintf(
			"✅| <b>ITEMSHOP UPDATED!</b>\n\n"+
				"📊| <i>TABS COUNT:</i> <code>%d</code>\n"+
				"🧬| <i>HASH:</i> <code>%s</code>\n\n"+
				"🔥| <i>developed by Seyyed</i>",
			itemShopTabs, hex.EncodeToString(main.shopHash),
		)),
	)

	if err == nil && post != nil {
		defer main.client.GetRawClient().PinChatMessage(target, post.GetID(), false, false)
	}
}

func (main *rawTracker) PostUpdate(reader io.Reader, tabCount, allTabsCount int, generateTime time.Duration, target int64) {
	doc := gogram.FileFromReaderWithPattern(reader, "AccountLand-*.png")

	main.client.GetRawClient().GetChat(target)

	_, err := main.client.SendMessage(
		target,
		gogram.Document(
			doc, nil, gogram.Text(fmt.Sprintf(
				"> Tab %d/%d\nGenerated in: %v",
				tabCount, allTabsCount, generateTime,
			)), false,
		),
	)
	if err != nil {
		fmt.Println(err)
	}

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
				} else if columnP > column-83 && columnP <= column-43 {
					img.Set(startX+rowP, startY+columnP, color.White)
				} else if columnP > column-43 && columnP < column-3 {
					img.Set(startX+rowP, startY+columnP, color.Black)
					// } else if columnP > column-83 && columnP < column-38 {
					// img.Set(startX+rowP, startY+columnP, color.White) //RGBA{144, 173, 216, 255}
					// } else if columnP > column-38 && columnP < column-3 {
					// 	img.Set(startX+rowP, startY+columnP, color.Black)
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
