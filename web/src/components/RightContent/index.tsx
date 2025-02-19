import { Tag, Space, Menu, Button } from 'antd';
import { BgColorsOutlined, QuestionCircleOutlined } from '@ant-design/icons';
import React from 'react';
import { useModel } from 'umi';
import Avatar from './AvatarDropdown';
import HeaderDropdown from '../HeaderDropdown';
//import HeaderSearch from '../HeaderSearch';
import FullScreen from '../FullScreen';// 引入全屏组件
import styles from './index.less';
import { setLocale, getLocale, FormattedMessage } from 'umi';

export type SiderTheme = 'light' | 'dark';

const ENVTagColor = {
  dev: 'orange',
  test: 'green',
  pre: '#87d068',
};

const GlobalHeaderRight: React.FC = () => {
  // const { initialState } = useModel('@@initialState');
  const { initialState, loading, refresh, setInitialState } = useModel('@@initialState');

  if (!initialState || !initialState.settings) {
    return null;
  }

  const { navTheme, layout } = initialState.settings;
  let className = styles.right;

  if ((navTheme === 'dark' && layout === 'top') || layout === 'mix') {
    className = `${styles.right}  ${styles.dark}`;
  }

  const changLang = () => {
    const locale = getLocale();
    if (!locale || locale === 'zh-CN') {
      setLocale('en-US');
    } else {
      setLocale('zh-CN');
    }
  };

  return (
    <Space className={className}>
      {/* <HeaderSearch
        className={`${styles.action} ${styles.search}`}
        placeholder="站内搜索"
        defaultValue="lepus"
        options={[
          { label: <a href="https://www.lepus.cc">lepus</a>, value: 'lepus' },
          {
            label: <a href="/alarm/event">告警事件</a>,
            value: '告警事件',
          },
          {
            label: <a href="/user/manager/">用户管理</a>,
            value: '用户管理',
          },
        ]}
      // onSearch={value => {
      //   console.log('input', value);
      // }}
      /> */}
      <div>
        <Button
          size="small"
          style={{
            margin: '0 8px',
          }}
          onClick={() => {
            changLang();
          }}
        >
          <FormattedMessage id="navBar.lang" />
        </Button>
      </div>
      <FullScreen />
      <HeaderDropdown
        overlay={
          <Menu>
            <Menu.Item
              onClick={() => {
                setInitialState((preInitialState) => ({
                  ...preInitialState,
                  settings: {
                    ...initialState.settings,
                    navTheme: 'light'
                  }
                }))
              }}
            >
              开灯
            </Menu.Item>
            <Menu.Item
              onClick={() => {
                setInitialState((preInitialState) => ({
                  ...preInitialState,
                  settings: {
                    ...initialState.settings,
                    navTheme: 'realDark'
                  }
                }))
              }}
            >
              关灯
            </Menu.Item>
          </Menu>
        }
      >
        <span className={styles.action}>
          <BgColorsOutlined />
        </span>
      </HeaderDropdown>
      <HeaderDropdown
        overlay={
          <Menu>
            <Menu.Item
              onClick={() => {
                window.open('https://www.lepus.cc');
              }}
            >
              Lepus官网
            </Menu.Item>
            <Menu.Item
              onClick={() => {
                window.open('https://discuss.lepus.cc');
              }}
            >
              Lepus社区
            </Menu.Item>
          </Menu>
        }
      >
        <span className={styles.action}>
          <QuestionCircleOutlined />
        </span>
      </HeaderDropdown>

      <Avatar />
      {REACT_APP_ENV && (
        <span>
          <Tag color={ENVTagColor[REACT_APP_ENV]}>{REACT_APP_ENV}</Tag>
        </span>
      )}
      {/*<SelectLang className={styles.action} />*/}
    </Space>
  );
};
export default GlobalHeaderRight;
