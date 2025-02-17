import { request } from "@@/plugin-request/request";
import { TableListParams } from './data.d';

export async function query(params?: TableListParams) {
  return request('/api/v1/task/heartbeat', {
    params,
  });
}
