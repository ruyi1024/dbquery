import { PageContainer } from '@ant-design/pro-components';
import { Card, Tabs, Form, Select, Button, Input, message } from 'antd';

import React, { useState, useEffect } from 'react';

import styles from './index.less';

const { TabPane } = Tabs;
const { TextArea } = Input;

const onChange = (key: string) => {
  console.log(key);
};

const Index: React.FC = () => {
  const [form] = Form.useForm();

  const [data, setData] = useState({
    connectionsChartList: [],
    threadsChartList: [],
    queriesChartList: [],
    slowQueriesChartList: [],
    bytesChartList: [],
    dmlChartList: [],
    trxChartList: [],
    innodbChartList: [],
  });

  const [dbSourceList, setDbSourceList] = useState([]);
  const [current, setCurrent] = useState("");

  useEffect(() => {
    setCurrent("sql");
    fetch('/api/v1/meta/node/list_search?module=mysql')
      .then((response) => response.json())
      .then((json) => setDbSourceList(json.data))
      .catch((error) => {
        console.log('fetch instances failed', error);
      });
  }, []);


  const asyncFetch = (values: {}) => {
    const params = { ...formValues, ...values };
    const headers = new Headers();
    headers.append('Content-Type', 'application/json');
    fetch('/api/v1/performance/mysql/chart', {
      method: 'post',
      headers: headers,
      body: JSON.stringify(params),
    })
      .then((response) => response.json())
      .then((json) => setData(json.data))
      .catch((error) => {
        console.log('fetch data failed', error);
      });
  };

  const onFinish = (fieldValue: []) => {
    const values = {
      host: fieldValue["instance"].split(":")[0],
      port: fieldValue["instance"].split(":")[1],
      start_time: fieldValue['time_range'][0].format('YYYY-MM-DD HH:mm:ss'),
      end_time: fieldValue['time_range'][1].format('YYYY-MM-DD HH:mm:ss'),
    };
    setFormValues(values);
    asyncFetch(values);
  };

  const onFinishFailed = (errorInfo: any) => {
    console.info(errorInfo);
    message.error('查询失败');
  };


  return (
    <PageContainer>
      <Card>
        <Tabs onChange={onChange} type="card">
          <TabPane tab="查询数据" key="sql">

            <Form
              labelCol={{ span: 4 }}
              wrapperCol={{ span: 20 }}
              style={{ marginTop: 8 }}
              form={form}
              onFinish={onFinish}
              onFinishFailed={onFinishFailed}
              initialValues={{}}
              name={'sqlForm'}
            >

              <Form.Item
                name={'database'}
                label="选择数据库"
                rules={[{ required: true, message: '请选择MySQL数据库' }]}
              >
                <Select showSearch style={{ width: 260 }} placeholder="请选择MySQL数据库">

                </Select>
              </Form.Item>

              <Form.Item
                name={'database'}
                label="输入查询SQL"
                rules={[{ required: true, message: '请选择MySQL数据库' }]}
              >
                <TextArea
                  placeholder="请输入需要查询的SQL语句"
                  autoSize={{ minRows: 4, maxRows: 8 }}
                />
              </Form.Item>

              <Form.Item wrapperCol={{ offset: 4, span: 16 }}>
                <Button type="primary" htmlType="submit">
                  查询提交
                </Button>
              </Form.Item>


            </Form>

          </TabPane>
          <TabPane tab="查询表结构" key="2">
            查询元数据
          </TabPane>
          <TabPane tab="查询表索引" key="3">
            查询表索引
          </TabPane>
          <TabPane tab="查询建表SQL" key="4">
            查询建表语句
          </TabPane>

        </Tabs>
      </Card>

    </PageContainer>
  );
};

export default Index;
