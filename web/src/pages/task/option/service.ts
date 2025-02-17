import { request } from "@@/plugin-request/request";
import { TableListParams } from './data.d';

export async function query(params?: TableListParams) {
  return request('/api/v1/task/option', {
    params,
  });
}

export async function remove(params: { key: number[] }) {
  return request('/api/v1/task/option', {
    method: 'DELETE',
    data: {
      ...params,
    },
  });
}

export async function add(params: TableListParams) {
  return request('/api/v1/task/option', {
    method: 'POST',
    data: {
      ...params,
    },
  });
}

export async function update(params: TableListParams) {
  return request('/api/v1/task/option', {
    method: 'PUT',
    data: {
      ...params,
    },
  });
}
