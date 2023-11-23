package fortnite

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ItemShop struct {
	Status int `json:"status"`
	Data   struct {
		Featured struct {
			Entries []ItemShopItem `json:"entries"`
		} `json:"featured"`
		Daily struct {
			Name    string         `json:"name"`
			Entries []ItemShopItem `json:"entries"`
		} `json:"daily"`
	} `json:"data"`
}

type ItemShopItem struct {
	FinalPrice int `json:"finalPrice"`
	Bundle     any `json:"bundle,omitempty"`
	Items      []struct {
		Type struct {
			Value string `json:"value"`
		} `json:"type"`
	} `json:"items"`
	NewDisplayAsset struct {
		MaterialInstances []struct {
			ID     string `json:"id"`
			Images struct {
				Background string `json:"Background"`
			} `json:"images"`
		} `json:"materialInstances"`
	} `json:"newDisplayAsset"`
}

func GetItemshop() (data *ItemShop, err error) {
	resp, err := http.Get("https://fortnite-api.com/v2/shop/br/combined")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code %d", resp.StatusCode)
	}

	data = &ItemShop{}
	if err = json.NewDecoder(resp.Body).Decode(data); err != nil {
		return nil, err
	}

	return data, nil
}
