// src/pages/Register.jsx
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card, Button, Input, Form, message } from 'antd';
import { register } from '../api/authApi';

const Register = () => {
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const [form] = Form.useForm();

  const handleRegister = async () => {
    try {
      const values = await form.validateFields();
      setLoading(true);

      // 调用注册接口（仅普通用户）
      const result = await register(values.account, values.password);

      if (result.success) {
        message.success('注册成功，请登录');
        navigate('/login'); // 跳转到登录页
      } else {
        message.error(result.message || '注册失败');
      }
    } catch (error) {
      if (error.name !== 'ValidateError') {
        message.error('注册失败，请重试');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'flex-start', marginTop: 0, backgroundColor: '#E9F3F6', paddingTop: 0, height: '100vh', overflow: 'hidden', boxSizing: 'border-box' }}>
      <Card title="注册" style={{ width: 400 }}>
        <Form form={form} layout="vertical">
          <Form.Item
            name="account"
            label="账号"
            rules={[{ required: true, message: '请输入账号' },
            { min: 3, message: '账号至少3个字符'}
            ]}
          >
            <Input placeholder="请输入账号" />
          </Form.Item>

          <Form.Item
            name="password"
            label="密码"
            rules={[{ required: true, message: '请输入密码' },
            { min: 8, message: '密码至少8个字符'}
            ]}
          >
            <Input.Password placeholder="请输入密码" />
          </Form.Item>

          <Form.Item
            name="checkPassword"
            label="确认密码"
            rules={[
              { required: true, message: '请确认密码' },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue('password') === value) {
                    return Promise.resolve();
                  }
                  return Promise.reject(new Error('两次密码输入不一致'));
                },
              }),
            ]}
          >
            <Input.Password placeholder="请再次输入密码" />
          </Form.Item>

          <Form.Item>
            <Button
              type="primary"
              onClick={handleRegister}
              loading={loading}
              style={{ backgroundColor: '#00687B', borderColor: '#00687B', width: '85px', height: '35px', margin: '0 auto', display: 'block' }}
            >
              注册
            </Button>
          </Form.Item>

          <div style={{ textAlign: 'center' }}>
            已有账号？{' '}
            <a href="/login" style={{ color: '#00687B' }}>
              立即登录
            </a>
          </div>
        </Form>
      </Card>
    </div>
  );
};

export default Register;
