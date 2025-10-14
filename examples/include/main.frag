
// the purpose of this is to test includes, and folder level start up

uniform sampler2D tex0;
uniform float iTime;
uniform float deltaTime;

out vec4 fragColor;

in vec2 uv;

uniform vec2 res;

void main() {
    fragColor = vec4(testInclude(), 1.0f);
}
