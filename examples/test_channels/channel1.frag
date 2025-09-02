#version 410

uniform sampler2D tex0;
uniform float iTime;
uniform float deltaTime;

out vec4 fragColor;

in vec2 uv;

uniform vec2 res;

void main() {
    fragColor = vec4(1.0, 0.0, 0.0, 1.0);
}
