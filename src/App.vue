<template>
  <div id="app">
    <!-- 路由出口：所有匹配的路由组件会在这里渲染 -->
    <router-view />
  </div>
</template>

<script>
// 当使用路由后，不需要在这里直接导入页面组件
export default {
  name: 'App',
  components: {
    // 路由组件不需要在这里注册，通过路由配置自动加载
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
    