<template>
  <div class="bg-white p-6 shadow-sm min-h-[calc(100vh-10rem)]">
    <div class="max-w-3xl mx-auto">
      
      <!-- 1. 头像区域（放最上方，居中展示） -->
      <div class="mb-8 text-center">
        <div class="relative w-40 h-40 rounded-full overflow-hidden border-4 border-gray-100 mx-auto mb-4 transition-transform hover:scale-105">
          <img :src="userInfo.avatar" alt="用户头像" class="w-full h-full object-cover">
          <label class="absolute bottom-0 left-0 right-0 bg-black/50 text-white text-sm p-2 text-center cursor-pointer hover:bg-black/70 transition-colors">
            更换头像
            <input type="file" accept="image/*" class="hidden" @change="handleAvatarChange">
          </label>
        </div>
      </div>
      
      <!-- 2. 学号展示（紧跟头像下方） -->
      <div class="mb-8 bg-gray-50 p-6 rounded-lg text-center">
        <p class="text-gray-500 text-sm mb-1">学号</p>
        <p class="font-medium text-xl">{{ userInfo.studentId }}</p>
      </div>
      
      <!-- 3. 已绑定邮箱展示（注册时已绑定，仅展示不提供修改） -->
      <div class="mb-8 bg-gray-50 p-6 rounded-lg text-center">
        <p class="text-gray-500 text-sm mb-1">绑定的邮箱</p>
        <p class="font-medium text-lg text-blue-600">{{ userInfo.email }}</p>
      </div>
      
      <!-- 4. 反馈建议模块（最下方，仅保留输入框+提交按钮） -->
      <div>
        <h3 class="text-xl font-semibold mb-4 text-center">反馈与建议</h3>
        <div class="bg-gray-50 p-6 rounded-lg">
          <form @submit.prevent="submitFeedback">
            <!-- 仅保留反馈内容输入框（核心功能） -->
            <div class="mb-6">
              <label for="feedbackContent" class="block text-gray-700 mb-2">请输入您的问题或建议</label>
              <textarea 
                id="feedbackContent" 
                v-model="feedbackContent" 
                rows="5"
                class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all resize-none"
                placeholder="例如：功能使用问题、优化建议等..."
                required
              ></textarea>
            </div>
            
            <!-- 仅保留反馈提交按钮 -->
            <button 
              type="submit"
              :disabled="submitting"
              class="w-full px-6 py-3 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors disabled:bg-blue-300 font-medium"
            >
              {{ submitting ? '提交中...' : '提交反馈' }}
            </button>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, reactive } from 'vue';
import config from '../config'; // 引入配置文件
import { ElMessage } from 'element-plus';
import axios from 'axios'
export default {
  name: 'StudentHome',
  
  setup() {
    // 使用reactive创建响应式的用户信息对象
    const userInfo = reactive({
      avatar: '', // 默认头像
      studentId: '', // 学号将通过网络请求获取
      email: '', // 邮箱将通过网络请求获取
      registDate: '',// 注册日期将通过网络请求获取
    });
    
    // 反馈内容和提交状态
    const feedbackContent = ref('');
    const submitting = ref(false);
    
    // 从网络获取用户信息
    const fetchUserInfo = async () => {
      try {
        // 实际项目中替换为真实API地址
        const response = await fetch(
            `${config.base_url}/home/studentInfo`, 
            {
              method: 'GET',
              credentials: 'include' // 如果需要携带cookie
            }
        );
        const data = await response.json();
        if (!response.ok || !data.status) {
          throw new Error(data.msg||'获取用户信息失败');
        }
        // 更新用户信息
        if (data.data.avatar) {
          //拿到链接
          const response = await fetch(
            `${config.download_url}/${data.data.avatar}`, 
            {
              method: 'GET',
              credentials: 'include' // 如果需要携带cookie
            }
          );
          const res=await response.json()
          userInfo.avatar=res.data.url;
        }
        if (data.data.id) userInfo.studentId= data.data.id;
        if (data.data.email) userInfo.email = data.data.email;
        if(data.data.regist_date) userInfo.registDate = data.data.regist_date;
      } catch (error) {
        ElMessage.error(error.message)
        // 可以在这里添加错误提示给用户
      }
    };
    
    // 在组件挂载后获取用户信息
    onMounted(() => {
      fetchUserInfo();
    });
    
    // 头像上传处理函数
    const handleAvatarChange = async (e) => {
        const file = e.target.files[0];
        if (!file) return;

        // 先进行本地预览
        const reader = new FileReader();
        reader.onload = async () => {
            // 保存之前的
            try {
            // 显示上传中状态
            ElMessage.success('正在上传头像，请稍候...');
            // 发送上传请求获取上传链接
            const response = await fetch(`${config.upload_url}/avatar`, {
                method: 'POST',
                // 上传文件时不要手动设置Content-Type，浏览器会自动处理
                credentials: 'include' // 若需要携带cookie（如身份验证）
            });

            const result = await response.json();

            if (!response.ok || !result.status) {
                throw new Error(result.msg || '上传失败，请重试');
            }
            // 上传到指定链接
            const uploadUrl=result.data.urls[0];
            const resp=await axios.put(uploadUrl,file,{
              headers:{
                'Content-Type':file.type
              }
            })
            if (!resp.status){
              throw new Error('上传失败，请重试',resp.status);
            }
            //告诉服务器上传成功了
            const resp0=await fetch(`${config.uploadConfirm_url}/avatar`,{
              method: 'POST',
              credentials: 'include',
              body:JSON.stringify({
                urls:[uploadUrl]
              })
            });
            const resp00=await resp0.json()
            if (!resp0.ok || !resp00.status){
              throw new Error(resp00.msg || '上传失败，请重试');
            }
            //
            ElMessage.success('头像上传成功！');

            } catch (error) {
            // 上传失败：恢复原来的头像（假设之前有有效的头像地址）
            ElMessage.error('上传失败：' + error.message);
            }
            //刷新界面
            fetchUserInfo()
        };

        // 读取文件为DataURL（用于本地预览）
        reader.readAsDataURL(file);
        };

    
    // 反馈提交处理函数
    const submitFeedback = async () => {
      // 简单验证：反馈内容不能为空
      if (!feedbackContent.value.trim()) {
        alert('请输入反馈内容后再提交');
        return;
      }
      
      // 提交加载状态
      submitting.value = true;
      
      try {
        //首先获取上传链接
        const respForUploadUrl=await fetch(`${config.base_url}/home/feedback`,{
          method:'POST',
          credentials:'include'
        })
        const resp0=await respForUploadUrl.json()
        if (!respForUploadUrl.ok || !resp0.status){
           throw new Error('网络错误，提交失败' | resp0.msg);
        }
        //获取链接并且上传
        const uploadUrl=resp0.data.urls[0]
        const text=JSON.stringify({
          "date":Date(),
          "id":userInfo.studentId,
          "content":feedbackContent.value
        })
        const respForSucess=await axios.put(uploadUrl,text,{
           headers:{
            'Content-Type': 'application/json'
           },
           responseType:"text"
        })
        if (!respForSucess.status){
           throw new Error('网络错误，提交失败');
        }
        //
        alert('反馈提交成功！感谢您的支持～');
        feedbackContent.value = '';
      } catch (error) {
        console.error('提交反馈出错:', error);
        alert('提交失败，请稍后再试');
      } finally {
        submitting.value = false;
      }
    };
    
    // 将需要在模板中使用的变量和方法返回
    return {
      userInfo,
      feedbackContent,
      submitting,
      handleAvatarChange,
      submitFeedback
    };
  }
};
</script>
