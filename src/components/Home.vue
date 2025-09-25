<template>
  <div class="flex h-screen">
    <!-- 左侧导航栏 -->
    <aside class="w-56 bg-white shadow-md flex flex-col">
      <div class="p-5 border-b border-gray-100 font-bold text-lg">
        <h1>管理系统</h1>
      </div>
      <nav class="flex-1">
        <ul class="list-none p-0 m-0">
          <li 
            v-for="(item, index) in menuItems" 
            :key="index" 
            :class="['p-4 cursor-pointer transition-colors duration-200']"
          >
            <!-- 导航链接必须使用完整路径 -->
            <router-link 
              :to="item.path"
              :class="{ 
                'bg-gray-100 font-bold': $route.path === item.path,
                'hover:bg-gray-50': $route.path !== item.path
              }"
              class="block w-full h-full"
            >
              {{ item.name }}
            </router-link>
          </li>
        </ul>
      </nav>
      <div class="p-5 border-t border-gray-100 text-sm text-gray-600">
        <p>管理员</p>
        <p>2049983474@qq.com</p>
      </div>
    </aside>

    <!-- 右侧内容区 -->
    <main class="flex-1 bg-gray-50 p-6 overflow-y-auto">
      <div class="mb-6">
        <h2 class="text-2xl font-semibold mb-2">{{ currentMenu.name }}</h2>
        <p class="text-gray-600">{{ currentMenu.description }}</p>
      </div>
      <div class="bg-white p-6 shadow-sm min-h-[calc(100vh-10rem)]">
        <!-- 子路由组件会显示在这里 -->
        <router-view></router-view>
      </div>
    </main>
  </div>
</template>

<script>
export default {
  name: 'HomePage',
  data() {
    return {
      menuItems: [
        { 
          name: '个人中心', 
          description: '查看/编辑个人信息、修改密码',
          path: '/home/student-home' // 第一个菜单的完整路径
        },
        { 
          name: '历史记录', 
          description: '查看系统操作日志、数据变更记录',
          path: '/home/history'
        },
        { 
          name: '学生考核', 
          description: '管理考核任务、评分及结果统计',
          path: '/home/assessment'
        },
        { 
          name: '系统设置', 
          description: '配置系统参数、权限及模块开关',
          path: '/home/settings'
        }
      ]
    };
  },
  computed: {
    currentMenu() {
      // 优先匹配当前路由，无匹配时默认取第一个菜单（双重保险）
      return this.menuItems.find(item => item.path === this.$route.path) || this.menuItems[0];
    }
  },
  mounted() {
    // 组件挂载后：若当前路由不是任何子路由（如仅进入父路由“/home”），自动跳转到第一个菜单
    const isCurrentRouteValid = this.menuItems.some(item => item.path === this.$route.path);
    if (!isCurrentRouteValid) {
      this.$router.push(this.menuItems[0].path);
    }
  }
};
</script>