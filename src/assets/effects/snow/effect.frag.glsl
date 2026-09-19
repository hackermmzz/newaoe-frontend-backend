precision highp float;

uniform float u_time;
uniform vec2 u_resolution;
uniform sampler2D u_randomTexture;
uniform vec2 u_randomTextureSize;

float hash21(vec2 p) {
    p = fract(p * vec2(123.34, 456.21));
    p += dot(p, p + 45.32);
    return fract(p.x * p.y);
}

float noise(vec2 p) {
    vec2 i = floor(p);
    vec2 f = fract(p);
    f = f * f * (3.0 - 2.0 * f);
    float a = hash21(i);
    float b = hash21(i + vec2(1.0, 0.0));
    float c = hash21(i + vec2(0.0, 1.0));
    float d = hash21(i + vec2(1.0, 1.0));
    return mix(mix(a, b, f.x), mix(c, d, f.x), f.y);
}

float fbm(vec2 p) {
    float value = 0.0;
    float amplitude = 0.5;
    for (int i = 0; i < 4; i++) {
        value += noise(p) * amplitude;
        p = p * 2.0 + vec2(17.3, 9.1);
        amplitude *= 0.5;
    }
    return value;
}

float randomValue(float index) {
    // One normalized red-channel value is stored in each texel.
    float texelIndex = mod(floor(index), u_randomTextureSize.x * u_randomTextureSize.y);
    vec2 texel = vec2(
        mod(texelIndex, u_randomTextureSize.x),
        floor(texelIndex / u_randomTextureSize.x)
    );
    return texture2D(u_randomTexture, (texel + 0.5) / u_randomTextureSize).r;
}

float randomSeed() {
    // The host replaces these values every ten seconds, giving each snowfall cycle a new layout.
    return randomValue(0.0) * 3.7
        + randomValue(1.0) * 5.3
        + randomValue(2.0) * 7.1
        + randomValue(3.0) * 11.9
        + randomValue(4.0) * 13.7;
}

float snowflakeShape(vec2 point, float radius) {
    float distanceToCenter = length(point);
    float core = 1.0 - smoothstep(radius * 0.05, radius * 0.28, distanceToCenter);
    float halo = 1.0 - smoothstep(radius * 0.20, radius, distanceToCenter);

    // Six radial arms make the foreground flakes read as crystals instead of blurred circles.
    float angle = atan(point.y, point.x);
    float armDistance = abs(sin(angle * 3.0)) * distanceToCenter;
    float armWidth = 1.0 - smoothstep(radius * 0.025, radius * 0.12, armDistance);
    float armLength = 1.0 - smoothstep(radius * 0.42, radius, distanceToCenter);
    float arms = armWidth * armLength;

    // Small branch marks add an irregular, hand-cut snowflake silhouette.
    float branchDistance = abs(sin(angle * 6.0 + 0.45)) * distanceToCenter;
    float branches = (1.0 - smoothstep(radius * 0.018, radius * 0.075, branchDistance))
        * smoothstep(radius * 0.24, radius * 0.5, distanceToCenter)
        * (1.0 - smoothstep(radius * 0.48, radius * 0.82, distanceToCenter));

    return max(core, max(arms * 0.92, branches * 0.62)) + halo * 0.12;
}

// Each cell gets its own position, phase, size and wind value from the random texture.
float snowLayer(vec2 uv, float time, float scale, float speed, float size, float seed) {
    float aspect = u_resolution.x / u_resolution.y;
    vec2 p = vec2(uv.x * aspect, uv.y) * scale;
    p.y += time * speed;

    vec2 id = floor(p);
    vec2 cell = fract(p) - 0.5;
    // Each cell addresses a different part of the 65536-value random pool.
    float randomIndex = mod(abs(id.x * 307.0 + id.y * 181.0 + seed * 43.0), 65536.0);
    float randomX = randomValue(randomIndex);
    float randomY = randomValue(randomIndex + 1.0);
    float randomSize = randomValue(randomIndex + 2.0);
    float randomPhase = randomValue(randomIndex + 3.0);

    float wind = sin(time * (0.35 + randomPhase * 0.95) + randomPhase * 18.0 + id.y * 0.23)
        * (0.06 + randomPhase * 0.14);
    vec2 flakeCenter = vec2((randomX - 0.5) * 0.78 + wind, (randomY - 0.5) * 0.78);
    vec2 point = cell - flakeCenter;
    float radius = size * (0.65 + randomSize * 0.95);
    float sparkle = 0.72 + 0.28 * sin(time * (0.8 + randomPhase * 2.0) + randomPhase * 25.0);

    return snowflakeShape(point, radius) * sparkle * smoothstep(0.08, 0.72, randomSize);
}

float starField(vec2 uv, float time) {
    float aspect = u_resolution.x / u_resolution.y;
    vec2 p = vec2(uv.x * aspect, uv.y) * vec2(42.0, 25.0);
    vec2 id = floor(p);
    vec2 cell = fract(p) - 0.5;
    float random = hash21(id + 91.0);
    float star = 1.0 - smoothstep(0.012, 0.045, length(cell));
    float twinkle = 0.7 + 0.3 * sin(time * (1.0 + random * 2.0) + random * 20.0);
    return star * smoothstep(0.72, 0.96, random) * twinkle;
}

void main() {
    vec2 uv = gl_FragCoord.xy / u_resolution.xy;
    float time = u_time;

    // A blue hour gradient: darker overhead, softly lit near the horizon.
    vec3 zenith = vec3(0.012, 0.027, 0.075);
    vec3 horizon = vec3(0.105, 0.185, 0.285);
    float skyGradient = smoothstep(0.05, 0.9, uv.y);
    vec3 color = mix(horizon, zenith, skyGradient);

    // Slow, translucent aurora ribbons keep the background alive without competing with the snow.
    float auroraNoise = fbm(vec2(uv.x * (2.0 + randomValue(8.0) * 0.5) + time * 0.018, uv.y * 3.0));
    float auroraWave = sin(uv.x * (4.5 + randomValue(9.0) * 1.5) + auroraNoise * 3.0 + time * 0.12 + randomValue(10.0));
    float auroraMask = smoothstep(0.48, 0.82, uv.y) * smoothstep(0.25, 0.75, auroraWave);
    color += vec3(0.025, 0.09, 0.12) * auroraMask * 0.55;

    // Moon and its broad halo.
    vec2 moonPosition = vec2(0.78, 0.73);
    float aspect = u_resolution.x / u_resolution.y;
    // UV 的横纵单位对应不同的像素数量，先按宽高比校正横坐标才能保持正圆。
    float moonDistance = length((uv - moonPosition) * vec2(aspect, 1.0));
    float moonGlow = exp(-moonDistance * moonDistance / 0.035);
    float moon = 1.0 - smoothstep(0.066, 0.078, moonDistance);
    color += vec3(0.16, 0.22, 0.38) * moonGlow;
    color += vec3(0.86, 0.93, 1.0) * moon;

    // A sparse star field is mostly visible above the mountain line.
    float stars = starField(uv, time) * smoothstep(0.32, 0.85, uv.y);
    color += vec3(0.62, 0.76, 1.0) * stars * 0.75;

    // Layered mountain silhouettes with snow caps.
    float ridgeFar = 0.37 + fbm(vec2(uv.x * (2.5 + randomValue(11.0) * 0.6) + 4.0, 0.0)) * 0.15;
    float farMask = 1.0 - smoothstep(ridgeFar, ridgeFar + 0.012, uv.y);
    color = mix(color, vec3(0.045, 0.095, 0.15), farMask * 0.9);
    float farCap = farMask * smoothstep(ridgeFar - 0.045, ridgeFar - 0.008, uv.y);
    color = mix(color, vec3(0.30, 0.42, 0.54), farCap * 0.72);

    float ridgeNear = 0.27 + fbm(vec2(uv.x * (3.5 + randomValue(12.0) * 0.8) - 7.0, 0.0)) * 0.18;
    float nearMask = 1.0 - smoothstep(ridgeNear, ridgeNear + 0.014, uv.y);
    color = mix(color, vec3(0.025, 0.062, 0.105), nearMask);
    float nearCap = nearMask * smoothstep(ridgeNear - 0.055, ridgeNear - 0.01, uv.y);
    color = mix(color, vec3(0.18, 0.28, 0.38), nearCap * 0.8);

    // Blue snowfield at the bottom of the frame.
    float ground = 1.0 - smoothstep(0.27, 0.48, uv.y);
    float groundTexture = fbm(vec2(uv.x * 5.0, uv.y * 4.0 + time * 0.01));
    color += vec3(0.055, 0.095, 0.14) * ground;
    color += vec3(0.025, 0.04, 0.06) * groundTexture * ground;

    // Falling snow, from tiny distant flakes to bright foreground flakes.
    float snowFar = snowLayer(uv, time, 15.0, 0.30 + randomValue(13.0) * 0.08, 0.030, 2.0);
    float snowMid = snowLayer(uv, time, 22.0, 0.52 + randomValue(14.0) * 0.12, 0.050, 17.0);
    float snowNear = snowLayer(uv, time, 32.0, 0.78 + randomValue(15.0) * 0.16, 0.080, 43.0);
    float snowFront = snowLayer(uv, time, 45.0, 1.05, 0.110, 71.0);
    color += vec3(0.58, 0.72, 0.90) * snowFar * 0.42;
    color += vec3(0.76, 0.86, 0.98) * snowMid * 0.62;
    color += vec3(0.92, 0.97, 1.0) * snowNear * 0.88;
    color += vec3(1.0) * snowFront * 1.05;

    // Subtle vignette keeps the page content readable in the middle.
    vec2 vignettePoint = uv - 0.5;
    float vignette = 1.0 - smoothstep(0.42, 0.78, length(vignettePoint * vec2(0.86, 0.95)));
    color *= mix(0.68, 1.0, vignette);

    color = 1.0 - exp(-color * 1.18);
    color = pow(color, vec3(0.94));
    gl_FragColor = vec4(color, 1.0);
}
