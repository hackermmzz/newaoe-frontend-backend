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
          @click="loadStudents(currentPage)"
        >
          {{ isLoading ? '加载中...' : '刷新列表' }}
        </button>
      </div>
    </div>

    <div v-if="errorMessage" class="rounded-lg bg-red-50 p-4 text-sm text-red-700" role="alert">
      {{ errorMessage }}
    </div>

    <!-- 分页控件放在学生列表上方，便于直接翻页。 -->
    <div
      v-if="students.length && !isLoading"
      class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-y border-gray-200 py-4"
    >
      <p class="text-sm text-gray-500">
        第 {{ currentPage }} 页，本页 {{ students.length }} 条
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
            class="w-16 px-2 py-1.5 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-1 focus:ring-blue-500"
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

    <div v-if="isLoading" class="py-16 text-center text-gray-500">
      正在加载学生信息...
    </div>

    <div v-else-if="!students.length" class="py-16 text-center text-gray-500">
      暂无学生数据
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

  </section>
</template>

<script>
import { onMounted, ref } from 'vue';
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
    const targetPage = ref(1);
    const hasNextPage = ref(true);
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

    const loadStudents = async (page = currentPage.value) => {
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
          if (page === currentPage.value + 1) {
            hasNextPage.value = false;
          }
          targetPage.value = currentPage.value;
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

        students.value = nextStudents;
        currentPage.value = page;
        targetPage.value = page;
        hasNextPage.value = payload.length >= pageSize;
      } catch (error) {
        targetPage.value = currentPage.value;
        errorMessage.value = error?.message || '加载学生信息失败';
        ElMessage.error(errorMessage.value);
      } finally {
        isLoading.value = false;
      }
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

    const handlePrevPage = () => {
      if (currentPage.value > 1) loadStudents(currentPage.value - 1);
    };

    const handleNextPage = () => {
      if (hasNextPage.value) loadStudents(currentPage.value + 1);
    };

    const handlePageJump = () => {
      const page = Number(targetPage.value);
      if (!Number.isInteger(page) || page < 1) {
        errorMessage.value = '请输入合法的页码';
        targetPage.value = currentPage.value;
        return;
      }
      if (page !== currentPage.value) loadStudents(page);
    };

    onMounted(loadStudents);

    return {
      errorMessage,
      currentPage,
      exportStudentInfo,
      formatDate,
      handleNextPage,
      handlePageJump,
      handlePrevPage,
      hasNextPage,
      isLoading,
      isExporting,
      loadStudents,
      openStudentHistory,
      targetPage,
      students
    };
  }
};
</script>
