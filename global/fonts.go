package global

import (
	_ "embed"

	"github.com/LlamaNite/llamaimage"
)

var Fonts rawFonts

//go:embed assets/fonts/burbank.ttf
var burbank []byte

type rawFonts struct {
	Burbank *llamaimage.LlamaFont
}

func loadFonts() rawFonts {
	data := rawFonts{}

	fontFace, err := llamaimage.OpenFont(burbank)
	checkErr(err)
	data.Burbank = fontFace

	return data
}
