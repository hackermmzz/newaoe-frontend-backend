<template>
  <section class="space-y-6">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h2 class="text-2xl font-semibold text-gray-800">学生统计</h2>
        <p class="mt-1 text-sm text-gray-500">点击学生可查看该学生的提交历史记录。</p>
      </div>
      <div class="flex items-center gap-2">
        <button
          type="button"
          class="px-4 py-2 rounded-lg bg-emerald-600 text-white hover:bg-emerald-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          :disabled="isExporting"
          @click="exportStudentInfo"
        >
          {{ isExporting ? '导出中...' : '导出信息' }}
        </button>
        <button
          type="button"
          class="px-4 py-2 rounded-lg bg-blue-600 text-white hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          :disabled="isLoading"
          @click="refreshStudents"
        >
          {{ isLoading ? '加载中...' : '刷新列表' }}
        </button>
      </div>
    </div>

    <form
      class="flex flex-col sm:flex-row gap-2"
      role="search"
      @submit.prevent="searchStudent"
    >
      <label for="student-search-id" class="sr-only">按学生 ID 搜索</label>
      <input
        id="student-search-id"
        v-model="searchId"
        type="text"
        autocomplete="off"
        placeholder="请输入学生 ID"
        class="flex-1 px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
      <button
        type="submit"
        class="px-4 py-2 rounded-lg bg-indigo-600 text-white hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        :disabled="isLoading"
      >
        {{ isLoading ? '搜索中...' : '搜索' }}
      </button>
      <button
        v-if="activeSearchId"
        type="button"
        class="px-4 py-2 rounded-lg border border-gray-300 text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        :disabled="isLoading"
        @click="clearSearch"
      >
        清除搜索
      </button>
    </form>

    <div v-if="errorMessage" class="rounded-lg bg-red-50 p-4 text-sm text-red-700" role="alert">
      {{ errorMessage }}
    </div>

    <div v-if="isLoading && !students.length" class="py-16 text-center text-gray-500">
      正在加载学生信息...
    </div>

    <div v-else-if="!students.length" class="py-16 text-center text-gray-500">
      {{ activeSearchId ? '未找到该学生' : '暂无学生数据' }}
    </div>

    <div v-else class="grid grid-cols-1 gap-4">
      <button
        v-for="student in students"
        :key="student.id"
        type="button"
        class="bg-white rounded-xl border border-gray-200 p-5 text-left shadow-sm hover:shadow-md hover:border-blue-300 transition-all"
        @click="openStudentHistory(student.id)"
      >
        <div class="flex items-center gap-4">
          <img
            v-if="student.avatarUrl"
            :src="student.avatarUrl"
            :alt="`${student.id} 的头像`"
            class="w-16 h-16 rounded-full object-cover border border-gray-200 bg-gray-50"
          >
          <div
            v-else
            class="w-16 h-16 rounded-full flex items-center justify-center bg-gray-100 text-gray-400"
            aria-hidden="true"
          >
            <i class="fa fa-user text-2xl"></i>
          </div>
          <div class="min-w-0">
            <p class="font-semibold text-gray-800 truncate">{{ student.id }}</p>
            <p class="mt-1 text-sm text-gray-500">注册时间</p>
            <p class="text-sm text-gray-600 truncate">{{ formatDate(student.registDate) }}</p>
          </div>
        </div>
      </button>
    </div>

    <div
      ref="loadMoreTrigger"
      v-show="!activeSearchId && hasNextPage"
      class="py-5 text-center text-sm text-gray-500"
      aria-live="polite"
    >
      {{ isLoading ? '正在加载更多学生...' : '下滑加载更多' }}
    </div>

  </section>
</template>

<script>
import { onBeforeUnmount, onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import config from '../config';
import { getDownloadUrl } from '../utils/download';

export default {
  name: 'StudentStatistics',

  setup() {
    const router = useRouter();
    const students = ref([]);
    const isLoading = ref(false);
    const isExporting = ref(false);
    const errorMessage = ref('');
    const currentPage = ref(1);
    const hasNextPage = ref(true);
    const searchId = ref('');
    const activeSearchId = ref('');
    const loadMoreTrigger = ref(null);
    let loadMoreObserver = null;
    const pageSize = config.ManagerStudentRecordPerPage || 10;

    const formatDate = value => {
      if (value === null || value === undefined || value === '') return '未知';
      const date = new Date(value);
      return Number.isNaN(date.getTime())
        ? String(value)
        : date.toLocaleString('zh-CN', { hour12: false });
    };

    const normalizeStudent = item => ({
      id: String(item?.id ?? '').trim(),
      registDate: item?.registDate ?? item?.registData ?? item?.regist_date ?? '',
      avatar: String(item?.avatar ?? '').trim(),
      avatarUrl: ''
    });

    const loadStudents = async (page = currentPage.value, append = false) => {
      if (isLoading.value || !Number.isInteger(page) || page < 1) return;
      isLoading.value = true;
      errorMessage.value = '';
      try {
        const beg = (page - 1) * pageSize;
        const end = beg + pageSize - 1;
        const requestUrl = new URL(config.manager_student_url);
        requestUrl.searchParams.set('range', `${beg}:${end}`);

        const response = await fetch(requestUrl.toString(), {
          method: 'GET',
          credentials: 'include'
        });
        const data = await response.json();
        if (!response.ok || !data?.status) {
          throw new Error(data?.msg || `加载学生信息失败（HTTP ${response.status}）`);
        }

        if (!Array.isArray(data?.data)) {
          throw new Error('不可跳转');
        }
        const payload = data.data;

        if (payload.length === 0 && page > 1) {
          hasNextPage.value = false;
          return;
        }

        const normalized = payload
          .map(normalizeStudent)
          .filter(student => student.id);

        const nextStudents = await Promise.all(normalized.map(async student => {
          if (!student.avatar) return student;
          try {
            return {
              ...student,
              avatarUrl: await getDownloadUrl(student.avatar)
            };
          } catch (error) {
            return student;
          }
        }));

        students.value = append
          ? [...students.value, ...nextStudents]
          : nextStudents;
        currentPage.value = page;
        hasNextPage.value = payload.length >= pageSize;
      } catch (error) {
        errorMessage.value = error?.message || '加载学生信息失败';
        ElMessage.error(errorMessage.value);
      } finally {
        isLoading.value = false;
      }
    };

    const searchStudent = async () => {
      const studentId = String(searchId.value ?? '').trim();
      if (!studentId) {
        clearSearch();
        return;
      }
      if (isLoading.value) return;

      isLoading.value = true;
      errorMessage.value = '';
      try {
        const requestUrl = new URL(config.manager_student_search_url);
        requestUrl.searchParams.set('id', studentId);
        const response = await fetch(requestUrl.toString(), {
          method: 'GET',
          credentials: 'include'
        });
        const data = await response.json();
        if (!response.ok || !data?.status) {
          throw new Error(data?.msg || `搜索学生失败（HTTP ${response.status}）`);
        }

        if (!data?.data) {
          throw new Error('学生不存在或暂无信息');
        }

        const student = normalizeStudent(data.data);
        if (!student.id) {
          throw new Error('未找到该学生');
        }
        if (student.avatar) {
          try {
            student.avatarUrl = await getDownloadUrl(student.avatar);
          } catch (error) {
            // 头像下载失败不影响学生信息展示。
          }
        }

        students.value = [student];
        activeSearchId.value = studentId;
        currentPage.value = 1;
        hasNextPage.value = false;
      } catch (error) {
        students.value = [];
        activeSearchId.value = studentId;
        errorMessage.value = error?.message || '搜索学生失败';
        ElMessage.error(errorMessage.value);
      } finally {
        isLoading.value = false;
      }
    };

    const clearSearch = () => {
      if (isLoading.value) return;
      searchId.value = '';
      activeSearchId.value = '';
      currentPage.value = 1;
      hasNextPage.value = true;
      loadStudents(1);
    };

    const refreshStudents = () => {
      if (activeSearchId.value) {
        searchStudent();
        return;
      }
      currentPage.value = 1;
      hasNextPage.value = true;
      loadStudents(1);
    };

    const loadNextStudents = () => {
      if (activeSearchId.value || !hasNextPage.value || isLoading.value) return;
      loadStudents(currentPage.value + 1, true);
    };

    const openStudentHistory = async id => {
      const studentId = String(id ?? '').trim();
      if (!studentId) {
        ElMessage.error('不可跳转');
        return;
      }

      try {
        await router.push({
          name: 'ManagerStudentHistory',
          params: { studentId }
        });
      } catch (error) {
        ElMessage.error('不可跳转');
      }
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

        if (!downloadReference) {
          throw new Error(result?.msg || '下载链接无效');
        }

        const downloadUrl = new URL(downloadReference, `${config.base_url}/`).toString();
        triggerDownload(downloadUrl);
        ElMessage.success('学生信息导出成功，正在下载');
      } catch (error) {
        ElMessage.error(error?.message || '网络异常');
      } finally {
        isExporting.value = false;
      }
    };

    onMounted(() => {
      loadMoreObserver = new IntersectionObserver(
        entries => {
          if (entries.some(entry => entry.isIntersecting)) {
            loadNextStudents();
          }
        },
        { rootMargin: '240px 0px' }
      );
      if (loadMoreTrigger.value) {
        loadMoreObserver.observe(loadMoreTrigger.value);
      }
      loadStudents(1);
    });

    onBeforeUnmount(() => {
      loadMoreObserver?.disconnect();
    });

    return {
      errorMessage,
      activeSearchId,
      currentPage,
      exportStudentInfo,
      formatDate,
      hasNextPage,
      isLoading,
      isExporting,
      loadStudents,
      openStudentHistory,
      refreshStudents,
      searchId,
      searchStudent,
      clearSearch,
      loadMoreTrigger,
      students
    };
  }
};
</script>
