import { request } from "@@/plugin-request/request";
import { TableListParams } from './data.d';

export async function queryTable(params?: TableListParams) {
  return request('/api/v1/meta/table/list', {
    params,
  });
}

