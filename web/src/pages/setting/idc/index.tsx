import { PlusOutlined, FormOutlined, DeleteOutlined } from '@ant-design/icons';
import { Button, Divider, message, Popconfirm, Select } from 'antd';
import React, { useState, useRef, useEffect } from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import ProTable, { ProColumns, ActionType } from '@ant-design/pro-table';
import CreateForm from './components/CreateForm';
import UpdateForm from './components/UpdateForm';
import { TableListItem } from './data.d';
import { queryIdc, updateIdc, addIdc, removeIdc } from './service';
import { useAccess, FormattedMessage } from 'umi';
/**
 * 添加节点
 * @param fields
 */
const handleAdd = async (fields: TableListItem) => {
  const hide = message.loading('正在添加');
  try {
    await addIdc({ ...fields });
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
    await updateIdc({
      ...fields,
      "id": id,
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
    await removeIdc({
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

const formInitValue = { "id": 0, "idck_key": "", "idc_name": "", "city": "", "description": "" }

const TableList: React.FC<{}> = () => {
  const [createModalVisible, handleModalVisible] = useState<boolean>(false);
  const [updateModalVisible, handleUpdateModalVisible] = useState<boolean>(false);
  const [formValues, setFormValues] = useState(formInitValue);
  const actionRef = useRef<ActionType>();
  const access = useAccess();


  const columns: ProColumns<TableListItem>[] = [
    {
      title: <FormattedMessage id="pages.searchTable.column.idcKey" />,
      dataIndex: 'idc_key',
      initialValue: formValues.idc_key,
      sorter: false,
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
      title: <FormattedMessage id="pages.searchTable.column.idcName" />,
      dataIndex: 'idc_name',
      initialValue: formValues.idc_name,
      sorter: false,
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
      title: <FormattedMessage id="pages.searchTable.column.city" />,
      dataIndex: 'city',
      initialValue: formValues.city,
      sorter: false,
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
      title: <FormattedMessage id="pages.searchTable.column.description" />,
      dataIndex: 'description',
      initialValue: formValues.description,
      sorter: false,
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
            title={`确认要删除数据【${record.idc_name}】,删除后不可恢复，是否继续？`}
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
        search={true}
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
        request={(params, sorter, filter) => queryIdc({ ...params, sorter, filter })}
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
