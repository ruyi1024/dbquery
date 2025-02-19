import { PlusOutlined, FormOutlined, DeleteOutlined } from '@ant-design/icons';
import { Button, Divider, message, Popconfirm, Select } from 'antd';
import React, { useState, useRef } from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import ProTable, { ProColumns, ActionType } from '@ant-design/pro-table';
import { TableListItem } from './data.d';
import { query } from './service';
import { useAccess, FormattedMessage } from 'umi';

const TableList: React.FC<{}> = () => {
  const actionRef = useRef<ActionType>();

  const columns: ProColumns<TableListItem>[] = [
    {
      title: '任务标识',
      dataIndex: 'heartbeat_key',
      sorter: true,
    },
    {
      title: '心跳开始时间',
      dataIndex: 'heartbeat_time',
      sorter: true,
      valueType: 'dateTime',
      search: false,
    },
    {
      title: '心跳结束时间',
      dataIndex: 'heartbeat_end_time',
      sorter: true,
      valueType: 'dateTime',
      search: false,
    },
    {
      title: <FormattedMessage id="pages.searchTable.column.gmtCreated" />,
      dataIndex: 'gmt_created',
      sorter: true,
      valueType: 'dateTime',
      search: false,
    },
    {
      title: <FormattedMessage id="pages.searchTable.column.gmtUpdated" />,
      dataIndex: 'gmt_updated',
      sorter: true,
      valueType: 'dateTime',
      search: false,
    },

  ];

  return (
    <PageContainer>
      <ProTable<TableListItem>
        headerTitle={<FormattedMessage id="pages.searchTable.datalist" />}
        actionRef={actionRef}
        rowKey="id"
        request={(params, sorter, filter) => query({ ...params, sorter, filter })}
        columns={columns}
      />

    </PageContainer>
  );
};

export default TableList;
