package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
)

var paused bool = false
var lastKey glfw.Action = -1

func main() {

	/*webcam, _ := gocv.VideoCaptureDevice(0)
	window := gocv.NewWindow("Hello")
	img := gocv.NewMat()

	for {
		webcam.Read(&img)
		window.IMShow(img)
		window.WaitKey(1)
	}*/

	file := "main.json"

	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	in, err := ParseJsonToInputFile(file)
	//fmt.Println(in)
	if err != nil {
		fmt.Println("ERROR: issue parsing input file " + err.Error())
		return
	}

	runtime.LockOSThread()

	window := glInit(&in)
	defer glTerminate()

	program := initGLProgram(&in)

	// vao := makeVao(quad)
	for !window.ShouldClose() {

		if !paused {
			glDraw(window, program)
		}
		checkInputs(window, &program, &in)
	}

	glTerminate()
	CleanUp(&program)

	for _, vid := range program.videos {
		endVideo(vid)
	}
	videoWriter.Close()

}

func checkInputs(window *glfw.Window, program *OpenGLProgram, in *InputFile) {
	if window.GetKey(glfw.KeyR) == glfw.Press {
		CleanUp(program)
		*program = initGLProgram(in)
		fmt.Println("RELOADED!")
	} else if window.GetKey(glfw.KeyF) == glfw.Press {
		// FULLSCREEN

	} else if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	} else if window.GetKey(glfw.KeyLeftShift) == glfw.Press {
		if window.GetKey(glfw.KeyLeftShift) != lastKey {
			fmt.Println("PAUSED")
			paused = !paused
		}
	}

	lastKey = window.GetKey(glfw.KeyLeftShift)

}
