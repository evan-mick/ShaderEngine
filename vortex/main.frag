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

#define MOVE_RAD 0.2
#define MOVE_SPEED_X 2.0
#define MOVE_SPEED_Y 2.0
#define PI 3.1415
#define HALF_BANDS 10.0

float brightness(vec4 col) {
    return (col.r*0.2f + col.b*0.6f + col.g*0.1f) * col.a; 
}

void main() {

    
    vec2 offset = MOVE_RAD * vec2(cos(MOVE_SPEED_X * sin(iTime) + PI/2), sin(MOVE_SPEED_Y * cos(iTime))) - 1.0;
    vec2 mod_uv = (uv * 2.0) + offset;

    float dist = length(mod_uv);

    float dist1 = .1f;
    float dist2 = .2f;

    float cur_step = 0.0f;

    cur_step = step(dist1, dist);

    cur_step = cur_step > 0.01f ? step(dist2, dist) : 1.0 - cur_step;

    float angle = atan(mod_uv.y, mod_uv.x);

    cur_step = mod(angle + cos(iTime + uv.x*uv.x*2.0) + 0.01* sin(iTime * mod_uv.y), PI/HALF_BANDS) * 4.0 * (sin(iTime) + 1.5);//smoothstep(0.5, f0.6, angle);



    vec4 col = vec4(0.502, 0.0, 1.0, 1.0);

   // y = mx + b
  //  0 = mx + b - y
    
    // if (mod_uv.y > dist1 && abs(10.0 * mod_uv.x - mod_uv.y) < 0.01f) {
    //     cur_step = 0.0f; 
    // }

    // vec2 new_uv = vec2(mod(iTime, 2.0));

    //vec4 col = texture(tex0, new_uv);

    //col = vec4((1.0 - (floor(col.r * 5.0)/5)) * 0.25, 0, 0, 1.0); 

    // fragColor = vec4(vec3(step(0.5, dist)), 1.0);
    fragColor = vec4(vec3(cur_step), 1.0) * col;
}