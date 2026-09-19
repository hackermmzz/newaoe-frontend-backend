<template>
  <!-- 所有特效都集中在这个文件，后续可以直接替换 WebGL 着色器。 -->
  <canvas ref="canvas" class="effect-canvas" aria-hidden="true"></canvas>
</template>

<script setup>
import { onBeforeUnmount, onMounted, ref } from 'vue';

const canvas = ref(null);
let animationFrame = 0;
let resizeHandler = null;

const vertexShaderSource = [
  'attribute vec2 a_position;',
  'void main() {',
  '  gl_Position = vec4(a_position, 0.0, 1.0);',
  '}'
].join('\n');

// WebGL 特效入口：黑洞、吸积盘、恒星、星空和流星都在这里调整。
const fragmentShaderSource = [
  'precision highp float;',
  'uniform float u_time;',
  'uniform vec2 u_resolution;',
  '',
  'float hash21(vec2 p) {',
  '  p = fract(p * vec2(123.34, 345.45));',
  '  p += dot(p, p + 34.345);',
  '  return fract(p.x * p.y);',
  '}',
  '',
  'void main() {',
  '  vec2 uv = gl_FragCoord.xy / u_resolution.xy;',
  '  vec2 p = uv * 2.0 - 1.0;',
  '  p.x *= u_resolution.x / u_resolution.y;',
  '  float time = u_time * 0.55;',
  '  float radius = length(p);',
  '',
  '  // 深空底色和微弱星云。',
  '  vec3 color = vec3(0.002, 0.004, 0.012);',
  '  color += vec3(0.015, 0.025, 0.07) * exp(-radius * 1.2);',
  '',
  '  // 背景星点：随机分布、闪烁，并避开黑洞中心。',
  '  for (int i = 0; i < 34; i++) {',
  '    float fi = float(i);',
  '    float seed = hash21(vec2(fi, fi * 2.71 + 4.0));',
  '    vec2 starPos = vec2(',
  '      fract(seed * 17.31) * 2.4 - 1.2,',
  '      fract(seed * 41.17) * 1.8 - 0.9',
  '    );',
  '    float starDistance = length(p - starPos);',
  '    float twinkle = 0.55 + 0.45 * sin(time * (1.2 + seed * 3.0) + seed * 20.0);',
  '    float star = exp(-starDistance * starDistance / (0.00045 + seed * 0.0012));',
  '    star *= smoothstep(0.13, 0.26, length(starPos));',
  '    color += vec3(0.55 + seed * 0.4, 0.7 + seed * 0.25, 1.0) * star * twinkle * 0.8;',
  '  }',
  '',
  '  // 恒星沿着吸积流靠近黑洞，靠近事件视界时逐渐消失。',
  '  float fall = fract(time * 0.055);',
  '  vec2 starPos = mix(vec2(-0.92, 0.34), vec2(-0.10, 0.035), smoothstep(0.0, 1.0, fall));',
  '  float starRadius = mix(0.035, 0.012, fall);',
  '  float starGlow = exp(-length(p - starPos) * length(p - starPos) / (starRadius * starRadius));',
  '  color += vec3(1.0, 0.58, 0.16) * starGlow * (1.0 - smoothstep(0.72, 0.98, fall));',
  '  vec2 streamDir = normalize(-starPos);',
  '  vec2 streamStart = starPos + streamDir * starRadius;',
  '  vec2 streamOffset = p - streamStart;',
  '  float streamAlong = dot(streamOffset, streamDir);',
  '  float streamSide = abs(streamOffset.x * streamDir.y - streamOffset.y * streamDir.x);',
  '  float stream = exp(-streamSide * streamSide / 0.0018) * smoothstep(0.0, 0.72, streamAlong);',
  '  stream *= (1.0 - smoothstep(0.0, 0.95, streamAlong));',
  '  color += vec3(1.0, 0.22, 0.035) * stream * 0.85;',
  '',
  '  // 旋转吸积盘：内圈高温明亮，外圈偏红，向事件视界卷入。',
  '  vec2 diskPoint = vec2(p.x, p.y * 2.35);',
  '  float diskRadius = length(diskPoint);',
  '  float diskAngle = atan(diskPoint.y, diskPoint.x);',
  '  float diskBand = exp(-pow((diskRadius - 0.30) / 0.105, 2.0));',
  '  diskBand += 0.42 * exp(-pow((diskRadius - 0.52) / 0.15, 2.0));',
  '  float turbulence = 0.55 + 0.45 * sin(diskAngle * 13.0 - time * 4.0 - diskRadius * 26.0);',
  '  vec3 hotDisk = mix(vec3(1.0, 0.11, 0.015), vec3(1.0, 0.82, 0.22), smoothstep(0.20, 0.48, diskRadius));',
  '  color += hotDisk * diskBand * turbulence * 1.45;',
  '',
  '  // 黑洞引力透镜般的蓝色边缘和完全不透光的事件视界。',
  '  float photonRing = exp(-pow((radius - 0.145) / 0.018, 2.0));',
  '  color += vec3(1.0, 0.48, 0.08) * photonRing * 1.4;',
  '  float lensGlow = exp(-pow((radius - 0.19) / 0.075, 2.0));',
  '  color += vec3(0.08, 0.25, 1.0) * lensGlow * 0.48;',
  '  color *= 1.0 - smoothstep(0.105, 0.135, radius);',
  '',
  '  // 斜向动态流星，头部发亮、尾迹逐渐消散。',
  '  for (int i = 0; i < 4; i++) {',
  '    float fi = float(i);',
  '    float seed = hash21(vec2(fi + 21.0, fi * 4.13 + 8.0));',
  '    vec2 meteorStart = vec2(-1.25 + seed * 0.55, 0.78 - seed * 0.15);',
  '    vec2 meteorDir = normalize(vec2(0.92, -0.42 - seed * 0.36));',
  '    float progress = fract(time * (0.10 + seed * 0.07) + seed);',
  '    vec2 head = meteorStart + meteorDir * progress * 2.5;',
  '    vec2 headToPoint = p - head;',
  '    float headGlow = exp(-dot(headToPoint, headToPoint) / 0.0018);',
  '    vec2 fromStart = p - meteorStart;',
  '    float along = clamp(dot(fromStart, meteorDir), 0.0, progress * 2.5);',
  '    float distanceToTrail = length(p - (meteorStart + meteorDir * along));',
  '    float tail = exp(-distanceToTrail * distanceToTrail / 0.0012) * smoothstep(0.0, 0.8, along);',
  '    color += vec3(0.4, 0.75, 1.0) * (headGlow * 1.5 + tail * 0.42);',
  '  }',
  '',
  '  gl_FragColor = vec4(color, 1.0);',
  '}'
].join('\n');

const createShader = (gl, type, source) => {
  const shader = gl.createShader(type);
  gl.shaderSource(shader, source);
  gl.compileShader(shader);
  if (!gl.getShaderParameter(shader, gl.COMPILE_STATUS)) {
    gl.deleteShader(shader);
    throw new Error(gl.getShaderInfoLog(shader) || 'WebGL shader 编译失败');
  }
  return shader;
};

const createProgram = gl => {
  const vertexShader = createShader(gl, gl.VERTEX_SHADER, vertexShaderSource);
  const fragmentShader = createShader(gl, gl.FRAGMENT_SHADER, fragmentShaderSource);
  const program = gl.createProgram();
  gl.attachShader(program, vertexShader);
  gl.attachShader(program, fragmentShader);
  gl.linkProgram(program);
  gl.deleteShader(vertexShader);
  gl.deleteShader(fragmentShader);
  if (!gl.getProgramParameter(program, gl.LINK_STATUS)) {
    gl.deleteProgram(program);
    throw new Error(gl.getProgramInfoLog(program) || 'WebGL 程序链接失败');
  }
  return program;
};

const mountWebglEffect = () => {
  const element = canvas.value;
  const gl = element?.getContext('webgl', { alpha: false, antialias: true });
  if (!gl) {
    element?.classList.add('effect-canvas--fallback');
    return;
  }

  try {
    const program = createProgram(gl);
    const positionLocation = gl.getAttribLocation(program, 'a_position');
    const timeLocation = gl.getUniformLocation(program, 'u_time');
    const resolutionLocation = gl.getUniformLocation(program, 'u_resolution');
    const buffer = gl.createBuffer();

    gl.bindBuffer(gl.ARRAY_BUFFER, buffer);
    gl.bufferData(
      gl.ARRAY_BUFFER,
      new Float32Array([-1, -1, 1, -1, -1, 1, -1, 1, 1, -1, 1, 1]),
      gl.STATIC_DRAW
    );
    gl.useProgram(program);
    gl.enableVertexAttribArray(positionLocation);
    gl.vertexAttribPointer(positionLocation, 2, gl.FLOAT, false, 0, 0);

    const resize = () => {
      const pixelRatio = Math.min(window.devicePixelRatio || 1, 2);
      const width = Math.max(1, Math.floor(window.innerWidth * pixelRatio));
      const height = Math.max(1, Math.floor(window.innerHeight * pixelRatio));
      if (element.width !== width || element.height !== height) {
        element.width = width;
        element.height = height;
        gl.viewport(0, 0, width, height);
      }
    };

    const render = timestamp => {
      resize();
      gl.uniform1f(timeLocation, timestamp * 0.001);
      gl.uniform2f(resolutionLocation, element.width, element.height);
      gl.drawArrays(gl.TRIANGLES, 0, 6);
      animationFrame = window.requestAnimationFrame(render);
    };

    resizeHandler = resize;
    window.addEventListener('resize', resizeHandler);
    animationFrame = window.requestAnimationFrame(render);
  } catch (error) {
    console.warn('WebGL 特效初始化失败，已切换到 CSS 降级效果。', error);
    element.classList.add('effect-canvas--fallback');
  }
};

onMounted(mountWebglEffect);

onBeforeUnmount(() => {
  window.cancelAnimationFrame(animationFrame);
  if (resizeHandler) window.removeEventListener('resize', resizeHandler);
});
</script>

<style>
html[data-theme='effect'] body {
  overflow-x: hidden;
  background: #071426;
}

.effect-canvas {
  position: fixed;
  inset: 0;
  /* 画布只作为背景，不能盖住路由页面的普通内容。 */
  z-index: -1;
  display: block;
  width: 100vw;
  height: 100vh;
  pointer-events: none;
}

.effect-canvas--fallback {
  background:
    radial-gradient(circle at 18% 18%, rgba(14, 165, 233, 0.28), transparent 32%),
    radial-gradient(circle at 82% 12%, rgba(99, 102, 241, 0.26), transparent 30%),
    linear-gradient(135deg, #071426, #102a46 55%, #071426);
}
</style>
