#version 410
uniform sampler2D textureSampler;
out vec4 fragColor;

uniform vec2 res; 

void main() {
    fragColor = texture(textureSampler, vec2(gl_FragCoord.x/res.x / 2.0, 1.0 - (gl_FragCoord.y/res.y/2.0)));
    //fragColor = vec4(float(gl_FragCoord.x)/res.x, float(gl_FragCoord.y)/res.y, 1.0, 1.0);//vec4(0.0, 0.0, 1.0, 1.0); //vec4(gl_FragCoord.x/res.x, gl_FragCoord.y/res.y, gl_FragCoord.z/500.f, 1.0);
}