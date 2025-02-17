import { Settings as LayoutSettings } from '@ant-design/pro-components';

const Settings: LayoutSettings & {
  pwa?: boolean;
  logo?: string;
} = {
  navTheme: 'light',
  //navTheme: 'dark',
  primaryColor: '#006699',
  layout: 'top',
  contentWidth: 'Fluid',
  fixedHeader: true,
  fixSiderbar: true,
  colorWeak: false,
  headerHeight: 48,
  //collapsed: true,
  splitMenus: true,
  title: 'DBQuery',
  pwa: false,
  logo: '/logo.png',
  iconfontUrl: 'https://www.lepus.cc',
};

export default Settings;

