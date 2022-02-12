package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

const (
	endpoint  = "http://elasticmq:9324"
	region    = "ap-northeast-1"
	queueName = "local.fifo"
)

func RetrieveMessage() error {
	time.Sleep(3 * time.Second)

	fmt.Println("batch job start")

	client := getConnection()

	output, err := client.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		// ローカルで起動が間に合わない場合があるのでループしとく
		fmt.Println(err)
		return nil
	}

	params := &sqs.ReceiveMessageInput{
		QueueUrl: output.QueueUrl,
		// 一度に取得する最大メッセージ数。最大でも5まで。
		MaxNumberOfMessages: aws.Int64(5),
		// これでキューが空の場合はロングポーリング(20秒間繋ぎっぱなし)になる。
		WaitTimeSeconds: aws.Int64(20),
	}

	resp, err := client.ReceiveMessage(params)

	if err != nil {
		return err
	}

	if len(resp.Messages) == 0 {
		fmt.Println("empty queue.")
		return nil
	}

	// メッセージの数だけgoroutineで並列処理
	var wg sync.WaitGroup
	for _, m := range resp.Messages {
		wg.Add(1)
		go func(msg *sqs.Message) {
			defer wg.Done()
			if err := DeleteMessage(msg, output, client); err != nil {
				fmt.Println(err)
			}
		}(m)
	}

	wg.Wait()

	return nil
}

func getConnection() *sqs.SQS {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	client := sqs.New(sess, &aws.Config{
		Endpoint: aws.String(endpoint),
		Region:   aws.String(region),
	})
	return client
}

func DeleteMessage(msg *sqs.Message, url *sqs.GetQueueUrlOutput, conn *sqs.SQS) error {
	fmt.Println("start DeleteMessage.")
	time.Sleep(3 * time.Second)
	params := &sqs.DeleteMessageInput{
		QueueUrl:      url.QueueUrl,
		ReceiptHandle: aws.String(*msg.ReceiptHandle),
	}
	_, err := conn.DeleteMessage(params)

	if err != nil {
		return err
	}
	return nil
}
