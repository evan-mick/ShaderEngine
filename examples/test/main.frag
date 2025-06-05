#version 410

uniform sampler2D tex0;
uniform float iTime;
uniform float deltaTime; 

out vec4 fragColor;

in vec2 uv; 

uniform vec2 res; 

void main() {
   vec4 a = texture(tex0,uv + vec2(0.1*sin(iTime + uv.y*10.0), 0.0));
   fragColor = a;
}
