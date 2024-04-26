package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/prompts"
)

func main() {
	g := gin.Default()
	v1 := g.Group("/api/v1")
	v1.POST("/translate", translator)

	g.Run()
}

func translator(c *gin.Context) {
	var requestData struct {
		OutputLang string `json:"output_lang"`
		Text       string `json:"text"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	prompt := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
		prompts.NewSystemMessagePromptTemplate("你是一个只能翻译文本的翻译引擎，不需要解释。", nil),
		prompts.NewHumanMessagePromptTemplate(`翻译这段文字到 {{.outputLang}}:{{.text}}`, []string{"outputLang", "text"}),
	})

	vals := map[string]any{
		"outputLang": requestData.OutputLang,
		"text":       requestData.Text,
	}

	msg, err := prompt.FormatMessages(vals)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	llm, err := ollama.New(ollama.WithModel("qwen"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	content := []llms.MessageContent{
		llms.TextParts(msg[0].GetType(), msg[0].GetContent()),
		llms.TextParts(msg[1].GetType(), msg[1].GetContent()),
	}

	resp, err := llm.GenerateContent(context.Background(), content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": resp.Choices[0].Content,
	})
}
