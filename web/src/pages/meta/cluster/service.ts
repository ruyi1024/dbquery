import {request} from "@@/plugin-request/request";
import { TableListParams } from './data.d';

export async function queryCluster(params?: TableListParams) {
  return request('/api/v1/meta/cluster/list', {
    params,
  });
}

export async function removeCluster(params: { key: number[] }) {
  return request('/api/v1/meta/cluster/list', {
    method: 'DELETE',
    data: {
      ...params,
    },
  });
}

export async function addCluster(params: TableListParams) {
  return request('/api/v1/meta/cluster/list', {
    method: 'POST',
    data: {
      ...params,
    },
  });
}

export async function updateCluster(params: TableListParams) {
  return request('/api/v1/meta/cluster/list', {
    method: 'PUT',
    data: {
      ...params,
    },
  });
}
