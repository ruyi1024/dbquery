(self.webpackChunkant_design_pro=self.webpackChunkant_design_pro||[]).push([[308],{56621:function(F){F.exports={pre:"pre___34IB6"}},26909:function(F,W,s){"use strict";s.r(W);var $t=s(57338),Ie=s(273),Xt=s(66456),Be=s(15885),qt=s(57663),u=s(71577),ea=s(49111),P=s(19650),ta=s(17462),g=s(76772),aa=s(54421),U=s(38272),sa=s(13062),K=s(71230),na=s(89032),M=s(15746),oa=s(58024),y=s(91894),ra=s(43358),O=s(34041),H=s(11849),la=s(34792),f=s(48086),ia=s(9715),E=s(55843),d=s(2824),Re=s(92570),Ae=s(40504),Fe=s(57206),We=s(34707),w=s(26001),Ue=s(67265),Ke=s(36839),c=s(67294),He=s(30381),we=s.n(He),ze=s(56621),Qe=s.n(ze),Ne=s(74981),ke=s(49332),ua=s.n(ke),Ge=s(24203),ca=s.n(Ge),Ye=s(82679),da=s.n(Ye),Ve=s(15500),Je=s(89899),_a=s.n(Je),z=s(15434),ha=s.n(z),Q=s(93162),fa=s.n(Q),r=s(85029),e=s(85893),$e=[{title:(0,e.jsx)(r._H,{id:"pages.execute.favoriteTime"}),dataIndex:"gmt_created"},{title:(0,e.jsx)(r._H,{id:"pages.execute.favoriteContent"}),dataIndex:"content",copyable:!0,tip:(0,e.jsx)(r._H,{id:"pages.execute.favoriteTip"})}],Xe=function(){var qe=E.Z.useForm(),et=(0,d.Z)(qe,1),x=et[0],tt=(0,c.useState)({datasource:"",database:"",table:"",sql:""}),N=(0,d.Z)(tt,2),xa=N[0],at=N[1],st=(0,c.useState)([{id:0,cluster_name:""}]),k=(0,d.Z)(st,2),G=k[0],nt=k[1],ot=(0,c.useState)([]),Y=(0,d.Z)(ot,2),V=Y[0],rt=Y[1],lt=(0,c.useState)([]),J=(0,d.Z)(lt,2),$=J[0],j=J[1],it=(0,c.useState)([]),X=(0,d.Z)(it,2),Z=X[0],b=X[1],ut=(0,c.useState)([]),q=(0,d.Z)(ut,2),ct=q[0],ee=q[1],dt=(0,c.useState)(!1),te=(0,d.Z)(dt,2),_t=te[0],ae=te[1],ht=(0,c.useState)(""),se=(0,d.Z)(ht,2),l=se[0],ft=se[1],Et=(0,c.useState)(""),ne=(0,d.Z)(Et,2),m=ne[0],xt=ne[1],mt=(0,c.useState)(""),oe=(0,d.Z)(mt,2),h=oe[0],L=oe[1],pt=(0,c.useState)(""),re=(0,d.Z)(pt,2),I=re[0],B=re[1],vt=(0,c.useState)(""),le=(0,d.Z)(vt,2),p=le[0],D=le[1],Dt=(0,c.useState)(!1),ie=(0,d.Z)(Dt,2),Ct=ie[0],S=ie[1],gt=(0,c.useState)(0),ue=(0,d.Z)(gt,2),ce=ue[0],de=ue[1],yt=(0,c.useState)(),_e=(0,d.Z)(yt,2),he=_e[0],fe=_e[1],St=(0,c.useState)(),Ee=(0,d.Z)(St,2),xe=Ee[0],me=Ee[1],Tt=(0,c.useState)(!1),pe=(0,d.Z)(Tt,2),R=pe[0],ve=pe[1],Pt=(0,c.useState)(""),De=(0,d.Z)(Pt,2),T=De[0],Ce=De[1],Mt=(0,c.useState)(0),ge=(0,d.Z)(Mt,2),Ot=ge[0],ye=ge[1],jt=(0,c.useState)({chineseName:"",username:""}),Se=(0,d.Z)(jt,2),ma=Se[0],Zt=Se[1],bt=(0,c.useState)(""),Te=(0,d.Z)(bt,2),pa=Te[0],Lt=Te[1],It=c.createRef();(0,c.useEffect)(function(){var a=we()().format("YYYYMMDD");Lt(a),fetch("/api/v1/currentUser").then(function(t){return t.json()}).then(function(t){Zt(t.data)}).catch(function(t){console.log("Fetch current userinfo failed",t)}),fetch("/api/v1/query/datasource_type").then(function(t){return t.json()}).then(function(t){nt(t.data);var n={};t.data.forEach(function(i){n[i.id]=i.name})}).catch(function(t){console.log("Fetch type list failed",t)})},[]);var Bt=function(t){j([]),b([]),L(""),B(""),D(""),x.setFieldsValue({datasource:"",database:"",table:"",sql:""});var n=x.getFieldsValue(),i=n.type;ft(t),fetch("/api/v1/query/datasource?type="+i).then(function(o){return o.json()}).then(function(o){return rt(o.data)}).catch(function(o){console.log("fetch datasource list failed",o)})},Rt=function(t){j([]),b([]),L(""),B(""),D(""),x.setFieldsValue({database:"",table:"",sql:""}),xt(t),fetch("/api/v1/query/database?datasource="+t+"&type="+l).then(function(n){return n.json()}).then(function(n){return j(n.data)}).catch(function(n){console.log("fetch database list failed",n)})},Pe=function(t){L(t),D(""),x.setFieldsValue({table:"",sql:""}),fetch("/api/v1/query/table?datasource="+m+"&database="+t+"&type="+l).then(function(n){return n.json()}).then(function(n){return n.data==null?[]:n.data}).then(function(n){return b(n)}).catch(function(n){console.log("fetch table list failed",n)})},At=function(t){Ft(t)},Ft=function(t){B(t);var n="";(l=="MySQL"||l=="TiDB"||l=="Doris"||l=="MariaDB"||l=="GreatSQL"||l=="OceanBase"||l=="ClickHouse"||l=="PostgreSQL")&&(n="select * from "+t+" limit 100"),l=="Oracle"&&(n="select * from "+h+"."+t+" where rownum<=100"),l=="SQLServer"&&(n="select top 100 * from "+t),l=="MongoDB"&&(n="select.from('"+t+"').where('_id','!=','').limit(100)"),D(n),x.setFieldsValue({sql:n})},Wt=function(t,n){var i=n.map(function(o){return{name:o.table_name,value:o.table_name,score:100,meta:""}});console.log(i),t.completers.push({getCompletions:function(v,Ze,be,Le,A){A(null,i)}})},Ut=function(t){x.setFieldsValue({sql:t}),D(t)},Kt=function(){if(l=="Redis"){f.ZP.warning("Redis\u6570\u636E\u6E90\u4E0D\u652F\u6301\u8BE5\u529F\u80FD");return}if(l==""||h==""||p==""){f.ZP.warning("\u6570\u636E\u6E90/\u6570\u636E\u5E93/SQL\u4E0D\u5B8C\u6574\uFF0C\u65E0\u6CD5\u683C\u5F0F\u5316SQL");return}D((0,Ve.WU)(p))},Ht=function(){if(l==""||m==""||p==""){f.ZP.warning("\u6570\u636E\u6E90/SQL\u4E0D\u5B8C\u6574\uFF0C\u65E0\u6CD5\u6536\u85CFSQL");return}var t=new Headers,n={datasource_type:l,datasource:m,database_name:h,content:p};t.append("Content-Type","application/json"),fetch("/api/v1/favorite/list",{method:"post",headers:t,body:JSON.stringify(n)}).then(function(i){return i.json()}).then(function(i){i.success==!0?f.ZP.success("\u52A0\u5165\u6536\u85CF\u5939\u6210\u529F."):f.ZP.success("\u52A0\u5165\u6536\u85CF\u5939\u5931\u8D25.")}).catch(function(i){console.log("fetch data failed",i)})},wt=function(){if(l==""||m==""){f.ZP.warning("\u9009\u62E9\u6570\u636E\u6E90\u540E\u624D\u80FD\u6253\u5F00\u6536\u85CF\u5939");return}fetch("/api/v1/favorite/list?datasource="+m+"&datasource_type="+l+"&database="+h).then(function(t){return t.json()}).then(function(t){return ee(t.data==null?[]:t.data)}).catch(function(t){console.log("fetch favorite list failed",t)}),ae(!0)},Me=function(){ee([]),ae(!1)},zt=function(t){console.info(t),S(!0);var n=(0,H.Z)((0,H.Z)({},t),{},{query_type:"execute"}),i=new Headers;i.append("Content-Type","application/json"),fetch("/api/v1/query/doQuery",{method:"post",headers:i,body:JSON.stringify(n)}).then(function(o){return o.json()}).then(function(o){return console.info(o.data),S(!1),ve(o.success),Ce(o.msg),fe(o.data),me(o.columns),de(o.total),ye(o.times)}).catch(function(o){console.log("fetch data failed",o)})},Oe=function(t){var n={datasource_type:t.type,datasource:t.datasource,database:t.database,table:t.table,sql:t.sql};at(n),zt(n)},je=function(t){console.info(t),f.ZP.error("\u6267\u884C\u67E5\u8BE2\u672A\u5B8C\u6210.")},_=function(t){if(t!="doExplain"&&(I==""||I==null)){f.ZP.error("\u8BF7\u5148\u70B9\u51FB\u5DE6\u4FA7\u8868\u540D\u79F0\u9009\u62E9\u8868.");return}S(!0);var n={datasource_type:l,datasource:m,database:h,table:I,sql:p,query_type:t},i=new Headers;i.append("Content-Type","application/json"),fetch("/api/v1/query/doQuery",{method:"post",headers:i,body:JSON.stringify(n)}).then(function(o){return o.json()}).then(function(o){return S(!1),ve(o.success),Ce(o.msg),fe(o.data),me(o.columns),de(o.total),ye(o.times)}).catch(function(o){console.log("fetch data failed",o),f.ZP.error("\u6267\u884C\u67E5\u8BE2\u5931\u8D25")})},Qt=function(t){return t.map(function(n){var i={header:n.title,key:n.dataIndex,width:n.width/5||20};return i})},Nt=function(t,n){t.xlsx.writeBuffer().then(function(i){var o=new Blob([i],{type:""});(0,Q.saveAs)(o,n)})},kt=function(){var t=new z.Workbook,n=t.addWorksheet("Result");n.properties.defaultRowHeight=20,n.columns=Qt(xe);var i=n.addRows(he);i==null||i.forEach(function(C){C.font={size:11,name:"\u5B8B\u4F53"},C.alignment={vertical:"middle",horizontal:"left",wrapText:!1}});var o=n.getRow(1);o.eachCell(function(C,va){C.fill={type:"pattern",pattern:"solid",fgColor:{argb:"0099CC"}},C.font={bold:!0,italic:!1,size:11,name:"\u5B8B\u4F53",color:{argb:"FFFFFF"}},C.alignment={vertical:"middle",horizontal:"center",wrapText:!1}});var v=new Date,Ze=v.getFullYear().toString(),be=(v.getMonth()+1).toString(),Le=v.getDate().toString(),A=v.getHours().toString(),Yt=v.getMinutes().toString(),Vt=v.getSeconds().toString(),Jt=l+"-"+Ze+be+Le+A+Yt+Vt+".xlsx";Nt(t,Jt),Gt("exportExcel")},Gt=function(t){var n={datasource_type:l,datasource:m,database:h,sql:p,query_type:t},i=new Headers;i.append("Content-Type","application/json"),fetch("/api/v1/query/writeLog",{method:"post",headers:i,body:JSON.stringify(n)}).then(function(o){return o.json()}).then(function(o){return o.success==!0}).catch(function(o){return!1})};return(0,e.jsxs)(Ue.ZP,{children:[(0,e.jsx)(K.Z,{style:{marginTop:"10px"},children:(0,e.jsx)(M.Z,{span:24,children:(0,e.jsx)(y.Z,{children:(0,e.jsxs)(E.Z,{style:{marginTop:0},form:x,onFinish:Oe,onFinishFailed:je,initialValues:{},name:"sqlForm",layout:"inline",children:[(0,e.jsx)(E.Z.Item,{name:"type",label:(0,e.jsx)(r._H,{id:"pages.execute.datasourceType"}),rules:[{required:!0,message:(0,e.jsx)(r._H,{id:"pages.execute.selectDatasourceType"})}],children:(0,e.jsx)(O.Z,{showSearch:!0,style:{width:240},placeholder:(0,e.jsx)(r._H,{id:"pages.execute.selectDatasourceType"}),onChange:function(t){Bt(t)},children:G&&G.map(function(a){return(0,e.jsx)(Option,{value:a.name,children:a.name},a.name)})})}),(0,e.jsx)(E.Z.Item,{name:"datasource",label:(0,e.jsx)(r._H,{id:"pages.execute.datasource"}),rules:[{required:!0,message:(0,e.jsx)(r._H,{id:"pages.execute.selectDatasource"})}],children:(0,e.jsx)(O.Z,{showSearch:!0,style:{width:320},placeholder:(0,e.jsx)(r._H,{id:"pages.execute.selectDatasource"}),value:m,onChange:function(t){Rt(t)},children:V&&V.map(function(a){return(0,e.jsxs)(Option,{value:a.host+":"+a.port,children:[a.name,"[",a.status==1?"\u53EF\u7528":"\u4E0D\u53EF\u7528","] "]},a.host+":"+a.port)})})}),l!=="Redis"&&(0,e.jsx)(E.Z.Item,{name:"database",label:(0,e.jsx)(r._H,{id:"pages.execute.database"}),rules:[{required:!0,message:(0,e.jsx)(r._H,{id:"pages.execute.selectDatabase"})}],children:(0,e.jsx)(O.Z,{showSearch:!0,style:{width:240},placeholder:(0,e.jsx)(r._H,{id:"pages.execute.selectDatabase"}),value:h,onChange:function(t){Pe(t)},children:$&&$.map(function(a){return(0,e.jsx)(Option,{value:a.database_name,children:a.database_name},a.database_name)})})})]})})})}),(0,e.jsxs)(K.Z,{children:[l!="Redis"&&(0,e.jsx)(M.Z,{span:4,children:(0,e.jsx)(y.Z,{size:"small",title:(0,e.jsx)(r._H,{id:"pages.execute.table"}),extra:(0,e.jsx)("a",{href:"javascript:void(0)",onClick:function(t){return Pe(h)},children:(0,e.jsx)(r._H,{id:"pages.execute.refresh"})}),style:{width:"100%",height:"750px",overflow:"auto"},children:(0,e.jsx)(U.ZP,{size:"small",dataSource:Z,renderItem:Z!=null&&function(a){return(0,e.jsx)(U.ZP.Item,{children:(0,e.jsxs)("a",{href:"javascript:void(0)",onClick:function(n){return At(a.table_name)},children:[(0,e.jsx)(Re.Z,{})," ",a.table_name]})})}})})}),(0,e.jsxs)(M.Z,{span:20,children:[(0,e.jsxs)(y.Z,{children:[h&&h.length>0&&(0,e.jsx)(g.Z,{message:"\u5F53\u524D\u67E5\u8BE2\u5F15\u64CE:"+l+", \u5F53\u524D\u6570\u636E\u5E93:"+h,type:"info",showIcon:!0,closable:!0}),l=="Redis"&&(0,e.jsx)(P.Z,{direction:"vertical",children:(0,e.jsx)(g.Z,{message:"\u8BF7\u9009\u62E9\u67E5\u8BE2\u6570\u636E\u6E90\uFF0C\u518D\u8F93\u5165\u547D\u4EE4\uFF0C\u5F53\u524D\u652F\u6301\u7684\u547D\u4EE4\u6709\uFF1ARANDOMKEY\u3001EXISTS\u3001TYPE\u3001TTL\u3001GET\u3001HLEN\u3001HKEYS\u3001HGET\u3001HGETALL\u3001LLEN\u3001LINDEX\u3001LRANGE\u3001SCARD\u3001SMEMBERS\u3001SISMEMBER\u3001ZCARD\u3001ZCOUNT\u3001ZRANGE",type:"info",showIcon:!0,closable:!0})}),(0,e.jsxs)(E.Z,{style:{marginTop:8},form:x,onFinish:Oe,onFinishFailed:je,initialValues:{},name:"sqlForm",layout:"horizontal",children:[(0,e.jsxs)(E.Z.Item,{name:"sql",rules:[{required:!0,message:"\u8BF7\u8F93\u5165SQL\u8BED\u53E5"}],children:[(0,e.jsx)(Ne.ZP,{ref:It,placeholder:"\u8BF7\u8F93\u5165SQL\u8BED\u53E5",mode:"mysql",theme:"textmate",name:"blah2",fontSize:14,showPrintMargin:!0,showGutter:!0,highlightActiveLine:!0,style:{width:"100%",height:"200px",border:"1px solid #ccc"},value:p,editorProps:{$blockScrolling:!1},onChange:function(t){return Ut(t)},onLoad:function(t){return Wt(t,Z)},setOptions:{useWorker:!1,enableBasicAutocompletion:!0,enableLiveAutocompletion:!0,enableSnippets:!0,showLineNumbers:!0,tabSize:1}}),(0,e.jsx)(u.Z,{htmlType:"button",type:"dashed",icon:(0,e.jsx)(Ae.Z,{}),size:"small",onClick:function(){return Kt()},children:(0,e.jsx)(r._H,{id:"pages.execute.formatSql"})}),(0,e.jsx)(u.Z,{htmlType:"button",type:"dashed",icon:(0,e.jsx)(Fe.Z,{}),size:"small",onClick:function(){return Ht()},children:(0,e.jsx)(r._H,{id:"pages.execute.favoriteSql"})}),(0,e.jsx)(u.Z,{htmlType:"button",type:"dashed",icon:(0,e.jsx)(We.Z,{}),size:"small",onClick:function(){return wt()},children:(0,e.jsx)(r._H,{id:"pages.execute.openFavorite"})})]}),(0,e.jsx)(E.Z.Item,{wrapperCol:{offset:0,span:16},children:(0,e.jsxs)(P.Z,{children:[(0,e.jsx)(u.Z,{type:"primary",htmlType:"submit",icon:(0,e.jsx)(w.Z,{}),children:(0,e.jsx)(r._H,{id:"pages.execute.executeSql"})}),(l=="MySQL"||l=="TiDB"||l=="Doris"||l=="MariaDB"||l=="GreatSQL"||l=="OceanBase")&&(0,e.jsxs)(e.Fragment,{children:[(0,e.jsx)(u.Z,{type:"default",htmlType:"button",onClick:function(){return _("doExplain")},children:(0,e.jsx)(r._H,{id:"pages.execute.showExplain"})}),(0,e.jsx)(u.Z,{type:"default",htmlType:"button",onClick:function(){return _("showIndex")},children:(0,e.jsx)(r._H,{id:"pages.execute.showIndex"})}),(0,e.jsx)(u.Z,{type:"default",htmlType:"button",onClick:function(){return _("showColumn")},children:(0,e.jsx)(r._H,{id:"pages.execute.showColumn"})}),(0,e.jsx)(u.Z,{type:"default",htmlType:"button",onClick:function(){return _("showCreate")},children:(0,e.jsx)(r._H,{id:"pages.execute.showCreate"})}),(0,e.jsx)(u.Z,{type:"default",htmlType:"button",onClick:function(){return _("showTableSize")},children:(0,e.jsx)(r._H,{id:"pages.execute.showTableSize"})})]}),l=="Oracle"&&(0,e.jsxs)(e.Fragment,{children:[(0,e.jsx)(u.Z,{type:"default",htmlType:"button",onClick:function(){return _("doExplain")},children:(0,e.jsx)(r._H,{id:"pages.execute.showExplain"})}),(0,e.jsx)(u.Z,{type:"default",htmlType:"button",onClick:function(){return _("showIndex")},children:(0,e.jsx)(r._H,{id:"pages.execute.showIndex"})}),(0,e.jsx)(u.Z,{type:"default",htmlType:"button",onClick:function(){return _("showColumn")},children:(0,e.jsx)(r._H,{id:"pages.execute.showColumn"})}),(0,e.jsx)(u.Z,{type:"default",htmlType:"button",onClick:function(){return _("showCreate")},children:(0,e.jsx)(r._H,{id:"pages.execute.showCreate"})}),(0,e.jsx)(u.Z,{type:"default",htmlType:"button",onClick:function(){return _("showTableSize")},children:(0,e.jsx)(r._H,{id:"pages.execute.showTableSize"})})]}),l=="PostgreSQL"&&(0,e.jsxs)(e.Fragment,{children:[(0,e.jsx)(u.Z,{type:"default",htmlType:"button",onClick:function(){return _("doExplain")},children:(0,e.jsx)(r._H,{id:"pages.execute.showExplain"})}),(0,e.jsx)(u.Z,{type:"default",htmlType:"button",onClick:function(){return _("showIndex")},children:(0,e.jsx)(r._H,{id:"pages.execute.showIndex"})}),(0,e.jsx)(u.Z,{type:"default",htmlType:"button",onClick:function(){return _("showColumn")},children:(0,e.jsx)(r._H,{id:"pages.execute.showColumn"})}),(0,e.jsx)(u.Z,{type:"default",htmlType:"button",onClick:function(){return _("showTableSize")},children:(0,e.jsx)(r._H,{id:"pages.execute.showTableSize"})})]}),l=="ClickHouse"&&(0,e.jsxs)(e.Fragment,{children:[(0,e.jsx)(u.Z,{type:"default",htmlType:"button",onClick:function(){return _("showColumn")},children:(0,e.jsx)(r._H,{id:"pages.execute.showColumn"})}),(0,e.jsx)(u.Z,{type:"default",htmlType:"button",onClick:function(){return _("showCreate")},children:(0,e.jsx)(r._H,{id:"pages.execute.showCreate"})}),(0,e.jsx)(u.Z,{type:"default",htmlType:"button",onClick:function(){return _("showTableSize")},children:(0,e.jsx)(r._H,{id:"pages.execute.showTableSize"})})]})]})})]})]}),(0,e.jsxs)(y.Z,{children:[R==!1&&T!=""&&(0,e.jsx)(g.Z,{type:"error",message:(0,e.jsx)(r._H,{id:"pages.execute.queryFailed"})+T,banner:!0}),R==!0&&T!=""&&(0,e.jsx)(g.Z,{type:"success",message:"\u6267\u884C\u6210\u529F\uFF0C\u8017\u65F6\uFF1A"+Ot+"\u6BEB\u79D2,"+T,banner:!0}),R==!0&&ce>=0&&(0,e.jsxs)("div",{style:{whiteSpace:"pre-wrap",marginTop:"10px"},children:[(0,e.jsxs)("div",{style:{width:"100%",float:"right",marginBottom:"10px"},children:["\u67E5\u8BE2\u5230"+ce+"\u6761\u6570\u636E"," ",(0,e.jsx)(u.Z,{icon:(0,e.jsx)(w.Z,{}),onClick:kt,children:"\u67E5\u8BE2\u7ED3\u679C\u5BFC\u51FAExcel"})]}),(0,e.jsx)(Be.Z,{bordered:!0,loading:Ct,scroll:{scrollToFirstRowOnChange:!0,x:100},className:Qe().tableStyle,dataSource:he,columns:xe,size:"small"})]})]})]})]}),(0,e.jsx)(Ie.Z,{title:(0,e.jsx)(r._H,{id:"pages.execute.favorite"}),placement:"right",width:800,onClose:Me,visible:_t,extra:(0,e.jsx)(P.Z,{children:(0,e.jsx)(u.Z,{onClick:Me,children:(0,e.jsx)(r._H,{id:"pages.execute.close"})})}),children:(0,e.jsx)(Ke.Z,{rowKey:"id",search:!1,dataSource:ct,columns:$e,size:"middle"})})]})};W.default=Xe}}]);
