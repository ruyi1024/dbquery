import {request} from "@@/plugin-request/request";
import { TableListParams } from './data.d';

export async function queryTask(params?: TableListParams) {
  return request('/api/v1/task/list', {
    params,
  });
}

export async function removeInstance(params: { key: number[] }) {
  return request('/api/v1/meta/instance', {
    method: 'DELETE',
    data: {
      ...params,
    },
  });
}

export async function addInstance(params: TableListParams) {
  return request('/api/v1/meta/instance', {
    method: 'POST',
    data: {
      ...params,
    },
  });
}

export async function updateInstance(params: TableListParams) {
  return request('/api/v1/meta/instance', {
    method: 'PUT',
    data: {
      ...params,
    },
  });
}
