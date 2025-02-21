import { BorderLeftOutlined, HeartOutlined, RightSquareOutlined, TableOutlined, UnorderedListOutlined } from '@ant-design/icons';
import { PageContainer, ProColumns, ProTable } from '@ant-design/pro-components';
import { Card, Tabs, Form, Select, Button, Input, message, Table, Alert, Space, List, Row, Col, AutoComplete, Drawer } from 'antd';

import React, { useState, useEffect, useRef } from 'react';
import moment from 'moment';
import styles from './index.less';

// 引入ACE编辑器
import AceEditor from 'react-ace'
// 引入对应的mode
import 'ace-builds/src-noconflict/mode-mysql'
// 引入对应的theme
import 'ace-builds/src-noconflict/theme-github'
// 如果要有代码提示，下面这句话必须引入!!!
import 'ace-builds/src-noconflict/ext-language_tools'
// js中实现SQL格式化
import { format } from 'sql-formatter'
import 'ace-builds/src-noconflict/mode-sql'
import { Ace } from 'ace-builds';

//导出excel
import * as Exceljs from 'exceljs';
import { Workbook } from 'exceljs';
import { saveAs } from 'file-saver';
import { FormattedMessage } from 'umi';

const favoriteColumns: ProColumns[] = [
  {
    title: <FormattedMessage id="pages.execute.favoriteTime" />,
    dataIndex: 'gmt_created',
  },
  {
    title: <FormattedMessage id="pages.execute.favoriteContent" />,
    dataIndex: 'content',
    copyable: true,
    tip: <FormattedMessage id="pages.execute.favoriteTip" />
  }
]


const Index: React.FC = () => {

  const [form] = Form.useForm();

  const [formValues, setFormValues] = useState({
    datasource: "",
    database: "",
    table: "",
    sql: "",
  });

  const [typeList, setTypeList] = useState<any[]>([{ id: 0, cluster_name: '' }]);
  const [datasourceList, setDatasourceList] = useState([]);
  const [databaseList, setDatabaseList] = useState([]);
  const [tableList, setTableList] = useState([]);
  const [favoriteList, setFavoriteList] = useState([]);
  const [openFavorite, setOpenFavorite] = useState(false);

  const [type, setType] = useState<string>('');
  const [datasource, setDatasource] = useState<string>('');
  const [database, setDatabase] = useState<string>('');
  const [table, setTable] = useState<string>('');
  const [sqlContent, setSqlContent] = useState<string>('');

  const [loading, setLoading] = useState<boolean>(false);
  const [tableDataTotal, setTableDataTotal] = useState<number>(0)
  const [tableDataList, setTableDataList] = useState<any>();
  const [tableDataColumn, setTableDataColumn] = useState<any>();
  const [tableDataSuccess, setTableDataSuccess] = useState<boolean>(false);
  const [tableDataMsg, setTableDataMsg] = useState<any>("");
  const [queryTimes, setQueryTimes] = useState<number>(0);

  const [currentUserinfo, setCurrentUserinfo] = useState({ "chineseName": "", "username": "" });
  const [currentDate, setCurrentDate] = useState<string>("");

  const editorRef = React.createRef()

  useEffect(() => {
    const currentDate = moment().format("YYYYMMDD");
    setCurrentDate(currentDate)
    //获取登录用户信息
    fetch('/api/v1/currentUser')
      .then((response) => response.json())
      .then((json) => {
        setCurrentUserinfo(json.data);
      })
      .catch((error) => {
        console.log('Fetch current userinfo failed', error);
      });

    //获取数据源类型  
    fetch('/api/v1/query/datasource_type')
      .then((response) => response.json())
      .then((json) => {
        setTypeList(json.data);
        const valueDict: { [key: number]: string } = {};
        json.data.forEach((record: { id: string | number; name: string; }) => {
          valueDict[record.id] = record.name;
        });
      })
      .catch((error) => {
        console.log('Fetch type list failed', error);
      });

  }, []);


  //获取数据源
  const didQueryDatasource = (val: string) => {
    setDatabaseList([]);
    setTableList([]);
    setDatabase("");
    setTable("");
    setSqlContent("");
    form.setFieldsValue({ "datasource": "", "database": "", "table": "", "sql": "" });
    const formValue = form.getFieldsValue();
    const type = formValue.type;
    setType(val);
    fetch('/api/v1/query/datasource?type=' + type)
      .then((response) => response.json())
      .then((json) => setDatasourceList(json.data))
      .catch((error) => {
        console.log('fetch datasource list failed', error);
      });
  };

  //获取数据库
  const didQueryDatabase = (val: string) => {
    setDatabaseList([]);
    setTableList([]);
    setDatabase("");
    setTable("");
    setSqlContent("");
    form.setFieldsValue({ "database": "", "table": "", "sql": "" });
    setDatasource(val);
    fetch('/api/v1/query/database?datasource=' + val + '&type=' + type)
      .then((response) => response.json())
      .then((json) => setDatabaseList(json.data))
      .catch((error) => {
        console.log('fetch database list failed', error);
      });
  };

  //获取数据表
  const didQueryTable = (val: string) => {
    setDatabase(val);
    setSqlContent("");
    form.setFieldsValue({ "table": "", "sql": "" });
    fetch('/api/v1/query/table?datasource=' + datasource + '&database=' + val + '&type=' + type)
      .then((response) => response.json())
      .then((json) => (json.data == null ? [] : json.data))
      .then((data) => (setTableList(data)))
      .catch((error) => {
        console.log('fetch table list failed', error);
      });
  };


  //点击表名事件
  const onClickTable = (val: string) => {
    didSetTable(val);
  };

  //点击表名后填充SQL内容
  const didSetTable = (val: string) => {
    setTable(val);
    let sql = ''
    if (type == 'MySQL' || type == 'TiDB' || type == 'Doris' || type == "MariaDB" || type == "GreatSQL" || type == "OceanBase" || type == 'ClickHouse' || type == 'PostgreSQL') {
      sql = "select * from " + val + " limit 100"
    }
    if (type == 'Oracle') {
      sql = "select * from " + database + '.' + val + " where rownum<=100"
    }
    if (type == 'SQLServer') {
      sql = "select top 100 * from " + val
    }
    if (type == 'MongoDB') {
      sql = "select.from('" + val + "')" + ".where('_id','!=','').limit(100)"
    }
    setSqlContent(sql);
    form.setFieldsValue({
      sql: sql
    });
  };

  //自动提示
  const complete = (editor: Ace.Editor, tableDataList: any[]) => {
    const completers = tableDataList.map(item => ({
      name: item.table_name,
      value: item.table_name,
      score: 100,
      meta: '',
    }));
    console.log(completers)
    editor.completers.push({
      getCompletions(editor, session, pos, prefix, callback) {
        callback(null, completers);
      },
    });

  }

  //编辑器内容改变
  const onChangeContent = (value: string) => {
    form.setFieldsValue({ "sql": value });
    setSqlContent(value);
  }

  //选择内容改变
  // const onSelectContent = (value: string) => {
  //   console.info("111111");
  //   alert(value);
  //   form.setFieldsValue({ "sql": value });
  //   setSqlContent(value);
  // }

  //格式化SQL
  const beautifySql = () => {
    if (type == "Redis") {
      message.warning("Redis数据源不支持该功能");
      return;
    }
    if (type == "" || database == "" || sqlContent == "") {
      message.warning("数据源/数据库/SQL不完整，无法格式化SQL");
      return;
    }
    setSqlContent(format(sqlContent));
  }

  //收藏SQL
  const favoriteSql = () => {
    if (type == "" || datasource == "" || sqlContent == "") {
      message.warning("数据源/SQL不完整，无法收藏SQL");
      return;
    }
    const headers = new Headers();
    const params = { "datasource_type": type, "datasource": datasource, "database_name": database, "content": sqlContent };
    headers.append('Content-Type', 'application/json');
    fetch('/api/v1/favorite/list', {
      method: 'post',
      headers: headers,
      body: JSON.stringify(params),
    })
      .then((response) => response.json())
      .then((json) => {
        if (json.success == true) {
          message.success("加入收藏夹成功.")
        } else {
          message.success("加入收藏夹失败.")
        }
      })
      .catch((error) => {
        console.log('fetch data failed', error);
      });
  }

  //打开收藏夹
  const showDrawer = () => {
    if (type == "" || datasource == "") {
      message.warning("选择数据源后才能打开收藏夹");
      return;
    }
    fetch('/api/v1/favorite/list?datasource=' + datasource + '&datasource_type=' + type + '&database=' + database)
      .then((response) => response.json())
      .then((json) => setFavoriteList(json.data == null ? [] : json.data))
      .catch((error) => {
        console.log('fetch favorite list failed', error);
      });
    setOpenFavorite(true);

  }
  //关闭收藏夹
  const closeDrawer = () => {
    setFavoriteList([]);
    setOpenFavorite(false);
  }


  //表单提交查询执行请求
  const asyncFetch = (values: {}) => {
    console.info(values);
    setLoading(true);
    const params = { ...values, "query_type": "execute" };
    const headers = new Headers();
    headers.append('Content-Type', 'application/json');
    fetch('/api/v1/query/doQuery', {
      method: 'post',
      headers: headers,
      body: JSON.stringify(params),
    })
      .then((response) => response.json())
      .then((json) => {
        console.info(json.data);
        setLoading(false);
        return (
          setTableDataSuccess(json.success),
          setTableDataMsg(json.msg),
          setTableDataList(json.data),
          setTableDataColumn(json.columns),
          setTableDataTotal(json.total),
          setQueryTimes(json.times)
        );
      })
      .catch((error) => {
        console.log('fetch data failed', error);
      });
  };


  const onFinish = (fieldValue: []) => {
    const values = {
      datasource_type: fieldValue["type"],
      datasource: fieldValue["datasource"],
      database: fieldValue["database"],
      table: fieldValue["table"],
      sql: fieldValue["sql"],
    };
    setFormValues(values);
    asyncFetch(values);
  };

  const onFinishFailed = (errorInfo: any) => {
    console.info(errorInfo);
    message.error('执行查询未完成.');
  };

  //点击按钮提交
  const queryPost = (query_type: any) => {
    if (query_type != 'doExplain' && (table == "" || table == null)) {
      message.error('请先点击左侧表名称选择表.');
      return;
    }
    setLoading(true);
    const params = { "datasource_type": type, "datasource": datasource, "database": database, "table": table, "sql": sqlContent, "query_type": query_type };
    const headers = new Headers();
    headers.append('Content-Type', 'application/json');
    fetch('/api/v1/query/doQuery', {
      method: 'post',
      headers: headers,
      body: JSON.stringify(params),
    })
      .then((response) => response.json())
      .then((json) => {
        setLoading(false);
        return (
          setTableDataSuccess(json.success),
          setTableDataMsg(json.msg),
          setTableDataList(json.data),
          setTableDataColumn(json.columns),
          setTableDataTotal(json.total),
          setQueryTimes(json.times)
        );
      })
      .catch((error) => {
        console.log('fetch data failed', error);
        message.error('执行查询失败');
      });
  }

  //导出excel模块
  const generateHeaders = (columns: any) => {
    return columns.map((col: { title: any; dataIndex: any; width: number; }) => {
      const obj: ITableHeaer = {
        header: col.title,
        key: col.dataIndex,
        width: col.width / 5 || 20,
      }
      return obj;
    }
    )
  }
  const saveWorkBook = (workbook: Workbook, fileName: string) => {
    workbook.xlsx.writeBuffer().then((data: BlobPart) => {
      const blob = new Blob([data], { type: '' });
      saveAs(blob, fileName);
    })
  }
  const exportExcel = () => {
    //创建工作簿
    const workbook = new Exceljs.Workbook();
    //添加sheet
    const worksheet = workbook.addWorksheet("Result");
    //设置sheet默认行高
    worksheet.properties.defaultRowHeight = 20;
    //设置列
    worksheet.columns = generateHeaders(tableDataColumn);
    //添加行
    let rows = worksheet.addRows(tableDataList);
    //设置字体和对齐方式
    rows?.forEach(row => {
      row.font = {
        size: 11,
        name: '宋体',
      }
      row.alignment = { vertical: 'middle', 'horizontal': 'left', wrapText: false };
    })
    //设置首行样式
    let headerRow = worksheet.getRow(1);
    headerRow.eachCell((cell, _colNum) => {
      //设置背景
      cell.fill = {
        type: 'pattern',
        pattern: 'solid',
        fgColor: { argb: '0099CC' }
      }
      //设置字体
      cell.font = {
        bold: true,
        italic: false,
        size: 11,
        name: '宋体',
        color: { argb: 'FFFFFF' }
      }
      //设置对齐
      cell.alignment = { vertical: 'middle', 'horizontal': 'center', wrapText: false };
    })

    //生成文件名
    const date = new Date();
    const year = date.getFullYear().toString();
    const month = (date.getMonth() + 1).toString();
    const day = date.getDate().toString();
    const hour = date.getHours().toString();
    const minute = date.getMinutes().toString();
    const second = date.getSeconds().toString();
    const exportFileName = type + "-" + year + month + day + hour + minute + second + '.xlsx'
    //导出文件
    saveWorkBook(workbook, exportFileName)
    //记录日志
    writeLog("exportExcel");

  }

  //前端调用记录日志方法
  const writeLog = (doType: string) => {
    const params = { "datasource_type": type, "datasource": datasource, "database": database, "sql": sqlContent, "query_type": doType };
    const headers = new Headers();
    headers.append('Content-Type', 'application/json');
    fetch('/api/v1/query/writeLog', {
      method: 'post',
      headers: headers,
      body: JSON.stringify(params),
    })
      .then((response) => response.json())
      .then((json) => {
        if (json.success == true) {
          return true;
        }
        return false;
      })
      .catch((error) => {
        return false;
      });
  }


  return (
    <PageContainer>
      <Row style={{ marginTop: '10px' }}><Col span={24}><Card>
        <Form
          style={{ marginTop: 0 }}
          form={form}
          onFinish={onFinish}
          onFinishFailed={onFinishFailed}
          initialValues={{}}
          name={'sqlForm'}
          layout="inline"
        >

          <Form.Item
            name={'type'}
            label={<FormattedMessage id="pages.execute.datasourceType" />}
            rules={[{ required: true, message: <FormattedMessage id="pages.execute.selectDatasourceType" /> }]}
          >
            {/* <Radio.Group defaultValue="" onChange={(val) => {didQueryDatasource(val); }} >
            {typeList && typeList.map(item => <Radio.Button value={item.name}>{item.name}</Radio.Button>)}
          </Radio.Group> */}
            <Select
              showSearch style={{ width: 240 }}
              placeholder={<FormattedMessage id="pages.execute.selectDatasourceType" />}
              onChange={(val) => { didQueryDatasource(val); }}
            >
              {typeList && typeList.map(item => <Option key={item.name} value={item.name}>{item.name}</Option>)}
            </Select>
          </Form.Item>

          <Form.Item
            name={'datasource'}
            label={<FormattedMessage id="pages.execute.datasource" />}
            rules={[{ required: true, message: <FormattedMessage id="pages.execute.selectDatasource" /> }]}
          >
            <Select
              showSearch style={{ width: 320 }}
              placeholder={<FormattedMessage id="pages.execute.selectDatasource" />}
              value={datasource}
              onChange={(val) => {
                didQueryDatabase(val);
              }}
            >
              {datasourceList && datasourceList.map(item => <Option key={item.host + ":" + item.port} value={item.host + ":" + item.port}>{item.name}[{item.status == 1 ? "可用" : "不可用"}] </Option>)}
            </Select>
          </Form.Item>

          {type !== "Redis" &&
            <Form.Item
              name={'database'}
              label={<FormattedMessage id="pages.execute.database" />}
              rules={[{ required: true, message: <FormattedMessage id="pages.execute.selectDatabase" /> }]}
            >
              <Select showSearch style={{ width: 240 }} placeholder={<FormattedMessage id="pages.execute.selectDatabase" />} value={database}
                onChange={(val) => {
                  didQueryTable(val);
                }}
              >
                {databaseList && databaseList.map(item => <Option key={item.database_name} value={item.database_name}>{item.database_name}</Option>)}
              </Select>
            </Form.Item>
          }
        </Form>
      </Card></Col></Row>

      <Row>
        {type != "Redis" &&
          <Col span={4}>
            <Card size='small' title={<FormattedMessage id="pages.execute.table" />} extra={<a href='javascript:void(0)' onClick={event => didQueryTable(database)}><FormattedMessage id="pages.execute.refresh" /></a>} style={{ width: '100%', height: '750px', overflow: 'auto' }}>
              <List
                size="small"
                dataSource={tableList}
                renderItem={tableList != null && (item => <List.Item><a href='javascript:void(0)' onClick={event => onClickTable(item.table_name)}><TableOutlined /> {item.table_name}</a></List.Item>)}
              />
            </Card>
          </Col>
        }

        <Col span={20}>
          <Card>
            {database && database.length > 0 &&
              <Alert message={"当前查询引擎:" + type + ", 当前数据库:" + database} type="info" showIcon closable />
            }
            {type == "Redis" &&
              <Space direction='vertical'>
                <Alert message="请选择查询数据源，再输入命令，当前支持的命令有：RANDOMKEY、EXISTS、TYPE、TTL、GET、HLEN、HKEYS、HGET、HGETALL、LLEN、LINDEX、LRANGE、SCARD、SMEMBERS、SISMEMBER、ZCARD、ZCOUNT、ZRANGE" type="info" showIcon closable />
              </Space>
            }
            <Form
              style={{ marginTop: 8 }}
              form={form}
              onFinish={onFinish}
              onFinishFailed={onFinishFailed}
              initialValues={{}}
              name={'sqlForm'}
              layout="horizontal"
            >

              <Form.Item
                name={'sql'}
                rules={[{ required: true, message: "请输入SQL语句" }]}
              >
                {/* <TextArea
                autoSize={{minRows: 4, maxRows: 8}}
                defaultValue={sqlContent}
                value={sqlContent}
              /> */}
                <AceEditor
                  ref={editorRef}
                  placeholder="请输入SQL语句"
                  mode="mysql"
                  theme="textmate"
                  name="blah2"
                  fontSize={14}
                  showPrintMargin={true}
                  showGutter={true}
                  highlightActiveLine={true}
                  style={{ width: '100%', height: '200px', border: '1px solid #ccc' }}
                  value={sqlContent}
                  editorProps={{
                    $blockScrolling: false,
                  }}
                  onChange={(value) => onChangeContent(value)} //获取输入框的内容
                  //onPaste={(value)=>onChangeContent(value)}
                  onLoad={editor => complete(editor, tableList)}
                  //onSelection={(selectedText: string, event?: any) => onSelectContent(selectedText)}

                  // 设置编辑器格式化和代码提示 
                  setOptions={{
                    useWorker: false,
                    enableBasicAutocompletion: true,
                    enableLiveAutocompletion: true,
                    // 自动提词此项必须设置为true
                    enableSnippets: true,
                    showLineNumbers: true,
                    tabSize: 1,
                  }}
                />
                <Button htmlType='button' type='dashed' icon={<BorderLeftOutlined />} size='small' onClick={() => beautifySql()}>{<FormattedMessage id="pages.execute.formatSql" />}</Button>
                <Button htmlType='button' type='dashed' icon={<HeartOutlined />} size='small' onClick={() => favoriteSql()}>{<FormattedMessage id="pages.execute.favoriteSql" />}</Button>
                <Button htmlType='button' type='dashed' icon={<UnorderedListOutlined />} size='small' onClick={() => showDrawer()}>{<FormattedMessage id="pages.execute.openFavorite" />}</Button>
              </Form.Item>

              <Form.Item wrapperCol={{ offset: 0, span: 16 }}>
                <Space>
                  <Button type="primary" htmlType="submit" icon={<RightSquareOutlined />}>{<FormattedMessage id="pages.execute.executeSql" />}</Button>

                  {(type == "MySQL" || type == "TiDB" || type == "Doris" || type == "MariaDB" || type == "GreatSQL" || type == "OceanBase") &&
                    <>
                      <Button type="default" htmlType="button" onClick={() => queryPost("doExplain")}>
                        {<FormattedMessage id="pages.execute.showExplain" />}
                      </Button>
                      <Button type="default" htmlType="button" onClick={() => queryPost("showIndex")}>
                        {<FormattedMessage id="pages.execute.showIndex" />}
                      </Button>
                      <Button type="default" htmlType="button" onClick={() => queryPost("showColumn")}>
                        {<FormattedMessage id="pages.execute.showColumn" />}
                      </Button>
                      <Button type="default" htmlType="button" onClick={() => queryPost("showCreate")}>
                        {<FormattedMessage id="pages.execute.showCreate" />}
                      </Button>
                      <Button type="default" htmlType="button" onClick={() => queryPost("showTableSize")}>
                        {<FormattedMessage id="pages.execute.showTableSize" />}
                      </Button>
                    </>
                  }

                  {(type == "Oracle") &&
                    <>
                      <Button type="default" htmlType="button" onClick={() => queryPost("doExplain")}>
                        {<FormattedMessage id="pages.execute.showExplain" />}
                      </Button>
                      <Button type="default" htmlType="button" onClick={() => queryPost("showIndex")}>
                        {<FormattedMessage id="pages.execute.showIndex" />}
                      </Button>
                      <Button type="default" htmlType="button" onClick={() => queryPost("showColumn")}>
                        {<FormattedMessage id="pages.execute.showColumn" />}
                      </Button>
                      <Button type="default" htmlType="button" onClick={() => queryPost("showCreate")}>
                        {<FormattedMessage id="pages.execute.showCreate" />}
                      </Button>
                      <Button type="default" htmlType="button" onClick={() => queryPost("showTableSize")}>
                        {<FormattedMessage id="pages.execute.showTableSize" />}
                      </Button>
                    </>
                  }
                  {(type == "PostgreSQL") &&
                    <>
                      <Button type="default" htmlType="button" onClick={() => queryPost("doExplain")}>
                        {<FormattedMessage id="pages.execute.showExplain" />}
                      </Button>
                      <Button type="default" htmlType="button" onClick={() => queryPost("showIndex")}>
                        {<FormattedMessage id="pages.execute.showIndex" />}
                      </Button>
                      <Button type="default" htmlType="button" onClick={() => queryPost("showColumn")}>
                        {<FormattedMessage id="pages.execute.showColumn" />}
                      </Button>
                      {/* <Button type="default" htmlType="button"  onClick={()=>queryPost("showCreate")}>
                查看建表语句
              </Button> */}
                      <Button type="default" htmlType="button" onClick={() => queryPost("showTableSize")}>
                        {<FormattedMessage id="pages.execute.showTableSize" />}
                      </Button>
                    </>
                  }
                  {(type == "ClickHouse") &&
                    <>
                      <Button type="default" htmlType="button" onClick={() => queryPost("showColumn")}>
                        {<FormattedMessage id="pages.execute.showColumn" />}
                      </Button>
                      <Button type="default" htmlType="button" onClick={() => queryPost("showCreate")}>
                        {<FormattedMessage id="pages.execute.showCreate" />}
                      </Button>
                      <Button type="default" htmlType="button" onClick={() => queryPost("showTableSize")}>
                        {<FormattedMessage id="pages.execute.showTableSize" />}
                      </Button>
                    </>
                  }
                </Space>
              </Form.Item>
            </Form>
          </Card>

          <Card>
            {tableDataSuccess == false && tableDataMsg != "" &&
              <Alert type="error" message={<FormattedMessage id="pages.execute.queryFailed" /> + tableDataMsg} banner />
            }
            {tableDataSuccess == true && tableDataMsg != "" &&
              <Alert type="success" message={"执行成功，耗时：" + queryTimes + "毫秒," + tableDataMsg} banner />
            }
            {tableDataSuccess == true && tableDataTotal >= 0 &&
              <div style={{ whiteSpace: 'pre-wrap', marginTop: '10px' }} >
                <div style={{ width: '100%', float: 'right', marginBottom: '10px' }}>{"查询到" + tableDataTotal + "条数据"} <Button icon={<RightSquareOutlined />} onClick={exportExcel}>查询结果导出Excel</Button></div>
                <Table
                  bordered
                  loading={loading}
                  scroll={{ scrollToFirstRowOnChange: true, x: 100 }}
                  className={styles.tableStyle}
                  dataSource={tableDataList}
                  columns={tableDataColumn}
                  size={'small'}
                />
              </div>
            }
          </Card>
          {/* <Alert type="info" message="支持MySQL/MariaDB/GreatSQL/TiDB/Doris/OceanBase/ClickHouse/Oracle/PostgreSQL/SQLServer/MongoDB/Redis数据查询导出。" banner closable /> */}
        </Col>
      </Row>

      <Drawer
        title={<FormattedMessage id="pages.execute.favorite" />}
        placement='right'
        width={800}
        onClose={closeDrawer}
        visible={openFavorite}
        extra={
          <Space>
            <Button onClick={closeDrawer}>{<FormattedMessage id="pages.execute.close" />}</Button>
          </Space>
        }
      >
        <ProTable rowKey="id" search={false} dataSource={favoriteList} columns={favoriteColumns} size="middle" />

      </Drawer>

    </PageContainer>
  );
};

export default Index;

