import {request} from "@@/plugin-request/request";
import { TableListParams } from './data.d';

export async function queryWeb(params?: TableListParams) {
  return request('/api/v1/meta/web/list', {
    params,
  });
}

export async function removeWeb(params: { key: number[] }) {
  return request('/api/v1/meta/web/list', {
    method: 'DELETE',
    data: {
      ...params,
    },
  });
}

export async function addWeb(params: TableListParams) {
  return request('/api/v1/meta/web/list', {
    method: 'POST',
    data: {
      ...params,

    },
  });
}

export async function updateWeb(params: TableListParams) {
  return request('/api/v1/meta/web/list', {
    method: 'PUT',
    data: {
      ...params,
    },
  });
}
