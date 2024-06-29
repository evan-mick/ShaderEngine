package main

import (
	"fmt"
	"image/jpeg"
	"os"

	"github.com/cheggaaa/pb/v3"
	"github.com/go-gl/gl/v2.1/gl"
	"gocv.io/x/gocv"
)

type VideoData struct {
	video *gocv.VideoCapture
	//writer       *gocv.VideoWriter
	texture          uint32
	fps              float64
	frames           int
	currentFrame     int
	width            int
	height           int
	material         *gocv.Mat
	allFramesRead    bool
	allFrames        []*gocv.Mat
	webcam           bool
	removeBackground bool
	filename         string
	// currentMatIndex int
}

var videoWriter *gocv.VideoWriter
var writerData *VideoData

func CreateVideoFromFile(file string) (*VideoData, error) {
	var video *gocv.VideoCapture
	var err error = nil

	if file == "WEBCAM" {
		video, err = gocv.OpenVideoCapture(0)
	} else {
		video, err = gocv.OpenVideoCapture(file)
	}

	if err != nil {
		fmt.Println("VIDEO OPEN ERROR " + err.Error())
		return nil, err
	}

	setMat := gocv.NewMat()
	frames := int(video.Get(gocv.VideoCaptureFrameCount))
	dat := VideoData{
		video: video,
		// writer:       nil,
		fps:              video.Get(gocv.VideoCaptureFPS),
		frames:           frames,
		width:            int(video.Get(gocv.VideoCaptureFrameWidth)),
		height:           int(video.Get(gocv.VideoCaptureFrameHeight)),
		currentFrame:     -1,
		material:         &setMat,
		allFrames:        make([]*gocv.Mat, 0, frames),
		webcam:           (file == "WEBCAM"),
		removeBackground: true,
		filename:         file,
		// materials:       []gocv.Mat{gocv.Mat{}, gocv.Mat{}, gocv.Mat{}},
		// currentMatIndex: 0,
	}

	//

	// dat.material.
	// used to be apart of video data
	// videoWriter = setupVideoWriter(program)

	//dat.ReadFrame(0)

	//fmt.Printf("COLS %d CHANNELS %d", dat.material.Cols(), dat.material.Channels())

	return &dat, nil

}

func (dat *VideoData) GetData() []int8 {
	return dat.material.DataPtrInt8()
}

func (dat *VideoData) ReadFrame(frame int) {
	if !dat.webcam && frame == dat.currentFrame {
		return
	}

	if dat.allFramesRead {
		dat.material = dat.allFrames[frame]
		return
	}

	read := dat.video.Read(dat.material)

	if !read {

		if !dat.webcam {
			dat.video.Set(gocv.VideoCapturePosFrames, 0)
		}

		read := dat.video.Read(dat.material)
		if !read {
			fmt.Println("VIDEO READ FRAME ERROR")
			return
		}
	}
	dat.currentFrame = frame
	gocv.CvtColor(*dat.material, dat.material, gocv.ColorBGRAToRGBA)

}

func (dat *VideoData) ReadAllFrames() {

	mat := gocv.NewMat()

	j := 0
	fmt.Printf("Reading %s\n", dat.filename)
	bar := pb.StartNew(dat.frames)
	for i := range dat.frames {
		dat.video.Set(gocv.VideoCapturePosFrames, float64(i))
		attempts := 0
		for !dat.video.Read(&mat) && attempts < 10 {
			attempts++
		}

		if attempts > 1000 {
			break
		}

		//percent := int(float32(i) / float32(dat.frames) * 100.0)
		//if percent%50 == 0 {
		//fmt.Printf("Video %d percent done", percent)
		//}
		bar.Increment()
		if mat.Empty() {
			continue
		}

		newMat := gocv.NewMat()
		dat.video.Read(&newMat)
		gocv.CvtColor(mat, &newMat, gocv.ColorBGRAToRGBA)
		dat.allFrames = append(dat.allFrames, &newMat)
		j++
	}
	bar.Finish()
	fmt.Printf("Video read\n")
	dat.frames = min(j, min(dat.frames, len(dat.allFrames)))
	dat.allFramesRead = true

}

func setupVideoWriter(data *OpenGLProgram) *gocv.VideoWriter {
	writer, err := gocv.VideoWriterFile("testout.avi", "MPEG", float64(data.recordFPS), data.width, data.height, true)

	if err != nil {
		fmt.Println("Video writer creation error " + err.Error())
		return nil
	}

	return writer
}

func writeData(data *VideoData) {

	// newMat := gocv.NewMat()
	// data.material.CopyTo(&newMat)
	gl.ReadPixels(0, 0, int32(data.width), int32(data.height), gl.BGRA, gl.UNSIGNED_BYTE, gl.Ptr(writerData.GetData()))

	//if  {
	//fmt.Println("EMPTY BADKAJDFKLDAJKLFJAKLFJDAK")
	//}
	f, err := os.Create("img.jpeg")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if writerData.material.Empty() {

		panic("IS EMPTY")
	}

	img, err := writerData.material.ToImage()
	if err != nil {
		panic(err)
	}

	jpeg.Encode(f, img, nil)

	//videoWriter.Write(*data.material)
}

func endVideo(data *VideoData) {
	// data.writer.Close()
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
