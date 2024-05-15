package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"log"
	"os"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Thank you sweet prince
// https://kylewbanks.com/blog/tutorial-opengl-with-golang-part-1-hello-opengl

const (
	width  = 1000
	height = 1000

	vertexShaderSource = `
		#version 410
		layout (location = 0) in vec3 position;
		void main() {
			gl_Position = vec4(position, 1.0);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 410
		uniform sampler2D textureSampler;
		out vec4 fragColor;

		uniform vec2 res; 

		void main() {
			fragColor = texture(textureSampler, vec2(gl_FragCoord.x/res.x / 2.0, 1.0 - (gl_FragCoord.y/res.y/2.0)));
			//fragColor = vec4(float(gl_FragCoord.x)/res.x, float(gl_FragCoord.y)/res.y, 1.0, 1.0);//vec4(0.0, 0.0, 1.0, 1.0); //vec4(gl_FragCoord.x/res.x, gl_FragCoord.y/res.y, gl_FragCoord.z/500.f, 1.0);
		}
	` + "\x00"
)

var (
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

	// fragmentShaderSource = "\x00"
)

var textures []uint32
var cur_text uint32

// var deleteShaders [uint32]

func glInit() *glfw.Window {

	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Test", nil, nil)

	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	return window

}

func glTerminate() {
	glfw.Terminate()
}

func createShaders() (vertexShader uint32, fragmentShader uint32) {

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err = compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	return vertexShader, fragmentShader
}

func initGLProgram() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	prog := gl.CreateProgram()

	vertex, frag := createShaders()
	gl.AttachShader(prog, vertex)
	gl.AttachShader(prog, frag)
	gl.LinkProgram(prog)

	gl.UseProgram(prog)

	textureUniform := gl.GetUniformLocation(prog, gl.Str("textureSampler\x00"))
	gl.Uniform1i(textureUniform, 0)
	cur_text = loadPictureAsTexture("test.png")

	loc := gl.GetUniformLocation(prog, gl.Str("res\x00"))
	gl.Uniform2f(loc, float32(width), float32(height))
	fmt.Printf("LOC %d", loc)

	// gl.bindfra(prog, 0, gl.Str("outputColor\x00"))

	return prog
}

func glDraw(vao uint32, window *glfw.Window, program uint32) {

	gl.UseProgram(program)

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	//glSetShaderData(program)

	gl.BindVertexArray(vao)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, cur_text)

	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(quad)/3))

	glfw.PollEvents()
	window.SwapBuffers()
}

/*func glSetShaderData(dat uint32) {
	res_str, free := gl.Strs("res")
	gl.Uniform2i(gl.GetUniformLocation(dat, *res_str), width, height)
	free()
}*/

// Compiles shader with open GL
// returns shader id, on error shader id is 0 and the error has the open gl message
func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)

	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))

		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}
	return shader, nil

}

func makeVao(points []float32) uint32 {

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	return vao
}

// copied from new texture function
// https://github.com/go-gl/example/blob/master/gl21-cube/cube.go
func loadPictureAsTexture(file string) uint32 {

	imgFile, err := os.Open(file)
	if err != nil {
		fmt.Println("Texture error at path " + file + " " + err.Error())
		return 0
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		fmt.Println("Image decode error at path " + file + " " + err.Error())
		return 0
	}
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		if err != nil {
			fmt.Println("Texture error at path " + file + " Unsupported stride")
			return 0
		}
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	textures = append(textures, texture)

	return texture

}

// FREEING FUNCTIONS

func CleanUp() {
	// May need to use for loops?
	gl.DeleteTextures(int32(len(textures)), &textures[0])
}

//func genBuffers
