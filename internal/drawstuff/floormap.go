package drawstuff

import (
	"encoding/json"
	"os"
)

type TileMapLayerJson struct {
	Data   []int `json:"data"`
	Height int   `json:"height"`
	Width  int   `json:"width"`
}

type TileMapJson struct {
	Layers []TileMapLayerJson `json:"layers"`
}

func NewTileMapJson(filepath string) (*TileMapJson, error) {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var TileMapJson TileMapJson
	err = json.Unmarshal(contents, &TileMapJson)
	if err != nil {
		return nil, err
	}
	return &TileMapJson, nil
}
