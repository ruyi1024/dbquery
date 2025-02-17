import {request} from "@@/plugin-request/request";
import { TableListParams } from './data.d';

export async function queryHost(params?: TableListParams) {
  return request('/api/v1/meta/host/list', {
    params,
  });
}

export async function removeHost(params: { key: number[] }) {
  return request('/api/v1/meta/host/list', {
    method: 'DELETE',
    data: {
      ...params,
    },
  });
}

export async function addHost(params: TableListParams) {
  return request('/api/v1/meta/host/list', {
    method: 'POST',
    data: {
      ...params,

    },
  });
}

export async function updateHost(params: TableListParams) {
  return request('/api/v1/meta/host/list', {
    method: 'PUT',
    data: {
      ...params,
    },
  });
}
