export interface TableListItem {
  id: number;
  idc: string;
  env: string;
  type: string;
  name: string;
  host: string;
  port: string;
  user: string;
  pass: string;
  dbid: string;
  enable: number;
  dbmeta_enable: number;
  sensitive_enable: number;
  execute_enable: number;
  monitor_enable: number;
  alarm_enable: number;
  gmt_created: date;
  gmt_updated: date;
}

export interface TableListPagination {
  total: number;
  pageSize: number;
  current: number;
}

export interface TableListData {
  list: TableListItem[];
  pagination: Partial<TableListPagination>;
}

export interface TableListParams {
  id?: number;
  type?: string;
  name?: string;
  host?: string;
  port?: string;
  user?: string;
  pass?: string;
  enable?: number;
  pageSize?: number;
  currentPage?: number;
  filter?: { [key: string]: any[] };
  sorter?: { [key: string]: any };
}
