import React, { createContext, useContext } from "react";
import { api } from "../services/api";

type AuthContextType = {
  login: (account: string, password: string) => Promise<any>;
  register: (account: string, password: string) => Promise<any>;
  updatePassword: (userId: string, newPassword: string) => Promise<any>;
  deleteUser: (userId: string, operatorAccount: string) => Promise<any>;
};

const AuthContext = createContext<AuthContextType | null>(null);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  async function login(account: string, password: string) {
    // 后端比较 bcrypt；前端直接传明文
    return api.login({ account, password });
  }

  async function register(account: string, password: string) {
    return api.register({ account, password });
  }

  async function updatePassword(userId: string, newPassword: string) {
    // 注意路径参数必须传正确的 userId
    return api.updatePassword(userId, { password: newPassword });
  }

  async function deleteUser(userId: string, operatorAccount: string) {
    // 操作者账号从请求体读取（后端已改为请求体方式）
    return api.deleteUser(userId, { account: operatorAccount });
  }

  const value: AuthContextType = { login, register, updatePassword, deleteUser };
  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
}