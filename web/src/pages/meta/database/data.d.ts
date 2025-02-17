export interface TableListItem {
  id: number;
  datasource_type: string;
  host: string;
  port: string;
  database_name: string;
  characters: string;
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
  datasource_type?: string;
  host?: string;
  port?: string;
  database_name?: string;
  characters?: string;
  pageSize?: number;
  currentPage?: number;
  filter?: { [key: string]: any[] };
  sorter?: { [key: string]: any };
}
