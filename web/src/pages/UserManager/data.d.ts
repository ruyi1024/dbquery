export type UserListItem = {
  id: number;
  username?: string;
  chineseName: string;
  password: string;
  updatedAt: Date;
  createdAt: Date;
  admin?: boolean;
  modify?: boolean;
  remark?: string;
};

export type UserListData = {
  list: UserListItem[];
};
