<template>
  <div id="app">
    <div class="theme-switcher" role="group" aria-label="主题模式">
      <span class="theme-switcher__label">主题</span>
      <button
        v-for="mode in themeModes"
        :key="mode.id"
        type="button"
        :class="['theme-switcher__button', { 'theme-switcher__button--active': theme === mode.id }]"
        :aria-pressed="theme === mode.id"
        @click="changeTheme(mode.id)"
      >
        {{ mode.label }}
      </button>
    </div>
    <EffectTheme v-if="theme === 'effect'" />
    <!-- 路由出口：所有匹配的路由组件会在这里渲染 -->
    <router-view />
    <FloatingBall />
  </div>
</template>

<script>
// 当使用路由后，不需要在这里直接导入页面组件
import FloatingBall from '@/components/FloatingBall'
import EffectTheme from '@/components/EffectTheme.vue';
import { applyTheme, getInitialTheme, THEME_MODES } from '@/utils/theme';
export default {
  name: 'App',
  components: {
    FloatingBall,
    EffectTheme
  },
  data() {
    return {
      theme: getInitialTheme(),
      themeModes: THEME_MODES
    };
  },
  methods: {
    changeTheme(theme) {
      this.theme = applyTheme(theme);
    }
  },
  mounted() {
    this.theme = applyTheme(this.theme);
  }
}
</script>

<style>
html,
body {
  min-height: 100%;
  margin: 0;
}

body {
  background: #f8fafc;
  transition: background-color 0.35s ease, color 0.35s ease;
}

#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  min-height: 100vh;
  margin-top: 0;
}

.theme-switcher {
  position: fixed;
  top: 14px;
  right: 18px;
  z-index: 10000;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px;
  border: 1px solid rgba(148, 163, 184, 0.45);
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.88);
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.14);
  backdrop-filter: blur(12px);
}

.theme-switcher__label {
  padding: 0 6px;
  color: #64748b;
  font-size: 12px;
  font-weight: 600;
}

.theme-switcher__button {
  border: 0;
  border-radius: 999px;
  padding: 5px 9px;
  color: #64748b;
  background: transparent;
  cursor: pointer;
  font-size: 12px;
  transition: color 0.2s ease, background-color 0.2s ease, transform 0.2s ease;
}

.theme-switcher__button:hover {
  color: #2563eb;
  transform: translateY(-1px);
}

.theme-switcher__button--active {
  color: #ffffff;
  background: #2563eb;
}

/* 深色模式：保留现有 Tailwind 页面结构，只统一基础色。 */
html[data-theme='dark'] body {
  background: #0b1120;
  color: #e5e7eb;
}

html[data-theme='dark'] #app,
html[data-theme='effect'] #app {
  color: #e5e7eb;
}

html[data-theme='dark'] .bg-white,
html[data-theme='effect'] .bg-white {
  background-color: rgba(17, 24, 39, 0.92) !important;
}

html[data-theme='dark'] .bg-gray-50,
html[data-theme='effect'] .bg-gray-50 {
  background-color: transparent !important;
}

html[data-theme='dark'] .bg-gray-100,
html[data-theme='effect'] .bg-gray-100 {
  background-color: rgba(51, 65, 85, 0.55) !important;
}

/* 覆盖 Tailwind 的 hover 背景，避免深色/特效模式下悬停导航变成刺眼的白块。 */
html[data-theme='dark'] .hover\:bg-gray-50:hover,
html[data-theme='effect'] .hover\:bg-gray-50:hover {
  background-color: rgba(51, 65, 85, 0.38) !important;
}

html[data-theme='dark'] .hover\:bg-gray-100:hover,
html[data-theme='effect'] .hover\:bg-gray-100:hover {
  background-color: rgba(71, 85, 105, 0.5) !important;
}

html[data-theme='dark'] .text-gray-800,
html[data-theme='dark'] .text-gray-900,
html[data-theme='effect'] .text-gray-800,
html[data-theme='effect'] .text-gray-900 {
  color: #f1f5f9 !important;
}

html[data-theme='dark'] .text-gray-600,
html[data-theme='dark'] .text-gray-700,
html[data-theme='dark'] .text-gray-500,
html[data-theme='effect'] .text-gray-600,
html[data-theme='effect'] .text-gray-700,
html[data-theme='effect'] .text-gray-500 {
  color: #cbd5e1 !important;
}

html[data-theme='dark'] .border-gray-100,
html[data-theme='dark'] .border-gray-200,
html[data-theme='dark'] .border-gray-300,
html[data-theme='effect'] .border-gray-100,
html[data-theme='effect'] .border-gray-200,
html[data-theme='effect'] .border-gray-300 {
  border-color: rgba(148, 163, 184, 0.3) !important;
}

html[data-theme='dark'] input,
html[data-theme='dark'] textarea,
html[data-theme='effect'] input,
html[data-theme='effect'] textarea {
  color: #e5e7eb;
  background-color: rgba(15, 23, 42, 0.75);
}

html[data-theme='effect'] #app {
  position: relative;
  z-index: 1;
  isolation: isolate;
}

html[data-theme='dark'] .theme-switcher,
html[data-theme='effect'] .theme-switcher {
  border-color: rgba(125, 211, 252, 0.35);
  background: rgba(15, 23, 42, 0.82);
}

html[data-theme='dark'] .theme-switcher__label,
html[data-theme='effect'] .theme-switcher__label {
  color: #94a3b8;
}

@media (max-width: 640px) {
  .theme-switcher {
    top: 8px;
    right: 8px;
  }

  .theme-switcher__label {
    display: none;
  }
}
</style>
