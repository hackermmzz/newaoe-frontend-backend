<template>
  <div class="min-h-screen bg-gray-50 flex flex-col">
    <!-- 顶部固定导航栏 -->
    <header class="bg-white shadow-md sticky top-0 z-50 transition-all duration-300">
      <div class="container mx-auto px-4 py-4 flex justify-between items-center">
        <h1 class="text-2xl font-bold text-gray-800">
          {{ historyTitle }}
        </h1>
        <router-link
          v-if="backPath"
          :to="backPath"
          class="text-sm text-blue-600 hover:text-blue-800"
        >
          返回学生统计
        </router-link>
        <button 
          v-if="!isReadOnly"
          @click="toggleUploadForm"
          class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-lg transition-all duration-200 transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
        >
          <i class="fa fa-upload mr-2"></i>提交新文件
        </button>
      </div>
      
      <!-- 上传表单区域 - 条件显示 -->
      <div 
        v-if="showUploadForm && !isReadOnly"
        class="bg-gray-50 border-t border-gray-200 px-4 py-4 animate-fadeIn"
      >
        <div class="container mx-auto">
          <h2 class="text-xl font-semibold text-gray-700 mb-4">上传 .h 和 .cpp 文件</h2>
          
          <div class="bg-white p-4 rounded-lg shadow-sm border border-gray-200">
            <div class="mb-4">
              <label class="block text-gray-700 mb-2 font-medium">选择文件</label>
              <div class="flex flex-col md:flex-row gap-4">
                <div class="flex-1">
                  <label class="flex items-center justify-center w-full h-32 border-2 border-dashed border-gray-300 rounded-lg cursor-pointer bg-gray-50 hover:bg-gray-100 transition-colors">
                    <div class="flex flex-col items-center justify-center pt-5 pb-6">
                      <i class="fa fa-file-code-o text-3xl text-gray-400 mb-2"></i>
                      <p class="mb-1 text-sm text-gray-600"><span class="font-semibold">点击上传头文件 (.h)</span></p>
                      <p class="text-xs text-gray-500">或拖放文件到此处</p>
                    </div>
                    <input 
                      type="file" 
                      class="hidden" 
                      accept=".h" 
                      @change="handleFileUpload($event, 'header')"
                    />
                  </label>
                  <p v-if="headerFile" class="mt-2 text-sm text-gray-600 truncate">{{ headerFile.name }}</p>
                </div>
                
                <div class="flex-1">
                  <label class="flex items-center justify-center w-full h-32 border-2 border-dashed border-gray-300 rounded-lg cursor-pointer bg-gray-50 hover:bg-gray-100 transition-colors">
                    <div class="flex flex-col items-center justify-center pt-5 pb-6">
                      <i class="fa fa-file-code-o text-3xl text-gray-400 mb-2"></i>
                      <p class="mb-1 text-sm text-gray-600"><span class="font-semibold">点击上传源文件 (.cpp)</span></p>
                      <p class="text-xs text-gray-500">或拖放文件到此处</p>
                    </div>
                    <input 
                      type="file" 
                      class="hidden" 
                      accept=".cpp" 
                      @change="handleFileUpload($event, 'source')"
                    />
                  </label>
                  <p v-if="sourceFile" class="mt-2 text-sm text-gray-600 truncate">{{ sourceFile.name }}</p>
                </div>
              </div>
            </div>
            
            <div class="mb-4">
              <label for="description" class="block text-gray-700 mb-2 font-medium">描述 (可选)</label>
              <textarea 
                id="description"
                v-model="description"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-colors"
                rows="2"
                placeholder="请输入本次提交的描述..."
              ></textarea>
            </div>
            
            <div class="flex justify-end gap-3">
              <button 
                @click="toggleUploadForm"
                class="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors"
              >
                取消
              </button>
              <button 
                @click="submitFiles"
                :disabled="!headerFile || !sourceFile || isSubmitting"
                class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <i v-if="isSubmitting" class="fa fa-spinner fa-spin mr-2"></i>
                <span v-if="isSubmitting">提交中...</span>
                <span v-else>确认提交</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </header>

    <!-- 主要内容区域 - 历史记录列表 -->
    <main class="flex-grow container mx-auto px-4 py-8">
      <!-- 空状态显示 -->
      <div v-if="historyList.length === 0 && !isLoading" class="text-center py-16">
        <div class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-gray-100 mb-4">
          <i class="fa fa-history text-2xl text-gray-400"></i>
        </div>
        <h3 class="text-xl font-medium text-gray-700 mb-2">暂无历史记录</h3>
        <p class="text-gray-500 mb-6">
          {{ isReadOnly ? '暂无历史记录' : '点击上方"提交新文件"按钮开始上传你的第一个文件' }}
        </p>
        <button 
          v-if="!isReadOnly"
          @click="toggleUploadForm"
          class="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
        >
          <i class="fa fa-upload mr-2"></i>立即提交
        </button>
      </div>
      
      <!-- 加载状态 -->
      <div v-if="isLoading" class="text-center py-16">
        <div class="inline-block animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
        <p class="mt-4 text-gray-600">加载第 {{ currentPage }} 页记录中...</p>
      </div>
      
      <!-- 历史记录列表（核心） -->
      <div v-if="historyList.length > 0 && !isLoading">
        <!-- 搜索+分页控制区 -->
        <div class="flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
          <h2 class="text-xl font-semibold text-gray-800">提交记录（第 {{ currentPage }} 页）</h2>
          
          <!-- 搜索框 -->
          <div class="relative w-full md:w-64">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="搜索文件名/描述/时间..."
              class="w-full pl-9 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-colors"
            />
            <i class="fa fa-search absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"></i>
          </div>
        </div>
        
        <!-- 分页控件（核心新增） -->
        <div class="flex items-center justify-between mb-6 gap-4 flex-wrap">
          <div class="text-sm text-gray-600">
            每页显示 {{ pageSize }} 条，本页 {{ historyList.length }} 条
          </div>
          <div class="flex items-center gap-2">
            <!-- 上一页 -->
            <button 
              @click="handlePrevPage"
              :disabled="isLoading || currentPage <= 1"
              class="px-3 py-1.5 border border-gray-300 rounded-md text-sm hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              <i class="fa fa-chevron-left mr-1 text-xs"></i>上一页
            </button>
            
            <!-- 页码显示 -->
            <span class="text-sm text-gray-700 px-2">
              第 {{ currentPage }} 页
            </span>
            
            <!-- 下一页 -->
            <button 
              @click="handleNextPage"
              :disabled="isLoading || !hasNextPage"
              class="px-3 py-1.5 border border-gray-300 rounded-md text-sm hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              下一页<i class="fa fa-chevron-right ml-1 text-xs"></i>
            </button>
            
            <!-- 页码跳转 -->
            <div class="flex items-center gap-1">
              <input
                v-model.number="targetPage"
                type="number"
                :min="1"
                class="w-16 px-2 py-1.5 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-1 focus:ring-blue-500"
                placeholder="页码"
              >
              <button 
                @click="handlePageJump"
                class="px-2 py-1.5 border border-gray-300 rounded-md text-sm hover:bg-gray-50 transition-colors"
              >
                跳转
              </button>
            </div>
          </div>
        </div>
        
        <!-- 记录列表 -->
        <div class="space-y-4">
          <div 
            v-for="(item, index) in filteredHistory" 
            :key="item.indices || index"  
            class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden hover:shadow-md transition-shadow"
          >
            <!-- 记录编号与状态 -->
            <div class="p-4 border-b border-gray-100 flex justify-between items-center">
              <h3 class="font-medium text-gray-900">提交 #{{ item.indices }}</h3>
              <span class="px-2 py-1 text-xs font-medium bg-green-100 text-green-800 rounded-full">
                已提交
              </span>
            </div>
            
            <!-- 核心信息：时间、文件大小 -->
            <div class="p-4 border-b border-gray-100">
              <div class="grid grid-cols-1 md:grid-cols-3 gap-2 text-sm text-gray-600">
                <div class="flex items-center">
                  <i class="fa fa-calendar-o mr-2 text-gray-400"></i>
                  <span><strong>提交时间:</strong> {{ convertUtcToCts(item.submittime) }}</span>
                </div>
                <div class="flex items-center">
                  <i class="fa fa-file-o mr-2 text-blue-400"></i>
                  <span><strong>头文件大小:</strong> {{ (item.headersize/1024).toFixed(2)}} KB</span>
                </div>
                <div class="flex items-center">
                  <i class="fa fa-file-code-o mr-2 text-green-400"></i>
                  <span><strong>源文件大小:</strong> {{ (item.sourcesize/1024).toFixed(2) }} KB</span>
                </div>
              </div>
            </div>
            
            <!-- 描述与运行状态 -->
            <div class="p-4">
              <div v-if="item.description" class="mb-3 text-sm text-gray-700">
                <strong>描述:</strong> {{ item.description }}
              </div>
              
              <div class="mb-4">
                <div class="flex flex-wrap items-center justify-between gap-2">
                  <div class="flex items-center text-sm text-gray-600 bg-gray-50 px-3 py-1.5 rounded w-full">
                    <i class="fa fa-file-code-o mr-2 text-blue-500 shrink-0"></i>
                    <span class="run-status-text" :class="{ 'expanded': expandedItems[index] || isCrash(item), 'crash-reason': isCrash(item) }">
                      <span class="status-content" :ref="(element) => setStatusElement(element, index)">
                        <span>运行状态: {{ getStatusLabel(item) }}</span>
                        <template v-if="isCrash(item) && getCrashReason(item)">
                          <span class="block mt-1">{{ getCrashReason(item) }}</span>
                        </template>
                      </span>
                    </span>
                  </div>
                  
                  <button 
                    v-if="!isCrash(item) && hasStatusOverflow(index)"
                    @click="toggleExpand(index)"
                    class="text-sm text-blue-600 hover:text-blue-800 transition-colors flex items-center"
                  >
                    <i :class="['fa', expandedItems[index] ? 'fa-chevron-up' : 'fa-chevron-down', 'mr-1']"></i>
                    {{ expandedItems[index] ? "收起" : "显示全部" }}
                  </button>
                </div>
              </div>

              <!-- 编译失败时，Data 是编译日志文件标识，需通过后端获取真实下载链接 -->
              <div v-if="isCompileFail(item)" class="mb-4">
                <button
                  @click="downloadFile(getCompileErrorLog(item))"
                  :disabled="!getCompileErrorLog(item)"
                  class="text-sm px-3 py-1.5 text-red-600 hover:text-red-800 hover:bg-red-50 rounded transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <i class="fa fa-download mr-1"></i>下载编译日志
                </button>
              </div>

              <!-- 运行成功/失败时，status.data 为录像文件标识 -->
              <div v-if="isGameResult(item)" class="mb-4">
                <button
                  @click="downloadFile(getVideoLink(item))"
                  :disabled="!getVideoLink(item)"
                  class="text-sm px-3 py-1.5 text-blue-600 hover:text-blue-800 hover:bg-blue-50 rounded transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <i class="fa fa-download mr-1"></i>下载录像
                </button>
              </div>

              <!-- 崩溃状态的新结构中，crash_log_file 不为空时提供日志下载 -->
              <div v-if="isCrash(item) && getCrashLogFile(item)" class="mb-4">
                <button
                  @click="downloadFile(getCrashLogFile(item))"
                  class="text-sm px-3 py-1.5 text-red-600 hover:text-red-800 hover:bg-red-50 rounded transition-colors"
                >
                  <i class="fa fa-download mr-1"></i>下载崩溃日志
                </button>
              </div>

              <!-- 运行中及运行结束时展示比赛结果数据 -->
              <div v-if="isGameStats(item)" class="mb-4 grid grid-cols-2 md:grid-cols-3 gap-2 text-sm text-gray-600 bg-gray-50 px-3 py-2 rounded">
                <div><strong>Food:</strong> {{ getStatusInfo(item).food }}</div>
                <div><strong>Wood:</strong> {{ getStatusInfo(item).wood }}</div>
                <div><strong>Gold:</strong> {{ getStatusInfo(item).gold }}</div>
                <div><strong>Stone:</strong> {{ getStatusInfo(item).stone }}</div>
                <div><strong>Frame:</strong> {{ getStatusInfo(item).frame }}</div>
                <div><strong>Score:</strong> {{ getStatusInfo(item).score }}</div>
              </div>
              
              <div class="mt-2 flex gap-2">
                <button 
                  @click="downloadFile(item.header)"
                  class="text-sm px-3 py-1.5 text-blue-600 hover:text-blue-800 hover:bg-blue-50 rounded transition-colors"
                  :disabled="!item.header" 
                >
                  <i class="fa fa-download mr-1"></i>下载头文件
                </button>
                <button 
                  @click="downloadFile(item.source)"
                  class="text-sm px-3 py-1.5 text-green-600 hover:text-green-800 hover:bg-green-50 rounded transition-colors"
                  :disabled="!item.source"  
                >
                  <i class="fa fa-download mr-1"></i>下载源文件
                </button>
                <button 
                  v-if="!isReadOnly"
                  @click="handleRun(item)" 
                  class="text-sm px-3 py-1.5 text-purple-600 hover:text-purple-800 hover:bg-purple-50 rounded transition-colors"
                  :disabled="!item.header || !item.source" 
                >
                  <span >重新运行</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- 页脚 -->
    <footer class="bg-white border-t border-gray-200 py-6">
      <div class="container mx-auto px-4 text-center text-gray-500 text-sm">
        <p>代码提交历史记录系统 &copy; {{ new Date().getFullYear() }}</p>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, computed, watch, nextTick, defineProps } from 'vue';
import config from '../config.js';
import { ElMessage } from 'element-plus';
import axios from 'axios';
import { downloadFile as requestDownload } from '../utils/download';

const props = defineProps({
  readOnly: {
    type: Boolean,
    default: false
  },
  historyTitle: {
    type: String,
    default: '代码提交历史'
  },
  backPath: {
    type: String,
    default: ''
  },
  historyUrl: {
    type: String,
    default: ''
  },
  requestParams: {
    type: Object,
    default: () => ({})
  }
});

// ======================== 核心新增：分页状态管理 ========================
const currentPage = ref(1); // 当前页码（默认第1页）
const pageSize = ref(config.HistoryRecordPerPage || 10); // 每页条数（优先从config取，默认10）
const hasNextPage = ref(true);
const targetPage = ref(1); // 跳转目标页码（绑定输入框）
const isReadOnly = computed(() => props.readOnly === true);
const historyTitle = computed(() => props.historyTitle || '代码提交历史');
const backPath = computed(() => props.backPath || '');

// ======================== 原有状态保留 ========================
const showUploadForm = ref(false);
const headerFile = ref(null);
const sourceFile = ref(null);
const description = ref('');
const isSubmitting = ref(false);
const isLoading = ref(true);
const historyList = ref([]);
const searchQuery = ref('');
const expandedItems = ref({});
const statusElements = ref({});
const statusOverflow = ref({});

watch(isReadOnly, readOnly => {
  if (readOnly) showUploadForm.value = false;
});

const setStatusElement = (element, index) => {
  if (element) {
    statusElements.value[index] = element;
  } else {
    delete statusElements.value[index];
  }
};

const updateStatusOverflow = async () => {
  await nextTick();
  const overflow = {};
  Object.entries(statusElements.value).forEach(([index, element]) => {
    if (!element) return;
    const styles = window.getComputedStyle(element);
    const lineHeight = parseFloat(styles.lineHeight);
    const fontSize = parseFloat(styles.fontSize);
    const oneLineHeight = Number.isFinite(lineHeight)
      ? lineHeight
      : (Number.isFinite(fontSize) ? fontSize * 1.5 : 21);
    // status-content 不受折叠容器 max-height 影响，用实际内容高度判断是否超过一行。
    overflow[index] = element.scrollHeight > oneLineHeight + 1;
  });
  statusOverflow.value = overflow;
};

const hasStatusOverflow = (index) => Boolean(statusOverflow.value[index]);

watch([historyList, searchQuery], updateStatusOverflow, { flush: 'post' });
onMounted(() => {
  window.addEventListener('resize', updateStatusOverflow);
  updateStatusOverflow();
});
onBeforeUnmount(() => {
  window.removeEventListener('resize', updateStatusOverflow);
});

// 后端 CodeRunStatusInfo.Status 的可读文案。Data 仅在特定状态下作为附加信息使用。
const statusLabels = {
  [config.Code_Status_Error]: '服务器异常',
  [config.Code_Status_Wait]: '等待中',
  [config.Code_Status_Compile]: '编译中',
  [config.Code_Status_Compile_Success]: '编译成功',
  [config.Code_Status_Compile_Fail]: '编译失败',
  [config.Code_Status_Running]: '运行中',
  [config.Code_Status_Success]: '运行成功',
  [config.Code_Status_Fail]: '运行失败',
  [config.Code_Status_Crash]: '游戏崩溃',
};

const getStatusInfo = (item) => {
  const rawStatus = item?.status;
  // 正常接口返回 CodeRunStatusInfo 对象；同时兼容被序列化成 JSON 字符串的情况。
  if (rawStatus && typeof rawStatus === 'object' && !Array.isArray(rawStatus)) {
    return rawStatus;
  }
  if (typeof rawStatus === 'string') {
    try {
      const parsedStatus = JSON.parse(rawStatus);
      if (parsedStatus && typeof parsedStatus === 'object' && !Array.isArray(parsedStatus)) {
        return parsedStatus;
      }
    } catch (error) {
      // 解析失败时直接把原始 status.data（此处为原始值）展示出来。
      return { status: rawStatus, data: rawStatus };
    }
  }
  return { status: rawStatus };
};

const getStatusCode = (item) => {
  const rawCode = getStatusInfo(item).status;
  if (typeof rawCode === 'number' && Number.isInteger(rawCode)) {
    return rawCode;
  }
  if (typeof rawCode === 'string' && /^-?\d+$/.test(rawCode.trim())) {
    return Number(rawCode.trim());
  }
  return NaN;
};

const getLegacyStatusText = (item) => {
  const statusInfo = getStatusInfo(item);
  if (statusInfo.data) {
    return String(statusInfo.data);
  }
  return '未运行';
};

const getStatusLabel = (item) => {
  const code = getStatusCode(item);
  return Object.prototype.hasOwnProperty.call(statusLabels, code)
    ? statusLabels[code]
    : getLegacyStatusText(item);
};

const getStatusData = (item) => {
  const data = getStatusInfo(item).data;
  return data === undefined || data === null ? '' : String(data);
};

const getCrashDetails = (item) => {
  let data = getStatusInfo(item).data;
  if (typeof data === 'string') {
    try {
      const parsed = JSON.parse(data);
      if (parsed && typeof parsed === 'object' && !Array.isArray(parsed)) data = parsed;
    } catch (error) {
      // 不是 JSON 结构时，按旧格式直接显示原始 data。
    }
  }
  if (data && typeof data === 'object' && !Array.isArray(data)) {
    return {
      reason: data.crash_reason === undefined || data.crash_reason === null
        ? ''
        : String(data.crash_reason),
      logLink: data.crash_log_file === undefined || data.crash_log_file === null
        ? ''
        : String(data.crash_log_file).trim()
    };
  }
  return { reason: data === undefined || data === null ? '' : String(data), logLink: '' };
};

const getCrashReason = (item) => getCrashDetails(item).reason;
const getCrashLogFile = (item) => getCrashDetails(item).logLink;
const getCompileErrorLog = (item) => {
  let data = getStatusInfo(item).data;
  if (typeof data === 'string') {
    try {
      const parsed = JSON.parse(data);
      if (parsed && typeof parsed === 'object' && !Array.isArray(parsed)) data = parsed;
    } catch (error) {
      return '';
    }
  }
  if (!data || typeof data !== 'object' || Array.isArray(data)) return '';
  return String(data.compile_error_log ?? '').trim();
};
const getVideoLink = (item) => {
  let data = getStatusInfo(item).data;
  if (typeof data === 'string') {
    try {
      const parsed = JSON.parse(data);
      if (parsed && typeof parsed === 'object' && !Array.isArray(parsed)) data = parsed;
    } catch (error) {
      return '';
    }
  }
  if (!data || typeof data !== 'object' || Array.isArray(data)) return '';
  return String(data.video_file ?? '').trim();
};

const isCompileFail = (item) => getStatusCode(item) === config.Code_Status_Compile_Fail;
const isGameResult = (item) => [config.Code_Status_Success, config.Code_Status_Fail].includes(getStatusCode(item));
const isGameStats = (item) => isGameResult(item) || getStatusCode(item) === config.Code_Status_Running;
const isCrash = (item) => getStatusCode(item) === config.Code_Status_Crash;

// ======================== 核心修改：同步页码输入框与当前页 ========================
watch(currentPage, (newPage) => {
  targetPage.value = newPage; // 切换页码时，输入框自动同步当前页
});

watch(() => props.requestParams, async (newParams, oldParams) => {
  if (newParams === oldParams) return;
  currentPage.value = 1;
  targetPage.value = 1;
  hasNextPage.value = true;
  historyList.value = [];
  const records = await GetHistory(1, pageSize.value);
  if (Array.isArray(records)) historyList.value = records;
});

// ======================== 时间格式转换（不变） ========================
const convertUtcToCts = (utcTime) => {
  if (!utcTime) return '未知时间';
  
  const utcDate = new Date(utcTime);
  if (isNaN(utcDate.getTime())) return '无效时间';
  
  const ctsTime = new Date(utcDate.getTime() + 8 * 60 * 60 * 1000);
  const year = ctsTime.getUTCFullYear();
  const month = String(ctsTime.getUTCMonth() + 1).padStart(2, '0');
  const day = String(ctsTime.getUTCDate()).padStart(2, '0');
  const hours = String(ctsTime.getUTCHours()).padStart(2, '0');
  const minutes = String(ctsTime.getUTCMinutes()).padStart(2, '0');
  const seconds = String(ctsTime.getUTCSeconds()).padStart(2, '0');
  
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
};

// ======================== 核心修改：分页获取历史记录 ========================
const GetHistory = async (page = 1, pageSize = config.HistoryRecordPerPage) => {
  isLoading.value = true;
  try {
    const beg = (page - 1) * pageSize;
    const end = beg + pageSize - 1;
    const historyUrl = props.historyUrl || `${config.history_url}`;
    const requestUrl = new URL(historyUrl);
    requestUrl.searchParams.set('range', `${beg}:${end}`);
    Object.entries(props.requestParams || {}).forEach(([key, value]) => {
      if (value !== undefined && value !== null && String(value).trim()) {
        requestUrl.searchParams.set(key, String(value));
      }
    });

    const response = await fetch(requestUrl.toString(), {
      method: 'GET',
      credentials: 'include',
    });

    let resData = await response.json();
    if (!response.ok || !resData.status) {
      throw new Error(resData.msg || `HTTP错误: ${response.status}`);
    }

    const payload = resData?.data;
    const rawRecords = Array.isArray(payload)
      ? payload
      : (Array.isArray(payload?.record)
        ? payload.record
        : (Array.isArray(payload?.records) ? payload.records : []));

    // 空页表示越界；保持当前页并停止下一页，和排行榜行为一致。
    if (rawRecords.length === 0 && page > 1) {
      if (page === currentPage.value + 1) hasNextPage.value = false;
      targetPage.value = currentPage.value;
      return null;
    }

    let records = [...rawRecords];
    // 按提交时间倒序（最新在前）
    records = records.sort((a, b) => {
      const timeA = new Date(a.submittime).getTime();
      const timeB = new Date(b.submittime).getTime();
      return timeB - timeA;
    });
    // 校验必要字段（移除无效记录）
    const requiredFields = ['header', 'source', 'description', 'submittime', 'headersize', 'sourcesize', 'status', 'indices','class'];
    records = records.filter(item => {
      const missingFields = requiredFields.filter(field => !(field in item));
      if (missingFields.length > 0) {
        console.warn(`过滤无效记录（缺少字段）:`, missingFields);
        return false;
      }
      return true;
    });
    currentPage.value = page;
    targetPage.value = page;
    hasNextPage.value = rawRecords.length >= pageSize;
    ElMessage.success(`成功加载 ${records.length} 条记录`);
    return records;

  } catch (error) {
    ElMessage.error(error.message||'分页获取历史记录失败');
    targetPage.value = currentPage.value;
    return null;
  } finally {
    isLoading.value = false; // 结束加载状态
  }
};

// ======================== 原有工具方法（不变） ========================
const getFileNameFromUrl = (url) => {
  if (!url) return '未知文件名';
  const urlWithoutParams = url.split('?')[0];
  const fileName = urlWithoutParams.split('/').pop();
  return fileName.includes('.') ? fileName : '未命名文件';
};

const downloadFile = async (fileUrl) => {
  try {
    const { fileName } = await requestDownload(fileUrl);
    ElMessage.success(`开始下载：${fileName}`);
  } catch (error) {
    if (error?.message === '文件链接无效') {
      ElMessage.warning(error.message);
    } else {
      console.error('下载失败：', error);
      ElMessage.error(error?.message || '下载失败，请重试');
    }
  }
};

const toggleExpand = (index) => {
  expandedItems.value[index] = !expandedItems.value[index];
  updateStatusOverflow();
};

const handleRun = async (item) => {
  if (isReadOnly.value) return;
  try {
    historyList.value = [...historyList.value];
    const resp = await fetch(`${config.codeRun_url}/coderun`, {
      method: 'POST',
      credentials: 'include',
      body: JSON.stringify({ 
        indices:item.indices,
        class:config.Code_ReRunSubmit
      })
    });

    let data = await resp.json();
    if (!resp.ok || !data.status) {
      throw new Error(data.msg || "运行失败");
    }

    ElMessage.success("运行成功");
    // 重新获取当前页数据（保证状态同步）
    const newRecords=await GetHistory(currentPage.value, pageSize.value);
    historyList.value=newRecords
  } catch (error) {
    ElMessage.error(error.message || '运行出错');

  } finally {
    historyList.value = [...historyList.value];
    // 重置展开状态
    for (const key in expandedItems.value) {
      if (historyList.value[Number(key)] === item) {
        expandedItems.value[key] = false;
        break;
      }
    }
  }
};

const toggleUploadForm = () => {
  if (isReadOnly.value) return;
  showUploadForm.value = !showUploadForm.value;
  if (showUploadForm.value) {
    headerFile.value = null;
    sourceFile.value = null;
    description.value = '';
    document.querySelectorAll('input[type="file"]').forEach(input => input.value = '');
  }
};

const handleFileUpload = (event, type) => {
  const file = event.target.files[0];
  if (!file) return;

  if (type === 'header' && !file.name.endsWith('.h')) {
    ElMessage.warning('请上传 .h 格式的头文件');
    event.target.value = '';
    return;
  }
  if (type === 'source' && !file.name.endsWith('.cpp')) {
    ElMessage.warning('请上传 .cpp 格式的源文件');
    event.target.value = '';
    return;
  }

  if (type === 'header') headerFile.value = file;
  if (type === 'source') sourceFile.value = file;
};

// ======================== 核心修改：提交后刷新当前页 ========================
const submitFiles = async () => {
  if (isReadOnly.value) return;
  if (!headerFile.value || !sourceFile.value) return;
  
  isSubmitting.value = true;
  try {
    //获取三个上传链接
    const uploadURLGet=await fetch(`${config.codeSubmit_url}/codecommonsubmit`, {
      method: 'GET',
      credentials: 'include'
    });
    const resData = await uploadURLGet.json();
    if (!uploadURLGet.ok || !resData.status) {
      throw new Error(resData.msg || '文件提交失败');
    }
    //获取url
    const key=resData.data.key;
    const urls=resData.data.urls;
    const headerURL=urls[0]
    const sourceURL=urls[1]
    const descURL=urls[2]
    //上传三个文件
    const p1=axios.put(headerURL,headerFile.value,{
      headers:{
        'Content-Type':headerFile.value.type
      }
    })
    const p2=axios.put(sourceURL,sourceFile.value,{
      headers:{
        'Content-Type':sourceFile.value.type
      }
    })
    const p3=axios.put(descURL,description.value ||"原神启动!" ,{
      headers:{
        'Content-Type':'text/plain'
      }
    })
    const [r0,r1,r2]=await Promise.all([p1,p2,p3])
    if (!r0.status || !r1.status || !r2.status){
       throw new Error('文件提交失败，请重试');
    }
    //告诉服务器上传成功了
    const tellServer=await fetch(`${config.codeSubmit_url}/codecommonsubmitACK`,{
      method:'POST',
      credentials:'include',
      body:JSON.stringify({
        key:key
      })
    })
    const dt=await tellServer.json()
    if (!tellServer.ok || !dt.status){
       throw new Error ('文件提交失败，请重试' | dt.msg);
    }
    //
    ElMessage.success('文件提交成功！');
    //运行代码
     const coderun=await fetch(`${config.codeRun_url}/coderun`, {
      method: 'POST',
      credentials: 'include',
      body:JSON.stringify({
        indices:dt.data.indices,
        class:config.Code_CommonSubmit
      })
    });
    const coderunData=await coderun.json()
    if (!coderun.ok || !coderunData.status){
       throw new Error ('运行失败，请重试' | coderunData.msg);
    }
    ElMessage.success('运行成功！');
    // 提交后重新获取当前页数据（保证新记录显示）
    const newHistory = await GetHistory(currentPage.value, pageSize.value);
    historyList.value = newHistory;
    expandedItems.value = {};
    toggleUploadForm();

  } catch (error) {
    ElMessage.error(error.message || '文件提交失败，请重试');
  } finally {
    isSubmitting.value = false;
  }
};

// ======================== 搜索过滤（基于当前页数据） ========================
const filteredHistory = computed(() => {
  if (!searchQuery.value.trim()) return historyList.value;
  
  const query = searchQuery.value.toLowerCase();
  return historyList.value.filter(item => {
    const headerName = getFileNameFromUrl(item.header).toLowerCase();
    const sourceName = getFileNameFromUrl(item.source).toLowerCase();
    const descMatch = item.description ? item.description.toLowerCase().includes(query) : false;
    const dateMatch = item.submittime.toLowerCase().includes(query);
    const statusMatch = `${getStatusLabel(item)} ${getStatusData(item)}`.toLowerCase().includes(query);
    
    return headerName.includes(query) || sourceName.includes(query) || descMatch || dateMatch || statusMatch;
  });
});

// ======================== 核心新增：分页控制方法 ========================
// 上一页
const handlePrevPage = async () => {
  if (currentPage.value > 1) {
    const prevPage = currentPage.value - 1;
    const records = await GetHistory(prevPage, pageSize.value);
    if (Array.isArray(records)) {
      historyList.value = records;
      expandedItems.value = {}; // 切换页重置展开状态
    }
  }
};

// 下一页
const handleNextPage = async () => {
  if (hasNextPage.value) {
    const nextPage = currentPage.value + 1;
    const records = await GetHistory(nextPage, pageSize.value);
    if (Array.isArray(records)) {
      historyList.value = records;
      expandedItems.value = {}; // 切换页重置展开状态
    }
  }
};

// 页码跳转
const handlePageJump = async () => {
  // 总页数未知，越界由接口返回空数据来判断。
  const target = Number(targetPage.value);
  if (!Number.isInteger(target) || target < 1 || target === currentPage.value) {
    ElMessage.warning('请输入合法的页码');
    return;
  }
  const records = await GetHistory(target, pageSize.value);
  if (Array.isArray(records)) {
    historyList.value = records;
    expandedItems.value = {}; // 切换页重置展开状态
  }
};

// ======================== 组件挂载：加载第1页数据 ========================
onMounted(async () => {
  try {
    const initialHistory = await GetHistory(1, pageSize.value);
    historyList.value = Array.isArray(initialHistory) ? initialHistory : [];
  } catch (error) {
    ElMessage.error(error.message || '加载历史记录失败 ');
    historyList.value = [];
  }
});
</script>

<script>
export default {
  name: 'HistoryRecords'
}
</script>

<style>
/* 自定义动画：表单显示过渡 */
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-10px); }
  to { opacity: 1; transform: translateY(0); }
}
.animate-fadeIn {
  animation: fadeIn 0.3s ease-out forwards;
}

/* 优化文件上传区域hover效果 */
label.border-dashed:hover {
  border-color: #3b82f6;
  background-color: #f0f9ff;
}

/* 适配小屏幕的记录布局 */
@media (max-width: 768px) {
  .grid-cols-3 {
    grid-template-columns: 1fr;
    gap: 1rem !important;
  }
  .flex-wrap {
    flex-direction: column;
    align-items: flex-start !important;
  }
}

/* 运行状态文本样式：默认一行，展开后多行 */
.run-status-text {
  flex: 1;
  max-height: 1.5em; /* 一行高度（根据字体调整） */
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: pre-wrap;
  transition: max-height 0.3s ease, white-space 0.3s ease;
}
.status-content {
  display: block;
}
.run-status-text.expanded {
  max-height: 10em; /* 最多显示10行，可根据需求调整 */
  white-space: pre-wrap;
  text-overflow: unset;
}
.run-status-text.crash-reason {
  max-height: none;
  overflow: visible;
  white-space: pre-wrap;
  overflow-wrap: anywhere;
}
</style>
