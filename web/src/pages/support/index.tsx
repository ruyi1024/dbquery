
import { PageContainer } from '@ant-design/pro-components';
import { Collapse, Row, Col, Card } from 'antd';

const { Panel } = Collapse;


export default () => (
  <PageContainer>
    <Row gutter={[16, 24]} style={{ marginTop: '10px' }}>
      <Col span={16}>
        <Card>
          <Collapse>
            <Panel header="Lepus是什么？" key="1">
              Lepus,中文译为天兔，致力于打造简洁、智能、强大的开源企业级数据库监控系统，致力于数据库一站式监控管理，让数据库监控和运维管理更简单。 <a href="" target={"_blank"}>进入官网了解更多</a>
            </Panel>
            <Panel header="Lepus是开源免费的吗？" key="2">
              Lepus 是完全开源和免费的，您在遵守开源协议和Lepus规范的前提下，可以免费使用。
            </Panel>
            <Panel header="Lepus是什么开源协议？" key="3">
              Lepus采用的开源协议为GPLV3，您可以通过https://www.gnu.org/licenses/gpl-3.0.html 获取完整协议内容。特别注意：您可以下载使用源代码用于内部学习研究，但是禁止任何形式的商业行为，包括但不限于二次开发的商业行为。
            </Panel>
            <Panel header="使用Lepus需要注意什么？" key="4">
              Lepus是开源免费产品，Lepus开发团队和使用者无合同和责任关系，Lepus团队不承担因产品或者使用问题造成的任何损失。
            </Panel>
            <Panel header="如何加入Lepus开发贡献代码？" key="5">
              请在我们的项目Git库提交Pull Requests即可：https://gitee.com/lepus-group
            </Panel>
            <Panel header="如何获得社区帮助？" key="6">
              1.参考官网网站和Git仓库的文档和手册（优先）;
              2.加入Lepus微信社区群（推荐）,添加作者微信 Andy_Ruyi 后邀请入群（加微信请备注Lepus加微信群）;
              3.加入QQ交流群沟通和解决问题，QQ群号码：149648217 。
            </Panel>
            <Panel header="如何联络Lepus团队？" key="8">
              除社区文档，QQ、微信群，Lepus开发团队不提供其他技术类的免费支持服务。如有合作需求，请发邮件到：ruyi@139.com，或添加微信：Andy_Ruyi (需备注目的)。
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
