import {request} from "@@/plugin-request/request";
import { stringify } from 'qs';

export async function queryLog(params?: string) {
  return request(`/api/v1/audit/query_log?${stringify(params)}`);
}


