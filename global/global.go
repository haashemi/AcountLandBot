package global

import (
	"os"

	"github.com/LlamaNite/llamalog"
)

func init() {
	Config = loadConfig()
	Colors = loadColors(Config)

	Fonts = loadFonts()
}

func checkErr(err error) {
	if err != nil {
		llamalog.NewLogger("global").Error(err.Error())
		os.Exit(1)
	}
}
