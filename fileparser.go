package main

import (
	"encoding/json"
	"os"
)

type InputFile struct {
	Width      int      `json:"width"`
	Height     int      `json:"height"`
	ShaderPath string   `json:"shader"`
	Textures   []string `json:"textures"`
}

func ParseJsonToInputFile(filepath string) (InputFile, error) {

	dat, err := os.ReadFile(filepath)

	if err != nil {
		return InputFile{}, err
	}

	var ret InputFile
	err = json.Unmarshal(dat, &ret)

	if err != nil {
		return InputFile{}, err
	}

	return ret, nil
}
