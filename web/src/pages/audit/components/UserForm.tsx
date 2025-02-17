import React, { useState, useEffect } from 'react';
import {Form, Input,Checkbox } from 'antd';
import {Modal} from "antd";
import type {UserListItem } from "../data";

type UsersFormProps = {
  updateModalVisible: boolean;
  onSubmit: (values: UserListItem) => void;
  onCancel: () => void;
  values: UserListItem;
}


const layout = {
  labelCol: { span: 5 },
  wrapperCol: { span: 16 },
};

const UsersForm: React.FC<UsersFormProps> = ({
  updateModalVisible,
  onSubmit,
  onCancel,
  values
  }) => {
  // const intl = useIntl();
  const [form] = Form.useForm();
  const [formValues, setFormValues] = useState<UserListItem>()

  useEffect(() => {
    if(values !== null){
      console.log("values:", values)
      setFormValues(values);
      form.setFieldsValue({...values})
    } else {
      form.resetFields();
    }
  }, [values]);

  return <Modal
    destroyOnClose
    width={500}
    title={values.modify ? `修改用户${values.username}` : '新增用户'}
    visible={updateModalVisible}
    onCancel={onCancel}
    onOk={() => {
      form
        .validateFields()
        .then(vals => {
          form.resetFields();
          const data = {...vals, modify: values.modify};
          if (values.modify) {
            // @ts-ignore
            data.id = values.id || 0
          }
          onSubmit({...data});
        })
        .catch(info => {
          console.log('Validate Failed:', info);
        });
    }}
  >
      <Form
        {...layout}
        form={form}
        initialValues={formValues ? {...formValues} : {}}
        preserve={false}
      >
        <Form.Item name="username" label="用户" rules={[{ required: true }]}>
          <Input style={{width: 180}} />
        </Form.Item>
        <Form.Item name="chineseName" label="姓名" rules={[{ required: true }]}>
          <Input style={{width: 180}} />
        </Form.Item>
        <Form.Item name="password" label="密码" rules={[{ required: !values.modify }]}>
          <Input.Password />
        </Form.Item>
        <Form.Item name="admin" valuePropName="checked" label="管理员">
          <Checkbox />
        </Form.Item>
        {/*<Form.Item name="admin" label="是否管理员" rules={[{ required: !values.modify }]}>*/}
        {/*<Select  style={{ width: 180 }} >*/}
        {/*  <Option value={0}>否</Option><Option value={1}>是</Option>*/}
        {/*</Select>*/}
        {/*</Form.Item>*/}
        <Form.Item name="remark" label="备注">
          <Input.TextArea />
        </Form.Item>
      </Form>

  </Modal>
}

export default UsersForm;
