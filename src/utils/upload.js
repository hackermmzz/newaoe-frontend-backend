import { ElMessage } from 'element-plus';

const getFormField = (formData, fieldName) => {
  const matchedEntry = Object.entries(formData).find(
    ([key]) => key.toLowerCase() === fieldName.toLowerCase()
  );
  return matchedEntry?.[1];
};

const decodePolicy = encodedPolicy => {
  const normalizedPolicy = encodedPolicy
    .replace(/-/g, '+')
    .replace(/_/g, '/');
  const paddedPolicy = normalizedPolicy.padEnd(
    normalizedPolicy.length + ((4 - normalizedPolicy.length % 4) % 4),
    '='
  );
  return JSON.parse(window.atob(paddedPolicy));
};

const getContentLengthRange = uploadPolicy => {
  const formData = uploadPolicy?.formdata;
  if (!formData || typeof formData !== 'object' || Array.isArray(formData)) {
    throw new Error('上传凭证缺少 formdata');
  }

  const encodedPolicy = getFormField(formData, 'policy');
  if (typeof encodedPolicy !== 'string' || !encodedPolicy) {
    throw new Error('上传凭证缺少文件大小策略');
  }

  let policy;
  try {
    policy = decodePolicy(encodedPolicy);
  } catch (error) {
    throw new Error('上传凭证中的文件大小策略无效');
  }

  const condition = policy?.conditions?.find(item => (
    Array.isArray(item)
    && String(item[0]).toLowerCase() === 'content-length-range'
  ));
  const minSize = Number(condition?.[1]);
  const maxSize = Number(condition?.[2]);
  if (
    !Number.isFinite(minSize)
    || !Number.isFinite(maxSize)
    || minSize < 0
    || maxSize < minSize
  ) {
    throw new Error('上传凭证缺少有效的文件大小限制');
  }

  return { minSize, maxSize };
};

const formatFileSize = bytes => {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  if (bytes < 1024 * 1024 * 1024) return `${(bytes / 1024 / 1024).toFixed(1)} MB`;
  return `${(bytes / 1024 / 1024 / 1024).toFixed(1)} GB`;
};

const createNotifiedError = message => {
  ElMessage.error(message);
  const error = new Error(message);
  error.notified = true;
  return error;
};

/**
 * 使用服务端返回的 POST Policy 批量上传文件。
 * 在发起任何请求前，会先校验全部文件是否符合各自 policy 的大小限制。
 */
export const uploadFiles = async (files, uploadPolicies, { signal } = {}) => {
  if (!Array.isArray(files) || !Array.isArray(uploadPolicies)) {
    throw new TypeError('文件和上传凭证必须为数组');
  }
  if (files.length !== uploadPolicies.length) {
    throw new Error('文件数量与上传凭证数量不一致');
  }

  const preparedUploads = files.map((file, index) => {
    if (!file || !Number.isFinite(file.size)) {
      throw new TypeError(`第 ${index + 1} 个文件无效`);
    }

    const uploadPolicy = uploadPolicies[index];
    const uploadUrl = typeof uploadPolicy?.url === 'string'
      ? uploadPolicy.url.trim()
      : '';
    if (!uploadUrl) {
      throw new Error(`第 ${index + 1} 个上传地址无效`);
    }

    return {
      file,
      uploadUrl,
      formDataFields: uploadPolicy.formdata,
      ...getContentLengthRange(uploadPolicy)
    };
  });

  const sizeMismatch = preparedUploads.find(({ file, minSize, maxSize }) => (
    file.size < minSize || file.size > maxSize
  ));
  if (sizeMismatch) {
    const { file, minSize, maxSize } = sizeMismatch;
    const fileName = file.name || '未命名文件';
    if (file.size > maxSize) {
      throw createNotifiedError(
        `文件“${fileName}”太大（${formatFileSize(file.size)}），最大允许 ${formatFileSize(maxSize)}`
      );
    }
    throw createNotifiedError(
      `文件“${fileName}”太小（${formatFileSize(file.size)}），最小要求 ${formatFileSize(minSize)}`
    );
  }

  return Promise.all(preparedUploads.map(async ({ file, uploadUrl, formDataFields }) => {
    const requestBody = new FormData();
    Object.entries(formDataFields).forEach(([key, value]) => {
      requestBody.append(key, value);
    });
    requestBody.append('file', file, file.name || 'file');

    const response = await fetch(uploadUrl, {
      method: 'POST',
      body: requestBody,
      signal
    });
    if (!response.ok) {
      throw new Error(`文件“${file.name || '未命名文件'}”上传失败（HTTP ${response.status}）`);
    }
    return response;
  }));
};

export default {
  uploadFiles
};
