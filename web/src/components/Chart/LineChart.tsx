import React from 'react';
import { Line } from '@ant-design/charts';

interface LineChartProps {
  data: never[]
  unit: string
}


const LineChart: React.FC<LineChartProps> = (props) => {
  let { data, unit, } = props;
  unit = unit && unit != "" ? unit : ""
  //max = max && max != "" ? max : ""
  let config = {
    data: data,
    xField: 'time',
    yField: 'value',
    seriesField: 'category',
    yAxis: {
      label: {
        formatter: function formatter(v: any) {
          return ''.concat(v).replace(/\d{1,3}(?=(\d{3})+$)/g, function (s) {
            return ''.concat(s, ',');
          }) + unit;
        },
      },
      //max: "",
    },
    //color: ['#99CCFF', '#99CCCC', '#CCCCFF', '#FFCC99', '#FFCCCC'],
    slider: {
      start: 0,
      end: 1,
      trendCfg: { isArea: true },
    },
  };

  return (
    <Line {...config} />
  );
};

export default LineChart;
