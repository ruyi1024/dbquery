import React, { useState, useRef, useEffect } from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import ProTable, { ProColumns, ActionType } from '@ant-design/pro-table';
import { TableListItem } from './data.d';
import { queryTable } from './service';


const TableList: React.FC<{}> = () => {

  const actionRef = useRef<ActionType>();

  const columns: ProColumns<TableListItem>[] = [

    {
      title: '数据表名',
      dataIndex: 'table_name',
      sorter: true,
    },
    {
      title: '表类型',
      dataIndex: 'table_type',
      hideInSearch: true,
    },
    {
      title: '表字符集',
      dataIndex: 'characters',
      hideInSearch: true,
    },
    {
      title: '表备注',
      dataIndex: 'table_comment',
      hideInSearch: true,
    },
    {
      title: '所属数据库',
      dataIndex: 'database_name',
      sorter: true,
    },
    {
      title: '数据库类型',
      dataIndex: 'datasource_type',
      sorter: true,
    },
    {
      title: '所属主机',
      dataIndex: 'host',
    },
    {
      title: '所属端口',
      dataIndex: 'port',
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
        headerTitle="数据表列表"
        actionRef={actionRef}
        rowKey="id"
        search={{
          labelWidth: 120,
        }}
        request={(params, sorter, filter) => queryTable({ ...params, sorter, filter })}
        columns={columns}
        pagination={{
          pageSize: 10,
        }}
      />
    </PageContainer>
  );
};

export default TableList;
