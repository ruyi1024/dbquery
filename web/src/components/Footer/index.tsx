import { GithubOutlined } from '@ant-design/icons';
import { DefaultFooter } from '@ant-design/pro-components';

const Footer: React.FC = () => {
  const defaultMessage = 'LEPUS数据库服务平台， Power by Lepus.cc，Version:6.0';

  const currentYear = new Date().getFullYear();

  return (
    <DefaultFooter
      copyright={`2014-${currentYear} ${defaultMessage}`}
      links={[
        {
          key: 'Lepus-cc',
          title: 'Lepus官方网站',
          href: 'https://www.lepus.cc',
          blankTarget: true,
        },
        {
          key: 'github',
          title: <GithubOutlined />,
          href: 'https://gitee.com/lepus-group',
          blankTarget: true,
        },
        {
          key: 'discuss-lepus',
          title: 'Lepus交流社区',
          href: 'https://discuss.lepus.cc',
          blankTarget: true,
        },
      ]}
    />
  );
};

export default Footer;
