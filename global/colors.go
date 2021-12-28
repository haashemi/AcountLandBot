package global

import (
	"image/color"

	"github.com/LlamaNite/llamaimage"
)

var Colors rawColors

type rawColors struct {
	Background struct{ Start, End color.RGBA }
	Rarities   map[string]struct{ Start, End, Overlay color.RGBA }
}

// ToDo: Add rarities
func loadColors(config *rawConfig) rawColors {
	var err error
	data := rawColors{}

	data.Background.Start, err = llamaimage.HexToRGBA(config.Colors.DefaultBackground.Start)
	checkErr(err)
	data.Background.End, err = llamaimage.HexToRGBA(config.Colors.DefaultBackground.End)
	checkErr(err)

	return data
}
