import { PlusOutlined, FormOutlined, DeleteOutlined } from '@ant-design/icons';
import { Button, Divider, message, Popconfirm, Select } from 'antd';
import React, { useState, useRef, useEffect } from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import ProTable, { ProColumns, ActionType } from '@ant-design/pro-table';
import CreateForm from './components/CreateForm';
import UpdateForm from './components/UpdateForm';
import { TableListItem } from './data.d';
import { queryNode, updateNode, addNode, removeNode } from './service';
import { useAccess } from 'umi';

/**
 * 添加节点
 * @param fields
 */
const handleAdd = async (fields: TableListItem) => {
  const hide = message.loading('正在添加');
  try {
    await addNode({ ...fields });
    hide();
    message.success('添加成功');
    return true;
  } catch (error) {
    hide();
    message.error('添加失败请重试！');
    return false;
  }
};

/**
 * 更新节点
 * @param fields
 */
const handleUpdate = async (fields: FormValueType, id: number) => {
  const hide = message.loading('正在配置');
  try {
    await updateNode({
      ...fields,
      id: id,
    });
    hide();
    message.success('修改成功');
    return true;
  } catch (error) {
    hide();
    message.error('修改失败请重试！');
    return false;
  }
};

/**
 *  删除节点
 * @param selectedRows
 */
const handleRemove = async (id: number) => {
  const hide = message.loading('正在删除');
  try {
    await removeNode({
      id: id,
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

const formInitValue = {
  cluster_id: '',
  ip: '',
  domain: '',
  port: '',
  user: '',
  pass: '',
  dbid: '',
  role: '',
  monitor: '',
};

const TableList: React.FC<{}> = () => {
  const [createModalVisible, handleModalVisible] = useState<boolean>(false);
  const [updateModalVisible, handleUpdateModalVisible] = useState<boolean>(false);
  const [formValues, setFormValues] = useState(formInitValue);
  const actionRef = useRef<ActionType>();
  const access = useAccess();

  const [clusterList, setClusterList] = useState<any[]>([{ id: 0, cluster_name: '' }]);
  const [clusterEnum, setClusterEnum] = useState<{}>({});
  const [hostList, setHostList] = useState<any[]>([{ ip: '' }]);
  const [hostEnum, setHostEnum] = useState<{}>({});

  useEffect(() => {
    fetch('/api/v1/meta/cluster/list')
      .then((response) => response.json())
      .then((json) => {
        setClusterList(json.data);
        const valueDict: { [key: number]: string } = {};
        json.data.forEach((record) => {
          valueDict[record.id] = record.cluster_name;
        });
        setClusterEnum(valueDict);
      })
      .catch((error) => {
        console.log('Fetch cluster list failed', error);
      });

    fetch('/api/v1/meta/host/list')
      .then((response) => response.json())
      .then((json) => {
        setHostList(json.data);
        const valueDict: { [key: number]: string } = {};
        json.data.forEach((record) => {
          valueDict[record.ip_address] = record.ip_address;
        });
        setHostEnum(valueDict);
      })
      .catch((error) => {
        console.log('Fetch cluster list failed', error);
      });
  }, []);

  const columns: ProColumns<TableListItem>[] = [
    {
      title: '集群',
      dataIndex: 'cluster_id',
      initialValue: formValues.cluster_id,
      formItemProps: {
        rules: [
          {
            required: true,
            message: '此项为必填项',
          },
        ],
      },
      renderFormItem: (_, { type, defaultRender, ...rest }, form) => {
        return (
          <Select showSearch>
            {clusterList &&
              clusterList.map((item) => (
                <Option key={item.cluster_name} value={item.id}>
                  {item.cluster_name}
                </Option>
              ))}
          </Select>
        );
      },
      valueEnum: clusterEnum,
    },
    {
      title: '主机IP',
      dataIndex: 'ip',
      initialValue: formValues.ip,
      formItemProps: {
        rules: [
          {
            required: true,
            message: '此项为必填项',
          },
        ],
      },
      renderFormItem: (_, { type, defaultRender, ...rest }, form) => {
        return (
          <Select showSearch>
            {hostList &&
              hostList.map((item) => (
                <Option key={item.ip_address} value={item.ip_address}>
                  {item.ip_address}
                </Option>
              ))}
          </Select>
        );
      },
    },
    {
      title: '域名',
      dataIndex: 'domain',
      initialValue: formValues.domain,
    },
    {
      title: '端口',
      dataIndex: 'port',
      initialValue: formValues.port,
      formItemProps: {
        rules: [
          {
            required: true,
            message: '此项为必填项',
          },
        ],
      },
    },
    {
      title: '用户名',
      dataIndex: 'user',
      hideInTable: true,
      hideInSearch: true,
      initialValue: formValues.user,
    },
    {
      title: '密码',
      dataIndex: 'pass',
      hideInTable: true,
      hideInSearch: true,
      initialValue: '',
    },
    {
      title: 'DBID',
      dataIndex: 'dbid',
      hideInTable: true,
      hideInSearch: true,
      initialValue:  formValues.dbid,
    },
    {
      title: '角色',
      dataIndex: 'role',
      initialValue: formValues.role,
      formItemProps: {
        rules: [
          {
            required: true,
            message: '此项为必填项',
          },
        ],
      },
      renderFormItem: (_, { type, defaultRender, ...rest }, form) => {
        return (
          <Select>
            <Option key={1} value={1}>
              主
            </Option>
            <Option key={2} value={2}>
              备
            </Option>
          </Select>
        );
      },
      valueEnum: {
        1: { text: '主节点', status: 'Success' },
        2: { text: '备节点', status: 'Warning' },
      },
    },

    {
      title: '是否监控',
      dataIndex: 'monitor',
      filters: true,
      onFilter: true,
      valueEnum: {
        0: { text: '否', status: 'Default' },
        1: { text: '是', status: 'Success' },
      },
      sorter: true,
      initialValue: formValues.monitor,
      formItemProps: {
        rules: [
          {
            required: true,
            message: '此项为必填项',
          },
        ],
      },
      renderFormItem: (_, { type, defaultRender, ...rest }, form) => {
        return (
          <Select>
            <Option key={0} value={0}>
              否
            </Option>
            <Option key={1} value={1}>
              是
            </Option>
          </Select>
        );
      },
    },
    {
      title: '创建时间',
      dataIndex: 'gmt_created',
      sorter: true,
      valueType: 'dateTime',
      hideInForm: true,
      search: false,
    },
    {
      title: '修改时间',
      dataIndex: 'gmt_updated',
      sorter: true,
      valueType: 'dateTime',
      hideInForm: true,
      search: false,
    },
    {
      title: '操作',
      dataIndex: 'option',
      valueType: 'option',
      render: (_, record) => (
        <>
          <a
            onClick={() => {
              handleUpdateModalVisible(true);
              setFormValues(record);
            }}
          >
            <FormOutlined />
            修改
          </a>
          <Divider type="vertical" />
          <Popconfirm
            title={`确认要删除数据【${record.id}】,删除后不可恢复，是否继续？`}
            placement={'left'}
            onConfirm={async () => {
              if (!access.canAdmin) {
                message.error('操作权限受限，请联系平台管理员');
                return;
              }
              const success = await handleRemove(record.id);
              if (success) {
                if (actionRef.current) {
                  actionRef.current.reload();
                }
              }
            }}
          >
            <a>
              <DeleteOutlined />
              删除
            </a>
          </Popconfirm>
        </>
      ),
    },
  ];

  return (
    <PageContainer>
      <ProTable<TableListItem>
        headerTitle="数据列表"
        actionRef={actionRef}
        rowKey="id"
        search={{
          labelWidth: 120,
        }}
        toolBarRender={() => [
          <Button
            type="primary"
            onClick={() => {
              handleModalVisible(true);
              setFormValues(formInitValue);
            }}
          >
            <PlusOutlined /> 新建
          </Button>,
        ]}
        request={(params, sorter, filter) => queryNode({ ...params, sorter, filter })}
        columns={columns}
        pagination={{
          pageSize: 10,
        }}
      />

      <CreateForm onCancel={() => handleModalVisible(false)} modalVisible={createModalVisible}>
        <ProTable<TableListItem, TableListItem>
          onSubmit={async (value) => {
            if (!access.canAdmin) {
              message.error('操作权限受限，请联系平台管理员');
              return;
            }
            const success = await handleAdd(value);
            if (success) {
              handleModalVisible(false);
              if (actionRef.current) {
                actionRef.current.reload();
              }
            }
          }}
          rowKey="id"
          type="form"
          columns={columns}
        />
      </CreateForm>

      <UpdateForm
        onCancel={() => handleUpdateModalVisible(false)}
        updateModalVisible={updateModalVisible}
      >
        <ProTable<TableListItem, TableListItem>
          onSubmit={async (value) => {
            if (!access.canAdmin) {
              message.error('操作权限受限，请联系平台管理员');
              return;
            }
            const success = await handleUpdate(value, formValues.id);
            if (success) {
              handleUpdateModalVisible(false);
              if (actionRef.current) {
                actionRef.current.reload();
              }
            }
          }}
          rowKey="id"
          type="form"
          columns={columns}
        />
      </UpdateForm>
    </PageContainer>
  );
};

export default TableList;
