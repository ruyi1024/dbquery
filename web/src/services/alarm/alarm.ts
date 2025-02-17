import {request} from "umi";

export async function batchUpdateAlarmStatus(body: any) {
  return request<any>(`/api/v1/alarm/batchUpdateStatus`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
  });
}
