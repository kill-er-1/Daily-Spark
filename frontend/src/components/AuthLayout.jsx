import React from "react";

export default function AuthLayout({ title, subtitle, children, footer }) {
  return (
    <div className="min-h-screen w-full bg-gradient-to-br from-amber-50 via-rose-50 to-orange-50 dark:from-zinc-900 dark:via-zinc-900 dark:to-zinc-900 flex items-center justify-center p-6">
      <div className="w-full max-w-4xl grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="hidden md:flex flex-col justify-center rounded-2xl bg-white/60 dark:bg-white/5 backdrop-blur p-8 border border-white/40 dark:border-white/10 shadow-sm">
          <div className="flex items-center gap-3 mb-4">
            <span className="text-2xl">🌟</span>
            <h2 className="text-xl font-semibold text-amber-800">心情笔记管理系统</h2>
          </div>
          <p className="text-sm text-neutral-600 dark:text-neutral-300">
            记录每一天最快乐的一件事，提醒自己过去的每一天都有意义。
            支持上传、时间排序查看、编辑和删除快乐事件；管理员可维护公开数据与用户信息。
          </p>
        </div>

        <div className="rounded-2xl bg-white dark:bg-zinc-950 border border-neutral-200/70 dark:border-neutral-800 shadow-xl">
          <div className="p-6 md:p-8">
            <div className="mb-6">
              <h1 className="text-2xl md:text-3xl font-bold tracking-tight">{title}</h1>
              {subtitle && (
                <p className="mt-2 text-sm text-neutral-600 dark:text-neutral-400">{subtitle}</p>
              )}
            </div>

            {children}

            {footer && <div className="mt-6">{footer}</div>}
          </div>
        </div>
      </div>
    </div>
  );
}