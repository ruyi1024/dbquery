import { GithubOutlined } from '@ant-design/icons';
import { DefaultFooter } from '@ant-design/pro-components';
import { FormattedMessage } from 'umi';
const Footer: React.FC = () => {
  const defaultMessage = 'DBQuery， Power by DB-Query.com，Version:1.0';

  const currentYear = new Date().getFullYear();

  return (
    <DefaultFooter
      copyright={`2014-${currentYear} ${defaultMessage}`}
      links={[
        {
          key: 'Lepus-cc',
          title: <FormattedMessage id="layout.dbquerySite" />,
          href: 'https://db-query.com',
          blankTarget: true,
        },
        {
          key: 'github',
          title: <GithubOutlined />,
          href: 'https://github.com/ruyi1024/dbquery',
          blankTarget: true,
        },
        {
          key: 'discuss-lepus',
          title: <FormattedMessage id="layout.lepusSite" />,
          href: 'https://www.lepus.cc',
          blankTarget: true,
        },
      ]}
    />
  );
};

export default Footer;
