#version 410

uniform sampler2D tex0;
uniform sampler2D tex1;
uniform sampler2D tex2;
uniform sampler2D tex3; 
uniform float iTime;
uniform float deltaTime; 

out vec4 fragColor;

in vec2 uv; 

uniform vec2 res; 


float brightness(vec4 col) {
    return (col.r*0.2f + col.b*0.6f + col.g*0.1f) * col.a; 
}

void main() {

    
    vec2 new_uv = vec2(mod(iTime, 2.0));
    vec4 col = texture(tex0, new_uv);

    //col = vec4((1.0 - (floor(col.r * 5.0)/5)) * 0.25, 0, 0, 1.0); 

    fragColor = col;
}