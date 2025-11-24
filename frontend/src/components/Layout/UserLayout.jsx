// src/components/Layout/UserLayout.jsx（用户布局）
import { Outlet } from 'react-router-dom'; // 用于渲染子路由
import { Layout } from 'antd';

const { Header, Content, Footer } = Layout;

const UserLayout = () => {
  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ background: '#e9f3f6', padding: '0 20px' }}>
        <h2 style={{ color: '#00687b', marginTop: 15 }}>心情笔记 - 用户端</h2>
      </Header>
      <Content style={{ padding: '20px' }}>
        <Outlet /> {/* 子路由页面会在这里渲染 */}
      </Content>
      <Footer style={{ textAlign: 'center' }}>心情笔记 ©{new Date().getFullYear()}</Footer>
    </Layout>
  );
};

export default UserLayout;