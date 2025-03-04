import { PlusOutlined, FormOutlined, DeleteOutlined } from '@ant-design/icons';
import { Button, Divider, message, Popconfirm, Select } from 'antd';
import React, { useState, useRef, useEffect } from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import ProTable, { ProColumns, ActionType } from '@ant-design/pro-table';
import CreateForm from './components/CreateForm';
import UpdateForm from './components/UpdateForm';
import { TableListItem } from './data.d';
import { queryHost, updateHost, addHost, removeHost } from './service';
import { useAccess, FormattedMessage } from 'umi';


const { Option } = Select;
/**
 * 添加节点
 * @param fields
 */
const handleAdd = async (fields: TableListItem) => {
  const hide = message.loading('正在添加');
  try {
    await addHost({
      ...fields,
      online: Number(fields.online),
    });
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
    await updateHost({
      ...fields,
      "id": id,
      online: Number(fields.online),
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
    await removeHost({
      "id": id,
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

const formInitValue = { "id": 0, "idc_id": "", "env_id": "", "ip_address": "", "hostname": "", "online": "1", "description": "" }

const TableList: React.FC<{}> = () => {
  const [createModalVisible, handleModalVisible] = useState<boolean>(false);
  const [updateModalVisible, handleUpdateModalVisible] = useState<boolean>(false);
  const [formValues, setFormValues] = useState(formInitValue);
  const actionRef = useRef<ActionType>();
  const [idcList, setIdcList] = useState([{ "id": 0, "idc_name": "" }]);
  const [envList, setEnvList] = useState([{ "id": 0, "env_name": "" }]);
  const [idcEnum, setIdcEnum] = useState<{}>({})
  const [envEnum, setEnvEnum] = useState<{}>({})
  const access = useAccess();

  useEffect(() => {

    fetch('/api/v1/meta/idc/list')
      .then((response) => response.json())
      .then((json) => {
        setIdcList(json.data);
        const valueDict: { [key: number]: string } = {}
        json.data.forEach(record => { valueDict[record.id] = record.idc_name });
        setIdcEnum(valueDict)
      })
      .catch((error) => {
        console.log('Fetch idc list failed', error);
      });

    fetch('/api/v1/meta/env/list')
      .then((response) => response.json())
      .then((json) => {
        setEnvList(json.data);
        const valueDict: { [key: number]: string } = {}
        json.data.forEach(record => { valueDict[record.id] = record.env_name });
        setEnvEnum(valueDict)
      })
      .catch((error) => {
        console.log('Fetch env list failed', error);
      });

  }, []);

  const columns: ProColumns<TableListItem>[] = [
    {
      title: 'IP地址',
      dataIndex: 'ip_address',
      initialValue: formValues.ip_address,
      sorter: true,
      formItemProps: {
        rules: [
          {
            required: true,
            message: <FormattedMessage id="pages.searchTable.form.requireItem" />,
          },
        ],
      },
    },
    {
      title: '主机名',
      dataIndex: 'hostname',
      initialValue: formValues.hostname,
      search: false,
    },
    {
      title: '备注信息',
      dataIndex: 'description',
      initialValue: formValues.description,
      search: false,
    },
    {
      title: <FormattedMessage id="pages.searchTable.column.ic" />,
      dataIndex: 'idc_id',
      filters: true,
      onFilter: true,
      sorter: true,
      initialValue: formValues.idc_id,
      formItemProps: {
        rules: [
          {
            required: true,
            message: <FormattedMessage id="pages.searchTable.form.requireItem" />,
          },
        ],
      },
      valueEnum: idcEnum,
      renderFormItem: (_, { type, defaultRender, ...rest }, form) => {
        return <Select>
          {idcList && idcList.map(item => <Option key={item.id} value={item.id}>{item.idc_name}</Option>)}
        </Select>
      },

    },
    {
      title: <FormattedMessage id="pages.searchTable.column.env" />,
      dataIndex: 'env_id',
      filters: true,
      onFilter: true,
      sorter: true,
      initialValue: formValues.env_id,
      formItemProps: {
        rules: [
          {
            required: true,
            message: <FormattedMessage id="pages.searchTable.form.requireItem" />,
          },
        ],
      },
      valueEnum: envEnum,
      renderFormItem: (_, { type, defaultRender, ...rest }, form) => {
        return <Select>
          {envList && envList.map(item => <Option key={item.id} value={item.id}>{item.env_name}</Option>)}
        </Select>
      },

    },
    {
      title: '上线',
      dataIndex: 'online',
      filters: true,
      onFilter: true,
      valueEnum: {
        0: { text: '否', status: 'Error' },
        1: { text: '是', status: 'Success' },
      },
      sorter: true,
      search: false,
      initialValue: String(formValues.online),
      formItemProps: {
        rules: [
          {
            required: true,
            message: <FormattedMessage id="pages.searchTable.form.requireItem" />,
          },
        ],
      },
    },
    {
      title: <FormattedMessage id="pages.searchTable.column.gmtCreated" />,
      dataIndex: 'gmt_created',
      sorter: true,
      valueType: 'dateTime',
      hideInForm: true,
      search: false,
    },
    {
      title: <FormattedMessage id="pages.searchTable.column.gmtUpdated" />,
      dataIndex: 'gmt_updated',
      sorter: true,
      valueType: 'dateTime',
      hideInForm: true,
      search: false,
    },
    {
      title: <FormattedMessage id="pages.searchTable.column.operate" />,
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
            <FormOutlined /><FormattedMessage id="pages.searchTable.operate.edit" />
          </a>
          <Divider type="vertical" />
          <Popconfirm
            title={`确认要删除数据【${record.ip_address}】,删除后不可恢复，是否继续？`}
            placement={"left"}
            onConfirm={async () => {
              if (!access.canAdmin) { message.error('操作权限受限，请联系平台管理员'); return }
              const success = await handleRemove(record.id);
              if (success) {
                if (actionRef.current) {
                  actionRef.current.reload();
                }
              }
            }}
          >
            <a><DeleteOutlined /><FormattedMessage id="pages.searchTable.operate.delete" /></a>
          </Popconfirm>
        </>
      ),
    },
  ];

  return (
    <PageContainer>
      <ProTable<TableListItem>
        headerTitle={<FormattedMessage id="pages.searchTable.datalist" />}
        actionRef={actionRef}
        rowKey="id"
        search={{
          labelWidth: 120,
        }}
        toolBarRender={() => [
          <Button type="primary"
            onClick={() => {
              handleModalVisible(true);
              setFormValues(formInitValue);
            }}
          >
            <PlusOutlined />
            <FormattedMessage id="pages.searchTable.operate.create" />
          </Button>,
        ]}
        request={(params, sorter, filter) => queryHost({ ...params, sorter, filter })}
        columns={columns}
        pagination={{
          pageSize: 10,
        }}
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
          rowKey="id"
          type="form"
          columns={columns}
        />
      </CreateForm>

      <UpdateForm onCancel={() => handleUpdateModalVisible(false)} updateModalVisible={updateModalVisible}>
        <ProTable<TableListItem, TableListItem>
          onSubmit={async (value) => {
            if (!access.canAdmin) { message.error('操作权限受限，请联系平台管理员'); return }
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
