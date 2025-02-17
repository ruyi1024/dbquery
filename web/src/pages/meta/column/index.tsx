import React, { useState, useRef, useEffect } from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import ProTable, { ProColumns, ActionType } from '@ant-design/pro-table';
import { TableListItem } from './data.d';
import { queryColumn } from './service';


const TableList: React.FC<{}> = () => {

  const actionRef = useRef<ActionType>();

  const columns: ProColumns<TableListItem>[] = [

    {
      title: '字段名',
      dataIndex: 'column_name',
      sorter: true,
    },
    {
      title: '数据类型',
      dataIndex: 'data_type',
      hideInSearch: true,
    },
    {
      title: '允许为空',
      dataIndex: 'is_nullable',
      hideInSearch: true,
    },
    {
      title: '默认值',
      dataIndex: 'default_value',
      hideInSearch: true,
    },
    {
      title: '字段备注',
      dataIndex: 'column_comment',
      hideInSearch: true,
    },
    {
      title: '字符集',
      dataIndex: 'characters',
      hideInSearch: true,
    },
    {
      title: '所属表',
      dataIndex: 'table_name',
      sorter: true,
    },
    {
      title: '所属库',
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
        headerTitle="数据字段列表"
        actionRef={actionRef}
        rowKey="id"
        search={{
          labelWidth: 120,
        }}
        request={(params, sorter, filter) => queryColumn({ ...params, sorter, filter })}
        columns={columns}
        pagination={{
          pageSize: 10,
        }}
      />
    </PageContainer>
  );
};

export default TableList;
