// src/pages/Login.jsx
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card, Button, Input, Radio, Form, message } from 'antd';
import { login } from '../api/authApi';
import { setToken, setUserRole, setUserId } from '../utils/auth';

const Login = () => {
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const [form] = Form.useForm();

  // 处理登录提交
  const handleSubmit = async () => {
    try {
      const values = await form.validateFields(); // 获取表单值
      setLoading(true);

      const result = await login(
        values.account,
        values.password,
        values.role
      );

      if (result.success) {
        const user = result.data?.user;
        setToken('signed-in');
        if (user?.id) setUserId(user.id);
        setUserRole(user?.is_admin ? 'admin' : 'user');
        if (user?.is_admin) {
          navigate('/admin/users');
        } else {
          navigate('/user/home');
        }
        message.success(result.message || '登录成功！');
      } else {
        message.error(result.message || '登录失败');
      }
    } catch (error) {
      if (error.name !== 'ValidateError') {
        message.error(error?.message || '登录失败，请重试');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'flex-start', marginTop: 0, backgroundColor: '#E9F3F6', paddingTop: 0, height: '100vh', overflow: 'hidden', boxSizing: 'border-box' }}>
      <Card title="登录" style={{ width: 400 }}>
        <Form form={form} layout="vertical" initialValues={{ role: 'user' }}>
          {/* 用户名 */}
          <Form.Item
            name="account"
            label="账号"
            rules={[{ required: true, message: '请输入账号' }]}
          >
            <Input placeholder="请输入账号" />
          </Form.Item>

          {/* 密码 */}
          <Form.Item
            name="password"
            label="密码"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password placeholder="请输入密码" />
          </Form.Item>

          {/* 身份选择 */}
          <Form.Item name="role" label="身份">
            <Radio.Group>
              <Radio.Button value="user">普通用户</Radio.Button>
              <Radio.Button value="admin">管理员</Radio.Button>
            </Radio.Group>
          </Form.Item>

          {/* 登录按钮 */}
          <Form.Item>
            <Button
              type="primary"
              onClick={handleSubmit}
              loading={loading}
              style={{ backgroundColor: '#00687B', borderColor: '#00687B', width: 85, height: 35, margin: '0 auto', display: 'block' }}
            >
              登录
            </Button>
          </Form.Item>

          {/* 注册链接 */}
          <div style={{ textAlign: 'center' }}>
            还没有账号？{' '}
            <a href="/register" style={{ color: '#00687B' }}>
              立即注册
            </a>
          </div>
        </Form>
      </Card>
    </div>
  );
};

export default Login;
