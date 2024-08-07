package main

import (
	"fmt"
	"image"

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

// var formatCtx *gmf.FmtCtx //:= gmf.NewCtx()
// var outputCtx *gmf.FmtCtx

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
	writer, err := gocv.VideoWriterFile(data.fileName+".avi", "MPEG", float64(data.recordFPS), data.width, data.height, true)

	if err != nil {
		fmt.Println("Video writer creation error " + err.Error())
		return nil
	}

	return writer
}

// TODO: Change this so its setup based on the program
// not even fully sure why it is video data rn
// f stands for "framebuffer," o stands for "output"
func writeData(fWidth int32, fHeight int32, oWidth int32, oHeight int32) {

	// pixels := make([]uint8, width*height*3)

	// Read the pixels from the framebuffer
	// gl.ReadPixels(0, 0, int32(width), int32(height), gl.RGB, gl.UNSIGNED_BYTE, gl.Ptr(&pixels[0]))

	img := image.NewRGBA(image.Rect(0, 0, int(fWidth), int(fHeight)))
	gl.ReadPixels(0, 0, fWidth, fHeight, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
	FlipVertically(img)

	mat, _ := gocv.ImageToMatRGB(img)
	resizedMat := gocv.NewMat()
	gocv.Resize(mat, &resizedMat, image.Point{int(oWidth), int(oHeight)}, 0, 0, gocv.InterpolationLinear)
	videoWriter.Write(resizedMat)

	mat.Close()
	resizedMat.Close()
	// draw.FlipVertically(img)
	// return img
	// frame := gmf.NewFrameFromBytes(frameData, 1920, 1080)
	// // codecCtx.Encode(frame)
	// outputCtx.WriteFrame(frame)
	// Create an OpenCV Mat with the pixel data
	/*mat, err := openGLToCVMat(pixels, width, height)
	if err != nil {
		fmt.Println(err.Error())
		return
	}*/

	// Save the Mat to an image file
	// gocv.IMWrite("output.png", mat)
	//videoWriter.Write(mat)

}

// Flips an RGBA image vertically
func FlipVertically(img *image.RGBA) {
	bounds := img.Bounds()
	for y := 0; y < bounds.Dy()/2; y++ {
		for x := 0; x < bounds.Dx(); x++ {
			i1 := img.PixOffset(x, y)
			i2 := img.PixOffset(x, bounds.Dy()-y-1)
			for k := 0; k < 4; k++ {
				img.Pix[i1+k], img.Pix[i2+k] = img.Pix[i2+k], img.Pix[i1+k]
			}
		}
	}
}

func openGLToCVMat(pixels []byte, width, height int) (gocv.Mat, error) {
	// Create an empty OpenCV Mat with the correct dimensions
	mat := gocv.NewMatWithSize(height, width, gocv.MatTypeCV8UC3)

	// Use pointers to directly set the pixel data in the Mat
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := (y*width + x) * 3
			r, g, b := pixels[idx], pixels[idx+1], pixels[idx+2]

			// OpenCV uses BGR format
			mat.SetUCharAt(y, x*3, b)
			mat.SetUCharAt(y, x*3+1, g)
			mat.SetUCharAt(y, x*3+2, r)
		}
	}

	return mat, nil
}

func endVideo(data *VideoData) {
	// data.writer.Close()
	videoWriter.Close()
	// outputCtx.CloseOutputAndRelease()
	// data.video.Close()
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
