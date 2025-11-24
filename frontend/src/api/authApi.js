// src/api/authApi.js
import request from '../utils/request';

// 用户登录
export const login = async (account, password, role) => {
  // 调用后端登录接口，参数包含身份（role: 'user' 或 'admin'）
  const data = await request.post('api/v1/users/signin', {
    account,
    password
  });
  return data;
};

// 用户注册（仅普通用户）
export const register = async (account,password) => {
  const data = await request.post('api/v1/users/signup', {
    account,
    password,
  });
  return data;
};

export const queryUsers = async () => {
  const data = await request.get('api/v1/user/query');
  return data;
};

export const updateUser = async (id, payload) => {
  const data = await request.post(`api/v1/users/update/${id}` , payload);
  return data;
};
