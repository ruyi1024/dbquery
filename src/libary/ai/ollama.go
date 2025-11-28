package ai

import (
	"dbmcloud/setting"
	"log"
	"strings"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/llm"
)

func Text2SQL(text string) (string, error) {
	ollamaUrl := setting.Setting.OllamaUrl
	model := setting.Setting.OllamaModel

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature: 0.2,
	})

	firstQuestion := llm.GenQuery{
		Model:   model,
		Prompt:  "你是一名专业DBA，负责执行SQL语句，你的任务是将用户执行SQL的大白话转换成具体的SQL语句，如果可以转换，则只返回SQL即可,如果已经是SQL语句，就什么都不做，原样返回SQL字符即可，如果用户执行的SQL不是查询数据库，则返回0就行。用户的输入为:" + text + "。",
		Options: options,
	}

	answer, err := completion.Generate(ollamaUrl, firstQuestion)
	if err != nil {
		log.Fatal("😡:", err)
		return "", err
	}
	//fmt.Println(strings.Split(answer.Response, "</think>")[1])

	return strings.Split(answer.Response, "</think>")[1], nil

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
