
declare namespace API {

  interface DidParams {
    id?: number;
    limit?: number;
    pageSize?: number;
    sorterField?: stirng;
    sorterOrder?: string;
    eventKeyword?: string;
    typeKeyword?: string;
    groupKeyword?: string;
    eventEntityKeyword?: string;
    eventKeyKeyword?: string;
    startTime?: string;
    endTime?: string;
    reset?: boolean;
  }

  interface ResData {
    success: boolean;
    errorCode?: number,
    errorMsg?: string;
    data?: any;
    total?: number;
  }

  type EventListRes = ResData & {
    data?: EventItem[];
  }

  export interface EventItem {
    id: number;
    title: string;
    event_type: string;
    event_group: string;
    event_key: string;
    event_entity: string;
    alarm_rule: string;
    alarm_value: string;
    alarm_level: string;
    alarm_sleep: number;
    alarm_times: number;
    channel_id: number;
    gmt_created: date;
    gmt_updated: date;
  }

  export interface TableListPagination {
    total: number;
    pageSize: number;
    current: number;
  }

  export interface TableListData {
    list: EventItem[];
    pagination: Partial<TableListPagination>;
  }

  export interface TableListParams {
    id?: number;
    event_type?: string;
    event_group?: string;
    pageSize?: number;
    currentPage?: number;
    filter?: { [key: string]: any[] };
    sorter?: { [key: string]: any };
  }
}
