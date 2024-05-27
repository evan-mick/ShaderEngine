package main

import (
	"encoding/json"
	"os"
)

type InputFile struct {
	width    int
	height   int
	textures []string
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
