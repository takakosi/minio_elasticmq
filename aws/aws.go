package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// セッションを返す
func CreateSession() *session.Session {
	// 特に設定しなくても環境変数にセットしたクレデンシャル情報を利用して接続してくれる
	cfg := aws.Config{
		Region: aws.String("ap-northeast-1"),
	}
	return session.Must(session.NewSession(&cfg))
}
