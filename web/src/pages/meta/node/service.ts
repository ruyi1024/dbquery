import {request} from "@@/plugin-request/request";
import { TableListParams } from './data.d';

export async function queryNode(params?: TableListParams) {
  return request('/api/v1/meta/node/list', {
    params,
  });
}

export async function removeNode(params: { key: number[] }) {
  return request('/api/v1/meta/node/list', {
    method: 'DELETE',
    data: {
      ...params,
    },
  });
}

export async function addNode(params: TableListParams) {
  return request('/api/v1/meta/node/list', {
    method: 'POST',
    data: {
      ...params,
    },
  });
}

export async function updateNode(params: TableListParams) {
  return request('/api/v1/meta/node/list', {
    method: 'PUT',
    data: {
      ...params,
    },
  });
}
