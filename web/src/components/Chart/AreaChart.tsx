import React from 'react';
import { Area } from '@ant-design/charts';

interface AreaChartProps {
  data: never[];
  unit: string;
}

const AreaChart: React.FC<AreaChartProps> = (props) => {
  let { data, unit } = props;
  unit = unit && unit != '' ? unit : '';
  let config = {
    data: data,
    xField: 'time',
    yField: 'value',
    seriesField: 'category',
    yAxis: {
      label: {
        formatter: function formatter(v:any) {
          return (
            ''.concat(v).replace(/\d{1,3}(?=(\d{3})+$)/g, function (s) {
              return ''.concat(s, ',');
            }) + unit
          );
        },
      },
    },

    slider: {
      start: 0,
      end: 1,
      trendCfg: { isArea: true },
    },
  };

  return <Area {...config} />;
};

export default AreaChart;
