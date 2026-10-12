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
            <div>
              <el-dropdown
                trigger="click"
                placement="bottom-end"
                @command="runManagerAction"
                @visible-change="actionMenuOpen = $event"
              >
                <button
                  type="button"
                  class="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 font-medium text-white shadow-sm transition-colors hover:bg-blue-700"
                  aria-haspopup="menu"
                  :aria-expanded="actionMenuOpen"
                >
                  管理操作
                  <svg
                    class="h-4 w-4 transition-transform"
                    :class="{ 'rotate-180': actionMenuOpen }"
                    viewBox="0 0 20 20"
                    fill="currentColor"
                    aria-hidden="true"
                  >
                    <path fill-rule="evenodd" d="M5.23 7.21a.75.75 0 011.06.02L10 11.168l3.71-3.938a.75.75 0 111.08 1.04l-4.25 4.5a.75.75 0 01-1.08 0l-4.25-4.5a.75.75 0 01.02-1.06z" clip-rule="evenodd" />
                  </svg>
                </button>

                <template #dropdown>
                  <el-dropdown-menu class="w-72 max-w-[calc(100vw-2rem)]">
                    <div class="px-4 pb-1 pt-2 text-xs font-semibold text-gray-400">内容与账户</div>
                    <el-dropdown-item command="openAnnouncementDialog">发布公告</el-dropdown-item>
                    <el-dropdown-item command="openOjVersionDialog">OJ版本推送</el-dropdown-item>
                    <el-dropdown-item command="openLatestPackageDialog" :disabled="uploadingLatestPackage">
                      {{ uploadingLatestPackage ? '发版中...' : '发版新版包' }}
                    </el-dropdown-item>
                    <el-dropdown-item command="openTeacherDialog">添加考核教师</el-dropdown-item>
                    <el-dropdown-item command="openEnhanceVipDialog">提升管理员</el-dropdown-item>
                    <el-dropdown-item command="openResetAssessmentDialog" :disabled="resettingAssessmentSubmit">
                      重置用户考核提交记录
                    </el-dropdown-item>
                    <div class="my-1 border-t border-gray-100"></div>
                    <div class="px-4 pb-1 pt-2 text-xs font-semibold text-gray-400">提交控制</div>
                    <el-dropdown-item command="runAllAssessmentSubmit" :disabled="runningAllAssessmentSubmit">
                      {{ runningAllAssessmentSubmit ? '运行中...' : '运行所有考核提交' }}
                    </el-dropdown-item>
                    <el-dropdown-item command="rerunAnomalRecord" :disabled="rerunningAnomalRecord">
                      {{ rerunningAnomalRecord ? '重新运行中...' : '重新运行异常记录' }}
                    </el-dropdown-item>
                    <el-dropdown-item command="blockAllSubmit" :disabled="blockingAllSubmit">
                      <span class="text-red-600">{{ blockingAllSubmit ? '阻塞中...' : '阻塞提交' }}</span>
                    </el-dropdown-item>
                    <el-dropdown-item command="cancelSubmitBlock" :disabled="cancellingSubmitBlock">
                      <span class="text-green-600">{{ cancellingSubmitBlock ? '取消中...' : '取消阻塞' }}</span>
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>

              <input
                ref="latestPackageInput"
                type="file"
                class="hidden"
                :disabled="uploadingLatestPackage"
                @change="handleLatestPackageFileChange"
              >
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
            <router-link
              to="/home/manager/assessment-ranking"
              class="text-gray-600 hover:text-blue-600 transition-colors"
              active-class="text-blue-600 font-medium"
            >
              考核排行榜
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
          v-if="latestPackageDialogOpen"
          class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4"
          role="presentation"
          @click.self="closeLatestPackageDialog"
        >
          <div
            class="w-full max-w-xl rounded-2xl bg-white p-6 shadow-2xl"
            role="dialog"
            aria-modal="true"
            aria-labelledby="latest-package-title"
          >
            <div class="flex items-start justify-between gap-4">
              <div>
                <div class="flex items-center gap-3">
                  <span class="flex h-10 w-10 items-center justify-center rounded-xl bg-cyan-100 text-cyan-700">
                    <i class="fa fa-cloud-upload text-lg" aria-hidden="true"></i>
                  </span>
                  <div>
                    <h2 id="latest-package-title" class="text-xl font-semibold text-gray-900">发版新版包</h2>
                    <p class="mt-1 text-sm text-gray-500">上传后会立即作为用户可下载的最新版包。</p>
                  </div>
                </div>
              </div>
              <button
                type="button"
                class="text-2xl leading-none text-gray-400 hover:text-gray-600 disabled:opacity-50"
                aria-label="关闭"
                :disabled="uploadingLatestPackage"
                @click="closeLatestPackageDialog"
              >
                &times;
              </button>
            </div>

            <form class="mt-6 space-y-5" @submit.prevent="uploadLatestPackage">
              <div>
                <label class="mb-2 block text-sm font-medium text-gray-700">新版包文件</label>
                <div class="flex items-center gap-3 rounded-xl border border-dashed border-cyan-300 bg-cyan-50/60 p-4">
                  <span class="flex h-11 w-11 shrink-0 items-center justify-center rounded-lg bg-white text-cyan-600 shadow-sm">
                    <i class="fa fa-file-archive-o text-xl" aria-hidden="true"></i>
                  </span>
                  <div class="min-w-0 flex-1">
                    <p class="truncate text-sm font-medium text-gray-800">
                      {{ latestPackageFile ? latestPackageFile.name : '尚未选择文件' }}
                    </p>
                    <p class="mt-1 text-xs text-gray-500">请选择要发布的 AOE 安装包或压缩包</p>
                  </div>
                  <button
                    type="button"
                    class="shrink-0 rounded-lg border border-cyan-600 px-3 py-2 text-sm font-medium text-cyan-700 transition-colors hover:bg-cyan-100 disabled:opacity-50"
                    :disabled="uploadingLatestPackage"
                    @click="selectLatestPackage"
                  >
                    {{ latestPackageFile ? '重新选择' : '选择文件' }}
                  </button>
                </div>
              </div>

              <div>
                <label for="latest-package-detail" class="mb-2 block text-sm font-medium text-gray-700">版本详情</label>
                <textarea
                  id="latest-package-detail"
                  v-model="latestPackageDetail"
                  rows="4"
                  maxlength="500"
                  required
                  :disabled="uploadingLatestPackage"
                  class="w-full resize-y rounded-xl border border-gray-300 px-3 py-2.5 text-sm focus:border-cyan-500 focus:outline-none focus:ring-2 focus:ring-cyan-500"
                  placeholder="例如：修复启动问题，更新地图资源……"
                ></textarea>
                <p class="mt-1.5 text-xs text-gray-500">这段文字会作为 detail 参数提交给发版接口。</p>
              </div>

              <div class="flex justify-end gap-3 border-t border-gray-100 pt-4">
                <button
                  type="button"
                  class="rounded-lg border border-gray-300 px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 disabled:opacity-50"
                  :disabled="uploadingLatestPackage"
                  @click="closeLatestPackageDialog"
                >
                  取消
                </button>
                <button
                  type="submit"
                  class="rounded-lg bg-cyan-600 px-4 py-2 text-sm font-medium text-white hover:bg-cyan-700 disabled:cursor-not-allowed disabled:opacity-50"
                  :disabled="uploadingLatestPackage || !latestPackageFile"
                >
                  <i v-if="uploadingLatestPackage" class="fa fa-spinner fa-spin mr-1" aria-hidden="true"></i>
                  {{ uploadingLatestPackage ? '发布中...' : '确认发版' }}
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

        <div
          v-if="enhanceVipDialogOpen"
          class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4"
          role="presentation"
          @click.self="closeEnhanceVipDialog"
        >
          <div
            class="w-full max-w-lg rounded-xl bg-white p-6 shadow-xl"
            role="dialog"
            aria-modal="true"
            aria-labelledby="enhance-vip-title"
          >
            <div class="flex items-start justify-between gap-4">
              <div>
                <h2 id="enhance-vip-title" class="text-xl font-semibold text-gray-900">提升管理员</h2>
                <p class="mt-1 text-sm text-gray-500">请输入需要提升管理员权限的用户 ID。</p>
              </div>
              <button
                type="button"
                class="text-2xl leading-none text-gray-400 hover:text-gray-600 disabled:opacity-50"
                aria-label="关闭"
                :disabled="enhancingVip"
                @click="closeEnhanceVipDialog"
              >
                &times;
              </button>
            </div>

            <form class="mt-6 space-y-4" @submit.prevent="enhanceVip">
              <div>
                <label for="enhance-vip-id-input" class="mb-1 block text-sm font-medium text-gray-700">用户 ID</label>
                <input
                  id="enhance-vip-id-input"
                  v-model="enhanceVipUserId"
                  type="text"
                  maxlength="100"
                  required
                  autocomplete="off"
                  :disabled="enhancingVip"
                  class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-amber-500 focus:outline-none focus:ring-2 focus:ring-amber-500"
                  placeholder="请输入用户 ID 或者邮箱"
                >
              </div>
              <div class="flex justify-end gap-3 pt-2">
                <button
                  type="button"
                  class="rounded-lg border border-gray-300 px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 disabled:opacity-50"
                  :disabled="enhancingVip"
                  @click="closeEnhanceVipDialog"
                >
                  取消
                </button>
                <button
                  type="submit"
                  class="rounded-lg bg-amber-600 px-4 py-2 text-sm font-medium text-white hover:bg-amber-700 disabled:cursor-not-allowed disabled:opacity-50"
                  :disabled="enhancingVip"
                >
                  {{ enhancingVip ? '提升中...' : '确认提升' }}
                </button>
              </div>
            </form>
          </div>
        </div>

        <div
          v-if="resetAssessmentDialogOpen"
          class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4"
          role="presentation"
          @click.self="closeResetAssessmentDialog"
        >
          <div
            class="w-full max-w-lg rounded-xl bg-white p-6 shadow-xl"
            role="dialog"
            aria-modal="true"
            aria-labelledby="reset-assessment-title"
          >
            <div class="flex items-start justify-between gap-4">
              <div>
                <h2 id="reset-assessment-title" class="text-xl font-semibold text-gray-900">重置用户考核提交记录</h2>
                <p class="mt-1 text-sm text-gray-500">请输入用户 ID 或邮箱。重置后该用户可以重新提交考核。</p>
              </div>
              <button
                type="button"
                class="text-2xl leading-none text-gray-400 hover:text-gray-600 disabled:opacity-50"
                aria-label="关闭"
                :disabled="resettingAssessmentSubmit"
                @click="closeResetAssessmentDialog"
              >
                &times;
              </button>
            </div>

            <form class="mt-6 space-y-4" @submit.prevent="resetAssessmentSubmit">
              <div>
                <label for="reset-assessment-id-input" class="mb-1 block text-sm font-medium text-gray-700">用户 ID 或邮箱</label>
                <input
                  id="reset-assessment-id-input"
                  v-model="resetAssessmentUserId"
                  type="text"
                  maxlength="100"
                  required
                  autocomplete="off"
                  :disabled="resettingAssessmentSubmit"
                  class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-rose-500 focus:outline-none focus:ring-2 focus:ring-rose-500"
                  placeholder="请输入用户 ID 或邮箱"
                >
              </div>
              <div class="flex justify-end gap-3 pt-2">
                <button
                  type="button"
                  class="rounded-lg border border-gray-300 px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 disabled:opacity-50"
                  :disabled="resettingAssessmentSubmit"
                  @click="closeResetAssessmentDialog"
                >
                  取消
                </button>
                <button
                  type="submit"
                  class="rounded-lg bg-rose-600 px-4 py-2 text-sm font-medium text-white hover:bg-rose-700 disabled:cursor-not-allowed disabled:opacity-50"
                  :disabled="resettingAssessmentSubmit"
                >
                  {{ resettingAssessmentSubmit ? '重置中...' : '确认重置' }}
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
import { uploadFiles } from '../utils/upload';

export default {
  name: 'ManagerPage',
  components: {
    VipOnly
  },
  data() {
    return {
      actionMenuOpen: false,
      announcementDialogOpen: false,
      announcementTitle: '',
      announcementContent: '',
      publishingAnnouncement: false,
      ojVersionDialogOpen: false,
      latestVersion: '',
      updatingOjVersion: false,
      uploadingLatestPackage: false,
      latestPackageDialogOpen: false,
      latestPackageFile: null,
      latestPackageDetail: '',
      teacherDialogOpen: false,
      teacherName: '',
      addingTeacher: false,
      enhanceVipDialogOpen: false,
      enhanceVipUserId: '',
      enhancingVip: false,
      resetAssessmentDialogOpen: false,
      resetAssessmentUserId: '',
      resettingAssessmentSubmit: false,
      runningAllAssessmentSubmit: false,
      rerunningAnomalRecord: false,
      blockingAllSubmit: false,
      cancellingSubmitBlock: false
    };
  },
  methods: {
    runManagerAction(actionName) {
      const action = this[actionName];
      this.actionMenuOpen = false;
      if (typeof action === 'function') action.call(this);
    },
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
    openLatestPackageDialog() {
      if (this.uploadingLatestPackage) return;
      this.latestPackageDialogOpen = true;
      this.latestPackageFile = null;
      this.latestPackageDetail = '';
    },
    closeLatestPackageDialog() {
      if (this.uploadingLatestPackage) return;
      this.latestPackageDialogOpen = false;
      this.latestPackageFile = null;
      this.latestPackageDetail = '';
    },
    selectLatestPackage() {
      if (this.uploadingLatestPackage) return;
      this.$refs.latestPackageInput?.click();
    },
    handleLatestPackageFileChange(event) {
      this.latestPackageFile = event.target.files?.[0] || null;
      event.target.value = '';
    },
    async uploadLatestPackage() {
      const file = this.latestPackageFile;
      const detail = this.latestPackageDetail.trim();
      if (!file || !detail || this.uploadingLatestPackage) return;

      this.uploadingLatestPackage = true;
      try {
        const requestUrl = new URL(config.manager_upload_latest_aoe_package_url);
        requestUrl.searchParams.set('detail', detail);
        const response = await fetch(requestUrl.toString(), {
          method: 'GET',
          credentials: 'include'
        });
        const result = await response.json();
        if (!response.ok || result?.status === false) {
          throw new Error(result?.msg || `获取新版包上传链接失败（HTTP ${response.status}）`);
        }

        await uploadFiles([file], [result?.data]);

        ElMessage.success('新版包发布成功');
        this.latestPackageDialogOpen = false;
        this.latestPackageFile = null;
        this.latestPackageDetail = '';
      } catch (error) {
        if (!error?.notified) {
          ElMessage.error(error?.message || '新版包发布失败');
        }
      } finally {
        this.uploadingLatestPackage = false;
      }
    },
    openTeacherDialog() {
      this.teacherDialogOpen = true;
    },
    closeTeacherDialog() {
      if (this.addingTeacher) return;
      this.teacherDialogOpen = false;
    },
    openEnhanceVipDialog() {
      this.enhanceVipDialogOpen = true;
    },
    closeEnhanceVipDialog() {
      if (this.enhancingVip) return;
      this.enhanceVipDialogOpen = false;
    },
    openResetAssessmentDialog() {
      this.resetAssessmentDialogOpen = true;
    },
    closeResetAssessmentDialog() {
      if (this.resettingAssessmentSubmit) return;
      this.resetAssessmentDialogOpen = false;
    },
    async resetAssessmentSubmit() {
      const id = this.resetAssessmentUserId.trim();
      if (!id || this.resettingAssessmentSubmit) {
        ElMessage.warning('请输入用户 ID 或邮箱');
        return;
      }

      this.resettingAssessmentSubmit = true;
      try {
        const requestUrl = new URL(config.manager_reset_assessment_submit_url);
        requestUrl.searchParams.set('id', id);
        const response = await fetch(requestUrl.toString(), {
          method: 'GET',
          credentials: 'include'
        });
        let result = {};
        try {
          result = await response.json();
        } catch (error) {
          if (!response.ok) throw new Error(`重置考核提交记录失败（HTTP ${response.status}）`);
        }
        if (!response.ok || result?.status === false) {
          throw new Error(result?.msg || `重置考核提交记录失败（HTTP ${response.status}）`);
        }

        ElMessage.success(result?.msg || '考核提交记录已重置');
        this.resetAssessmentUserId = '';
        this.resetAssessmentDialogOpen = false;
      } catch (error) {
        ElMessage.error(error?.message || '重置考核提交记录失败');
      } finally {
        this.resettingAssessmentSubmit = false;
      }
    },
    async enhanceVip() {
      const id = this.enhanceVipUserId.trim();
      if (!id || this.enhancingVip) {
        ElMessage.warning('请输入用户 ID');
        return;
      }

      this.enhancingVip = true;
      try {
        const requestUrl = new URL(config.manager_enhance_vip_url);
        requestUrl.searchParams.set('id', id);
        const response = await fetch(requestUrl.toString(), {
          method: 'GET',
          credentials: 'include'
        });
        let result = {};
        try {
          result = await response.json();
        } catch (error) {
          if (!response.ok) throw new Error(`提升管理员失败（HTTP ${response.status}）`);
        }
        if (!response.ok || result?.status === false) {
          throw new Error(result?.msg || `提升管理员失败（HTTP ${response.status}）`);
        }

        ElMessage.success(result?.msg || '管理员权限提升成功');
        this.enhanceVipUserId = '';
        this.enhanceVipDialogOpen = false;
      } catch (error) {
        ElMessage.error(error?.message || '提升管理员失败');
      } finally {
        this.enhancingVip = false;
      }
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
    async rerunAnomalRecord() {
      if (this.rerunningAnomalRecord) return;
      this.rerunningAnomalRecord = true;
      try {
        const response = await fetch(config.manager_rerun_anomal_record_url, {
          method: 'GET',
          credentials: 'include'
        });
        let result = {};
        try {
          result = await response.json();
        } catch (error) {
          if (!response.ok) throw new Error(`重新运行异常记录失败（HTTP ${response.status}）`);
        }
        if (!response.ok || result?.status === false) {
          throw new Error(result?.msg || `重新运行异常记录失败（HTTP ${response.status}）`);
        }
        ElMessage.success('异常记录已重新运行');
      } catch (error) {
        ElMessage.error(error?.message || '重新运行异常记录失败');
      } finally {
        this.rerunningAnomalRecord = false;
      }
    },
    async blockAllSubmit() {
      if (this.blockingAllSubmit) return;
      this.blockingAllSubmit = true;
      try {
        const response = await fetch(config.manager_block_all_submit_url, {
          method: 'GET',
          credentials: 'include'
        });
        let result = {};
        try {
          result = await response.json();
        } catch (error) {
          if (!response.ok) throw new Error(`阻塞提交失败（HTTP ${response.status}）`);
        }
        if (!response.ok || result?.status === false) {
          throw new Error(result?.msg || `阻塞提交失败（HTTP ${response.status}）`);
        }
        ElMessage.success('已阻塞所有提交');
      } catch (error) {
        ElMessage.error(error?.message || '阻塞提交失败');
      } finally {
        this.blockingAllSubmit = false;
      }
    },
    async cancelSubmitBlock() {
      if (this.cancellingSubmitBlock) return;
      this.cancellingSubmitBlock = true;
      try {
        const response = await fetch(config.manager_cancel_submit_block_url, {
          method: 'GET',
          credentials: 'include'
        });
        let result = {};
        try {
          result = await response.json();
        } catch (error) {
          if (!response.ok) throw new Error(`取消阻塞失败（HTTP ${response.status}）`);
        }
        if (!response.ok || result?.status === false) {
          throw new Error(result?.msg || `取消阻塞失败（HTTP ${response.status}）`);
        }
        ElMessage.success('已取消提交阻塞');
      } catch (error) {
        ElMessage.error(error?.message || '取消阻塞失败');
      } finally {
        this.cancellingSubmitBlock = false;
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
