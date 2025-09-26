<template>
  <div class="min-h-screen bg-gray-50 flex flex-col">
    

    <!-- 主要内容 -->
    <main class="flex-grow max-w-7xl w-full mx-auto py-6 sm:px-6 lg:px-8">
      <div class="bg-white shadow-md rounded-lg overflow-hidden">
        <div class="p-6">
          <form @submit.prevent="handleSubmit" class="space-y-6">
            <!-- 头文件上传（仅支持.h格式） -->
            <div>
              <label for="headerFile" class="block text-sm font-medium text-gray-700 mb-1">
                头文件
              </label>
              <div class="mt-1 flex justify-center px-6 pt-5 pb-6 border-2 border-gray-300 border-dashed rounded-md">
                <div class="space-y-1 text-center">
                  <svg class="mx-auto h-12 w-12 text-gray-400" stroke="currentColor" fill="none" viewBox="0 0 48 48" aria-hidden="true">
                    <path d="M28 8H12a4 4 0 00-4 4v20m32-12v8m0 0v8a4 4 0 01-4 4H12a4 4 0 01-4-4v-4m32-4l-3.172-3.172a4 4 0 00-5.656 0L28 28M8 32l9.172-9.172a4 4 0 015.656 0L28 28m0 0l4 4m4-24h8m-4-4v8m-12 4h.02" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
                  </svg>
                  <div class="flex text-sm text-gray-600">
                    <label for="headerFile" class="relative cursor-pointer bg-white rounded-md font-medium text-indigo-600 hover:text-indigo-500 focus-within:outline-none">
                      <span>上传文件</span>
                      <input id="headerFile" name="headerFile" type="file" class="sr-only" @change="handleHeaderFileChange" accept=".h">
                    </label>
                    <p class="pl-1">或拖放文件</p>
                  </div>
                  <!-- 提示文字修改：仅支持.h格式 -->
                  <p class="text-xs text-gray-500">
                    仅支持 .h 格式（其他格式将无法上传）
                  </p>
                  <p v-if="headerFileName" class="text-sm text-green-600 mt-2">
                    已选择: {{ headerFileName }}
                  </p>
                </div>
              </div>
            </div>

            <!-- 源文件上传（仅支持.cpp格式） -->
            <div>
              <label for="sourceFile" class="block text-sm font-medium text-gray-700 mb-1">
                源文件
              </label>
              <div class="mt-1 flex justify-center px-6 pt-5 pb-6 border-2 border-gray-300 border-dashed rounded-md">
                <div class="space-y-1 text-center">
                  <svg class="mx-auto h-12 w-12 text-gray-400" stroke="currentColor" fill="none" viewBox="0 0 48 48" aria-hidden="true">
                    <path d="M28 8H12a4 4 0 00-4 4v20m32-12v8m0 0v8a4 4 0 01-4 4H12a4 4 0 01-4-4v-4m32-4l-3.172-3.172a4 4 0 00-5.656 0L28 28M8 32l9.172-9.172a4 4 0 015.656 0L28 28m0 0l4 4m4-24h8m-4-4v8m-12 4h.02" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
                  </svg>
                  <div class="flex text-sm text-gray-600">
                    <label for="sourceFile" class="relative cursor-pointer bg-white rounded-md font-medium text-indigo-600 hover:text-indigo-500 focus-within:outline-none">
                      <span>上传文件</span>
                      <input id="sourceFile" name="sourceFile" type="file" class="sr-only" @change="handleSourceFileChange" accept=".cpp">
                    </label>
                    <p class="pl-1">或拖放文件</p>
                  </div>
                  <!-- 提示文字修改：仅支持.cpp格式 -->
                  <p class="text-xs text-gray-500">
                    仅支持 .cpp 格式（其他格式将无法上传）
                  </p>
                  <p v-if="sourceFileName" class="text-sm text-green-600 mt-2">
                    已选择: {{ sourceFileName }}
                  </p>
                </div>
              </div>
            </div>

            <!-- 授课老师选择（数据从resp.data获取，仅用name字段） -->
            <div>
              <label for="teacher" class="block text-sm font-medium text-gray-700 mb-1">
                授课老师
              </label>
              <select id="teacher" 
                      v-model="selectedTeacher"
                      class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md"
                      required>
                <option value="">请选择授课老师</option>
                <!-- 循环resp.data中的老师列表，仅使用name字段 -->
                <option v-for="teacher in teachers" :key="teacher.name" :value="teacher.name">
                  {{ teacher.name }}
                </option>
              </select>
            </div>

            <!-- 提交按钮 -->
            <div class="flex justify-end">
              <button type="submit" 
                      class="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
                提交考核
              </button>
            </div>
          </form>
        </div>
      </div>
    </main>

    <!-- 页脚 -->
    <footer class="bg-white border-t border-gray-200 py-4">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <p class="text-center text-sm text-gray-500">
          &copy; {{ new Date().getFullYear() }} 学生考核系统
        </p>
      </div>
    </footer>

    <!-- 确认提示框 -->
    <div v-if="showConfirmation" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-xl max-w-md w-full p-6">
        <h3 class="text-lg font-medium text-gray-900 mb-4">确认提交</h3>
        <p class="text-gray-700 mb-6">
          在截至日期前你可以提交多次，请确认所填信息正确无误，否则后果自负
        </p>
        <div class="flex justify-end space-x-3">
          <button @click="showConfirmation = false" 
                  class="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
            取消
          </button>
          <button @click="confirmSubmission" 
                  class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
            确认提交
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ElMessage } from 'element-plus';
import config from '../config.js';
export default {
  name: 'AssessmentSubmission',
  data() {
    return {
      // 头文件信息
      headerFile: null,
      headerFileName: '',
      
      // 源文件信息
      sourceFile: null,
      sourceFileName: '',
      
      // 授课老师列表（从接口resp.data获取，每个元素仅含name字段）
      teachers: [],
      selectedTeacher: '',
      
      // 确认提示框显示状态
      showConfirmation: false
    }
  },
  mounted() {
    // 页面加载时获取授课老师列表（数据来源：resp.data）
    this.fetchTeachers();
  },
  methods: {
    // 处理头文件选择（仅允许.h格式）
    handleHeaderFileChange(event) {
      const file = event.target.files[0];
      if (!file) return;

      // 校验文件格式：仅允许.h
      const fileExt = file.name.split('.').pop().toLowerCase();
      if (fileExt !== 'h') {
        ElMessage.error('头文件仅支持 .h 格式，请重新选择！');
        // 清空选择
        this.headerFile = null;
        this.headerFileName = '';
        event.target.value = '';
        return;
      }

      // 格式正确，保存文件信息
      this.headerFile = file;
      this.headerFileName = file.name;
    },
    
    // 处理源文件选择（仅允许.cpp格式）
    handleSourceFileChange(event) {
      const file = event.target.files[0];
      if (!file) return;

      // 校验文件格式：仅允许.cpp
      const fileExt = file.name.split('.').pop().toLowerCase();
      if (fileExt !== 'cpp') {
        ElMessage.error('源文件仅支持 .cpp 格式，请重新选择！');
        // 清空选择
        this.sourceFile = null;
        this.sourceFileName = '';
        event.target.value = '';
        return;
      }

      // 格式正确，保存文件信息
      this.sourceFile = file;
      this.sourceFileName = file.name;
    },
    
    // 处理表单提交，显示确认提示框
    handleSubmit() {
      // 完整信息校验
      if (!this.headerFile || !this.sourceFile || !this.selectedTeacher) {
        ElMessage.warning('请填写完整信息（头文件、源文件、授课老师）后再提交！');
        return;
      }
      
      // 显示确认提示框
      this.showConfirmation = true;
    },
    
    // 确认提交
    confirmSubmission() {
      this.showConfirmation = false;
      this.submitAssessment();
    },
    
    // 网络接口：获取授课老师列表（数据从resp.data提取）
    async fetchTeachers() {
      try {
        const response = await fetch(`${config.base_url}/home/getteacher`, {
          method: "GET",
          credentials: 'include'
        });
        const resp = await response.json(); // resp为接口完整返回数据

        // 校验接口状态
        if (!response.ok || !resp.status) {
          ElMessage.error(resp.msg || "获取老师列表失败：网络异常");
          return;
        }

        // 从resp.data获取老师列表（每个元素仅含name字段）
        this.teachers = resp.data || [];
        // 若resp.data为空，提示用户
        if (this.teachers.length === 0) {
          ElMessage.info("当前暂无授课老师数据，请联系管理员");
        }

      } catch (error) {
        console.error('获取授课老师列表失败:', error);
        ElMessage.error('获取老师列表失败，请刷新页面重试');
      }
    },
    
    // 网络接口：提交考核内容
    async submitAssessment() {
      try {
        // 创建FormData对象用于文件上传
        const formData = new FormData();
        formData.append('header', this.headerFile);
        formData.append('source', this.sourceFile);
        formData.append('teacher', this.selectedTeacher);
        
        const response = await fetch(`${config.base_url}/upload/AssessmentSubmissionUpload`, {
          method: 'POST',
          credentials: 'include',
          body: formData
        });
        
        const data = await response.json();

        if (!response.ok || !data.status) {
          ElMessage.error(data.msg || "提交失败，请稍后重试");
          return;
        }

        // 提交成功提示+重置表单
        ElMessage.success("提交成功！可在截止日期前重新提交更新内容");
        this.resetForm();

      } catch (error) {
        console.error('提交考核失败:', error);
        ElMessage.error('提交考核失败，请检查网络后重试');
      }
    },
    
    // 重置表单
    resetForm() {
      // 清空文件信息
      this.headerFile = null;
      this.headerFileName = '';
      this.sourceFile = null;
      this.sourceFileName = '';
      this.selectedTeacher = '';
      
      // 重置文件输入框（避免重复选择同一文件时不触发change事件）
      document.getElementById('headerFile').value = '';
      document.getElementById('sourceFile').value = '';
    }
  }
}
</script>

<style scoped>
/* 可以添加额外的组件样式 */
</style>