// src/utils/request.js
import axios from 'axios';
import { getToken, logout } from './auth';
import { message } from 'antd'; // 从 antd 引入提示组件

// 创建 axios 实例
const request = axios.create({
  baseURL: 'http://10.83.132.135:8080', // 后端接口基础地址（根据实际修改）
  timeout: 5000,
});

// 请求拦截器：添加 token 到请求头
request.interceptors.request.use(
  (config) => {
    const token = getToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`; // 按后端要求的格式传递
    }
    config.headers['Content-Type'] = 'application/json'; // 设置请求头，根据后端要求
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器：处理错误
request.interceptors.response.use(
  (response) => {
    const { status, data } = response;
    return {
      success: status >= 200 && status < 300,
      message: data?.message ?? data?.error ?? '',
      data: data?.data ?? data,
      status,
    };
  },
  (error) => {
    const status = error.response?.status ?? 0;
    const errData = error.response?.data;
    let msg = errData?.message ?? errData?.error ?? '请求失败';
    const url = String(error.config?.url || '');
    if (!errData?.message && !errData?.error) {
      if (status === 400) msg = 'invalid request body';
      else if (status === 404) msg = 'account not exists';
      else if (status === 409) msg = 'account already exists';
      else if (status === 500) msg = 'internal server error description';
    }
    if (status === 401 && !url.includes('/users/signin')) {
      message.error('登录已过期，请重新登录');
      logout();
      window.location.href = '/login';
    } else {
      if (status === 401 && !errData?.message && !errData?.error) {
        msg = 'password not match';
      }
      message.error(msg);
    }
    return Promise.reject({
      success: false,
      message: msg,
      status,
      data: errData,
    });
  }
);

export default request;
