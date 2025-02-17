import { request } from "@@/plugin-request/request";
import { TableListParams } from './data.d';

export async function queryDatabase(params?: TableListParams) {
  return request('/api/v1/meta/database/list', {
    params,
  });
}

export async function removeDatabase(params: { key: number[] }) {
  return request('/api/v1/meta/database/list', {
    method: 'DELETE',
    data: {
      ...params,
    },
  });
}

export async function addDatabase(params: TableListParams) {
  return request('/api/v1/meta/database/list', {
    method: 'POST',
    data: {
      ...params,
    },
  });
}

export async function updateDatabase(params: TableListParams) {
  return request('/api/v1/meta/database/list', {
    method: 'PUT',
    data: {
      ...params,
    },
  });
}
