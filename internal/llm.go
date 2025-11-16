package internal

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"go.uber.org/zap"
	"net/http"
	"os"
)

type llmClient struct {
	client openai.Client
	logger *zap.Logger

	model string
}

func newLLMClient(logger *zap.Logger) *llmClient {
	llmBaseURL := os.Getenv("LLM_BASE_URL")
	llmAPIKey := os.Getenv("LLM_API_KEY")
	llmModel := os.Getenv("LLM_MODEL")

	if llmBaseURL == "" || llmAPIKey == "" || llmModel == "" {
		panic("invalid LLM_BASE_URL or LLM_API_KEY or LLM_MODEL")
	}

	c := &llmClient{
		model: llmModel,
		client: openai.NewClient(
			option.WithAPIKey(llmAPIKey),
			option.WithBaseURL(llmBaseURL)),
		logger: logger,
	}

	logger.Info("LLM client created",
		zap.String("base_url", llmBaseURL),
		zap.String("model", llmModel),
	)
	return c
}

func (c *llmClient) encodeImageToDataURL(filepath string) (string, error) {
	fileBytes, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	mimeType := http.DetectContentType(fileBytes)
	base64String := base64.StdEncoding.EncodeToString(fileBytes)
	dataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, base64String)
	return dataURL, nil
}

func (c *llmClient) extractVerbsFromImage(imagePath string) []*VerbRecord {
	ctx := context.Background()

	var (
		err      error
		response *openai.ChatCompletion
	)

	c.logger.Debug("extracting verbs from image", zap.String("image_path", imagePath))
	dataURL, err := c.encodeImageToDataURL(imagePath)
	if err != nil {
		c.logger.Error("failed to encode image", zap.Error(err), zap.String("imagePath", imagePath))
	}

	response, err = c.client.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(PromptExtractVerbsAndTranslate),
				openai.UserMessage(
					[]openai.ChatCompletionContentPartUnionParam{
						{
							OfImageURL: &openai.ChatCompletionContentPartImageParam{
								ImageURL: openai.ChatCompletionContentPartImageImageURLParam{
									URL: dataURL,
								},
							},
						},
					},
				),
			},
			Model: c.model,
		},
	)
	if err != nil {
		c.logger.Error("failed to process image in llm", zap.Error(err), zap.String("imagePath", imagePath))
	}

	r := response.Choices[0].Message.Content

	res := []*VerbRecord{}
	if err := json.Unmarshal([]byte(r), &res); err != nil {
		c.logger.Error("can't unmarshal results",
			zap.Error(err),
		)
	}

	c.logger.Debug("image processed",
		zap.String("imagePath", imagePath),
		zap.Int("verbs extracted", len(res)),
	)

	return res
}

func (c *llmClient) addExampleSentences(verb *VerbRecord) (*VerbRecord, error) {
	var (
		err      error
		response *openai.ChatCompletion
	)

	ctx := context.Background()

	verbStr, _ := json.Marshal(verb)
	response, err = c.client.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(PromptAddExampleSentences),
				openai.UserMessage(string(verbStr)),
			},
			Model: c.model,
		},
	)
	if err != nil {
		return nil, err
	}
	
	res := VerbRecord{}
	err = json.Unmarshal([]byte(response.Choices[0].Message.Content), &res)
	if err != nil {
		c.logger.Error("can't unmarshal results",
			zap.Error(err),
		)
	}
	return &res, nil
}
