package minio

import (
	"go-social-network/configs"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client

func Default() *minio.Client {
	return minioClient
}

func InitMinio() {
	var err error
	conf := configs.Data().Minio
	minioClient, err = minio.New(conf.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.AccessKey, conf.AccessSecret, ""),
		Secure: conf.UseSSL,
	})
	if err != nil {
		log.Fatal("minio client create fail, err %s", err)
	}
}
