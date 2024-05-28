package main

import (
	"fmt"

	"github.com/go-gl/gl/v2.1/gl"
	"gocv.io/x/gocv"
)

type VideoData struct {
	video        *gocv.VideoCapture
	writer       *gocv.VideoWriter
	texture      uint32
	fps          float64
	frames       int
	currentFrame int
	width        int
	height       int
	material     gocv.Mat
	// currentMatIndex int
}

func CreateVideoFromFile(file string) (*VideoData, error) {
	video, err := gocv.OpenVideoCapture(file)

	if err != nil {
		fmt.Println("VIDEO OPEN ERROR " + err.Error())
		return nil, err
	}

	dat := VideoData{
		video:        video,
		writer:       nil,
		fps:          video.Get(gocv.VideoCaptureFPS),
		frames:       int(video.Get(gocv.VideoCaptureFrameCount)),
		width:        int(video.Get(gocv.VideoCaptureFrameWidth)),
		height:       int(video.Get(gocv.VideoCaptureFrameHeight)),
		currentFrame: -1,
		material:     gocv.NewMat(),
		// materials:       []gocv.Mat{gocv.Mat{}, gocv.Mat{}, gocv.Mat{}},
		// currentMatIndex: 0,
	}

	//

	// dat.material.
	dat.writer = setupVideoWriter(&dat)

	dat.ReadFrame(0)

	//fmt.Printf("COLS %d CHANNELS %d", dat.material.Cols(), dat.material.Channels())

	return &dat, nil

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
	gocv.CvtColor(dat.material, &dat.material, gocv.ColorBGRAToRGBA)

}

func setupVideoWriter(data *VideoData) *gocv.VideoWriter {
	writer, err := gocv.VideoWriterFile("testout.avi", "MPEG", float64(data.fps), data.width, data.height, true)

	if err != nil {
		fmt.Println("Video writer creation error " + err.Error())
		return nil
	}

	return writer
}

func writeData(data *VideoData) {

	// newMat := gocv.NewMat()
	// data.material.CopyTo(&newMat)
	gl.ReadPixels(0, 0, int32(data.width), int32(data.height), gl.BGR, gl.UNSIGNED_BYTE, gl.Ptr(data.GetData()))
	data.writer.Write(data.material)
}

func endVideo(data *VideoData) {
	data.writer.Close()
	data.video.Close()
	/*err := ffmpeg_go.Input("testout.avi").
		Filter("transpose", ffmpeg_go.Args{"0"}).
		Filter("transpose", ffmpeg_go.Args{"2"}).
		Output("test.avi").
		Run()

	if err != nil {
		fmt.Println(err.Error())
	}*/
}

//	video.Set(gocv.VideoCapturePosFrames, 0)
