package main

// https://github.com/go-audio
// helpful packages

// Process
// Load audio file
//		Debating whether to load all at once, or stream chunk, should allow for both, will load all to start
//		Should have a getter then for PCM data and required fields
//		Use general audio interface
// File should be played in processed chunks
//		Go routine is constantly trying to load more chunks (or at least up to a point)
//		Loading chunks
//			Get PCM data for current position up to chunk size
//			Input that into first row of texture
//			Run FFT
//			Input that into second row of texture
// Every "Frame"
//		Realtime: Get chunk from current frame and process it live
//			could also assume a fixed framerate? like yeah we're running hundreds of fps but maybe the audio only runs 60 at a time or whatever
//		Frame: Get chunk from current fixed frame
//		after chunk obtained, put it in texture
//			Possibly have 2 textures and flip between them
//		then go on to draw frame

// I think the move is to get a basic version of this, then profile and see if its actually a bottleneck compared with opengl stuff

type AudioLoadType uint8

const (
	LoadAll = 0
	Stream  = 1
)

func CreateAudioTexture() {

}

func LoadAudioFile() {

}
