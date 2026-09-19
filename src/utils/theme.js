export const THEME_MODES = [
  { id: 'white', label: '纯白' },
  { id: 'dark', label: '深色' },
  { id: 'effect', label: '特效' }
];

const THEME_STORAGE_KEY = 'aoe-theme';

export const getInitialTheme = () => {
  if (typeof window === 'undefined') return 'effect';
  const storedTheme = window.localStorage.getItem(THEME_STORAGE_KEY);
  return THEME_MODES.some(theme => theme.id === storedTheme) ? storedTheme : 'effect';
};

export const applyTheme = theme => {
  const nextTheme = THEME_MODES.some(item => item.id === theme) ? theme : 'effect';
  if (typeof document !== 'undefined') {
    document.documentElement.dataset.theme = nextTheme;
    document.documentElement.style.colorScheme = nextTheme === 'white' ? 'light' : 'dark';
  }
  if (typeof window !== 'undefined') {
    window.localStorage.setItem(THEME_STORAGE_KEY, nextTheme);
  }
  return nextTheme;
};

// 在应用启动前设置主题，减少首次渲染时的闪白。
applyTheme(getInitialTheme());
