package main

import (
	"fmt"
	"github.com/xinggaoya/qwen-sdk/qwen"
	qwenmodel "github.com/xinggaoya/qwen-sdk/qwenmodel"
	"testing"
)

func TestAdd(t *testing.T) {

	// 初始化QWEN聊天机器人客户端，使用您的API密钥
	apiKey := "sk-2dd0b09e04e44bda8acb42ec0219a81e"
	qwenclient := qwen.NewWithDefaultChat(apiKey)

	//qwenclient.QWenModel = "new model"

	// 定义一条消息对话的历史记录
	messages := []qwenmodel.Messages{
		{Role: qwenmodel.ChatUser, Content: "你好"},
		{Role: qwenmodel.ChatBot, Content: "你好！有什么我能为你做的吗？"},
		{Role: qwenmodel.ChatUser, Content: "我想买一件衬衫"},
	}

	// 获取AI对消息的回复
	resp := qwenclient.GetAIReply(messages)

	// 打印收到的回复
	fmt.Printf("收到的回复：%v\n", resp.Output.Text)

}
