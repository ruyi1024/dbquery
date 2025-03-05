package ai

import (
  "context"
	"fmt"
	"net/http"
	"log"

	"github.com/gin-gonic/gin"
	
  deepseek "github.com/cohesion-org/deepseek-go"
)

func AiRunDeepseek(c *gin.Context) {
  params := make(map[string]string)
  c.BindJSON(&params)
  if len(params) == 0 {
	  c.JSON(http.StatusOK, gin.H{"success": false, "msg": "params error."})
	  return
  }

  //datasourceType := params["datasource_type"]
  //datasource := params["datasource"]
  //databaseName := params["database"]
  //table := params["table"]
  sql := params["sql"]
  datanum:="10"
  table := "CREATE TABLE `datasource` ( `id` bigint(20) NOT NULL AUTO_INCREMENT,`name` varchar(50) DEFAULT NULL,`group_name` varchar(50) DEFAULT NULL,`idc` varchar(30) DEFAULT NULL,`env` varchar(30) DEFAULT NULL, `type` varchar(30) DEFAULT NULL,`host` varchar(100) DEFAULT NULL,PRIMARY KEY (`id`),UNIQUE KEY `idx_datasource_name` (`name`),UNIQUE KEY `uniq_host_port_dbid` (`host`,`port`,`dbid`) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8"

	client := deepseek.NewClient("xxxxxx")
  

	// Create a chat completion request
	request := &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleSystem, Content: "你是一名专业DBA，擅长SQL的优化，数据库类型为MySQL，请结合表结构的索引，数据量大小，字段类型等信息简单明了的给出优化建议，如果涉及到改成SQL，请返回优化后的SQL语句，要保证返回的SQL语法是正确的."},
			{Role: deepseek.ChatMessageRoleUser, Content: "需要优化的SQL语句为:" + sql+",表结构为："+table+",表里的数据量为:"+datanum},
		},
	}

	// Send the request and handle the response
	ctx := context.Background()
	response, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Print the response
	fmt.Println("Response:", response.Choices[0].Message.Content)
  c.JSON(http.StatusOK, gin.H{"success": true, "msg": response.Choices[0].Message.Content})
  return

}
