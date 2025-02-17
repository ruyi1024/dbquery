import { request } from "@@/plugin-request/request";
import { TableListParams } from './data.d';

export async function queryColumn(params?: TableListParams) {
  return request('/api/v1/meta/column/list', {
    params,
  });
}

