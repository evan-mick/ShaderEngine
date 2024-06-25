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

#define OCTAVES 8


vec4 hsv2rgb(vec3 hsv) {

    float hue = fract(hsv.x); //only use fractional part of hue, making it loop
    float r = abs(hue * 6.0 - 3.0) - 1.0; //red
    float g = 2.0 - abs(hue * 6.0 - 2.0); //green
    float b = 2.0 - abs(hue * 6.0 - 4.0); //blue
    vec4 rgb = vec4(r,g,b, 1.0); //combine components
    rgb = clamp(rgb, 0.0, 1.0); //clamp between 0 and 1
   
    rgb = mix(vec4(1.0), rgb, 1.0 - hsv.y); //apply saturation
    rgb = rgb * hsv.z; //apply value
    return rgb;
}

float brightness(vec4 col) {
    return (col.r*0.2f + col.b*0.6f + col.g*0.1f) * col.a; 
}

float random(vec2 p) {
    return fract(sin(dot(p.xy, vec2(12.9898, 78.233))) * 43758.5453123);
}

float noise(vec2 p) {
    vec2 s = floor(p);

    float a = random(s);
    float b = random(s + vec2(1.0, 0.0));
    float c = random(s + vec2(0.0, 1.0));
    float d = random(s + vec2(1.0, 1.0));

    //step(1.0, fract(p));//smoothstep(0.0, 1.0, fract(p));
    vec2 f = smoothstep(0.0, 1.0, fract(p));

    float ab = mix(a, b, f.x);
    float cd = mix(c, d, f.x);

    float o = mix(ab, cd, f.y);
    return o; 
}

float fractal(vec2 p) {
    float o = 0.0;
    float strength = 0.5; 
    vec2 pos = p;
    for (int i = 0; i < OCTAVES; i++) {
        o += noise(p) * strength;
        pos *= 2.0;
        strength *= 0.5; 
    }
    o /= 1.0 - 0.5 * pow(0.5, float(OCTAVES - 1));
    return o; 
}

mat2 rotate(float r) {
    return mat2(
        cos(r), -sin(r),
        sin(r), cos(r)
    );
}


void main() {
    vec4 main_col = vec4(122.f/255.f, 28.f/255.f, 23.f/255.f, 1.0f);;

    vec2 new_uv = uv; 

    float n = fractal((new_uv + iTime / 200.0) * 5.0);
    float r = n * 3.1415 * 2.0;
    
    new_uv += vec2(0.0, 1.0) * rotate(r);
    
    float n_0 = fractal((new_uv + iTime / 100.0) * 8.0);
    float r_0 = n_0 * 3.1415 * 2.0 *sin(iTime/8.0);
    
    float n_1 = fractal((new_uv + vec2(0.0, 0.1) * rotate(r_0) - iTime / 200.0) * 8.0);
    float r_1 = n_1 * 3.1415 * 2.0 *sin(iTime/8.0);
    
    new_uv += vec2(0.0, 0.1) * rotate(r_1);
    
    float m1 = sin(new_uv.y * 60.0 + iTime / 1.0) / 2.0 + 0.5;
    float m2 = sin(new_uv.y * 5.0 + iTime / 2.0) / 2.0 + 0.5;
    float m3 = sin(new_uv.y * 10.0 + iTime / 5.0) / 2.0 + 0.5;

    vec3 c1 = hsv2rgb(vec3(new_uv.x, 0.9, 1.0)).rgb;
    vec3 c2 = hsv2rgb(vec3(new_uv.y, 0.9, 1.0)).rgb;
    vec3 c3 = hsv2rgb(vec3(new_uv.x/uv.y, 0.9, 1.0)).rgb;

    vec3 new_col = vec3(0.0) + (m1 * c1) + (m2 * c2) + (m3 * c3);
    new_col /= c1 + c2 + c3;


    if (brightness(vec4(new_col, 1.0)) > 0.4) {
        new_col = vec3(122.f/255.f, 28.f/255.f, 23.f/255.f);
    } else {
        new_col = vec3(1.0, 1.0, 1.0);
    }

    fragColor = vec4(vec3(new_col), 1.0);
}