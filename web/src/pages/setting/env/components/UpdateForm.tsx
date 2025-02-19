import React from 'react';
import { Modal } from 'antd';
import { FormattedMessage } from 'umi';

interface UpdateFormProps {
  updateModalVisible: boolean;
  onCancel: () => void;
}

const UpdateForm: React.FC<UpdateFormProps> = (props) => {
  const { updateModalVisible, onCancel } = props;

  return (
    <Modal
      destroyOnClose
      title={<FormattedMessage id="pages.searchTable.operate.update" />}
      visible={updateModalVisible}
      onCancel={() => onCancel()}
      footer={null}
    >
      {props.children}
    </Modal>
  );
};

export default UpdateForm;
