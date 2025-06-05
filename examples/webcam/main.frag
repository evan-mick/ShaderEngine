#version 410

uniform sampler2D tex0;
uniform float iTime;
uniform float deltaTime; 

out vec4 fragColor;

in vec2 uv; 

uniform vec2 res; 

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

float brightness(vec4 col) {
    return (col.r*0.2f + col.b*0.6f + col.g*0.1f) * col.a; 
}

vec4 webcamSection() {
    vec4 new_col = background();
    if (step(0.2, edge(uv, tex3)*20.0) > 0.0) {
        return vec4(floor(((brightness(new_col) + 0.5)*16.0))/16.0, 0.0, 0.0, 1.0);
    }
    return new_col;
}

void main() {
   fragColor = webcamSection();
}
