package main

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"path"
	"strconv"
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
		layout (location = 0) in vec3 position;

		out vec2 uv; 

		void main() {
			uv = (position.xy * 0.5) + 0.5;
			gl_Position = vec4(position.x, -position.y, 0.0, 1.0);
		}
	` + "\x00"
)

var (
	quad = []float32{
		//-0.5, 0.5, 0,
		//0.5, -0.5, 0,
		//0.5, 0.5, 0,

		//-0.5, 0.5, 0,
		//-0.5, -0.5, 0,
		//0.5, -0.5, 0,
		-1, 1, 0,
		-1, -1, 0,
		1, -1, 0,

		1, -1, 0,
		1, 1, 0,
		-1, 1, 0,
	}

	previousTime float64
)

type Channel struct {
	shaderFileName string

	programID  uint32
	vertexID   uint32
	fragmentID uint32
	textures   []uint32
	vbo        uint32
	vao        uint32
	fbo        uint32
	fboTexture uint32
}

type OpenGLProgram struct {
	jsonFileName string
	directory    string

	includes []string

	mainChannel   *Channel
	extraChannels []*Channel

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
	//renderFBO  uint32
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

func getFullFragSrc(program *OpenGLProgram, shaderPath string) string {
	fullSrc := "#version 410\n"

	fragSrc, err := getTextFromFile(program.directory + program.folder + shaderPath)

	if err != nil {
		panic("Error getting file " + err.Error())
	}

	// includes
	for _, includePath := range program.includes {
		includeSrc, err := getTextFromFile(program.directory + program.folder + includePath)
		fmt.Println("INCLUDE " + includePath)

		if err != nil {
			fmt.Println("Error with include: " + err.Error())
			continue
		}

		fullSrc += includeSrc
	}
	fullSrc += fragSrc
	fullSrc += "\x00"

	return fullSrc
}

func initializeChannel(program *OpenGLProgram, channelData ChannelJson) Channel {

	vao, vbo := makeQuadVaoVbo(quad)
	prog := gl.CreateProgram()

	fullSrc := getFullFragSrc(program, channelData.ShaderPath)

	fmt.Print(fullSrc)

	vertex, frag := createShaders(fullSrc, vertexShaderSource)
	gl.AttachShader(prog, vertex)
	gl.AttachShader(prog, frag)
	gl.LinkProgram(prog)
	gl.UseProgram(prog)

	var fbo uint32 = 0
	//if program.recordFPS > 0 {
	gl.GenFramebuffers(1, &fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, fbo)
	//globalDat.renderFBO = fbo

	var texture uint32 = 0
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(program.width), int32(program.height), 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, texture, 0)
	//}

	retChannel := Channel{
		programID:      prog,
		vertexID:       vertex,
		fragmentID:     frag,
		fbo:            fbo,
		fboTexture:     texture,
		shaderFileName: channelData.ShaderPath,
		vao:            vao,
		vbo:            vbo,
		textures:       []uint32{},
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.UseProgram(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	return retChannel
}

func initGLProgram(in *InputFile, full_filepath string) OpenGLProgram {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	dir, filename := path.Split(full_filepath)

	prefix := dir
	folder := in.Folder + "/"

	glfw.SetTime(0)
	retProg := OpenGLProgram{
		jsonFileName: filename,
		directory:    prefix,
		folder:       folder,
		//renderFBO:    fbo,

		mainChannel:   nil,
		extraChannels: []*Channel{},

		videos:        []*VideoData{},
		recordFPS:     in.RecordFPS,
		recordSeconds: in.RecordSeconds,
		time:          0,
		width:         in.Width,
		height:        in.Height,
		timesRendered: 0,
		includes:      in.Includes,
		//data:       []GLData{GLData{id: texture, dataType: T_TEXTURE}},
	}
	// But actually, unless we're recording, main channel should go straight to window (CHANGE FBO)
	mainChanData := ChannelJson{ShaderPath: in.ShaderPath, Textures: in.Textures}
	mainChan := initializeChannel(&retProg, mainChanData)
	retProg.mainChannel = &mainChan

	if retProg.recordFPS < 0 {
		retProg.mainChannel.fbo = 0
	}

	for _, info := range in.Channels {
		newChan := initializeChannel(&retProg, info)
		retProg.extraChannels = append(retProg.extraChannels, &newChan)
	}

	// Have to load data on second pass, because channel textures haven't been made yet before
	LoadOpenGLDataFromInputFile(&retProg, &mainChanData, retProg.mainChannel)
	for i, _ := range in.Channels {
		LoadOpenGLDataFromInputFile(&retProg, &(in.Channels[i]), retProg.extraChannels[i])
	}

	return retProg
}

func safeReloadChannelProgram(program *OpenGLProgram, channel *Channel) {

	//fragSrc, err := getTextFromFile(program.directory + program.folder + channel.shaderFileName)
	fragSrc := getFullFragSrc(program, channel.shaderFileName)

	//if err != nil || !fragShaderCompilable(fragSrc) {
	//	panic("Program directory failure in safe reload (THIS SHOUlD NEVER HAPPEN!)")
	//}

	newFragmentShader, err := compileShader(fragSrc, gl.FRAGMENT_SHADER)
	if err != nil {
		fmt.Println("Frag shader compilation error: " + err.Error())
		return
	}
	gl.DetachShader(channel.programID, channel.fragmentID)
	gl.AttachShader(channel.programID, newFragmentShader)
	gl.LinkProgram(channel.programID)

	channel.fragmentID = newFragmentShader
}

func regenerateShader(program *OpenGLProgram) {
	safeReloadChannelProgram(program, program.mainChannel)

	for _, channel := range program.extraChannels {
		safeReloadChannelProgram(program, channel)
	}
}

func LoadOpenGLDataFromInputFile(prog *OpenGLProgram, channelData *ChannelJson, outChannel *Channel) {
	// TODO: video texture caching
	// If two channels have the same vid, they should use the same texture id

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
		fps: float64(prog.recordFPS),
		// frames:       int(video.Get(gocv.VideoCaptureFrameCount)),
		width:        prog.width,
		height:       prog.height,
		currentFrame: -1,
		material:     &mat,
		// materials:       []gocv.Mat{gocv.Mat{}, gocv.Mat{}, gocv.Mat{}},
		// currentMatIndex: 0,
	}

	//gl.BindVertexArray(prog.vao)

	type TextureType uint8
	const (
		INVALID TextureType = 0
		PHOTO               = 1
		VIDEO               = 2
		WEBCAM              = 3
		CHANNEL             = 4
	)

	gl.ActiveTexture(gl.TEXTURE0)
	// TODO: What if out of range? More than 32 textures?
	for _, fileTexturePath := range channelData.Textures {

		textureType := INVALID

		// Add the file prefix thing?
		// Want to be able to run json from anywhere and have the whole thing just work
		texturePath := /*prog.filePrefix + "/" +*/ prog.directory + prog.folder + "/" + fileTexturePath
		chanVal, chanErr := strconv.Atoi(fileTexturePath)

		if fileTexturePath == "WEBCAM" {
			textureType = WEBCAM
		} else if strings.HasSuffix(texturePath, ".jpg") || strings.HasSuffix(texturePath, ".png") || strings.HasSuffix(texturePath, ".jpeg") || strings.HasSuffix(texturePath, ".mkv") {
			textureType = PHOTO
		} else if strings.HasSuffix(texturePath, ".webm") || strings.HasSuffix(texturePath, ".mov") || strings.HasSuffix(texturePath, ".aiff") || strings.HasSuffix(texturePath, ".mp4") || strings.HasSuffix(texturePath, ".mpeg") {
			textureType = VIDEO
		} else if chanErr == nil || fileTexturePath == "MAIN" {
			textureType = CHANNEL
		}

		// TODO: some smarter active texture stuff
		// Index of all the textures that have been loaded thus far
		// Then just iteratively add active textures everytime this function is called
		// Then smartly assign the uniforms to the right textures based on channel location n stuff
		// for now, too bad!
		//gl.ActiveTexture(gl.TEXTURE0 + uint32(i))
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
		case CHANNEL:
			if fileTexturePath == "MAIN" {
				texture = prog.mainChannel.fboTexture
			} else if chanVal < len(prog.extraChannels) {
				texture = prog.extraChannels[chanVal].fboTexture
			}
		}
		newTextures = append(newTextures, texture)

		//vidData.video.Set(gocv.VideoCapturePosFrames, 0)
		//vidData.video.Read(&writerData.material)

		//str := fmt.Sprintf("tex%d", i)
		//textureUniform := gl.GetUniformLocation(prog.programID, gl.Str(str+"\x00"))
		//gl.Uniform1i(textureUniform, int32(i))
	}

	outChannel.textures = newTextures

	for _, vid := range newVideos {
		prog.videos = append(prog.videos, vid)
	}
}

func drawChannel(program *OpenGLProgram, channel *Channel, width int32, height int32, elapsed float32) {

	if channel == nil {
		panic("nil channel in draw")
	}

	//fmt.Printf("drawing channel: %s %f\n", channel.shaderFileName, program.time)

	gl.BindFramebuffer(gl.FRAMEBUFFER, channel.fbo)
	gl.UseProgram(channel.programID)
	gl.Viewport(0, 0, width, height)
	gl.Uniform1f(gl.GetUniformLocation(channel.programID, gl.Str("deltaTime\x00")), elapsed)
	gl.Uniform1f(gl.GetUniformLocation(channel.programID, gl.Str("iTime\x00")), float32(program.time))

	loc := gl.GetUniformLocation(channel.programID, gl.Str("res\x00"))
	gl.Uniform2f(loc, float32(program.width), float32(program.height))

	// GO THROUGH TEXTURES AND ACTIVATE THEM
	for i, _ := range channel.textures {
		gl.ActiveTexture(gl.TEXTURE0 + uint32(i))
		gl.BindTexture(gl.TEXTURE_2D, channel.textures[i])
		str := fmt.Sprintf("tex%d", i)
		textureUniform := gl.GetUniformLocation(channel.programID, gl.Str(str+"\x00"))
		gl.Uniform1i(textureUniform, int32(i))
		//fmt.Printf("TEXTURES: %d\n", channel.textures[i])
	}

	// Draw to screen or to rendering
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.BindVertexArray(channel.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(quad)/3))

	//gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

}

func glDraw(window *glfw.Window, program *OpenGLProgram) bool {

	isLiveVideo := program.recordFPS < 0

	if !isLiveVideo && program.time >= float64(program.recordSeconds) {
		return true
	}

	if paused {
		// Still need to track input, namely for unpausing
		glfw.PollEvents()
		return false
	}

	var elapsed float64
	var w, h int

	// Update appropriate variables
	if isLiveVideo {
		elapsed := glfw.GetTime() - previousTime
		previousTime = glfw.GetTime()
		program.time += elapsed

		w, h = window.GetFramebufferSize()
	} else {
		program.time = float64(program.timesRendered) / float64(program.recordFPS)
		elapsed = 1.0 / float64(program.recordFPS)

		fmt.Printf("time %f %d %d\n", program.time, program.recordFPS, program.timesRendered)

		w = program.width
		h = program.height
	}

	for _, vid := range program.videos {
		updateVideo(program.time, vid)
		//fmt.Println(vid.material.Type().String())
	}

	for _, channel := range program.extraChannels {
		drawChannel(program, channel, int32(program.width), int32(program.height), float32(elapsed))
	}
	drawChannel(program, program.mainChannel, int32(w), int32(h), float32(elapsed))

	if !isLiveVideo {
		//fWidth, fHeight := window.GetFramebufferSize()
		writeData(int32(program.width), int32(program.height))
	}

	program.timesRendered++
	window.SwapBuffers()
	glfw.PollEvents()
	return false
}

func cleanChannel(channel *Channel) {
	gl.DeleteProgram(channel.programID)

	for _, tex := range channel.textures {
		gl.DeleteTextures(1, &tex)
	}
	gl.DeleteFramebuffers(1, &channel.fbo)
	gl.DeleteVertexArrays(1, &channel.vao)
	gl.DeleteBuffers(1, &channel.vbo)
	// Potential issue: double delete the same fbo texture? bc that's what's accessed from other channels
	gl.DeleteTextures(1, &channel.fboTexture)
}

// FREEING FUNCTIONS
func CleanUp(prog *OpenGLProgram) {

	cleanChannel(prog.mainChannel)
	for _, channel := range prog.extraChannels {
		cleanChannel(channel)
	}

	for _, vid := range prog.videos {
		vid.material.Close()
		vid.video.Close()
	}

	if videoWriter != nil {
		videoWriter.Close()
	}
}
