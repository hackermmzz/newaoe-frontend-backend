<template>
  <section class="forum-editor-page">
    <header class="editor-page-header">
      <button class="back-button" type="button" @click="goBack">← 返回论坛</button>
      <div class="header-copy">
        <span>NEW POST</span>
        <h1>发布帖子</h1>
        <p>像编辑文档一样编排正文，图片、视频和文件会出现在放置的位置。</p>
      </div>
      <button class="publish-button header-publish" type="button" :disabled="submitting" @click="publishPost">
        {{ submitting ? '发布中…' : '发布帖子' }}
      </button>
    </header>

    <main class="editor-sheet">
      <label class="field-label" for="forum-title">帖子标题</label>
      <input
        id="forum-title"
        v-model.trim="draftTitle"
        class="title-input"
        type="text"
        maxlength="120"
        placeholder="请输入帖子标题"
      >

      <div class="content-label-row">
        <span class="field-label">正文内容</span>
        <span class="content-tip">支持直接粘贴或拖入图片、视频和文件</span>
      </div>
      <div
        ref="editorRef"
        :class="['document-editor', { 'is-dragging-file': draggingOverEditor }]"
        contenteditable="true"
        data-placeholder="从这里开始写正文…"
        @input="handleEditorInput"
        @keydown="handleEditorKeydown"
        @keyup="rememberSelection"
        @mouseup="rememberSelection"
        @click="handleEditorClick"
        @paste="handleEditorPaste"
        @dragenter.prevent="draggingOverEditor = true"
        @dragover.prevent="handleEditorDragOver"
        @dragleave="handleEditorDragLeave"
        @drop.prevent="handleEditorDrop"
      ></div>
      <div class="editor-note">
        <span aria-hidden="true">ⓘ</span>
        资源会插入光标或拖放位置；连续选择的图片和视频会并排显示，文件独占一行。本地会自动生成唯一文件名。
      </div>
    </main>

    <footer class="editor-page-actions">
      <button class="cancel-button" type="button" :disabled="submitting" @click="goBack">取消</button>
      <button class="publish-button" type="button" :disabled="submitting" @click="publishPost">
        {{ submitting ? '正在上传并发布…' : '发布帖子' }}
      </button>
    </footer>
  </section>
</template>

<script>
export default {
  name: 'ForumEditorPage'
};
</script>

<script setup>
import { nextTick, onBeforeUnmount, onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import config from '../config';
import { uploadFiles } from '../utils/upload';

const router = useRouter();
const draftTitle = ref('');
const editorRef = ref(null);
const draggingOverEditor = ref(false);
const submitting = ref(false);
const attachments = ref([]);

let savedRange = null;
let attachmentCounter = 0;
let editorHistory = [];
let editorHistoryIndex = -1;
let historyTimer = null;

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
  if (!response.ok || result?.status === false) {
    throw new Error(result?.msg || `请求失败（HTTP ${response.status}）`);
  }
  return result?.data ?? result;
};

const goBack = () => {
  if (submitting.value) return;
  router.push('/home/forum');
};

const rememberSelection = () => {
  const selection = window.getSelection();
  if (!selection?.rangeCount || !editorRef.value) return;
  const range = selection.getRangeAt(0);
  if (editorRef.value.contains(range.commonAncestorContainer)) {
    savedRange = range.cloneRange();
  }
};

const commitEditorHistory = () => {
  if (!editorRef.value) return;
  if (historyTimer) {
    window.clearTimeout(historyTimer);
    historyTimer = null;
  }
  const snapshot = editorRef.value.innerHTML;
  if (editorHistory[editorHistoryIndex] === snapshot) return;
  editorHistory = editorHistory.slice(0, editorHistoryIndex + 1);
  editorHistory.push(snapshot);
  if (editorHistory.length > 100) editorHistory.shift();
  editorHistoryIndex = editorHistory.length - 1;
};

const queueEditorHistory = () => {
  if (historyTimer) window.clearTimeout(historyTimer);
  historyTimer = window.setTimeout(commitEditorHistory, 300);
};

const restoreEditorHistory = index => {
  const editor = editorRef.value;
  if (!editor || index < 0 || index >= editorHistory.length) return;
  editorHistoryIndex = index;
  editor.innerHTML = editorHistory[index];
  editor.focus();
  const range = document.createRange();
  range.selectNodeContents(editor);
  range.collapse(false);
  const selection = window.getSelection();
  selection.removeAllRanges();
  selection.addRange(range);
  savedRange = range.cloneRange();
};

const handleEditorInput = () => {
  rememberSelection();
  queueEditorHistory();
};

const handleEditorKeydown = event => {
  if (!(event.ctrlKey || event.metaKey) || event.altKey) return;
  const key = event.key.toLowerCase();
  const undo = key === 'z' && !event.shiftKey;
  const redo = key === 'y' || (key === 'z' && event.shiftKey);
  if (!undo && !redo) return;
  event.preventDefault();
  commitEditorHistory();
  if (undo && editorHistoryIndex > 0) restoreEditorHistory(editorHistoryIndex - 1);
  if (redo && editorHistoryIndex < editorHistory.length - 1) restoreEditorHistory(editorHistoryIndex + 1);
};

const getExtension = fileName => {
  const match = String(fileName || '').match(/(\.[a-zA-Z0-9]{1,12})$/);
  return match ? match[1].toLowerCase() : '';
};

const createUniqueUploadFile = file => {
  attachmentCounter += 1;
  const storedName = `forum-${Date.now()}-${attachmentCounter}${getExtension(file.name)}`;
  return {
    storedName,
    uploadFile: new File([file], storedName, {
      type: file.type,
      lastModified: file.lastModified
    })
  };
};

const placeNodeAtSelection = node => {
  const editor = editorRef.value;
  if (!editor) return;
  editor.focus();
  const range = savedRange && editor.contains(savedRange.commonAncestorContainer)
    ? savedRange
    : document.createRange();
  if (!savedRange || !editor.contains(range.commonAncestorContainer)) {
    range.selectNodeContents(editor);
    range.collapse(false);
  }
  range.deleteContents();
  range.insertNode(node);
  const spacer = document.createElement('p');
  spacer.appendChild(document.createElement('br'));
  node.after(spacer);
  range.setStart(spacer, 0);
  range.collapse(true);
  const selection = window.getSelection();
  selection.removeAllRanges();
  selection.addRange(range);
  savedRange = range.cloneRange();
};

const createAttachmentNode = attachment => {
  const wrapper = document.createElement(attachment.kind === 'file' ? 'div' : 'figure');
  wrapper.className = attachment.kind === 'file' ? 'forum-file-card' : 'forum-embed';
  wrapper.dataset.attachmentId = attachment.id;
  wrapper.contentEditable = 'false';

  if (attachment.kind === 'image') {
    const image = document.createElement('img');
    image.src = attachment.previewUrl;
    image.alt = attachment.originalName;
    image.dataset.forumAsset = attachment.storedName;
    image.dataset.assetKind = 'image';
    wrapper.appendChild(image);
  } else if (attachment.kind === 'video') {
    const video = document.createElement('video');
    video.src = attachment.previewUrl;
    video.controls = true;
    video.preload = 'metadata';
    video.dataset.forumAsset = attachment.storedName;
    video.dataset.assetKind = 'video';
    wrapper.appendChild(video);
  } else {
    const link = document.createElement('a');
    link.href = attachment.storedName;
    link.dataset.forumAsset = attachment.storedName;
    link.dataset.assetKind = 'file';
    link.textContent = `📎 ${attachment.originalName}`;
    wrapper.appendChild(link);
  }

  if (attachment.kind !== 'file') {
    const caption = document.createElement('figcaption');
    caption.textContent = attachment.originalName;
    wrapper.appendChild(caption);
  }
  const removeButton = document.createElement('button');
  removeButton.type = 'button';
  removeButton.dataset.removeAttachment = attachment.id;
  removeButton.textContent = '移除';
  wrapper.appendChild(removeButton);
  return wrapper;
};

const createMediaRow = mediaAttachments => {
  const row = document.createElement('div');
  row.className = 'forum-media-row';
  row.contentEditable = 'false';
  mediaAttachments.forEach(attachment => row.appendChild(createAttachmentNode(attachment)));
  return row;
};

const getMediaRowAtSelection = () => {
  const editor = editorRef.value;
  if (!editor || !savedRange) return null;
  let selectionNode = savedRange.startContainer;
  if (selectionNode.nodeType !== Node.ELEMENT_NODE) selectionNode = selectionNode.parentElement;
  if (!selectionNode || !editor.contains(selectionNode)) return null;

  const containingRow = selectionNode.closest?.('.forum-media-row');
  if (containingRow) return containingRow;

  const currentBlock = selectionNode.closest?.('p, div');
  if (
    currentBlock
    && currentBlock !== editor
    && !currentBlock.textContent.trim()
    && currentBlock.previousElementSibling?.classList.contains('forum-media-row')
  ) {
    return currentBlock.previousElementSibling;
  }

  if (selectionNode === editor && savedRange.startOffset > 0) {
    const previousNode = editor.childNodes[savedRange.startOffset - 1];
    if (previousNode?.nodeType === Node.ELEMENT_NODE && previousNode.classList.contains('forum-media-row')) {
      return previousNode;
    }
  }
  return null;
};

const getAttachmentKind = file => {
  if (file.type.startsWith('image/')) return 'image';
  if (file.type.startsWith('video/')) return 'video';
  return 'file';
};

const insertFilesIntoEditor = files => {
  if (!files?.length) return;
  commitEditorHistory();
  const newAttachments = Array.from(files || []).map(file => {
    const kind = getAttachmentKind(file);
    const { storedName, uploadFile } = createUniqueUploadFile(file);
    const attachment = {
      id: `${Date.now()}-${attachmentCounter}`,
      kind,
      originalName: file.name,
      storedName,
      uploadFile,
      previewUrl: URL.createObjectURL(file)
    };
    attachments.value.push(attachment);
    return attachment;
  });

  let mediaGroup = [];
  const insertMediaGroup = () => {
    if (!mediaGroup.length) return;
    const existingRow = getMediaRowAtSelection();
    if (existingRow) {
      mediaGroup.forEach(attachment => existingRow.appendChild(createAttachmentNode(attachment)));
    } else {
      placeNodeAtSelection(createMediaRow(mediaGroup));
    }
    mediaGroup = [];
  };

  newAttachments.forEach(attachment => {
    if (attachment.kind !== 'file') {
      mediaGroup.push(attachment);
      return;
    }
    insertMediaGroup();
    placeNodeAtSelection(createAttachmentNode(attachment));
  });
  insertMediaGroup();
  commitEditorHistory();
};

const handleEditorPaste = event => {
  const files = Array.from(event.clipboardData?.files || []);
  if (!files.length) return;
  event.preventDefault();
  rememberSelection();
  insertFilesIntoEditor(files);
};

const setSelectionFromPoint = (clientX, clientY) => {
  let range = null;
  if (document.caretRangeFromPoint) {
    range = document.caretRangeFromPoint(clientX, clientY);
  } else if (document.caretPositionFromPoint) {
    const position = document.caretPositionFromPoint(clientX, clientY);
    if (position) {
      range = document.createRange();
      range.setStart(position.offsetNode, position.offset);
      range.collapse(true);
    }
  }
  if (range && editorRef.value?.contains(range.commonAncestorContainer)) {
    savedRange = range.cloneRange();
  }
};

const handleEditorDragOver = event => {
  draggingOverEditor.value = true;
  setSelectionFromPoint(event.clientX, event.clientY);
};

const handleEditorDragLeave = event => {
  if (!editorRef.value?.contains(event.relatedTarget)) draggingOverEditor.value = false;
};

const handleEditorDrop = event => {
  draggingOverEditor.value = false;
  const files = Array.from(event.dataTransfer?.files || []);
  if (!files.length) return;
  setSelectionFromPoint(event.clientX, event.clientY);
  insertFilesIntoEditor(files);
};

const removeAttachment = id => {
  if (!attachments.value.some(item => item.id === id)) return;
  commitEditorHistory();
  const attachmentNode = editorRef.value?.querySelector(`[data-attachment-id="${id}"]`);
  const mediaRow = attachmentNode?.closest('.forum-media-row');
  attachmentNode?.remove();
  if (mediaRow && !mediaRow.querySelector('[data-attachment-id]')) mediaRow.remove();
  commitEditorHistory();
};

const handleEditorClick = event => {
  const removeButton = event.target.closest('[data-remove-attachment]');
  if (removeButton) {
    removeAttachment(removeButton.dataset.removeAttachment);
    return;
  }
  if (event.target.closest('a[data-forum-asset]')) event.preventDefault();
};

const buildPostHtml = activeAttachments => {
  const clone = editorRef.value.cloneNode(true);
  clone.querySelectorAll('[data-remove-attachment]').forEach(button => button.remove());
  clone.querySelectorAll('[contenteditable]').forEach(element => element.removeAttribute('contenteditable'));
  clone.querySelectorAll('[data-attachment-id]').forEach(element => element.removeAttribute('data-attachment-id'));
  clone.querySelectorAll('[data-forum-asset]').forEach(element => {
    const attachment = activeAttachments.find(item => item.storedName === element.dataset.forumAsset);
    if (!attachment) {
      element.closest('.forum-embed, .forum-file-card')?.remove();
      return;
    }
    element.setAttribute(element.tagName === 'A' ? 'href' : 'src', attachment.storedName);
  });
  return `<!doctype html>
<html lang="zh-CN">
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"></head>
<body><article class="forum-document">${clone.innerHTML}</article></body>
</html>`;
};

const getPolicies = (payload, key) => Array.isArray(payload?.[key]) ? payload[key] : [];

const publishPost = async () => {
  if (!draftTitle.value) {
    ElMessage.warning('请输入帖子标题');
    return;
  }
  const textContent = editorRef.value?.innerText?.trim() || '';
  const activeNames = new Set(
    Array.from(editorRef.value?.querySelectorAll('[data-forum-asset]') || [])
      .map(element => element.dataset.forumAsset)
  );
  const activeAttachments = attachments.value.filter(item => activeNames.has(item.storedName));
  if (!textContent && !activeAttachments.length) {
    ElMessage.warning('请输入帖子内容');
    return;
  }

  const images = activeAttachments.filter(item => item.kind === 'image');
  const videos = activeAttachments.filter(item => item.kind === 'video');
  const files = activeAttachments.filter(item => item.kind === 'file');
  submitting.value = true;
  try {
    const response = await fetch(config.forum_announcement_url, {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        content: buildPostHtml(activeAttachments),
        title: draftTitle.value,
        files: files.map(item => item.storedName),
        videos: videos.map(item => item.storedName),
        images: images.map(item => item.storedName)
      })
    });
    const payload = await parseApiResponse(response);
    const filePolicies = getPolicies(payload, 'fileUrls');
    const videoPolicies = getPolicies(payload, 'videoUrls');
    const imagePolicies = getPolicies(payload, 'imageUrls');
    if (!payload?.key) throw new Error('服务器未返回帖子 key');
    if (
      filePolicies.length !== files.length
      || videoPolicies.length !== videos.length
      || imagePolicies.length !== images.length
    ) {
      throw new Error('服务器返回的附件上传链接数量不正确');
    }
    await uploadFiles(
      [...files, ...videos, ...images].map(item => item.uploadFile),
      [...filePolicies, ...videoPolicies, ...imagePolicies]
    );
    const ackResponse = await fetch(config.forum_announcement_ack_url, {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ key: payload.key })
    });
    await parseApiResponse(ackResponse);
    ElMessage.success('帖子发布成功');
    await router.push('/home/forum');
  } catch (error) {
    ElMessage.error(error?.message || '帖子发布失败');
  } finally {
    submitting.value = false;
  }
};

onMounted(() => nextTick(() => {
  editorRef.value?.focus();
  commitEditorHistory();
}));
onBeforeUnmount(() => {
  if (historyTimer) window.clearTimeout(historyTimer);
  attachments.value.forEach(item => URL.revokeObjectURL(item.previewUrl));
});
</script>

<style scoped>
.forum-editor-page {
  --editor-surface: #ffffff;
  --editor-surface-soft: #f8fafc;
  --editor-border: #e2e8f0;
  --editor-text: #1e293b;
  --editor-muted: #64748b;
  max-width: 1100px;
  margin: 0 auto;
  color: var(--editor-text);
  text-align: left;
}
.editor-page-header { margin-bottom: 20px; padding: 24px 28px; display: grid; grid-template-columns: auto minmax(0, 1fr) auto; align-items: center; gap: 22px; color: #fff; background: linear-gradient(135deg, #1d4ed8, #4338ca); border-radius: 18px; box-shadow: 0 16px 38px rgba(37, 99, 235, 0.2); }
.header-copy span { display: block; margin-bottom: 4px; color: #bfdbfe; font-size: 11px; font-weight: 800; letter-spacing: 0.16em; }
.header-copy h1 { margin: 0 0 5px; font-size: 27px; }
.header-copy p { margin: 0; color: #dbeafe; font-size: 13px; }
.back-button { padding: 9px 12px; color: #dbeafe; background: rgba(255, 255, 255, 0.1); border: 1px solid rgba(255, 255, 255, 0.2); border-radius: 9px; cursor: pointer; }
.publish-button, .cancel-button { padding: 11px 20px; border: 0; border-radius: 10px; cursor: pointer; font-weight: 800; }
.publish-button { color: #fff; background: #2563eb; box-shadow: 0 8px 20px rgba(37, 99, 235, 0.24); }
.header-publish { color: #1d4ed8; background: #fff; }
.publish-button:disabled, .cancel-button:disabled { cursor: not-allowed; opacity: 0.55; }
.editor-sheet { padding: 30px; background: var(--editor-surface); border: 1px solid var(--editor-border); border-radius: 18px; box-shadow: 0 12px 34px rgba(15, 23, 42, 0.07); backdrop-filter: blur(16px); }
.field-label { display: block; margin-bottom: 9px; color: var(--editor-text); font-size: 13px; font-weight: 800; }
.title-input { box-sizing: border-box; width: 100%; margin-bottom: 26px; padding: 15px 2px; color: var(--editor-text); background: transparent; border: 0; border-bottom: 1px solid var(--editor-border); outline: none; font-size: 25px; font-weight: 700; }
.title-input::placeholder { color: #94a3b8; }
.content-label-row { display: flex; align-items: center; justify-content: space-between; gap: 16px; }
.content-tip { color: var(--editor-muted); font-size: 12px; }
.document-editor { box-sizing: border-box; min-height: 520px; padding: 24px; color: var(--editor-text); overflow-wrap: anywhere; background: var(--editor-surface-soft); border: 1px solid var(--editor-border); border-radius: 13px; outline: none; font-size: 15px; line-height: 1.8; }
.document-editor:empty::before { color: #94a3b8; content: attr(data-placeholder); pointer-events: none; }
.document-editor:focus { border-color: #60a5fa; box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.12); }
.document-editor.is-dragging-file { background: rgba(37, 99, 235, 0.1); border-color: #3b82f6; box-shadow: inset 0 0 0 2px rgba(37, 99, 235, 0.12); }
.document-editor :deep(.forum-media-row) { margin: 20px 0; display: grid; grid-auto-columns: minmax(0, 1fr); grid-auto-flow: column; gap: 8px; }
.document-editor :deep(.forum-embed) { position: relative; min-width: 0; margin: 0; padding: 7px; overflow: hidden; background: var(--editor-surface); border: 1px solid var(--editor-border); border-radius: 12px; }
.document-editor :deep(.forum-embed img), .document-editor :deep(.forum-embed video) { display: block; width: 100%; max-height: 600px; object-fit: contain; background: #0f172a; border-radius: 8px; }
.document-editor :deep(.forum-embed:not(:only-child) img), .document-editor :deep(.forum-embed:not(:only-child) video) { aspect-ratio: 1 / 1; object-fit: cover; }
.document-editor :deep(figcaption) { margin-top: 6px; overflow: hidden; color: var(--editor-muted); text-align: center; text-overflow: ellipsis; white-space: nowrap; font-size: 11px; }
.document-editor :deep(.forum-file-card) { position: relative; margin: 16px 0; padding: 16px 72px 16px 16px; color: #60a5fa; background: rgba(37, 99, 235, 0.1); border: 1px solid rgba(96, 165, 250, 0.35); border-radius: 11px; }
.document-editor :deep([data-remove-attachment]) { position: absolute; top: 7px; right: 7px; padding: 5px 7px; color: #fff; background: rgba(15, 23, 42, 0.7); border: 0; border-radius: 6px; cursor: pointer; font-size: 11px; }
.editor-note { margin-top: 11px; display: flex; gap: 7px; color: var(--editor-muted); font-size: 12px; }
.editor-page-actions { margin-top: 17px; padding: 16px 0; display: flex; justify-content: flex-end; gap: 10px; }
.cancel-button { color: var(--editor-text); background: var(--editor-surface-soft); border: 1px solid var(--editor-border); }
@media (max-width: 720px) {
  .editor-page-header { padding: 20px; grid-template-columns: 1fr auto; }
  .header-copy { grid-column: 1 / -1; grid-row: 1; }
  .back-button { grid-column: 1; grid-row: 2; justify-self: start; }
  .header-publish { grid-column: 2; grid-row: 2; }
  .editor-sheet { padding: 18px; }
  .document-editor { min-height: 420px; padding: 17px; }
  .content-label-row { align-items: flex-start; flex-direction: column; gap: 0; }
}
</style>

<style>
html[data-theme='dark'] .forum-editor-page {
  --editor-surface: rgba(15, 23, 42, 0.94);
  --editor-surface-soft: rgba(2, 6, 23, 0.66);
  --editor-border: rgba(148, 163, 184, 0.28);
  --editor-text: #e2e8f0;
  --editor-muted: #94a3b8;
}
html[data-theme='effect'] .forum-editor-page {
  --editor-surface: rgba(15, 23, 42, 0.66);
  --editor-surface-soft: rgba(2, 6, 23, 0.5);
  --editor-border: rgba(148, 163, 184, 0.3);
  --editor-text: #e2e8f0;
  --editor-muted: #a8b5c7;
}
html[data-theme='custom'] .forum-editor-page {
  --editor-surface: rgba(255, 255, 255, 0.9);
  --editor-surface-soft: rgba(248, 250, 252, 0.84);
}
</style>
