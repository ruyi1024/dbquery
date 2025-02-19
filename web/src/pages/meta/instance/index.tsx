import React, { useState, useRef, useEffect } from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import ProTable, { ProColumns, ActionType } from '@ant-design/pro-table';
import { TableListItem } from './data.d';
import { queryInstance } from './service';
import { useAccess, FormattedMessage } from 'umi';

const TableList: React.FC<{}> = () => {

  const actionRef = useRef<ActionType>();

  const columns: ProColumns<TableListItem>[] = [
    {
      title: <FormattedMessage id="pages.searchTable.column.type" />,
      dataIndex: 'type',
      sorter: true,
    },
    {
      title: <FormattedMessage id="pages.searchTable.column.name" />,
      dataIndex: 'name',
      sorter: true,
    },
    {
      title: <FormattedMessage id="pages.searchTable.column.host" />,
      dataIndex: 'host',
      sorter: true,
    },
    {
      title: <FormattedMessage id="pages.searchTable.column.port" />,
      dataIndex: 'port',
    },
    {
      title: <FormattedMessage id="pages.searchTable.column.role" />,
      dataIndex: 'role',
      filters: true,
      onFilter: true,
      valueEnum: {
        1: { text: <FormattedMessage id="pages.searchTable.column.master" /> },
        2: { text: <FormattedMessage id="pages.searchTable.column.slave" /> },
      },
    },

    {
      title: <FormattedMessage id="pages.searchTable.column.enable" />,
      dataIndex: 'enable',
      filters: true,
      onFilter: true,
      valueEnum: {
        0: { text: <FormattedMessage id="pages.searchTable.column.no" />, status: 'Default' },
        1: { text: <FormattedMessage id="pages.searchTable.column.yes" />, status: 'Success' },
      },
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
