
Design:

Ideally, the program should run with an argument to a path to a file with data:

https://kylewbanks.com/blog/tutorial-opengl-with-golang-part-1-hello-opengl

Uses json or xml files to setup then display opengl

Has all the shader toy stuff in the background

Json has
- list of filepaths that will act as the different "channels"
- resolution data
- full shadertoy shader support (no sound stuff)

Then in the terminal, you just run the app then the path to the json, and it will automatically parse then open a window for it which you can then full screen

Later on, also support recording

Stored something like

```
Resolution: (1920, 1080)
Buffers: [
	// This would be Buffer0
	{
		Code: "(GLSL CODE HERE)",
		Channels: [
			// ORDERED BY NUMBER OF channel
			// Potential types, Buffers (could be self or others), videos, photos
			// Eventually, music???
			"CHANNEL_TYPE",
			"CHANNEL_TYPE",
		]
	}
	// This would be Buffer1
	{
		Code: "(GLSL CODE HERE)",
		Channels: [
			"CHANNEL_TYPE",
		]
	}
]
```



Backend design
-> parser
    -> creates structs from stuff
-> runner
    -> draws the shader
    -> plays music 

3 parts, sound player, drawer, parser, CLI


Parser reads files
- sets up channels
- sets up programs
- buffers, etc.



Somewhere big loop of 

Create buffers, compile program, etc.
(MAIN LOOP)
for each channel in currentShader'sChannels
    run shader
Draw final shader to screen


