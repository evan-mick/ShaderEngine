#version 410
uniform sampler2D textureSampler;
uniform float iTime;
uniform float deltaTime; 

out vec4 fragColor;

in vec2 uv; 

uniform vec2 res; 


void main() {
    vec2 centered = (uv*2.0) - 1.0;

    float angle = atan(uv.y, uv.x);
    float len = length(center);

    fragColor = vec4(step(0.1, len));


}