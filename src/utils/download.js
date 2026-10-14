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
  let identifier = normalizeIdentifier(fileIdentifier);

  // 部分接口会直接返回已签名的完整下载地址，而不是统一下载接口的文件标识。
  if (/^https?:\/\//i.test(identifier)) {
    return {
      url: identifier,
      download_type: config.DownloadURLType_UnKnown,
      toString() {
        return this.url;
      },
      valueOf() {
        return this.url;
      }
    };
  }

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
  const download_type=result?.data?.download_type??config.DownloadURLType_UnKnown;
  if (!url) {
    throw new Error(`下载链接无效：${result?.msg || '未知错误'}`);
  }
  //
  return {
    url: url,
    download_type: download_type,
    toString() {
      return this.url;
    },
    valueOf() {
      return this.url;
    }
  };
};

/**
 * 获取真实地址并触发浏览器下载。
 */
export const downloadFile = async (fileIdentifier, options = {}) => {
  let result = await getDownloadUrl(fileIdentifier, {
    ...options,
    forceDownload: true
  });
  let url=result.url;
  if (result.download_type===config.DownloadURLType_CDN) {
    const urlObj = new URL(url);
    urlObj.searchParams.set('attachment', 'true');
    url=urlObj.toString();
  }
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

/**
 * 通过统一下载接口取得文件内容，供 HTML、日志等需要在页面内展示的文本文件使用。
 * 与 downloadFile 相同，先解析 ContentPath 对应的真实下载地址，再下载文件；
 * 区别是这里返回文件内容，而不是触发浏览器“另存为”。
 */
export const downloadTextFile = async (fileIdentifier, { signal } = {}) => {
  const result = await getDownloadUrl(fileIdentifier, { signal });
  const response = await fetch(result.url, {
    method: 'GET',
    signal,
    credentials: 'include',
  });
  if (!response.ok) {
    throw new Error(`下载文件失败（HTTP ${response.status}）`);
  }
  const fileBlob = await response.blob();
  return {
    content: await fileBlob.text(),
    url: result.url,
    contentType: fileBlob.type
  };
};

export default {
  getDownloadUrl,
  downloadFile,
  downloadTextFile
};
