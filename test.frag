#version 410
uniform sampler2D textureSampler;
uniform float iTime;
uniform float deltaTime; 

out vec4 fragColor;

uniform vec2 res; 


float brightness(vec4 col) {
    return col.r/3.0 + col.b/3.0 + col.g/3.0; 
}

void main() {

    vec2 coord = vec2(gl_FragCoord.x/res.x / 2.0, 1.0 - (gl_FragCoord.y/res.y/2.0));
    vec2 pix = vec2(1.0/res.x, 1.0/res.y)/2.0;
    fragColor = texture(textureSampler, coord);
    float up = brightness(texture(textureSampler, coord + vec2(0.0, pix.y)));
    float down = brightness(texture(textureSampler, coord + vec2(0.0, -pix.y)));
    float left = brightness(texture(textureSampler, coord + vec2(pix.x, 0.0)));
    float right = brightness(texture(textureSampler, coord + vec2(-pix.x, 0.0)));


    float dist = sqrt((left - right) * (left - right) + (up - down) * (up - down));

    float outline = ((dist > 0.05f) ? 1.f : 0.f) * ((0.5*sin(iTime)) + 1.f);

    fragColor = vec4(outline, 0.f, 1.f, 1.f);
    // fragColor = vec4(float(gl_FragCoord.x)/res.x, float(gl_FragCoord.y)/res.y, 1.0, 1.0);//vec4(0.0, 0.0, 1.0, 1.0); //vec4(gl_FragCoord.x/res.x, gl_FragCoord.y/res.y, gl_FragCoord.z/500.f, 1.0);
    // fragColor = vec4(sin(iTime), 0.0, 0.0, 1.0);
}