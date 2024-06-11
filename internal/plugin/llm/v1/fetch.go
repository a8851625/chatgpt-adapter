package v1

import (
	"github.com/bincooo/chatgpt-adapter/internal/common"
	"github.com/bincooo/chatgpt-adapter/pkg"
	"github.com/bincooo/emit.io"
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
)

func fetch(ctx *gin.Context, proxies, token string, completion pkg.ChatCompletion) (*http.Response, error) {
	
	toolCall := pkg.Config.GetStringMap("toolCall")
	fmt.Print("%+v\n", toolCall)

	baseUrl, ok := toolCall["baseurl"].(string)
	if !ok {
		panic("toolCall baseUrl is not set")
	}

	// 从header 获取 authori zation token
	// 再根据token 路由到不同的 base url

	if completion.TopP == 0 {
		completion.TopP = 1
	}

	if completion.Temperature == 0 {
		completion.Temperature = 0.7
	}

	if completion.MaxTokens == 0 {
		completion.MaxTokens = 1024
	}

	tokens := 0
	for _, message := range completion.Messages {
		tokens += common.CalcTokens(message.GetString("content"))
	}
	ctx.Set(ginTokens, token)

	completion.Stream = true
	return emit.ClientBuilder().
		Context(ctx.Request.Context()).
		Proxies(proxies).
		POST(baseUrl+"/v1/chat/completions").
		Header("Authorization", "Bearer "+token).
		JHeader().
		Body(completion).
		DoC(emit.Status(http.StatusOK), emit.IsSTREAM)
}
