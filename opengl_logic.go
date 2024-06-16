package main

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"gocv.io/x/gocv"
)

// Thank you sweet prince
// https://kylewbanks.com/blog/tutorial-opengl-with-golang-part-1-hello-opengl

const (
	vertexShaderSource = `
		#version 410
		// attribute vec2 vertexIn;
		layout (location = 0) in vec3 position;

		out vec2 uv; 

		void main() {
			uv = (position.xy * 0.5) + 0.5;
			// uv.y = 1.0 - uv.y; 
			gl_Position = vec4(position.x, -position.y, 0.0, 1.0);
		}
	` + "\x00"

	/*fragmentShaderSource = `
		#version 410
		uniform sampler2D textureSampler;
		out vec4 fragColor;

		uniform vec2 res;

		void main() {
			fragColor = texture(textureSampler, vec2(gl_FragCoord.x/res.x / 2.0, 1.0 - (gl_FragCoord.y/res.y/2.0)));
			//fragColor = vec4(float(gl_FragCoord.x)/res.x, float(gl_FragCoord.y)/res.y, 1.0, 1.0);//vec4(0.0, 0.0, 1.0, 1.0); //vec4(gl_FragCoord.x/res.x, gl_FragCoord.y/res.y, gl_FragCoord.z/500.f, 1.0);
		}
	` + "\x00"*/
)

var (
	//width  = 1000
	//height = 1000
	// triangle = []float32{
	// 	0, 0.5, 0,
	// 	-0.5, -0.5, 0,
	// 	0.5, -0.5, 0,
	// }
	/*quad = []float32{
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,

		0.5, -0.5, 0,
		0.5, 0.5, 0,
		-0.5, 0.5, 0,
	}*/
	quad = []float32{

		1, -1, 0,
		1, 1, 0,
		-1, 1, 0,

		-1, 1, 0,
		-1, -1, 0,
		1, -1, 0,
	}

	previousTime float64

	// fragmentShaderSource = "\x00"
)

// type GLDataType uint8

// const T_TEXTURE GLDataType = 1
// const T_VAO GLDataType = 2
// const T_VBO GLDataType = 3

// type GLData struct {
// 	id       uint32
// 	dataType GLDataType
// }

type OpenGLProgram struct {
	programID  uint32
	fileName   string
	vertexID   uint32
	fragmentID uint32
	// data       []GLData

	// This may not be best, oh well for now
	textures []uint32
	vbo      uint32
	vao      uint32

	videos []*VideoData

	recordFPS     int32
	timesRendered uint64
	width         int
	height        int
	folder        string
}

type GlobalGLData struct {
	fullscreen bool
	window     *glfw.Window
}

var globalDat GlobalGLData

// var video *vidio.Video
// var rgbaMain *image.RGBA

// var deleteShaders [uint32]

func glInit(in *InputFile) *glfw.Window {

	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// monitor := glfw.GetPrimaryMonitor() // The primary monitor.. Later Occulus?..

	// mode := monitor.GetVideoMode()
	// width = //mode.Width
	// height = //mode.Height
	window, err := glfw.CreateWindow(in.Width, in.Height, "Test", nil, nil)
	globalDat.window = window

	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	return window

}

func glTerminate() {
	glfw.Terminate()
}

func initGLProgram(in *InputFile) OpenGLProgram {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vao, vbo := makeVao(quad)
	prog := gl.CreateProgram()

	prefix := in.Folder + "/"

	fragSrc, err := getTextFromFile(prefix + in.ShaderPath)

	if err != nil {
		panic("Error getting file ")
	}

	vertex, frag := createShaders(fragSrc, vertexShaderSource)
	gl.AttachShader(prog, vertex)
	gl.AttachShader(prog, frag)
	gl.LinkProgram(prog)

	gl.UseProgram(prog)

	//gl.Enable(gl.BLEND)
	//gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	retProg := OpenGLProgram{
		programID:  prog,
		fileName:   in.ShaderPath,
		vertexID:   vertex,
		fragmentID: frag,
		textures:   []uint32{},
		vao:        vao,
		vbo:        vbo,
		videos:     []*VideoData{},
		recordFPS:  in.RecordFPS,
		width:      in.Width,
		height:     in.Height,
		//data:       []GLData{GLData{id: texture, dataType: T_TEXTURE}},
	}

	LoadOpenGLDataFromInputFile(&retProg, in)

	return retProg
}

func LoadOpenGLDataFromInputFile(prog *OpenGLProgram, input *InputFile) {
	var newTextures []uint32
	var newVideos []*VideoData

	videoWriter = setupVideoWriter(prog)
	mat := gocv.NewMatWithSize(prog.height, prog.width, gocv.MatTypeCV8SC3)
	writerData = &VideoData{
		// video: video,
		// writer:       nil,
		fps: float64(prog.recordFPS),
		// frames:       int(video.Get(gocv.VideoCaptureFrameCount)),
		width:        prog.width,
		height:       prog.height,
		currentFrame: -1,
		material:     &mat,
		// materials:       []gocv.Mat{gocv.Mat{}, gocv.Mat{}, gocv.Mat{}},
		// currentMatIndex: 0,
	}

	gl.BindVertexArray(prog.vao)

	loc := gl.GetUniformLocation(prog.programID, gl.Str("res\x00"))
	gl.Uniform2f(loc, float32(input.Width), float32(input.Height))

	for i, texturePath := range input.Textures {

		isWebcam := (texturePath == "WEBCAM")

		texturePath = input.Folder + "/" + texturePath

		isPhoto := strings.HasSuffix(texturePath, ".jpg") || strings.HasSuffix(texturePath, ".png") || strings.HasSuffix(texturePath, ".jpeg")
		isVideo := strings.HasSuffix(texturePath, ".webm") || strings.HasSuffix(texturePath, ".mov") || strings.HasSuffix(texturePath, ".aiff") || strings.HasSuffix(texturePath, ".mp4") || strings.HasSuffix(texturePath, ".mpeg")

		if !isPhoto && !isVideo && !isWebcam {
			fmt.Println("ERROR, invalid file format for " + texturePath)
			continue
		}

		gl.ActiveTexture(gl.TEXTURE0 + uint32(i))

		var texture uint32

		if isPhoto {
			texture = loadPictureAsTexture(texturePath)
		} else if isVideo || isWebcam {
			var vidData *VideoData

			//*writerData.material = gocv.NewMat()

			if isVideo {
				texture, vidData = setupVideo(texturePath)
				fmt.Printf("NULL VIDEO: %t", vidData.video == nil)
				vidData.ReadAllFrames()
				fmt.Printf("%d", len(vidData.allFrames))
			} else {
				texture, vidData = setupVideo("WEBCAM")
			}

			// IT IS NOT RECORDING RIGHT
			//vidData.video.Set(gocv.VideoCapturePosFrames, 0)
			//vidData.video.Read(&writerData.material)

			newVideos = append(newVideos, vidData)
		}
		newTextures = append(newTextures, texture)

		// TODO: What if out of range? More than 32 textures?
		// for _, text := range newTextures {

		// gl.BindTexture(gl.TEXTURE_2D, texture)

		str := fmt.Sprintf("tex%d", i)
		textureUniform := gl.GetUniformLocation(prog.programID, gl.Str(str+"\x00"))
		gl.Uniform1i(textureUniform, int32(i))
		fmt.Printf("TEXTURE %d %d %s\n", texture, i, str)
		// }

	}

	prog.textures = newTextures
	prog.videos = newVideos
}

func glDraw(window *glfw.Window, program OpenGLProgram) {

	gl.UseProgram(program.programID)

	time := glfw.GetTime()
	elapsed := time - previousTime
	previousTime = time

	// updateVideo(time, video)
	for _, vid := range program.videos {
		updateVideo(time, vid)
	}

	// fmt.Printf("FPS: %f", 1.0/elapsed)

	if program.recordFPS < 0 {
		gl.Uniform1f(gl.GetUniformLocation(program.programID, gl.Str("iTime\x00")), float32(time))
		gl.Uniform1f(gl.GetUniformLocation(program.programID, gl.Str("deltaTime\x00")), float32(elapsed))
	} else {
		gl.Uniform1f(gl.GetUniformLocation(program.programID, gl.Str("iTime\x00")), float32(program.timesRendered)/float32(program.recordFPS))
		gl.Uniform1f(gl.GetUniformLocation(program.programID, gl.Str("deltaTime\x00")), 1.0/float32(program.recordFPS))
	}

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	//glSetShaderData(program)

	gl.BindVertexArray(program.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(quad)/3))

	if program.recordFPS > 0 && program.timesRendered == 0 {
		// for _, vid := range program.videos {
		writeData(writerData)
		// }
	}

	program.timesRendered++

	glfw.PollEvents()
	window.SwapBuffers()
}

/*func glSetShaderData(dat uint32) {
	res_str, free := gl.Strs("res")
	gl.Uniform2i(gl.GetUniformLocation(dat, *res_str), width, height)
	free()
}*/

// FREEING FUNCTIONS

func CleanUp(prog *OpenGLProgram) {
	// May need to use for loops?

	for _, texture := range prog.textures {
		//gl.DeleteTextures(int32(len(prog.textures)), &prog.textures[0])
		gl.DeleteTextures(1, &texture)
	}

	for _, vid := range prog.videos {
		vid.material.Close()
		vid.video.Close()
	}
	gl.DeleteVertexArrays(1, &prog.vao)
	gl.DeleteBuffers(1, &prog.vbo)
	gl.DeleteProgram(prog.programID)
}

//func genBuffers
