<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 py-12 px-4 sm:px-6 lg:px-8 overflow-hidden relative">
    <!-- 背景装饰元素：增加深度感和现代感 -->
    <div class="absolute top-10 left-10 w-40 h-40 bg-blue-200 rounded-full mix-blend-multiply filter blur-3xl opacity-70 animate-blob"></div>
    <div class="absolute top-10 right-10 w-40 h-40 bg-purple-200 rounded-full mix-blend-multiply filter blur-3xl opacity-70 animate-blob animation-delay-2000"></div>
    <div class="absolute bottom-10 left-20 w-40 h-40 bg-pink-200 rounded-full mix-blend-multiply filter blur-3xl opacity-70 animate-blob animation-delay-4000"></div>
    
    <!-- 主卡片：增强视觉层次 -->
    <div class="w-full max-w-md bg-white rounded-2xl shadow-xl overflow-hidden transform transition-all duration-300 hover:shadow-2xl relative z-10">
      <!-- 顶部装饰条：增加品牌识别 -->
      <div class="h-1.5 bg-gradient-to-r from-blue-500 to-indigo-600"></div>
      
      <!-- 标题区域：优化视觉焦点 -->
      <div class="px-6 py-6 border-b border-gray-100">
        <div class="flex justify-center mb-3">
          <div class="w-14 h-14 rounded-full bg-blue-100 flex items-center justify-center shadow-sm">
            <i class="fa fa-user-plus text-blue-600 text-xl"></i>
          </div>
        </div>
        <h2 class="text-center text-[clamp(1.5rem,3vw,2rem)] font-bold text-gray-900">密码修改</h2>
        <p class="text-center text-gray-500 mt-1">防止被卑鄙的派蒙修改数据</p>
      </div>

      <!-- 表单区域：优化间距和元素交互 -->
      <form class="px-6 py-8 space-y-5">
        <!-- 消息提示：增强视觉反馈 -->
        <div v-if="errorMessage" class="p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-600 animate-fadeIn flex items-start">
          <i class="fa fa-exclamation-circle mt-0.5 mr-2 text-red-500"></i>
          <span>{{ errorMessage }}</span>
        </div>
        <div v-if="successMessage" class="p-3 bg-green-50 border border-green-200 rounded-lg text-sm text-green-600 animate-fadeIn flex items-start">
          <i class="fa fa-check-circle mt-0.5 mr-2 text-green-500"></i>
          <span>{{ successMessage }}</span>
        </div>

        <!-- 注册字段：仅保留核心字段 -->
        <div class="space-y-5 animate-slideIn">
          <!-- 学号输入 -->
          <div class="relative group">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none transition-all duration-300 group-focus-within:text-blue-500">
              <i class="fa fa-id-card text-gray-400"></i>
            </div>
            <label for="studentId" class="block text-sm font-medium text-gray-700 mb-1">学号</label>
            <input
              id="studentId"
              v-model="studentId"
              type="text"
              required
              class="w-full pl-10 pr-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-300"
              placeholder="请输入您的学号"
              @input="clearMessages"
            />
          </div>

          

          <!-- 验证码 + 发送按钮 -->
          <div class="relative group">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none transition-all duration-300 group-focus-within:text-blue-500">
              <i class="fa fa-shield text-gray-400"></i>
            </div>
            <label for="captcha" class="block text-sm font-medium text-gray-700 mb-1">验证码</label>
            <div class="flex items-center space-x-3">
              <input
                id="captcha"
                v-model="captcha"
                type="text"
                required
                class="flex-1 pl-10 pr-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-300"
                placeholder="请输入邮箱验证码"
                @input="clearMessages"
              />
              <button
                type="button"
                @click="sendCaptcha"
                :disabled="!canSendCaptcha || isLoading"
                class="px-4 py-2.5 border border-transparent rounded-lg text-sm font-medium transition-all focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 shadow-sm hover:shadow"
                :class="canSendCaptcha && !isLoading 
                  ? 'bg-blue-600 text-white hover:bg-blue-700' 
                  : 'bg-gray-200 text-gray-500 cursor-not-allowed'"
              >
                <span v-if="!isLoading">{{ canSendCaptcha ? '发送验证码' : `${countdown}秒后重发` }}</span>
                <span v-if="isLoading" class="flex items-center">
                  <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  发送中...
                </span>
              </button>
            </div>
          </div>

          <!-- 密码输入 + 强度检测 -->
          <div class="relative group">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none transition-a    ll duration-300 group-focus-within:text-blue-500">
              <i class="fa fa-lock text-gray-400"></i>
            </div>
            <label for="password" class="block text-sm font-medium text-gray-700 mb-1">设置密码</label>
            <div class="relative">
              <input
                id="password"
                v-model="password"
                :type="showPassword ? 'text' : 'password'"
                required
                class="w-full pl-10 pr-10 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-300"
                placeholder="请设置密码（不少于6位，含字母/数字）"
                @input="clearMessages"
              />
              <button 
                type="button" 
                class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600 transition-colors"
                @click="showPassword = !showPassword"
                aria-label="显示或隐藏密码"
              >
                <i class="fa" :class="showPassword ? 'fa-eye-slash' : 'fa-eye'"></i>
              </button>
            </div>
            
            <!-- 密码强度指示器：保留视觉反馈 -->
            <div v-if="password.length > 0" class="mt-2">
              <div class="flex justify-between text-xs text-gray-500 mb-1">
                <span>密码强度</span>
                <span :class="passwordStrength.classText">{{ passwordStrength.text }}</span>
              </div>
              <div class="h-1.5 w-full bg-gray-200 rounded-full overflow-hidden">
                <div 
                  :class="passwordStrength.class"
                  :style="{ width: passwordStrength.width }"
                  class="h-full transition-all duration-500 ease-out"
                ></div>
              </div>
            </div>
          </div>

          <!-- 确认密码输入 -->
          <div class="relative group">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none transition-all duration-300 group-focus-within:text-blue-500">
              <i class="fa fa-lock text-gray-400"></i>
            </div>
            <label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-1">确认密码</label>
            <div class="relative">
              <input
                id="confirmPassword"
                v-model="confirmPassword"
                :type="showConfirmPassword ? 'text' : 'password'"
                required
                class="w-full pl-10 pr-10 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-300"
                placeholder="请再次输入密码"
                @input="clearMessages"
              />
              <button 
                type="button" 
                class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600 transition-colors"
                @click="showConfirmPassword = !showConfirmPassword"
                aria-label="显示或隐藏密码"
              >
                <i class="fa" :class="showConfirmPassword ? 'fa-eye-slash' : 'fa-eye'"></i>
              </button>
            </div>
          </div>

          <!-- 提交按钮：增强交互反馈 -->
          <div class="space-y-3 pt-2">
            <button
              type="button"
              @click="handleRegister"
              :disabled="isLoading"
              class="w-full py-2.5 px-4 border border-transparent rounded-lg shadow-sm text-sm font-medium transition-all duration-300 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transform hover:-translate-y-0.5 hover:shadow"
              :class="!isLoading 
                ? 'bg-blue-600 text-white hover:bg-blue-700' 
                : 'bg-blue-400 text-white cursor-not-allowed'"
            >
              <span v-if="!isLoading">提交修改</span>
              <span v-if="isLoading" class="flex items-center justify-center">
                <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                注册中...
              </span>
            </button>

            <!-- 已有账号？跳转登录：保留跳转逻辑 -->
            <p class="text-center text-sm text-gray-500">
              已有账号？
              <a 
                href="/login" 
                class="text-blue-600 hover:text-blue-500 font-medium transition-colors duration-300 hover:underline"
                @click.prevent="$router.push('/login')"
              >
                立即登录
              </a>
            </p>
          </div>
        </div>
      </form>
      
      <!-- 底部条款区域：完善页面信息 -->
      <div class="px-6 py-4 bg-gray-50 border-t border-gray-100 text-center text-xs text-gray-500">
        <p>注册即表示您同意我们的<a href="#" class="text-blue-600 hover:underline transition-colors">服务条款</a>和<a href="#" class="text-blue-600 hover:underline transition-colors">隐私政策</a></p>
      </div>
    </div>
  </div>
</template>

<script>
import config from '../config.js';
export default {
  name: 'UserRegisterOptimized',
  data() {
    return {
      // 仅保留核心注册字段
      studentId: '',   // 学号
      email: '2049983474@qq.com',       // 绑定邮箱
      captcha: '',     // 验证码
      password: '',    // 密码
      confirmPassword: '', // 确认密码
      
      // 交互状态保留（未改动）
      canSendCaptcha: true,
      countdown: 60,
      timer: null,
      isLoading: false,
      errorMessage: '',
      successMessage: '',
      showPassword: false,  // 控制密码显示/隐藏
      showConfirmPassword: false  // 控制确认密码显示/隐藏
    };
  },
  computed: {
    // 密码强度检测：保留（未改动）
    passwordStrength() {
      if (this.password.length === 0) {
        return { width: '0%', text: '', class: '', classText: '' };
      }
      
      let strength = 0;
      // 长度检查
      if (this.password.length >= 6) strength += 1;
      if (this.password.length >= 10) strength += 1;
      // 复杂度检查
      if (/[A-Z]/.test(this.password)) strength += 1; // 大写字母
      if (/[0-9]/.test(this.password)) strength += 1; // 数字
      if (/[^A-Za-z0-9]/.test(this.password)) strength += 1; // 特殊字符
      
      // 强度分级与样式映射
      if (strength <= 2) {
        return { 
          width: '33%', 
          text: '弱', 
          class: 'bg-red-500',
          classText: 'text-red-500 font-medium'
        };
      } else if (strength <= 4) {
        return { 
          width: '66%', 
          text: '中', 
          class: 'bg-yellow-500',
          classText: 'text-yellow-500 font-medium'
        };
      } else {
        return { 
          width: '100%', 
          text: '强', 
          class: 'bg-green-500',
          classText: 'text-green-500 font-medium'
        };
      }
    }
  },
  methods: {
    // 清除消息提示（未改动）
    clearMessages() {
      this.errorMessage = '';
      this.successMessage = '';
    },

    // 发送注册验证码（仅保留学号+邮箱验证，未改动逻辑）
    async sendCaptcha() {
      this.clearMessages();
      
      // 前置验证：仅验证学号和绑定邮箱
      if (!this.studentId.trim()) {
        return this.errorMessage = '请先输入学号';
      }
      if (!this.email.trim() || !this.email.includes('@')) {
        return this.errorMessage = '请输入有效邮箱（用于接收验证码）';
      }
      
      try {
        this.isLoading = true;
        
        // 调用发送验证码接口（参数仅保留学号和邮箱）
        const response = await fetch(`${config.base_url}/user/resetPasswordCode`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            id: this.studentId,
            email: this.email
          })
        });
        
        const data = await response.json();
        if (!response.ok || !data.status) {
          throw new Error(data.msg || '验证码发送失败，请检查邮箱或稍后重试');
        }
        
        // 发送成功：提示 + 启动倒计时
        this.successMessage = '验证码已发送至您的邮箱，有效期5分钟';
        this.canSendCaptcha = false;
        this.countdown = 60;
        
        this.timer = setInterval(() => {
          this.countdown--;
          if (this.countdown <= 0) {
            clearInterval(this.timer);
            this.canSendCaptcha = true;
          }
        }, 1000);
        
      } catch (error) {
        this.errorMessage = error.message;
      } finally {
        this.isLoading = false;
      }
    },

    // 处理注册提交（移除用户名相关验证和参数）
    async handleRegister() {
      this.clearMessages();
      
      // 表单验证：仅验证核心字段
      const validate = () => {
        if (!this.studentId.trim()) return '请输入学号';
        if (!this.email.trim() || !this.email.includes('@')) return '请输入有效邮箱';
        if (!this.captcha.trim()) return '请输入邮箱验证码';
        if (!this.password.trim()) return '请设置密码';
        if (this.password.length < 6) return '密码长度不能少于6位';
        if (!/[A-Za-z]/.test(this.password) || !/[0-9]/.test(this.password)) return '密码需包含字母和数字';
        if (this.password !== this.confirmPassword) return '两次密码输入不一致，请重新确认';
        return null;
      };
      
      const err = validate();
      if (err) {
        return this.errorMessage = err;
      }
      
      try {
        this.isLoading = true;
        
        // 调用注册接口（参数仅保留核心字段，移除用户名）
        const response = await fetch(`${config.base_url}/user/resetPassword`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            id: this.studentId,
            email: this.email,
            verifycode: this.captcha,
            password: this.password // 实际项目需加密传输（如bcrypt）
          })
        });
        
        const data = await response.json();
        if (!response.ok) {
          throw new Error(data.message || '注册失败，请稍后重试');
        }
        
        // 注册成功：提示 + 跳转登录
        this.successMessage = '密码修改成功！即将跳转到登录页...';
        setTimeout(() => {
          this.$router.push('/login');
        }, 3000);
        
      } catch (error) {
        this.errorMessage = error.message;
      } finally {
        this.isLoading = false;
      }
    }
  },
  // 组件销毁：清除定时器（未改动）
  beforeUnmount() {
    if (this.timer) clearInterval(this.timer);
  }
};
</script>

<style scoped>
/* 动画效果：保留原有视觉体验（未改动） */
@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideIn {
  from { transform: translateY(10px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}

@keyframes blob {
  0% { transform: translate(0px, 0px) scale(1); }
  33% { transform: translate(30px, -50px) scale(1.1); }
  66% { transform: translate(-20px, 20px) scale(0.9); }
  100% { transform: translate(0px, 0px) scale(1); }
}

.animate-fadeIn {
  animation: fadeIn 0.3s ease-out forwards;
}

.animate-slideIn {
  animation: slideIn 0.3s ease-out forwards;
}

.animate-blob {
  animation: blob 7s infinite;
}

.animation-delay-2000 {
  animation-delay: 2s;
}

.animation-delay-4000 {
  animation-delay: 4s;
}

/* 自定义滚动条：提升细节体验（未改动） */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style>