<template>
  <section class="feedback-records space-y-6">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h2 class="text-2xl font-semibold text-gray-800">反馈记录</h2>
        <p class="mt-1 text-sm text-gray-500">按提交时间倒序显示所有反馈内容。</p>
      </div>
      <button
        type="button"
        class="px-4 py-2 rounded-lg bg-blue-600 text-white hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        :disabled="isLoading"
        @click="loadFeedback(currentPage)"
      >
        {{ isLoading ? '加载中...' : '刷新列表' }}
      </button>
    </div>

    <div v-if="errorMessage" class="feedback-records__error rounded-lg bg-red-50 p-4 text-sm text-red-700" role="alert">
      {{ errorMessage }}
    </div>

    <!-- 分页控件放在记录上方，避免翻页时需要滚动到页面底部。 -->
    <div
      v-if="!isLoading"
      class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-y border-gray-200 py-4"
    >
      <p class="feedback-records__muted text-sm text-gray-500">
        第 {{ currentPage }} 页，本页 {{ feedbackList.length }} 条
      </p>

      <div class="flex flex-wrap items-center gap-2">
        <button
          type="button"
          class="px-3 py-1.5 border border-gray-300 rounded-md text-sm hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          :disabled="isLoading || currentPage <= 1"
          @click="handlePrevPage"
        >
          上一页
        </button>

        <div class="flex items-center gap-1">
          <input
            v-model.number="targetPage"
            type="number"
            :min="1"
            class="feedback-records__page-input w-16 px-2 py-1.5 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-1 focus:ring-blue-500"
            aria-label="跳转页码"
          >
          <button
            type="button"
            class="px-2 py-1.5 border border-gray-300 rounded-md text-sm hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
            :disabled="isLoading"
            @click="handlePageJump"
          >
            跳转
          </button>
        </div>

        <button
          type="button"
          class="px-3 py-1.5 border border-gray-300 rounded-md text-sm hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          :disabled="isLoading || !hasNextPage"
          @click="handleNextPage"
        >
          下一页
        </button>
      </div>
    </div>

    <div v-if="isLoading" class="feedback-records__muted py-16 text-center text-gray-500">
      正在加载第 {{ currentPage }} 页反馈记录...
    </div>

    <div v-else-if="!feedbackList.length" class="feedback-records__muted py-16 text-center text-gray-500">
      暂无反馈记录
    </div>

    <div v-else class="space-y-5">
      <article
        v-for="(item, index) in feedbackList"
        :key="item.indices ?? index"
        class="feedback-card overflow-hidden rounded-xl border-2 border-amber-200 bg-white shadow-sm transition-shadow hover:shadow-md"
      >
        <header class="feedback-card__header flex flex-col gap-2 border-b border-amber-100 bg-amber-50 px-5 py-4 sm:flex-row sm:items-center sm:justify-between">
          <h3 class="feedback-card__title font-semibold text-gray-900">反馈 #{{ item.indices ?? '未知' }}</h3>
          <span class="feedback-card__time text-sm text-gray-600">提交时间：{{ formatDate(item.submittime) }}</span>
        </header>

        <div class="space-y-4 px-5 py-4">
          <dl class="feedback-card__details grid grid-cols-1 gap-3 text-sm text-gray-700 sm:grid-cols-2">
            <div>
              <dt class="feedback-card__label font-medium text-gray-500">基础目录</dt>
              <dd class="feedback-card__value mt-1 break-all">{{ item.basefolder || '未提供' }}</dd>
            </div>
            <div>
              <dt class="feedback-card__label font-medium text-gray-500">记录编号</dt>
              <dd class="feedback-card__value mt-1">{{ item.indices ?? '未知' }}</dd>
            </div>
          </dl>

          <div class="feedback-card__links rounded-lg border-2 border-blue-200 bg-blue-50/40 p-3">
            <div class="feedback-card__links-title mb-2 flex items-center gap-2 text-sm font-semibold text-blue-800">
              <span class="inline-flex h-2.5 w-2.5 rounded-full bg-blue-500" aria-hidden="true"></span>
              反馈链接
            </div>
            <a
              v-if="getLink(item)"
              :href="getLink(item)"
              target="_blank"
              rel="noopener noreferrer"
              class="feedback-card__link block break-all rounded-md border border-blue-300 bg-white px-3 py-2 text-sm font-medium text-blue-700 underline decoration-blue-300 underline-offset-2 hover:border-blue-500 hover:bg-blue-50 hover:text-blue-900"
            >
              {{ getLink(item) }}
            </a>
            <p v-else class="feedback-card__empty py-3 text-center text-sm text-gray-500">该记录没有反馈链接</p>
          </div>
        </div>
      </article>
    </div>
  </section>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { ElMessage } from 'element-plus';
import config from '../config';

const pageSize = config.ManagerFeedbackRecordPerPage || config.HistoryRecordPerPage || 10;
const currentPage = ref(1);
const targetPage = ref(1);
const hasNextPage = ref(true);
const isLoading = ref(false);
const errorMessage = ref('');
const feedbackList = ref([]);

const formatDate = value => {
  if (value === null || value === undefined || value === '') return '未知时间';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return String(value);
  return date.toLocaleString('zh-CN', { hour12: false, timeZone: 'Asia/Shanghai' });
};

// 后端结构体的 JSON 标签目前是 "string"，同时兼容常见的 Link/link 返回形式。
const getLink = item => String(
  item?.string
  ?? item?.Link
  ?? item?.link
  ?? item?.String
  ?? ''
);

const extractRecords = payload => {
  if (Array.isArray(payload)) return payload;
  if (Array.isArray(payload?.record)) return payload.record;
  if (Array.isArray(payload?.records)) return payload.records;
  return [];
};

const sortBySubmitTime = records => records
  .map((record, index) => ({ record, index }))
  .sort((left, right) => {
    const leftTime = new Date(left.record?.submittime).getTime();
    const rightTime = new Date(right.record?.submittime).getTime();
    const leftValid = Number.isFinite(leftTime);
    const rightValid = Number.isFinite(rightTime);
    if (leftValid && rightValid && leftTime !== rightTime) return rightTime - leftTime;
    if (leftValid !== rightValid) return leftValid ? -1 : 1;
    return left.index - right.index;
  })
  .map(({ record }) => record);

const loadFeedback = async (page = currentPage.value) => {
  if (isLoading.value || !Number.isInteger(page) || page < 1) return;
  isLoading.value = true;
  errorMessage.value = '';

  try {
    const beg = (page - 1) * pageSize;
    const end = beg + pageSize - 1;
    const requestUrl = new URL(config.manager_feedback_url);
    requestUrl.searchParams.set('range', `${beg}:${end}`);

    const response = await fetch(requestUrl.toString(), {
      method: 'GET',
      credentials: 'include'
    });
    const result = await response.json();
    // 兼容管理接口常见的 { status, data } 包装，以及接口直接返回数组的形式。
    if (!response.ok || (!Array.isArray(result) && result?.status === false)) {
      throw new Error(result?.msg || `加载反馈记录失败（HTTP ${response.status}）`);
    }

    const rawRecords = extractRecords(Array.isArray(result) ? result : (result?.data ?? result));
    if (!rawRecords.length && page > 1) {
      if (page === currentPage.value + 1) hasNextPage.value = false;
      targetPage.value = currentPage.value;
      return;
    }

    feedbackList.value = sortBySubmitTime(rawRecords);
    currentPage.value = page;
    targetPage.value = page;
    hasNextPage.value = rawRecords.length >= pageSize;
  } catch (error) {
    errorMessage.value = error?.message || '加载反馈记录失败';
    targetPage.value = currentPage.value;
    ElMessage.error(errorMessage.value);
  } finally {
    isLoading.value = false;
  }
};

const handlePrevPage = () => {
  if (currentPage.value > 1) loadFeedback(currentPage.value - 1);
};

const handleNextPage = () => {
  if (hasNextPage.value) loadFeedback(currentPage.value + 1);
};

const handlePageJump = () => {
  const page = Number(targetPage.value);
  if (!Number.isInteger(page) || page < 1) {
    ElMessage.warning('请输入合法的页码');
    targetPage.value = currentPage.value;
    return;
  }
  if (page !== currentPage.value) loadFeedback(page);
};

onMounted(() => loadFeedback(1));
</script>

<style>
/* 反馈记录使用独立的高对比度暗色面板，避免全局主题覆盖后出现浅底浅字。 */
html[data-theme='dark'] .feedback-records,
html[data-theme='effect'] .feedback-records {
  color: #e2e8f0;
}

html[data-theme='dark'] .feedback-records .feedback-card,
html[data-theme='effect'] .feedback-records .feedback-card {
  border-color: rgba(148, 163, 184, 0.34) !important;
  background: rgba(15, 23, 42, 0.86) !important;
}

html[data-theme='dark'] .feedback-records .feedback-card__header,
html[data-theme='effect'] .feedback-records .feedback-card__header {
  border-color: rgba(148, 163, 184, 0.28) !important;
  background: rgba(30, 41, 59, 0.92) !important;
}

html[data-theme='dark'] .feedback-records .feedback-card__title,
html[data-theme='effect'] .feedback-records .feedback-card__title {
  color: #f8fafc !important;
}

html[data-theme='dark'] .feedback-records .feedback-card__time,
html[data-theme='effect'] .feedback-records .feedback-card__time,
html[data-theme='dark'] .feedback-records .feedback-card__value,
html[data-theme='effect'] .feedback-records .feedback-card__value {
  color: #e2e8f0 !important;
}

html[data-theme='dark'] .feedback-records .feedback-card__label,
html[data-theme='effect'] .feedback-records .feedback-card__label,
html[data-theme='dark'] .feedback-records .feedback-records__muted,
html[data-theme='effect'] .feedback-records .feedback-records__muted {
  color: #a8b7ca !important;
}

html[data-theme='dark'] .feedback-records .feedback-card__links,
html[data-theme='effect'] .feedback-records .feedback-card__links {
  border-color: rgba(96, 165, 250, 0.52) !important;
  background: rgba(15, 23, 42, 0.72) !important;
}

html[data-theme='dark'] .feedback-records .feedback-card__links-title,
html[data-theme='effect'] .feedback-records .feedback-card__links-title {
  color: #93c5fd !important;
}

html[data-theme='dark'] .feedback-records .feedback-card__link,
html[data-theme='effect'] .feedback-records .feedback-card__link {
  border-color: rgba(96, 165, 250, 0.72) !important;
  background: rgba(30, 41, 59, 0.92) !important;
  color: #bfdbfe !important;
  text-decoration-color: rgba(147, 197, 253, 0.72) !important;
}

html[data-theme='dark'] .feedback-records .feedback-card__link:hover,
html[data-theme='effect'] .feedback-records .feedback-card__link:hover {
  background: rgba(37, 99, 235, 0.28) !important;
  color: #dbeafe !important;
}

html[data-theme='dark'] .feedback-records .feedback-card__empty,
html[data-theme='effect'] .feedback-records .feedback-card__empty {
  color: #a8b7ca !important;
}

html[data-theme='dark'] .feedback-records .border-gray-300,
html[data-theme='effect'] .feedback-records .border-gray-300 {
  border-color: rgba(148, 163, 184, 0.46) !important;
}

html[data-theme='dark'] .feedback-records .border-gray-300:not(:disabled),
html[data-theme='effect'] .feedback-records .border-gray-300:not(:disabled) {
  color: #e2e8f0;
  background: rgba(30, 41, 59, 0.82);
}

html[data-theme='dark'] .feedback-records .feedback-records__page-input,
html[data-theme='effect'] .feedback-records .feedback-records__page-input {
  border-color: rgba(148, 163, 184, 0.46) !important;
  color: #f8fafc !important;
  background: rgba(15, 23, 42, 0.86) !important;
}

html[data-theme='dark'] .feedback-records .feedback-records__error,
html[data-theme='effect'] .feedback-records .feedback-records__error {
  border: 1px solid rgba(248, 113, 113, 0.45);
  color: #fecaca !important;
  background: rgba(127, 29, 29, 0.42) !important;
}
</style>

<script>
export default {
  name: 'FeedbackRecords'
};
</script>
