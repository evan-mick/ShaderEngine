package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/go-gl/glfw/v3.3/glfw"
)

var paused bool = false
var lastKey glfw.Action = -1
var autoHotReload bool = false

func main() {

	/*webcam, _ := gocv.VideoCaptureDevice(0)
	window := gocv.NewWindow("Hello")
	img := gocv.NewMat()

	for {
		webcam.Read(&img)
		window.IMShow(img)
		window.WaitKey(1)
	}*/

	//file := "main.json"
	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments, please pass in a json file")
		return
	}
	var fullscreen bool = false
	file := os.Args[1]
	if len(os.Args) > 2 {
		var err error
		fullscreen, err = strconv.ParseBool(os.Args[2])
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	in, err := ParseJsonToInputFile(file)
	//fmt.Println(in)
	if err != nil {
		fmt.Println("ERROR: issue parsing input file " + err.Error())
		return
	}

	runtime.LockOSThread()

	window := glInitFull(&in, fullscreen)
	defer glTerminate()

	program := initGLProgram(&in, file)

	var lastReloadTime float64 = 0

	// vao := makeVao(quad)
	var finished = false
	for !window.ShouldClose() && !finished {
		finished = glDraw(window, &program)
		checkInputs(window, &program, &in)

		// Hot reload
		elapsed := glfw.GetTime() - float64(lastReloadTime)
		if autoHotReload && elapsed > 1.0 {
			lastReloadTime = glfw.GetTime()
			regenerateShader(&program)
		}
	}

	glTerminate()
	CleanUp(&program)
}

func checkInputs(window *glfw.Window, program *OpenGLProgram, in *InputFile) {

	if window.GetKey(glfw.KeyR) == glfw.Press {
		// this is broken unsure why
		CleanUp(program)
		*program = initGLProgram(in, program.directory+program.jsonFileName)
		lastKey = window.GetKey(glfw.KeyR)
	} else if window.GetKey(glfw.KeyS) == glfw.Press {
		regenerateShader(program)
	} else if window.GetKey(glfw.KeyT) == glfw.Press {
		program.time = 0
	} else if window.GetKey(glfw.KeyC) == glfw.Press {

	} else if window.GetKey(glfw.KeyH) == glfw.Press {
		autoHotReload = !autoHotReload
		fmt.Println("Hot reload toggled " + strconv.FormatBool(autoHotReload))

	} else if window.GetKey(glfw.KeyF) == glfw.Press {
		// FULLSCREEN
		//goFullScreen(in.Width, in.Height, !globalDat.fullscreen)
		//globalDat.fullscreen = !globalDat.fullscreen

	} else if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
		// lastKey = window.GetKey(glfw.KeyP)
	} else if window.GetKey(glfw.KeyP) == glfw.Press {
		if window.GetKey(glfw.KeyP) != lastKey {
			fmt.Println("PAUSED")
			paused = !paused
			lastKey = window.GetKey(glfw.KeyP)
		}
	} else {
		lastKey = 0
	}

}
