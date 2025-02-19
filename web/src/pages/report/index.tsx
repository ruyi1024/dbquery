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
import { FormattedMessage } from 'umi';

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
                title: <FormattedMessage id="pages.searchTable.column.type" />,
                value: dashboardData && dashboardData.datasourceTypeCount,
                status: 'processing',
              }}
            />
            <Divider />
            <StatisticCard
              statistic={{
                title: <FormattedMessage id="pages.searchTable.column.idc" />,
                value: dashboardData && dashboardData.datasourceIdcCount,
                status: 'processing',
              }}
            />
            <Divider />
            <StatisticCard
              statistic={{
                title: <FormattedMessage id="pages.searchTable.column.env" />,
                value: dashboardData && dashboardData.datasourceEnvCount,
                status: 'processing',
              }}
            />
            <Divider />
            <StatisticCard
              statistic={{
                title: <FormattedMessage id="pages.report.instance" />,
                value: dashboardData && dashboardData.datasourceCount,
                status: 'processing',
              }}
            />
            <Divider />
            <StatisticCard
              statistic={{
                title: <FormattedMessage id="pages.report.database" />,
                value: dashboardData && dashboardData.databaseCount,
                status: 'warning',
              }}
            />
            <StatisticCard
              statistic={{
                title: <FormattedMessage id="pages.report.table" />,
                value: dashboardData && dashboardData.tableCount,
                status: 'success',
              }}
            />
            <StatisticCard
              statistic={{
                title: <FormattedMessage id="pages.report.column" />,
                value: dashboardData && dashboardData.columnCount,
                status: 'error',
              }}
            />
          </StatisticCard.Group>
        </Col>
      </Row>

      <Row gutter={[16, 24]} style={{ marginTop: '15px' }}>
        <Col span={12}>
          <Card title={<span><PieChartTwoTone />&nbsp;<FormattedMessage id="pages.report.datasourceDistribution" /></span>} bordered={false}>
            <PieChart data={datasourcePieData} loading={loading} height={330} />
          </Card>
        </Col>
        <Col span={12}>
          <Card title={<span><PieChartTwoTone />&nbsp;<FormattedMessage id="pages.report.databaseDistribution" /></span>} bordered={false}>
            <PieChart data={databasePieData} loading={loading} height={330} />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 24]} style={{ marginTop: '15px' }}>
        <Col span={12}>
          <Card title={<span><PieChartTwoTone />&nbsp;<FormattedMessage id="pages.report.tableDistribution" /></span>} bordered={false}>
            <PieChart data={tablePieData} loading={loading} height={330} />
          </Card>
        </Col>
        <Col span={12}>
          <Card title={<span><PieChartTwoTone />&nbsp;<FormattedMessage id="pages.report.columnDistribution" /></span>} bordered={false}>
            <PieChart data={columnPieData} loading={loading} height={330} />
          </Card>
        </Col>
      </Row>


    </PageContainer>
  );
};
