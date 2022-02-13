package s3

import (
	"bufio"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	la "github.com/takakosi/minio_elasticmq/aws"
)

const (
	endpoint   = "http://minio:9000"
	bucketName = "my-bucket-1"
)

var (
	svc *s3.S3
)

// 特定バケット配下のオブジェクト一覧
func GetObjects() error {

	// オブジェクト取得
	obj, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("test_data.csv"),
	})

	if err != nil {
		return err
	}

	// 1行ずつ読み込みメモリに全読みしない。
	scanner := bufio.NewScanner(obj.Body)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return nil
}

func getS3Connection() {
	sess := la.CreateSession()
	svc = s3.New(sess, &aws.Config{
		Endpoint:         aws.String(endpoint), // コンテナ内からアクセスする場合はホストをサービス名で指定
		S3ForcePathStyle: aws.Bool(true),       // ローカルで動かす場合は必須

	})
}

func InitS3() {
	fmt.Println("Init start")
	getS3Connection()
}
