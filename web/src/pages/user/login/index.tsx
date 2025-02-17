import {
  LockTwoTone,
  UserOutlined,
} from '@ant-design/icons';
import { Alert, message } from 'antd';
import React, { useState } from 'react';
import ProForm, { ProFormCheckbox, ProFormText } from '@ant-design/pro-form';
import { useIntl, Link, history, FormattedMessage, SelectLang, useModel } from 'umi';
import Footer from '@/components/Footer';
import type { LoginParamsType } from '@/services/login';
import { fakeAccountLogin, } from '@/services/login';
import CryptoJS from 'crypto-js'


import styles from './index.less';

const LoginMessage: React.FC<{
  content: string;
}> = ({ content }) => (
  <Alert
    style={{
      marginBottom: 24,
    }}
    message={content}
    type="error"
    showIcon
  />
);

/**
 * 此方法会跳转到 redirect 参数所在的位置
 */
const goto = () => {
  if (!history) return;
  setTimeout(() => {
    const { query } = history.location;
    const { redirect } = query as { redirect: string };
    history.push(redirect || '/');
  }, 10);
};

const Login: React.FC = () => {
  const [submitting, setSubmitting] = useState(false);
  //const [userLoginStatus, setUserLoginStatus] = useState<API.LoginStatus>({success:true,msg:""});
  const [userLoginStatus, setUserLoginStatus] = useState<boolean>(true);
  const [type, setType] = useState<string>('account');
  const { initialState, setInitialState } = useModel('@@initialState');

  const intl = useIntl();

  const fetchUserInfo = async () => {
    const userInfo = await initialState?.fetchUserInfo?.();
    if (userInfo) {
      setInitialState({
        ...initialState,
        currentUser: userInfo,
      });
    }
  };

  const PaddingLeft = (key: string, length: number) => {
    let pkey = key.toString();
    const l = pkey.length;
    if (l < length) {
      pkey = new Array(length - l + 1).join('0') + pkey;
    } else if (l > length) {
      pkey = pkey.slice(length);
    }
    return pkey;
  }

  const handleSubmit = async (values: LoginParamsType) => {
    console.info("values:", values)
    setSubmitting(true);
    try {
      // 登录
      let key = "1234567890abcdef";
      // 证key的长度为16byte,进行'0'补位
      key = PaddingLeft(key, 16);
      key = CryptoJS.enc.Utf8.parse(key);
      const encrypted = CryptoJS.AES.encrypt(JSON.stringify({ ...values, type }), key, {
        iv: key,
        mode: CryptoJS.mode.CBC,
        padding: CryptoJS.pad.Pkcs7
      });
      console.info("login post...")
      const loginState = await fakeAccountLogin(encrypted.ciphertext.toString(CryptoJS.enc.Hex));
      console.info("loginState", loginState)
      if (loginState.successLogin) {
        message.success('登录成功！');
        await fetchUserInfo();
        goto();
        return;
      }
      // 如果失败去设置用户错误信息
      message.warning(loginState.msg);
      setUserLoginStatus(loginState.success);
    } catch (error) {
      message.error('登录失败，请重试');
    }
    setSubmitting(false);
  };

  const loginType = "account";

  return (
    <div className={styles.container}>
      <div className={styles.lang}>{SelectLang && <SelectLang />}</div>
      <div className={styles.content}>
        <div className={styles.top}>
          <div className={styles.header}>
            <Link to="/">
              <img alt="logo" className={styles.logo} src="/logo.png" />
              <span className={styles.title}>LEPUS</span>
            </Link>
          </div>
          <div className={styles.desc}>数据库SQL查询、安全、监控、管理平台</div>
        </div>

        <div className={styles.main}>
          <ProForm
            initialValues={{
              autoLogin: true,
            }}
            submitter={{
              searchConfig: {
                submitText: intl.formatMessage({
                  id: 'pages.login.submit',
                  defaultMessage: '登录',
                }),
              },
              render: (_, dom) => dom.pop(),
              submitButtonProps: {
                loading: submitting,
                size: 'large',
                style: {
                  width: '100%',
                },
              },
            }}
            onFinish={async (values) => {
              handleSubmit(values as LoginParamsType);
            }}
          >

            {!userLoginStatus && loginType === 'account' && (
              <LoginMessage
                content="账号或密码错误"
              />
            )}
            {type === 'account' && (
              <>
                <ProFormText
                  name="username"
                  fieldProps={{
                    size: 'large',
                    prefix: <UserOutlined className={styles.prefixIcon} />,
                  }}
                  placeholder={intl.formatMessage({
                    id: 'pages.login.username.placeholder',
                    defaultMessage: '访客账号: guest',
                  })}
                  rules={[
                    {
                      required: true,
                      message: (
                        <FormattedMessage
                          id="pages.login.username.required"
                          defaultMessage="请输入账号!"
                        />
                      ),
                    },
                  ]}
                />
                <ProFormText.Password
                  name="password"
                  fieldProps={{
                    size: 'large',
                    prefix: <LockTwoTone className={styles.prefixIcon} />,
                  }}
                  placeholder={intl.formatMessage({
                    id: 'pages.login.password.placeholder',
                    defaultMessage: '访客密码: guest',
                  })}
                  rules={[
                    {
                      required: true,
                      message: (
                        <FormattedMessage
                          id="pages.login.password.required"
                          defaultMessage="请输入密码！"
                        />
                      ),
                    },
                  ]}
                />
              </>
            )}

            {/* <div
              style={{
                marginBottom: 24,
              }}
            >
              <ProFormCheckbox noStyle name="autoLogin">
                <FormattedMessage id="pages.login.rememberMe" defaultMessage="自动登录" />
              </ProFormCheckbox>
              <a
                style={{
                  float: 'right',
                }}
              >
                <FormattedMessage id="pages.login.forgotPassword" defaultMessage="忘记密码" />
              </a>
            </div> */}
          </ProForm>
        </div>
      </div>
      <Footer />
    </div>
  );
};

export default Login;
