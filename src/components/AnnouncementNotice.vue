<template>
  <transition name="announcement-pop">
    <div
      v-if="announcement"
      class="announcement-overlay"
      role="presentation"
    >
      <article
        class="announcement-card"
        role="dialog"
        aria-modal="true"
        aria-labelledby="global-announcement-title"
      >
        <div class="announcement-card__accent" aria-hidden="true"></div>
        <header class="announcement-card__header">
          <div class="announcement-card__decoration announcement-card__decoration--one" aria-hidden="true"></div>
          <div class="announcement-card__decoration announcement-card__decoration--two" aria-hidden="true"></div>

          <div class="announcement-card__heading">
            <span class="announcement-card__icon" aria-hidden="true">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                <path stroke-linecap="round" stroke-linejoin="round" d="M10.34 3.94A2 2 0 0112 3a2 2 0 011.66.94 7 7 0 014.34 6.48V14l1.4 2.1A1.2 1.2 0 0118.4 18H5.6a1.2 1.2 0 01-1-1.9L6 14v-3.58a7 7 0 014.34-6.48zM9.75 21h4.5" />
              </svg>
            </span>
            <div class="min-w-0">
              <div class="announcement-card__eyebrow">
                <span class="announcement-card__eyebrow-dot"></span>
                系统公告
              </div>
              <h2 id="global-announcement-title" class="announcement-card__title">
                {{ announcement.title || '最新公告' }}
              </h2>
            </div>
          </div>

          <div class="announcement-card__meta">
            <div class="announcement-card__meta-item">
              <span class="announcement-card__meta-label">发布于</span>
              <span>{{ announcement.createdAt }}</span>
            </div>
            <div class="announcement-card__meta-divider" aria-hidden="true"></div>
            <div class="announcement-card__meta-item">
              <span class="announcement-card__meta-label">更新于</span>
              <span>{{ announcement.updatedAt }}</span>
            </div>
          </div>
        </header>

        <div class="announcement-card__body">
          <div class="announcement-card__content">
            <p>{{ announcement.content }}</p>
          </div>
        </div>

        <footer class="announcement-card__footer">
          <p class="announcement-card__hint">
            <svg viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
              <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7.25-3.25a.75.75 0 00-1.5 0v3.5c0 .414.336.75.75.75h2.25a.75.75 0 000-1.5h-1.5V6.75zm.125 7a.875.875 0 11-1.75 0 .875.875 0 011.75 0z" clip-rule="evenodd" />
            </svg>
            阅读后，本版本公告将不再提示
          </p>
          <button
            type="button"
            class="announcement-card__button"
            @click="acknowledgeAnnouncement"
          >
            <svg viewBox="0 0 20 20" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" d="M4 10.5l3.5 3.5L16 5.5" />
            </svg>
            我已知晓
          </button>
        </footer>
      </article>
    </div>
  </transition>
</template>

<script setup>
import { onBeforeUnmount, ref, watch } from 'vue';
import { useRoute } from 'vue-router';
import config from '../config';

const route = useRoute();
const announcement = ref(null);
let requestController = null;
let requestSequence = 0;

const isAnnouncementEnabled = value => (
  value !== false
  && value !== 0
  && String(value ?? '').toLowerCase() !== 'false'
  && String(value ?? '') !== '0'
);

const getAnnouncementPayload = result => {
  let payload = result?.data ?? result;
  if (Array.isArray(payload)) payload = payload[0];
  return payload && typeof payload === 'object' ? payload : null;
};

const formatAnnouncementDate = value => {
  if (value === null || value === undefined || value === '') return '未知';

  const numericValue = Number(value);
  const date = typeof value === 'number' || /^\d+$/.test(String(value).trim())
    ? new Date(numericValue < 100000000000 ? numericValue * 1000 : numericValue)
    : new Date(value);
  if (Number.isNaN(date.getTime())) return String(value);
  return date.toLocaleString('zh-CN', { hour12: false });
};

const fetchAnnouncement = async () => {
  requestController?.abort();
  announcement.value = null;
  const currentSequence = ++requestSequence;

  if (route.meta.hideAnnouncement) return;

  requestController = new AbortController();

  try {
    const response = await fetch(config.fetch_announcement_url, {
      method: 'GET',
      credentials: 'include',
      signal: requestController.signal
    });
    const result = await response.json();
    if (!response.ok || result?.status === false) {
      throw new Error(result?.msg || `获取公告失败（HTTP ${response.status}）`);
    }

    const payload = getAnnouncementPayload(result);
    if (
      currentSequence !== requestSequence
      || !payload
      || !isAnnouncementEnabled(payload.enabled)
    ) return;

    const indices = String(payload.indices ?? '').trim();
    const version = String(payload.version ?? '').trim();
    if (!indices || !version) return;

    const storageKey = `Announcement_${indices}_${version}`;
    try {
      if (window.localStorage.getItem(storageKey) !== null) return;
    } catch (error) {
      // 本地存储不可用时仍显示公告，避免用户错过重要信息。
    }

    announcement.value = {
      storageKey,
      title: String(payload.title ?? '').trim(),
      content: String(payload.content ?? ''),
      createdAt: formatAnnouncementDate(
        payload.create_time ?? payload.createTime ?? payload.created_at
      ),
      updatedAt: formatAnnouncementDate(
        payload.update_time ?? payload.updateTime ?? payload.updated_at
      )
    };
  } catch (error) {
    if (error?.name !== 'AbortError') {
      // 公告获取失败不阻断当前页面的其它功能。
      console.error('获取公告失败:', error);
    }
  }
};

const acknowledgeAnnouncement = () => {
  const storageKey = announcement.value?.storageKey;
  if (storageKey) {
    try {
      window.localStorage.setItem(storageKey, '1');
    } catch (error) {
      // 本地存储不可用时仍允许关闭当前公告。
    }
  }
  announcement.value = null;
};

watch(() => route.fullPath, fetchAnnouncement, { immediate: true });

onBeforeUnmount(() => {
  requestController?.abort();
});
</script>

<style>
.announcement-overlay {
  position: fixed;
  inset: 0;
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1.25rem;
  background: rgba(15, 23, 42, 0.55);
  backdrop-filter: blur(8px);
}

.announcement-card {
  position: relative;
  width: 100%;
  max-width: 36rem;
  overflow: hidden;
  border: 1px solid rgba(148, 163, 184, 0.22);
  border-radius: 1.5rem;
  color: #1e293b;
  background: #ffffff;
  box-shadow: 0 30px 80px rgba(15, 23, 42, 0.28), 0 8px 24px rgba(15, 23, 42, 0.14);
  text-align: left;
}

.announcement-card__accent {
  height: 0.3rem;
  background: linear-gradient(90deg, #2563eb, #4f46e5 52%, #7c3aed);
}

.announcement-card__header {
  position: relative;
  overflow: hidden;
  padding: 1.75rem 1.75rem 1.35rem;
  border-bottom: 1px solid #e2e8f0;
  background: linear-gradient(145deg, #eff6ff 0%, #eef2ff 58%, #faf5ff 100%);
}

.announcement-card__decoration {
  position: absolute;
  border-radius: 9999px;
  pointer-events: none;
}

.announcement-card__decoration--one {
  top: -5rem;
  right: -3rem;
  width: 11rem;
  height: 11rem;
  background: rgba(99, 102, 241, 0.1);
}

.announcement-card__decoration--two {
  right: 5rem;
  bottom: -4rem;
  width: 7rem;
  height: 7rem;
  background: rgba(59, 130, 246, 0.08);
}

.announcement-card__heading {
  position: relative;
  display: flex;
  align-items: flex-start;
  gap: 1rem;
}

.announcement-card__icon {
  display: inline-flex;
  width: 3rem;
  height: 3rem;
  flex: 0 0 auto;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(99, 102, 241, 0.18);
  border-radius: 1rem;
  color: #ffffff;
  background: linear-gradient(145deg, #3b82f6, #4f46e5);
  box-shadow: 0 10px 22px rgba(79, 70, 229, 0.25);
}

.announcement-card__icon svg {
  width: 1.55rem;
  height: 1.55rem;
}

.announcement-card__eyebrow {
  display: flex;
  align-items: center;
  gap: 0.45rem;
  margin-bottom: 0.35rem;
  color: #4f46e5;
  font-size: 0.75rem;
  font-weight: 700;
  letter-spacing: 0.12em;
}

.announcement-card__eyebrow-dot {
  width: 0.45rem;
  height: 0.45rem;
  border-radius: 9999px;
  background: #6366f1;
  box-shadow: 0 0 0 0.25rem rgba(99, 102, 241, 0.12);
}

.announcement-card__title {
  margin: 0;
  overflow-wrap: anywhere;
  color: #0f172a;
  font-size: 1.45rem;
  font-weight: 750;
  line-height: 1.35;
}

.announcement-card__meta {
  position: relative;
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-top: 1.25rem;
  padding: 0.8rem 1rem;
  border: 1px solid rgba(148, 163, 184, 0.2);
  border-radius: 0.85rem;
  color: #475569;
  background: rgba(255, 255, 255, 0.65);
  font-size: 0.75rem;
}

.announcement-card__meta-item {
  display: flex;
  min-width: 0;
  flex: 1;
  flex-direction: column;
  gap: 0.15rem;
}

.announcement-card__meta-label {
  color: #94a3b8;
  font-size: 0.68rem;
  font-weight: 600;
}

.announcement-card__meta-divider {
  width: 1px;
  height: 1.8rem;
  background: #dbe4f0;
}

.announcement-card__body {
  max-height: min(46vh, 25rem);
  overflow-y: auto;
  padding: 1.5rem 1.75rem;
  background: #ffffff;
}

.announcement-card__content {
  position: relative;
  padding: 1.15rem 1.25rem 1.15rem 1.45rem;
  border: 1px solid #e2e8f0;
  border-radius: 1rem;
  color: #334155;
  background: #f8fafc;
  font-size: 0.95rem;
  line-height: 1.85;
}

.announcement-card__content::before {
  position: absolute;
  top: 1rem;
  bottom: 1rem;
  left: 0;
  width: 0.25rem;
  border-radius: 0 9999px 9999px 0;
  background: linear-gradient(180deg, #3b82f6, #8b5cf6);
  content: '';
}

.announcement-card__content p {
  margin: 0;
  white-space: pre-wrap;
  overflow-wrap: anywhere;
}

.announcement-card__footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 1rem 1.75rem 1.25rem;
  border-top: 1px solid #eef2f7;
  background: #ffffff;
}

.announcement-card__hint {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  margin: 0;
  color: #94a3b8;
  font-size: 0.72rem;
}

.announcement-card__hint svg {
  width: 1rem;
  height: 1rem;
  flex: 0 0 auto;
}

.announcement-card__button {
  display: inline-flex;
  min-width: 7.5rem;
  align-items: center;
  justify-content: center;
  gap: 0.45rem;
  border: 0;
  border-radius: 0.75rem;
  padding: 0.72rem 1.25rem;
  color: #ffffff;
  background: linear-gradient(135deg, #2563eb, #4f46e5);
  box-shadow: 0 8px 18px rgba(37, 99, 235, 0.24);
  cursor: pointer;
  font-size: 0.875rem;
  font-weight: 650;
  transition: transform 0.2s ease, box-shadow 0.2s ease, filter 0.2s ease;
}

.announcement-card__button:hover {
  filter: brightness(1.06);
  box-shadow: 0 11px 24px rgba(37, 99, 235, 0.3);
  transform: translateY(-1px);
}

.announcement-card__button:active {
  transform: translateY(0);
}

.announcement-card__button svg {
  width: 1.05rem;
  height: 1.05rem;
}

.announcement-pop-enter-active,
.announcement-pop-leave-active {
  transition: opacity 0.22s ease;
}

.announcement-pop-enter-active .announcement-card,
.announcement-pop-leave-active .announcement-card {
  transition: transform 0.22s ease, opacity 0.22s ease;
}

.announcement-pop-enter-from,
.announcement-pop-leave-to {
  opacity: 0;
}

.announcement-pop-enter-from .announcement-card,
.announcement-pop-leave-to .announcement-card {
  opacity: 0;
  transform: translateY(0.75rem) scale(0.975);
}

html[data-theme='dark'] .announcement-card,
html[data-theme='effect'] .announcement-card {
  border-color: rgba(148, 163, 184, 0.28);
  color: #e2e8f0;
  background: #0f172a;
}

html[data-theme='dark'] .announcement-card__header,
html[data-theme='effect'] .announcement-card__header {
  border-color: rgba(148, 163, 184, 0.2);
  background: linear-gradient(145deg, #172554 0%, #1e1b4b 58%, #2e1065 100%);
}

html[data-theme='dark'] .announcement-card__title,
html[data-theme='effect'] .announcement-card__title {
  color: #f8fafc;
}

html[data-theme='dark'] .announcement-card__meta,
html[data-theme='effect'] .announcement-card__meta {
  border-color: rgba(148, 163, 184, 0.24);
  color: #cbd5e1;
  background: rgba(15, 23, 42, 0.5);
}

html[data-theme='dark'] .announcement-card__meta-divider,
html[data-theme='effect'] .announcement-card__meta-divider {
  background: rgba(148, 163, 184, 0.28);
}

html[data-theme='dark'] .announcement-card__body,
html[data-theme='dark'] .announcement-card__footer,
html[data-theme='effect'] .announcement-card__body,
html[data-theme='effect'] .announcement-card__footer {
  border-color: rgba(148, 163, 184, 0.18);
  background: #0f172a;
}

html[data-theme='dark'] .announcement-card__content,
html[data-theme='effect'] .announcement-card__content {
  border-color: rgba(148, 163, 184, 0.25);
  color: #dbeafe;
  background: rgba(30, 41, 59, 0.72);
}

@media (max-width: 640px) {
  .announcement-overlay {
    align-items: flex-end;
    padding: 0.75rem;
  }

  .announcement-card {
    border-radius: 1.25rem;
  }

  .announcement-card__header,
  .announcement-card__body {
    padding-right: 1.15rem;
    padding-left: 1.15rem;
  }

  .announcement-card__meta {
    align-items: stretch;
    flex-direction: column;
    gap: 0.6rem;
  }

  .announcement-card__meta-divider {
    width: 100%;
    height: 1px;
  }

  .announcement-card__footer {
    align-items: stretch;
    flex-direction: column;
    padding-right: 1.15rem;
    padding-left: 1.15rem;
  }

  .announcement-card__hint {
    justify-content: center;
  }

  .announcement-card__button {
    width: 100%;
  }
}
</style>
