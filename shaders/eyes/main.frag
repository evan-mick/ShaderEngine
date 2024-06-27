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

float brightness(vec4 col) {
    return (col.r*0.2f + col.b*0.6f + col.g*0.1f) * col.a; 
}

#define SCALE 5.0

float noise_edge(vec2 new_uv) {

    vec2 p = SCALE * vec2(1,1)/res.xy;
    vec2 p1 = SCALE * vec2(1,0)/res.xy;
    // was 0, 4 here?
    vec2 p2 = SCALE * vec2(0,1)/res.xy;
    
    vec3 lum = vec3(0.2, 0.7, 0.1); 
    
    float a = fbm(new_uv);
    float b = fbm(new_uv+p1);
    float bb = fbm(new_uv-p1);
    float c = fbm(new_uv+p2);
    float cc = fbm(new_uv-p2);
    
    float dx = bb-b;
    float dy = cc-c;
    
    float edge = sqrt(dx*dx + dy*dy);

    return edge; 

}

vec4 background() {
    vec2 new_uv = uv * SCALE + vec2(0.1 * cos(iTime + uv.y*10.0), -iTime/10.0); 

    float edge = noise_edge(new_uv);
    // Scale the coordinate system to see
    // some noise in action
    vec2 pos = vec2(new_uv);

    // Use the noise function
    float n = fbm(pos);

    return vec4(vec3(0.0, (edge > 0.1 ? 0.7 : 0.0), n+ edge*100.0), 1.0);
}


vec4 QEye() {
    vec4 col = texture(tex0, uv);
    if (col.g > 0.15 && col.r < 0.5) {
        col = background();
    } else {
        float bright = brightness(col);

        col = vec4(floor(((bright + 0.3)*4.0))/4.0, 0.0, 0.0, 1.0) - edge(uv, tex0)*20.0; 
    }
    return col; 
}

vec4 Testsuo() {
    vec4 col = texture(tex1, uv);
    float bright = floor(brightness(col)*6.0)/6.0;
    float edge = edge(uv, tex1);
    if (bright < 0.15) {
        return vec4(floor(((bright+0.4)*4.0))/4.0, 0.0, 0.0, 1.0);
    }
    return vec4(vec3(0.0, (edge > 0.08 ? 0.8 : 0.0), bright+edge*10.0), 1.0);; 
}

vec4 TF2() {
    float zoom = 2.0; 
    vec2 new_uv = ((((uv*2) - 1.0) / zoom) + 1.0)/2.0;
    vec4 col = texture(tex2, new_uv);
    float bright = floor(brightness(col)*6.0)/6.0;
    float edge = edge(new_uv, tex2);

    if (col.b > .3 && col.r > .25) {
        return vec4(floor(((bright + 0.3)*4.0))/4.0, 0.0, 0.0, 1.0);
    } else if (col.b < .2 && col.r < .2 && col.g < .2) {
        return vec4(1.0);
    } else {
        return background();// vec4(vec3(0.0, (edge > 0.08 ? 0.8 : 0.0), bright+edge*10.0), 1.0);
    }

    
}

vec4 webcamSection() {
    
    vec4 new_col = background();
    if (step(0.2, edge(uv, tex3)*20.0) > 0.0) {
        return vec4(floor(((brightness(new_col) + 0.5)*16.0))/16.0, 0.0, 0.0, 1.0);
    }
    return new_col;
}


void main() {
    
    float modTime = mod(iTime, 20.0);

    fragColor = webcamSection();
    


    if (modTime < 8.0) {
        fragColor =  webcamSection();
    } else if (modTime < 10.0) {
        fragColor = QEye();
   } else if (modTime < 13.0) {
        fragColor = Testsuo(); 
   } else if (modTime < 20.0) {
        fragColor = TF2(); 
   }

}
