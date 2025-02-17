export interface TableListItem {
  id: number;
  cluster_id: number;
  env_id: number;
  name: string;
  url: string;
  method: string;
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
  env_id?: number;
  name?: string;
  url?: string;
  method?: string;
  monitor?: number;
  pageSize?: number;
  currentPage?: number;
  filter?: { [key: string]: any[] };
  sorter?: { [key: string]: any };
}
