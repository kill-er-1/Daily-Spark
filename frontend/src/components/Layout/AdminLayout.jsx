// src/components/Layout/AdminLayout.jsx（管理员布局）
import { Outlet } from 'react-router-dom';
import { Layout } from 'antd';

const { Header, Content, Footer } = Layout;

const AdminLayout = () => {
  return (
    <Layout style={{ minHeight: '100vh', background: '#E9F3F6' }}>
      <Header style={{ background: '#E9F3F6', padding: '0 20px' }}>
        <h2 style={{ color: '#00687b', marginTop: 15 }}>心情笔记 - 管理员端</h2>
      </Header>
      <Content style={{ padding: '20px', background: '#E9F3F6' }}>
        <Outlet />
      </Content>
      <Footer style={{ textAlign: 'center', background: '#E9F3F6' }}>心情笔记 ©{new Date().getFullYear()}</Footer>
    </Layout>
  );
};

export default AdminLayout;
