package main

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Thank you sweet prince
// https://kylewbanks.com/blog/tutorial-opengl-with-golang-part-1-hello-opengl

const (
	vertexShaderSource = `
		#version 410
		layout (location = 0) in vec3 position;
		void main() {
			gl_Position = vec4(position, 1.0);
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
	width  = 1000
	height = 1000
	// triangle = []float32{
	// 	0, 0.5, 0,
	// 	-0.5, -0.5, 0,
	// 	0.5, -0.5, 0,
	// }
	quad = []float32{
		-1, 1, 0,
		-1, -1, 0,
		1, -1, 0,

		1, -1, 0,
		1, 1, 0,
		-1, 1, 0,
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
	vertexID   uint32
	fragmentID uint32
	// data       []GLData

	// This may not be best, oh well for now
	textures []uint32
	vbo      uint32
	vao      uint32

	videoData *VideoData
}

type GlobalGLData struct {
	fullscreen bool
	window     *glfw.Window
}

var globalDat GlobalGLData

// var video *vidio.Video
// var rgbaMain *image.RGBA

// var deleteShaders [uint32]

func glInit() *glfw.Window {

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
	window, err := glfw.CreateWindow(width, height, "Test", nil, nil)
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

func initGLProgram() OpenGLProgram {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vao, vbo := makeVao(quad)
	prog := gl.CreateProgram()

	fragSrc := getTextFromFile("test.frag")

	vertex, frag := createShaders(fragSrc, vertexShaderSource)
	gl.AttachShader(prog, vertex)
	gl.AttachShader(prog, frag)
	gl.LinkProgram(prog)

	gl.UseProgram(prog)

	textureUniform := gl.GetUniformLocation(prog, gl.Str("textureSampler\x00"))
	gl.Uniform1i(textureUniform, 0)
	texture := loadPictureAsTexture("test.png")

	vidData := setupVideo("testvid.mov")

	loc := gl.GetUniformLocation(prog, gl.Str("res\x00"))
	gl.Uniform2f(loc, float32(width), float32(height))
	fmt.Printf("LOC %d", loc)

	// gl.bindfra(prog, 0, gl.Str("outputColor\x00"))

	return OpenGLProgram{
		programID:  prog,
		vertexID:   vertex,
		fragmentID: frag,
		textures:   []uint32{texture},
		vao:        vao,
		vbo:        vbo,
		videoData:  vidData,
		//data:       []GLData{GLData{id: texture, dataType: T_TEXTURE}},
	}
}

func glDraw(window *glfw.Window, program OpenGLProgram) {

	gl.UseProgram(program.programID)

	time := glfw.GetTime()
	elapsed := time - previousTime
	previousTime = time

	// updateVideo(time, video)
	updateVideo(time, program.videoData)

	gl.Uniform1f(gl.GetUniformLocation(program.programID, gl.Str("iTime\x00")), float32(time))
	gl.Uniform1f(gl.GetUniformLocation(program.programID, gl.Str("deltaTime\x00")), float32(elapsed))

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	//glSetShaderData(program)

	gl.BindVertexArray(program.vao)

	// TODO: bind all textures
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, program.textures[0])

	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(quad)/3))

	writeData(program.videoData)

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
