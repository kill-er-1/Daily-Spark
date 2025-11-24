// src/utils/auth.js
// 存储键名（避免硬编码）
const TOKEN_KEY = 'happy_note_token';
const USER_ROLE_KEY = 'happy_note_role'; // 'user' 或 'admin'
const USER_ID_KEY = 'happy_note_user_id';

// 存储 token 到本地存储
export const setToken = (token) => {
  localStorage.setItem(TOKEN_KEY, token);
};

// 获取本地存储的 token
export const getToken = () => {
  return localStorage.getItem(TOKEN_KEY);
};

// 存储用户角色
export const setUserRole = (role) => {
  localStorage.setItem(USER_ROLE_KEY, role);
};

// 获取用户角色
export const getUserRole = () => {
  return localStorage.getItem(USER_ROLE_KEY) || '';
};

export const setUserId = (id) => {
  localStorage.setItem(USER_ID_KEY, id);
};

export const getUserId = () => {
  return localStorage.getItem(USER_ID_KEY) || '';
};

// 退出登录（清除本地存储）
export const logout = () => {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(USER_ROLE_KEY);
  localStorage.removeItem(USER_ID_KEY);
};

// 判断是否为管理员
export const isAdmin = () => {
  return getUserRole() === 'admin';
};
