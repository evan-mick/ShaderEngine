# Shader Engine


Shader Engine is a tool for locally making OpenGL quad shaders. 
It is heavily inspired by shader toy, with more functionality for recording and outputting videos.
I work on this project quite casually, it is not feature complete, but it is very servicable for making videos out of shaders.


## Running Shader Engine

To run, clone the respository, and run 
```
go build
```

It should output a binary entitled "shaderEngine"



Shader engine is a program that takes in a json file, currently with the following properties
```
{
    "width" : (int),
    "height" : (int),
    "shader" : (string),
    "folder" : (string),
    "fullscreen" : (bool),
    "recordfps" : (int),
    "textures" : [(string)] (ie list of strings representing file pathes)
}
```

**Width and Height:** The width and height of the viewport for the output video
**Shader:** The filepath of the shader .frag file relative to "folder"
**Folder:** The filepath to the folder were the shader and textures are stored
**Fullscreen:** Dictates if the viewport is fullscreen or not
**Record FPS:** This dictates how many fps the output video will have. If you only want to see the window and do not want to output a video, set this to -1.
**Textures:** a list of filepaths, each filepath is relative to "folder."  
- Supports jpg, png, jpeg, mkv, webm, mov, aiff, mp4, and mpeg file formats
- "WEBCAM" is a special keyword that allows you to use the webcam as a texture


#### Example
Given this JSON file, and assuming there was a folder in the directory where shader engine was run named "shaderEngineExample," with the files "example_image.jpeg," "video.mp4," and "main.frag" this would work. 
```
shaderEngineExample.json
{
    "width": 1920,
    "height": 1080,
    "fullscreen": true,
    "shader": "main.frag",
    "folder": "shaderEngineExample",
    "recordfps": -1,
    "textures": [
	"example_image.jpeg",
	"video.mp4",
        "WEBCAM"
    ]
}
```

Upon the user running ```./shaderEngine shaderEngineExample.json``` where the json, the folder, and the shaderEngine binary are all in the same folder, this would create a fullscreen window with a 1920x1080 resolution. It would run main.frag and input the 3 textures into said program as "tex0," "tex1," and "tex2" respectively. 

Notably, if the user switched the -1 in recordfps to any number greater than 0, they would only see a blank window, and after a few seconds (or minutes) a video with the name of their json would be outputted. This would of course yield some unprectible results given that they are using the webcam for it.

## Dependencies

The biggest external dependency that may cause headache for this project is OpenCV. It is mainly used for video reading and writing as well as its webcam support.
Part of the future plan is to add its effects as well, which is why it is being left in and I will not try to replace it with something more lightweight.
OpenGL is also a major dependency, but that should come with most OS's. 

## TODO

There are a number of potential features I would like to add, here is a list so one can get the idea of what direction I would want this project to head

### Minor
- [ ] Fullscreen switching in viewport
- [ X ] "Hot reloading" shaders, no need to restart the image loading sequence when shader gets recompiled or loaded in
- [ ] Remappable keywords
- [ ] Ex. ala shadertoy time is "iTime" right now, but the json could rename that to "time" or "t" or whatever someone wants
- [ ] Extra effects, including real time background removal (opencv)

### Major, long term
- Direct Music integration
	- Sampling music like shader toy
	- Direct beat input and beats per minute information
- Sequencing
	- Timestamping (or "beatstamping") different shaders together
	- "Events" at certain times
- (Long, long term) GUI
- (Long, long term) program sharing website
