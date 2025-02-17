import React, { useState, useRef, useEffect } from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import ProTable, { ProColumns, ActionType } from '@ant-design/pro-table';
import { TableListItem } from './data.d';
import { queryInstance } from './service';


const TableList: React.FC<{}> = () => {

  const actionRef = useRef<ActionType>();

  const columns: ProColumns<TableListItem>[] = [
    {
      title: '实例类型',
      dataIndex: 'type',
      sorter: true,
    },
    {
      title: '实例名',
      dataIndex: 'name',
      sorter: true,
    },
    {
      title: '实例主机',
      dataIndex: 'host',
      sorter: true,
    },
    {
      title: '实例端口',
      dataIndex: 'port',
    },
    {
      title: '角色',
      dataIndex: 'role',
      filters: true,
      onFilter: true,
      valueEnum: {
        1: { text: '主' },
        2: { text: '备' },
      },
    },

    {
      title: '是否启用',
      dataIndex: 'enable',
      filters: true,
      onFilter: true,
      valueEnum: {
        0: { text: '否', status: 'Default' },
        1: { text: '是', status: 'Success' },
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

  ];

  return (
    <PageContainer>
      <ProTable<TableListItem>
        headerTitle="数据库实例列表"
        actionRef={actionRef}
        rowKey="id"
        search={{
          labelWidth: 120,
        }}
        request={(params, sorter, filter) => queryInstance({ ...params, sorter, filter })}
        columns={columns}
        pagination={{
          pageSize: 10,
        }}
      />
    </PageContainer>
  );
};

export default TableList;
