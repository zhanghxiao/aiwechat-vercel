package config

import (
	"os"
	"strings"
)

const (
	Wx_Token_key           = "WX_TOKEN"
	Wx_App_Id_key          = "WX_APP_ID"
	Wx_App_Secret_key      = "WX_APP_SECRET"
	Wx_Subscribe_Reply_key = "WX_SUBSCRIBE_REPLY"
	Wx_Help_Reply_key      = "WX_HELP_REPLY"

	Wx_Event_Key_Chat_Gpt_key   = "AI_CHAT_GPT"
	Wx_Event_Key_Chat_Spark_key = "AI_CHAT_SPARK"
	Wx_Event_Key_Chat_Qwen_key  = "AI_CHAT_QWEN"

	Wx_Command_Gemini    = "/gemini"
	Wx_Command_Prompt    = "/prompt"
	Wx_Command_RmPrompt  = "/cpt"
	Wx_Command_GetPrompt = "/getpt"
	Wx_Command_SetModel  = "/setmodel"
	Wx_Command_Clear     = "/clear"

	Wx_Todo_Add  = "/ta"
	Wx_Todo_Del  = "/td"
	Wx_Todo_List = "/tl"

	Wx_Coin = "/cb"
	    Wx_Command_Help   = "/1"
	    Wx_Command_Gpt    = "/2"
	    Wx_Command_Spark  = "/3"
	    Wx_Command_Qwen   = "/4"
	    Wx_Command_Zhipu  = "/5"
	    Wx_Command_Baidu  = "/6"
	    Wx_Command_GetModel = "/7"
)

var (
	Wx_Commands = []string{Wx_Command_Help, Wx_Command_Gpt, Wx_Command_Spark, Wx_Command_Qwen, Wx_Command_Gemini}
)

func GetWxToken() string {
	return os.Getenv(Wx_Token_key)
}
func GetWxAppId() string {
	return os.Getenv(Wx_App_Id_key)
}
func GetWxAppSecret() string {
	return os.Getenv(Wx_App_Secret_key)
}
func GetWxSubscribeReply() string {
	subscribeMsg := os.Getenv(Wx_Subscribe_Reply_key)
	return strings.ReplaceAll(subscribeMsg, "\\n", "\n")
}
func GetWxHelpReply() string {
    helpMsg := "输入以下命令进行对话\n/1：查看帮助\n/2：与GPT对话\n/3：与讯飞星火对话\n/4: 与通义千问对话\n/5:与智谱清言对话\n" +
        "/6 :与百度文心一言对话\n/7:告诉我当前是哪个模型跟我聊天"
    return strings.ReplaceAll(helpMsg, "\\n", "\n")
}
func GetWxEventKeyChatGpt() string {
	return os.Getenv(Wx_Event_Key_Chat_Gpt_key)
}
func GetWxEventKeyChatSpark() string {
	return os.Getenv(Wx_Event_Key_Chat_Spark_key)
}
func GetWxEventKeyChatQwen() string {
	return os.Getenv(Wx_Event_Key_Chat_Qwen_key)
}

func GetBotWelcomeReply(botType string) string {
	switch botType {
	case Bot_Type_Gpt:
		return GetGptWelcomeReply()
	case Bot_Type_Gemini:
		return GetGeminiWelcomeReply()
	case Bot_Type_Spark:
		return GetSparkWelcomeReply()
	case Bot_Type_Qwen:
		return GetQwenWelcomeReply()
	}

	return botType
}
