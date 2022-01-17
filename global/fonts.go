package global

import (
	_ "embed"

	"github.com/LlamaNite/llamaimage"
)

var Fonts rawFonts

//go:embed assets/fonts/burbank.ttf
var burbank []byte

//go:embed assets/fonts/Kalameh-Bold.ttf
var KalamehBold []byte

type rawFonts struct {
	Burbank     *llamaimage.LlamaFont
	KalamehBold *llamaimage.LlamaFont
}

func loadFonts() rawFonts {
	data := rawFonts{}

	fontFace, err := llamaimage.OpenFont(burbank)
	checkErr(err)
	data.Burbank = fontFace

	fontFace, err = llamaimage.OpenFont(KalamehBold)
	checkErr(err)
	data.KalamehBold = fontFace

	return data
}
