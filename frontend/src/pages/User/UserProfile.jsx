import { useState } from 'react'
import { Card, Form, Input, Button, message } from 'antd'
import { updateUser } from '../../api/authApi'
import { getUserId } from '../../utils/auth'

const UserProfile = () => {
  const [loading, setLoading] = useState(false)
  const [form] = Form.useForm()

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      const { nickname, password } = values
      if (!nickname && !password) {
        message.error('请至少填写昵称或密码')
        return
      }
      const id = getUserId()
      if (!id) {
        message.error('未获取到用户ID，请重新登录')
        return
      }
      setLoading(true)
      const result = await updateUser(id, { nickname, password })
      if (result.success) {
        message.success(result.message || '用户信息已更新')
        form.resetFields(['password'])
      } else {
        message.error(result.message || '更新失败')
      }
    } catch (err) {
      if (err.name !== 'ValidateError') {
        message.error(err?.message || '更新失败，请重试')
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={{ display: 'flex', justifyContent: 'center' }}>
      <Card title="修改用户信息" style={{ width: 480 }}>
        <Form form={form} layout="vertical">
          <Form.Item name="nickname" label="新的昵称">
            <Input placeholder="请输入新的昵称（可选）" />
          </Form.Item>
          <Form.Item name="password" label="新的密码">
            <Input.Password placeholder="请输入新的密码（可选）" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" loading={loading} onClick={handleSubmit} style={{ backgroundColor: '#00687B', borderColor: '#00687B' }}>
              更新
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  )
}

export default UserProfile
