package main

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
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

	fragSrc, err := getTextFromFile(in.ShaderPath)

	if err != nil {
		panic("Error getting file ")
	}

	vertex, frag := createShaders(fragSrc, vertexShaderSource)
	gl.AttachShader(prog, vertex)
	gl.AttachShader(prog, frag)
	gl.LinkProgram(prog)

	gl.UseProgram(prog)

	retProg := OpenGLProgram{
		programID:  prog,
		fileName:   in.ShaderPath,
		vertexID:   vertex,
		fragmentID: frag,
		textures:   []uint32{},
		vao:        vao,
		vbo:        vbo,
		videos:     []*VideoData{},
		//data:       []GLData{GLData{id: texture, dataType: T_TEXTURE}},
	}

	LoadOpenGLDataFromInputFile(&retProg, in)

	return retProg
}

func LoadOpenGLDataFromInputFile(prog *OpenGLProgram, input *InputFile) {
	var newTextures []uint32
	var newVideos []*VideoData

	gl.BindVertexArray(prog.vao)

	loc := gl.GetUniformLocation(prog.programID, gl.Str("res\x00"))
	gl.Uniform2f(loc, float32(input.Width), float32(input.Height))

	for i, texturePath := range input.Textures {
		isPhoto := strings.HasSuffix(texturePath, ".jpg") || strings.HasSuffix(texturePath, ".png") || strings.HasSuffix(texturePath, ".jpeg")
		isVideo := strings.HasSuffix(texturePath, ".mov") || strings.HasSuffix(texturePath, ".aiff") || strings.HasSuffix(texturePath, ".mp4") || strings.HasSuffix(texturePath, ".mpeg")

		if !isPhoto && !isVideo {
			fmt.Println("ERROR, invalid file format for " + texturePath)
			continue
		}

		var texture uint32

		if isPhoto {
			texture = loadPictureAsTexture(texturePath)
		} else if isVideo {
			var vidData *VideoData
			texture, vidData = setupVideo(texturePath)
			newVideos = append(newVideos, vidData)
		}
		newTextures = append(newTextures, texture)

		// TODO: What if out of range? More than 32 textures?
		// for _, text := range newTextures {
		gl.ActiveTexture(gl.TEXTURE0 + uint32(i))
		gl.BindTexture(gl.TEXTURE_2D, texture)

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

	gl.Uniform1f(gl.GetUniformLocation(program.programID, gl.Str("iTime\x00")), float32(time))
	gl.Uniform1f(gl.GetUniformLocation(program.programID, gl.Str("deltaTime\x00")), float32(elapsed))

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	//glSetShaderData(program)

	gl.BindVertexArray(program.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(quad)/3))

	for _, vid := range program.videos {
		writeData(vid)
	}

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
	gl.DeleteTextures(int32(len(prog.textures)), &prog.textures[0])
	gl.DeleteVertexArrays(1, &prog.vao)
	gl.DeleteBuffers(1, &prog.vbo)
	gl.DeleteProgram(prog.programID)
}

//func genBuffers
