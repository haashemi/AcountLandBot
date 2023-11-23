package generator

import (
	"github.com/LlamaNite/llamaimage"
	"gopkg.in/yaml.v3"
)

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

func loadFonts() (*Fonts, error) {
	fonts := &Fonts{}

	if file, err := assets.ReadFile("assets/fonts/burbank.ttf"); err != nil {
		return nil, err
	} else {
		if fonts.Burbank, err = llamaimage.OpenFont(file); err != nil {
			return nil, err
		}
	}

	if file, err := assets.ReadFile("assets/fonts/Kalameh-Bold.ttf"); err != nil {
		return nil, err
	} else {
		if fonts.KalamehBold, err = llamaimage.OpenFont(file); err != nil {
			return nil, err
		}
	}

	return fonts, nil
}

func loadColors() (*Colors, error) {
	colors := &Colors{}
	file, err := assets.ReadFile("assets/colors.yaml")
	if err != nil {
		return nil, err
	}

	conf := &colorsConfig{}
	if err := yaml.Unmarshal(file, conf); err != nil {
		return nil, err
	}

	startColor, err := llamaimage.HexToRGBA(conf.DefaultBackground.Start)
	if err != nil {
		return nil, err
	}
	colors.DefaultBackground.Start = startColor

	endColor, err := llamaimage.HexToRGBA(conf.DefaultBackground.End)
	if err != nil {
		return nil, err
	}
	colors.DefaultBackground.End = endColor

	colors.Rarities = map[string]GradientColor{}
	for _, colorData := range conf.RarityColors {
		startColor, err = llamaimage.HexToRGBA(colorData.Start)
		if err != nil {
			return nil, err
		}
		endColor, err = llamaimage.HexToRGBA(colorData.End)
		if err != nil {
			return nil, err
		}
		overlayColor, err := llamaimage.HexToRGBA(colorData.Overlay)
		if err != nil {
			return nil, err
		}

		colors.Rarities[colorData.Rarity] = GradientColor{
			Start:   startColor,
			End:     endColor,
			Overlay: overlayColor,
		}
	}

	return colors, nil
}
