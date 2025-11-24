import { useEffect, useState } from 'react'
import { Table, message } from 'antd'
import { queryUsers } from '../../api/authApi'

const UserList = () => {
  const [loading, setLoading] = useState(false)
  const [users, setUsers] = useState([])

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        setLoading(true)
        const result = await queryUsers()
        if (result.success) {
          setUsers(result.data?.users || [])
        } else {
          message.error(result.message || '获取用户失败')
        }
      } catch (err) {
        message.error(err?.message || '获取用户失败')
      } finally {
        setLoading(false)
      }
    }
    fetchUsers()
  }, [])

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id' },
    { title: '账号', dataIndex: 'account', key: 'account' },
    { title: '昵称', dataIndex: 'nickname', key: 'nickname' },
    {
      title: '管理员',
      dataIndex: 'is_admin',
      key: 'is_admin',
      render: (val) => (val ? '是' : '否'),
    },
    { title: '创建时间', dataIndex: 'created_at', key: 'created_at' },
    { title: '更新时间', dataIndex: 'updated_at', key: 'updated_at' },
  ]

  return (
    <Table
      rowKey="id"
      columns={columns}
      dataSource={users}
      loading={loading}
      pagination={{ pageSize: 10 }}
    />
  )
}

export default UserList
