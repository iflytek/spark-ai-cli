package llm

import (
	"context"
	"fmt"
	"github.com/iflytek/spark-ai-go/sparkai/llms/spark"
	"github.com/iflytek/spark-ai-go/sparkai/llms/spark/client/sparkclient"
	"github.com/iflytek/spark-ai-go/sparkai/messages"
	"github.com/joho/godotenv"
	"time"
)

func init() {
	godotenv.Load(".env")

}

func LlmInit() error {

	return nil
}

func LLMRequest(question string, resultChan chan string, finish chan bool) (finResult string) {
	SPARK_API_KEY := ""
	SPARK_API_SECRET := ""
	SPARK_API_BASE := ""
	SPARK_APP_ID := ""
	SPARK_DOMAIN := ""

	defer func() {
		time.Sleep(1 * time.Second)
		finish <- true
		/*close(resultChan)
		close(finish)*/
	}()

	_, client, err := spark.NewClient(spark.WithBaseURL(SPARK_API_BASE), spark.WithApiKey(SPARK_API_KEY), spark.WithApiSecret(SPARK_API_SECRET), spark.WithAppId(SPARK_APP_ID), spark.WithAPIDomain(SPARK_DOMAIN))
	if err != nil {
		panic(err.Error())
		return
	}

	ctx := context.Background()
	r := &sparkclient.ChatRequest{
		Domain: &SPARK_DOMAIN,
		Messages: []messages.ChatMessage{
			&messages.GenericChatMessage{
				Role:    "user",
				Content: question,
			},
		},
	}

	_, err = client.CreateChatWithCallBack(ctx, r, func(msg messages.ChatMessage) error {
		// result process:
		finResult += msg.GetContent()
		resultChan <- msg.GetContent()
		return nil
	})

	if err != nil {
		fmt.Print(err.Error())
		return
	}
	return
}
