package ai

import (
	"fmt"
	"net/http"
	"log"

	"github.com/gin-gonic/gin"
	
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/enums/option"
)

func AiRun(c *gin.Context) {
	ollamaUrl := "http://localhost:11434"
	model := "deepseek-r1:1.5b"

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature: 0.5,
  })

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


  firstQuestion := llm.GenQuery{
	  Model: model,
	  Prompt: "你是一名专业DBA，擅长SQL的优化，数据库类型为MySQL，请首先简单明了的给出优化建议，如果涉及到改成SQL，请返回优化后的SQL语句，要保证返回的SQL语法是正确的，需要优化的SQL语句为:" + sql + "。",
	  Options: options,
  }

  answer, err := completion.Generate(ollamaUrl, firstQuestion)
  if err != nil {
	  log.Fatal("😡:", err)
  }
  fmt.Println(answer.Response)

  c.JSON(http.StatusOK, gin.H{"success": true, "msg": answer.Response})
  return

  /*
  fmt.Println()

  secondQuestion := llm.GenQuery{
	  Model: model,
	  Prompt: "Who is his best friend?",
	  Context: answer.Context,
	  Options: options,
  }

  answer, err = completion.Generate(ollamaUrl, secondQuestion)
  if err != nil {
	  log.Fatal("😡:", err)
  }
  fmt.Println(answer.Response)
  */
}
