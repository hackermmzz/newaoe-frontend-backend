<template>
  <!-- 所有特效都集中在这个文件，后续可以直接替换 WebGL 着色器。 -->
  <canvas ref="canvas" class="effect-canvas" aria-hidden="true"></canvas>
</template>
<script setup>
import { onBeforeUnmount, onMounted, ref } from 'vue';
import fragmentShaderSource from '!!raw-loader!@/assets/effects/thunder/effect.frag.glsl'

const vertexShaderSource = `
attribute vec2 a_position;
void main() {
  gl_Position = vec4(a_position, 0.0, 1.0);
}
`


const canvas = ref(null);
let animationFrame = 0;
let resizeHandler = null;

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
