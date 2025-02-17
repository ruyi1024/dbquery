import React from 'react';
import { Pie, G2 } from '@ant-design/plots';

interface PieChartProps {
  data: { type: string; value: number; }[]
  loading: boolean
  height: number
}


const PieChart: React.FC<PieChartProps> = (props) => {
  const G = G2.getEngine('canvas');
  let { data, loading, height } = props;
  //console.log(typeof (data));
  let config = {
    appendPadding: 10,
    data: data,
    angleField: 'value',
    colorField: 'type',
    radius: 0.8,
    label: {
      type: 'spider',
      labelHeight: 40,
      formatter: (data: { type: any; value: any; percent: number; }, mappingData: { color: any; }) => {
        const group = new G.Group({});
        group.addShape({
          type: 'circle',
          attrs: {
            x: 0,
            y: 0,
            width: 40,
            height: 50,
            r: 5,
            fill: mappingData.color,
          },
        });
        group.addShape({
          type: 'text',
          attrs: {
            x: 10,
            y: 8,
            text: `${data.type}`,
            fill: mappingData.color,
          },
        });
        group.addShape({
          type: 'text',
          attrs: {
            x: 0,
            y: 25,
            text: `共${data.value}个，占比${(Math.round((data.percent + Number.EPSILON) * 100))}%`,
            fill: 'rgba(0, 0, 0, 0.65)',
            fontWeight: 800,
          },
        });
        return group;
      },
    },
    interactions: [{ type: 'pie-legend-active' }, { type: 'element-active' }],
  };

  return (
    <Pie {...config} loading={loading} height={height} />
  );
};

export default PieChart;
