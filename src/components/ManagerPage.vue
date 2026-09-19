<template>
  <VipOnly>
    <template #default>
      <div class="min-h-screen bg-gray-50 flex flex-col">
        <header class="bg-white shadow-sm border-b border-gray-200">
          <div class="container mx-auto px-4 py-4 flex items-center justify-between">
            <div>
              <h1 class="text-2xl font-bold text-gray-800">管理平台</h1>
              <p class="mt-1 text-sm text-gray-500">查看学生信息和提交记录</p>
            </div>
            <router-link
              to="/home/manager/student-statistics"
              class="px-4 py-2 rounded-lg bg-blue-600 text-white hover:bg-blue-700 transition-colors"
            >
              学生统计
            </router-link>
          </div>

          <nav class="container mx-auto px-4 pb-3 flex items-center gap-4 text-sm">
            <router-link
              to="/home/manager/student-statistics"
              class="text-gray-600 hover:text-blue-600 transition-colors"
              active-class="text-blue-600 font-medium"
            >
              学生统计
            </router-link>
            <router-link
              to="/home/manager/feedback"
              class="text-gray-600 hover:text-blue-600 transition-colors"
              active-class="text-blue-600 font-medium"
            >
              反馈记录
            </router-link>
          </nav>
        </header>

        <main class="flex-grow container mx-auto px-4 py-8">
          <router-view></router-view>
        </main>
      </div>
    </template>

    <template #loading>
      <div class="py-16 text-center text-gray-500">正在检查管理平台权限...</div>
    </template>

    <template #fallback>
      <div class="py-16 text-center text-gray-500">当前账号没有管理平台权限</div>
    </template>
  </VipOnly>
</template>

<script>
import VipOnly from './VipOnly.vue';

export default {
  name: 'ManagerPage',
  components: {
    VipOnly
  }
};
</script>
