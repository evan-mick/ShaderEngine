package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type ChannelJson struct {
	ShaderPath string   `json:"shader"`
	Textures   []string `json:"textures"`
}

type InputFile struct {
	Width         int    `json:"width"`
	Height        int    `json:"height"`
	Folder        string `json:"folder"`
	RecordFPS     int32  `json:"recordfps"`
	RecordSeconds int64  `json:recordSeconds`

	// ShaderPath and Textures are default, but overriden by Channels field
	ShaderPath string        `json:"shader"`
	Textures   []string      `json:"textures"`
	Channels   []ChannelJson `json:channels`
}

func ParseJsonToInputFile(filepath string) (InputFile, error) {

	dat, err := os.ReadFile(filepath)

	if err != nil {
		return InputFile{}, err
	}

	var ret InputFile
	err = json.Unmarshal(dat, &ret)

	fmt.Print(ret)

	if err != nil {
		return InputFile{}, err
	}

	return ret, nil
}
