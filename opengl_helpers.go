package main

import (
	"fmt"
	"image"
	"image/draw"
	"os"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
)

func createShaders(fragSrc string, vertSrc string) (vertexShader uint32, fragmentShader uint32) {

	vertexShader, err := compileShader(vertSrc, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err = compileShader(fragSrc, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	return vertexShader, fragmentShader
}

func getTextFromFile(filePath string) (string, error) {
	rawData, err := os.ReadFile(filePath)

	if err != nil {
		fmt.Println("ERROR GETTING STRING FROM FILE " + filePath + " " + err.Error())
		return "", err
	}
	str := string(rawData)
	str += "\x00"

	return str, nil

}

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

func makeVao(points []float32) (vao uint32, vbo uint32) {

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	return vao, vbo
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
	gl.ActiveTexture(texture)
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

	//rgbaMain = rgba
	return texture
}

func setupVideo(file string) (uint32, *VideoData) {

	video, err := CreateVideoFromFile(file)

	if err != nil {
		fmt.Println("ERROR with reading frames " + err.Error())
		return 0, nil
		// return nil
	}

	var texture uint32
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(video.width),
		int32(video.height),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(video.GetData()))

	video.texture = texture
	return texture, video
	// rgbaMain = rgba
	// return video
}

func updateVideo(seconds float64, video *VideoData) {

	if video == nil {
		return
	}

	frame := int(seconds*float64(video.fps)) % int(video.frames)
	video.ReadFrame(frame)
	gl.ActiveTexture(video.texture)
	gl.BindTexture(gl.TEXTURE_2D, video.texture)
	// fmt.Printf("video text %d", video.texture)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(video.width),
		int32(video.height),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(video.GetData()))
}
