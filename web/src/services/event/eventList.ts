import { request } from 'umi';
import { stringify } from 'qs';

export async function getEventList(params?: {
  sorterOrder?: string;
  typeKeyword: string;
  offset: number;
  pageSize?: number;
  groupKeyword: string;
  sorterField?: stirng;
  eventKeyword: string;
  eventKeyKeyword: string;
  limit?: number;
  reset?: boolean;
  startTime: moment.Moment | string;
  endTime: moment.Moment | string;
  id?: number;
  eventEntityKeyword: string;
}) {
  return request<API.EventListRes>(`/api/v1/event?${stringify({ ...params })}`, {
    method: 'GET',
  });
}

export async function getEventFilterItems() {
  return request<API.EventListRes>('/api/v1/event/filterItems', {
    method: 'GET',
  });
}

export async function getEventCharts(params?: API.DidParams) {
  return request<API.EventListRes>(`/api/v1/event/charts?${stringify({ ...params })}`, {
    method: 'GET',
  });
}

export async function getEventChartsFull(params?: {
  sorterOrder?: string;
  typeKeyword: string;
  pageSize?: number;
  groupKeyword: string;
  sorterField?: stirng;
  eventKeyword: string;
  eventKeyKeyword: string;
  limit?: number;
  reset?: boolean;
  startTime: moment.Moment | string;
  endTime: moment.Moment | string;
  id?: number;
  eventEntityKeyword: string;
}) {
  return request<API.EventListRes>(`/api/v1/event/chartsFull?${stringify({ ...params })}`, {
    method: 'GET',
  });
}
export async function getEventAllDescription() {
  return request<API.EventListRes>(`/api/v1/event/all/list`, {
    method: 'GET',
  });
}
