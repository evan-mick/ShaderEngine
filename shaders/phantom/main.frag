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

// 2D Random
float random (in vec2 st) {
    return fract(sin(dot(st.xy,
                         vec2(12.9898,78.233)))
                 * 43758.5453123);
}

// 2D Noise based on Morgan McGuire @morgan3d
// https://www.shadertoy.com/view/4dS3Wd
float noise (in vec2 st) {
    vec2 i = floor(st);
    vec2 f = fract(st);

    // Four corners in 2D of a tile
    float a = random(i);
    float b = random(i + vec2(1.0, 0.0));
    float c = random(i + vec2(0.0, 1.0));
    float d = random(i + vec2(1.0, 1.0));

    vec2 u = f*f*(3.0-2.0*f);
    // u = smoothstep(0.,1.,f);

    // Mix 4 coorners percentages
    return mix(a, b, u.x) +
            (c - a)* u.y * (1.0 - u.x) +
            (d - b) * u.x * u.y;
}

float ramp(float input_val) {
    return 1.0 - abs(((fract(input_val)-0.5)*2.0)); 
}

// #define OCTAVES 6
float fbm (in vec2 st) {
    // Initial values
    float value = 0.0;
    float amplitude = .5;
    float frequency = 0.;
    //
    // Loop of octaves
    for (int i = 0; i < OCTAVES; i++) {
        value += amplitude * noise(st);
        st *= 2.;
        amplitude *= .5;
    }

    float scale = 2.0*sin(iTime*3.14)*sin(iTime*2.14) + 5.0;
    return 0.2 + floor(value * scale)/scale;
}

vec4 album() {

    vec2 mid_uv = (uv*2.0) - 1.0; 
    float angle = atan(mid_uv.y, mid_uv.x);
    float dist = length(mid_uv);

    float cos_ = cos(iTime/2.0);

    float new_angle = angle + 6.28*cos_*cos_*cos_ + cos_*dist*3.0;//+ dist*random(uv+iTime);

    vec2 mod_uv = dist* vec2(cos(new_angle), sin(new_angle));
    mod_uv =(mod_uv + 1.0) / 2.0;


    // this kinda dope too
    //vec4 col = texture(tex0, mod_uv + vec2(sin(iTime+uv.y), cos(iTime+uv.x)));

    vec4 col = texture(tex0, mod_uv + 0.1*vec2(sin(iTime+uv.y), cos(iTime+uv.x)));


    if (col.r < 0.1 && col.g < 0.1 && col.b < 0.1) {
        float noise = fbm((uv+20.0 + iTime/3.0)*10);
        col = vec4(0.0, 0.2, noise, 1.0);
    } else {
        float noise = fbm((uv+20.0 - iTime/3.0)*10);
        vec4(0.6196, 0.5451, 1.0, 1.0);
        col = vec4(0.4 + noise * 0.3, 0.5, 1.0, 1.0);

    }
    return col; 
}

float brightness(vec4 col) {
    return (col.r*0.2f + col.b*0.6f + col.g*0.1f) * col.a; 
}

vec4 bridge() {

    vec4 col = texture(tex3, uv);
    

    float bright = brightness(col);

    col = vec4(0.6, 0.5, 1.0, 1.0) * floor(bright*5.0)/5.0 * 1.5;

    if (uv.y < 0.7 + 0.1*sin(cos((-iTime + uv.x)) * 5.0)*random(floor(vec2(uv.x, uv.y)*100.0)/100.0 - iTime) && bright > 0.6) {
        float noise = fbm((uv+20.0 + iTime/3.0)*10);
        col = vec4(0.0, 0.2, noise, 1.0);
    }

    return col; 
}


vec4 img() {
    vec4 col = texture(tex1, uv);

    if (col.r < 0.1 && col.g < 0.1 && col.b < 0.1) {
        float noise = fbm((uv+20.0 + iTime/3.0)*10);
        col = vec4(0.0, 0.2, noise, 1.0);
    } else {
        col = vec4(1.0);
    }
    return col; 
}

float edge(vec2 uv, sampler2D text) {
    vec2 p = vec2(1,1)/res.xy;
    vec2 p1 = vec2(1,0)/res.xy;
    // was 0, 4 here?
    vec2 p2 = vec2(0,1)/res.xy;
    
    vec3 lum = vec3(0.2, 0.7, 0.1); 
    
    vec3 a = texture(text,uv).xyz * lum;
    vec3 b = texture(text,uv+p1).xyz * lum;
    vec3 bb = texture(text,uv-p1).xyz * lum;
    vec3 c = texture(text,uv+p2).xyz * lum;
    vec3 cc = texture(text,uv-p2).xyz * lum;
    
    float dx = length(bb)-length(b);
    float dy = length(cc)-length(c);
    
    float edge = sqrt(dx*dx + dy*dy);
    
    return edge; 

}

vec4 highschool() {
    vec4 col = texture(tex2, uv);
    float edgepart = edge(uv, tex2);

    if (col.g > 0.6 && col.r < 0.7) {
        float noise = fbm((uv+20.0 + iTime/3.0)*10);
        col = vec4(0.0, 0.2, noise, 1.0);
    } else {
        col = vec4(1.0) - (edgepart > 0.008 ? 1.0 : 0.0);
    }
    return col; 
}


void main() {
    //fragColor = album();


    float modTime = mod(iTime, 36.0);

    // fragColor = webcamSection();
    fragColor = album();
    if (modTime < 10.0) {
        fragColor = album();
    } else if (modTime < 16.0) {
        fragColor = img();
   } else if (modTime < 24.0) {
        fragColor = bridge();
    } else if (iTime > 10.0 && modTime < 28.0) {
        fragColor = highschool();
   }
//    fragColor = img();
}