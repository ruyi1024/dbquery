export interface TableListItem {
  id: number;
  cluster_id: number;
  ip: string;
  domain: string;
  port: string;
  user: string;
  pass: string;
  role: number;
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
  cluster_id?: number;
  ip?: string;
  domain?: string;
  port?: string;
  user?: string;
  pass?: string;
  role?: number;
  monitor?: number;
  pageSize?: number;
  currentPage?: number;
  filter?: { [key: string]: any[] };
  sorter?: { [key: string]: any };
}
