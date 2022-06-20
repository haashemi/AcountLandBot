package global

import (
	_ "embed"

	"github.com/LlamaNite/llamaimage"
)

//go:embed assets/fonts/burbank.ttf
var burbank []byte

//go:embed assets/fonts/Kalameh-Bold.ttf
var KalamehBold []byte

type fonts struct {
	Burbank     *llamaimage.LlamaFont
	KalamehBold *llamaimage.LlamaFont
}

func loadFonts(data *fonts) error {
	fontFace, err := llamaimage.OpenFont(burbank)
	if err != nil {
		return err
	}
	data.Burbank = fontFace

	fontFace, err = llamaimage.OpenFont(KalamehBold)
	if err != nil {
		return err
	}
	data.KalamehBold = fontFace

	return nil
}
