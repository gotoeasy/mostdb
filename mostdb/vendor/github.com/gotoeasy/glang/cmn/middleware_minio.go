package cmn

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Minio客户端结构体
type Minio struct {
	endpoint string
	username string
	password string
	bucket   string
}

func NewMinio(endpoint string, username string, password string, bucket string) *Minio {
	return &Minio{
		endpoint: endpoint,
		username: username,
		password: password,
		bucket:   bucket,
	}
}

// 上传文件
func (m *Minio) Upload(localPathFile string, minioObjectName string) error {
	ctx := context.Background()

	// 初始化
	minioClient, err := minio.New(m.endpoint, &minio.Options{Creds: credentials.NewStaticV4(m.username, m.password, ""), Secure: false})
	if err != nil {
		return err
	}

	// 检查创建Bucket
	exists, err := minioClient.BucketExists(ctx, m.bucket)
	if !(err == nil && exists) {
		err = minioClient.MakeBucket(ctx, m.bucket, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}

	// 上传
	_, err = minioClient.FPutObject(ctx, m.bucket, minioObjectName, localPathFile, minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

// 下载文件
func (m *Minio) Download(minioObjectName string, localPathFile string) error {
	ctx := context.Background()

	// 初始化
	minioClient, err := minio.New(m.endpoint, &minio.Options{Creds: credentials.NewStaticV4(m.username, m.password, ""), Secure: false})
	if err != nil {
		return err
	}

	// 下载
	err = minioClient.FGetObject(ctx, m.bucket, minioObjectName, localPathFile, minio.GetObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}
