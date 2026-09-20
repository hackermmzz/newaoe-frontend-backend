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
              v-for="(item, index) in visibleMenuItems"
              :key="index"
              :class="['p-4 cursor-pointer transition-colors duration-200']"
          >
            <!-- 导航链接必须使用完整路径 -->
            <router-link 
              :to="item.path"
              :class="{ 
                'bg-gray-100 font-bold': $route.path === item.path || $route.path.startsWith(`${item.path}/`),
                'hover:bg-gray-50': !($route.path === item.path || $route.path.startsWith(`${item.path}/`))
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
        <router-view v-if="userVip !== null"></router-view>
      </div>
    </main>
  </div>
</template>

<script>
import config from '../config';

export default {
  name: 'HomePage',
  data() {
    return {
      userVip: null,
      menuItems: [
        { 
          name: '个人中心', 
          description: '查看/编辑个人信息,反馈Bug',
          path: '/home/student-home' // 第一个菜单的完整路径
        },
        { 
          name: '历史记录', 
          description: '查看历史提交记录',
          path: '/home/history'
        },
        { 
          name: '学生考核', 
          description: '管理考核任务、评分及结果统计',
          path: '/home/assessment'
        },
        { 
          name: '排行榜',
          description: '查看对战胜负、分数及运行记录排名',
          path: '/home/ranking'
        },
        {
          name: '管理平台',
          description: '查看学生信息、提交历史及反馈记录',
          path: '/home/manager',
          requiresVip: true
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
    isVipUser() {
      return Number(this.userVip) >= Number(config.VIP_SUPER);
    },
    isTouristOrLower() {
      return this.userVip !== null
        && Number(this.userVip) <= Number(config.VIP_TOURIST);
    },
    visibleMenuItems() {
      if (this.userVip === null) return [];
      return this.menuItems.filter(item => {
        if (item.requiresVip && !this.isVipUser) return false;
        if (this.isTouristOrLower) {
          return ['/home/student-home', '/home/ranking', '/home/settings'].includes(item.path);
        }
        return true;
      });
    },
    currentMenu() {
      // 优先匹配当前路由，无匹配时默认取第一个菜单（双重保险）
      return this.visibleMenuItems.find(item => this.$route.path === item.path
        || this.$route.path.startsWith(`${item.path}/`))
        || this.visibleMenuItems[0]
        || this.menuItems[0];
    }
  },
  mounted() {
    this.loadUserVip();
  },
  watch: {
    '$route.path'() {
      this.ensureAllowedRoute();
    }
  },
  methods: {
    ensureAllowedRoute() {
      if (this.userVip === null) return;
      const isCurrentRouteValid = this.visibleMenuItems.some(item => this.$route.path === item.path
        || this.$route.path.startsWith(`${item.path}/`));
      if (!isCurrentRouteValid && this.visibleMenuItems.length) {
        this.$router.push(this.visibleMenuItems[0].path);
      }
    },
    async loadUserVip() {
      try {
        const response = await fetch(`${config.base_url}/home/studentInfo`, {
          method: 'GET',
          credentials: 'include'
        });
        const data = await response.json();
        if (response.ok && data?.status) {
          this.userVip = data?.data?.vip ?? data?.user?.vip ?? 0;
        } else {
          this.userVip = 0;
        }
        this.ensureAllowedRoute();
      } catch (error) {
        this.userVip = 0;
        this.ensureAllowedRoute();
      }
    }
  }
};
</script>
