#version 410

uniform sampler2D tex0;
uniform sampler2D tex1;
uniform sampler2D tex2;
uniform sampler2D tex3; 
uniform sampler2D tex4; 
uniform float iTime;
uniform float deltaTime; 

out vec4 fragColor;

in vec2 uv; 

uniform vec2 res; 

#define OCTAVES 8


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


float edge(vec2 uv, sampler2D text) {
    vec2 p = vec2(1,1)/res.xy;
    vec2 p1 = vec2(1,0)/res.xy;
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

vec4 checkEdge(vec3 old_col, sampler2D tex) {
    vec3 new_col = old_col;
    float edge_val = edge(uv, tex) * 4.0; 
    float sensitivity = 0.1;
    if (edge_val > sensitivity) {
        new_col = vec3(0.0);
        if (edge_val < sensitivity * 1.5) {
            new_col = vec3(1.0);
        }
    }

    return vec4(vec3(new_col), 1.0);
}


float plot(vec2 st, float pct) {
    return step(pct - 0.01, st.y) - step(pct+0.01, st.y);
}


vec2 zoomUV(vec2 inuv, float zoom) {
    return ((((inuv*2) - 1.0) / zoom) + 1.0)/2.0;
}

float brightness(vec4 col) {
    return (col.r*0.2f + col.b*0.6f + col.g*0.1f) * col.a; 
}

vec4 punishment() {
    vec4 col = texture(tex1, uv);
    vec4 face = texture(tex0, vec2(uv.x, ((((uv.y*2) - 1.0) / 1.5) + 1.0)/2.0) + vec2(0, 0.1*sin(iTime+uv.x)) + vec2(0.0, 0.05*floor(random(uv+iTime)*2.0)/2.0));

    if (col.r > 0.5 && col.b < 0.3) {
        //float sinval = 0.5 + 0.5 * sin(uv.x*10.0 + iTime*2.0);
        //float val = plot(uv, sinval);//step(0.0, sin(uv.x*20.0 + uv.y + iTime));
        col = vec4(floor(face.g*3.0)/3.0);//vec4(val, 0.0, 0.0, 1.0) + face; 
    } else {
        col = vec4(1.0);
    }
    //col = max(face, col);
    return col; 
}

vec4 webcamSection() {
    float redstp = smoothstep(.4, .6, sin(iTime*10.0));
    vec4 edge_col = checkEdge(vec3(0.0), tex2);
    return vec4(edge_col.r, edge_col.g - redstp, edge_col.b - redstp, 1.0);
}

vec4 danceAlien() {
    vec4 col = texture(tex3, uv);

   // float edge_val = edge(uv, tex3);
    

    if (col.r < 0.1 && col.g < 0.1 && col.b < 0.1) {
        col = vec4(0.0);
    } else {
        float edgeval = edge(uv, tex3);
        if (edgeval > 0.2) {
            col = vec4(0.0);
        } else {
            col = vec4(1.0);
        }
       //col = checkEdge(col.rgb, tex3);
    }

    return col; 
}


vec4 Testsuo() {
    vec4 col = texture(tex4, uv);
    float bright = brightness(col);//floor(brightness(col)*6.0)/6.0;
    float edge = edge(uv, tex4);
    if (edge > 0.08) {
        return vec4(0.0);
    }
    return vec4(floor(bright*5.0)/5.0);
    //vec4(vec3(0.0, (edge > 0.08 ? 0.8 : 0.0), bright+edge*10.0), 1.0);; 
}


void main() {
    
    if (modTime < 18.0) {
        fragColor = Testsuo();
    } else if (modTime < 24.0) {
        fragColor = punishment();
    } else if (modTime < 32.0) {
        fragColor = danceAlien();
    } else if (modTime < 36.0) {
        fragColor = webcamSection();
    }
    //fragColor = Testsuo();//danceAlien();//punishment();// webcamSection(); 
}