import { request } from 'umi';

export type LoginParamsType = {
  username: string;
  password: string;
};

export async function fakeAccountLogin(params: any) {
  return request<any>('/api/v1/login/account', {
    method: 'POST',
    data: params,
  });
}
export async function outLogin() {
  return request('/api/v1/login/outLogin');
}
