<template>
  <section class="space-y-6" :aria-busy="isLoading">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h3 class="text-xl font-semibold text-gray-800">学生排行榜</h3>
        <p class="mt-2 text-sm text-gray-500">
          排行榜顺序由后端提供，页面按接口返回顺序展示。
        </p>
      </div>

      <button
        type="button"
        class="ranking-button bg-blue-600 text-white hover:bg-blue-700"
        :disabled="isLoading"
        @click="loadRanking(currentPage)"
      >
        {{ isLoading ? '加载中...' : '刷新榜单' }}
      </button>
    </div>

    <!-- 错误提示 -->
    <div
      v-if="errorMessage"
      class="bg-red-50 text-red-700 rounded-lg p-4 text-sm"
      role="alert"
    >
      {{ errorMessage }}

      <button
        type="button"
        class="ml-3 underline"
        :disabled="isLoading"
        @click="loadRanking(requestedPage)"
      >
        重新加载
      </button>
    </div>

    <div class="bg-white rounded-lg border border-gray-200 overflow-hidden">
      <!-- 标题栏 -->
      <div
        class="p-4 border-b border-gray-100 flex flex-col sm:flex-row sm:items-center justify-between gap-2"
      >
        <h4 class="font-medium text-gray-800">对战成绩榜单</h4>

        <p class="text-sm text-gray-500">
          最近更新：{{ updatedAt || '—' }}
        </p>
      </div>

      <!-- 加载状态 -->
      <div
        v-if="isLoading"
        class="py-16 text-center text-gray-500"
        role="status"
      >
        <div
          class="inline-block h-8 w-8 animate-spin rounded-full border-2 border-gray-200 border-t-blue-600"
        ></div>

        <p class="mt-3">
          正在加载第 {{ requestedPage }} 页...
        </p>
      </div>

      <!-- 空数据 -->
      <div
        v-else-if="!rankedRecords.length"
        class="py-16 px-4 text-center text-gray-500"
        role="status"
      >
        {{ errorMessage ? '榜单加载失败，请重试' : '本页暂无排行数据' }}
      </div>

      <!-- 排行榜 -->
      <div v-else class="overflow-x-auto">
        <table class="w-full text-sm text-left">
          <caption class="sr-only">
            学生对战排行榜
          </caption>

          <thead class="bg-gray-50 text-gray-500 whitespace-nowrap">
            <tr>
              <th scope="col" class="px-5 py-4 font-medium">
                排名
              </th>

              <th scope="col" class="px-5 py-4 font-medium">
                ID
              </th>

              <th scope="col" class="px-5 py-4 font-medium">
                胜利/失败
              </th>

              <th scope="col" class="px-5 py-4 font-medium">
                Score
              </th>

              <th scope="col" class="px-5 py-4 font-medium">
                Frame
              </th>

              <th scope="col" class="px-5 py-4 font-medium">
                提交时间
              </th>

              <th scope="col" class="px-5 py-4 font-medium">
                消息
              </th>
            </tr>
          </thead>

          <tbody class="divide-y divide-gray-100">
            <tr
              v-for="(student, index) in rankedRecords"
              :key="student.id"
              class="hover:bg-gray-50 transition-colors align-top"
            >
              <!-- 全局排名 -->
              <td class="px-5 py-4 font-semibold text-gray-700 whitespace-nowrap">
                {{ (currentPage - 1) * pageSize + index + 1 }}
              </td>

              <!-- 学生 ID -->
              <td class="px-5 py-4 font-medium text-gray-800 whitespace-nowrap">
                {{ student.id }}
              </td>

              <!-- 胜负 -->
              <td class="px-5 py-4 whitespace-nowrap">
                <span
                  class="inline-block px-2 py-1 rounded-full text-xs font-medium"
                  :class="
                    student.win
                      ? 'bg-green-100 text-green-700'
                      : 'bg-red-100 text-red-700'
                  "
                >
                  {{ student.win ? '胜利' : '失败' }}
                </span>
              </td>

              <!-- 分数 -->
              <td class="px-5 py-4 font-semibold text-blue-600">
                {{ student.score }}
              </td>

              <!-- Frame -->
              <td class="px-5 py-4 text-gray-600">
                {{ student.frame }}
              </td>

              <!-- 提交时间 -->
              <td class="px-5 py-4 text-gray-600 whitespace-nowrap">
                {{ formatTime(student.submittime) }}
              </td>

              <!-- 消息 -->
              <td class="px-5 py-4 text-gray-600">
                <div class="description-text">
                  {{
                    expandedItems[index]
                      ? student.msg
                      : truncateDescription(student.msg)
                  }}
                </div>

                <button
                  v-if="student.msg.length > descriptionLimit"
                  type="button"
                  class="mt-2 text-blue-600 hover:text-blue-800"
                  :aria-expanded="!!expandedItems[index]"
                  @click="
                    expandedItems[index] = !expandedItems[index]
                  "
                >
                  {{
                    expandedItems[index]
                      ? '收起'
                      : '展开全部'
                  }}
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- 分页 -->
      <div
        class="p-4 border-t border-gray-100 flex flex-col lg:flex-row lg:items-center justify-between gap-4"
      >
        <p class="text-sm text-gray-500">
          第 {{ currentPage }} 页，本页 {{ records.length }} 条
        </p>

        <div class="flex flex-wrap items-center gap-3">
          <button
            type="button"
            class="ranking-button border border-gray-300 hover:bg-gray-50"
            :disabled="isLoading || currentPage <= 1"
            @click="loadRanking(currentPage - 1)"
          >
            上一页
          </button>

          <span class="text-sm text-gray-600">
            第 {{ currentPage }} 页
          </span>

          <button
            type="button"
            class="ranking-button border border-gray-300 hover:bg-gray-50"
            :disabled="isLoading || !hasNextPage"
            @click="loadRanking(currentPage + 1)"
          >
            下一页
          </button>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import {
  computed,
  onBeforeUnmount,
  onMounted,
  ref
} from 'vue';

import config from '../config';

export default {
  name: 'StudentRanking',

  setup() {
    const records = ref([]);

    const isLoading = ref(false);

    const errorMessage = ref('');

    const updatedAt = ref('');

    const currentPage = ref(1);

    const requestedPage = ref(1);

    const hasNextPage = ref(true);

    const pageSize =
      config.RankingRecordPerPage || 10;

    const descriptionLimit = 100;

    const expandedItems = ref({});

    let controller = null;

    let disposed = false;

    /*
     * 后端已经负责排序。
     *
     * 后端排序规则例如：
     *
     * win DESC
     * score DESC
     * frame ASC
     * id ASC
     *
     * 所以前端绝对不要再次排序，
     * 直接按照接口返回顺序显示。
     */
    const rankedRecords = computed(() => {
      return records.value;
    });

    /*
     * 检查并规范后端返回的数据。
     *
     * 正确格式：
     *
     * {
     *   id: "923106840404",
     *   win: true,
     *   score: 100,
     *   frame: 98,
     *   submittime: "2026-09-12T20:07:41Z",
     *   msg: "Accepted"
     * }
     */
    const normalizeRecord = item => {
      if (!item || typeof item !== 'object') {
        throw new Error('排行榜记录格式不正确');
      }

      if (
        typeof item.id !== 'string' ||
        !item.id.trim()
      ) {
        throw new Error(
          '排行榜记录 id 格式不正确'
        );
      }

      /*
       * 注意：
       *
       * 你原来的代码这里检查的是：
       *
       * item.success
       *
       * 但后端实际返回的是：
       *
       * item.win
       *
       * 所以必须改成 win。
       */
      if (typeof item.win !== 'boolean') {
        throw new Error(
          '排行榜记录 win 格式不正确'
        );
      }

      if (!Number.isSafeInteger(item.score)) {
        throw new Error(
          '排行榜记录 score 格式不正确'
        );
      }

      if (
        !Number.isSafeInteger(item.frame) ||
        item.frame < 0
      ) {
        throw new Error(
          '排行榜记录 frame 格式不正确'
        );
      }

      return {
        id: item.id,
        win: item.win,
        score: item.score,
        frame: item.frame,
        submittime:
          item.submittime ?? '',
        msg: String(item.msg ?? '')
      };
    };

    /*
     * 获取排行榜。
     */
    const loadRanking = async (
      page = currentPage.value
    ) => {
      if (
        isLoading.value ||
        disposed ||
        !Number.isInteger(page) ||
        page < 1
      ) {
        return;
      }

      isLoading.value = true;

      requestedPage.value = page;

      errorMessage.value = '';

      /*
       * 如果前一个请求还存在，
       * 先取消。
       */
      controller?.abort();

      const requestController =
        new AbortController();

      controller = requestController;

      /*
       * 15 秒超时。
       */
      const timeout = setTimeout(() => {
        requestController.abort();
      }, 15000);

      try {
        /*
         * 假设每页 10 条：
         *
         * 第 1 页：
         * beg = 0
         * end = 10
         *
         * 第 2 页：
         * beg = 10
         * end = 20
         *
         * 因为你的 Go 后端是：
         *
         * Limit(end - beg, beg)
         *
         * 所以 end 应该是开区间，
         * 不能写 beg + pageSize - 1。
         */
        const beg =
          (page - 1) * pageSize;

        const end =
          beg + pageSize;

        const requestUrl =
          new URL(
            config.ranking_url +
            '/fetchrank'
          );

        requestUrl.searchParams.set(
          'range',
          `${beg}:${end}`
        );

        const response = await fetch(
          requestUrl.toString(),
          {
            method: 'GET',

            credentials: 'include',

            signal:
              requestController.signal
          }
        );

        if (!response.ok) {
          throw new Error(
            `加载失败（HTTP ${response.status}），请稍后重试`
          );
        }

        const result =
          await response.json();

        if (disposed) {
          return;
        }

        /*
         * 假设后端响应：
         *
         * {
         *   "data": [...]
         * }
         */
        const nextData =
          result?.data;

        console.log(
          '排行榜接口返回：',
          result
        );

        if (!Array.isArray(nextData)) {
          throw new Error(
            '排行榜响应中的 data 必须是数组'
          );
        }

        /*
         * 如果请求下一页，
         * 结果已经为空，
         * 那么仍然停留在当前页。
         */
        if (
          nextData.length === 0 &&
          page > 1
        ) {
          hasNextPage.value = false;

          return;
        }

        /*
         * 数据格式检查。
         */
        const nextRecords =
          nextData.map(normalizeRecord);

        records.value =
          nextRecords;

        currentPage.value =
          page;

        /*
         * 如果这一页不足 pageSize，
         * 基本可以确定没有下一页了。
         *
         * 例如 pageSize = 10，
         * 当前只返回 6 条，
         * 下一页肯定为空。
         */
        hasNextPage.value =
          nextData.length === pageSize;

        /*
         * 翻页后关闭之前展开的 msg。
         */
        expandedItems.value = {};

        /*
         * 更新时间。
         */
        updatedAt.value =
          new Date().toLocaleString(
            'zh-CN',
            {
              hour12: false
            }
          );
      } catch (error) {
        if (disposed) {
          return;
        }

        if (
          error.name === 'AbortError'
        ) {
          errorMessage.value =
            '请求超时，请重新加载';
        } else {
          errorMessage.value =
            error.message ||
            '加载失败，请检查网络后重试';
        }
      } finally {
        clearTimeout(timeout);

        /*
         * 防止已经被新请求替换后，
         * 旧请求错误地修改状态。
         */
        if (
          controller ===
          requestController
        ) {
          controller = null;

          isLoading.value = false;
        }
      }
    };

    /*
     * msg 太长时进行截断。
     */
    const truncateDescription = text => {
      const value =
        String(text ?? '');

      if (!value) {
        return '—';
      }

      if (
        value.length >
        descriptionLimit
      ) {
        return (
          value.slice(
            0,
            descriptionLimit
          ) + '…'
        );
      }

      return value;
    };

    /*
     * 格式化后端时间。
     */
    const formatTime = value => {
      if (
        value === null ||
        value === undefined ||
        value === ''
      ) {
        return '—';
      }

      const date =
        new Date(value);

      if (
        Number.isNaN(
          date.getTime()
        )
      ) {
        return String(value);
      }

      return date.toLocaleString(
        'zh-CN',
        {
          hour12: false
        }
      );
    };

    /*
     * 页面加载时获取第一页。
     */
    onMounted(() => {
      loadRanking(1);
    });

    /*
     * 页面销毁时终止请求。
     */
    onBeforeUnmount(() => {
      disposed = true;

      controller?.abort();
    });

    return {
      records,

      rankedRecords,

      isLoading,

      errorMessage,

      updatedAt,

      currentPage,

      requestedPage,

      hasNextPage,

      pageSize,

      descriptionLimit,

      expandedItems,

      loadRanking,

      truncateDescription,

      formatTime
    };
  }
};
</script>

<style scoped>
.ranking-button {
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
  font-size: 0.875rem;
  transition:
    background-color 0.2s,
    opacity 0.2s;
}

.ranking-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.ranking-button:focus-visible {
  outline: 2px solid #3b82f6;
  outline-offset: 2px;
}

.description-text {
  min-width: 12rem;
  max-width: 28rem;
  white-space: pre-wrap;
  overflow-wrap: anywhere;
}
</style>