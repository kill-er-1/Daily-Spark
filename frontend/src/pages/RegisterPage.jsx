import React, { useState } from "react";
import { useAuth } from "../context/AuthContext.jsx";
import AuthLayout from "../components/AuthLayout.jsx";

export default function RegisterPage({ onSwitch }) {
  const { register } = useAuth(); // 简化版 AuthContext：只暴露方法
  const [account, setAccount] = useState("");
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const [msg, setMsg] = useState("");
  const [loading, setLoading] = useState(false);

  async function onSubmit(e) {
    e.preventDefault();
    setMsg("");
    const acc = account.trim();
    if (!acc || !password) {
      return setMsg("账号或密码不能为空");
    }
    if (acc.length < 4) {
      return setMsg("账号至少 4 位");
    }
    if (password.length < 8) {
      return setMsg("密码至少 8 位");
    }
    if (password !== confirm) {
      return setMsg("两次输入的密码不一致");
    }
    setLoading(true);
    try {
      const res = await register(acc, password);
      setMsg("注册成功");
      console.log("注册响应：", res);
    } catch (err) {
      setMsg(err?.message || "注册失败");
    } finally {
      setLoading(false);
    }
  }

  return (
    <AuthLayout
      title="创建你的账号"
      subtitle="注册后记录每天的快乐瞬间，留住美好"
      footer={
        <div className="text-sm text-neutral-600 dark:text-neutral-400">
          已有账号？
          <button
            type="button"
            onClick={() => onSwitch && onSwitch("login")}
            className="ml-1 text-amber-700 hover:text-amber-800"
          >
            去登录
          </button>
        </div>
      }
    >
      <form onSubmit={onSubmit} className="space-y-4">
        {msg && (
          <div className="rounded-lg border border-red-200/70 bg-red-50 px-3 py-2 text-sm text-red-700">
            {msg}
          </div>
        )}

        <div>
          <label className="text-sm font-medium text-neutral-700 dark:text-neutral-200">账号</label>
          <input
            className="mt-1 w-full rounded-xl border border-neutral-300 dark:border-neutral-700 bg-white dark:bg-zinc-950 px-3 py-2 outline-none focus:ring-2 focus:ring-amber-500"
            value={account}
            onChange={(e) => setAccount(e.target.value)}
            placeholder="请输入账号"
          />
        </div>

        <div>
          <label className="text-sm font-medium text-neutral-700 dark:text-neutral-200">密码</label>
          <input
            type="password"
            className="mt-1 w-full rounded-xl border border-neutral-300 dark:border-neutral-700 bg-white dark:bg-zinc-950 px-3 py-2 outline-none focus:ring-2 focus:ring-amber-500"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="请输入密码（至少 8 位）"
          />
        </div>

        <div>
          <label className="text-sm font-medium text-neutral-700 dark:text-neutral-200">确认密码</label>
          <input
            type="password"
            className="mt-1 w-full rounded-xl border border-neutral-300 dark:border-neutral-700 bg-white dark:bg-zinc-950 px-3 py-2 outline-none focus:ring-2 focus:ring-amber-500"
            value={confirm}
            onChange={(e) => setConfirm(e.target.value)}
            placeholder="请再次输入密码"
          />
        </div>

        <button
          type="submit"
          disabled={loading}
          className="w-full rounded-xl bg-amber-600 text-white py-2.5 font-medium hover:bg-amber-700 disabled:opacity-60 transition-colors"
        >
          {loading ? "注册中..." : "注册"}
        </button>

        <p className="text-xs text-neutral-500">
          我们只用于记录与管理快乐事件，不会公开你的私有数据。
        </p>
      </form>
    </AuthLayout>
  );
}