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
            <div class="flex flex-wrap items-center justify-end gap-3">
              <button
                type="button"
                class="px-4 py-2 rounded-lg bg-indigo-600 text-white hover:bg-indigo-700 transition-colors"
                @click="openAnnouncementDialog"
              >
                发布公告
              </button>
              <button
                type="button"
                class="px-4 py-2 rounded-lg bg-purple-600 text-white hover:bg-purple-700 transition-colors"
                @click="openOjVersionDialog"
              >
                OJ版本推送
              </button>
              <button
                type="button"
                class="px-4 py-2 rounded-lg bg-emerald-600 text-white hover:bg-emerald-700 transition-colors"
                @click="openTeacherDialog"
              >
                添加考核教师
              </button>
              <button
                type="button"
                class="px-4 py-2 rounded-lg bg-orange-600 text-white hover:bg-orange-700 transition-colors disabled:cursor-not-allowed disabled:opacity-50"
                :disabled="runningAllAssessmentSubmit"
                @click="runAllAssessmentSubmit"
              >
                {{ runningAllAssessmentSubmit ? '运行中...' : '运行所有考核提交' }}
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
              普通提交记录
            </router-link>
            <router-link
              to="/home/manager/assessment-history"
              class="text-gray-600 hover:text-blue-600 transition-colors"
              active-class="text-blue-600 font-medium"
            >
              考核提交记录
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

        <div
          v-if="ojVersionDialogOpen"
          class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4"
          role="presentation"
          @click.self="closeOjVersionDialog"
        >
          <div
            class="w-full max-w-lg rounded-xl bg-white p-6 shadow-xl"
            role="dialog"
            aria-modal="true"
            aria-labelledby="oj-version-update-title"
          >
            <div class="flex items-start justify-between gap-4">
              <div>
                <h2 id="oj-version-update-title" class="text-xl font-semibold text-gray-900">OJ版本推送</h2>
                <p class="mt-1 text-sm text-gray-500">请输入要推送的最新 OJ 版本编号。</p>
              </div>
              <button
                type="button"
                class="text-2xl leading-none text-gray-400 hover:text-gray-600"
                aria-label="关闭"
                :disabled="updatingOjVersion"
                @click="closeOjVersionDialog"
              >
                &times;
              </button>
            </div>

            <form class="mt-6 space-y-4" @submit.prevent="updateOjVersion">
              <div>
                <label for="oj-version-input" class="mb-1 block text-sm font-medium text-gray-700">最新版本编号</label>
                <input
                  id="oj-version-input"
                  v-model="latestVersion"
                  type="text"
                  maxlength="100"
                  required
                  :disabled="updatingOjVersion"
                  class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="请输入最新版本编号"
                >
              </div>
              <div class="flex justify-end gap-3 pt-2">
                <button
                  type="button"
                  class="rounded-lg border border-gray-300 px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 disabled:opacity-50"
                  :disabled="updatingOjVersion"
                  @click="closeOjVersionDialog"
                >
                  取消
                </button>
                <button
                  type="submit"
                  class="rounded-lg bg-purple-600 px-4 py-2 text-sm font-medium text-white hover:bg-purple-700 disabled:cursor-not-allowed disabled:opacity-50"
                  :disabled="updatingOjVersion"
                >
                  {{ updatingOjVersion ? '推送中...' : '确认推送' }}
                </button>
              </div>
            </form>
          </div>
        </div>

        <div
          v-if="teacherDialogOpen"
          class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4"
          role="presentation"
          @click.self="closeTeacherDialog"
        >
          <div
            class="w-full max-w-lg rounded-xl bg-white p-6 shadow-xl"
            role="dialog"
            aria-modal="true"
            aria-labelledby="add-teacher-title"
          >
            <div class="flex items-start justify-between gap-4">
              <div>
                <h2 id="add-teacher-title" class="text-xl font-semibold text-gray-900">添加考核教师</h2>
                <p class="mt-1 text-sm text-gray-500">请输入教师账号或名称，添加后即可用于学生考核。</p>
              </div>
              <button
                type="button"
                class="text-2xl leading-none text-gray-400 hover:text-gray-600"
                aria-label="关闭"
                :disabled="addingTeacher"
                @click="closeTeacherDialog"
              >
                &times;
              </button>
            </div>

            <form class="mt-6 space-y-4" @submit.prevent="addTeacher">
              <div>
                <label for="teacher-input" class="mb-1 block text-sm font-medium text-gray-700">教师</label>
                <input
                  id="teacher-input"
                  v-model="teacherName"
                  type="text"
                  maxlength="100"
                  required
                  :disabled="addingTeacher"
                  class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="请输入教师账号或名称"
                >
              </div>
              <div class="flex justify-end gap-3 pt-2">
                <button
                  type="button"
                  class="rounded-lg border border-gray-300 px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 disabled:opacity-50"
                  :disabled="addingTeacher"
                  @click="closeTeacherDialog"
                >
                  取消
                </button>
                <button
                  type="submit"
                  class="rounded-lg bg-emerald-600 px-4 py-2 text-sm font-medium text-white hover:bg-emerald-700 disabled:cursor-not-allowed disabled:opacity-50"
                  :disabled="addingTeacher"
                >
                  {{ addingTeacher ? '添加中...' : '确认添加' }}
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
      publishingAnnouncement: false,
      ojVersionDialogOpen: false,
      latestVersion: '',
      updatingOjVersion: false,
      teacherDialogOpen: false,
      teacherName: '',
      addingTeacher: false,
      runningAllAssessmentSubmit: false
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
    openOjVersionDialog() {
      this.ojVersionDialogOpen = true;
    },
    closeOjVersionDialog() {
      if (this.updatingOjVersion) return;
      this.ojVersionDialogOpen = false;
    },
    openTeacherDialog() {
      this.teacherDialogOpen = true;
    },
    closeTeacherDialog() {
      if (this.addingTeacher) return;
      this.teacherDialogOpen = false;
    },
    async runAllAssessmentSubmit() {
      if (this.runningAllAssessmentSubmit) return;

      this.runningAllAssessmentSubmit = true;
      try {
        const response = await fetch(config.manager_run_all_assessment_submit_url, {
          method: 'GET',
          credentials: 'include'
        });
        let result = {};
        try {
          result = await response.json();
        } catch (error) {
          if (!response.ok) throw new Error(`运行所有考核提交失败（HTTP ${response.status}）`);
        }
        if (!response.ok || result?.status === false) {
          throw new Error(result?.msg || `运行所有考核提交失败（HTTP ${response.status}）`);
        }

        ElMessage.success(result?.msg || '已开始运行所有考核提交');
      } catch (error) {
        ElMessage.error(error?.message || '运行所有考核提交失败');
      } finally {
        this.runningAllAssessmentSubmit = false;
      }
    },
    async addTeacher() {
      const teacher = this.teacherName.trim();
      if (!teacher || this.addingTeacher) {
        ElMessage.warning('请输入教师账号或名称');
        return;
      }

      this.addingTeacher = true;
      try {
        const requestUrl = new URL(config.teacherAdd_url);
        requestUrl.searchParams.set('teacher', teacher);
        const response = await fetch(requestUrl.toString(), {
          method: 'GET',
          credentials: 'include'
        });
        let result = {};
        try {
          result = await response.json();
        } catch (error) {
          if (!response.ok) throw new Error(`添加考核教师失败（HTTP ${response.status}）`);
        }
        if (!response.ok || result?.status === false) {
          throw new Error(result?.msg || `添加考核教师失败（HTTP ${response.status}）`);
        }

        ElMessage.success('考核教师添加成功');
        this.teacherName = '';
        this.teacherDialogOpen = false;
      } catch (error) {
        ElMessage.error(error?.message || '添加考核教师失败');
      } finally {
        this.addingTeacher = false;
      }
    },
    async updateOjVersion() {
      const latestVersion = this.latestVersion.trim();
      if (!latestVersion || this.updatingOjVersion) {
        ElMessage.warning('请输入最新版本编号');
        return;
      }

      this.updatingOjVersion = true;
      try {
        const requestUrl = new URL(config.manager_oj_version_update_url);
        requestUrl.searchParams.set('latestVersion', latestVersion);
        const response = await fetch(requestUrl.toString(), {
          method: 'GET',
          credentials: 'include'
        });
        let result = {};
        try {
          result = await response.json();
        } catch (error) {
          if (!response.ok) throw new Error(`OJ版本推送失败（HTTP ${response.status}）`);
        }
        if (!response.ok || result?.status === false) {
          throw new Error(result?.msg || `OJ版本推送失败（HTTP ${response.status}）`);
        }

        ElMessage.success('OJ版本推送成功');
        this.latestVersion = '';
        this.ojVersionDialogOpen = false;
      } catch (error) {
        ElMessage.error(error?.message || 'OJ版本推送失败');
      } finally {
        this.updatingOjVersion = false;
      }
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
