import { PageContainer } from '@ant-design/pro-components';
import { Card, Tabs, Form, Select, Button, Input, message, Table, Alert, Space } from 'antd';

import React, { useState, useEffect,useRef } from 'react';

import styles from './index.less';

const { TabPane } = Tabs;
const { TextArea } = Input;


const onChange = (key: string) => {
  console.log(key);
};

const Index: React.FC = () => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState<boolean>(false);
  const [pageSize, setPageSize] = useState<number>(20);
  const [current, setCurrent] = useState<number>(1);
  const [total, setTotal] = useState<number>();

  const [tableDataList, setTableDataList] = useState<any>();
  const [tableDataColumn, setTableDataColumn] = useState<any>();
  const [tableDataSuccess, setTableDataSuccess] = useState<boolean>(false);
  const [tableDataMsg, setTableDataMsg] = useState<any>("");
  const [queryTimes, setQueryTimes] = useState<number>(0);

  const [formValues, setFormValues] = useState({
    daasource: "",
    sql: "",
  });
  const [dbSourceList, setDbSourceList] = useState([]);
  const [tbSourceList, setTbSourceList] = useState([]);
  const [tbKeyword, setTbKeyword] = useState<string>('');
  const [datasource, setDatasource] = useState<string>('');
  const [tablesource, setTablesource] = useState<string>('');
  const [sqlContent, setSqlContent] = useState<string>('');



  useEffect(() => {
    setTbKeyword("");
    setTbSourceList([]);
    fetch('/api/v1/meta/database/list?db_type=MySQL')
      .then((response) => response.json())
      .then((json) => setDbSourceList(json.data))
      .catch((error) => {
        console.log('fetch database list failed', error);
      });
  }, []);


  const didQueryTables = (val: string) => {
    setTbKeyword("");
    setTbSourceList([]);
    setDatasource(val);
    fetch('/api/v1/meta/table/list?db_type=MySQL&db_name=' + val)
      .then((response) => response.json())
      .then((json) => setTbSourceList(json.data))
      .catch((error) => {
        console.log('fetch table list failed', error);
      });
  };

  const didSetTable = (val: string) => {
    setTablesource(val);
    const sql = "select * from "+val+" limit 100"
    setSqlContent(sql);
    form.setFieldsValue({
      sql: sql
    });
  };

  const asyncFetch = (values: {}) => {
    console.info(values);
    const params = { ...values };
    const headers = new Headers();
    headers.append('Content-Type', 'application/json');
    fetch('/api/v1/query/mysql', {
      method: 'post',
      headers: headers,
      body: JSON.stringify(params),
    })
      .then((response) => response.json())
      .then((json) => {
        console.info(json.data)
        return (
          setTableDataSuccess(json.success),
          setTableDataMsg(json.msg),
          setTableDataList(json.data),
          setTableDataColumn(json.columns),
          setQueryTimes(json.times)
        );
      })
      .catch((error) => {
        console.log('fetch data failed', error);
      });
  };



  const onFinish = (fieldValue: []) => {
    console.info(fieldValue["sql"]);
    const values = {
      datasource: fieldValue["datasource"],
      sql: fieldValue["sql"],
    };
    setFormValues(values);
    asyncFetch(values);
  };

  const onFinishFailed = (errorInfo: any) => {
    console.info(errorInfo);
    message.error('查询失败');
  };

  const queryPost = (type: any)=>{
    const params = {"datasource":datasource,"tablesource":tablesource,"sql":sqlContent,"type":type };

    const headers = new Headers();
    headers.append('Content-Type', 'application/json');
    fetch('/api/v1/query/mysql', {
      method: 'post',
      headers: headers,
      body: JSON.stringify(params),
    })
      .then((response) => response.json())
      .then((json) => {
        console.info(json.data)
        return (
          setTableDataSuccess(json.success),
            setTableDataMsg(json.msg),
            setTableDataList(json.data),
            setTableDataColumn(json.columns)
        );
      })
      .catch((error) => {
        console.log('fetch data failed', error);
        message.error('查询失败');
      });
  }


  return (
    <PageContainer>
      <Alert message="MySQL查询变更平台是用于企业数据库查询和变更线上数据的安全管理工具，可用于线上数据查询和数据修改，支持权限配置、安全风险检查拦截、查询审计功能。" type="info" showIcon closable />
      <Card>

        <Form
          labelCol={{ span: 3 }}
          wrapperCol={{ span: 31 }}
          style={{ marginTop: 8 }}
          form={form}
          onFinish={onFinish}
          onFinishFailed={onFinishFailed}
          initialValues={{}}
          name={'sqlForm'}
        >

          <Form.Item
            name={'datasource'}
            label="选择数据库"
            rules={[{ required: true, message: '请选择MySQL数据库' }]}
          >
            <Select
              showSearch style={{ width: 260 }}
              placeholder="请选择MySQL数据库"
              onChange={(val) => {
                didQueryTables(val);
              }}
            >
              {dbSourceList && dbSourceList.map(item => <Option key={item.db_name + ":" + item.host + ":" + item.port} value={item.db_name + ":" + item.host + ":" + item.port}>{item.host + ":" + item.port + "/" + item.db_name}</Option>)}
            </Select>
          </Form.Item>

          <Form.Item
            name={'tablesource'}
            label="选择数据表"
            rules={[{ required: true, message: '请选择MySQL数据表' }]}
          >
            <Select showSearch style={{ width: 260 }} placeholder="请选择MySQL数据表" value={tbKeyword}
              onChange={(val) => {
                didSetTable(val);
              }}
            >
              {tbSourceList && tbSourceList.map(item => <Option key={item.table_name} value={item.table_name}>{item.table_name}</Option>)}
            </Select>
          </Form.Item>

          {sqlContent ?
            <Form.Item
              name={'sql'}
              label="输入查询SQL"
              rules={[{required: true, message: '请选择MySQL数据库'}]}
            >
              <TextArea
                autoSize={{minRows: 4, maxRows: 8}}
                defaultValue={sqlContent}
                value={sqlContent}
              />
            </Form.Item>
          : ""}
          <Form.Item wrapperCol={{ offset: 3, span: 16 }}>
            <Space>
              <Button type="primary" htmlType="submit">
                执行数据查询
              </Button>
              <Button type="primary" htmlType="button"  onClick={()=>queryPost("doExplain")}>
                查询执行计划
              </Button>
              <Button type="primary" htmlType="button"  onClick={()=>queryPost("showColumn")}>
                查询表结构
              </Button>
              <Button type="primary" htmlType="button"  onClick={()=>queryPost("showIndex")}>
                查询表索引
              </Button>
              <Button type="primary" htmlType="button"  onClick={()=>queryPost("showCreate")}>
                查询建表SQL
              </Button>
              <Button type="primary" htmlType="button"  onClick={()=>queryPost("showSize")}>
                查询表大小
              </Button>
              <Button type="primary" danger htmlType="button"  onClick={()=>queryPost("doDml")}>
                执行数据变更
              </Button>

            </Space>
          </Form.Item>
        </Form>
        {tableDataSuccess == false && tableDataMsg!="" &&
          <Alert type="error" message={"查询执行数据库错误：" + tableDataMsg} banner />
        }
        {tableDataSuccess == true &&
          <>
        <Alert type="success" message={"查询执行数据库成功，执行用时：" +queryTimes+"毫秒" } banner />
          <Table
            className={styles.tableStyle}
            dataSource={tableDataList}
            columns={tableDataColumn}
            size={'small'}
          />
          </>
        }


      </Card>

    </PageContainer>
  );
};

export default Index;

