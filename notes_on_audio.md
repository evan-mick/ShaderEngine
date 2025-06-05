


Audio textures in 512x2 (col by row) textures
Top row is frequency data, the fft for that moment of the audio file, where each col is a bucket where (i+1)/512 is the percentage for up to 48khz
so for instance, the first bucket 0, is 1/512*48k frequency, then 2/512 * 48k freq for 1, etc. 

The mono information is once per frame, or once per buffer period, a window of audio samples from the current time onward



So, the process should be 
- load audio
	- create texture
- every frame
	- run fft for [0, 511] frequencies
	- also sample mono, by default, 512 actual samples, but give option in config to make time or frame buffered

