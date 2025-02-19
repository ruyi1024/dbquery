import React, { useEffect, useRef, useState } from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import { Button, Card, Col, Divider, Input, message, Popconfirm, Row, Space, Table, Tag, Tooltip } from 'antd';
import type { UserListData, UserListItem } from './data';
import { queryLog } from './service';
import { PlusOutlined, ReloadOutlined } from '@ant-design/icons';
import moment from 'moment';
import { ActionType } from "@ant-design/pro-table";
import { useAccess, FormattedMessage } from 'umi';

const { Search } = Input;
const handleSearchKeyword = async (val: string) => {
  console.log('search val:', val);
};

const query = async (params: string) => {
  try {
    return await queryLog(params);
  } catch (e) {
    return { success: false, msg: e }
  }
}



const Audit: React.FC = () => {
  const [list, setList] = useState<UserListData[]>([]);
  const [total, setTotal] = useState<number>(0);
  const [loading, setLoading] = useState<boolean>(false);
  const [keyword, setKeyword] = useState<string>();
  const [currentPage, setCurrentPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(10);
  const actionRef = useRef<ActionType>();
  const access = useAccess();


  const did = (params: any) => {
    setLoading(true);
    const data = {
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
      title: '记录时间',
      dataIndex: 'gmt_created',
      sorter: true,
      render: (text: string) => moment(text).format('YYYY-MM-DD HH:mm:ss'),
      with: "100",
    },
    {
      title: <FormattedMessage id="pages.searchTable.column.user" />,
      dataIndex: 'username',
      with: '60',
    },
    {
      title: '数据源',
      dataIndex: 'datasource_type',
      with: 80,
    },
    {
      title: '数据库',
      dataIndex: 'database',
    },
    {
      title: <FormattedMessage id="pages.searchTable.column.operate" />,
      dataIndex: 'query_type',
      sorter: true,
    },
    {
      title: 'SQL类型',
      dataIndex: 'sql_type',
      sorter: true,
      with: 80,
    },
    {
      title: '执行内容',
      dataIndex: 'content',
      sorter: false,
      ellipsis: {
        showTitle: false,
      },
      render: text => (
        <Tooltip placement='topLeft' title={text}>
          {text}
        </Tooltip>
      ),
      width: 200,
    },
    {
      title: '执行状态',
      dataIndex: 'status',
      sorter: true,
      with: 80,
    },
    {
      title: '执行结果',
      dataIndex: 'result',
      with: "250",
      ellipsis: {
        showTitle: false,
      },
      render: text => (
        <Tooltip placement='topLeft' title={text}>
          {text}
        </Tooltip>
      ),
      textWrap: 'word-break',
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
                did({ keyword: val });
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
            <Button type="link" icon={<PlusOutlined />} onClick={() => handleUpdateModalVisible(true)}>
              导出数据
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


    </PageContainer>
  );
};
export default Audit;
