package global

import (
	"embed"
	"image"

	"github.com/LlamaNite/llamaimage"
)

//go:embed assets/images/*
var images embed.FS

func GetImage(filename string) (image.Image, error) {
	return llamaimage.OpenImageFromEFS(images, "assets/images/"+filename+".png")
}
