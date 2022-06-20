package global

import (
	_ "embed"
	"image/color"

	"github.com/LlamaNite/llamaimage"
	"gopkg.in/yaml.v3"
)

//go:embed assets/colors.yaml
var colorsConfigFile []byte

type colors struct {
	DefaultBackground gradientColor
	Rarities          map[string]gradientColor
}

type gradientColor struct{ Start, End, Overlay color.RGBA }

type colorsConfig struct {
	DefaultBackground struct {
		Start   string `yaml:"start"`
		End     string `yaml:"end"`
		Overlay string `yaml:"overlay"`
	} `yaml:"default_background"`
	RarityColors []struct {
		Rarity  string `yaml:"rarity"`
		Start   string `yaml:"start"`
		End     string `yaml:"end"`
		Overlay string `yaml:"overlay"`
	} `yaml:"rarity_colors"`
}

func (main *colors) Of(rarity string) gradientColor {
	if colorData, isExist := main.Rarities[rarity]; isExist {
		return colorData
	}
	return main.DefaultBackground
}

func loadColors(data *colors) error {
	conf := &colorsConfig{}
	if err := yaml.Unmarshal(colorsConfigFile, conf); err != nil {
		return err
	}

	startColor, err := llamaimage.HexToRGBA(conf.DefaultBackground.Start)
	if err != nil {
		return err
	}
	data.DefaultBackground.Start = startColor

	endColor, err := llamaimage.HexToRGBA(conf.DefaultBackground.End)
	if err != nil {
		return err
	}
	data.DefaultBackground.End = endColor

	data.Rarities = map[string]gradientColor{}
	for _, colorData := range conf.RarityColors {
		startColor, err = llamaimage.HexToRGBA(colorData.Start)
		if err != nil {
			return err
		}
		endColor, err = llamaimage.HexToRGBA(colorData.End)
		if err != nil {
			return err
		}
		overlayColor, err := llamaimage.HexToRGBA(colorData.Overlay)
		if err != nil {
			return err
		}

		data.Rarities[colorData.Rarity] = gradientColor{
			Start:   startColor,
			End:     endColor,
			Overlay: overlayColor,
		}
	}

	return nil
}
