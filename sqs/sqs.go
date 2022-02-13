package sqs

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	la "github.com/takakosi/minio_elasticmq/aws"
)

const (
	endpoint  = "http://elasticmq:9324"
	queueName = "local.fifo"
)

var (
	svc            *sqs.SQS
	err            error
	queueUrlOutput *sqs.GetQueueUrlOutput
)

func RetrieveMessage() (*sqs.ReceiveMessageOutput, error) {

	params := &sqs.ReceiveMessageInput{
		QueueUrl: queueUrlOutput.QueueUrl,
		// 一度に取得する最大メッセージ数。最大でも5まで。
		MaxNumberOfMessages: aws.Int64(5),
		// ロングポーリング設定。20秒繋ぎっぱなし。
		WaitTimeSeconds: aws.Int64(20),
	}

	resp, err := svc.ReceiveMessage(params)

	if err != nil {
		return nil, err
	}

	if len(resp.Messages) == 0 {
		fmt.Println("empty queue.")
		return nil, nil
	}

	return resp, nil
}

func getSQSConnection() {
	sess := la.CreateSession()
	svc = sqs.New(sess, &aws.Config{
		Endpoint: aws.String(endpoint),
	})
}

func getQueueUrl() (*sqs.GetQueueUrlOutput, error) {
	return svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
}

func DeleteMessage(msg *sqs.Message) error {
	fmt.Println("start DeleteMessage.")
	params := &sqs.DeleteMessageInput{
		QueueUrl:      queueUrlOutput.QueueUrl,
		ReceiptHandle: aws.String(*msg.ReceiptHandle),
	}
	_, err := svc.DeleteMessage(params)

	if err != nil {
		return err
	}
	return nil
}

func InitSQS() {
	fmt.Println("Init start")
	// SQS のコネクション取得
	getSQSConnection()

	// SQSキュー取得
	// ローカルの起動が間に合わない場合に備え、ループ
	for {
		queueUrlOutput, err = getQueueUrl()
		if err == nil {
			break
		}
	}
}
