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

      <!-- 4. 账号安全 -->
      <div class="mb-8 bg-gray-50 p-6 rounded-lg flex items-center justify-between gap-4">
        <div>
          <p class="font-medium text-gray-900">账号密码</p>
          <p class="mt-1 text-sm text-gray-500">通过绑定邮箱验证码修改密码</p>
        </div>
        <router-link
          to="/passwordforget"
          class="shrink-0 rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700"
        >
          修改密码
        </router-link>
      </div>
      
      <!-- 5. 反馈建议模块 -->
      <div>
        <h3 class="text-xl font-semibold mb-4 text-center">反馈与建议</h3>
        <div class="bg-gray-50 p-6 rounded-lg">
          <form @submit.prevent="submitFeedback">
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

            <div class="mb-6">
              <label for="feedbackImages" class="block text-gray-700 mb-2">反馈图片（可多选）</label>
              <input
                id="feedbackImages"
                ref="feedbackImageInput"
                type="file"
                accept="image/*"
                multiple
                :disabled="submitting"
                class="block w-full text-sm text-gray-600 file:mr-4 file:rounded-md file:border-0 file:bg-blue-50 file:px-4 file:py-2 file:text-blue-700 hover:file:bg-blue-100"
                @change="handleFeedbackImages"
              >
              <ul v-if="feedbackImages.length" class="mt-3 space-y-2">
                <li
                  v-for="(image, index) in feedbackImages"
                  :key="`${image.name}-${image.size}-${image.lastModified}`"
                  class="flex items-center justify-between gap-3 rounded-md bg-white px-3 py-2 text-sm"
                >
                  <span class="min-w-0 truncate">{{ image.name }}（{{ formatFileSize(image.size) }}）</span>
                  <button type="button" class="shrink-0 text-red-500 hover:text-red-700" :disabled="submitting" @click="removeFeedbackImage(index)">
                    删除
                  </button>
                </li>
              </ul>
            </div>

            <div class="mb-6">
              <label for="feedbackVideos" class="block text-gray-700 mb-2">反馈视频（可多选）</label>
              <input
                id="feedbackVideos"
                ref="feedbackVideoInput"
                type="file"
                accept="video/*"
                multiple
                :disabled="submitting"
                class="block w-full text-sm text-gray-600 file:mr-4 file:rounded-md file:border-0 file:bg-blue-50 file:px-4 file:py-2 file:text-blue-700 hover:file:bg-blue-100"
                @change="handleFeedbackVideos"
              >
              <ul v-if="feedbackVideos.length" class="mt-3 space-y-2">
                <li
                  v-for="(video, index) in feedbackVideos"
                  :key="`${video.name}-${video.size}-${video.lastModified}`"
                  class="flex items-center justify-between gap-3 rounded-md bg-white px-3 py-2 text-sm"
                >
                  <span class="min-w-0 truncate">{{ video.name }}（{{ formatFileSize(video.size) }}）</span>
                  <button type="button" class="shrink-0 text-red-500 hover:text-red-700" :disabled="submitting" @click="removeFeedbackVideo(index)">
                    删除
                  </button>
                </li>
              </ul>
            </div>

            <div class="mb-6">
              <label for="feedbackFiles" class="block text-gray-700 mb-2">其他文件（可多选）</label>
              <input
                id="feedbackFiles"
                ref="feedbackFileInput"
                type="file"
                multiple
                :disabled="submitting"
                class="block w-full text-sm text-gray-600 file:mr-4 file:rounded-md file:border-0 file:bg-blue-50 file:px-4 file:py-2 file:text-blue-700 hover:file:bg-blue-100"
                @change="handleFeedbackFiles"
              >
              <p class="mt-1 text-xs text-gray-500">图片和视频请分别放入对应区域。</p>
              <ul v-if="feedbackFiles.length" class="mt-3 space-y-2">
                <li
                  v-for="(file, index) in feedbackFiles"
                  :key="`${file.name}-${file.size}-${file.lastModified}`"
                  class="flex items-center justify-between gap-3 rounded-md bg-white px-3 py-2 text-sm"
                >
                  <span class="min-w-0 truncate">{{ file.name }}（{{ formatFileSize(file.size) }}）</span>
                  <button type="button" class="shrink-0 text-red-500 hover:text-red-700" :disabled="submitting" @click="removeFeedbackFile(index)">
                    删除
                  </button>
                </li>
              </ul>
            </div>

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
import { getDownloadUrl } from '../utils/download';
export default {
  name: 'StudentHome',
  
  setup() {
    // 使用reactive创建响应式的用户信息对象
    const userInfo = reactive({
      avatar: '', // 默认头像
      studentId: '', // 学号将通过网络请求获取
      email: '', // 邮箱将通过网络请求获取
      registDate: '',// 注册日期将通过网络请求获取,
      vip: 0,
    });
    
    // 反馈内容、附件和提交状态
    const feedbackContent = ref('');
    const feedbackImages = ref([]);
    const feedbackVideos = ref([]);
    const feedbackFiles = ref([]);
    const feedbackImageInput = ref(null);
    const feedbackVideoInput = ref(null);
    const feedbackFileInput = ref(null);
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
          userInfo.avatar = await getDownloadUrl(data.data.avatar);
        }
        if (data.data.id) userInfo.studentId= data.data.id;
        if (data.data.email) userInfo.email = data.data.email;
        if(data.data.registDate) userInfo.registDate = data.data.registDate;
        if(data.data.vip) userInfo.vip=data.data.vip;
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

    
    const appendUniqueFiles = (currentFiles, incomingFiles) => {
      const fileKeys = new Set(currentFiles.map(file => `${file.name}-${file.size}-${file.lastModified}`));
      return [
        ...currentFiles,
        ...incomingFiles.filter(file => {
          const key = `${file.name}-${file.size}-${file.lastModified}`;
          if (fileKeys.has(key)) return false;
          fileKeys.add(key);
          return true;
        })
      ];
    };

    const hasFileExtension = (file, extensions) => {
      const name = String(file?.name || '').toLowerCase();
      return extensions.some(extension => name.endsWith(extension));
    };

    const isImageFile = file => (
      file.type.startsWith('image/')
      || hasFileExtension(file, ['.jpg', '.jpeg', '.png', '.gif', '.webp', '.bmp', '.svg'])
    );

    const isVideoFile = file => (
      file.type.startsWith('video/')
      || hasFileExtension(file, ['.mp4', '.webm', '.ogg', '.mov', '.m4v', '.avi', '.mkv'])
    );

    const handleFeedbackImages = (event) => {
      const selectedFiles = Array.from(event.target.files || []);
      const images = selectedFiles.filter(isImageFile);
      if (images.length !== selectedFiles.length) {
        ElMessage.warning('反馈图片区域只能选择图片文件');
      }
      feedbackImages.value = appendUniqueFiles(feedbackImages.value, images);
      event.target.value = '';
    };

    const handleFeedbackFiles = (event) => {
      const selectedFiles = Array.from(event.target.files || []);
      const normalFiles = selectedFiles.filter(file => (
        !isImageFile(file) && !isVideoFile(file)
      ));
      if (normalFiles.length !== selectedFiles.length) {
        ElMessage.warning('图片和视频请放入对应区域');
      }
      feedbackFiles.value = appendUniqueFiles(feedbackFiles.value, normalFiles);
      event.target.value = '';
    };

    const handleFeedbackVideos = (event) => {
      const selectedFiles = Array.from(event.target.files || []);
      const videos = selectedFiles.filter(isVideoFile);
      if (videos.length !== selectedFiles.length) {
        ElMessage.warning('反馈视频区域只能选择视频文件');
      }
      feedbackVideos.value = appendUniqueFiles(feedbackVideos.value, videos);
      event.target.value = '';
    };

    const removeFeedbackImage = index => {
      feedbackImages.value.splice(index, 1);
    };

    const removeFeedbackVideo = index => {
      feedbackVideos.value.splice(index, 1);
    };

    const removeFeedbackFile = index => {
      feedbackFiles.value.splice(index, 1);
    };

    const formatFileSize = size => {
      if (size < 1024) return `${size} B`;
      if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`;
      return `${(size / 1024 / 1024).toFixed(1)} MB`;
    };

    const escapeHtml = value => String(value ?? '')
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;')
      .replace(/'/g, '&#039;');

    const buildFeedbackHtml = ({ indices, studentId, submittedAt, content, images, videos, files }) => {
      const imageHtml = images.length
        ? `<section><h2>相关图片</h2>${images.map(image => `
          <figure>
            <img src="${escapeHtml(image.url)}" alt="${escapeHtml(image.name)}">
            <figcaption>${escapeHtml(image.name)}</figcaption>
          </figure>`).join('')}</section>`
        : '';
      const fileHtml = files.length
        ? `<section><h2>相关文件</h2><ul>${files.map(file => `
          <li><a href="${escapeHtml(file.url)}" target="_blank" rel="noopener noreferrer">${escapeHtml(file.name)}</a></li>`).join('')}</ul></section>`
        : '';
      const videoHtml = videos.length
        ? `<section><h2>相关视频</h2>${videos.map(video => `
          <figure>
            <video controls preload="metadata" src="${escapeHtml(video.url)}"></video>
            <figcaption>${escapeHtml(video.name)}</figcaption>
          </figure>`).join('')}</section>`
        : '';

      return `<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>反馈 #${escapeHtml(indices)}</title>
  <style>
    :root { color-scheme: light; }
    body { margin: 0; padding: 32px 20px; color: #1f2937; background: #f3f6fb; font-family: Arial, "Microsoft YaHei", sans-serif; line-height: 1.7; }
    article { max-width: 900px; margin: 0 auto; overflow: hidden; background: #ffffff; border: 1px solid #e5e7eb; border-radius: 16px; box-shadow: 0 10px 30px rgba(15, 23, 42, 0.08); }
    .feedback-header { padding: 28px 32px; color: #ffffff; background: linear-gradient(135deg, #2563eb, #4f46e5); }
    h1 { margin: 0 0 14px; font-size: 28px; line-height: 1.3; }
    .meta { display: flex; flex-wrap: wrap; gap: 8px 24px; color: #dbeafe; font-size: 14px; }
    .content { padding: 28px 32px; }
    section { margin-bottom: 28px; }
    section:last-child { margin-bottom: 0; }
    h2 { margin: 0 0 12px; color: #111827; font-size: 20px; }
    .feedback-text { margin: 0; padding: 16px 18px; white-space: normal; overflow-wrap: anywhere; color: #374151; background: #f9fafb; border-left: 4px solid #60a5fa; border-radius: 8px; }
    figure { margin: 16px 0; padding: 14px; background: #f9fafb; border: 1px solid #e5e7eb; border-radius: 10px; }
    img { display: block; max-width: 100%; max-height: 720px; height: auto; margin: 0 auto; border-radius: 8px; object-fit: contain; }
    video { display: block; width: 100%; max-height: 720px; margin: 0 auto; border-radius: 8px; background: #111827; }
    figcaption { margin-top: 8px; color: #6b7280; font-size: 13px; text-align: center; overflow-wrap: anywhere; }
    ul { margin: 0; padding-left: 22px; }
    li + li { margin-top: 8px; }
    a { color: #2563eb; overflow-wrap: anywhere; }
    @media (max-width: 640px) { body { padding: 12px; } .feedback-header, .content { padding: 22px 18px; } }
  </style>
</head>
<body>
  <article>
    <header class="feedback-header">
      <h1>反馈 #${escapeHtml(indices)}</h1>
      <div class="meta">
        <span>学生 ID：${escapeHtml(studentId || '未知')}</span>
        <span>提交时间：${escapeHtml(submittedAt)}</span>
      </div>
    </header>
    <main class="content">
      <section>
        <h2>反馈内容</h2>
        <p class="feedback-text">${escapeHtml(content).replace(/\r?\n/g, '<br>')}</p>
      </section>
      ${imageHtml}
      ${videoHtml}
      ${fileHtml}
    </main>
  </article>
</body>
</html>`;
    };

    const createApiError = message => {
      const error = new Error(message || '网络异常');
      error.apiMessage = message || '';
      return error;
    };

    const parseJsonResponse = async (response) => {
      const responseText = await response.text();
      let result = {};
      if (responseText) {
        try {
          result = JSON.parse(responseText);
        } catch (error) {
          if (!response.ok) throw createApiError();
        }
      }
      if (!response.ok || result?.status === false) {
        throw createApiError(result?.msg || result?.data?.msg);
      }
      return result?.data && typeof result.data === 'object' ? result.data : result;
    };

    const uploadAttachment = async (url, file) => {
      await axios.put(url, file, {
        headers: {
          'Content-Type': file.type || 'application/octet-stream'
        }
      });
    };

    // 反馈提交处理函数
    const submitFeedback = async () => {
      if (!feedbackContent.value.trim()) {
        ElMessage.warning('请输入反馈内容后再提交');
        return;
      }

      const images = [...feedbackImages.value];
      const videos = [...feedbackVideos.value];
      const files = [...feedbackFiles.value];
      submitting.value = true;

      try {
        // 第一步：提交客户端文件名数组，获取反馈编号和每个附件的上传链接。
        const attachmentResponse = await fetch(`${config.feedback_url}/feedbackUploadAttachment`, {
          method: 'POST',
          credentials: 'include',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            images: images.map(image => image.name),
            videos: videos.map(video => video.name),
            files: files.map(file => file.name)
          })
        });
        const attachmentData = await parseJsonResponse(attachmentResponse);
        const indices = Number(attachmentData.indices);
        const imageUploadUrls = Array.isArray(attachmentData.uploadimageurls)
          ? attachmentData.uploadimageurls
          : [];
        const fileUploadUrls = Array.isArray(attachmentData.uploadfileurls)
          ? attachmentData.uploadfileurls
          : [];
        const videoUploadUrls = Array.isArray(attachmentData.uploadvideourls)
          ? attachmentData.uploadvideourls
          : [];
        const imageDownloadUrls = Array.isArray(attachmentData.downloadimageurls)
          ? attachmentData.downloadimageurls
          : [];
        const fileDownloadUrls = Array.isArray(attachmentData.downloadfileurls)
          ? attachmentData.downloadfileurls
          : [];
        const videoDownloadUrls = Array.isArray(attachmentData.downloadvideourls)
          ? attachmentData.downloadvideourls
          : [];

        if (!Number.isInteger(indices)) {
          throw new Error('服务器未返回有效的反馈编号');
        }
        if (
          imageUploadUrls.length !== images.length
          || fileUploadUrls.length !== files.length
          || videoUploadUrls.length !== videos.length
          || imageDownloadUrls.length !== images.length
          || fileDownloadUrls.length !== files.length
          || videoDownloadUrls.length !== videos.length
        ) {
          throw new Error('服务器返回的附件上传链接数量不正确');
        }

        // 第二步：把图片和普通文件上传到各自对应的链接。
        await Promise.all([
          ...images.map((image, index) => uploadAttachment(imageUploadUrls[index], image)),
          ...videos.map((video, index) => uploadAttachment(videoUploadUrls[index], video)),
          ...files.map((file, index) => uploadAttachment(fileUploadUrls[index], file))
        ]);

        // HTML 直接引用服务器返回的下载地址。
        const embeddedImages = images.map((image, index) => ({
          name: image.name,
          url: imageDownloadUrls[index]
        }));
        const feedbackHtml = buildFeedbackHtml({
          indices,
          studentId: userInfo.studentId,
          submittedAt: new Date().toLocaleString('zh-CN', { hour12: false }),
          content: feedbackContent.value,
          images: embeddedImages,
          videos: videos.map((video, index) => ({
            name: video.name,
            url: videoDownloadUrls[index]
          })),
          files: files.map((file, index) => ({
            name: file.name,
            url: fileDownloadUrls[index]
          }))
        });

        // 第三步：提交包含文字、图片和普通文件链接的完整 HTML。
        const feedbackResponse = await fetch(`${config.feedback_url}/feedbackUpload`, {
          method: 'POST',
          credentials: 'include',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            indices,
            text: feedbackHtml
          })
        });
        await parseJsonResponse(feedbackResponse);

        ElMessage.success('反馈提交成功！感谢您的支持～');
        feedbackContent.value = '';
        feedbackImages.value = [];
        feedbackVideos.value = [];
        feedbackFiles.value = [];
        if (feedbackImageInput.value) feedbackImageInput.value.value = '';
        if (feedbackVideoInput.value) feedbackVideoInput.value.value = '';
        if (feedbackFileInput.value) feedbackFileInput.value.value = '';
      } catch (error) {
        console.error('提交反馈出错:', error);
        ElMessage.error(error?.apiMessage || error?.response?.data?.msg || '网络异常');
      } finally {
        submitting.value = false;
      }
    };
    
    // 将需要在模板中使用的变量和方法返回
    return {
      userInfo,
      feedbackContent,
      feedbackImages,
      feedbackVideos,
      feedbackFiles,
      feedbackImageInput,
      feedbackVideoInput,
      feedbackFileInput,
      submitting,
      handleAvatarChange,
      handleFeedbackImages,
      handleFeedbackVideos,
      handleFeedbackFiles,
      removeFeedbackImage,
      removeFeedbackVideo,
      removeFeedbackFile,
      formatFileSize,
      submitFeedback
    };
  }
};
</script>
