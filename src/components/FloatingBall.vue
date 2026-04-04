<template>
  <div class="ai-assistant-container">
    <div 
       ref="ballRef"
       class="floating-ball"
       :class="{ 'is-side': isAtSide }"
       :style="ballStyle"
       @mousedown="handleStart"
       @touchstart.passive="handleStart" 
       @click="handleClick"
    >
    <img src="@/assets/Expert_lou.jpg" class="ball-avatar" alt="AI Avatar" />
    </div>

    <transition name="chat-fade">
      <div 
        v-if="showChat" 
        class="chat-panel" 
        :style="chatStyle"
      >
        <div class="chat-header">
          <div class="header-status">
            <div class="status-dot"></div>
            <span>娄老师</span>
          </div>
          <button class="close-btn" @click="showChat = false">×</button>
        </div>

        <div class="chat-body" ref="chatBody">
          <div v-for="msg in messages" :key="msg.id" :class="['message-row', msg.role]">
            <div class="message-bubble">
              {{ msg.content }}
            </div>
          </div>
          
          <div v-if="isTyping" class="message-row ai">
            <div class="message-bubble typing-indicator">
              <span></span><span></span><span></span>
            </div>
          </div>
        </div>

        <div class="chat-footer">
          <input 
            v-model="input" 
            @keyup.enter="send" 
            :disabled="isTyping"
            placeholder="问问关于帝国时代的事..." 
          />
          <button class="send-btn" @click="send" :disabled="isTyping || !input">
            <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
             <line x1="22" y1="2" x2="11" y2="13"></line>
             <polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>
            </svg>
          </button>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, reactive, computed, nextTick } from 'vue';

const showChat = ref(false);
const isAtSide = ref(false);
const input = ref('');
const chatBody = ref(null);
const isTyping = ref(false); // 新增：是否正在等待后端返回
const hideDepth = 20;

const messages = ref([
  { id: 1, role: 'ai', content: '准备好开始你的帝国时代了吗？' }
]);

// --- 逻辑函数 ---

// 滚动到底部
const scrollToBottom = () => {
  nextTick(() => {
    if (chatBody.value) {
      chatBody.value.scrollTo({
        top: chatBody.value.scrollHeight,
        behavior: 'smooth'
      });
    }
  });
};

// 发送消息并请求后端
const send = async () => {
  if (!input.value || isTyping.value) return;

  const userContent = input.value;
  // 1. 推送用户消息
  messages.value.push({ id: Date.now(), role: 'user', content: userContent });
  input.value = '';
  isTyping.value = true;
  scrollToBottom();

  try {
    // 【预留接口位置】
    // 这里替换成后端真实 API 地址
    const response = await fetch('http://localhost:3000/api/chat', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ 
        prompt: userContent,
        history: messages.value.slice(-5) // 可选：把最近几轮对话传给后端
      })
    });

    if (!response.ok) throw new Error('网络请求失败');
    
    const data = await response.json();

    // 2. 推送 AI 消息
    messages.value.push({ 
      id: Date.now() + 1, 
      role: 'ai', 
      content: data.reply || data.content 
    });

  } catch (error) {
    console.error("API Error:", error);
    messages.value.push({ 
      id: Date.now() + 1, 
      role: 'ai', 
      content: '娄老师现在正在玩quake3，稍后再来找他聊吧。' 
    });
  } finally {
    isTyping.value = false;
    scrollToBottom();
  }
};

// --- 拖拽与样式逻辑---
const pos = reactive({ x: window.innerWidth - 80, y: window.innerHeight - 150 });
const ballStyle = computed(() => ({
  transform: `translate(${pos.x}px, ${pos.y}px)`,
  opacity: isAtSide.value ? 0.7 : 1 
}));

const chatStyle = computed(() => {
  const ballSize = 56;
  const panelWidth = 360;
  const panelHeight = 520;
  const gap = 15;
  let style = { transition: isDragging ? 'none' : 'all 0.3s ease' };
  if (pos.x + ballSize / 2 > window.innerWidth / 2) {
    style.left = `${pos.x - panelWidth - gap}px`;
  } else {
    style.left = `${pos.x + ballSize + gap}px`;
  }
  let topY = pos.y;
  if (topY + panelHeight > window.innerHeight) topY = window.innerHeight - panelHeight - 20;
  if (topY < 20) topY = 20;
  style.top = `${topY}px`;
  return style;
});

// --- 修改后的拖拽逻辑 ---
let isDragging = false;
let startPos = { x: 0, y: 0 };

const handleStart = (e) => {
  isDragging = false;
  
  // 兼容写法：获取初始点击/触摸坐标
  const clientX = e.type.includes('touch') ? e.touches[0].clientX : e.clientX;
  const clientY = e.type.includes('touch') ? e.touches[0].clientY : e.clientY;
  
  startPos = { x: clientX - pos.x, y: clientY - pos.y };

  const onMove = (ev) => {
    isDragging = true;
    // 获取移动中的坐标
    const moveX = ev.type.includes('touch') ? ev.touches[0].clientX : ev.clientX;
    const moveY = ev.type.includes('touch') ? ev.touches[0].clientY : ev.clientY;

    let newX = moveX - startPos.x;
    let newY = moveY - startPos.y;

    const ballSize = 56;
    const padding = 10;

    // Y 轴限制
    if (newY < padding) newY = padding;
    if (newY > window.innerHeight - ballSize - padding) newY = window.innerHeight - ballSize - padding;

    // X 轴限制（允许部分隐入墙内）
    if (newX < -hideDepth) newX = -hideDepth;
    if (newX > window.innerWidth - ballSize + hideDepth) {
      newX = window.innerWidth - ballSize + hideDepth;
    }

    pos.x = newX;
    pos.y = newY;
    isAtSide.value = false;

    // 移动端防止滚动页面
    if (ev.cancelable) ev.preventDefault();
  };

  const onEnd = () => {
    const threshold = 80;
    const ballSize = 56;

    // 磁吸逻辑
    if (pos.x < threshold) {
      pos.x = -hideDepth;
      isAtSide.value = true;
    } else if (pos.x > window.innerWidth - ballSize - threshold) {
      pos.x = window.innerWidth - ballSize + hideDepth;
      isAtSide.value = true;
    }

    // 移除所有监听
    document.removeEventListener('mousemove', onMove);
    document.removeEventListener('mouseup', onEnd);
    document.removeEventListener('touchmove', onMove);
    document.removeEventListener('touchend', onEnd);
  };

  // 同时注册鼠标和触摸的移动/结束事件
  document.addEventListener('mousemove', onMove);
  document.addEventListener('mouseup', onEnd);
  document.addEventListener('touchmove', onMove, { passive: false }); 
  document.addEventListener('touchend', onEnd);
};

const handleClick = () => { if (!isDragging) showChat.value = !showChat.value; };
</script>

<style scoped>
/* 悬浮球基础样式 */
.floating-ball {
  position: fixed;
  top: 0; left: 0;
  width: 56px; height: 56px;
  background: rgba(255, 255, 255, 0.1); 
  backdrop-filter: blur(5px); 
  border-radius: 50%;
  cursor: grab;
  z-index: 9999;
  display: flex; align-items: center; justify-content: center;
  box-shadow: 0 8px 24px rgba(0,0,0,0.3);
  transition: opacity 0.3s, transform 0.1s ease-out;
  overflow: hidden;
  touch-action: none; 
  user-select: none;
}
.ball-avatar { width: 100%; height: 100%; object-fit: cover; }
.floating-ball:active { cursor: grabbing; }

.chat-panel {
  position: fixed;
  width: 360px; height: 520px;
  background: rgba(28, 28, 30, 0.85);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 24px;
  display: flex; flex-direction: column;
  box-shadow: 0 20px 60px rgba(0,0,0,0.6);
  z-index: 10000;
  overflow: hidden;
  max-width: 90vw; /* 确保在任何屏幕都不会横向撑破 */
  max-height: 80vh; /* 确保不会纵向撑破 */
}

@media (max-width: 600px) {
  .chat-panel {
    /* 在手机上，对话框不再死死跟着球，而是稍微居中一些 */
    width: 95vw !important; 
    height: 60vh !important;
    left: 2.5vw !important; 
    border-radius: 20px; 
  }
}

.chat-header { 
  padding: 18px 20px; 
  display: flex; 
  justify-content: space-between; 
  align-items: center; 
  border-bottom: 1px solid rgba(255, 255, 255, 0.08); 
}
.header-status { 
  display: flex; 
  align-items: center; 
  color: #fff; 
  font-weight: 600; 
}
.status-dot { 
  width: 8px; 
  height: 8px; 
  background: #4ade80; 
  border-radius: 50%;
  margin-right: 10px; 
  box-shadow: 0 0 8px #4ade80; 
}
.close-btn { 
  background: none; 
  border: none; 
  color: #fff; 
  font-size: 20px; 
  cursor: pointer; 
  opacity: 0.5; 
}
.close-btn:hover { opacity: 1; }

.chat-body { 
  flex: 1; 
  padding: 20px; 
  overflow-y: auto; 
  display: flex; 
  flex-direction: column; 
  gap: 16px; 
}

.message-row { 
  display: flex; 
  width: 100%; 
}
.message-row.ai { justify-content: flex-start; }
.message-row.user { justify-content: flex-end; }

.message-bubble { 
  max-width: 80%; 
  padding: 12px 16px; 
  font-size: 14px; 
  line-height: 1.5; 
  word-wrap: break-word;
  text-align: left;
}

.ai .message-bubble { 
  background: rgba(255, 255, 255, 0.08); 
  color: #efefef; 
  border-radius: 18px 18px 18px 4px; 
}
.user .message-bubble { 
  background: #3b82f6; 
  color: white; 
  border-radius: 18px 18px 4px 18px; 
}

.typing-indicator span {
  display: inline-block; 
  width: 6px; 
  height: 6px; 
  background: #aaa; 
  border-radius: 50%; 
  margin-right: 3px;
  animation: typing 1s infinite;
}
.typing-indicator span:nth-child(2) { animation-delay: 0.2s; }
.typing-indicator span:nth-child(3) { animation-delay: 0.4s; }
@keyframes typing { 0%, 100% { transform: translateY(0); } 50% { transform: translateY(-5px); } }

.chat-footer { 
  padding: 16px;
  background: rgba(0,0,0,0.2); 
  display: flex; 
  gap: 12px; 
}
.chat-footer input { 
  flex: 1; 
  background: rgba(255, 255, 255, 0.05); 
  border: 1px solid rgba(255, 255, 255, 0.1); 
  border-radius: 20px; 
  padding: 10px 16px; 
  color: white; 
  outline: none; 
  font-size: 16px;
}
.chat-footer input:disabled { cursor: not-allowed; opacity: 0.5; }
.send-btn { 
  display: flex;
  align-items: center;
  justify-content: center;
  background: #3b82f6; 
  color: white; 
  border: none; 
  width: 40px; 
  height: 40px; 
  border-radius: 50%; 
  cursor: pointer; 
}
.send-btn:disabled { background: #555; cursor: not-allowed; }

.chat-fade-enter-active, .chat-fade-leave-active { transition: all 0.3s ease; }
.chat-fade-enter-from, .chat-fade-leave-to { 
  opacity: 0; 
  transform: translateY(10px) scale(0.95); 
  }
</style>