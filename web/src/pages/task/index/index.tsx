import { PageContainer } from '@ant-design/pro-layout';
import { LoadingOutlined, ReloadOutlined } from '@ant-design/icons';
import ProTable, { ProColumns, ActionType } from '@ant-design/pro-table';
import { StatisticCard } from '@ant-design/pro-components';
import { Button, Card } from 'antd';
import moment from 'moment';
import React, { useEffect, useState, useRef } from 'react';
import { TableListItem, TableListParams } from './data.d';
import { queryTask } from './service';


const { Divider } = StatisticCard;


const valueEnum = {
  0: 'close',
  1: 'running',
  2: 'online',
  3: 'error',
};

export type TableListItem = {
  key: number;
  name: string;
  status: string;
  updatedAt: number;
  createdAt: number;
  progress: number;
  money: number;
};




const columns: ProColumns<TableListItem>[] = [

  {
    title: '任务名',
    dataIndex: 'task_name',
    sorter: true,
  },

  {
    title: '任务组',
    dataIndex: 'category_name',
    search: false,
  },
  {
    title: '调度类型',
    dataIndex: 'schedule_type',
    valueEnum: {
      'crontab': { text: '定时任务' },
      'period': { text: '周期任务' },
      'manual': { text: '单次任务' },
    },
    search: false,
  },

  {
    title: '执行状态',
    dataIndex: 'status',
    valueEnum: {
      'waiting': { text: '等待运行', status: 'Default' },
      'success': { text: '执行成功', status: 'Success' },
      'failed': { text: '执行失败', status: 'Error' },
      'running': { text: '正在执行', status: 'Warning' },
    },
  },

  {
    title: '下次执行时间',
    dataIndex: 'next_time',
    search: false,
  },
  {
    title: '启用状态',
    dataIndex: 'enable',
    filters: true,
    onFilter: true,
    valueEnum: {
      '0': { text: '禁用', status: 'Error' },
      '1': { text: '启用', status: 'Success' },
    },
    sorter: true,
    search: false,
  },
];


const TableList: React.FC<{}> = () => {

  const actionRef = useRef<ActionType>();

  const [time, setTime] = useState(() => Date.now());
  const [polling, setPolling] = useState<number | undefined>(2000);
  const [data, setData] = useState({});
  const [taskCount, setTaskCount] = useState<number>(0);
  const [taskSuccessCount, setTaskSuccessCount] = useState<number>(0);
  const [taskRunningCount, setTaskRunningCount] = useState<number>(0);
  const [taskFailedCount, setTaskFailedCount] = useState<number>(0);
  const [taskWaitCount, setTaskWaitCount] = useState<number>(0);
  const [taskSuccessPct, setTaskSuccessPct] = useState<number>(0);


  const didQuery = async () => {
    try {
      didQuery
      const data = await queryTask();
      setData(data);
      setTaskCount(data.taskCount);
      setTaskSuccessCount(data.taskSuccessCount);
      setTaskRunningCount(data.taskRunningCount);
      setTaskFailedCount(data.taskFailedCount);
      setTaskWaitCount(data.taskWaitCount);
      setTaskSuccessPct(data.taskSuccessPct);
      return
    } catch (e) {
      return { success: false, msg: e }
    }
  }

  useEffect(() => {
    didQuery();
  });

  const timeAwait = (waitTime: number): Promise<void> =>
    new Promise((res) =>
      window.setTimeout(() => {
        res();
        didQuery();
      }, waitTime),
    );


  // @ts-ignore
  return (
    <PageContainer>
      <Card>
        <StatisticCard.Group>
          <StatisticCard
            statistic={{
              title: '今日执行成功率',
              tip: '近1小时任务执行成功率',
              value: taskSuccessPct + '%',
            }}
          />
          <Divider />
          <StatisticCard
            statistic={{
              title: '今日执行总数',
              value: taskCount,
              status: 'processing',
            }}
          />
          <StatisticCard
            statistic={{
              title: '等待执行',
              value: taskWaitCount,
              status: 'default',
            }}
          />
          <StatisticCard
            statistic={{
              title: '正在执行',
              value: taskRunningCount,
              status: 'warning',
            }}
          />
          <StatisticCard
            statistic={{
              title: '执行成功',
              value: taskSuccessCount,
              status: 'success',
            }}
          />
          <StatisticCard
            statistic={{
              title: '执行失败',
              value: taskFailedCount,
              status: 'error',
            }}
          />

        </StatisticCard.Group>
      </Card>

      <ProTable<TableListItem>
        actionRef={actionRef}
        rowKey="id"
        search={false}
        polling={polling || undefined}
        request={async () => {
          await timeAwait(1000);
          setTime(Date.now());
          return data;
        }}
        //request={(params, sorter, filter) => queryTask({ ...params, sorter, filter })}
        columns={columns}
        pagination={false}
        dateFormatter="string"
        headerTitle={`最新查询时间：${moment(time).format('YYYY-MM-DD HH:mm:ss')}`}
        toolBarRender={() => [
          <Button
            key="3"
            type="primary"
            onClick={() => {
              if (polling) {
                setPolling(undefined);
                return;
              }
              setPolling(2000);
            }}
          >
            {polling ? <LoadingOutlined /> : <ReloadOutlined />}
            {polling ? '停止自动刷新' : '开始自动刷新'}
          </Button>,
        ]}
      />

    </PageContainer>
  );
};

export default TableList;
