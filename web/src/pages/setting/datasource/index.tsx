import { PlusOutlined, FormOutlined, DeleteOutlined, QuestionCircleOutlined } from '@ant-design/icons';
import { Button, Divider, message, Popconfirm, Select, Tooltip, Badge } from 'antd';
import React, { useState, useRef, useEffect } from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import ProTable, { ProColumns, ActionType } from '@ant-design/pro-table';
import CreateForm from './components/CreateForm';
import UpdateForm from './components/UpdateForm';
import { TableListItem } from './data.d';
import { query, update, add, remove, check } from './service';
import { useAccess } from 'umi';

const tableProps = {
  layout: 'horizontal',
  formItemLayout: {
    labelCol: {
      xs: { span: 24 },
      sm: { span: 4 },
    },
    wrapperCol: {
      xs: { span: 24 },
      sm: { span: 20 },
    },
  },
}

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
const handleUpdate = async (fields: FormValueType, id: number) => {
  const hide = message.loading('正在配置');
  try {
    await update({
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
    await remove({
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

/*
检查连接
*/

const handleCheck = async (fields: TableListItem) => {
  const hide = message.loading('正在检查连接');
  try {
    await check({ ...fields });
    hide();
    return true;
  } catch (error) {
    hide();
    return false;
  }
};


const formInitValue = {
  idc: '',
  env: '',
  type: '',
  host: '',
  name: '',
  port: '',
  user: '',
  pass: '',
  dbid: '',
  enable: '',
  execute_enable: '',
  dbmeta_enable: '',
  sensitive_enable: '',
  monitor_enable: '',
  alarm_enable: '',
};

const TableList: React.FC<{}> = () => {
  const [createModalVisible, handleModalVisible] = useState<boolean>(false);
  const [updateModalVisible, handleUpdateModalVisible] = useState<boolean>(false);
  const [formValues, setFormValues] = useState(formInitValue);
  const actionRef = useRef<ActionType>();
  const access = useAccess();

  const [idcList, setIdcList] = useState<any[]>([{ id: 0, idc_key: '', idc_name: '' }]);
  const [idcEnum, setIdcEnum] = useState<{}>({});
  const [envList, setEnvList] = useState<any[]>([{ id: 0, env_key: '', env_name: '' }]);
  const [envEnum, setEnvEnum] = useState<{}>({});
  const [typeList, setTypeList] = useState<any[]>([{ id: 0, name: '' }]);
  const [typeEnum, setTypeEnum] = useState<{}>({});

  useEffect(() => {
    fetch('/api/v1/datasource_idc/list')
      .then((response) => response.json())
      .then((json) => {
        setIdcList(json.data);
        const valueDict: { [key: number]: string } = {};
        json.data.forEach((record: { id: string | number; idc_key: string; }) => {
          valueDict[record.id] = record.idc_key;
        });
        setIdcEnum(valueDict);
      })
      .catch((error) => {
        console.log('Fetch cluster list failed', error);
      });
    fetch('/api/v1/datasource_env/list')
      .then((response) => response.json())
      .then((json) => {
        setEnvList(json.data);
        const valueDict: { [key: number]: string } = {};
        json.data.forEach((record: { id: string | number; env_key: string; }) => {
          valueDict[record.id] = record.env_key;
        });
        setEnvEnum(valueDict);
      })
      .catch((error) => {
        console.log('Fetch cluster list failed', error);
      });
    fetch('/api/v1/datasource_type/list')
      .then((response) => response.json())
      .then((json) => {
        setTypeList(json.data);
        const valueDict: { [key: number]: string } = {};
        json.data.forEach((record: { id: string | number; name: string; }) => {
          valueDict[record.id] = record.name;
        });
        setTypeEnum(valueDict);
      })
      .catch((error) => {
        console.log('Fetch cluster list failed', error);
      });
  }, []);

  const columns: ProColumns<TableListItem>[] = [
    {
      title: '数据源',
      dataIndex: 'name',
      initialValue: formValues.name,
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
      title: '类型',
      dataIndex: 'type',
      initialValue: formValues.type,
      filters: true,
      onFilter: true,
      sorter: true,
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
            {typeList &&
              typeList.map((item) => (
                <Option key={item.name} value={item.name}>
                  {item.name}
                </Option>
              ))}
          </Select>
        );
      },
      valueEnum: typeEnum,
    },
    {
      title: '主机',
      dataIndex: 'host',
      initialValue: formValues.host,
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
      title: '用户',
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
      valueType: 'password',
    },
    {
      title: 'DBID',
      dataIndex: 'dbid',
      hideInTable: true,
      hideInSearch: true,
      initialValue: formValues.dbid,
    },
    {
      title: '机房',
      dataIndex: 'idc',
      initialValue: formValues.idc,
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
            {idcList &&
              idcList.map((item) => (
                <Option key={item.idc_key} value={item.idc_key}>
                  {item.idc_name}
                </Option>
              ))}
          </Select>
        );
      },
      valueEnum: idcEnum,
    },
    {
      title: '环境',
      dataIndex: 'env',
      initialValue: formValues.env,
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
            {envList &&
              envList.map((item) => (
                <Option key={item.env_key} value={item.env_key}>
                  {item.env_name}
                </Option>
              ))}
          </Select>
        );
      },
      valueEnum: envEnum,
    },

    {
      title: '启用',
      dataIndex: 'enable',
      filters: false,
      onFilter: false,
      valueEnum: {
        0: { text: '', status: 'Default' },
        1: { text: '', status: 'Success' },
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
      title: '查询',
      dataIndex: 'execute_enable',
      filters: false,
      onFilter: false,
      valueEnum: {
        0: { text: '', status: 'Default' },
        1: { text: '', status: 'Success' },
      },
      sorter: false,
      initialValue: formValues.execute_enable,
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
      title: '元数据',
      dataIndex: 'dbmeta_enable',
      filters: false,
      onFilter: false,
      valueEnum: {
        0: { text: '', status: 'Default' },
        1: { text: '', status: 'Success' },
      },
      sorter: false,
      initialValue: formValues.dbmeta_enable,
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
      title: '探敏',
      dataIndex: 'sensitive_enable',
      filters: false,
      onFilter: false,
      valueEnum: {
        0: { text: '', status: 'Default' },
        1: { text: '', status: 'Success' },
      },
      sorter: false,
      initialValue: formValues.sensitive_enable,
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
      title: '监控',
      dataIndex: 'monitor_enable',
      filters: false,
      onFilter: false,
      valueEnum: {
        0: { text: '', status: 'Default' },
        1: { text: '', status: 'Success' },
      },
      sorter: false,
      initialValue: formValues.monitor_enable,
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
      title: '告警',
      dataIndex: 'alarm_enable',
      filters: false,
      onFilter: false,
      valueEnum: {
        0: { text: '', status: 'Default' },
        1: { text: '', status: 'Success' },
      },
      sorter: false,
      initialValue: formValues.alarm_enable,
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
      title: '状态',
      dataIndex: 'status',
      filters: true,
      onFilter: true,
      render: (text: string, value: any) => {
        if (text == '1') {
          return <Tooltip title={value.status_text} ><Badge status={"success"} /></Tooltip>
        } else {
          return <Tooltip title={value.status_text} ><Badge status={"error"} /></Tooltip>
        }
      },
      hideInForm: true,

    },
    {
      title: '创建时间',
      dataIndex: 'gmt_created',
      sorter: false,
      valueType: 'dateTime',
      hideInForm: true,
      hideInTable: true,
      hideInSearch: true,
    },
    {
      title: '修改时间',
      dataIndex: 'gmt_updated',
      sorter: false,
      valueType: 'dateTime',
      hideInForm: true,
      hideInTable: true,
      hideInSearch: true,
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
        {...tableProps}
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
        request={(params, sorter, filter) => query({ ...params, sorter, filter })}
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
            //检查连接状态
            const checkConnection = await handleCheck(value);
            if (!checkConnection) {
              message.error('数据库连接检查失败，请检查数据源配置是否正确');
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
            //检查连接状态
            const checkConnection = await handleCheck(value);
            if (!checkConnection) {
              message.error('数据库连接检查失败，请检查数据源配置是否正确');
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
