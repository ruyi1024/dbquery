export interface TableListItem {
  id: number;
  idc_id: number;
  env_id: number;
  ip_address: string;
  hostname: string;
  description: string;
  online: number;
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
  idc_id?: number;
  env_id?: number;
  ip_address?: string;
  hostname?: string;
  description?: string;
  online?: number;
  pageSize?: number;
  currentPage?: number;
  filter?: { [key: string]: any[] };
  sorter?: { [key: string]: any };
}
