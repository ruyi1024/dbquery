(self.webpackChunkant_design_pro=self.webpackChunkant_design_pro||[]).push([[861],{38373:function(Re,B,r){"use strict";r.r(B),r.d(B,{default:function(){return Ee}});var Ve=r(58024),ae=r(91894),ze=r(66456),te=r(15885),Oe=r(13062),K=r(71230),Ae=r(89032),I=r(15746),Be=r(57663),N=r(71577),Ke=r(49111),se=r(19650),Ne=r(62350),ne=r(24565),Le=r(22385),Y=r(45777),We=r(48736),ue=r(27049),Ge=r(71153),le=r(60331),Z=r(2824),m=r(11849),Je=r(34792),g=r(48086),o=r(39428),y=r(3182),Qe=r(47673),E=r(4107),p=r(67294),ie=r(67265),$=r(21704),de=r(80129);function oe(i){return w.apply(this,arguments)}function w(){return w=(0,y.Z)((0,o.Z)().mark(function i(n){return(0,o.Z)().wrap(function(t){for(;;)switch(t.prev=t.next){case 0:return t.abrupt("return",(0,$.WY)("/api/v1/users/manager/lists?".concat((0,de.stringify)(n))));case 1:case"end":return t.stop()}},i)})),w.apply(this,arguments)}function ce(i){return R.apply(this,arguments)}function R(){return R=(0,y.Z)((0,o.Z)().mark(function i(n){return(0,o.Z)().wrap(function(t){for(;;)switch(t.prev=t.next){case 0:return t.abrupt("return",(0,$.WY)("/api/v1/users/manager/lists",{method:n.modify?"PUT":"POST",data:(0,m.Z)({},n)}));case 1:case"end":return t.stop()}},i)})),R.apply(this,arguments)}function me(i){return V.apply(this,arguments)}function V(){return V=(0,y.Z)((0,o.Z)().mark(function i(n){return(0,o.Z)().wrap(function(t){for(;;)switch(t.prev=t.next){case 0:return t.abrupt("return",(0,$.WY)("/api/v1/users/manager/lists",{method:"DELETE",data:(0,m.Z)({},n)}));case 1:case"end":return t.stop()}},i)})),V.apply(this,arguments)}var he=r(59879),fe=r(49101),ve=r(30381),L=r.n(ve),Xe=r(71194),pe=r(50146),ke=r(63185),Ze=r(9676),qe=r(9715),S=r(55843),e=r(85893),ge={labelCol:{span:5},wrapperCol:{span:16}},ye=function(n){var l=n.updateModalVisible,t=n.onSubmit,u=n.onCancel,s=n.values,d=S.Z.useForm(),z=(0,Z.Z)(d,1),x=z[0],O=(0,p.useState)(),P=(0,Z.Z)(O,2),H=P[0],D=P[1];return(0,p.useEffect)(function(){s!==null?(console.log("values:",s),D(s),x.setFieldsValue((0,m.Z)({},s))):x.resetFields()},[s]),(0,e.jsx)(pe.Z,{destroyOnClose:!0,width:500,title:s.modify?"\u4FEE\u6539\u7528\u6237".concat(s.username):"\u65B0\u589E\u7528\u6237",visible:l,onCancel:u,onOk:function(){x.validateFields().then(function(T){x.resetFields();var M=(0,m.Z)((0,m.Z)({},T),{},{modify:s.modify});s.modify&&(M.id=s.id||0),t((0,m.Z)({},M))}).catch(function(T){console.log("Validate Failed:",T)})},children:(0,e.jsxs)(S.Z,(0,m.Z)((0,m.Z)({},ge),{},{form:x,initialValues:H?(0,m.Z)({},H):{},preserve:!1,children:[(0,e.jsx)(S.Z.Item,{name:"username",label:"\u7528\u6237",rules:[{required:!0}],children:(0,e.jsx)(E.Z,{style:{width:180}})}),(0,e.jsx)(S.Z.Item,{name:"chineseName",label:"\u59D3\u540D",rules:[{required:!0}],children:(0,e.jsx)(E.Z,{style:{width:180}})}),(0,e.jsx)(S.Z.Item,{name:"password",label:"\u5BC6\u7801",rules:[{required:!s.modify}],children:(0,e.jsx)(E.Z.Password,{})}),(0,e.jsx)(S.Z.Item,{name:"admin",valuePropName:"checked",label:"\u7BA1\u7406\u5458",children:(0,e.jsx)(Ze.Z,{})}),(0,e.jsx)(S.Z.Item,{name:"remark",label:"\u5907\u6CE8",children:(0,e.jsx)(E.Z.TextArea,{})})]}))})},je=ye,h=r(85029),Se=E.Z.Search,_e=null,xe=function(){var i=(0,y.Z)((0,o.Z)().mark(function n(l){return(0,o.Z)().wrap(function(u){for(;;)switch(u.prev=u.next){case 0:return u.prev=0,u.next=3,oe(l);case 3:return u.abrupt("return",u.sent);case 6:return u.prev=6,u.t0=u.catch(0),u.abrupt("return",{success:!1,msg:u.t0});case 9:case"end":return u.stop()}},n,null,[[0,6]])}));return function(l){return i.apply(this,arguments)}}(),Ce=function(){var i=(0,y.Z)((0,o.Z)().mark(function n(l){var t,u;return(0,o.Z)().wrap(function(d){for(;;)switch(d.prev=d.next){case 0:return t=g.ZP.loading("\u6B63\u5728\u914D\u7F6E"),d.prev=1,d.next=4,ce((0,m.Z)({},l));case 4:return u=d.sent,t(),g.ZP.success("\u914D\u7F6E\u6210\u529F"),d.abrupt("return",u);case 10:return d.prev=10,d.t0=d.catch(1),t(),g.ZP.error("\u914D\u7F6E\u5931\u8D25\u8BF7\u91CD\u8BD5\uFF01"),d.abrupt("return",{success:!1,msg:d.t0});case 15:case"end":return d.stop()}},n,null,[[1,10]])}));return function(l){return i.apply(this,arguments)}}(),Fe=function(){var i=(0,y.Z)((0,o.Z)().mark(function n(l){var t;return(0,o.Z)().wrap(function(s){for(;;)switch(s.prev=s.next){case 0:return t=g.ZP.loading("\u6B63\u5728\u5220\u9664"),s.prev=1,s.next=4,me({username:l});case 4:return t(),g.ZP.success("\u5220\u9664\u6210\u529F\uFF0C\u5373\u5C06\u5237\u65B0"),s.abrupt("return",!0);case 9:return s.prev=9,s.t0=s.catch(1),t(),g.ZP.error("\u5220\u9664\u5931\u8D25\uFF0C\u8BF7\u91CD\u8BD5"),s.abrupt("return",!1);case 14:case"end":return s.stop()}},n,null,[[1,9]])}));return function(l){return i.apply(this,arguments)}}(),Te=function(){var n=(0,p.useState)([]),l=(0,Z.Z)(n,2),t=l[0],u=l[1],s=(0,p.useState)(0),d=(0,Z.Z)(s,2),z=d[0],x=d[1],O=(0,p.useState)(!1),P=(0,Z.Z)(O,2),H=P[0],D=P[1],W=(0,p.useState)(!1),T=(0,Z.Z)(W,2),M=T[0],b=T[1],Pe=(0,p.useState)(),G=(0,Z.Z)(Pe,2),Ue=G[0],A=G[1],He=(0,p.useState)(),J=(0,Z.Z)(He,2),Q=J[0],De=J[1],Me=(0,p.useState)(1),X=(0,Z.Z)(Me,2),k=X[0],be=X[1],Ie=(0,p.useState)(10),q=(0,Z.Z)(Ie,2),_=q[0],Ye=q[1],ee=(0,p.useRef)(),re=(0,h.md)(),U=function(a){D(!0);var c=(0,m.Z)({offset:_*(k>=2?k-1:0),limit:_,keyword:a&&a.keyword?a.keyword:Q},a);console.log("debug did data -->",c),xe(c).then(function(v){v.success&&(u(v.data),x(v.total)),D(!1)})},$e=[{title:(0,e.jsx)(h._H,{id:"pages.searchTable.column.username"}),dataIndex:"username",sorter:!0,render:function(a){return(0,e.jsx)("a",{children:a})}},{title:(0,e.jsx)(h._H,{id:"pages.searchTable.column.chineseName"}),dataIndex:"chineseName",sorter:!0},{title:(0,e.jsx)(h._H,{id:"pages.searchTable.column.admin"}),dataIndex:"admin",sorter:!0,render:function(a){return(0,e.jsx)(le.Z,{color:a?"green":"",children:a?(0,e.jsx)(h._H,{id:"pages.searchTable.column.yes"}):(0,e.jsx)(h._H,{id:"pages.searchTable.column.no"})})}},{title:(0,e.jsx)(h._H,{id:"pages.searchTable.column.gmtCreated"}),dataIndex:"createdAt",sorter:!0,render:function(a){return L()(a).format("YYYY-MM-DD HH:mm:ss")}},{title:(0,e.jsx)(h._H,{id:"pages.searchTable.column.gmtUpdated"}),dataIndex:"updatedAt",sorter:!0,render:function(a){return L()(a).format("YYYY-MM-DD HH:mm:ss")}},{title:(0,e.jsx)(h._H,{id:"pages.searchTable.column.operate"}),dataIndex:"id",key:"id",fixed:"right",width:150,render:function(a,c){return(0,e.jsx)(e.Fragment,{children:(0,e.jsxs)(se.Z,{split:(0,e.jsx)(ue.Z,{type:"vertical"}),children:[(0,e.jsx)(Y.Z,{title:"\u4FEE\u6539\u7528\u6237\u3010".concat(c.username,"\u3011"),children:(0,e.jsx)("a",{onClick:function(){console.log("debug ---> ",c),b(!0),A((0,m.Z)((0,m.Z)({},c),{},{modify:!0}))},children:(0,e.jsx)(h._H,{id:"pages.searchTable.operate.edit"})})}),(0,e.jsx)(Y.Z,{title:"\u5220\u9664\u7528\u6237\u3010".concat(c.username,"\u3011"),children:(0,e.jsx)(ne.Z,{title:"\u5220\u9664\u3010".concat(c.username,"\u3011\uFF0C\u5220\u9664\u540E\u6570\u636E\u4E0D\u53EF\u6062\u590D\u3002\u662F\u5426\u7EE7\u7EED\uFF1F"),placement:"left",onConfirm:(0,y.Z)((0,o.Z)().mark(function v(){var C;return(0,o.Z)().wrap(function(F){for(;;)switch(F.prev=F.next){case 0:if(re.canAdmin){F.next=3;break}return g.ZP.error("\u64CD\u4F5C\u6743\u9650\u53D7\u9650\uFF0C\u8BF7\u8054\u7CFB\u5E73\u53F0\u7BA1\u7406\u5458"),F.abrupt("return");case 3:return F.next=5,Fe(c.username);case 5:C=F.sent,C&&ee.current&&ee.current.reload();case 7:case"end":return F.stop()}},v)})),children:(0,e.jsx)("a",{children:(0,e.jsx)(h._H,{id:"pages.searchTable.operate.delete"})})})})]})})}}];(0,p.useEffect)(function(){U("")},[]);var we=function(a,c,v){var C={offset:a.pageSize*(a.current>=2?a.current-1:0),limit:a.pageSize,keyword:Q,sorterField:"",sorterOrder:""};v.field&&(C.sorterField="".concat(v.field),C.sorterOrder="".concat(v.order)),be(a.current),Ye(a.pageSize),U(C)};return(0,e.jsxs)(ie.ZP,{children:[(0,e.jsxs)(ae.Z,{size:"small",bodyStyle:{padding:10},children:[(0,e.jsxs)(K.Z,{children:[(0,e.jsxs)(I.Z,{flex:"auto",children:[(0,e.jsx)(Se,{placeholder:(0,e.jsx)(h._H,{id:"pages.searchTable.operate.searchUser"}),onSearch:function(a){console.log("debug on search --> ",a),De(a),U({keyword:a})},style:{width:280}}),(0,e.jsx)(Y.Z,{placement:"top",title:"\u91CD\u8F7D\u5E76\u5237\u65B0\u8868\u683C\u6570\u636E",children:(0,e.jsx)(N.Z,{type:"link",icon:(0,e.jsx)(he.Z,{}),onClick:function(){return U("")}})})]}),(0,e.jsx)(I.Z,{span:2,children:(0,e.jsx)(N.Z,{type:"link",icon:(0,e.jsx)(fe.Z,{}),onClick:function(){return b(!0)},children:(0,e.jsx)(h._H,{id:"pages.searchTable.operate.create"})})})]}),(0,e.jsx)(K.Z,{style:{paddingTop:10},children:(0,e.jsx)(I.Z,{span:24,children:(0,e.jsx)(te.Z,{size:"small",rowKey:"id",loading:H,columns:$e,dataSource:t,onChange:we,pagination:{total:z,showSizeChanger:!0,pageSizeOptions:["10","20","50","100","200"],showQuickJumper:!0,showTotal:function(a,c){return"\u7B2C ".concat(c[0],"-").concat(c[1],"\u6761\uFF0C \u5171 ").concat(a,"\u6761")}}})})})]}),(0,e.jsx)(je,{onSubmit:function(){var f=(0,y.Z)((0,o.Z)().mark(function a(c){var v;return(0,o.Z)().wrap(function(j){for(;;)switch(j.prev=j.next){case 0:if(re.canAdmin){j.next=3;break}return g.ZP.error("\u64CD\u4F5C\u6743\u9650\u53D7\u9650\uFF0C\u8BF7\u8054\u7CFB\u5E73\u53F0\u7BA1\u7406\u5458"),j.abrupt("return");case 3:return j.next=5,Ce(c);case 5:v=j.sent,v.success&&(U(""),b(!1),A(void 0));case 7:case"end":return j.stop()}},a)}));return function(a){return f.apply(this,arguments)}}(),onCancel:function(){b(!1),A(void 0)},updateModalVisible:M,values:Ue||{}})]})},Ee=Te}}]);
