package main

import (
	"github.com/LlamaNite/llamalog"
	"github.com/haashemi/AcountLandBot/bot"
	"github.com/haashemi/AcountLandBot/config"
	"github.com/haashemi/AcountLandBot/generator"
)

func main() {
	log := llamalog.NewLogger("AKB")

	conf, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configs > %s", err.Error())
	}
	log.Info("Config loaded.")

	gen, err := generator.NewGenerator(conf)
	if err != nil {
		log.Fatal("Failed to initialize generator > %s", err.Error())
	}
	log.Info("Generator loaded.")

	bot.Start(conf, gen)
}
