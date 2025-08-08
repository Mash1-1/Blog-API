package infrastructure

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/genai"
)

type AI struct {
	model_name string
	thinking   int
	Ai_client  *genai.Client
	config     *genai.GenerateContentConfig
}

func InitAIClient() *AI {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText("You are a professional tech writer who creates well-structured and engaging blog posts using markdown formatting.", genai.RoleUser),
	}

	model_name := os.Getenv("GEMINI_MODEL")

	return &AI{
		model_name: model_name,
		config:     config,
		Ai_client:  client,
	}
}

func (ai *AI) Generate_blog_content(message string) (*string, error) {

	userPrompt := genai.Text(fmt.Sprintf(`
	Write a detailed and engaging blog post about "%s".
	- Include an introduction, 3 main sections, and a conclusion.
	- Make it informative, creative, and easy to read.
	- Use markdown formatting for headers and bullet points.
	`, message))

	result, err := ai.Ai_client.Models.GenerateContent(
		context.Background(),
		ai.model_name,
		userPrompt,
		ai.config,
	)

	if err != nil {
		return nil, err
	}

	result_message := result.Text()
	return &result_message, nil
}
