precision highp float;

uniform float u_time;
uniform vec2 u_resolution;

#define PI 3.14159265359

float hash11(float p) {
    p = fract(p * 0.1031);
    p *= p + 33.33;
    p *= p + p;
    return fract(p);
}

float hash21(vec2 p) {
    p = fract(p * vec2(123.34, 345.45));
    p += dot(p, p + 34.345);
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

    return mix(
        mix(a, b, f.x),
        mix(c, d, f.x),
        f.y
    );
}

float fbm(vec2 p) {
    float v = 0.0;
    float a = 0.5;

    for (int i = 0; i < 5; i++) {
        v += noise(p) * a;
        p *= 2.02;
        a *= 0.5;
    }

    return v;
}


// =====================================================
// 单层雨滴
// =====================================================

float rainLayer(
    vec2 uv,
    float time,
    float scale,
    float speed,
    float lengthScale,
    float seedOffset
) {
    vec2 p = uv;

    // 雨向右下倾斜
    p.x += p.y * 0.23;

    p *= scale;

    // 下落
    p.y += time * speed;

    vec2 id = floor(p);
    vec2 gv = fract(p) - 0.5;

    float rnd = hash21(id + seedOffset);

    // 每列稍微错开
    gv.y += rnd * 0.8;

    // 雨滴横向偏移
    gv.x += (rnd - 0.5) * 0.55;

    // 细长雨线
    float width = 0.025;
    float len = lengthScale * (0.45 + rnd * 0.75);

    float lineX =
        exp(
            -abs(gv.x) /
            width
        );

    float lineY =
        smoothstep(
            len,
            0.0,
            abs(gv.y)
        );

    float drop =
        lineX *
        lineY;

    // 随机减少密度
    drop *=
        smoothstep(
            0.25,
            0.75,
            rnd
        );

    return drop;
}


// =====================================================
// 闪电时间控制
// =====================================================

float lightningFlash(float time) {
    // 每约 5 秒一个时间块
    float block =
        floor(time / 5.0);

    float localTime =
        mod(time, 5.0);

    float rnd =
        hash11(block * 17.13);

    // 并不是每个时间块都有闪电
    float active =
        step(0.58, rnd);

    // 随机闪电开始时间
    float start =
        0.7 +
        hash11(block * 8.71) *
        2.8;

    float t =
        localTime -
        start;

    // 主闪
    float flash1 =
        exp(
            -abs(t) * 22.0
        );

    // 第二次短闪
    float flash2 =
        exp(
            -abs(t - 0.12) * 30.0
        )
        *
        0.65;

    // 第三次更弱的回闪
    float flash3 =
        exp(
            -abs(t - 0.27) * 36.0
        )
        *
        0.28;

    return
        active *
        clamp(
            flash1 +
            flash2 +
            flash3,
            0.0,
            1.0
        );
}


// =====================================================
// 闪电路径
// =====================================================

float lightningBolt(
    vec2 p,
    float time,
    float seed
) {
    // 从天空顶部向下
    float y =
        clamp(
            (1.0 - p.y) * 0.55,
            0.0,
            1.0
        );

    float block =
        floor(time / 5.0);

    float n1 =
        noise(
            vec2(
                y * 5.0 + block * 2.1,
                seed
            )
        );

    float n2 =
        noise(
            vec2(
                y * 14.0 + block * 4.7,
                seed + 10.0
            )
        );

    float centerX =
        0.15
        +
        (n1 - 0.5) * 0.32
        +
        (n2 - 0.5) * 0.10;

    float d =
        abs(
            p.x -
            centerX
        );

    // 越往下闪电稍微变细
    float width =
        mix(
            0.022,
            0.006,
            y
        );

    float bolt =
        exp(
            -d /
            width
        );

    // 限制在天空区域
    bolt *=
        smoothstep(
            -0.85,
            0.6,
            p.y
        );

    // 主干只延伸到一定距离
    bolt *=
        1.0 -
        smoothstep(
            0.55,
            0.90,
            y
        );

    return bolt;
}


void main() {

    vec2 uv =
        gl_FragCoord.xy /
        u_resolution.xy;

    vec2 p =
        uv * 2.0 - 1.0;

    p.x *=
        u_resolution.x /
        u_resolution.y;

    float time =
        u_time;


    // =====================================================
    // 阴雨天空背景
    // =====================================================

    vec3 color =
        vec3(
            0.018,
            0.025,
            0.040
        );

    // 上暗下稍亮
    float skyGrad =
        smoothstep(
            -1.0,
            1.0,
            p.y
        );

    color +=
        vec3(
            0.035,
            0.045,
            0.060
        )
        *
        skyGrad;


    // =====================================================
    // 动态乌云
    // =====================================================

    vec2 cloudUV =
        p * vec2(1.4, 0.75);

    cloudUV.x +=
        time * 0.025;

    float cloud1 =
        fbm(
            cloudUV * 1.6
        );

    float cloud2 =
        fbm(
            cloudUV * 3.0 +
            vec2(
                time * 0.015,
                5.0
            )
        );

    float cloud =
        cloud1 * 0.72 +
        cloud2 * 0.28;

    cloud =
        smoothstep(
            0.32,
            0.78,
            cloud
        );

    color -=
        vec3(
            0.028,
            0.030,
            0.032
        )
        *
        cloud;


    // 云层顶部更浓
    float upperCloud =
        smoothstep(
            -0.15,
            0.80,
            p.y
        );

    color -=
        vec3(
            0.02,
            0.025,
            0.035
        )
        *
        upperCloud *
        cloud;


    // =====================================================
    // 远处小雨
    // =====================================================

    float rainFar =
        rainLayer(
            p,
            time,
            8.0,
            2.4,
            0.18,
            2.0
        );

    color +=
        vec3(
            0.25,
            0.32,
            0.42
        )
        *
        rainFar *
        0.25;


    // =====================================================
    // 中层雨
    // =====================================================

    float rainMid =
        rainLayer(
            p,
            time,
            12.0,
            4.6,
            0.28,
            17.0
        );

    color +=
        vec3(
            0.40,
            0.52,
            0.68
        )
        *
        rainMid *
        0.48;


    // =====================================================
    // 近处暴雨
    // =====================================================

    float rainNear =
        rainLayer(
            p,
            time,
            18.0,
            8.2,
            0.38,
            31.0
        );

    color +=
        vec3(
            0.60,
            0.72,
            0.90
        )
        *
        rainNear *
        0.72;


    // =====================================================
    // 屏幕前景快速雨线
    // =====================================================

    vec2 rp =
        p;

    rp.x +=
        rp.y * 0.20;

    rp *=
        vec2(
            26.0,
            12.0
        );

    rp.y +=
        time * 13.0;

    vec2 rid =
        floor(rp);

    vec2 rgv =
        fract(rp) -
        0.5;

    float rr =
        hash21(rid);

    rgv.x +=
        (rr - 0.5) *
        0.75;

    float fastRain =
        exp(
            -abs(rgv.x) *
            55.0
        )
        *
        smoothstep(
            0.50,
            0.0,
            abs(rgv.y)
        );

    fastRain *=
        step(
            0.72,
            rr
        );

    color +=
        vec3(
            0.72,
            0.82,
            1.0
        )
        *
        fastRain *
        0.45;


    // =====================================================
    // 雾气 / 雨幕
    // =====================================================

    float mist =
        fbm(
            vec2(
                p.x * 1.6 +
                time * 0.03,

                p.y * 0.9
            )
        );

    mist *=
        smoothstep(
            0.5,
            -0.8,
            p.y
        );

    color +=
        vec3(
            0.05,
            0.065,
            0.085
        )
        *
        mist *
        0.22;


    // =====================================================
    // 闪电
    // =====================================================

    float flash =
        lightningFlash(
            time
        );

    float bolt =
        lightningBolt(
            p,
            time,
            3.7
        );


    // 闪电主干白蓝色
    color +=
        vec3(
            0.75,
            0.88,
            1.0
        )
        *
        bolt *
        flash *
        3.5;


    // 闪电周围光晕
    color +=
        vec3(
            0.25,
            0.45,
            1.0
        )
        *
        sqrt(bolt) *
        flash *
        1.2;


    // =====================================================
    // 闪电时整个天空瞬间被照亮
    // =====================================================

    color +=
        vec3(
            0.24,
            0.30,
            0.42
        )
        *
        flash *
        (
            0.65 +
            cloud * 0.35
        );


    // =====================================================
    // 闪电时雨滴瞬间反光
    // =====================================================

    float allRain =
        rainFar * 0.2 +
        rainMid * 0.45 +
        rainNear * 0.8 +
        fastRain * 0.6;

    color +=
        vec3(
            0.65,
            0.82,
            1.0
        )
        *
        allRain *
        flash *
        0.75;


    // =====================================================
    // 闪电后的轻微余光
    // =====================================================

    color +=
        vec3(
            0.04,
            0.06,
            0.10
        )
        *
        flash *
        0.4;


    // =====================================================
    // 暗角
    // =====================================================

    vec2 vignetteUV =
        uv * (1.0 - uv.yx);

    float vignette =
        vignetteUV.x *
        vignetteUV.y *
        16.0;

    vignette =
        pow(
            vignette,
            0.18
        );

    color *=
        mix(
            0.58,
            1.0,
            vignette
        );


    // =====================================================
    // Tone Mapping
    // =====================================================

    color =
        1.0 -
        exp(
            -color *
            1.15
        );

    color =
        pow(
            color,
            vec3(0.92)
        );


    gl_FragColor =
        vec4(
            color,
            1.0
        );
}