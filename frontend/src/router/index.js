// src/router/index.js
import { createBrowserRouter, RouterProvider, Navigate } from 'react-router-dom';
import { getToken, getUserRole } from '../utils/auth'; // 后续会实现的工具函数

// 公共页面
import Login from '../pages/Login';
import Register from '../pages/Register';

// 用户页面
import UserLayout from '../components/Layout/UserLayout'; // 用户布局（带导航）
import NoteHome from '../pages/User/NoteHome';
import NoteList from '../pages/User/NoteList';
import NoteEdit from '../pages/User/NoteEdit';
import UserProfile from '../pages/User/UserProfile';

// 管理员页面
import AdminLayout from '../components/Layout/AdminLayout'; // 管理员布局（带侧边栏）
import UserList from '../pages/Admin/UserList';

// 权限控制组件：未登录则跳转登录页
const PrivateRoute = ({ children, role }) => {
  const token = getToken();
  const userRole = getUserRole(); // 'user' 或 'admin'

  if (!token) {
    return <Navigate to="/login" replace />; // 未登录 → 登录页
  }

  // 管理员路由必须验证角色
  if (role === 'admin' && userRole !== 'admin') {
    return <Navigate to="/user/home" replace />; // 非管理员访问管理员页 → 用户首页
  }

  return children;
};

// 路由配置
const router = createBrowserRouter([
  // 未登录可访问的路由
  { path: '/login', element: <Login /> },
  { path: '/register', element: <Register /> },

  // 用户路由（需登录 + 普通用户身份）
  {
    path: '/user',
    element: (
      <PrivateRoute role="user">
        <UserLayout /> {/* 包裹用户所有页面的布局（导航栏等） */}
      </PrivateRoute>
    ),
    children: [
      { path: 'home', element: <NoteHome /> }, // 首页（今日笔记）
      { path: 'notes', element: <NoteList /> }, // 历史笔记列表
      { path: 'notes/edit/:id?', element: <NoteEdit /> }, // 编辑笔记（id 可选，用于新建）
      { path: 'profile', element: <UserProfile /> },
      { path: '', element: <Navigate to="home" replace /> }, // 默认跳转用户首页
    ],
  },

  // 管理员路由（需登录 + 管理员身份）
  {
    path: '/admin',
    element: (
      <PrivateRoute role="admin">
        <AdminLayout /> {/* 包裹管理员所有页面的布局（侧边栏等） */}
      </PrivateRoute>
    ),
    children: [
      { path: 'users', element: <UserList /> }, // 用户管理列表
      { path: '', element: <Navigate to="users" replace /> }, // 默认跳转用户管理
    ],
  },

  // 根路径默认跳转登录页
  { path: '/', element: <Navigate to="/login" replace /> },
]);

// 路由提供者组件，供 App.js 使用
export default function AppRouter() {
  return <RouterProvider router={router} />;
}
