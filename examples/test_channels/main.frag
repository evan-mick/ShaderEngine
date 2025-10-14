uniform sampler2D tex0;
uniform sampler2D tex1;
uniform sampler2D tex2;
uniform float iTime;
uniform float deltaTime;

out vec4 fragColor;

in vec2 uv;

uniform vec2 res;

void main() {
    vec2 flipUV = vec2(uv.x, 1.0 - uv.y);
    vec4 a = (uv.x > 0.5) ?
        ((uv.y > 0.5) ?
        texture(tex2, flipUV) : texture(tex0, flipUV)) : texture(tex1, uv);
    fragColor = a;
}
