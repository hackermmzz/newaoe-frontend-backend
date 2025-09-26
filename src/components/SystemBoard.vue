<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 transition-colors duration-300">
    
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css">
    <!-- 主内容区 -->
    <main class="container mx-auto px-4 pt-24 pb-16">
      <!-- 面包屑导航 -->
      <div class="text-sm text-gray-500 dark:text-gray-400 mb-8">
        <span class="hover:text-blue-500 cursor-pointer transition-colors">首页</span>
        <i class="fa fa-chevron-right inline-block mx-2 text-xs"></i>
        <span class="text-gray-700 dark:text-gray-300">系统设置</span>
      </div>

      <!-- 设置卡片 -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- 左侧导航 -->
        <div class="lg:col-span-1">
          <div class="bg-white dark:bg-gray-800 rounded-xl shadow-sm p-6 sticky top-24">
            <h2 class="text-lg font-medium text-gray-800 dark:text-white mb-6">设置选项</h2>
            
            <ul class="space-y-1">
              <li>
                <a 
                  href="#" 
                  class="flex items-center px-4 py-3 rounded-lg bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 font-medium transition-colors"
                >
                  <i class="fa fa-info-circle w-5 h-5 mr-3"></i>
                  <span>基本信息</span>
                </a>
              </li>
              <li>
                <a 
                  href="#" 
                  class="flex items-center px-4 py-3 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-600 dark:text-gray-300 transition-colors"
                >
                  <i class="fa fa-bell w-5 h-5 mr-3"></i>
                  <span>通知设置</span>
                </a>
              </li>
              <li>
                <a 
                  href="#" 
                  class="flex items-center px-4 py-3 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-600 dark:text-gray-300 transition-colors"
                >
                  <i class="fa fa-privacy-shield w-5 h-5 mr-3"></i>
                  <span>隐私设置</span>
                </a>
              </li>
              <li>
                <a 
                  href="#" 
                  class="flex items-center px-4 py-3 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-600 dark:text-gray-300 transition-colors"
                >
                  <i class="fa fa-security w-5 h-5 mr-3"></i>
                  <span>安全设置</span>
                </a>
              </li>
              <li>
                <a 
                  href="#" 
                  class="flex items-center px-4 py-3 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-600 dark:text-gray-300 transition-colors"
                >
                  <i class="fa fa-language w-5 h-5 mr-3"></i>
                  <span>语言设置</span>
                </a>
              </li>
            </ul>
          </div>
        </div>
        
        <!-- 右侧内容 -->
        <div class="lg:col-span-2">
          <!-- 基本信息卡片 -->
          <div class="bg-white dark:bg-gray-800 rounded-xl shadow-sm p-6 mb-8 transform hover:shadow-md transition-all duration-300">
            <div class="flex items-center justify-between mb-6">
              <h2 class="text-xl font-semibold text-gray-800 dark:text-white">基本信息</h2>
              
            </div>
            
            <div class="space-y-6">
              <!-- 原神UID -->
              <div class="border-b border-gray-200 dark:border-gray-700 pb-5">
                <div class="flex flex-col md:flex-row md:items-center">
                  <div class="md:w-1/4 mb-3 md:mb-0">
                    <span class="text-gray-500 dark:text-gray-400">原神UID</span>
                  </div>
                  <div class="md:w-3/4 flex items-center">
                    <div class="flex items-center bg-gray-100 dark:bg-gray-700 px-4 py-2 rounded-lg w-full">
                      <i class="fa fa-id-card text-blue-500 mr-3"></i>
                      <span class="text-gray-800 dark:text-gray-200 font-medium">279362562</span>
                      <button 
                        @click="copyToClipboard('279362562')"
                        class="ml-auto text-gray-400 hover:text-blue-500 transition-colors"
                        aria-label="复制UID"
                      >
                        <i class="fa fa-copy"></i>
                      </button>
                    </div>
                  </div>
                </div>
              </div>
              
              <!-- 维护人员信息 -->
              <div>
                <h3 class="text-lg font-medium text-gray-800 dark:text-white mb-4">维护人员信息</h3>
                
                <!-- 微信 -->
                <div class="flex flex-col md:flex-row md:items-center mb-5">
                  <div class="md:w-1/4 mb-3 md:mb-0">
                    <span class="text-gray-500 dark:text-gray-400">微信</span>
                  </div>
                  <div class="md:w-3/4 flex items-center">
                    <div class="flex items-center bg-gray-100 dark:bg-gray-700 px-4 py-2 rounded-lg w-full">
                      <i class="fa fa-weixin text-green-500 mr-3 text-xl"></i>
                      <span class="text-gray-800 dark:text-gray-200 font-medium">19516579726</span>
                      <button 
                        @click="copyToClipboard('19516579726')"
                        class="ml-auto text-gray-400 hover:text-blue-500 transition-colors"
                        aria-label="复制微信"
                      >
                        <i class="fa fa-copy"></i>
                      </button>
                    </div>
                  </div>
                </div>
                
                <!-- QQ -->
                <div class="flex flex-col md:flex-row md:items-center">
                  <div class="md:w-1/4 mb-3 md:mb-0">
                    <span class="text-gray-500 dark:text-gray-400">QQ</span>
                  </div>
                  <div class="md:w-3/4 flex items-center">
                    <div class="flex items-center bg-gray-100 dark:bg-gray-700 px-4 py-2 rounded-lg w-full">
                      <i class="fa fa-qq text-blue-400 mr-3 text-xl"></i>
                      <span class="text-gray-800 dark:text-gray-200 font-medium">2049983474</span>
                      <button 
                        @click="copyToClipboard('2049983474')"
                        class="ml-auto text-gray-400 hover:text-blue-500 transition-colors"
                        aria-label="复制QQ"
                      >
                        <i class="fa fa-copy"></i>
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          
          <!-- 系统版本信息 -->
          <div class="bg-white dark:bg-gray-800 rounded-xl shadow-sm p-6 transform hover:shadow-md transition-all duration-300">
            <h2 class="text-xl font-semibold text-gray-800 dark:text-white mb-4">系统信息</h2>
            
            <div class="space-y-4 text-gray-600 dark:text-gray-300">
              <div class="flex justify-between pb-3 border-b border-gray-200 dark:border-gray-700">
                <span>系统版本</span>
                <span class="font-medium">v1.0.0</span>
              </div>
              <div class="flex justify-between pb-3 border-b border-gray-200 dark:border-gray-700">
                <span>最后更新</span>
                <span class="font-medium">2025-9-27</span>
              </div>
              
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- 页脚 -->
    <footer class="bg-white dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700 py-6">
      <div class="container mx-auto px-4 text-center text-gray-500 dark:text-gray-400 text-sm">
        <p>© 2025 系统管理平台. 保留所有权利.</p>
      </div>
    </footer>

    <!-- 复制成功提示 -->
    <div 
      v-if="showCopyToast" 
      class="fixed bottom-6 right-6 bg-green-500 text-white px-4 py-2 rounded-lg shadow-lg transform transition-all duration-300 animate-fade-in-out"
    >
      <i class="fa fa-check mr-2"></i> 已复制到剪贴板
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      showCopyToast: false
    }
  },
  mounted() {
    // 检查用户偏好的主题模式
    if (localStorage.theme === 'dark' || (!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  },
  methods: {
    // 切换深色/浅色模式
    toggleDarkMode() {
      if (document.documentElement.classList.contains('dark')) {
        document.documentElement.classList.remove('dark')
        localStorage.theme = 'light'
      } else {
        document.documentElement.classList.add('dark')
        localStorage.theme = 'dark'
      }
    },
    
    // 复制文本到剪贴板
    copyToClipboard(text) {
      navigator.clipboard.writeText(text).then(() => {
        this.showCopyToast = true
        setTimeout(() => {
          this.showCopyToast = false
        }, 2000)
      })
    }
  }
}
</script>

<style scoped>
/* 自定义动画 */
@keyframes fadeInOut {
  0% { opacity: 0; transform: translateY(10px); }
  10% { opacity: 1; transform: translateY(0); }
  90% { opacity: 1; transform: translateY(0); }
  100% { opacity: 0; transform: translateY(10px); }
}

.animate-fade-in-out {
  animation: fadeInOut 2s ease-in-out forwards;
}
</style>
