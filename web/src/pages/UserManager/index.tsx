import React, {useEffect, useRef, useState} from 'react';
import {PageContainer} from '@ant-design/pro-layout';
import {Button, Card, Col, Divider, Input, message, Popconfirm, Row, Space, Table, Tag, Tooltip} from 'antd';
import type {UserListData, UserListItem} from './data';
import {queryUser, updateUser,removeUser} from './service';
import {PlusOutlined, ReloadOutlined} from '@ant-design/icons';
import moment from 'moment';
import UsersForm from './components/UserForm';
import {ActionType} from "@ant-design/pro-table";
import { useAccess } from 'umi';

const { Search } = Input;
const handleSearchKeyword = async (val: string) => {
  console.log('search val:', val);
};

const query = async (params: string) =>{
  try {
    return await queryUser(params);
  } catch (e) {
    return {success: false, msg: e}
  }
}


/**
 * 更新节点
 * @param fields
 */
const handleUpdate = async (fields: UserListItem) => {

  const hide = message.loading('正在配置');
  try {
    const res = await updateUser({...fields});
    hide();
    message.success('配置成功');
    return res;
  } catch (error) {
    hide();
    message.error('配置失败请重试！');
    return {success: false, msg: error}
  }
};

/**
 *  删除节点
 * @param selectedRows
 */
const handleRemove = async (username: string) => {
  const hide = message.loading('正在删除');
  try {
    await removeUser({
      "username": username,
    });
    hide();
    message.success('删除成功，即将刷新');
    return true;
  } catch (error) {
    hide();
    message.error('删除失败，请重试');
    return false;
  }
};

const UserManager: React.FC = () => {
  const [list, setList] = useState<UserListData[]>([]);
  const [total, setTotal] = useState<number>(0);
  const [loading, setLoading] = useState<boolean>(false);
  const [updateModalVisible, handleUpdateModalVisible] = useState<boolean>(false);
  const [currentRow, setCurrentRow] = useState<UserListItem>();
  const [keyword, setKeyword] = useState<string>();
  const [currentPage, setCurrentPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(10);
  const actionRef = useRef<ActionType>();
  const access = useAccess();


  const did = (params: any) => {
    setLoading(true);
    const data  = {
      offset: pageSize * (currentPage >= 2 ? currentPage - 1 : 0),
      limit: pageSize,
      keyword: params && params.keyword ? params.keyword : keyword,
      ...params
    }
    console.log("debug did data -->", data)
    query(data).then((res) => {
      if (res.success) {
        setList(res.data);
        setTotal(res.total);
      }
      setLoading(false);
    });
  }

  const columns = [
    {
      title: '用户名',
      dataIndex: 'username',
      sorter: true,
      render: (text: string) => <a>{text}</a>,
    },
    {
      title: '姓名',
      dataIndex: 'chineseName',
      sorter: true,
    },
    {
      title: '管理员',
      dataIndex: 'admin',
      sorter: true,
      render: (text: boolean) => <Tag color={text ? 'green' : '' }>{text ? '是' : '否'}</Tag>,
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      sorter: true,
      render: (text: string) => moment(text).format('YYYY-MM-DD HH:mm:ss'),
    },
    {
      title: '修改时间',
      dataIndex: 'updatedAt',
      sorter: true,
      render: (text: string) => moment(text).format('YYYY-MM-DD HH:mm:ss'),
    },
    {
      title: '操作',
      dataIndex: 'id',
      key: 'id',
      fixed: 'right',
      width: 150,
      render: (text: number, record: any) => (
        <>
          <Space split={<Divider type="vertical" />}>
            <Tooltip title={`修改用户【${record.username}】`}>
              <a onClick={() => {
                console.log("debug ---> ", record)
                handleUpdateModalVisible(true);
                setCurrentRow({...record, modify: true});
              }}>修改</a>
            </Tooltip>
            <Tooltip title={`删除用户【${record.username}】`}>
              <Popconfirm
                title={`删除【${record.username}】，删除后数据不可恢复。是否继续？`}
                placement="left"
                onConfirm={async ()=>{
                  if (!access.canAdmin) {message.error('操作权限受限，请联系平台管理员');return}
                  const success = await handleRemove(record.username);
                  if (success) {
                    if (actionRef.current) {
                      actionRef.current.reload();
                    }
                  }
                }}
              >
                <a>删除</a>
              </Popconfirm>
            </Tooltip>
          </Space>
        </>
      ),
    },
  ];

  useEffect(() => {
    did('')
  }, []);

  const handleStandardTableChange = (pagination: { pageSize: number; current: number; }, _: any, sorter: any) => {
    const params = {
      offset: pagination.pageSize * (pagination.current >= 2 ? pagination.current - 1 : 0),
      limit: pagination.pageSize,
      keyword,
      sorterField: "",
      sorterOrder: ""

    };
    if (sorter.field) {
      params.sorterField = `${sorter.field}`;
      params.sorterOrder = `${sorter.order}`;
    }
    setCurrentPage(pagination.current);
    setPageSize(pagination.pageSize)
    did(params);
  };

  // @ts-ignore
  return (
    <PageContainer>
      <Card size="small" bodyStyle={{ padding: 10 }}>
        <Row>
          <Col flex="auto">
            <Search
              placeholder="支持搜索账号、姓名"
              onSearch={(val) => {
                console.log("debug on search --> ", val)
                setKeyword(val);
                did({keyword: val});
              }}
              style={{ width: 280 }}
            />
            <Tooltip placement="top" title="重载并刷新表格数据">
              <Button
                type="link"
                icon={<ReloadOutlined />}
                onClick={() => did('')}
              />
            </Tooltip>
          </Col>
          <Col span={2}>
            <Button type="link" icon={<PlusOutlined />} onClick={()=> handleUpdateModalVisible(true)}>
              新增用户
            </Button>
          </Col>
        </Row>
        <Row style={{ paddingTop: 10 }}>
          <Col span={24}>
            <Table
              size="small"
              rowKey="id"
              loading={loading}
              columns={columns}
              dataSource={list}
              onChange={handleStandardTableChange}
              pagination={{
                total,
                showSizeChanger: true,
                pageSizeOptions: ['10', '20', '50', '100', '200'],
                showQuickJumper: true,
                showTotal: (total: number, range: number[]) => `第 ${range[0]}-${range[1]}条， 共 ${total}条`,
              }}
            />
          </Col>
        </Row>
      </Card>

      <UsersForm
        onSubmit={async (value) => {
          if (!access.canAdmin) {message.error('操作权限受限，请联系平台管理员');return}
          const res = await handleUpdate(value);
          if (res.success) {
            did('');
            handleUpdateModalVisible(false);
            setCurrentRow(undefined);
          }
        }}
        onCancel={() => {
          handleUpdateModalVisible(false);
          setCurrentRow(undefined);
        }}
        updateModalVisible={updateModalVisible}
        values={currentRow || {}}
      />

    </PageContainer>
  );
};
export default UserManager;
