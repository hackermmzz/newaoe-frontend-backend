<template>
  <section class="space-y-6" :aria-busy="isLoading">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h3 class="text-xl font-semibold text-gray-800">考核排行榜</h3>
        <p class="mt-2 text-sm text-gray-500">排行榜顺序由后端提供，页面按接口返回顺序展示。</p>
      </div>
      <div class="flex items-center gap-2">
        <button
          type="button"
          class="ranking-button bg-emerald-600 text-white hover:bg-emerald-700"
          :disabled="isExporting"
          @click="exportStudentInfo"
        >
          {{ isExporting ? '导出中...' : '导出信息' }}
        </button>
        <button
          type="button"
          class="ranking-button bg-blue-600 text-white hover:bg-blue-700"
          :disabled="isLoading"
          @click="refreshRanking"
        >
          {{ isLoading ? '加载中...' : '刷新榜单' }}
        </button>
      </div>
    </div>

    <div v-if="errorMessage" class="bg-red-50 text-red-700 rounded-lg p-4 text-sm" role="alert">
      {{ errorMessage }}
      <button type="button" class="ml-3 underline" :disabled="isLoading" @click="refreshRanking">
        重新加载
      </button>
    </div>

    <div class="bg-white rounded-lg border border-gray-200 overflow-hidden">
      <div class="p-4 border-b border-gray-100 flex flex-col sm:flex-row sm:items-center justify-between gap-2">
        <h4 class="font-medium text-gray-800">考核成绩榜单</h4>
        <p class="text-sm text-gray-500">最近更新：{{ updatedAt || '—' }}</p>
      </div>
      <div class="p-4 border-b border-gray-100">
        <p class="text-sm text-gray-500">已加载 {{ records.length }} 条</p>
      </div>

      <div v-if="championRecord" class="champion-card">
        <img
          v-if="championRecord.avatarUrl"
          :src="championRecord.avatarUrl"
          :alt="`${championRecord.id} 的头像`"
          class="champion-avatar"
        >
        <div v-else class="champion-avatar champion-avatar-placeholder" aria-hidden="true">
          <i class="fa fa-user"></i>
        </div>
        <span class="champion-card-title">🏆 冠军</span>
        <span class="champion-card-id">ID：{{ championRecord.id }}</span>
      </div>

      <div v-if="isLoading && !rankedRecords.length" class="py-16 text-center text-gray-500" role="status">
        <div class="inline-block h-8 w-8 animate-spin rounded-full border-2 border-gray-200 border-t-blue-600"></div>
        <p class="mt-3">正在加载排行榜...</p>
      </div>
      <div v-else-if="!rankedRecords.length" class="py-16 px-4 text-center text-gray-500" role="status">
        {{ errorMessage ? '榜单加载失败，请重试' : '本页暂无排行数据' }}
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-sm text-left">
          <caption class="sr-only">学生考核排行榜</caption>
          <thead class="bg-gray-50 text-gray-500 whitespace-nowrap">
            <tr>
              <th scope="col" class="px-5 py-4 font-medium">头像</th>
              <th scope="col" class="px-5 py-4 font-medium">排名</th>
              <th scope="col" class="px-5 py-4 font-medium">ID</th>
              <th scope="col" class="px-5 py-4 font-medium">状态</th>
              <th scope="col" class="px-5 py-4 font-medium">Score</th>
              <th scope="col" class="px-5 py-4 font-medium">Frame</th>
              <th scope="col" class="px-5 py-4 font-medium">提交时间</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr v-for="(student, index) in rankedRecords" :key="student.id" class="hover:bg-gray-50 transition-colors align-top">
              <td class="px-5 py-4 whitespace-nowrap">
                <img
                  v-if="student.avatarUrl"
                  :src="student.avatarUrl"
                  :alt="`${student.id} 的头像`"
                  class="h-10 w-10 rounded-full object-cover border border-gray-200"
                >
                <span v-else class="text-gray-400" aria-label="暂无头像">—</span>
              </td>
              <td class="px-5 py-4 font-semibold text-gray-700 whitespace-nowrap">
                <div class="flex items-center gap-2">
                  <span class="rank-medal" :aria-label="getMedalLabel(index + 1)">{{ getMedal(index + 1) }}</span>
                  <span v-if="index + 1 <= 3" class="medal-title">{{ getMedalLabel(index + 1) }}</span>
                  <span>{{ index + 1 }}</span>
                  <span v-if="getStars(index + 1)" class="rank-stars" aria-label="荣誉星级">{{ getStars(index + 1) }}</span>
                </div>
              </td>
              <td class="px-5 py-4 font-medium text-gray-800 whitespace-nowrap">
                {{ student.id }}
                <span v-if="index === 0" class="champion-logo" title="冠军">🏆 冠军</span>
              </td>
              <td class="px-5 py-4 whitespace-nowrap">
                <span class="inline-block px-2 py-1 rounded-full text-xs font-medium" :class="statusClass(student.outcome)">
                  {{ statusLabel(student.outcome) }}
                </span>
              </td>
              <td class="px-5 py-4 font-semibold text-blue-600">{{ student.score }}</td>
              <td class="px-5 py-4 text-gray-600">{{ student.frame }}</td>
              <td class="px-5 py-4 text-gray-600 whitespace-nowrap">{{ formatTime(student.submittime) }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div ref="loadMoreTrigger" v-show="hasNextPage" class="py-5 text-center text-sm text-gray-500" aria-live="polite">
        {{ isLoading ? '正在加载更多排行记录...' : '下滑加载更多' }}
      </div>
    </div>
  </section>
</template>

<script>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue';
import config from '../config';
import { ElMessage } from 'element-plus';
import { getDownloadUrl } from '../utils/download';

export default {
  name: 'AssessmentRank',
  setup() {
    const records = ref([]);
    const isLoading = ref(false);
    const isExporting = ref(false);
    const errorMessage = ref('');
    const updatedAt = ref('');
    const currentPage = ref(1);
    const hasNextPage = ref(true);
    const loadMoreTrigger = ref(null);
    const pageSize = config.RankingRecordPerPage || 10;
    let controller = null;
    let observer = null;
    let disposed = false;

    const rankedRecords = computed(() => records.value);
    const championRecord = computed(() => rankedRecords.value[0] || null);

    const getMedal = rank => rank === 1 ? '🥇' : rank === 2 ? '🥈' : rank === 3 ? '🥉' : '';
    const getMedalLabel = rank => rank === 1 ? '金牌' : rank === 2 ? '银牌' : rank === 3 ? '铜牌' : `第 ${rank} 名`;
    const getStars = rank => rank <= 10 ? '★★★' : rank <= 20 ? '★★' : rank <= 30 ? '★' : '';

    const outcomeByStatusCode = new Map([
      [config.Code_Status_Unknown, 'unknown'],
      [config.Code_Status_Error, 'error'],
      [config.Code_Status_Wait, 'wait'],
      [config.Code_Status_Compile, 'compile'],
      [config.Code_Status_Compile_Fail, 'compile-fail'],
      [config.Code_Status_Compile_Success, 'compile-success'],
      [config.Code_Status_Running, 'running'],
      [config.Code_Status_Crash, 'crash'],
      [config.Code_Status_Fail, 'fail'],
      [config.Code_Status_Success, 'win']
    ]);

    const statusPresentations = {
      unknown: { label: '未知', className: 'bg-gray-100 text-gray-700' },
      error: { label: '服务器异常', className: 'bg-rose-100 text-rose-700' },
      wait: { label: '等待中', className: 'bg-blue-100 text-blue-700' },
      compile: { label: '编译中', className: 'bg-indigo-100 text-indigo-700' },
      'compile-fail': { label: '编译错误', className: 'bg-orange-100 text-orange-700' },
      'compile-success': { label: '编译成功', className: 'bg-teal-100 text-teal-700' },
      running: { label: '运行中', className: 'bg-cyan-100 text-cyan-700' },
      crash: { label: '游戏崩溃', className: 'bg-yellow-100 text-yellow-700' },
      fail: { label: '失败', className: 'bg-red-100 text-red-700' },
      win: { label: '胜利', className: 'bg-green-100 text-green-700' }
    };

    const getStatusCode = value => {
      if (typeof value === 'number' && Number.isInteger(value)) return value;
      if (typeof value === 'string') {
        const normalized = value.trim();
        if (/^-?\d+$/.test(normalized)) return Number(normalized);
        try {
          return getStatusCode(JSON.parse(normalized));
        } catch (error) {
          return NaN;
        }
      }
      if (value && typeof value === 'object' && !Array.isArray(value)) {
        return getStatusCode(value.status ?? value.Status ?? value.code ?? value.Code);
      }
      return NaN;
    };

    const normalizeRecord = item => {
      if (!item || typeof item !== 'object') throw new Error('考核排行榜记录格式不正确');
      const rankInfo = item.rankinfo && typeof item.rankinfo === 'object' ? item.rankinfo : item;
      if (typeof rankInfo.id !== 'string' || !rankInfo.id.trim()) throw new Error('考核排行榜记录 id 格式不正确');

      const code = getStatusCode(rankInfo.status ?? item.status);
      const outcome = outcomeByStatusCode.get(code);
      if (!outcome) throw new Error(`考核排行榜记录 ${rankInfo.id} 状态格式不正确`);
      if (!Number.isSafeInteger(rankInfo.score)) throw new Error('考核排行榜记录 score 格式不正确');
      if (!Number.isSafeInteger(rankInfo.frame) || rankInfo.frame < 0) throw new Error('考核排行榜记录 frame 格式不正确');

      return {
        id: rankInfo.id,
        avatar: String(item.avatar ?? '').trim(),
        outcome,
        score: rankInfo.score,
        frame: rankInfo.frame,
        submittime: rankInfo.submittime ?? ''
      };
    };

    const fetchAvatar = async (avatar, signal) => {
      if (!avatar) return '';
      try {
        return await getDownloadUrl(avatar, { signal });
      } catch (error) {
        if (error.name === 'AbortError') throw error;
        console.warn('[考核排行榜] 头像下载失败：', { avatar, error });
        return '';
      }
    };

    const resolveAvatars = async (nextRecords, signal) => {
      const requests = new Map();
      const getUrl = avatar => {
        if (!requests.has(avatar)) requests.set(avatar, fetchAvatar(avatar, signal));
        return requests.get(avatar);
      };
      return Promise.all(nextRecords.map(async record => ({
        ...record,
        avatarUrl: await getUrl(record.avatar)
      })));
    };

    const loadRanking = async (page = currentPage.value, append = false) => {
      if (isLoading.value || disposed || !Number.isInteger(page) || page < 1) return;
      isLoading.value = true;
      errorMessage.value = '';
      controller?.abort();
      const requestController = new AbortController();
      controller = requestController;
      const timeout = setTimeout(() => requestController.abort(), 15000);

      try {
        const requestUrl = new URL(config.manager_assessment_rank_url);
        const beg = (page - 1) * pageSize;
        requestUrl.searchParams.set('range', `${beg}:${beg + pageSize - 1}`);
        const response = await fetch(requestUrl.toString(), {
          method: 'GET', credentials: 'include', signal: requestController.signal
        });
        const data = await response.json();
        if (!response.ok || (!Array.isArray(data) && !data?.status)) throw new Error(data?.msg || '加载考核排行榜失败');
        const nextData = Array.isArray(data) ? data : data?.data || [];
        if (!Array.isArray(nextData)) throw new Error('考核排行榜响应中的 data 必须是数组');
        if (!nextData.length && page > 1) {
          hasNextPage.value = false;
          return;
        }
        const nextRecords = await resolveAvatars(nextData.map(normalizeRecord), requestController.signal);
        if (disposed) return;
        records.value = append ? [...records.value, ...nextRecords] : nextRecords;
        currentPage.value = page;
        hasNextPage.value = nextData.length >= pageSize;
        updatedAt.value = new Date().toLocaleString('zh-CN', { hour12: false });
      } catch (error) {
        if (disposed) return;
        ElMessage.error(error?.message || '加载考核排行榜失败');
        errorMessage.value = error?.name === 'AbortError' ? '请求超时，请重新加载' : (error?.message || '加载失败，请重试');
      } finally {
        clearTimeout(timeout);
        if (controller === requestController) {
          controller = null;
          isLoading.value = false;
        }
      }
    };

    const refreshRanking = () => {
      if (isLoading.value) return;
      currentPage.value = 1;
      hasNextPage.value = true;
      loadRanking(1);
    };

    const triggerDownload = url => {
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', '');
      link.target = '_self';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    };

    const exportStudentInfo = async () => {
      if (isExporting.value) return;
      isExporting.value = true;
      try {
        const response = await fetch(config.manager_info_export_url, {
          method: 'GET',
          credentials: 'include'
        });
        let result = {};
        try {
          result = await response.json();
        } catch (error) {
          throw new Error('导出接口返回格式错误');
        }
        if (!response.ok || result?.status === false) {
          throw new Error(result?.msg || '网络异常');
        }

        const payload = result?.data ?? result;
        const downloadReference = typeof payload === 'string'
          ? payload.trim()
          : String(
            payload?.url
            || payload?.downloadUrl
            || payload?.downloadurl
            || payload?.path
            || payload?.file
            || ''
          ).trim();
        if (!downloadReference) throw new Error(result?.msg || '下载链接无效');

        const downloadUrl = new URL(downloadReference, `${config.base_url}/`).toString();
        triggerDownload(downloadUrl);
        ElMessage.success('学生信息导出成功，正在下载');
      } catch (error) {
        ElMessage.error(error?.message || '网络异常');
      } finally {
        isExporting.value = false;
      }
    };
    const statusLabel = outcome => statusPresentations[outcome]?.label || '未知';
    const statusClass = outcome => statusPresentations[outcome]?.className || statusPresentations.unknown.className;
    const formatTime = value => {
      if (value === null || value === undefined || value === '') return '—';
      const date = new Date(value);
      return Number.isNaN(date.getTime()) ? String(value) : date.toLocaleString('zh-CN', { hour12: false });
    };

    onMounted(() => {
      observer = new IntersectionObserver(entries => {
        if (entries.some(entry => entry.isIntersecting) && !isLoading.value && hasNextPage.value) {
          loadRanking(currentPage.value + 1, true);
        }
      }, { rootMargin: '240px 0px' });
      if (loadMoreTrigger.value) observer.observe(loadMoreTrigger.value);
      loadRanking(1);
    });
    onBeforeUnmount(() => {
      disposed = true;
      controller?.abort();
      observer?.disconnect();
    });

    return {
      records, rankedRecords, championRecord, isLoading, isExporting, errorMessage, updatedAt,
      hasNextPage, loadMoreTrigger, refreshRanking, getMedal, getMedalLabel, getStars,
      statusLabel, statusClass, formatTime, exportStudentInfo
    };
  }
};
</script>

<style scoped>
.ranking-button { padding: 0.5rem 1rem; border-radius: 0.5rem; font-size: 0.875rem; transition: background-color 0.2s, opacity 0.2s; }
.ranking-button:disabled { opacity: 0.5; cursor: not-allowed; }
.champion-card { display: flex; align-items: center; gap: 0.75rem; flex-wrap: nowrap; margin: 1rem; padding: 0.875rem 1rem; border: 1px solid #cbd5e1; border-radius: 0.75rem; background: #fff; overflow: hidden; }
.champion-avatar { width: 3.5rem; height: 3.5rem; flex: 0 0 auto; border-radius: 9999px; object-fit: cover; border: 1px solid #d1d5db; background: #f3f4f6; }
.champion-avatar-placeholder { display: inline-flex; align-items: center; justify-content: center; color: #9ca3af; font-size: 1.25rem; }
.champion-card-title { color: #facc15; font-size: 0.8rem; font-weight: 700; white-space: nowrap; }
.champion-card-id { color: #1f2937; font-weight: 700; white-space: nowrap; }
.rank-medal { display: inline-flex; width: 1.5rem; justify-content: center; font-size: 1.25rem; }
.rank-stars { color: #f59e0b; letter-spacing: 0.08em; font-size: 0.8rem; }
.medal-title { color: #92400e; font-size: 0.75rem; font-weight: 700; }
.champion-logo { display: inline-flex; align-items: center; margin-left: 0.5rem; padding: 0.1rem 0.25rem; color: #facc15; font-size: 0.75rem; font-weight: 700; }
</style>
