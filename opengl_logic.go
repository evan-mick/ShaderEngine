package main

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"path"
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
)

var (
	quad = []float32{

		1, -1, 0,
		1, 1, 0,
		-1, 1, 0,

		-1, 1, 0,
		-1, -1, 0,
		1, -1, 0,
	}

	previousTime float64
)

type OpenGLProgram struct {
	programID      uint32
	shaderFileName string
	jsonFileName   string
	directory      string

	vertexID   uint32
	fragmentID uint32
	// data       []GLData

	// This may not be best, oh well for now
	textures []uint32
	vbo      uint32
	vao      uint32

	videos []*VideoData

	recordFPS     int32
	recordSeconds int64
	time          float64
	timesRendered uint64
	width         int
	height        int
	folder        string
	filePrefix    string
}

type GlobalGLData struct {
	fullscreen bool
	window     *glfw.Window
	renderFBO  uint32
}

var globalDat GlobalGLData

func glInit(in *InputFile) *glfw.Window {
	return glInitFull(in, false)
}

func glInitFull(in *InputFile, fullscreen bool) *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	return glSetupNewWindow(in.Width, in.Height, fullscreen)
}

func glSetupNewWindow(width int, height int, fullscreen bool) *glfw.Window {

	if globalDat.window != nil {
		globalDat.window.SetShouldClose(true)
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	monitor := glfw.GetPrimaryMonitor()
	if !fullscreen {
		monitor = nil
	}

	window, err := glfw.CreateWindow(width, height, "Shader Engine", monitor, nil)

	globalDat.window = window

	if err != nil {
		panic(err)
	}
	globalDat.fullscreen = fullscreen
	window.MakeContextCurrent()

	return window
}

func glTerminate() {
	glfw.Terminate()
}

func initGLProgram(in *InputFile, full_filepath string) OpenGLProgram {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	vao, vbo := makeQuadVaoVbo(quad)
	prog := gl.CreateProgram()

	dir, filename := path.Split(full_filepath)

	prefix := dir
	folder := in.Folder + "/"

	fragSrc, err := getTextFromFile(prefix + folder + in.ShaderPath)

	if err != nil {
		panic("Error getting file ")
	}

	vertex, frag := createShaders(fragSrc, vertexShaderSource)
	gl.AttachShader(prog, vertex)
	gl.AttachShader(prog, frag)
	gl.LinkProgram(prog)

	gl.UseProgram(prog)

	if in.RecordFPS > 0 {
		var fbo uint32 = 0
		gl.GenFramebuffers(1, &fbo)
		gl.BindFramebuffer(gl.FRAMEBUFFER, fbo)
		globalDat.renderFBO = fbo

		var texture uint32 = 0
		gl.GenTextures(1, &texture)
		gl.BindTexture(gl.TEXTURE_2D, texture)
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(in.Width), int32(in.Height), 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
		gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, texture, 0)
	}

	retProg := OpenGLProgram{
		programID:      prog,
		shaderFileName: in.ShaderPath,
		jsonFileName:   filename,
		directory:      prefix,
		folder:         folder,

		vertexID:      vertex,
		fragmentID:    frag,
		textures:      []uint32{},
		vao:           vao,
		vbo:           vbo,
		videos:        []*VideoData{},
		recordFPS:     in.RecordFPS,
		time:          0,
		width:         in.Width,
		height:        in.Height,
		timesRendered: 0,
		//data:       []GLData{GLData{id: texture, dataType: T_TEXTURE}},
	}

	LoadOpenGLDataFromInputFile(&retProg, in)
	glfw.SetTime(0)

	return retProg
}

func regenerateShader(program *OpenGLProgram) {
	prog := program.programID

	fragSrc, err := getTextFromFile(program.directory + program.folder + program.shaderFileName)

	gl.DetachShader(prog, program.vertexID)
	gl.DetachShader(prog, program.fragmentID)

	if err != nil {
		panic("Error getting file ")
	}

	vertex, frag := createShaders(fragSrc, vertexShaderSource)

	program.vertexID = vertex
	program.fragmentID = frag

	gl.UseProgram(prog)
	gl.AttachShader(prog, vertex)
	gl.AttachShader(prog, frag)
	gl.LinkProgram(prog)
	program.time = 0
}

func LoadOpenGLDataFromInputFile(prog *OpenGLProgram, input *InputFile) {
	var newTextures []uint32
	var newVideos []*VideoData

	// Only write if we are trying to record
	if prog.recordFPS > 0 {
		videoWriter = setupVideoWriter(prog)
	}
	mat := gocv.NewMatWithSize(prog.height, prog.width, gocv.MatTypeCV8SC3)
	// Originally wanted to use multiple materials for videos to potentially make it quicker
	// For whatever reason, could not get it to work, may revisit
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

	type TextureType uint8
	const (
		INVALID TextureType = 0
		PHOTO               = 1
		VIDEO               = 2
		WEBCAM              = 3
	)

	// TODO: What if out of range? More than 32 textures?
	for i, fileTexturePath := range input.Textures {

		textureType := INVALID

		// Add the file prefix thing?
		// Want to be able to run json from anywhere and have the whole thing just work
		texturePath := /*prog.filePrefix + "/" +*/ prog.directory + input.Folder + "/" + fileTexturePath

		if fileTexturePath == "WEBCAM" {
			textureType = WEBCAM
		} else if strings.HasSuffix(texturePath, ".jpg") || strings.HasSuffix(texturePath, ".png") || strings.HasSuffix(texturePath, ".jpeg") || strings.HasSuffix(texturePath, ".mkv") {
			textureType = PHOTO
		} else if strings.HasSuffix(texturePath, ".webm") || strings.HasSuffix(texturePath, ".mov") || strings.HasSuffix(texturePath, ".aiff") || strings.HasSuffix(texturePath, ".mp4") || strings.HasSuffix(texturePath, ".mpeg") {
			textureType = VIDEO
		}

		gl.ActiveTexture(gl.TEXTURE0 + uint32(i))
		var texture uint32
		var vidData *VideoData

		switch textureType {
		case INVALID:
			fmt.Println("ERROR, invalid file format for " + texturePath)
			continue
		case PHOTO:
			texture = loadPictureAsTexture(texturePath)
		case VIDEO:
			texture, vidData = setupVideo(texturePath)
			vidData.ReadAllFrames()
			newVideos = append(newVideos, vidData)
		case WEBCAM:
			texture, vidData = setupVideo("WEBCAM")
			newVideos = append(newVideos, vidData)
		}
		newTextures = append(newTextures, texture)

		//vidData.video.Set(gocv.VideoCapturePosFrames, 0)
		//vidData.video.Read(&writerData.material)

		str := fmt.Sprintf("tex%d", i)
		textureUniform := gl.GetUniformLocation(prog.programID, gl.Str(str+"\x00"))
		gl.Uniform1i(textureUniform, int32(i))
	}

	prog.textures = newTextures
	prog.videos = newVideos
}

func glDraw(window *glfw.Window, program *OpenGLProgram) {

	// gl.Viewport(0, 0, int32(program.width*2), int32(program.height*2))
	gl.UseProgram(program.programID)

	if !paused {
		cur_time := glfw.GetTime()
		elapsed := cur_time - previousTime
		previousTime = cur_time

		if program.recordFPS < 0 {
			gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
			w, h := window.GetFramebufferSize()
			gl.Viewport(0, 0, int32(w), int32(h))
			gl.Uniform1f(gl.GetUniformLocation(program.programID, gl.Str("deltaTime\x00")), float32(elapsed))
		} else {
			// gl.BindFramebuffer(gl.FRAMEBUFFER, globalDat.renderFBO)
			gl.BindFramebuffer(gl.FRAMEBUFFER, globalDat.renderFBO)
			gl.Viewport(0, 0, int32(program.width), int32(program.height))
			program.time = float64(program.timesRendered) / float64(program.recordFPS)
			fmt.Printf("time %f %d %d\n", program.time, program.recordFPS, program.timesRendered)
			gl.Uniform1f(gl.GetUniformLocation(program.programID, gl.Str("deltaTime\x00")), 1.0/float32(program.recordFPS))
		}
		gl.Uniform1f(gl.GetUniformLocation(program.programID, gl.Str("iTime\x00")), float32(program.time))

		// updateVideo(time, video)
		for _, vid := range program.videos {
			updateVideo(program.time, vid)
			fmt.Println(vid.material.Type().String())
		}

		// fmt.Printf("FPS: %f", 1.0/elapsed)

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		//glSetShaderData(program)

		gl.BindVertexArray(program.vao)
		// gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
		// gl.DrawArrays(gl.TRIANGLES, 0, int32(len(quad)/3))

		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(quad)/3))

		gl.BindFramebuffer(gl.FRAMEBUFFER, globalDat.renderFBO)

		if program.recordFPS > 0 {
			fWidth, fHeight := window.GetFramebufferSize()
			// width, height := window.GetSize()
			// fmt.Println(window.GetFramebufferSize())
			writeData(int32(fWidth), int32(fHeight), int32(program.width), int32(program.height))
			// writeData()
		}

		program.timesRendered++
		program.time += elapsed
		window.SwapBuffers()
	}

	glfw.PollEvents()

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

	if videoWriter != nil {
		videoWriter.Close()
	}
}

//func genBuffers
