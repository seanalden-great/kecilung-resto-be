// package utils

// import (
// 	"context"
// 	"fmt"
// 	"mime/multipart"
// 	"os"

// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/aws/aws-sdk-go-v2/config"
// 	"github.com/aws/aws-sdk-go-v2/credentials"
// 	"github.com/aws/aws-sdk-go-v2/service/s3"
// 	"github.com/google/uuid"
// )

// func UploadToCleverCloud(file *multipart.FileHeader) (string, error) {
// 	// Buka file yang diupload
// 	src, err := file.Open()
// 	if err != nil {
// 		return "", err
// 	}
// 	defer src.Close()

// 	// Ambil kredensial dari Environment Variables Vercel/.env
// 	accessKey := os.Getenv("S3_KEY")
// 	secretKey := os.Getenv("S3_SECRET")
// 	endpoint := os.Getenv("S3_ENDPOINT") // Contoh: cellar-c2.services.clever-cloud.com
// 	bucketName := os.Getenv("S3_BUCKET")

// 	// Konfigurasi AWS S3 khusus untuk Clever Cloud Cellar
// 	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
// 		return aws.Endpoint{
// 			URL:           "https://" + endpoint,
// 			SigningRegion: "us-east-1", // Clever cloud biasanya menggunakan region default ini untuk S3
// 		}, nil
// 	})

// 	cfg, err := config.LoadDefaultConfig(context.TODO(),
// 		config.WithRegion("us-east-1"),
// 		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
// 		config.WithEndpointResolverWithOptions(customResolver),
// 	)
// 	if err != nil {
// 		return "", err
// 	}

// 	client := s3.NewFromConfig(cfg)

// 	// Buat nama file unik agar tidak bentrok
// 	filename := uuid.New().String() + "-" + file.Filename

// 	// Unggah ke Clever Cloud
// 	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
// 		Bucket: aws.String(bucketName),
// 		Key:    aws.String(filename),
// 		Body:   src,
// 		ACL:    "public-read", // Agar gambar bisa diakses publik
// 	})
// 	if err != nil {
// 		return "", err
// 	}

// 	// Kembalikan URL Publik Gambar
// 	imageURL := fmt.Sprintf("https://%s.%s/%s", bucketName, endpoint, filename)
// 	return imageURL, nil
// }

package utils

import (
	"bytes" // Tambahan baru
	"context"
	"fmt"
	"io"    // Tambahan baru
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

func UploadToCleverCloud(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// 1. SOLUSI: Salin isi stream file ke dalam memory buffer
	// Ini menjamin AWS SDK tidak kehilangan posisi kursor saat menghitung SHA256
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, src); err != nil {
		return "", fmt.Errorf("gagal menyalin file ke memory: %v", err)
	}
	// Buat reader baru dari buffer memori yang sudah utuh
	bodyReader := bytes.NewReader(buf.Bytes())

	accessKey := os.Getenv("S3_KEY")
	secretKey := os.Getenv("S3_SECRET")
	endpoint := os.Getenv("S3_ENDPOINT")
	bucketName := os.Getenv("S3_BUCKET")

	if accessKey == "" || secretKey == "" || endpoint == "" || bucketName == "" {
		return "", fmt.Errorf("kredensial S3 belum lengkap di Vercel")
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           "https://" + endpoint,
			SigningRegion: "us-east-1",
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return "", err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	filename := uuid.New().String() + "-" + file.Filename

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
		Body:   bodyReader, // 2. UBAH DI SINI: Gunakan bodyReader, bukan src
		ACL:    "public-read",
	})
	if err != nil {
		return "", err
	}

	imageURL := fmt.Sprintf("https://%s/%s/%s", endpoint, bucketName, filename)
	return imageURL, nil
}