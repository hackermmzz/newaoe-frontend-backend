<template>
  <section class="forum-post-page">
    <header class="post-page-header">
      <button class="back-button" type="button" @click="router.push('/home/forum')">← 返回论坛</button>
      <div v-if="post">
        <div class="post-author">
          <span class="author-avatar">
            <img v-if="post.avatarUrl" :src="post.avatarUrl" :alt="`${post.id || '用户'}的头像`">
            <span v-else>{{ getAvatarFallback(post.id) }}</span>
          </span>
          <div>
            <span>POST #{{ post.indices }}</span>
            <h1>{{ post.title || '未命名帖子' }}</h1>
            <p>ID：{{ post.id || '—' }} · {{ formatTime(post.submitTime) }}</p>
          </div>
        </div>
      </div>
    </header>

    <div v-if="postLoadError" class="state-card error-state">
      <strong>无法打开帖子</strong>
      <span>{{ postLoadError }}</span>
      <button type="button" @click="router.push('/home/forum')">返回帖子列表</button>
    </div>

    <template v-else-if="post">
      <article class="post-document">
        <div v-if="loadingPostContent" class="state-card">正在下载并检查帖子内容…</div>
        <div v-else-if="postContentError" class="state-card error-state">{{ postContentError }}</div>
        <div v-else class="safe-post-content" v-html="renderedContent"></div>
      </article>

      <section class="comments-section" aria-labelledby="comments-title">
        <div class="comments-heading">
          <div>
            <span>DISCUSSION</span>
            <h2 id="comments-title">评论区</h2>
          </div>
          <small>{{ comments.length }} 条已加载</small>
        </div>

        <div v-if="loadingComments && !comments.length" class="empty-comments">正在获取其他用户的评论…</div>
        <div v-else-if="commentFetchError" class="empty-comments error-state">{{ commentFetchError }}</div>
        <div v-else-if="!comments.length" class="empty-comments">还没有评论，来留下第一条回复吧。</div>
        <div v-else class="comment-list">
          <article v-for="comment in comments" :key="comment.indices" class="comment-card">
            <header>
              <div class="comment-author">
                <span class="comment-avatar">
                  <img v-if="comment.avatarUrl" :src="comment.avatarUrl" :alt="`${comment.id || '用户'}的头像`">
                  <span v-else>{{ getAvatarFallback(comment.id) }}</span>
                </span>
                <strong>{{ comment.id || '未知用户' }}</strong>
              </div>
              <time :datetime="comment.submittime">{{ formatTime(comment.submittime) }}</time>
            </header>
            <p>{{ comment.content }}</p>
          </article>
        </div>

        <button v-if="comments.length && hasMoreComments" class="load-more" type="button" :disabled="loadingComments" @click="loadMoreComments">
          {{ loadingComments ? '加载中…' : '加载更多评论' }}
        </button>

        <div class="comment-composer">
          <textarea v-model="commentContent" rows="4" maxlength="3000" placeholder="友善交流，写下你的评论…" :disabled="submittingComment" @keydown.ctrl.enter.prevent="submitComment"></textarea>
          <div>
            <small>Ctrl + Enter 快速发表</small>
            <button type="button" :disabled="submittingComment || !commentContent.trim()" @click="submitComment">
              {{ submittingComment ? '发表中…' : '发表评论' }}
            </button>
          </div>
        </div>
      </section>
    </template>
  </section>
</template>

<script>
export default { name: 'ForumPostPage' };
</script>

<script setup>
import { onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import DOMPurify from 'dompurify';
import { ElMessage } from 'element-plus';
import config from '../config';
import { downloadTextFile, getDownloadUrl } from '../utils/download';

const route = useRoute();
const router = useRouter();
const commentPageSize = Number(config.ForumCommentRecordPerPage || config.ForumRecordPerPage || 10);
const post = ref(null);
const postLoadError = ref('');
const loadingPostContent = ref(false);
const postContentError = ref('');
const renderedContent = ref('');
const comments = ref([]);
const commentContent = ref('');
const currentCommentPage = ref(1);
const hasMoreComments = ref(true);
const loadingComments = ref(false);
const submittingComment = ref(false);
const commentFetchError = ref('');
let commentRequestId = 0;

const parseApiResponse = async response => {
  const responseText = await response.text();
  let result = {};
  if (responseText) {
    try {
      result = JSON.parse(responseText);
    } catch (error) {
      throw new Error(`服务器返回了无法解析的数据（HTTP ${response.status}）`);
    }
  }
  if (!response.ok || result?.status === false) throw new Error(result?.msg || `请求失败（HTTP ${response.status}）`);
  return result?.data ?? result;
};

const extractRecords = payload => {
  if (Array.isArray(payload)) return payload;
  if (Array.isArray(payload?.record)) return payload.record;
  if (Array.isArray(payload?.records)) return payload.records;
  return [];
};

const formatTime = value => {
  if (!value) return '时间未知';
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? String(value) : date.toLocaleString('zh-CN', { hour12: false });
};

const getAvatarFallback = id => String(id || '?').trim().slice(0, 1).toUpperCase() || '?';

const resolveAvatarUrl = async avatar => {
  if (!avatar) return '';
  try {
    const download = await getDownloadUrl(avatar);
    return download.url;
  } catch (error) {
    return '';
  }
};

const loadPostMetadata = () => {
  const indices = String(route.params.indices || '');
  let candidate = window.history.state?.forumPost;
  if (!candidate || String(candidate.indices) !== indices) {
    try {
      candidate = JSON.parse(window.sessionStorage.getItem(`forum-post-${indices}`) || 'null');
    } catch (error) {
      candidate = null;
    }
  }
  if (!candidate || String(candidate.indices) !== indices) {
    postLoadError.value = '缺少帖子信息，请从论坛列表重新进入。';
    return false;
  }
  post.value = {
    indices: candidate.indices,
    id: candidate.id,
    title: candidate.title,
    contentPath: candidate.contentPath,
    submitTime: candidate.submitTime,
    avatar: candidate.avatar,
    avatarUrl: ''
  };
  return true;
};

const getRelativeAssetIdentifier = (contentPath, assetName) => {
  if (!assetName || /^(?:https?:|data:|blob:|#)/i.test(assetName)) return assetName;
  const cleanAssetName = assetName.replace(/^\.\//, '');
  if (cleanAssetName.includes('/')) return cleanAssetName;
  const normalizedPath = String(contentPath || '').split(/[?#]/)[0];
  const slashIndex = normalizedPath.lastIndexOf('/');
  return slashIndex >= 0 ? `${normalizedPath.slice(0, slashIndex + 1)}${cleanAssetName}` : cleanAssetName;
};

const sanitizeAndHydrateContent = async (html, contentPath) => {
  const safeHtml = DOMPurify.sanitize(html, {
    FORBID_TAGS: ['script', 'style', 'iframe', 'object', 'embed', 'form', 'input', 'button', 'base'],
    FORBID_ATTR: ['onerror', 'onload', 'onclick', 'onmouseover', 'style'],
    ADD_ATTR: ['controls', 'preload', 'data-forum-asset', 'data-asset-kind']
  });
  const template = document.createElement('template');
  template.innerHTML = safeHtml;

  const mediaParents = new Set(
    Array.from(template.content.querySelectorAll('.forum-embed:not(.forum-media-row .forum-embed)'))
      .map(element => element.parentElement)
      .filter(Boolean)
  );
  mediaParents.forEach(parent => {
    let activeRow = null;
    Array.from(parent.children).forEach(child => {
      if (child.classList.contains('forum-embed')) {
        if (!activeRow) {
          activeRow = document.createElement('div');
          activeRow.className = 'forum-media-row';
          parent.insertBefore(activeRow, child);
        }
        activeRow.appendChild(child);
        return;
      }
      const isEmptySpacer = child.tagName === 'P' && !child.textContent.trim() && !child.querySelector('img, video, a');
      if (activeRow && isEmptySpacer) {
        child.remove();
        return;
      }
      activeRow = null;
    });
  });

  const assetElements = Array.from(template.content.querySelectorAll('[data-forum-asset]'));
  await Promise.all(assetElements.map(async element => {
    const identifier = getRelativeAssetIdentifier(contentPath, element.dataset.forumAsset);
    try {
      const download = await getDownloadUrl(identifier);
      element.setAttribute(element.tagName === 'A' ? 'href' : 'src', download.url);
      if (element.tagName === 'A') {
        element.target = '_blank';
        element.rel = 'noopener noreferrer';
      }
    } catch (error) {
      element.removeAttribute(element.tagName === 'A' ? 'href' : 'src');
      element.classList.add('unavailable-asset');
    }
  }));
  return template.innerHTML;
};

const loadPostContent = async () => {
  loadingPostContent.value = true;
  try {
    if (!post.value?.contentPath) throw new Error('帖子缺少正文文件地址');
    const downloadedFile = await downloadTextFile(post.value.contentPath);
    const html = downloadedFile.content;
    if (!html.trim()) throw new Error('帖子正文为空');
    renderedContent.value = await sanitizeAndHydrateContent(html, post.value.contentPath);
  } catch (error) {
    postContentError.value = error?.message || '帖子正文加载失败';
  } finally {
    loadingPostContent.value = false;
  }
};

const fetchComments = async (page = 1, append = false) => {
  if (!post.value?.indices) return;
  const requestId = ++commentRequestId;
  loadingComments.value = true;
  commentFetchError.value = '';
  try {
    const begin = (page - 1) * commentPageSize;
    const requestUrl = new URL(config.forum_comment_fetch_url);
    requestUrl.searchParams.set('range', `${begin}:${begin + commentPageSize - 1}`);
    requestUrl.searchParams.set('indices', String(post.value.indices));
    requestUrl.searchParams.set('indicesForForumAnnoucement', String(post.value.indices));
    const response = await fetch(requestUrl.toString(), { method: 'GET', credentials: 'include' });
    const normalizedRecords = extractRecords(await parseApiResponse(response)).map(comment => ({
      indices: comment?.indices,
      id: comment?.id,
      content: comment?.content,
      submittime: comment?.submittime,
      avatar: comment?.avatar,
      avatarUrl: ''
    }));
    const records = await Promise.all(normalizedRecords.map(async comment => ({
      ...comment,
      avatarUrl: await resolveAvatarUrl(comment.avatar)
    })));
    if (requestId !== commentRequestId) return;
    comments.value = append ? [...comments.value, ...records] : records;
    currentCommentPage.value = page;
    hasMoreComments.value = records.length >= commentPageSize;
  } catch (error) {
    if (requestId === commentRequestId) commentFetchError.value = error?.message || '评论加载失败';
  } finally {
    if (requestId === commentRequestId) loadingComments.value = false;
  }
};

const loadMoreComments = () => fetchComments(currentCommentPage.value + 1, true);

const submitComment = async () => {
  const content = commentContent.value.trim();
  if (!content || !post.value?.indices || submittingComment.value) return;
  submittingComment.value = true;
  try {
    const response = await fetch(config.forum_comment_url, {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ content, indicesForForumAnnoucement: post.value.indices })
    });
    await parseApiResponse(response);
    commentContent.value = '';
    ElMessage.success('评论发表成功');
    await fetchComments(1);
  } catch (error) {
    ElMessage.error(error?.message || '评论发表失败');
  } finally {
    submittingComment.value = false;
  }
};

onMounted(() => {
  if (!loadPostMetadata()) return;
  if (post.value.avatar) {
    resolveAvatarUrl(post.value.avatar).then(url => {
      if (post.value) post.value.avatarUrl = url;
    });
  }
  loadPostContent();
  fetchComments(1);
});
</script>

<style scoped>
.forum-post-page { --post-surface: #fff; --post-soft: #f8fafc; --post-border: #e2e8f0; --post-text: #1e293b; --post-muted: #64748b; max-width: 980px; margin: 0 auto; color: var(--post-text); text-align: left; }
.post-page-header { margin-bottom: 20px; padding: 25px 28px; display: flex; align-items: flex-start; gap: 22px; color: #fff; background: linear-gradient(135deg,#1d4ed8,#4338ca); border-radius: 18px; }
.post-page-header span { color: #bfdbfe; font-size: 11px; font-weight: 800; letter-spacing: .14em; }
.post-page-header h1 { margin: 5px 0 7px; font-size: 28px; }
.post-page-header p { margin: 0; color: #dbeafe; font-size: 12px; }
.post-author { display: flex; align-items: center; gap: 14px; }
.author-avatar { width: 50px; height: 50px; flex: 0 0 50px; display: grid; place-items: center; overflow: hidden; color: #1d4ed8; background: #dbeafe; border-radius: 50%; font-weight: 800; }
.author-avatar img,.comment-avatar img { width: 100%; height: 100%; object-fit: cover; }
.back-button { flex: 0 0 auto; padding: 9px 12px; color: #dbeafe; background: rgba(255,255,255,.1); border: 1px solid rgba(255,255,255,.2); border-radius: 9px; cursor: pointer; }
.post-document,.comments-section { padding: 28px; background: var(--post-surface); border: 1px solid var(--post-border); border-radius: 17px; box-shadow: 0 10px 30px rgba(15,23,42,.06); backdrop-filter: blur(16px); }
.safe-post-content { min-height: 180px; color: var(--post-text); line-height: 1.8; overflow-wrap: anywhere; }
.safe-post-content :deep(h1) { color: var(--post-text); }
.safe-post-content :deep(img),.safe-post-content :deep(video) { display: block; max-width: 100%; max-height: 680px; margin: 16px auto; border-radius: 10px; object-fit: contain; }
.safe-post-content :deep(video) { width: 100%; background: #0f172a; }
.safe-post-content :deep(figure),.safe-post-content :deep(.forum-file-card) { margin: 18px 0; padding: 13px; background: var(--post-soft); border: 1px solid var(--post-border); border-radius: 11px; }
.safe-post-content :deep(.forum-media-row) { margin: 18px 0; display: grid; grid-auto-columns: minmax(0, 1fr); grid-auto-flow: column; gap: 8px; }
.safe-post-content :deep(.forum-media-row .forum-embed) { min-width: 0; margin: 0; padding: 7px; overflow: hidden; }
.safe-post-content :deep(.forum-media-row img),.safe-post-content :deep(.forum-media-row video) { width: 100%; margin: 0; }
.safe-post-content :deep(.forum-media-row .forum-embed:not(:only-child) img),.safe-post-content :deep(.forum-media-row .forum-embed:not(:only-child) video) { aspect-ratio: 1 / 1; object-fit: cover; }
.safe-post-content :deep(figcaption) { color: var(--post-muted); text-align: center; font-size: 12px; }
.safe-post-content :deep(a) { color: #60a5fa; }
.comments-section { margin-top: 18px; }
.comments-heading { display: flex; align-items: center; justify-content: space-between; }
.comments-heading span { color: #6366f1; font-size: 11px; font-weight: 800; letter-spacing: .14em; }
.comments-heading h2 { margin: 4px 0 0; }
.comments-heading small { color: var(--post-muted); }
.empty-comments,.state-card { margin: 18px 0; padding: 25px; color: var(--post-muted); text-align: center; background: var(--post-soft); border: 1px dashed var(--post-border); border-radius: 11px; }
.state-card { display: flex; flex-direction: column; align-items: center; gap: 9px; }
.state-card button { padding: 8px 12px; color: #fff; background: #2563eb; border: 0; border-radius: 8px; cursor: pointer; }
.error-state { color: #fb7185; }
.comment-list { margin: 18px 0; display: grid; gap: 10px; }
.comment-card { padding: 16px 18px; background: var(--post-soft); border: 1px solid var(--post-border); border-radius: 11px; }
.comment-card header { display: flex; justify-content: space-between; gap: 15px; }
.comment-author { display: flex; align-items: center; gap: 9px; }
.comment-avatar { width: 32px; height: 32px; flex: 0 0 32px; display: grid; place-items: center; overflow: hidden; color: #fff; background: linear-gradient(135deg,#3b82f6,#6366f1); border-radius: 50%; font-size: 11px; font-weight: 800; }
.comment-card strong { color: #60a5fa; font-size: 13px; }
.comment-card time { color: var(--post-muted); font-size: 11px; }
.comment-card p { margin: 9px 0 0; color: var(--post-text); line-height: 1.65; white-space: pre-wrap; overflow-wrap: anywhere; }
.load-more { margin: 0 auto 16px; padding: 8px 14px; display: block; color: #2563eb; background: #eff6ff; border: 0; border-radius: 8px; cursor: pointer; }
.comment-composer { padding: 14px; background: var(--post-soft); border: 1px solid var(--post-border); border-radius: 11px; }
.comment-composer textarea { box-sizing: border-box; width: 100%; padding: 12px; color: var(--post-text); background: var(--post-surface); border: 1px solid var(--post-border); border-radius: 9px; resize: vertical; }
.comment-composer > div { margin-top: 9px; display: flex; align-items: center; justify-content: space-between; }
.comment-composer small { color: var(--post-muted); }
.comment-composer button { padding: 8px 14px; color: #fff; background: #2563eb; border: 0; border-radius: 8px; cursor: pointer; font-weight: 700; }
.comment-composer button:disabled { cursor: not-allowed; opacity: .5; }
@media (max-width:700px) { .post-page-header { padding: 20px; flex-direction: column; } .post-document,.comments-section { padding: 18px; } }
</style>

<style>
html[data-theme='dark'] .forum-post-page { --post-surface: rgba(15,23,42,.94); --post-soft: rgba(30,41,59,.72); --post-border: rgba(148,163,184,.28); --post-text: #e2e8f0; --post-muted: #94a3b8; }
html[data-theme='effect'] .forum-post-page { --post-surface: rgba(15,23,42,.68); --post-soft: rgba(30,41,59,.58); --post-border: rgba(148,163,184,.3); --post-text: #e2e8f0; --post-muted: #a8b5c7; }
html[data-theme='custom'] .forum-post-page { --post-surface: rgba(255,255,255,.9); --post-soft: rgba(248,250,252,.85); }
</style>
