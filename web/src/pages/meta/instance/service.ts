import { request } from "@@/plugin-request/request";
import { TableListParams } from './data.d';

export async function queryInstance(params?: TableListParams) {
  return request('/api/v1/meta/instance/list', {
    params,
  });
}
