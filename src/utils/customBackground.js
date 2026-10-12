const DATABASE_NAME = 'newaoe-preferences';
const DATABASE_VERSION = 1;
const STORE_NAME = 'appearance';
const BACKGROUND_KEY = 'custom-background';

const openDatabase = () => new Promise((resolve, reject) => {
  if (!window.indexedDB) {
    reject(new Error('当前浏览器不支持保存自定义背景'));
    return;
  }

  const request = window.indexedDB.open(DATABASE_NAME, DATABASE_VERSION);
  request.onupgradeneeded = () => {
    const database = request.result;
    if (!database.objectStoreNames.contains(STORE_NAME)) {
      database.createObjectStore(STORE_NAME);
    }
  };
  request.onsuccess = () => resolve(request.result);
  request.onerror = () => reject(request.error || new Error('打开背景存储失败'));
});

const runStoreRequest = async (mode, operation) => {
  const database = await openDatabase();
  return new Promise((resolve, reject) => {
    const transaction = database.transaction(STORE_NAME, mode);
    const store = transaction.objectStore(STORE_NAME);
    const request = operation(store);

    request.onsuccess = () => resolve(request.result);
    request.onerror = () => reject(request.error || new Error('保存背景失败'));
    transaction.oncomplete = () => database.close();
    transaction.onerror = () => database.close();
    transaction.onabort = () => database.close();
  });
};

export const getCustomBackground = () => runStoreRequest(
  'readonly',
  store => store.get(BACKGROUND_KEY)
);

export const saveCustomBackground = file => runStoreRequest(
  'readwrite',
  store => store.put(file, BACKGROUND_KEY)
);

export default {
  getCustomBackground,
  saveCustomBackground
};
