package main

import (
	"runtime"
)

func main() {

	/*webcam, _ := gocv.VideoCaptureDevice(0)
	window := gocv.NewWindow("Hello")
	img := gocv.NewMat()

	for {
		webcam.Read(&img)
		window.IMShow(img)
		window.WaitKey(1)
	}*/

	runtime.LockOSThread()

	window := glInit()
	defer glTerminate()

	program := initGLProgram()

	// vao := makeVao(quad)
	for !window.ShouldClose() {
		glDraw(window, program)
	}
}
