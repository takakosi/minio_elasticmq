package main

import (
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/labstack/echo"
)

// レスポンス用
type Response struct {
	Objects []ObjInfo `json:"objects"`
}

// レスポンス用
type ObjInfo struct {
	Key          string    `json:"key"`
	LastModified time.Time `json:"lastModified"`
	Size         int64     `json:"size"`
}

func main() {
	for {
		if err := sqs.RetrieveMessage(); err != nil {
			log.Fatal(err)
		}
	}
}

// 特定バケット配下のオブジェクト一覧
func getObjects(c echo.Context) error {
	sess := createSession()
	svc := s3.New(sess)
	bucket := c.QueryParam("bucket")

	// オブジェクト取得
	res, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	})
	// 本来はもっときちんとエラーハンドリングした方が良いが、簡単のため今回はこれで良しとする
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	var objects []ObjInfo
	// オブジェクト情報をJSONにつめて返す
	for _, content := range res.Contents {
		objects = append(
			objects,
			ObjInfo{Key: *content.Key, LastModified: *content.LastModified, Size: *content.Size},
		)
	}
	return c.JSON(http.StatusOK, objects)
}

// セッションを返す
func createSession() *session.Session {
	// 特に設定しなくても環境変数にセットしたクレデンシャル情報を利用して接続してくれる
	cfg := aws.Config{
		Region:           aws.String("ap-northeast-1"),
		Endpoint:         aws.String("http://minio:9000"), // コンテナ内からアクセスする場合はホストをサービス名で指定
		S3ForcePathStyle: aws.Bool(true),                  // ローカルで動かす場合は必須
	}
	return session.Must(session.NewSession(&cfg))
}
