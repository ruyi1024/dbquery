import {request} from "@@/plugin-request/request";
import { TableListParams } from './data.d';

export async function queryModule(params?: TableListParams) {
  return request('/api/v1/meta/module/list', {
    params,
  });
}

export async function removeModule(params: { key: number[] }) {
  return request('/api/v1/meta/module/list', {
    method: 'DELETE',
    data: {
      ...params,
    },
  });
}

export async function addModule(params: TableListParams) {
  return request('/api/v1/meta/module/list', {
    method: 'POST',
    data: {
      ...params,
    },
  });
}

export async function updateModule(params: TableListParams) {
  return request('/api/v1/meta/module/list', {
    method: 'PUT',
    data: {
      ...params,
    },
  });
}
