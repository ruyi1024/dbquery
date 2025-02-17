declare namespace API {
  export type CurrentUser = {
    success: boolean;
    errorCode?: number;
    errorMsg?: string;
    data: {
      admin: boolean;
      id?: number;
      avatar?: string;
      username?: string;
      chineseName?: string;
    },
    access?: 'user' | 'guest' | 'admin';
  };

  export type LoginStateType = {
    status?: 'ok' | 'error';
    type?: string;
  };

  export type LoginStatus = {
    success?: boolean;
    msg?: string;
  };

  export type NoticeIconData = {
    id: string;
    key: string;
    avatar: string;
    title: string;
    datetime: string;
    type: string;
    read?: boolean;
    description: string;
    clickClose?: boolean;
    extra: any;
    status: string;
  };
}
