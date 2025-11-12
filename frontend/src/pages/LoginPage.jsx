import React, { useState } from "react";
import { useAuth } from "../context/AuthContext.jsx";
import AuthLayout from "../components/AuthLayout.jsx";

export default function LoginPage({ onSwitch }) {
  const { login } = useAuth();
  const [account, setAccount] = useState("");
  const [password, setPassword] = useState("");
  const [msg, setMsg] = useState("");
  const [loading, setLoading] = useState(false);
  const [showPw, setShowPw] = useState(false);

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

    setLoading(true);
    try {
      await login(acc, password);
      setMsg("登录成功");
    } catch (err) {
      setMsg(err?.message || "登录失败");
    } finally {
      setLoading(false);
    }
  }

  return (
    <AuthLayout
      title="欢迎回来"
      subtitle="登录后开始记录今天的快乐事件"
      footer={
        <div className="text-sm text-neutral-600 dark:text-neutral-400">
          没有账号？
          <button
            type="button"
            onClick={() => onSwitch && onSwitch("register")}
            className="ml-1 text-amber-700 hover:text-amber-800"
          >
            去注册
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
          <div className="mt-1 relative">
            <input
              type={showPw ? "text" : "password"}
              className="w-full rounded-xl border border-neutral-300 dark:border-neutral-700 bg-white dark:bg-zinc-950 px-3 py-2 pr-10 outline-none focus:ring-2 focus:ring-amber-500"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="请输入密码"
            />
            <button
              type="button"
              onClick={() => setShowPw((v) => !v)}
              className="absolute right-2 top-1/2 -translate-y-1/2 text-sm text-neutral-500 hover:text-neutral-700"
            >
              {showPw ? "隐藏" : "显示"}
            </button>
          </div>
        </div>

        <button
          type="submit"
          disabled={loading}
          className="w-full rounded-xl bg-amber-600 text-white py-2.5 font-medium hover:bg-amber-700 disabled:opacity-60 transition-colors"
        >
          {loading ? "登录中..." : "登录"}
        </button>

        <p className="text-xs text-neutral-500">
          登录即表示你同意本系统用于记录与管理快乐事件的数据处理。
        </p>
      </form>
    </AuthLayout>
  );
}