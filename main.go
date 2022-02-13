package main

import (
	"fmt"
	"log"
	"sync"

	as "github.com/aws/aws-sdk-go/service/sqs"
	"github.com/takakosi/minio_elasticmq/s3"
	"github.com/takakosi/minio_elasticmq/sqs"
)

func main() {

	// キューの接続
	sqs.InitSQS()
	s3.InitS3()

	// バッチの本体処理
	for {

		fmt.Println("batch start")
		resp, err := sqs.RetrieveMessage()
		if err != nil {
			log.Fatal(err)
		}

		//　対象がなかった場合ループ
		if resp == nil {
			continue
		}

		// メッセージの数だけgoroutineで並列処理
		var wg sync.WaitGroup
		for _, m := range resp.Messages {
			wg.Add(1)
			go func(msg *as.Message) {
				defer wg.Done()
				if err := s3.GetObjects(); err != nil {
					fmt.Println(err)
				}
				if err := sqs.DeleteMessage(msg); err != nil {
					fmt.Println(err)
				}
			}(m)
		}

		wg.Wait()
	}
}
