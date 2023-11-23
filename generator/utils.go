package generator

import (
	"fmt"
	"image"
	"math"
	"net/http"

	"github.com/LlamaNite/llamaimage"
)

func SplitSlice[T any](items []T, splitSize int) [][]T {
	tabs := [][]T{}

	for i := 0; i < int(math.Ceil(float64(len(items))/float64(splitSize))); i++ {
		tab := []T{}
		for index, item := range items[i*splitSize:] {
			if index == splitSize {
				break
			}
			tab = append(tab, item)
		}
		tabs = append(tabs, tab)
	}

	return tabs
}

func getImage(url string, width, height float64) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch > status %d", resp.StatusCode)
	}

	icon, err := llamaimage.OpenImage(resp.Body)
	if err != nil {
		return nil, err
	}

	icon = llamaimage.Resize(icon, width, height, llamaimage.ResizeFit)

	return icon, nil
}
