import {request} from "@@/plugin-request/request";
import { TableListParams } from './data.d';

export async function queryEnv(params?: TableListParams) {
  return request('/api/v1/meta/env/list', {
    params,
  });
}

export async function removeEnv(params: { key: number[] }) {
  return request('/api/v1/meta/env/list', {
    method: 'DELETE',
    data: {
      ...params,
    },
  });
}

export async function addEnv(params: TableListParams) {
  return request('/api/v1/meta/env/list', {
    method: 'POST',
    data: {
      ...params,
    },
  });
}

export async function updateEnv(params: TableListParams) {
  return request('/api/v1/meta/env/list', {
    method: 'PUT',
    data: {
      ...params,
    },
  });
}
