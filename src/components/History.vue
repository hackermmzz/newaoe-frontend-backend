<template>
  <div class="min-h-screen bg-gray-50 flex flex-col">
    <!-- 顶部固定导航栏 -->
    <header class="bg-white shadow-md sticky top-0 z-50 transition-all duration-300">
      <div class="container mx-auto px-4 py-4 flex justify-between items-center">
        <h1 class="text-2xl font-bold text-gray-800">代码提交历史</h1>
        <button 
          @click="toggleUploadForm"
          class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-lg transition-all duration-200 transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
        >
          <i class="fa fa-upload mr-2"></i>提交新文件
        </button>
      </div>
      
      <!-- 上传表单区域 - 条件显示 -->
      <div 
        v-if="showUploadForm"
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
      <div v-if="historyList.length === 0 && !isLoading && totalRecord === 0" class="text-center py-16">
        <div class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-gray-100 mb-4">
          <i class="fa fa-history text-2xl text-gray-400"></i>
        </div>
        <h3 class="text-xl font-medium text-gray-700 mb-2">暂无历史记录</h3>
        <p class="text-gray-500 mb-6">点击上方"提交新文件"按钮开始上传你的第一个文件</p>
        <button 
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
      <div v-if="totalRecord > 0 && !isLoading">
        <!-- 搜索+分页控制区 -->
        <div class="flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
          <h2 class="text-xl font-semibold text-gray-800">提交记录 ({{ totalRecord }} 条)</h2>
          
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
            每页显示 {{ pageSize }} 条，共 {{ totalPages }} 页
          </div>
          <div class="flex items-center gap-2">
            <!-- 上一页 -->
            <button 
              @click="handlePrevPage"
              :disabled="currentPage === 1"
              class="px-3 py-1.5 border border-gray-300 rounded-md text-sm hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              <i class="fa fa-chevron-left mr-1 text-xs"></i>上一页
            </button>
            
            <!-- 页码显示 -->
            <span class="text-sm text-gray-700 px-2">
              第 {{ currentPage }} / {{ totalPages }} 页
            </span>
            
            <!-- 下一页 -->
            <button 
              @click="handleNextPage"
              :disabled="currentPage === totalPages"
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
                :max="totalPages"
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
              <h3 class="font-medium text-gray-900">提交 #{{ (currentPage - 1) * pageSize + index + 1 }}</h3>
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
                    <span class="run-status-text" :class="{ 'expanded': expandedItems[index] }">
                      运行状态: {{ item.status.data || "未运行" }}
                    </span>
                  </div>
                  
                  <button 
                    @click="toggleExpand(index)"
                    class="text-sm text-blue-600 hover:text-blue-800 transition-colors flex items-center"
                  >
                    <i :class="['fa', expandedItems[index] ? 'fa-chevron-up' : 'fa-chevron-down', 'mr-1']"></i>
                    {{ expandedItems[index] ? "收起" : "显示全部" }}
                  </button>
                </div>
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
import { ref, onMounted, computed, watch } from 'vue';
import config from '../config.js';
import { ElMessage } from 'element-plus';
import axios from 'axios';

// ======================== 核心新增：分页状态管理 ========================
const currentPage = ref(1); // 当前页码（默认第1页）
const totalRecord = ref(0); // 总记录数（接口返回的 totalrecord）
const pageSize = ref(config.HistoryRecordPerPage || 10); // 每页条数（优先从config取，默认10）
const totalPages = ref(0); // 总页数（计算得出：totalRecord / pageSize 向上取整）
const targetPage = ref(1); // 跳转目标页码（绑定输入框）

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

// ======================== 核心修改：同步页码输入框与当前页 ========================
watch(currentPage, (newPage) => {
  targetPage.value = newPage; // 切换页码时，输入框自动同步当前页
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
    // 1. 拼接分页参数到请求URL（page：当前页，pageSize：每页条数）
    let beg=(page-1)*pageSize
    let end=beg+pageSize-1
    const requestUrl = new URL(`${config.base_url}/home/gethistory`);
    requestUrl.searchParams.append('range', `${beg}:${end}`);

    const response = await fetch(requestUrl.toString(), {
      method: 'GET',
      credentials: 'include',
    });

    let resData = await response.json();
    if (!response.ok || !resData.status) {
      throw new Error(resData.msg || `HTTP错误: ${response.status}`);
    }

    // 2. 从接口获取总记录数（关键：resData.totalrecord）
    totalRecord.value = resData.data.totalrecord || 0;
    // 3. 计算总页数（向上取整，避免小数页）
    totalPages.value = Math.ceil(totalRecord.value / pageSize);
    // 4. 同步当前页码
    currentPage.value = page;

    // 5. 处理当前页数据（保留原有排序和字段校验）
    let records = Array.isArray(resData.data.record) ? resData.data.record : [];
    // 按提交时间倒序（最新在前）
    records = records.sort((a, b) => {
      const timeA = new Date(a.submittime).getTime();
      const timeB = new Date(b.submittime).getTime();
      return timeB - timeA;
    });
    // 校验必要字段（移除无效记录）
    const requiredFields = ['header', 'source', 'description', 'submittime', 'headersize', 'sourcesize', 'status', 'indices'];
    records = records.filter(item => {
      const missingFields = requiredFields.filter(field => !(field in item));
      if (missingFields.length > 0) {
        console.warn(`过滤无效记录（缺少字段）:`, missingFields);
        return false;
      }
      return true;
    });
    //
    ElMessage.success(`成功加载 ${records.length} 条记录`);
    return records;

  } catch (error) {
    ElMessage.error(error.message||'分页获取历史记录失败');
    totalRecord.value = 0;
    totalPages.value = 0;
    return [];
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
  if (!fileUrl) {
    ElMessage.warning('文件链接无效，无法下载');
    return;
  }

  const apiUrl = `${config.download_url}/${fileUrl}?download=true`;

  try {
    // 1. 获取真实的下载链接
    const UrlGetResp = await fetch(apiUrl, {
      method: 'GET',
      credentials: 'include'
    });

    if (!UrlGetResp.ok) {
      ElMessage.error("获取下载链接失败");
      return;
    }

    const data = await UrlGetResp.json();
    if (!data?.data?.url) {
      ElMessage.error("下载链接无效：" + (data.msg || "未知错误"));
      return;
    }

    const realUrl = data.data.url;
    const fileName = getFileNameFromUrl(realUrl);

    // ==========================================
    // ✅ 核心：触发浏览器自带下载 + 右上角进度条
    // ==========================================
    const link = document.createElement('a');
    link.href = realUrl;
    link.download = fileName;  // 强制下载，不预览
    link.target = '_self';     // 不打开新窗口
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);

    ElMessage.success(`开始下载：${fileName}`);

  } catch (error) {
    console.error('下载失败：', error);
    ElMessage.error('下载失败，请重试');
  }
};

const toggleExpand = (index) => {
  expandedItems.value[index] = !expandedItems.value[index];
};

const handleRun = async (item) => {
  try {
    historyList.value = [...historyList.value];
    
    const resp = await fetch(`${config.code_url}/CodeReRun`, {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ 
        source: item.source,
        header: item.header,
        description: item.description
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
  if (!headerFile.value || !sourceFile.value) return;
  
  isSubmitting.value = true;
  try {
    //获取三个上传链接
    const uploadURLGet=await fetch(`${config.upload_url}/code`, {
      method: 'POST',
      credentials: 'include'
    });
    const resData = await uploadURLGet.json();
    if (!uploadURLGet.ok || !resData.status) {
      throw new Error(resData.msg || '文件提交失败');
    }
    //获取url
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
    const p3=axios.put(descURL,description.value,{
      headers:{
        'Content-Type':'text/plain'
      }
    })
    const [r0,r1,r2]=await Promise.all([p1,p2,p3])
    if (!r0.status || !r1.status || !r2.status){
       throw new Error('文件提交失败，请重试');
    }
    //告诉服务器上传成功了
    const tellServer=await fetch(`${config.uploadConfirm_url}/code`,{
      method:'POST',
      credentials:'include',
      body:JSON.stringify({
        urls:urls
      })
    })
    const dt=await tellServer.json()
    if (!tellServer.ok || !dt.status){
       throw new Error ('文件提交失败，请重试' | dt.msg);
    }
    //
    ElMessage.success('文件提交成功！');
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
    const statusMatch = (item.status.data || '').toLowerCase().includes(query);
    
    return headerName.includes(query) || sourceName.includes(query) || descMatch || dateMatch || statusMatch;
  });
});

// ======================== 核心新增：分页控制方法 ========================
// 上一页
const handlePrevPage = async () => {
  if (currentPage.value > 1) {
    const prevPage = currentPage.value - 1;
    const records = await GetHistory(prevPage, pageSize.value);
    historyList.value = records;
    expandedItems.value = {}; // 切换页重置展开状态
  }
};

// 下一页
const handleNextPage = async () => {
  if (currentPage.value < totalPages.value) {
    const nextPage = currentPage.value + 1;
    const records = await GetHistory(nextPage, pageSize.value);
    historyList.value = records;
    expandedItems.value = {}; // 切换页重置展开状态
  }
};

// 页码跳转
const handlePageJump = async () => {
  // 校验目标页合法性（必须是数字、在1~总页数之间、不等于当前页）
  const target = Number(targetPage.value);
  if (isNaN(target) || target < 1 || target > totalPages.value || target === currentPage.value) {
    ElMessage.warning('请输入合法的页码');
    return;
  }
  const records = await GetHistory(target, pageSize.value);
  historyList.value = records;
  expandedItems.value = {}; // 切换页重置展开状态
};

// ======================== 组件挂载：加载第1页数据 ========================
onMounted(async () => {
  try {
    const initialHistory = await GetHistory(1, pageSize.value);
    historyList.value = initialHistory;
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
.run-status-text.expanded {
  max-height: 10em; /* 最多显示10行，可根据需求调整 */
  white-space: pre-wrap;
  text-overflow: unset;
}
</style>