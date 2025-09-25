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
      <div v-if="historyList.length === 0 && !isLoading" class="text-center py-16">
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
        <p class="mt-4 text-gray-600">加载历史记录中...</p>
      </div>
      
      <!-- 历史记录列表（核心） -->
      <div v-if="historyList.length > 0 && !isLoading">
        <div class="flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
          <h2 class="text-xl font-semibold text-gray-800">提交记录 ({{ historyList.length }})</h2>
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
        
        <div class="space-y-4">
          <!-- 循环渲染记录 -->
          <div 
            v-for="(item, index) in filteredHistory" 
            :key="index"  
            class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden hover:shadow-md transition-shadow"
          >
            <!-- 记录编号与状态 -->
            <div class="p-4 border-b border-gray-100 flex justify-between items-center">
              <h3 class="font-medium text-gray-900">提交 #{{ index + 1 }}</h3>
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
            
            <!-- 描述与运行状态（核心修改：运行状态默认一行，点击展开全部） -->
            <div class="p-4">
              <!-- 描述（可选字段） -->
              <div v-if="item.description" class="mb-3 text-sm text-gray-700">
                <strong>描述:</strong> {{ item.description }}
              </div>
              
              <!-- 运行状态区域：默认一行显示，点击展开全部 -->
              <div class="mb-4">
                <!-- 顶部：运行状态标签 + 展开/收起按钮 -->
                <div class="flex flex-wrap items-center justify-between gap-2">
                  <!-- 运行状态标签：默认一行，超出省略 -->
                  <div class="flex items-center text-sm text-gray-600 bg-gray-50 px-3 py-1.5 rounded w-full">
                    <i class="fa fa-file-code-o mr-2 text-blue-500 shrink-0"></i>
                    <span class="run-status-text" :class="{ 'expanded': expandedItems[index] }">
                      运行状态: {{ item.status || "未运行" }}
                    </span>
                  </div>
                  
                  <!-- 展开/收起按钮 -->
                  <button 
                    @click="toggleExpand(index)"
                    class="text-sm text-blue-600 hover:text-blue-800 transition-colors flex items-center"
                  >
                    <i :class="['fa', expandedItems[index] ? 'fa-chevron-up' : 'fa-chevron-down', 'mr-1']"></i>
                    {{ expandedItems[index] ? "收起" : "显示全部" }}
                  </button>
                </div>
              </div>
              
              <!-- 下载按钮 -->
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
import { ref, onMounted, computed } from 'vue';
import config from '../config.js';
import { ElMessage } from 'element-plus';

// 状态管理
const showUploadForm = ref(false);
const headerFile = ref(null);
const sourceFile = ref(null);
const description = ref('');
const isSubmitting = ref(false);
const isLoading = ref(true);
const historyList = ref([]);
const searchQuery = ref('');
// 控制每条记录的运行状态展开/收起（核心状态）
const expandedItems = ref({});

/**
 * 时间格式转换：UTC转CTS（UTC+8）
 */
const convertUtcToCts = (utcTime) => {
  if (!utcTime) return '未知时间';
  
  const utcDate = new Date(utcTime);
  if (isNaN(utcDate.getTime())) return '无效时间';
  
  // 计算CTS时间（UTC+8）
  const ctsTime = new Date(utcDate.getTime() + 8 * 60 * 60 * 1000);
  
  // 格式化输出：YYYY-MM-DD HH:MM:SS
  const year = ctsTime.getUTCFullYear();
  const month = String(ctsTime.getUTCMonth() + 1).padStart(2, '0');
  const day = String(ctsTime.getUTCDate()).padStart(2, '0');
  const hours = String(ctsTime.getUTCHours()).padStart(2, '0');
  const minutes = String(ctsTime.getUTCMinutes()).padStart(2, '0');
  const seconds = String(ctsTime.getUTCSeconds()).padStart(2, '0');
  
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
};

/**
 * 获取历史记录：仅依赖接口，无模拟数据
 */
const GetHistory = async () => {
  try {
    const response = await fetch(`${config.base_url}/home/gethistory`, {
      method: 'GET',
      credentials: 'include',
    });

    let resData = await response.json();
    console.log(resData)
    if (!response.ok || !resData.status) {
      throw new Error(resData.msg || `HTTP错误: ${response.status}`);
    }

    // 处理非数组返回
    if (!Array.isArray(resData.data)) {
      resData.data = [];
    }

    // 按提交时间倒序排序（最新在前）
    resData.data = resData.data.sort((a, b) => {
      const timeA = new Date(a.submittime).getTime();
      const timeB = new Date(b.submittime).getTime();
      return timeB - timeA;
    });

    // 校验必要字段（仅保留原字段，无额外新增）
    const requiredFields = ['header', 'source', 'description', 'submittime', 'headersize', 'sourcesize', 'status', 'indices'];
    resData.data.forEach((item, index) => {
      const missingFields = requiredFields.filter(field => !(field in item));
      if (missingFields.length > 0) {
        throw new Error(`第${index+1}条记录缺少字段：${missingFields.join(', ')}`);
      }
    });

    return resData.data;

  } catch (error) {
    console.error('获取历史记录失败:', error);
    return []; // 接口失败时返回空数组（无模拟数据）
  }
};

/**
 * 从文件链接提取文件名
 */
const getFileNameFromUrl = (url) => {
  if (!url) return '未知文件名';
  
  // 处理带参数的链接（如 ?name=xxx.cpp）
  const urlWithoutParams = url.split('?')[0];
  // 从最后一个 '/' 后截取文件名
  const fileName = urlWithoutParams.split('/').pop();
  
  return fileName.includes('.') ? fileName : '未命名文件';
};

/**
 * 文件下载功能
 */
const downloadFile = (fileUrl) => {
  if (!fileUrl) {
    ElMessage.warning('文件链接无效，无法下载');
    return;
  }
  fileUrl = `${config.download_url}/${fileUrl}`;
  try {
    const fileName = getFileNameFromUrl(fileUrl);
    const a = document.createElement('a');
    a.href = fileUrl;
    a.download = fileName;
    a.style.display = 'none';
    
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);

    //
    let filename=""
    for(let i=fileName.length-1;i>=0;i-=1){
      if(fileName[i]!='\\' &&fileName[i]!='/')
      {
        filename=fileName[i]+filename
      }
      else{
        break
      }
    }
    ElMessage.success(`开始下载：${filename}`);

  } catch (error) {
    ElMessage.error(`下载失败：${error.message}`);
    console.error('下载出错:', error);
  }
};

/**
 * 切换运行状态的展开/收起
 * @param {number} index - 记录索引
 */
const toggleExpand = (index) => {
  // 切换对应索引的展开状态（默认未展开）
  expandedItems.value[index] = !expandedItems.value[index];
};

/**
 * 运行代码：仅更新运行状态，无额外字段
 */
const handleRun = async (item) => {
  try {
    // 运行前更新状态（优化用户体验）
    item.status = '正在运行';
    // 触发视图更新
    historyList.value = [...historyList.value];
    
    const resp = await fetch(`${config.admin_url}/coderun`, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ 
        source: item.source,
        header:item.header,
        description:item.description
      })
    });

    let data = await resp.json();
    if (!resp.ok || !data.status) {
      throw new Error(data.msg || "运行失败");
    }

    ElMessage.success("运行成功");
    // 运行成功后更新状态（仅保留原字段）
    item.status = '正在运行';
    GetHistory()
  } catch (error) {
    ElMessage.error(error.message || '运行出错');
    // 运行失败更新状态
    item.status = '运行失败';

  } finally {
    // 最终触发视图更新
    historyList.value = [...historyList.value];
    // 收起运行状态详情（避免失败后仍展开）
    for (const key in expandedItems.value) {
      if (historyList.value[Number(key)] === item) {
        expandedItems.value[key] = false;
        break;
      }
    }
  }
};

/**
 * 切换上传表单显示/隐藏
 */
const toggleUploadForm = () => {
  showUploadForm.value = !showUploadForm.value;
  if (showUploadForm.value) {
    // 重置表单状态
    headerFile.value = null;
    sourceFile.value = null;
    description.value = '';
    document.querySelectorAll('input[type="file"]').forEach(input => input.value = '');
  }
};

/**
 * 处理文件上传（格式校验）
 */
const handleFileUpload = (event, type) => {
  const file = event.target.files[0];
  if (!file) return;

  // 格式校验
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

  // 赋值对应文件
  if (type === 'header') headerFile.value = file;
  if (type === 'source') sourceFile.value = file;
};

/**
 * 生成描述文件（提交时附带）
 */
const createDescFile = () => {
  const timestamp = new Date().getTime();
  const fileName = `description_${timestamp}.txt`;
  const descBlob = new Blob([description.value], { type: 'text/plain' });
  return new File([descBlob], fileName, { type: 'text/plain' });
};

/**
 * 提交文件：提交后刷新历史记录
 */
const submitFiles = async () => {
  if (!headerFile.value || !sourceFile.value) return;
  
  isSubmitting.value = true;
  try {
    const formData = new FormData();
    formData.append('file', headerFile.value);
    formData.append('file', sourceFile.value);
    formData.append('file', createDescFile());
    
    const response = await fetch(`${config.upload_url}/code`, {
      method: 'POST',
      credentials: 'include',
      body: formData,
    });

    const resData = await response.json();
    if (!response.ok || !resData.status) {
      throw new Error(resData.msg || '文件提交失败');
    }

    ElMessage.success('文件提交成功！');
    // 重新获取历史记录
    const newHistory = await GetHistory();
    historyList.value = newHistory;
    // 重置展开状态
    expandedItems.value = {};
    // 关闭上传表单
    toggleUploadForm();

  } catch (error) {
    ElMessage.error(error.message || '文件提交失败，请重试');
  } finally {
    isSubmitting.value = false;
  }
};

/**
 * 搜索过滤：支持文件名、描述、时间、运行状态
 */
const filteredHistory = computed(() => {
  if (!searchQuery.value.trim()) return historyList.value;
  
  const query = searchQuery.value.toLowerCase();
  return historyList.value.filter(item => {
    const headerName = getFileNameFromUrl(item.header).toLowerCase();
    const sourceName = getFileNameFromUrl(item.source).toLowerCase();
    const descMatch = item.description ? item.description.toLowerCase().includes(query) : false;
    const dateMatch = item.submittime.toLowerCase().includes(query);
    const statusMatch = (item.status || '').toLowerCase().includes(query);
    
    return headerName.includes(query) || sourceName.includes(query) || descMatch || dateMatch || statusMatch;
  });
});

/**
 * 组件挂载：加载历史记录
 */
onMounted(async () => {
  isLoading.value = true;
  try {
    const initialHistory = await GetHistory();
    historyList.value = initialHistory;
    if (historyList.value.length > 0) {
      ElMessage.success(`成功加载 ${historyList.value.length} 条历史记录`);
    }
  } catch (error) {
    ElMessage.error('加载历史记录失败');
    historyList.value = [];
  } finally {
    isLoading.value = false;
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
  max-height: none; /* 足够显示多行（可根据需求调整） */
  white-space: pre-wrap;
  text-overflow: unset;
}

</style>