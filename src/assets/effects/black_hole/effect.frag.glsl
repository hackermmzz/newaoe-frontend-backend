precision highp float;

uniform float u_time;
uniform vec2 u_resolution;

#define PI 3.14159265359

float hash21(vec2 p) {
    p = fract(p * vec2(123.34, 345.45));
    p += dot(p, p + 34.345);
    return fract(p.x * p.y);
}

mat2 rot(float a) {
    float s = sin(a), c = cos(a);
    return mat2(c, -s, s, c);
}

float starField(vec2 p, float time) {
    float s = 0.0;

    for (int i = 0; i < 42; i++) {
        float fi = float(i);

        float seed = hash21(
            vec2(fi, fi * 2.73 + 5.1)
        );

        vec2 pos = vec2(
            fract(seed * 17.13) * 2.8 - 1.4,
            fract(seed * 39.71) * 2.0 - 1.0
        );

        float d = length(p - pos);

        float twinkle =
            0.55 +
            0.45 *
            sin(
                time * (1.1 + seed * 3.2)
                + seed * 18.0
            );

        float star =
            exp(
                -d * d /
                (0.00035 + seed * 0.0011)
            );

        s +=
            star *
            twinkle *
            smoothstep(
                0.18,
                0.38,
                length(pos)
            );
    }

    return s;
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
        u_time * 0.6;


    // =====================================================
    // 黑洞位置
    // =====================================================

    vec2 bh =
        vec2(
            0.18,
            0.00
        );

    vec2 bp =
        p - bh;

    float r =
        length(bp);


    // =====================================================
    // 深空背景
    // =====================================================

    vec3 color =
        vec3(
            0.002,
            0.004,
            0.012
        );

    color +=
        vec3(
            0.018,
            0.028,
            0.070
        )
        *
        exp(
            -length(p) *
            1.2
        );


    // =====================================================
    // 星云
    // =====================================================

    float nebula =
        0.5 +
        0.5 *
        sin(
            p.x * 3.0 +
            p.y * 2.0 -
            time * 0.25
        );

    color +=
        vec3(
            0.010,
            0.015,
            0.040
        )
        *
        nebula
        *
        exp(
            -length(
                p +
                vec2(
                    0.2,
                    -0.1
                )
            )
            *
            1.5
        );


    // =====================================================
    // 背景星点
    // =====================================================

    float stars =
        starField(
            p,
            time
        );

    color +=
        vec3(
            0.65,
            0.78,
            1.0
        )
        *
        stars
        *
        0.9;


    // =====================================================
    // 被黑洞吸食的恒星
    // =====================================================

    float fall =
        fract(
            time * 0.07
        );

    vec2 sunPos =
        mix(
            vec2(
                -1.10,
                0.48
            ),
            vec2(
                -0.28,
                0.06
            ),
            fall
        );

    sunPos +=
        vec2(
            sin(
                fall * 10.0
            )
            *
            0.07,

            cos(
                fall * 8.0
            )
            *
            0.05
        )
        *
        (1.0 - fall);


    float sunRadius =
        mix(
            0.095,
            0.028,
            fall
        );


    float fadeIn =
        smoothstep(
            0.00,
            0.05,
            fall
        );


    float fadeOut =
        1.0 -
        smoothstep(
            0.82,
            0.98,
            fall
        );


    float sunLife =
        fadeIn *
        fadeOut;


    vec2 sp =
        p - sunPos;


    float sunDist2 =
        dot(
            sp,
            sp
        );


    // 恒星核心
    float sunCore =
        exp(
            -sunDist2 /
            (
                sunRadius *
                sunRadius *
                0.16
            )
        );


    // 恒星光晕
    float sunGlow =
        exp(
            -sunDist2 /
            (
                sunRadius *
                sunRadius *
                2.8
            )
        );


    color +=
        vec3(
            1.00,
            0.97,
            0.86
        )
        *
        sunCore
        *
        3.0
        *
        sunLife;


    color +=
        vec3(
            1.00,
            0.68,
            0.18
        )
        *
        sunGlow
        *
        1.1
        *
        sunLife;


    // =====================================================
    // 恒星光线被黑洞拉扯
    // =====================================================

    vec2 toBH =
        normalize(
            bh -
            sunPos
        );


    float dSun =
        max(
            length(sp),
            0.0001
        );


    vec2 dirFromSun =
        sp /
        dSun;


    float forward =
        max(
            dot(
                dirFromSun,
                toBH
            ),
            0.0
        );


    float rayMask =
        pow(
            forward,
            10.0
        );


    float rayPulse1 =
        0.65 +
        0.35 *
        sin(
            dSun * 38.0 -
            time * 5.0
        );


    float rayPulse2 =
        0.65 +
        0.35 *
        sin(
            dSun * 26.0 -
            time * 4.0 +
            1.2
        );


    float solarRay =
        exp(
            -dSun * 3.8
        )
        *
        rayMask
        *
        (
            0.55 *
            rayPulse1
            +
            0.45 *
            rayPulse2
        );


    color +=
        vec3(
            1.00,
            0.88,
            0.35
        )
        *
        solarRay
        *
        1.8
        *
        sunLife;


    color +=
        vec3(
            1.00,
            0.45,
            0.10
        )
        *
        solarRay
        *
        0.6
        *
        sunLife;


    // =====================================================
    // 潮汐物质流
    // =====================================================

    vec2 streamEnd =
        bh +
        vec2(
            -0.14,
            0.00
        );


    vec2 streamVec =
        streamEnd -
        sunPos;


    float lineLen2 =
        max(
            dot(
                streamVec,
                streamVec
            ),
            0.000001
        );


    float t =
        clamp(
            dot(
                p - sunPos,
                streamVec
            )
            /
            lineLen2,

            0.0,
            1.0
        );


    vec2 curvePt =
        mix(
            sunPos,
            streamEnd,
            t
        );


    curvePt +=
        vec2(
            0.0,

            sin(
                t * PI
            )
            *
            0.18
            *
            (1.0 - fall)
        );


    curvePt +=
        vec2(
            sin(
                t * 8.0 +
                time * 1.6
            )
            *
            0.015,

            0.0
        );


    float streamWidth =
        mix(
            0.010,
            0.022,
            t
        );


    float streamDist =
        length(
            p -
            curvePt
        );


    float stream =
        exp(
            -streamDist *
            streamDist /
            (
                streamWidth *
                streamWidth
            )
        );


    stream *=
        0.72 +
        0.28 *
        sin(
            t * 42.0 -
            time * 8.0
        );


    vec3 streamColor =
        mix(
            vec3(
                1.00,
                0.82,
                0.35
            ),

            vec3(
                1.00,
                0.22,
                0.04
            ),

            t
        );


    color +=
        streamColor
        *
        stream
        *
        1.35
        *
        sunLife;


    // =====================================================
    // 巨型动态吸积盘
    // =====================================================

    vec2 dp =
        vec2(
            bp.x,
            bp.y * 2.35
        );


    float dr =
        length(dp);


    float da =
        atan(
            dp.y,
            dp.x
        );


    // 主吸积盘
    float disk =
        exp(
            -pow(
                (
                    dr -
                    0.46
                )
                /
                0.15,

                2.0
            )
        );


    // 外圈
    disk +=
        0.42
        *
        exp(
            -pow(
                (
                    dr -
                    0.76
                )
                /
                0.22,

                2.0
            )
        );


    // 动态旋转纹理
    float swirl =
        0.55 +
        0.45 *
        sin(
            da * 13.0 -
            time * 4.8 -
            dr * 26.0
        );


    // 第二层湍流
    float swirl2 =
        0.70 +
        0.30 *
        sin(
            da * 24.0 +
            time * 7.0 +
            dr * 40.0
        );


    // 多普勒增强
    float doppler =
        1.0 +
        0.65 *
        cos(
            da +
            0.65
        );


    vec3 diskColor =
        mix(
            vec3(
                1.00,
                0.10,
                0.02
            ),

            vec3(
                1.00,
                0.82,
                0.25
            ),

            smoothstep(
                0.25,
                0.80,
                dr
            )
        );


    color +=
        diskColor
        *
        disk
        *
        swirl
        *
        swirl2
        *
        doppler
        *
        1.7;


    // =====================================================
    // 吸积盘高温内圈
    // =====================================================

    float innerHeat =
        exp(
            -pow(
                (
                    dr -
                    0.32
                )
                /
                0.07,

                2.0
            )
        );


    color +=
        vec3(
            1.00,
            0.92,
            0.65
        )
        *
        innerHeat
        *
        1.05;


    // =====================================================
    // 光子环
    // =====================================================

    float photonRing =
        exp(
            -pow(
                (
                    r -
                    0.22
                )
                /
                0.020,

                2.0
            )
        );


    color +=
        vec3(
            1.00,
            0.58,
            0.12
        )
        *
        photonRing
        *
        1.8;


    // =====================================================
    // 引力透镜蓝色光晕
    // =====================================================

    float lensGlow =
        exp(
            -pow(
                (
                    r -
                    0.32
                )
                /
                0.12,

                2.0
            )
        );


    color +=
        vec3(
            0.08,
            0.25,
            1.00
        )
        *
        lensGlow
        *
        0.58;


    // =====================================================
    // 黑洞事件视界
    // =====================================================

    float eventHorizon =
        1.0 -
        smoothstep(
            0.160,
            0.205,
            r
        );


    color *=
        1.0 -
        eventHorizon;


    // =====================================================
    // 流星
    // =====================================================

    for (int i = 0; i < 4; i++) {

        float fi =
            float(i);


        float seed =
            hash21(
                vec2(
                    fi + 21.0,
                    fi * 4.17 + 8.3
                )
            );


        vec2 meteorStart =
            vec2(
                -1.30 +
                seed * 0.60,

                0.85 -
                seed * 0.20
            );


        vec2 meteorDir =
            normalize(
                vec2(
                    0.90,

                    -0.35 -
                    seed * 0.40
                )
            );


        float progress =
            fract(
                time *
                (
                    0.11 +
                    seed * 0.08
                )
                +
                seed
            );


        vec2 head =
            meteorStart
            +
            meteorDir *
            progress *
            2.7;


        vec2 headToPoint =
            p -
            head;


        float headGlow =
            exp(
                -dot(
                    headToPoint,
                    headToPoint
                )
                /
                0.0018
            );


        vec2 fromStart =
            p -
            meteorStart;


        float along =
            clamp(
                dot(
                    fromStart,
                    meteorDir
                ),

                0.0,

                progress *
                2.7
            );


        vec2 trailPoint =
            meteorStart
            +
            meteorDir *
            along;


        float trailDist =
            length(
                p -
                trailPoint
            );


        float tail =
            exp(
                -trailDist *
                trailDist /
                0.0012
            )
            *
            smoothstep(
                0.0,
                0.9,
                along
            );


        color +=
            vec3(
                0.45,
                0.78,
                1.0
            )
            *
            (
                headGlow *
                1.55
                +
                tail *
                0.45
            );
    }


    // =====================================================
    // 全局动态呼吸
    // =====================================================

    color *=
        0.97 +
        0.03 *
        sin(
            time * 1.2
        );


    // =====================================================
    // Tone Mapping
    // =====================================================

    color =
        1.0 -
        exp(
            -color * 1.1
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