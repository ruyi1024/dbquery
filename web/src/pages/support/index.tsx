
import { PageContainer } from '@ant-design/pro-components';
import { Collapse, Row, Col, Card } from 'antd';

const { Panel } = Collapse;
import { FormattedMessage } from 'umi';

export default () => (
  <PageContainer>
    <Row gutter={[16, 24]} style={{ marginTop: '10px' }}>
      <Col span={16}>
        <Card>
          <Collapse>
            <Panel header={<FormattedMessage id="pages.support.title.dbquery" />} key="1">
              <FormattedMessage id="pages.support.title.dbquery.desc" />
            </Panel>
            <Panel header={<FormattedMessage id="pages.support.title.dbquery.commercial" />} key="2">
              <FormattedMessage id="pages.support.title.dbquery.commercial.desc" />
            </Panel>
            <Panel header={<FormattedMessage id="pages.support.title.dbquery.license" />} key="3">
              <FormattedMessage id="pages.support.title.dbquery.license.desc" />
            </Panel>
            <Panel header={<FormattedMessage id="pages.support.title.dbquery.attention" />} key="4">
              <FormattedMessage id="pages.support.title.dbquery.attention.desc" />
            </Panel>
            <Panel header={<FormattedMessage id="pages.support.title.dbquery.join" />} key="5">
              <FormattedMessage id="pages.support.title.dbquery.join.desc" />
            </Panel>
            <Panel header={<FormattedMessage id="pages.support.title.dbquery.help" />} key="6">
              <FormattedMessage id="pages.support.title.dbquery.help.desc" />
            </Panel>
            <Panel header={<FormattedMessage id="pages.support.title.dbquery.support" />} key="8">
              <FormattedMessage id="pages.support.title.dbquery.support.desc" />
            </Panel>

          </Collapse>
        </Card>
      </Col>
      <Col span={8}>
        <Card>

        </Card>
      </Col>
    </Row>
  </PageContainer>
);
