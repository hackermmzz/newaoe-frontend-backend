<template>
  <div class="register-container">
    <!-- 背景装饰 -->
    <div class="register-bg"></div>
    
    <!-- 注册卡片 -->
    <div class="register-card">
      <div class="register-header">
        <h2 class="register-title">提瓦特通行证注册</h2>
        <p class="register-desc">请填写以下信息完成注册</p>
      </div>

      <!-- 注册表单：新增邮箱输入框（位于学号和密码之间） -->
      <form class="register-form" @submit.prevent="handleRegister">
        <!-- 1. 学号输入 -->
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
              placeholder="请输入学号（数字+字母组合）"
              class="form-input"
              :class="{ 'input-error': studentIdError }"
            >
          </div>
          <p v-if="studentIdError" class="error-message">{{ studentIdError }}</p>
        </div>

        <!-- 新增：2. 邮箱输入 -->
        <div class="form-group">
          <label for="email" class="form-label">邮箱</label>
          <div class="input-wrapper">
            <span class="input-icon">
              <i class="fas fa-envelope"></i> <!-- 邮箱专用图标 -->
            </span>
            <input
              id="email"
              type="email" 
              v-model="email"
              @input="validateEmail"
              maxlength="50"
              placeholder="请输入您的邮箱（用于接收验证码）"
              class="form-input"
              :class="{ 'input-error': emailError }"
            >
          </div>
          <p v-if="emailError" class="error-message">{{ emailError }}</p>
        </div>

        <!-- 3. 密码输入 -->
        <div class="form-group">
          <label for="password" class="form-label">密码</label>
          <div class="input-wrapper">
            <span class="input-icon">
              <i class="fas fa-lock"></i>
            </span>
            <input
              id="password"
              :type="showPassword ? 'text' : 'password'"
              v-model="password"
              @input="validatePassword"
              maxlength="20"
              placeholder="请输入6-20位密码"
              class="form-input"
              :class="{ 'input-error': passwordError }"
            >
            <span class="input-action" @click="showPassword = !showPassword">
              <i :class="showPassword ? 'fas fa-eye-slash' : 'fas fa-eye'"></i>
            </span>
          </div>
          <p v-if="passwordError" class="error-message">{{ passwordError }}</p>
        </div>

        <!-- 4. 确认密码输入 -->
        <div class="form-group">
          <label for="confirmPassword" class="form-label">确认密码</label>
          <div class="input-wrapper">
            <span class="input-icon">
              <i class="fas fa-lock"></i>
            </span>
            <input
              id="confirmPassword"
              :type="showConfirmPassword ? 'text' : 'password'"
              v-model="confirmPassword"
              @input="validateConfirmPassword"
              maxlength="20"
              placeholder="请再次输入密码"
              class="form-input"
              :class="{ 'input-error': confirmPasswordError }"
            >
            <span class="input-action" @click="showConfirmPassword = !showConfirmPassword">
              <i :class="showConfirmPassword ? 'fas fa-eye-slash' : 'fas fa-eye'"></i>
            </span>
          </div>
          <p v-if="confirmPasswordError" class="error-message">{{ confirmPasswordError }}</p>
        </div>

        <!-- 5. 验证码输入 -->
        <div class="form-group">
          <label for="code" class="form-label">验证码</label>
          <div class="input-code-wrapper">
            <div class="input-wrapper code-input-inner">
              <span class="input-icon">
                <i class="fas fa-shield-alt"></i>
              </span>
              <input
                id="code"
                type="text"
                v-model="code"
                @input="validateCode"
                maxlength="5"
                placeholder="请输入5位验证码"
                class="form-input"
                :class="{ 'input-error': codeError }"
              >
            </div>
            <button
              type="button"
              class="code-button"
              :disabled="isCodeLoading || countdown > 0"
              @click="sendCode"
            >
              <span v-if="isCodeLoading">
                <i class="fas fa-spinner fa-spin mr-1"></i>发送中
              </span>
              <span v-else-if="countdown > 0">
                {{ countdown }}秒后重发
              </span>
              <span v-else>
                获取验证码
              </span>
            </button>
          </div>
          <p v-if="codeError" class="error-message">{{ codeError }}</p>
        </div>

        <!-- 注册按钮 -->
        <button
          type="submit"
          :disabled="isLoading || !isFormValid"
          class="register-button"
        >
          <span v-if="isLoading">
            <i class="fas fa-spinner fa-spin mr-2"></i>注册中...
          </span>
          <span v-else>完成注册</span>
        </button>

        <!-- 跳转登录入口 -->
        <div class="login-section">
          <span>已有账号? </span>
          <button
            type="button"
            @click="goToLogin"
            class="login-button"
          >
            立即登录
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import { ref, computed, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import config from '../config.js';

export default {
  name: 'RegisterPage',
  setup() {
    const router = useRouter();
    let countdownTimer = null;

    // 1. 表单响应式数据：新增email字段
    const studentId = ref('');
    const email = ref(''); // 新增：邮箱响应式变量
    const password = ref('');
    const confirmPassword = ref('');
    const code = ref('');
    const showPassword = ref(false);
    const showConfirmPassword = ref(false);
    const isLoading = ref(false);
    const isCodeLoading = ref(false);
    const countdown = ref(0);

    // 2. 错误提示信息：新增emailError
    const studentIdError = ref('');
    const emailError = ref(''); // 新增：邮箱错误提示
    const passwordError = ref('');
    const confirmPasswordError = ref('');
    const codeError = ref('');

    // 3. 字段验证函数：新增邮箱验证（validateEmail）
    // 3.1 验证学号
    const validateStudentId = () => {
      if (!studentId.value.trim()) {
        studentIdError.value = '请输入学号';
        return false;
      }
      if (studentId.value.length > 20) {
        studentIdError.value = '学号长度不能超过20位';
        return false;
      }
      const reg = /^[A-Za-z0-9]+$/;
      if (!reg.test(studentId.value.trim())) {
        studentIdError.value = '学号仅支持数字和字母组合';
        return false;
      }
      studentIdError.value = '';
      return true;
    };

    // 新增：3.2 验证邮箱（非空+格式正确）
    const validateEmail = () => {
      const emailVal = email.value.trim();
      if (!emailVal) {
        emailError.value = '请输入邮箱';
        return false;
      }
      // 邮箱格式正则（覆盖绝大多数常见邮箱格式）
      const emailReg = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
      if (!emailReg.test(emailVal)) {
        emailError.value = '请输入正确的邮箱格式（如：xxx@xx.com）';
        return false;
      }
      // 限制邮箱长度（可选，根据后端需求调整）
      if (emailVal.length > 50) {
        emailError.value = '邮箱长度不能超过50位';
        return false;
      }
      emailError.value = '';
      return true;
    };

    // 3.3 验证密码
    const validatePassword = () => {
      if (!password.value) {
        passwordError.value = '请输入密码';
        return false;
      }
      if (password.value.length < 6 || password.value.length > 20) {
        passwordError.value = '密码长度需在6-20位之间';
        return false;
      }
      const pureNumReg = /^[0-9]+$/;
      const pureLetterReg = /^[A-Za-z]+$/;
      if (pureNumReg.test(password.value) || pureLetterReg.test(password.value)) {
        passwordError.value = '密码不能为纯数字或纯字母';
        return false;
      }
      passwordError.value = '';
      validateConfirmPassword(); // 同步验证确认密码
      return true;
    };

    // 3.4 验证确认密码
    const validateConfirmPassword = () => {
      if (!confirmPassword.value) {
        confirmPasswordError.value = '请确认密码';
        return false;
      }
      if (confirmPassword.value !== password.value) {
        confirmPasswordError.value = '两次输入的密码不一致';
        return false;
      }
      confirmPasswordError.value = '';
      return true;
    };

    // 3.5 验证验证码
    const validateCode = () => {
      const codeVal = code.value.trim();
      if (!codeVal) {
        codeError.value = '请输入验证码';
        return false;
      }
      if (codeVal.length !== 5) {
        codeError.value = '验证码需为5位';
        return false;
      }
      const reg = /^[0-9]+$/;
      if (!reg.test(codeVal)) {
        codeError.value = '验证码仅支持数字';
        return false;
      }
      codeError.value = '';
      return true;
    };

    // 4. 表单整体有效性：新增邮箱验证（validateEmail()）
    const isFormValid = computed(() => {
      return (
        validateStudentId() &&
        validateEmail() && // 新增：必须通过邮箱验证
        validatePassword() &&
        validateConfirmPassword() &&
        validateCode() 
      );
    });

    // 5. 发送验证码：新增传递email参数（不再传空字符串）
    const sendCode = async () => {
      // 新增：发送验证码前同时验证学号和邮箱
      if (!validateStudentId() || !validateEmail()) {
        alert('请先输入正确的学号和邮箱');
        return;
      }

      isCodeLoading.value = true;
      try {
        const response = await fetch(`${config.base_url}/user/registCode`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ 
            id: studentId.value.trim(),
            email: email.value.trim() // 新增：传递用户输入的邮箱（不再是空字符串）
          }),
          credentials: 'include'
        });

        const data = await response.json();
        console.log(data);
        if (!response.ok || !data.status) {
          throw new Error(data.message || '验证码发送失败');
        }

        // 倒计时逻辑不变
        countdown.value = 60;
        countdownTimer = setInterval(() => {
          countdown.value--;
          if (countdown.value <= 0) {
            clearInterval(countdownTimer);
          }
        }, 1000);

        alert(`验证码已发送至您的邮箱：${email.value.trim()}，请查收`); // 提示用户邮箱
      } catch (err) {
        alert(err.message);
      } finally {
        isCodeLoading.value = false;
      }
    };

    // 6. 注册提交：新增传递email参数给后端
    const handleRegister = async () => {
      if (!isFormValid.value) return;

      isLoading.value = true;
      try {
        const response = await fetch(`${config.base_url}/user/regist`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            id: studentId.value.trim(),
            email: email.value.trim(), // 新增：注册时携带邮箱
            password: password.value,
            verifycode: code.value.trim()
          }),
          credentials: 'include'
        });

        const data = await response.json();
        if (!response.ok || !data.status) {
          throw new Error(data.message || '注册失败，请稍后重试');
        }

        alert('注册成功！请登录');
        router.push('/login');
      } catch (err) {
        alert(err.message);
      } finally {
        isLoading.value = false;
      }
    };

    // 7. 跳转登录页
    const goToLogin = () => {
      router.push('/login');
    };

    // 8. 清除定时器
    onUnmounted(() => {
      if (countdownTimer) clearInterval(countdownTimer);
    });

    // 返回数据：新增email、emailError、validateEmail
    return {
      studentId,
      email, // 新增
      password,
      confirmPassword,
      code,
      showPassword,
      showConfirmPassword,
      isLoading,
      isCodeLoading,
      countdown,
      studentIdError,
      emailError, // 新增
      passwordError,
      confirmPasswordError,
      codeError,
      isFormValid,
      validateStudentId,
      validateEmail, // 新增
      validatePassword,
      validateConfirmPassword,
      validateCode,
      sendCode,
      handleRegister,
      goToLogin
    };
  }
};
</script>

<style scoped>
/* 原有样式不变（新增的邮箱输入框会复用现有类，无需额外写样式） */
* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

.register-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px;
  position: relative;
  background-color: #f5f7fa;
}

.register-bg {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, #43cbff 0%, #9708cc 100%);
  opacity: 0.05;
  z-index: 1;
}

.register-card {
  width: 100%;
  max-width: 420px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.08);
  padding: 30px;
  position: relative;
  z-index: 2;
  transition: transform 0.3s ease;
}

.register-card:hover {
  transform: translateY(-5px);
}

.register-header {
  text-align: center;
  margin-bottom: 30px;
}

.register-title {
  font-size: 24px;
  font-weight: 700;
  color: #1d2129;
  margin-bottom: 8px;
}

.register-desc {
  color: #86909c;
  font-size: 14px;
}

.register-form {
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

.input-wrapper {
  position: relative;
  width: 100%;
}

.input-code-wrapper {
  display: flex;
  gap: 10px;
  align-items: center;
}

.code-input-inner {
  flex: 1;
}

.form-input {
  width: 100%;
  padding: 12px 12px 12px 40px;
  border: 1px solid #d9d9d9;
  border-radius: 6px;
  font-size: 14px;
  transition: all 0.2s ease;
}

.input-action {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: #86909c;
  cursor: pointer;
  pointer-events: auto;
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
  pointer-events: none;
}

.error-message {
  color: #f5222d;
  font-size: 12px;
  margin: 0;
  height: 16px;
}

.code-button {
  padding: 12px 16px;
  background: #4096ff;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
  transition: background 0.2s;
  white-space: nowrap;
}

.code-button:disabled {
  background: #a8cfff;
  cursor: not-allowed;
}

.register-button {
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

.register-button:hover:not(:disabled) {
  background: #1890ff;
}

.register-button:disabled {
  background: #a8cfff;
  cursor: not-allowed;
}

.login-section {
  text-align: center;
  margin-top: 15px;
  color: #86909c;
  font-size: 14px;
}

.login-button {
  color: #4096ff;
  background: none;
  border: none;
  padding: 0;
  font-size: 14px;
  cursor: pointer;
  font-weight: 500;
  transition: color 0.2s;
}

.login-button:hover {
  color: #1890ff;
}
</style>