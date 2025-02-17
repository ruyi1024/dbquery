export interface TableListItem {
  id: number;
  db_type: string;
  db_group: string;
  ip: string;
  port: string;
  user: string;
  pass: string;
  tag: string;
  monitor: number;
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
  db_type: string;
  db_group: string;
  ip: string;
  port: string;
  user: string;
  pass: string;
  tag: string;
  monitor: number;
  pageSize?: number;
  currentPage?: number;
  filter?: { [key: string]: any[] };
  sorter?: { [key: string]: any };
}
