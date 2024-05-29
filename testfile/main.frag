#version 410
// uniform sampler2D textureSampler;

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



    float bpm = 122.0;
    int beat = int(floor(iTime*(bpm/60.0)));//int((floor(iTime)/60.0)* 50.0); // bpm

    int modBeat = beat % 3;
    


    /*vec2 coord = vec2(gl_FragCoord.x/res.x, 1.0 - (gl_FragCoord.y/res.y));
    coord.x += 0.1 * sin(5.0 * sin(coord.y * 10.0 + iTime));

    vec2 pix = vec2(1.0/res.x, 1.0/res.y)/2.0;
    fragColor = texture(textureSampler, coord);
    float up = brightness(texture(textureSampler, coord + vec2(0.0, pix.y)));
    float down = brightness(texture(textureSampler, coord + vec2(0.0, -pix.y)));
    float left = brightness(texture(textureSampler, coord + vec2(pix.x, 0.0)));
    float right = brightness(texture(textureSampler, coord + vec2(-pix.x, 0.0)))*/

    // vec2 tile_uv = uv*10.0*(uv.x/uv.y);


    vec4 col = vec4(46.0/255.0, 0, 132.0/255.0, 1.0);

    vec2 tile_uv = uv*8.0;

    float aspect = res.x/res.y;
    tile_uv.x *= aspect;
    tile_uv -= 0.1 * iTime * (modBeat + 1);

    float y_mov = 10.0 * res.y/1080.0;
    float x_mov = 100.0 * res.x/1920.0;
    tile_uv += vec2(0.01*sin(iTime + uv.y*y_mov), (sin(uv.x)*2.0) + 0.05*cos(iTime + uv.x*x_mov));

    vec2 frac = fract(tile_uv);
    vec2 denom = floor(tile_uv);

    vec2 coord = uv;

    vec2 pix = vec2(1.0/res.x, 1.0/res.y)/2.0;
    fragColor = texture(tex3, coord);
    float up = brightness(texture(tex3, coord + vec2(0.0, pix.y)));
    float down = brightness(texture(tex3, coord + vec2(0.0, -pix.y)));
    float left = brightness(texture(tex3, coord + vec2(pix.x, 0.0)));
    float right = brightness(texture(tex3, coord + vec2(-pix.x, 0.0)));

    float dist = sqrt((left - right) * (left - right) + (up - down) * (up - down));

    float outline = ((dist > 0.01f) ? 1.f : 0.f);// * ((0.5*sin(iTime)) + 1.f);
    
    vec4 outlineCol = outline > 0.0f ? vec4(1.0, 1.0, 1.0, 1.0) : vec4(1.0, 0.0, 0.0, 0.0);

    /*fragColor = texture(tex0, frac);

    float bright = brightness(fragColor);

    fragColor = col * fragColor.g * 10.0;

    if (frac.x < 0.01 || frac.x > 0.99 || frac.y < 0.01 || frac.y > 0.99) {
        fragColor = col;
    }*/

    float new_x = uv.x;//uv.x*aspect - 0.5;
    
    //fragColor = outlineCol;//vec4(outline);//texture(tex3, vec2(new_x, uv.y));

    fragColor = texture(tex3, vec2(new_x, uv.y));
    

    if (new_x > 1.0 || new_x < 0.0 || uv.y > 1.0 || uv.y < 0.0) {
        fragColor = vec4(1.0, 0.0, 0.0, 0.0);
    }

    if (fragColor.r > 195.0/255.0 && (fragColor.g < 110.0/255.0 && fragColor.b < 110.0/255.0)) {
        //fragColor = texture(backText, frac);

        fragColor = (modBeat == 0) ? texture(tex0, frac) : (modBeat == 1) ? texture(tex1, frac) : texture(tex2, frac);
        float bright = brightness(fragColor);

        fragColor = col * fragColor.g * 10.0;

        if (frac.x < 0.01 || frac.x > 0.99 || frac.y < 0.01 || frac.y > 0.99) {
            fragColor = col;
        }
        //fragColor = vec4(0.0);
        return;
    }


    float bright = brightness(fragColor);

    fragColor = floor((col * bright * 7.0) * 3.0)/3.0;

    if (outline > 0.0) {
        fragColor = 1.0 - outlineCol;
    }
    // fragColor = vec4(0.0, 0.0, 0.0, 1.0);


    // fragColor = (fract(iTime) * texture(tex1, uv)) + ((1.0 - fract(iTime)) * texture(tex0, uv));




    // float dist = sqrt((left - right) * (left - right) + (up - down) * (up - down));

    // float outline = ((dist > 0.07f) ? 1.f : 0.f) * ((0.5*sin(iTime)) + 1.f);

    // fragColor = vec4(outline, 0.f, 1.f, 1.f);

    // fragColor = vec4(float(gl_FragCoord.x)/res.x, float(gl_FragCoord.y)/res.y, 1.0, 1.0);//vec4(0.0, 0.0, 1.0, 1.0); //vec4(gl_FragCoord.x/res.x, gl_FragCoord.y/res.y, gl_FragCoord.z/500.f, 1.0);
    // fragColor = vec4(sin(iTime), 0.0, 0.0, 1.0);
}