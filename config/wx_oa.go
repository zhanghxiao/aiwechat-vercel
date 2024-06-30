package config

import (
	"fmt"
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

	Bot_Type_Key   = "botType"
	Bot_Type_Gpt   = "gpt"
	Bot_Type_Spark = "spark"
	Bot_Type_Qwen  = "qwen"
	Bot_Type_Gemini  = "gemini"

	Wx_Command_1 = "/1"
	Wx_Command_2 = "/2"
	Wx_Command_3 = "/3"
	Wx_Command_4 = "/4"
	Wx_Command_5 = "/5"
	Wx_Command_6 = "/6"
	Wx_Command_7 = "/7"
)

var (
	Cache = NewCache()
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
	helpMsg := "输入以下命令进行对话\n/1：查看帮助\n/2：与GPT对话\n/3：与讯飞星火对话\n/4: 与通义千问对话\n/5: 与智谱清言对话\n/6: 与百度文心一言对话\n/7: 告诉我当前是哪个模型跟我聊天"
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

func GetGptToken() string {
	return os.Getenv("GPT_TOKEN")
}

func GetGeminiKey() string {
	return os.Getenv("GEMINI_KEY")
}

func GetBotType() string {
	botType := os.Getenv("BOT_TYPE")
	if botType == "" {
		botType = Bot_Type_Gpt
	}
	return botType
}

func GetMaxTokens() int {
	maxTokens := os.Getenv("MAX_TOKENS")
	if maxTokens == "" {
		return 1000
	}
	tokens, err := strconv.Atoi(maxTokens)
	if err != nil {
		return 1000
	}
	return tokens
}

func CheckBotConfig(botType string) (string, error) {
	switch botType {
	case Bot_Type_Gpt:
		token := GetGptToken()
		if token == "" {
			return "", fmt.Errorf("GPT_TOKEN not set")
		}
		return botType, nil
	case Bot_Type_Gemini:
		key := GetGeminiKey()
		if key == "" {
			return "", fmt.Errorf("GEMINI_KEY not set")
		}
		return botType, nil
	case Bot_Type_Spark:
		_, err := GetSparkConfig()
		if err != nil {
			return "", err
		}
		return botType, nil
	case Bot_Type_Qwen:
		_, err := GetQwenConfig()
		if err != nil {
			return "", err
		}
		return botType, nil
	default:
		return "", fmt.Errorf("unknown bot type:%s", botType)
	}
}

func GetUserBotType(userId string) string {
	botType, err := db.GetValue(fmt.Sprintf("%v:%v", Bot_Type_Key, userId))
	if err != nil {
		return GetBotType()
	}
	return botType
}

func GetBotWelcomeReply(botType string) string {
	switch botType {
	case Bot_Type_Gpt:
		return "切换到 GPT 模型"
	case Bot_Type_Spark:
		return "切换到 讯飞星火 模型"
	case Bot_Type_Qwen:
		return "切换到 通义千问 模型"
	case Bot_Type_Gemini:
		return "切换到 谷歌Gemini 模型"
	default:
		return "未知的模型类型"
	}
}

func GetSystemPrompt(userId, botType string) string {
	prompt, err := db.GetPrompt(userId, botType)
	if err != nil || prompt == "" {
		switch botType {
		case Bot_Type_Gpt:
			return "你是一个友好的助手。"
		case Bot_Type_Spark:
			return "你是一个友好的助手。"
		case Bot_Type_Qwen:
			return "你是一个友好的助手。"
		case Bot_Type_Gemini:
			return "你是一个友好的助手。"
		default:
			return "你是一个友好的助手。"
		}
	}
	return prompt
}
