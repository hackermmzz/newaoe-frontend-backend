<template>
  <!-- 所有特效都集中在这个文件，后续可以直接替换 WebGL 着色器。 -->
  <canvas ref="canvas" class="effect-canvas" aria-hidden="true"></canvas>
</template>
<script setup>
import { onBeforeUnmount, onMounted, ref } from 'vue';
import ThunderFragmentShaderSource from '!!raw-loader!@/assets/effects/thunder/effect.frag.glsl'
import SnowFragmentShaderSource from '!!raw-loader!@/assets/effects/snow/effect.frag.glsl'
import BlackHoleFragmentShaderSource from '!!raw-loader!@/assets/effects/blackhole/effect.frag.glsl'

// Vue compiles defineProps as a script-setup macro; older ESLint configs do not know it.
// eslint-disable-next-line no-undef
const props = defineProps({
  effect: {
    type: String,
    default: 'thunder'
  }
});

//顶点着色器
const vertexShaderSource = `
attribute vec2 a_position;
void main() {
  gl_Position = vec4(a_position, 0.0, 1.0);
}
`
// 每个特效都是一份独立的片元着色器，避免把多个完整 shader 拼在一起。
const fragmentShaderSources = {
  thunder: ThunderFragmentShaderSource,
  snow: SnowFragmentShaderSource,
  blackhole: BlackHoleFragmentShaderSource
};


const canvas = ref(null);
let animationFrame = 0;
let resizeHandler = null;
let randomTexture = null;
let randomTextureContext = null;
let randomTextureInitialized = false;

// 单通道纹理中的每个 texel 都代表一个 0～1 的随机浮点数，共 256 * 256 个值。
const RANDOM_TEXTURE_WIDTH = 256;
const RANDOM_TEXTURE_HEIGHT = 256;
const RANDOM_VALUE_COUNT = RANDOM_TEXTURE_WIDTH * RANDOM_TEXTURE_HEIGHT;
const createRandomTextureData = () => {
  const values = new Uint8Array(RANDOM_VALUE_COUNT);
  for (let index = 0; index < values.length; index += 1) {
    values[index] = Math.floor(Math.random() * 256);
  }
  return values;
};

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

const createProgram = (gl, fragmentShaderSource) => {
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
    const fragmentShaderSource = fragmentShaderSources[props.effect] || fragmentShaderSources.thunder;
    const program = createProgram(gl, fragmentShaderSource);
    const positionLocation = gl.getAttribLocation(program, 'a_position');
    const timeLocation = gl.getUniformLocation(program, 'u_time');
    const resolutionLocation = gl.getUniformLocation(program, 'u_resolution');
    // 大量随机数通过纹理传输，避免超过 WebGL 的 fragment uniform 数量限制。
    const randomTextureLocation = gl.getUniformLocation(program, 'u_randomTexture');
    const randomTextureSizeLocation = gl.getUniformLocation(program, 'u_randomTextureSize');
    const buffer = gl.createBuffer();
    gl.bindBuffer(gl.ARRAY_BUFFER, buffer);
    gl.bufferData(
      gl.ARRAY_BUFFER,
      new Float32Array([-1, -1, 1, -1, -1, 1, -1, 1, 1, -1, 1, 1]),
      gl.STATIC_DRAW
    );
    gl.useProgram(program);
    if (randomTextureLocation) {
      randomTexture = gl.createTexture();
      randomTextureContext = gl;
      gl.activeTexture(gl.TEXTURE0);
      gl.bindTexture(gl.TEXTURE_2D, randomTexture);
      gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST);
      gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST);
      gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE);
      gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE);
      gl.uniform1i(randomTextureLocation, 0);
      if (randomTextureSizeLocation) {
        gl.uniform2f(randomTextureSizeLocation, RANDOM_TEXTURE_WIDTH, RANDOM_TEXTURE_HEIGHT);
      }
    }
    const updateRandomValues = () => {
      if (!randomTextureLocation || !randomTexture) return;
      gl.useProgram(program);
      gl.activeTexture(gl.TEXTURE0);
      gl.bindTexture(gl.TEXTURE_2D, randomTexture);
      const data = createRandomTextureData();
      if (randomTextureInitialized) {
        gl.texSubImage2D(
          gl.TEXTURE_2D,
          0,
          0,
          0,
          RANDOM_TEXTURE_WIDTH,
          RANDOM_TEXTURE_HEIGHT,
          gl.LUMINANCE,
          gl.UNSIGNED_BYTE,
          data
        );
      } else {
        gl.texImage2D(
          gl.TEXTURE_2D,
          0,
          gl.LUMINANCE,
          RANDOM_TEXTURE_WIDTH,
          RANDOM_TEXTURE_HEIGHT,
          0,
          gl.LUMINANCE,
          gl.UNSIGNED_BYTE,
          data
        );
        randomTextureInitialized = true;
      }
    };
    updateRandomValues();
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
  if (randomTexture && randomTextureContext) randomTextureContext.deleteTexture(randomTexture);
  randomTexture = null;
  randomTextureContext = null;
  randomTextureInitialized = false;
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
