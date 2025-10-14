package main

import (
	"encoding/json"
	"os"
	"path"
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

	// ShaderPath and Textures are main channel, Channels is extra
	ShaderPath string   `json:"shader"`
	Textures   []string `json:"textures"`

	Channels []ChannelJson `json:channels`

	// Though this is frag code, it is global
	Includes []string `json:includes`
}

func ParseJsonToInputFile(filepath string) (InputFile, error) {

	dat, err := os.ReadFile(filepath)
	// A little extra logic for supporting passing in folders
	isFolder := false
	folder := ""

	if err != nil {
		_, suffix := path.Split(filepath)
		dat, err = os.ReadFile(filepath + "/" + suffix + ".json") // try to find file version if folder
		isFolder = true
		folder = suffix

		if err != nil {
			return InputFile{}, err
		}
	}

	var ret InputFile
	if isFolder {
		ret.Folder = folder
	}

	err = json.Unmarshal(dat, &ret)

	if err != nil {
		return InputFile{}, err
	}

	return ret, nil
}
