package global

import (
	"embed"
	"image"
	"os"

	"github.com/LlamaNite/llamaimage"
	"github.com/LlamaNite/llamalog"
)

//go:embed assets/images/*
var images embed.FS

var Fonts = &fonts{}
var Config = &config{}
var Colors = &colors{}

func init() {
	checkError(loadFonts(Fonts))
	checkError(loadConfig(Config))
	checkError(loadColors(Colors))
}

func checkError(err error) {
	if err == nil {
		return
	}
	llamalog.NewLogger("global").Error(err.Error())
	os.Exit(1)
}

func GetImage(filename string) (image.Image, error) {
	return llamaimage.OpenImageFromEFS(images, "assets/images/"+filename+".png")
}
