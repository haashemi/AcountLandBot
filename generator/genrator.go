package generator

import (
	"embed"
	"image"
	"image/color"

	"github.com/LlamaNite/llamaimage"
	"github.com/LlamaNite/llamalog"
	"github.com/haashemi/AcountLandBot/config"
)

//go:embed assets
var assets embed.FS

type Generator struct {
	fonts  *Fonts
	colors *Colors
	config *config.Config

	log *llamalog.Logger
}

type Fonts struct {
	Burbank     *llamaimage.LlamaFont
	KalamehBold *llamaimage.LlamaFont
}

type Colors struct {
	DefaultBackground GradientColor
	Rarities          map[string]GradientColor
}

type GradientColor struct {
	Start   color.RGBA
	End     color.RGBA
	Overlay color.RGBA
}

func NewGenerator(config *config.Config) (generator *Generator, err error) {
	generator = &Generator{config: config, log: llamalog.NewLogger("Generator")}

	generator.fonts, err = loadFonts()
	if err != nil {
		return nil, err
	}

	generator.colors, err = loadColors()
	if err != nil {
		return nil, err
	}

	return generator, nil
}

func (g *Generator) GetImage(filename string) (image.Image, error) {
	return llamaimage.OpenImageFromEFS(assets, "assets/images/"+filename+".png")
}

func (g *Generator) ColorOf(rarity string) GradientColor {
	if colorData, isExist := g.colors.Rarities[rarity]; isExist {
		return colorData
	}
	return g.colors.DefaultBackground
}
