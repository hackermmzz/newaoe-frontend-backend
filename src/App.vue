<template>
  <div id="app">
    <!-- 路由出口：所有匹配的路由组件会在这里渲染 -->
    <router-view />
    <FloatingBall />
  </div>
</template>

<script>
// 当使用路由后，不需要在这里直接导入页面组件
import FloatingBall from '@/components/FloatingBall'
export default {
  name: 'App',
  components: {
    FloatingBall
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
    