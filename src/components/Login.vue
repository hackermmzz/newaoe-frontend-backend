<template>
  <div class="login-container">
    <!-- 背景装饰 -->
    <div class="login-bg"></div>
    
    <!-- 登录卡片 -->
    <div class="login-card">
      <div class="login-header">
        <h2 class="login-title">欢迎来到提瓦特大陆</h2>
        <p class="login-desc">请使用您的学号和密码登录</p>
      </div>

      <form class="login-form" @submit.prevent="handleLogin">
        <!-- 学号输入 -->
        <div class="form-group">
          <label for="studentId" class="form-label">学号</label>
          <div class="input-wrapper">
            <span class="input-icon">
              <i class="fas fa-id-card"></i>
            </span>
            <input
              id="studentId"
              type="text"
              v-model="studentId"
              @input="validateStudentId"
              maxlength="20"
              placeholder="请输入学号"
              class="form-input"
              :class="{ 'input-error': studentIdError }"
            >
          </div>
          <p v-if="studentIdError" class="error-message">{{ studentIdError }}</p>
        </div>

        <!-- 密码输入 -->
        <div class="form-group">
          <label for="password" class="form-label">密码</label>
          <div class="input-wrapper">
            <span class="input-icon">
              <i class="fas fa-lock"></i>
            </span>
            <input
              id="password"
              type="password"
              v-model="password"
              placeholder="请输入密码"
              class="form-input"
            >
          </div>
        </div>

        <!-- 记住我和忘记密码 -->
        <div class="form-actions">
          <label class="remember-me">
            <input
              type="checkbox"
              v-model="rememberMe"
              class="remember-checkbox"
            >
            <span>记住我</span>
          </label>
          <button
            type="button"
            @click="goToForgotPassword"
            class="forgot-password"
          >
            忘记密码?
          </button>
        </div>

        <!-- 登录按钮 -->
        <button
          type="submit"
          :disabled="isLoading || !isFormValid"
          class="login-button"
        >
          <span v-if="isLoading">
            <i class="fas fa-spinner fa-spin mr-2"></i>登录中...
          </span>
          <span v-else>登录</span>
        </button>

        <!-- 注册入口 -->
        <div class="register-section">
          <span>还没有账号? </span>
          <button
            type="button"
            @click="goToRegister"
            class="register-button"
          >
            立即注册
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import config from '../config.js';

export default {
  name: 'LoginPage',
  setup() {
    const router = useRouter();
    
    // 表单数据
    const studentId = ref('');
    const password = ref('');
    const rememberMe = ref(false);
    const isLoading = ref(false);
    const studentIdError = ref('');

    // 验证学号格式（数字+字母组合，20位以内）
    const validateStudentId = () => {
      if (!studentId.value) {
        studentIdError.value = '请输入学号';
        return false;
      }
      
      // 检查长度
      if (studentId.value.length > 20) {
        studentIdError.value = '学号长度不能超过20位';
        return false;
      }
      
      // 新增：验证数字+字母组合（可选，根据需求调整）
      const reg = /^[A-Za-z0-9]+$/;
      if (!reg.test(studentId.value)) {
        studentIdError.value = '学号仅支持数字和字母组合';
        return false;
      }
      
      studentIdError.value = '';
      return true;
    };

    // 检查表单是否有效
    const isFormValid = computed(() => {
      return studentId.value && password.value && !studentIdError.value;
    });

    // 处理登录
    const handleLogin = async () => {
      // 验证表单
      if (!validateStudentId()) return;
      
      isLoading.value = true;
      
      try {
        // 调用登录接口
        const response = await fetch(`${config.base_url}/user/login`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            id: studentId.value,
            password: password.value,
          }),
          credentials: 'include'
        });

        const data = await response.json();
        if (!response.ok ||!data.status) {
          throw new Error(data.msg || '登录失败，请检查学号和密码');
        }

        // 登录成功处理
        localStorage.setItem('user', JSON.stringify(data.user));
        router.push('/home');
      } catch (err) {
        alert(err.message);
      } finally {
        isLoading.value = false;
      }
    };

    // 跳转注册页面
    const goToRegister = () => {
      router.push('/regist');
    };

    // 跳转忘记密码页面
    const goToForgotPassword = () => {
      router.push('/passwordforget');
    };

    return {
      studentId,
      password,
      rememberMe,
      isLoading,
      studentIdError,
      isFormValid,
      handleLogin,
      goToRegister,
      goToForgotPassword,
      validateStudentId
    };
  }
};
</script>

<style scoped>
/* 关键修复1：统一盒模型，让width包含padding和border */
* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

.login-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px;
  position: relative;
  background-color: #f5f7fa;
}

.login-bg {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, #43cbff 0%, #9708cc 100%);
  opacity: 0.05;
  z-index: 1;
}

.login-card {
  width: 100%;
  max-width: 400px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.08);
  padding: 30px;
  position: relative;
  z-index: 2;
  transition: transform 0.3s ease;
}

.login-card:hover {
  transform: translateY(-5px);
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-title {
  font-size: 24px;
  font-weight: 700;
  color: #1d2129;
  margin-bottom: 8px;
}

.login-desc {
  color: #86909c;
  font-size: 14px;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-label {
  font-size: 14px;
  color: #1d2129;
  font-weight: 500;
}

/* 关键修复2：确保输入框父容器宽度100%，和卡片对齐 */
.input-wrapper {
  position: relative;
  width: 100%;
}

.form-input {
  width: 100%;
  padding: 12px 12px 12px 40px; /* 左侧留足图标空间，不影响宽度 */
  border: 1px solid #d9d9d9;
  border-radius: 6px;
  font-size: 14px;
  transition: all 0.2s ease;
  /* 已通过*统一box-sizing，此处无需重复写 */
}

.form-input:focus {
  outline: none;
  border-color: #4096ff;
  box-shadow: 0 0 0 2px rgba(64, 150, 255, 0.2);
}

.form-input.input-error {
  border-color: #f5222d;
}

.input-icon {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: #86909c;
  /* 关键修复3：图标不占用输入框空间，避免挤压 */
  pointer-events: none;
}

.error-message {
  color: #f5222d;
  font-size: 12px;
  margin: 0;
  height: 16px;
}

.form-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 5px;
}

.remember-me {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #4e5969;
  font-size: 14px;
  cursor: pointer;
}

.remember-checkbox {
  width: 14px;
  height: 14px;
  accent-color: #4096ff;
}

.forgot-password {
  color: #4096ff;
  background: none;
  border: none;
  padding: 0;
  font-size: 14px;
  cursor: pointer;
  transition: color 0.2s;
}

.forgot-password:hover {
  color: #1890ff;
}

.login-button {
  width: 100%;
  padding: 12px;
  background: #4096ff;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
  display: flex;
  justify-content: center;
  align-items: center;
}

.login-button:hover:not(:disabled) {
  background: #1890ff;
}

.login-button:disabled {
  background: #a8cfff;
  cursor: not-allowed;
}

.register-section {
  text-align: center;
  margin-top: 15px;
  color: #86909c;
  font-size: 14px;
}

.register-button {
  color: #4096ff;
  background: none;
  border: none;
  padding: 0;
  font-size: 14px;
  cursor: pointer;
  font-weight: 500;
  transition: color 0.2s;
}

.register-button:hover {
  color: #1890ff;
}
</style>