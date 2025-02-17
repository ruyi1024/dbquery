import { PlusOutlined,FormOutlined,DeleteOutlined } from '@ant-design/icons';
import { Button, Divider, message, Popconfirm, Select } from 'antd';
import React, {useState, useRef, useEffect} from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import ProTable, { ProColumns, ActionType } from '@ant-design/pro-table';
import CreateForm from './components/CreateForm';
import UpdateForm from './components/UpdateForm';
import { TableListItem } from './data.d';
import { queryWeb, updateWeb, addWeb, removeWeb } from './service';
import { useAccess } from 'umi';


const { Option } = Select;
/**
 * 添加节点
 * @param fields
 */
const handleAdd = async (fields: TableListItem) => {
  const hide = message.loading('正在添加');
  try {
    await addWeb({
      ...fields,
      monitor:Number(fields.monitor),
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
const handleUpdate = async (fields: FormValueType,id: number) => {
  const hide = message.loading('正在配置');
  try {
    await updateWeb({
      ...fields,
      "id":id,
      monitor:Number(fields.monitor),
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
    await removeWeb({
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

const formInitValue = {"id":0,"cluster_id":"","env_id":"","name":"","url":"","monitor":"1","method":"GET"}

const TableList: React.FC<{}> = () => {
  const [createModalVisible, handleModalVisible] = useState<boolean>(false);
  const [updateModalVisible, handleUpdateModalVisible] = useState<boolean>(false);
  const [formValues, setFormValues] = useState(formInitValue);
  const actionRef = useRef<ActionType>();
  const [clusterList,setClusterList] = useState([{"id":0,"cluster_name":""}]);
  const [envList,setEnvList] = useState([{"id":0,"env_name":""}]);
  const [clusterEnum,setClusterEnum] = useState<{}>({})
  const [envEnum,setEnvEnum] = useState<{}>({})
  const access = useAccess();

  useEffect(() => {

    fetch('/api/v1/meta/cluster/list')
      .then((response) => response.json())
      .then((json) => {
        setClusterList(json.data);
        const valueDict : {[key:number]:string}  = {}
        json.data.forEach(record=>{ valueDict[record.id]=record.cluster_name });
        setClusterEnum(valueDict)
      })
      .catch((error) => {
        console.log('Fetch idc list failed', error);
      });

    fetch('/api/v1/meta/env/list')
      .then((response) => response.json())
      .then((json) => {
        setEnvList(json.data);
        const valueDict : {[key:number]:string}  = {}
        json.data.forEach(record=>{ valueDict[record.id]=record.env_name });
        setEnvEnum(valueDict)
      })
      .catch((error) => {
        console.log('Fetch env list failed', error);
      });

  }, []);

  const columns: ProColumns<TableListItem>[] = [
    {
      title: '网站名',
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
      title: 'URL地址',
      dataIndex: 'url',
      initialValue: formValues.url,
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
      title: '请求方式',
      dataIndex: 'method',
      initialValue: formValues.method,
      search:false,
    },
    {
      title: '集群',
      dataIndex: 'cluster_id',
      filters: true,
      onFilter: true,
      sorter: true,
      initialValue: formValues.cluster_id,
      formItemProps: {
        rules: [
          {
            required: true,
            message: '此项为必填项',
          },
        ],
      },
      valueEnum: clusterEnum,
      renderFormItem: (_, {type, defaultRender, ...rest}, form) => {
        return <Select>
          {clusterList && clusterList.map(item => <Option key={item.id} value={item.id}>{item.cluster_name}</Option>)}
        </Select>
      },

    },
    {
      title: '环境',
      dataIndex: 'env_id',
      filters: true,
      onFilter: true,
      sorter: true,
      initialValue: formValues.env_id,
      formItemProps: {
        rules: [
          {
            required: true,
            message: '此项为必填项',
          },
        ],
      },
      valueEnum: envEnum,
      renderFormItem: (_, {type, defaultRender, ...rest}, form) => {
        return <Select>
          {envList && envList.map(item => <Option key={item.id} value={item.id}>{item.env_name}</Option>)}
        </Select>
      },

    },
    {
      title: '监控',
      dataIndex: 'monitor',
      filters: true,
      onFilter: true,
      valueEnum: {
        0: { text: '否', status: 'Error' },
        1: { text: '是', status: 'Success' },
      },
      sorter: true,
      search:false,
      initialValue: String(formValues.monitor),
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
            title={`确认要删除数据【${record.name}】,删除后不可恢复，是否继续？`}
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
        request={(params, sorter, filter) => queryWeb({ ...params, sorter, filter })}
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
