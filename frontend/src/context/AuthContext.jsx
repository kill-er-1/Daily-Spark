import React, { createContext, useContext } from "react";
import { api } from "../services/api";

const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
  async function login(account, password) {
    return api.login({ account, password });
  }

  async function register(account, password) {
    return api.register({ account, password });
  }

  async function updatePassword(userId, newPassword) {
    return api.updatePassword(userId, { password: newPassword });
  }

  async function deleteUser(userId, operatorAccount) {
    return api.deleteUser(userId, { account: operatorAccount });
  }

  const value = { login, register, updatePassword, deleteUser };
  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
}