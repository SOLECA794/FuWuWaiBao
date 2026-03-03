package oss

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"smart-teaching-backend/pkg/config"
)

type MinioClient struct {
	client   *minio.Client
	bucket   string
	endpoint string
	useSSL   bool
}

func NewMinioClient(cfg *config.OSSConfig) (*MinioClient, error) {
	// 初始化MinIO客户端
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("初始化MinIO失败: %w", err)
	}

	// 检查bucket是否存在，不存在则创建
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("检查bucket失败: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("创建bucket失败: %w", err)
		}

		// 设置bucket策略为公开读
		policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::` + cfg.Bucket + `/*"]}]}`
		client.SetBucketPolicy(ctx, cfg.Bucket, policy)
	}

	return &MinioClient{
		client:   client,
		bucket:   cfg.Bucket,
		endpoint: cfg.Endpoint,
		useSSL:   cfg.UseSSL,
	}, nil
}

// UploadFile 上传文件
func (m *MinioClient) UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	// 上传文件
	_, err := m.client.PutObject(ctx, m.bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("上传文件失败: %w", err)
	}

	// 生成访问URL（7天有效期）
	url, err := m.client.PresignedGetObject(ctx, m.bucket, objectName, 7*24*time.Hour, nil)
	if err != nil {
		// 如果生成预签名URL失败，返回一个基本路径
		scheme := "http"
		if m.useSSL {
			scheme = "https"
		}
		urlStr := fmt.Sprintf("%s://%s/%s/%s", scheme, m.endpoint, m.bucket, objectName)
		return urlStr, nil
	}

	return url.String(), nil
}

// GetFileURL 获取文件访问URL
func (m *MinioClient) GetFileURL(ctx context.Context, objectName string, expires time.Duration) (string, error) {
	url, err := m.client.PresignedGetObject(ctx, m.bucket, objectName, expires, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

// DeleteFile 删除文件
func (m *MinioClient) DeleteFile(ctx context.Context, objectName string) error {
	return m.client.RemoveObject(ctx, m.bucket, objectName, minio.RemoveObjectOptions{})
}
