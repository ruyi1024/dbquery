import {request} from "@@/plugin-request/request";
import { stringify } from 'qs';

export async function queryUser(params?: string) {
  return request(`/api/v1/users/manager/lists?${stringify(params)}`);
}

export async function updateUser(params: { modify?: boolean; admin?:boolean; createdAt?: Date; password?: string; chineseName?: string; id?: number; username?: string; updatedAt?: Date }) {
  return request(`/api/v1/users/manager/lists`, {
    method: params.modify ? 'PUT' : 'POST',
    data: {
      ...params,
    },
  });
}

export async function removeUser(params: { username?: string }) {
  return request('/api/v1/users/manager/lists', {
    method: 'DELETE',
    data: {
      ...params,
    },
  });
}
