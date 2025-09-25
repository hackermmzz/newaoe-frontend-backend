const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  transpileDependencies: true,
  // 开发服务器配置
  devServer: {
    port: 80,
    open: true, // 自动打开浏览器
    client:{
      overlay: {
        warnings: false,
        errors: true // 错误信息显示在页面上
      }
    }
  },
  // 生产环境优化
  productionSourceMap: false, // 生产环境不生成sourceMap
  chainWebpack: config => {
    // 移除预加载插件，优化首屏加载
    config.plugins.delete('preload')
    config.plugins.delete('prefetch')
  }
})
