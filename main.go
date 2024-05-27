package main

import (
	"fmt"
	"os"
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

	file := "test.frag"

	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	runtime.LockOSThread()

	window := glInit()
	defer glTerminate()

	program := initGLProgram(file)

	// vao := makeVao(quad)
	for !window.ShouldClose() {
		glDraw(window, program)
		checkInputs(window, &program)
	}

	glTerminate()
	CleanUp(&program)

	for _, vid := range program.videos {
		endVideo(vid)
	}

}

func checkInputs(window *glfw.Window, program *OpenGLProgram) {
	if window.GetKey(glfw.KeyR) == glfw.Press {
		CleanUp(program)
		*program = initGLProgram(program.fileName)
		fmt.Println("RELOADED!")
	} else if window.GetKey(glfw.KeyF) == glfw.Press {
		// FULLSCREEN
	} else if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}

}
