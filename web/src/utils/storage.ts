export function getStorageItem(key: string) {
  const str = typeof key !== 'undefined' && localStorage ? localStorage.getItem(key) : null;
  try {
    if (str) {
      return JSON.parse(str)
    }
  } catch (e) {
    return null;
  }
  return null;
}

export function setStorageItem(key: string, val: any) {
  localStorage.setItem(key, JSON.stringify(val));
}

export function removeStorageItem(key: string) {
  localStorage.removeItem(key);
}
