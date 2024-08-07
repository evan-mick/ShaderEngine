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
    //vec2 new_uv = vec2(mod(iTime, 2.0));
    vec4 col = texture(tex1, uv);
    vec4 col2 = texture(tex0, fract((uv + vec2(0.2*sin(0.5 * (iTime + uv.y)), 0) - (iTime/50.0))*5.0));
    // vec4 col = col2; 
    //float bright = brightness(col);
    bool isWhite = (col.r > 0.8 && col.b > 0.8 && col.g > 0.8) ? true : false;

    col2 = vec4(round(brightness(col2) *(3.0 + 0.5*sin(iTime))));
    /*if ((col.r > .5) && !isWhite) {
    //if (col.b < 0.1) {
        col = col2;//vec4(0.0);
    } else {
        col = (vec4(1.0, 0.0, 0.0, 1.0));
    }*/
    //col = vec4((1.0 - (floor(col.r * 5.0)/5)) * 0.25, 0, 0, 1.0); 
    fragColor = col2;//vec4(floor(bright * 4.0)/4.0, 0.0, 0.0, 1.0) * 3.0;
}