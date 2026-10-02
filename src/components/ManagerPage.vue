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
            <div class="flex items-center gap-3">
              <button
                type="button"
                class="px-4 py-2 rounded-lg bg-indigo-600 text-white hover:bg-indigo-700 transition-colors"
                @click="openAnnouncementDialog"
              >
                发布公告
              </button>
              <router-link
                to="/home/manager/student-statistics"
                class="px-4 py-2 rounded-lg bg-blue-600 text-white hover:bg-blue-700 transition-colors"
              >
                学生统计
              </router-link>
            </div>
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
            <router-link
              to="/home/manager/all-history"
              class="text-gray-600 hover:text-blue-600 transition-colors"
              active-class="text-blue-600 font-medium"
            >
              全部提交记录
            </router-link>
          </nav>
        </header>

        <div
          v-if="announcementDialogOpen"
          class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4"
          role="presentation"
          @click.self="closeAnnouncementDialog"
        >
          <div
            class="w-full max-w-lg rounded-xl bg-white p-6 shadow-xl"
            role="dialog"
            aria-modal="true"
            aria-labelledby="publish-announcement-title"
          >
            <div class="flex items-start justify-between gap-4">
              <div>
                <h2 id="publish-announcement-title" class="text-xl font-semibold text-gray-900">发布公告</h2>
                <p class="mt-1 text-sm text-gray-500">发布后会在用户进入个人中心时显示。</p>
              </div>
              <button
                type="button"
                class="text-2xl leading-none text-gray-400 hover:text-gray-600"
                aria-label="关闭"
                :disabled="publishingAnnouncement"
                @click="closeAnnouncementDialog"
              >
                &times;
              </button>
            </div>

            <form class="mt-6 space-y-4" @submit.prevent="publishAnnouncement">
              <div>
                <label for="announcement-title-input" class="mb-1 block text-sm font-medium text-gray-700">标题</label>
                <input
                  id="announcement-title-input"
                  v-model="announcementTitle"
                  type="text"
                  maxlength="200"
                  required
                  :disabled="publishingAnnouncement"
                  class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="请输入公告标题"
                >
              </div>
              <div>
                <label for="announcement-content-input" class="mb-1 block text-sm font-medium text-gray-700">内容</label>
                <textarea
                  id="announcement-content-input"
                  v-model="announcementContent"
                  rows="7"
                  required
                  :disabled="publishingAnnouncement"
                  class="w-full resize-y rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="请输入公告内容"
                ></textarea>
              </div>
              <div class="flex justify-end gap-3 pt-2">
                <button
                  type="button"
                  class="rounded-lg border border-gray-300 px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 disabled:opacity-50"
                  :disabled="publishingAnnouncement"
                  @click="closeAnnouncementDialog"
                >
                  取消
                </button>
                <button
                  type="submit"
                  class="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 disabled:cursor-not-allowed disabled:opacity-50"
                  :disabled="publishingAnnouncement"
                >
                  {{ publishingAnnouncement ? '发布中...' : '确认发布' }}
                </button>
              </div>
            </form>
          </div>
        </div>

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
import config from '../config';
import { ElMessage } from 'element-plus';

export default {
  name: 'ManagerPage',
  components: {
    VipOnly
  },
  data() {
    return {
      announcementDialogOpen: false,
      announcementTitle: '',
      announcementContent: '',
      publishingAnnouncement: false
    };
  },
  methods: {
    openAnnouncementDialog() {
      this.announcementDialogOpen = true;
    },
    closeAnnouncementDialog() {
      if (this.publishingAnnouncement) return;
      this.announcementDialogOpen = false;
    },
    async publishAnnouncement() {
      const title = this.announcementTitle.trim();
      const content = this.announcementContent.trim();
      if (!title || !content || this.publishingAnnouncement) {
        ElMessage.warning('请输入公告标题和内容');
        return;
      }

      this.publishingAnnouncement = true;
      try {
        const response = await fetch(config.publish_announcement_url, {
          method: 'POST',
          credentials: 'include',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ title, content })
        });
        let result = {};
        try {
          result = await response.json();
        } catch (error) {
          if (!response.ok) throw new Error(`发布公告失败（HTTP ${response.status}）`);
        }
        if (!response.ok || result?.status === false) {
          throw new Error(result?.msg || `发布公告失败（HTTP ${response.status}）`);
        }

        ElMessage.success('公告发布成功');
        this.announcementTitle = '';
        this.announcementContent = '';
        this.announcementDialogOpen = false;
      } catch (error) {
        ElMessage.error(error?.message || '发布公告失败');
      } finally {
        this.publishingAnnouncement = false;
      }
    }
  }
};
</script>
