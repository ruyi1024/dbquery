import { PlusOutlined,FormOutlined,DeleteOutlined } from '@ant-design/icons';
import { Button, Divider, message, Popconfirm, Select } from 'antd';
import React, {useState, useRef, useEffect} from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import ProTable, { ProColumns, ActionType } from '@ant-design/pro-table';
import CreateForm from './components/CreateForm';
import UpdateForm from './components/UpdateForm';
import { TableListItem } from './data.d';
import { queryCluster, updateCluster, addCluster, removeCluster } from './service';
import { useAccess } from 'umi';
const { Option } = Select;
/**
 * 添加节点
 * @param fields
 */
const handleAdd = async (fields: TableListItem) => {
  const hide = message.loading('正在添加');
  try {
    await addCluster({ ...fields });
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
const handleUpdate = async (fields: FormValueType,id: number) => {
  const hide = message.loading('正在配置');
  try {
    await updateCluster({
      ...fields,
      "id":id,
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
    await removeCluster({
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

const formInitValue = {"id":0,"cluster_name":"","module_id":"","description":""}

const TableList: React.FC<{}> = () => {
  const [createModalVisible, handleModalVisible] = useState<boolean>(false);
  const [updateModalVisible, handleUpdateModalVisible] = useState<boolean>(false);
  const [formValues, setFormValues] = useState(formInitValue);
  const actionRef = useRef<ActionType>();
  const access = useAccess();
  const [moduleList,setModuleList] = useState<any[]>([{"id":0,"module_name":""}]);
  const [moduleEnum,setModuleEnum] = useState<{}>({})

  useEffect(() => {
    fetch('/api/v1/meta/module/list')
      .then((response) => response.json())
      .then((json) => {
        setModuleList(json.data);
        const valueDict : {[key:number]:string}  = {}
        json.data.forEach(record=>{ valueDict[record.id]=record.module_name });
        setModuleEnum(valueDict)
      })
      .catch((error) => {
        console.log('Fetch module list failed', error);
      });
  }, []);

  const columns: ProColumns<TableListItem>[] = [

    {
      title: '集群名',
      dataIndex: 'cluster_name',
      initialValue: formValues.cluster_name,
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
      title: '集群描述',
      dataIndex: 'description',
      initialValue: formValues.description,
      search:false,
    },
    {
      title: '集群类型',
      dataIndex: 'module_id',
      filters: true,
      onFilter: true,
      sorter: true,
      initialValue: formValues.module_id,
      formItemProps: {
        rules: [
          {
            required: true,
            message: '此项为必填项',
          },
        ],
      },
      hideInTable:false,
      renderFormItem: (_, {type, defaultRender, ...rest}, form) => {
        return <Select>
          {moduleList && moduleList.map(item => <Option key={item.id} value={item.id}>{item.module_name}</Option>)}
        </Select>
      },
      valueEnum: moduleEnum,
    },
    {
      title: '创建时间',
      dataIndex: 'gmt_created',
      sorter: true,
      valueType: 'dateTime',
      hideInForm: true,
      search:false,
    },
    {
      title: '修改时间',
      dataIndex: 'gmt_updated',
      sorter: true,
      valueType: 'dateTime',
      hideInForm: true,
      search:false,
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
            title={`确认要删除数据【${record.cluster_name}】,删除后不可恢复，是否继续？`}
            placement={"left"}
            onConfirm={async ()=>{
              if (!access.canAdmin) {message.error('操作权限受限，请联系平台管理员');return}
              const success = await handleRemove(record.id);
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
        search={{
          labelWidth: 120,
          //filterType: 'light',
        }}
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
        request={(params, sorter, filter) => queryCluster({ ...params, sorter, filter })}
        columns={columns}
        pagination={{
          pageSize: 10,
        }}
      />

      <CreateForm onCancel={() => handleModalVisible(false)} modalVisible={createModalVisible}>
        <ProTable<TableListItem, TableListItem>
          onSubmit={async (value) => {
            if (!access.canAdmin) {message.error('操作权限受限，请联系平台管理员');return}
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
            if (!access.canAdmin) {message.error('操作权限受限，请联系平台管理员');return}
            const success = await handleUpdate(value,formValues.id);
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
