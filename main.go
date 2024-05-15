package main

import (
	"fmt"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
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

		if window.GetKey(glfw.KeyR) == glfw.Press {
			CleanUp(program)
			program = initGLProgram()
			fmt.Println("RELOADED!")
		}
		glDraw(window, program)
	}
}
