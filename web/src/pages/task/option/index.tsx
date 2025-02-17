import { PlusOutlined, FormOutlined, DeleteOutlined } from '@ant-design/icons';
import { Button, Divider, message, Popconfirm, Select } from 'antd';
import React, { useState, useRef } from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import ProTable, { ProColumns, ActionType } from '@ant-design/pro-table';
import CreateForm from './components/CreateForm';
import UpdateForm from './components/UpdateForm';
import { TableListItem } from './data.d';
import { query, update, add, remove } from './service';
import { useAccess } from 'umi';

/**
 * 添加节点
 * @param fields
 */
const handleAdd = async (fields: TableListItem) => {
  const hide = message.loading('正在添加');
  try {
    await add({ ...fields });
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
const handleUpdate = async (fields: FormValueType, task_key: string) => {
  const hide = message.loading('正在配置');
  try {
    await update({
      ...fields,
      "task_key": task_key,
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
const handleRemove = async (task_key: string) => {
  const hide = message.loading('正在删除');
  try {
    await remove({
      "task_key": task_key,
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

const formInitValue = { "task_key": "", "task_name": "", "task_description": "", "crontab": "", "enable": "" }

const TableList: React.FC<{}> = () => {
  const [createModalVisible, handleModalVisible] = useState<boolean>(false);
  const [updateModalVisible, handleUpdateModalVisible] = useState<boolean>(false);
  const [formValues, setFormValues] = useState(formInitValue);
  const actionRef = useRef<ActionType>();
  const access = useAccess();


  const columns: ProColumns<TableListItem>[] = [
    {
      title: '任务标识',
      dataIndex: 'task_key',
      initialValue: formValues.task_key,
      sorter: true,
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
      title: '任务名',
      dataIndex: 'task_name',
      initialValue: formValues.task_name,
      sorter: true,
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
      title: '任务描述',
      dataIndex: 'task_description',
      initialValue: formValues.task_description,
      sorter: false,
      search: false,
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
      title: '计划任务',
      dataIndex: 'crontab',
      initialValue: formValues.crontab,
      sorter: false,
      search: false,
    },
    {
      title: '启用',
      dataIndex: 'enable',
      filters: true,
      onFilter: true,
      valueEnum: {
        0: { text: '禁用', status: 'Default' },
        1: { text: '启用', status: 'Success' },
      },
      sorter: false,
      initialValue: formValues.enable,
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
            <FormOutlined />修改
          </a>
          <Divider type="vertical" />
          <Popconfirm
            title={`确认要删除数据【${record.task_name}】,删除后不可恢复，是否继续？`}
            placement={"left"}
            onConfirm={async () => {
              if (!access.canAdmin) { message.error('操作权限受限，请联系平台管理员'); return }
              const success = await handleRemove(record.task_key);
              if (success) {
                if (actionRef.current) {
                  actionRef.current.reload();
                }
              }
            }}
          >
            <a><DeleteOutlined />删除</a>
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
        search={true}
        toolBarRender={() => [
          <Button type="primary"
            onClick={() => {
              handleModalVisible(true);
              setFormValues(formInitValue);
            }}
          >
            <PlusOutlined /> 新建
          </Button>,
        ]}
        request={(params, sorter, filter) => query({ ...params, sorter, filter })}
        columns={columns}
      />

      <CreateForm onCancel={() => handleModalVisible(false)} modalVisible={createModalVisible}>
        <ProTable<TableListItem, TableListItem>
          onSubmit={async (value) => {
            if (!access.canAdmin) { message.error('操作权限受限，请联系平台管理员'); return }
            const success = await handleAdd(value);
            if (success) {
              handleModalVisible(false);
              if (actionRef.current) {
                actionRef.current.reload();
              }
            }
          }}
          rowKey="task_key"
          type="form"
          columns={columns}
        />
      </CreateForm>

      <UpdateForm onCancel={() => handleUpdateModalVisible(false)} updateModalVisible={updateModalVisible}>
        <ProTable<TableListItem, TableListItem>
          onSubmit={async (value) => {
            if (!access.canAdmin) { message.error('操作权限受限，请联系平台管理员'); return }
            const success = await handleUpdate(value, formValues.task_key);
            if (success) {
              handleUpdateModalVisible(false);
              if (actionRef.current) {
                actionRef.current.reload();
              }
            }
          }}
          rowKey="task_key"
          type="form"
          columns={columns}
        />
      </UpdateForm>
    </PageContainer>
  );
};

export default TableList;
