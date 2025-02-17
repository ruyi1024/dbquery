import React, { useEffect, useState } from 'react';
import { PageContainer } from '@ant-design/pro-layout';
import { Row, Col, Card, Alert, message, Tooltip, Table, Space } from 'antd';
import { InfoCircleOutlined, SmileTwoTone, PieChartTwoTone, ProfileTwoTone, SoundTwoTone } from '@ant-design/icons';
import styles from './index.less';
import { ChartCard, MiniArea, MiniBar } from './components/Charts';
import Trend from './components/Trend';
import { Gauge } from '@ant-design/charts';
import PieChart from '@/components/Chart/PieChart';
import { StatisticCard } from '@ant-design/pro-components';
import moment from "moment";

const { Divider } = StatisticCard;



const demoPie = [
  {
    type: '分类一',
    value: 27,
  },
  {
    type: '分类二',
    value: 25,
  },
  {
    type: '分类三',
    value: 18,
  },
  {
    type: '分类四',
    value: 15,
  },
  {
    type: '分类五',
    value: 10,
  },
  {
    type: '其他',
    value: 5,
  },
];

export default (): React.ReactNode => {
  const [dashboardData, setDashboardData] = useState<any>([]);
  const [wsState, setWsState] = useState<boolean>(false);
  const [seconds, setSeconds] = useState<number>(1);
  const [lastTime, setLastTime] = useState<any>(new Date());
  const [loading, setLoading] = useState<boolean>(true);
  const [eventList, setEventList] = useState<any>([]);
  const [alarmList, setAlarmList] = useState<any>([]);
  const [datasourcePieData, setDatasourcePieData] = useState<any>([{ type: 'noData', value: 1 }]);
  const [databasePieData, setDatabasePieData] = useState<any>([{ type: 'noData', value: 1 }]);
  const [tablePieData, setTablePieData] = useState<any>([{ type: 'noData', value: 1 }]);
  const [columnPieData, setColumnPieData] = useState<any>([{ type: 'noData', value: 1 }]);

  // const columns_event = [
  //   {
  //     title: '事件时间',
  //     dataIndex: 'event_time',
  //   },
  //   {
  //     title: '事件类型',
  //     dataIndex: 'event_type',
  //   },
  //   {
  //     title: '事件组',
  //     dataIndex: 'event_group',
  //   },
  //   {
  //     title: '事件实体',
  //     dataIndex: 'event_entity',
  //   },
  //   {
  //     title: '事件指标',
  //     dataIndex: 'event_key',
  //   },
  //   {
  //     title: '事件值',
  //     dataIndex: 'event_value',
  //     render: (_: any, record: any) => <>{record.event_value}</>,
  //   },
  // ];

  // const columns_alarm = [
  //   {
  //     title: '告警时间',
  //     dataIndex: 'event_time',
  //   },
  //   {
  //     title: '告警信息',
  //     dataIndex: 'alarm_title',
  //   },
  //   {
  //     title: '告警级别',
  //     dataIndex: 'alarm_level',
  //   },
  //   {
  //     title: '事件类型',
  //     dataIndex: 'event_type',
  //   },
  //   {
  //     title: '事件实体',
  //     dataIndex: 'event_entity',
  //   },
  // ];

  useEffect(() => {
    try {
      fetch(`/api/v1/meta/dashboard/info`)
        .then((response) => response.json())
        .then((json) => {
          console.info(json.data);
          return (
            setDashboardData(json.data),
            setDatasourcePieData(json.data.datasourcePieDataList),
            setDatabasePieData(json.data.databasePieDataList),
            setTablePieData(json.data.tablePieDataList),
            setColumnPieData(json.data.columnPieDataList),
            setLoading(false)
          );
        })
        .catch((error) => {
          console.log('fetch dashboard data failed', error);
        });
    } catch (e) {
      message.error(`get data error. ${e}`)
    }
  }, []);

  return (
    <PageContainer>
      <Row gutter={[16, 24]} style={{ marginTop: '10px' }}>
        <Col span={24}>
          <StatisticCard.Group>
            <StatisticCard
              statistic={{
                title: '类型',
                value: dashboardData && dashboardData.datasourceTypeCount,
                status: 'processing',
              }}
            />
            <Divider />
            <StatisticCard
              statistic={{
                title: '机房',
                value: dashboardData && dashboardData.datasourceIdcCount,
                status: 'processing',
              }}
            />
            <Divider />
            <StatisticCard
              statistic={{
                title: '环境',
                value: dashboardData && dashboardData.datasourceEnvCount,
                status: 'processing',
              }}
            />
            <Divider />
            <StatisticCard
              statistic={{
                title: '实例数',
                value: dashboardData && dashboardData.datasourceCount,
                status: 'processing',
              }}
            />
            <Divider />
            <StatisticCard
              statistic={{
                title: '数据库',
                value: dashboardData && dashboardData.databaseCount,
                status: 'warning',
              }}
            />
            <StatisticCard
              statistic={{
                title: '数据表',
                value: dashboardData && dashboardData.tableCount,
                status: 'success',
              }}
            />
            <StatisticCard
              statistic={{
                title: '字段',
                value: dashboardData && dashboardData.columnCount,
                status: 'error',
              }}
            />
          </StatisticCard.Group>
        </Col>
      </Row>

      <Row gutter={[16, 24]} style={{ marginTop: '15px' }}>
        <Col span={12}>
          <Card title={<span><PieChartTwoTone />&nbsp;数据源实例分布</span>} bordered={false}>
            <PieChart data={datasourcePieData} loading={loading} height={330} />
          </Card>
        </Col>
        <Col span={12}>
          <Card title={<span><PieChartTwoTone />&nbsp;数据库分布</span>} bordered={false}>
            <PieChart data={databasePieData} loading={loading} height={330} />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 24]} style={{ marginTop: '15px' }}>
        <Col span={12}>
          <Card title={<span><PieChartTwoTone />&nbsp;数据表分布</span>} bordered={false}>
            <PieChart data={tablePieData} loading={loading} height={330} />
          </Card>
        </Col>
        <Col span={12}>
          <Card title={<span><PieChartTwoTone />&nbsp;数据字段分布</span>} bordered={false}>
            <PieChart data={columnPieData} loading={loading} height={330} />
          </Card>
        </Col>
      </Row>

      {/* <Row gutter={[16, 24]} style={{ marginTop: '15px' }}>
        <Col span={12}>
          <Card title={<span><ProfileTwoTone />&nbsp;最新数据库</span>} bordered={false} size="small">
            <Table
              columns={columns_event}
              loading={loading}
              dataSource={eventList}
              size="small"
              pagination={false}
            />
          </Card>
        </Col>
        <Col span={12}>
          <Card title={<span><SoundTwoTone />&nbsp;最新数据表</span>} bordered={false} size="small">
            <Table
              columns={columns_alarm}
              loading={loading}
              dataSource={alarmList}
              size="small"
              pagination={false}
            />
          </Card>
        </Col>
      </Row> */}
    </PageContainer>
  );
};
