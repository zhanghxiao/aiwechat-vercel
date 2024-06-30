package chat

import (
	_ "errors"
	"fmt"
	"github.com/pwh-pwh/aiwechat-vercel/client"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/pwh-pwh/aiwechat-vercel/db"
	"github.com/sashabaranov/go-openai"
	"github.com/pwh-pwh/aiwechat-vercel/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

var actionMap = map[string]func(param, userId string) string{
	"/1": func(param, userId string) string {
		return config.GetWxHelpReply()
	},
	"/2": func(param, userId string) string {
		return SwitchUserBot(userId, config.Bot_Type_Gpt)
	},
	"/3": func(param, userId string) string {
		return SetModel("讯飞星火v3.5", userId)
	},
	"/4": func(param, userId string) string {
		return SetModel("通义千问-plus", userId)
	},
	"/5": func(param, userId string) string {
		return SetModel("智谱glm-4", userId)
	},
	"/6": func(param, userId string) string {
		return SetModel("百度文心", userId)
	},
	"/7": func(param, userId string) string {
		return GetModel("", userId)
	},
}

func DoAction(userId, msg string) (r string, flag bool) {
	action, param, flag := isAction(msg)
	if flag {
		f := actionMap[action]
		r = f(param, userId)
	}
	return
}

func isAction(msg string) (string, string, bool) {
	for key := range actionMap {
		if strings.HasPrefix(msg, key) {
			return msg[:len(key)], strings.TrimSpace(msg[len(key):]), true
		}
	}
	return "", "", false
}

type BaseChat interface {
	Chat(userID string, msg string) string
	HandleMediaMsg(msg *message.MixMessage) string
}
type SimpleChat struct {
}

func (s SimpleChat) Chat(userID string, msg string) string {
	panic("implement me")
}

func (s SimpleChat) HandleMediaMsg(msg *message.MixMessage) string {
	switch msg.MsgType {
	case message.MsgTypeImage:
		return msg.PicURL
	case message.MsgTypeEvent:
		if msg.Event == message.EventSubscribe {
			subText := config.GetWxSubscribeReply() + config.GetWxHelpReply()
			if subText == "" {
				subText = "哇，又有帅哥美女关注我啦😄"
			}
			return subText
		} else if msg.Event == message.EventClick {
			switch msg.EventKey {
			case config.GetWxEventKeyChatGpt():
				return SwitchUserBot(string(msg.FromUserName), config.Bot_Type_Gpt)
			case config.GetWxEventKeyChatSpark():
				return SwitchUserBot(string(msg.FromUserName), config.Bot_Type_Spark)
			case config.GetWxEventKeyChatQwen():
				return SwitchUserBot(string(msg.FromUserName), config.Bot_Type_Qwen)
			default:
				return fmt.Sprintf("unkown event key=%v", msg.EventKey)
			}
		} else {
			return "不支持的类型"
		}
	default:
		return "未支持的类型"
	}
}

func SwitchUserBot(userId string, botType string) string {
	if _, err := config.CheckBotConfig(botType); err != nil {
		return err.Error()
	}
	db.SetValue(fmt.Sprintf("%v:%v", config.Bot_Type_Key, userId), botType, 0)
	return config.GetBotWelcomeReply(botType)
}

func SetPrompt(param, userId string) string {
	botType := config.GetUserBotType(userId)
	switch botType {
	case config.Bot_Type_Gpt:
		db.SetPrompt(userId, botType, param)
	case config.Bot_Type_Qwen:
		db.SetPrompt(userId, botType, param)
	case config.Bot_Type_Spark:
		db.SetPrompt(userId, botType, param)
	default:
		return fmt.Sprintf("%s 不支持设置system prompt", botType)
	}
	return fmt.Sprintf("%s 设置prompt成功", botType)
}

func RmPrompt(param string, userId string) string {
	botType := config.GetUserBotType(userId)
	db.RemovePrompt(userId, botType)
	return fmt.Sprintf("%s 删除prompt成功", botType)
}

func GetPrompt(param string, userId string) string {
	botType := config.GetUserBotType(userId)
	prompt, err := db.GetPrompt(userId, botType)
	if err != nil {
		return fmt.Sprintf("%s 当前未设置prompt", botType)
	}
	return fmt.Sprintf("%s 获取prompt成功，prompt：%s", botType, prompt)
}

func GetTodoList(param string, userId string) string {
	list, err := db.GetTodoList(userId)
	if err != nil {
		return err.Error()
	}
	return list
}

func AddTodo(param, userId string) string {
	err := db.AddTodoList(userId, param)
	if err != nil {
		return err.Error()
	}
	return "添加成功"
}

func DelTodo(param, userId string) string {
	index, err := strconv.Atoi(param)
	if err != nil {
		return "传入索引必须为数字"
	}
	err = db.DelTodoList(userId, index)
	if err != nil {
		return err.Error()
	}
	return "删除todo成功"
}

func GetCoin(param, userId string) string {
	coinPrice, err := client.GetCoinPrice(param)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("代币对:%s 价格:%s", coinPrice.Symbol, coinPrice.Price)
}

func SetModel(param, userId string) string {
	botType := config.GetUserBotType(userId)
	if botType == config.Bot_Type_Gpt || botType == config.Bot_Type_Gemini || botType == config.Bot_Type_Qwen {
		if err := db.SetModel(userId, botType, param); err != nil {
			return fmt.Sprintf("%s 设置model失败", botType)
		}
		return fmt.Sprintf("%s 设置model成功", botType)
	}
	return fmt.Sprintf("%s 不支持设置model", botType)
}

func GetModel(param string, userId string) string {
	botType := config.GetUserBotType(userId)
	model, err := db.GetModel(userId, botType)
	if err != nil || model == "" {
		return fmt.Sprintf("%s 当前未设置model", botType)
	}
	return fmt.Sprintf("%s 获取model成功，model：%s", botType, model)
}

func ClearMsg(param string, userId string) string {
	botType := config.GetUserBotType(userId)
	db.DeleteMsgList(botType, userId)
	return fmt.Sprintf("%s 清除消息成功", botType)
}

// 加入超时控制
func WithTimeChat(userID, msg string, f func(userID, msg string) string) string {
	if _, ok := config.Cache.Load(userID + msg); ok {
		rAny, _ := config.Cache.Load(userID + msg)
		r := rAny.(string)
		config.Cache.Delete(userID + msg)
		return r
	}
	resChan := make(chan string)
	go func() {
		resChan <- f(userID, msg)
	}()
	select {
	case res := <-resChan:
		return res
	case <-time.After(5 * time.Second):
		config.Cache.Store(userID+msg, <-resChan)
		return ""
	}
}

type ErrorChat struct {
	errMsg string
}

func (e *ErrorChat) HandleMediaMsg(msg *message.MixMessage) string {
	return e.errMsg
}

func (e *ErrorChat) Chat(userID string, msg string) string {
	return e.errMsg
}

func GetChatBot(botType string) BaseChat {
	if botType == "" {
		botType = config.GetBotType()
	}
	var err error
	botType, err = config.CheckBotConfig(botType)
	if err != nil {
		return &ErrorChat{
			errMsg: err.Error(),
		}
	}
	maxTokens := config.GetMaxTokens()

	switch botType {
	case config.Bot_Type_Gpt:
		url := os.Getenv("GPT_URL")
		if url == "" {
			url = "https://api.openai.com/v1/"
		}
		return &SimpleGptChat{
			token:     config.GetGptToken(),
			url:       url,
			maxTokens: maxTokens,
			BaseChat:  SimpleChat{},
		}
	case config.Bot_Type_Gemini:
		return &GeminiChat{
			BaseChat:  SimpleChat{},
			key:       config.GetGeminiKey(),
			maxTokens: maxTokens,
		}
	case config.Bot_Type_Spark:
		config, _ := config.GetSparkConfig()
		return &SparkChat{
			BaseChat:  SimpleChat{},
			Config:    config,
			maxTokens: maxTokens,
		}
	case config.Bot_Type_Qwen:
		config, _ := config.GetQwenConfig()
		return &QwenChat{
			BaseChat:  SimpleChat{},
			Config:    config,
			maxTokens: maxTokens,
		}
	default:
		return &ErrorChat{
			errMsg: fmt.Sprintf("unknown bot type:%s", botType),
		}
	}
}

type SimpleGptChat struct {
	BaseChat
	token     string
	url       string
	maxTokens int
}

func (s *SimpleGptChat) Chat(userID string, msg string) string {
	return WithTimeChat(userID, msg, s.chat)
}

func (s *SimpleGptChat) chat(userID string, msg string) string {
	gptReq := openai.ChatCompletionRequest{
		Model:     "gpt-3.5-turbo",
		MaxTokens: s.maxTokens,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: config.GetSystemPrompt(userID, config.Bot_Type_Gpt),
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: msg,
			},
		},
	}
	gptRsp, err := client.GptChatCompletions(s.token, gptReq)
	if err != nil {
		return err.Error()
	}
	return gptRsp.Choices[0].Message.Content
}

func (s *SimpleGptChat) HandleMediaMsg(msg *message.MixMessage) string {
	return "暂不支持此类消息哦"
}

type GeminiChat struct {
	BaseChat
	key       string
	maxTokens int
}

func (s *GeminiChat) Chat(userID string, msg string) string {
	return WithTimeChat(userID, msg, s.chat)
}

func (s *GeminiChat) chat(userID string, msg string) string {
	client := genai.NewClient(s.key)
	completionRequest := &genai.GenerateMessageRequest{
		Messages: []genai.Message{
			{
				Author:  "system",
				Content: config.GetSystemPrompt(userID, config.Bot_Type_Gemini),
			},
			{
				Author:  "user",
				Content: msg,
			},
		},
		Temperature: 0.5,
		TopK:        40,
		TopP:        0.95,
	}
	response, err := client.GenerateMessage(completionRequest)
	if err != nil {
		return err.Error()
	}
	if len(response.Messages) > 0 {
		return response.Messages[0].Content
	} else {
		return "response is empty"
	}
}

func (s *GeminiChat) HandleMediaMsg(msg *message.MixMessage) string {
	return "暂不支持此类消息哦"
}

type QwenChat struct {
	BaseChat
	Config    config.QwenConfig
	maxTokens int
}

func (s *QwenChat) Chat(userID string, msg string) string {
	return WithTimeChat(userID, msg, s.chat)
}

func (s *QwenChat) chat(userID string, msg string) string {
	systemPrompt := config.GetSystemPrompt(userID, config.Bot_Type_Qwen)
	rsp, err := client.QwenChat(s.Config, msg, systemPrompt, s.maxTokens)
	if err != nil {
		return err.Error()
	}
	return rsp
}

func (s *QwenChat) HandleMediaMsg(msg *message.MixMessage) string {
	return "暂不支持此类消息哦"
}

type SparkChat struct {
	BaseChat
	Config    config.SparkConfig
	maxTokens int
}

func (s *SparkChat) Chat(userID string, msg string) string {
	return WithTimeChat(userID, msg, s.chat)
}

func (s *SparkChat) chat(userID string, msg string) string {
	systemPrompt := config.GetSystemPrompt(userID, config.Bot_Type_Spark)
	rsp, err := client.SparkChat(s.Config, msg, systemPrompt, s.maxTokens)
	if err != nil {
		return err.Error()
	}
	return rsp
}

func (s *SparkChat) HandleMediaMsg(msg *message.MixMessage) string {
	return "暂不支持此类消息哦"
}
