package generator

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"sync"

	"github.com/LlamaNite/llamaimage"
	"github.com/abdullahdiaa/garabic"
	"github.com/haashemi/AcountLandBot/fortnite"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	ErrTooMuchItem = errors.New("items are more than supported items in a tab")
)

func (g *Generator) GenerateItemshop(items []fortnite.ItemShopItem) (image.Image, error) {
	if len(items) > 12 {
		return nil, ErrTooMuchItem
	}

	mainImage := llamaimage.NewImage(1080, 1920)

	bg, err := g.GetImage("background")
	if err != nil {
		return nil, err
	}
	llamaimage.Paste(mainImage, bg, 0, 0)

	var wg sync.WaitGroup
	wg.Add(len(items))

	for index, item := range items {
		go func(index int, item fortnite.ItemShopItem) {
			defer wg.Done()

			icon, err := g.GenerateItemshopIcon(item)
			if err != nil {
				g.log.Error("Failed to generate the icon > %v", err)
				return
			}

			const (
				columns = 3

				xOffset = 75
				yOffset = 215

				iconWidth  = 290
				iconHeight = 370

				xMargin = 30
				yMargin = 20
			)

			x := xOffset + ((index % columns) * (iconWidth + xMargin))  // offset + x position based on index
			y := yOffset + ((index / columns) * (iconHeight + yMargin)) // offset + y position based on the index

			llamaimage.Paste(mainImage, icon, x, y)

		}(index, item)
	}

	wg.Wait()

	return mainImage, nil
}

func (g *Generator) GenerateItemshopIcon(item fortnite.ItemShopItem) (image.Image, error) {
	icon := llamaimage.NewImage(290, 370)

	// draw the icon's background
	llamaimage.FillGradient(icon, color.RGBA{0, 145, 250, 255}, color.RGBA{1, 77, 189, 255}, llamaimage.GradientOrientationVertical)
	draw.Draw(icon, image.Rect(2, 288, 288, 328), image.White, image.Point{}, draw.Over) // draw a 286x40 white box with a 2x2 margin
	draw.Draw(icon, image.Rect(2, 328, 288, 368), image.Black, image.Point{}, draw.Over) // draw a 286x40 black box with a 2x2 margin

	// fetch and paste the icon
	if instances := item.NewDisplayAsset.MaterialInstances; len(instances) > 0 {
		img, err := getImage(instances[0].Images.Background, 286, 286)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch icon %s > %s", instances[0].ID, err.Error())
		}
		llamaimage.Paste(icon, img, 2, 2) // there should be a 2x2 margin
	}

	conf := g.config.Itemshop
	printer := message.NewPrinter(language.English)
	priceFont := g.fonts.Burbank.NewFace(27)

	primaryPrice := printer.Sprintf("%d %s", item.FinalPrice*conf.PrimaryPrice, conf.PrimaryCurrency)
	secondaryPrice := printer.Sprintf("%d %s", item.FinalPrice*conf.SecondaryPrice, conf.SecondaryCurrency)

	llamaimage.Write(icon, primaryPrice, color.Black, priceFont, 10, 295)
	llamaimage.Write(icon, secondaryPrice, color.White, priceFont, 10, 333)

	const titleFontSize = 23
	titleFont := g.fonts.KalamehBold.NewFace(titleFontSize) // Note: It may not support ENG letters, so change it if you want

	priceType := garabic.Shape(conf.PrimaryTitle)
	margin := g.fonts.KalamehBold.GetWidth(titleFontSize, priceType) + 10
	llamaimage.Write(icon, priceType, color.Black, titleFont, 287-margin, 286)

	priceType = garabic.Shape(conf.SecondaryTitle)
	margin = g.fonts.KalamehBold.GetWidth(titleFontSize, priceType) + 10
	llamaimage.Write(icon, priceType, color.White, titleFont, 287-margin, 324)

	return icon, nil
}
