package global

import (
	"io/fs"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

var Config *rawConfig

type rawConfig struct {
	TelegramBot struct {
		CachePath  string  `yaml:"CachePath"`
		APIID      string  `yaml:"APIID"`
		APIHash    string  `yaml:"APIHash"`
		BotToken   string  `yaml:"BotToken"`
		SuperUsers []int64 `yaml:"SuperUsers"`
		Admins     []int64 `yaml:"Admins"`
	} `yaml:"TelegramBot"`

	Itemshop struct {
		Channel      int64 `yaml:"Channel"`
		PriceLegal   int   `yaml:"PriceLegal"`
		PriceIllegal int   `yaml:"PriceIllegal"`
	} `yaml:"Itemshop"`

	Colors struct {
		DefaultBackground struct {
			Start string `yaml:"Start"`
			End   string `yaml:"End"`
		} `yaml:"DefaultBackground"`
		RarityColors []struct {
			Rarity  string `yaml:"Rarity"`
			Start   string `yaml:"Start"`
			End     string `yaml:"End"`
			Overlay string `yaml:"Overlay"`
		} `yaml:"RarityColors"`
	} `yaml:"Colors"`
}

func loadConfig() *rawConfig {
	configBytes, err := ioutil.ReadFile("./config.yaml")
	checkErr(err)

	data := &rawConfig{}
	checkErr(yaml.Unmarshal(configBytes, data))

	return data
}

func (data *rawConfig) UpdateConfig() error {
	configBytes, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile("config.yaml", configBytes, fs.ModePerm)
}
