import React, { useState, useRef, useEffect } from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import ProTable, { ProColumns, ActionType } from '@ant-design/pro-table';
import { TableListItem } from './data.d';
import { queryDatabase } from './service';
import { useAccess, FormattedMessage } from 'umi';

const TableList: React.FC<{}> = () => {

  const actionRef = useRef<ActionType>();

  const columns: ProColumns<TableListItem>[] = [

    {
      title: <FormattedMessage id="pages.searchTable.column.databaseName" />,
      dataIndex: 'database_name',
      sorter: true,
    },
    {
      title: <FormattedMessage id="pages.searchTable.column.characters" />,
      dataIndex: 'characters',
      hideInSearch: true,
    },
    {
      title: <FormattedMessage id="pages.searchTable.column.databaseType" />,
      dataIndex: 'datasource_type',
      sorter: true,
    },
    {
      title: <FormattedMessage id="pages.searchTable.column.host" />,
      dataIndex: 'host',
    },
    {
      title: <FormattedMessage id="pages.searchTable.column.port" />,
      dataIndex: 'port',
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

  ];

  return (
    <PageContainer>
      <ProTable<TableListItem>
        headerTitle={<FormattedMessage id="pages.searchTable.datalist" />}
        actionRef={actionRef}
        rowKey="id"
        search={{
          labelWidth: 120,
        }}
        request={(params, sorter, filter) => queryDatabase({ ...params, sorter, filter })}
        columns={columns}
        pagination={{
          pageSize: 10,
        }}
      />
    </PageContainer>
  );
};

export default TableList;
