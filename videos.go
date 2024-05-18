package main

import (
	"fmt"

	"gocv.io/x/gocv"
)

type VideoData struct {
	video        *gocv.VideoCapture
	fps          int
	frames       int
	currentFrame int
	width        int
	height       int
	material     gocv.Mat
	// currentMatIndex int
}

func CreateVideoFromFile(file string) (VideoData, error) {
	video, err := gocv.OpenVideoCapture(file)

	if err != nil {
		fmt.Println("VIDEO OPEN ERROR " + err.Error())
		return VideoData{}, err
	}

	dat := VideoData{
		video:        video,
		fps:          int(video.Get(gocv.VideoCaptureFPS)),
		frames:       int(video.Get(gocv.VideoCaptureFrameCount)),
		width:        int(video.Get(gocv.VideoCaptureFrameWidth)),
		height:       int(video.Get(gocv.VideoCaptureFrameHeight)),
		currentFrame: -1,
		material:     gocv.NewMat(),
		// materials:       []gocv.Mat{gocv.Mat{}, gocv.Mat{}, gocv.Mat{}},
		// currentMatIndex: 0,
	}

	dat.ReadFrame(0)

	return dat, nil

}

func (dat *VideoData) GetData() []int8 {
	return dat.material.DataPtrInt8()
}

func (dat *VideoData) ReadFrame(frame int) {
	if frame == dat.currentFrame {
		return
	}

	read := dat.video.Read(&dat.material)

	if !read {
		dat.video.Set(gocv.VideoCapturePosFrames, 0)
		read := dat.video.Read(&dat.material)
		if !read {
			fmt.Println("VIDEO READ FRAME ERROR")
			return
		}
	}
	dat.currentFrame = frame
	gocv.CvtColor(dat.material, &dat.material, gocv.ColorBGRToRGB)

}

//	video.Set(gocv.VideoCapturePosFrames, 0)
