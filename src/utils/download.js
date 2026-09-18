import config from '../config';

const normalizeIdentifier = fileIdentifier => {
  const normalized = String(fileIdentifier ?? '').trim().replace(/^\/+/, '');
  if (!normalized) {
    throw new Error('文件链接无效');
  }
  return normalized;
};

/**
 * 通过统一下载接口获取文件的真实访问地址。
 * forceDownload 仅用于需要浏览器下载行为的场景，头像等资源保持默认值即可。
 */
export const getDownloadUrl = async (
  fileIdentifier,
  { signal, forceDownload = false } = {}
) => {
  const identifier = normalizeIdentifier(fileIdentifier);
  const requestUrl = new URL(`${config.download_url}/${identifier}`);
  if (forceDownload) {
    requestUrl.searchParams.set('attachment', 'true');
  }

  const response = await fetch(requestUrl.toString(), {
    method: 'GET',
    credentials: 'include',
    signal
  });

  let result;
  try {
    result = await response.json();
  } catch (error) {
    throw new Error(`获取下载链接失败（HTTP ${response.status}）`);
  }

  if (!response.ok || result?.status === false) {
    throw new Error(result?.msg || `获取下载链接失败（HTTP ${response.status}）`);
  }

  const url = typeof result?.data?.url === 'string'
    ? result.data.url.trim()
    : '';
  if (!url) {
    throw new Error(`下载链接无效：${result?.msg || '未知错误'}`);
  }
  return url;
};

/**
 * 获取真实地址并触发浏览器下载。
 */
export const downloadFile = async (fileIdentifier, options = {}) => {
  const url = await getDownloadUrl(fileIdentifier, {
    ...options,
    forceDownload: true
  });
  const fileName = url.split('?')[0].split('/').pop() || 'download';
  const link = document.createElement('a');
  link.href = url;
  link.download = fileName;
  link.target = '_self';
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  return { url, fileName };
};

export default {
  getDownloadUrl,
  downloadFile
};
