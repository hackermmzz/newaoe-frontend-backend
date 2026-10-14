<template>
  <section class="forum-page">
    <header class="forum-hero">
      <div class="hero-copy">
        <span class="forum-eyebrow">COMMUNITY</span>
        <h1>论坛</h1>
        <p>分享经验、记录想法，在这里和大家认真交流。</p>
      </div>
      <router-link class="primary-button" to="/home/forum/new">
        <span class="button-icon" aria-hidden="true">＋</span>
        <span>发布帖子</span>
      </router-link>
    </header>

    <div class="forum-toolbar">
      <div>
        <span class="section-kicker">LATEST POSTS</span>
        <h2>最新帖子</h2>
      </div>
      <span v-if="posts.length" class="post-count">已加载 {{ posts.length }} 篇</span>
    </div>

    <div v-if="loadingPosts && !posts.length" class="state-card">正在加载帖子…</div>
    <div v-else-if="!posts.length" class="state-card">
      <strong>还没有帖子</strong>
      <span>来发布论坛里的第一篇内容吧。</span>
    </div>
    <div v-else class="post-list">
      <article
        v-for="post in posts"
        :key="post.indices"
        class="post-card"
        role="link"
        tabindex="0"
        @click="openPost(post)"
        @keydown.enter.prevent="openPost(post)"
        @keydown.space.prevent="openPost(post)"
      >
        <span class="post-card-accent" aria-hidden="true"></span>
        <span class="post-card-content">
          <span class="post-meta-row">
            <span class="post-author">
              <span class="post-avatar">
                <img v-if="post.avatarUrl" :src="post.avatarUrl" :alt="`${post.id || '用户'}的头像`">
                <span v-else>{{ getAvatarFallback(post.id) }}</span>
              </span>
              <span class="author-copy">
                <strong>{{ post.id || '未知用户' }}</strong>
                <time :datetime="post.submitTime">{{ formatTime(post.submitTime) }}</time>
              </span>
            </span>
            <span class="post-index">#{{ post.indices }}</span>
          </span>

          <span class="post-copy">
            <strong class="post-title">{{ post.title || '未命名帖子' }}</strong>
            <span :class="['post-summary', { 'is-empty': !post.summary }]">
              {{ post.summary || '作者暂未提供正文摘要，点击查看完整内容。' }}
            </span>
          </span>

          <span class="post-footer">
            <button
              v-if="canDeletePost(post)"
              class="delete-post-button"
              type="button"
              :disabled="deletingPostIndex === post.indices"
              @click.stop="removePost(post)"
            >
              {{ deletingPostIndex === post.indices ? '删除中…' : '删除帖子' }}
            </button>
            <span class="read-more">
              <span>阅读全文</span>
              <span class="post-arrow" aria-hidden="true">→</span>
            </span>
          </span>
        </span>
      </article>
    </div>

    <button v-if="posts.length && hasMorePosts" class="load-more" type="button" :disabled="loadingPosts" @click="loadMorePosts">
      {{ loadingPosts ? '加载中…' : '加载更多' }}
    </button>
  </section>
</template>

<script>
export default { name: 'ForumPage' };
</script>

<script setup>
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import config from '../config';
import { getDownloadUrl } from '../utils/download';

const router = useRouter();
const pageSize = Number(config.ForumRecordPerPage || config.HistoryRecordPerPage || 10);
const posts = ref([]);
const currentPage = ref(1);
const hasMorePosts = ref(true);
const loadingPosts = ref(false);
const deletingPostIndex = ref(null);
const currentUserId = ref('');
const currentUserVip = ref(config.VIP_NONE);

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

const normalizePost = post => ({
  indices: post?.indices,
  id: post?.id,
  title: post?.title,
  summary: post?.summary,
  contentPath: post?.contentPath,
  submitTime: post?.submitTime,
  avatar: post?.avatar,
  avatarUrl: ''
});

const resolveAvatar = async item => {
  if (!item.avatar) return item;
  try {
    const download = await getDownloadUrl(item.avatar);
    return { ...item, avatarUrl: download.url };
  } catch (error) {
    return item;
  }
};

const extractRecords = payload => {
  if (Array.isArray(payload)) return payload;
  if (Array.isArray(payload?.record)) return payload.record;
  if (Array.isArray(payload?.records)) return payload.records;
  return [];
};

const fetchPosts = async (page = 1, append = false) => {
  if (loadingPosts.value) return;
  loadingPosts.value = true;
  try {
    const begin = (page - 1) * pageSize;
    const requestUrl = new URL(config.forum_fetch_url);
    requestUrl.searchParams.set('range', `${begin}:${begin + pageSize - 1}`);
    const response = await fetch(requestUrl.toString(), { method: 'GET', credentials: 'include' });
    const normalizedRecords = extractRecords(await parseApiResponse(response)).map(normalizePost);
    const records = await Promise.all(normalizedRecords.map(resolveAvatar));
    posts.value = append ? [...posts.value, ...records] : records;
    currentPage.value = page;
    hasMorePosts.value = records.length >= pageSize;
  } catch (error) {
    ElMessage.error(error?.message || '获取论坛帖子失败');
  } finally {
    loadingPosts.value = false;
  }
};

const loadMorePosts = () => fetchPosts(currentPage.value + 1, true);
const getAvatarFallback = id => String(id || '?').trim().slice(0, 1).toUpperCase() || '?';
const formatTime = value => {
  if (!value) return '时间未知';
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? String(value) : date.toLocaleString('zh-CN', { hour12: false });
};

const applyCurrentUser = user => {
  if (!user || typeof user !== 'object') return;
  const id = user.id ?? user.studentId ?? user.student_id;
  const vip = Number(user.vip);
  if (id !== undefined && id !== null) currentUserId.value = String(id);
  if ([config.VIP_TOURIST, config.VIP_NONE, config.VIP_SUPER].includes(vip)) currentUserVip.value = vip;
};

const loadCurrentUser = async () => {
  try {
    applyCurrentUser(JSON.parse(window.localStorage.getItem('user') || 'null'));
  } catch (error) {
    // Ignore malformed legacy login cache and use the server profile below.
  }
  try {
    const response = await fetch(`${config.base_url}/home/studentInfo`, {
      method: 'GET',
      credentials: 'include'
    });
    const result = await response.json();
    if (response.ok && result?.status !== false) applyCurrentUser(result?.data ?? result?.user ?? result);
  } catch (error) {
    // The cached identity is still sufficient for author-owned posts.
  }
};

const canDeletePost = post => Number(currentUserVip.value) === Number(config.VIP_SUPER)
  || (currentUserId.value && String(post?.id) === currentUserId.value);

const removePost = async post => {
  if (!canDeletePost(post) || deletingPostIndex.value !== null) return;
  try {
    await ElMessageBox.confirm(
      `确定删除帖子“${post.title || '未命名帖子'}”吗？删除后无法恢复。`,
      '删除帖子',
      { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' }
    );
  } catch (error) {
    return;
  }

  deletingPostIndex.value = post.indices;
  try {
    const requestUrl = new URL(config.forum_remove_announcement_url);
    requestUrl.searchParams.set('indices', String(post.indices));
    const response = await fetch(requestUrl.toString(), { method: 'GET', credentials: 'include' });
    await parseApiResponse(response);
    posts.value = posts.value.filter(item => item.indices !== post.indices);
    window.sessionStorage.removeItem(`forum-post-${post.indices}`);
    ElMessage.success('帖子已删除');
  } catch (error) {
    ElMessage.error(error?.message || '删除帖子失败');
  } finally {
    deletingPostIndex.value = null;
  }
};

const openPost = post => {
  window.sessionStorage.setItem(`forum-post-${post.indices}`, JSON.stringify(post));
  router.push({ path: `/home/forum/post/${post.indices}`, state: { forumPost: post } });
};

onMounted(() => {
  loadCurrentUser();
  fetchPosts(1);
});
</script>

<style scoped>
.forum-page {
  --forum-surface: rgba(255, 255, 255, 0.94);
  --forum-surface-soft: #f8fafc;
  --forum-border: #e2e8f0;
  --forum-text: #172033;
  --forum-muted: #64748b;
  --forum-primary: #2563eb;
  max-width: 1080px;
  margin: 0 auto;
  color: var(--forum-text);
  text-align: left;
}
.forum-hero {
  position: relative;
  isolation: isolate;
  margin-bottom: 28px;
  padding: 34px 38px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 28px;
  overflow: hidden;
  color: #fff;
  background: linear-gradient(130deg, #1d4ed8 0%, #4338ca 58%, #6d28d9 100%);
  border-radius: 24px;
  box-shadow: 0 20px 48px rgba(37, 99, 235, 0.22);
}
.forum-hero::before,
.forum-hero::after { position: absolute; z-index: -1; border-radius: 50%; content: ''; }
.forum-hero::before { top: -80px; right: 13%; width: 210px; height: 210px; background: rgba(255, 255, 255, 0.1); }
.forum-hero::after { right: -45px; bottom: -105px; width: 245px; height: 245px; border: 42px solid rgba(255, 255, 255, 0.07); }
.hero-copy { min-width: 0; }
.forum-eyebrow { display: block; margin-bottom: 7px; color: #bfdbfe; font-size: 11px; font-weight: 800; letter-spacing: 0.2em; }
.forum-hero h1 { margin: 0 0 8px; font-size: 34px; line-height: 1.15; letter-spacing: -0.02em; }
.forum-hero p { margin: 0; color: #dbeafe; font-size: 14px; }
.primary-button { position: relative; flex: none; padding: 11px 17px 11px 12px; display: flex; align-items: center; gap: 9px; color: #1d4ed8; text-decoration: none; background: #fff; border-radius: 12px; font-weight: 800; box-shadow: 0 9px 22px rgba(15, 23, 42, 0.18); transition: transform 0.2s, box-shadow 0.2s; }
.primary-button:hover { transform: translateY(-2px); box-shadow: 0 13px 26px rgba(15, 23, 42, 0.22); }
.button-icon { width: 27px; height: 27px; display: grid; place-items: center; color: #fff; background: #2563eb; border-radius: 8px; font-size: 20px; line-height: 1; }
.forum-toolbar { margin: 0 2px 14px; display: flex; align-items: flex-end; justify-content: space-between; gap: 16px; }
.section-kicker { color: var(--forum-primary); font-size: 10px; font-weight: 900; letter-spacing: 0.15em; }
.forum-toolbar h2 { margin: 3px 0 0; font-size: 21px; line-height: 1.2; }
.post-count { color: var(--forum-muted); font-size: 12px; }
.post-list { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 15px; }
.post-card { position: relative; box-sizing: border-box; width: 100%; min-height: 224px; padding: 0; overflow: hidden; color: inherit; text-align: left; background: var(--forum-surface); border: 1px solid var(--forum-border); border-radius: 18px; cursor: pointer; box-shadow: 0 8px 24px rgba(15, 23, 42, 0.055); backdrop-filter: blur(16px); transition: transform 0.22s, border-color 0.22s, box-shadow 0.22s; }
.post-card:hover { transform: translateY(-4px); border-color: #60a5fa; box-shadow: 0 16px 34px rgba(37, 99, 235, 0.12); }
.post-card:focus-visible { outline: 3px solid rgba(59, 130, 246, 0.28); outline-offset: 3px; }
.post-card-accent { position: absolute; top: 0; right: 0; left: 0; height: 3px; background: linear-gradient(90deg, #2563eb, #7c3aed); opacity: 0; transition: opacity 0.22s; }
.post-card:hover .post-card-accent { opacity: 1; }
.post-card-content { box-sizing: border-box; min-height: 224px; padding: 20px 21px 17px; display: flex; flex-direction: column; }
.post-meta-row { display: flex; align-items: center; justify-content: space-between; gap: 14px; }
.post-author { min-width: 0; display: flex; align-items: center; gap: 10px; }
.post-avatar { width: 38px; height: 38px; flex: none; display: grid; place-items: center; overflow: hidden; color: #fff; background: linear-gradient(135deg, #3b82f6, #7c3aed); border: 2px solid rgba(255, 255, 255, 0.72); border-radius: 50%; box-shadow: 0 3px 9px rgba(37, 99, 235, 0.2); font-size: 14px; font-weight: 800; }
.post-avatar img { width: 100%; height: 100%; object-fit: cover; }
.author-copy { min-width: 0; display: flex; flex-direction: column; gap: 3px; }
.author-copy strong { overflow: hidden; font-size: 13px; text-overflow: ellipsis; white-space: nowrap; }
.author-copy time { color: var(--forum-muted); font-size: 11px; }
.post-index { flex: none; padding: 6px 8px; color: var(--forum-primary); background: rgba(37, 99, 235, 0.09); border-radius: 8px; font-size: 11px; font-weight: 800; }
.post-copy { min-width: 0; margin: 18px 0 16px; display: flex; flex: 1; flex-direction: column; gap: 8px; }
.post-title { overflow: hidden; font-size: 18px; line-height: 1.42; text-overflow: ellipsis; white-space: nowrap; }
.post-summary { display: -webkit-box; overflow: hidden; color: var(--forum-muted); font-size: 13px; line-height: 1.72; overflow-wrap: anywhere; -webkit-box-orient: vertical; -webkit-line-clamp: 2; }
.post-summary.is-empty { opacity: 0.7; font-style: italic; }
.post-footer { padding-top: 13px; display: flex; align-items: center; justify-content: space-between; color: var(--forum-primary); border-top: 1px solid var(--forum-border); font-size: 12px; font-weight: 800; }
.read-more { margin-left: auto; display: inline-flex; align-items: center; gap: 8px; }
.delete-post-button { padding: 6px 9px; color: #dc2626; background: rgba(239, 68, 68, 0.08); border: 1px solid rgba(239, 68, 68, 0.18); border-radius: 8px; cursor: pointer; font: inherit; transition: background 0.2s, border-color 0.2s; }
.delete-post-button:hover:not(:disabled) { background: rgba(239, 68, 68, 0.14); border-color: rgba(239, 68, 68, 0.34); }
.delete-post-button:disabled { cursor: wait; opacity: 0.58; }
.post-arrow { font-size: 18px; transition: transform 0.2s; }
.post-card:hover .post-arrow { transform: translateX(4px); }
.state-card { padding: 62px 24px; display: flex; flex-direction: column; align-items: center; gap: 8px; color: var(--forum-muted); text-align: center; background: var(--forum-surface); border: 1px dashed var(--forum-border); border-radius: 18px; backdrop-filter: blur(14px); }
.state-card strong { color: var(--forum-text); font-size: 17px; }
.load-more { margin: 25px auto 0; padding: 11px 24px; display: block; color: var(--forum-primary); background: var(--forum-surface); border: 1px solid var(--forum-border); border-radius: 11px; cursor: pointer; font-weight: 800; box-shadow: 0 5px 15px rgba(15, 23, 42, 0.05); transition: border-color 0.2s, transform 0.2s; }
.load-more:hover:not(:disabled) { transform: translateY(-1px); border-color: #60a5fa; }
.load-more:disabled { cursor: wait; opacity: 0.6; }
@media (max-width: 760px) {
  .forum-hero { padding: 25px 23px; align-items: flex-start; flex-direction: column; border-radius: 19px; }
  .forum-hero h1 { font-size: 29px; }
  .post-list { grid-template-columns: 1fr; }
  .post-card, .post-card-content { min-height: 0; }
}
@media (max-width: 430px) {
  .forum-toolbar { align-items: flex-start; flex-direction: column; gap: 5px; }
  .post-card-content { padding: 17px; }
  .author-copy time { font-size: 10px; }
}
</style>

<style>
html[data-theme='dark'] .forum-page { --forum-surface: rgba(15,23,42,.94); --forum-surface-soft: rgba(30,41,59,.72); --forum-border: rgba(148,163,184,.28); --forum-text: #e2e8f0; --forum-muted: #94a3b8; --forum-primary: #60a5fa; }
html[data-theme='effect'] .forum-page { --forum-surface: rgba(15,23,42,.66); --forum-surface-soft: rgba(30,41,59,.52); --forum-border: rgba(148,163,184,.3); --forum-text: #e2e8f0; --forum-muted: #a8b5c7; --forum-primary: #7dd3fc; }
html[data-theme='custom'] .forum-page { --forum-surface: rgba(255,255,255,.88); --forum-surface-soft: rgba(248,250,252,.8); }
</style>
