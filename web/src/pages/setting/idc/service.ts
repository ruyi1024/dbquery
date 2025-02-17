import { request } from "@@/plugin-request/request";
import { TableListParams } from './data.d';

export async function queryIdc(params?: TableListParams) {
  return request('/api/v1/datasource_idc/list', {
    params,
  });
}

export async function removeIdc(params: { key: number[] }) {
  return request('/api/v1/datasource_idc/list', {
    method: 'DELETE',
    data: {
      ...params,
    },
  });
}

export async function addIdc(params: TableListParams) {
  return request('/api/v1/datasource_idc/list', {
    method: 'POST',
    data: {
      ...params,
    },
  });
}

export async function updateIdc(params: TableListParams) {
  return request('/api/v1/datasource_idc/list', {
    method: 'PUT',
    data: {
      ...params,
    },
  });
}
