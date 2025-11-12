import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwind from "@tailwindcss/vite";

// 通过代理将同源下的 /api/v1 转发到后端 http://localhost:8080
// 这样浏览器不再触发跨域预检（OPTIONS），避免 404。
export default defineConfig({
  plugins: [react(), tailwind()],
  server: {
    proxy: {
      "/api/v1": {
        target: "http://localhost:8080",
        changeOrigin: true,
        // 如后端使用自签名证书，可加：secure: false
      },
    },
  },
});
