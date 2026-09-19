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

// package utils

// import (
// 	"bytes" // Tambahan baru
// 	"context"
// 	"fmt"
// 	"io"    // Tambahan baru
// 	"mime/multipart"
// 	"os"

// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/aws/aws-sdk-go-v2/config"
// 	"github.com/aws/aws-sdk-go-v2/credentials"
// 	"github.com/aws/aws-sdk-go-v2/service/s3"
// 	"github.com/google/uuid"
// )

// func UploadToCleverCloud(file *multipart.FileHeader) (string, error) {
// 	src, err := file.Open()
// 	if err != nil {
// 		return "", err
// 	}
// 	defer src.Close()

// 	// 1. SOLUSI: Salin isi stream file ke dalam memory buffer
// 	// Ini menjamin AWS SDK tidak kehilangan posisi kursor saat menghitung SHA256
// 	buf := bytes.NewBuffer(nil)
// 	if _, err := io.Copy(buf, src); err != nil {
// 		return "", fmt.Errorf("gagal menyalin file ke memory: %v", err)
// 	}
// 	// Buat reader baru dari buffer memori yang sudah utuh
// 	bodyReader := bytes.NewReader(buf.Bytes())

// 	accessKey := os.Getenv("S3_KEY")
// 	secretKey := os.Getenv("S3_SECRET")
// 	endpoint := os.Getenv("S3_ENDPOINT")
// 	bucketName := os.Getenv("S3_BUCKET")

// 	if accessKey == "" || secretKey == "" || endpoint == "" || bucketName == "" {
// 		return "", fmt.Errorf("kredensial S3 belum lengkap di Vercel")
// 	}

// 	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
// 		return aws.Endpoint{
// 			URL:           "https://" + endpoint,
// 			SigningRegion: "us-east-1",
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

// 	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
// 		o.UsePathStyle = true
// 	})

// 	filename := uuid.New().String() + "-" + file.Filename

// 	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
// 		Bucket: aws.String(bucketName),
// 		Key:    aws.String(filename),
// 		Body:   bodyReader, // 2. UBAH DI SINI: Gunakan bodyReader, bukan src
// 		ACL:    "public-read",
// 	})
// 	if err != nil {
// 		return "", err
// 	}

// 	imageURL := fmt.Sprintf("https://%s/%s/%s", endpoint, bucketName, filename)
// 	return imageURL, nil
// }

// package utils

// import (
// 	"bytes"
// 	"context"
// 	"fmt"
// 	"io"
// 	"mime/multipart"
// 	"net/http" // Tambahan baru untuk mendeteksi tipe konten (MIME Type)
// 	"os"

// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/aws/aws-sdk-go-v2/config"
// 	"github.com/aws/aws-sdk-go-v2/credentials"
// 	"github.com/aws/aws-sdk-go-v2/service/s3"
// 	"github.com/google/uuid"
// )

// func UploadToCleverCloud(file *multipart.FileHeader) (string, error) {
// 	src, err := file.Open()
// 	if err != nil {
// 		return "", err
// 	}
// 	defer src.Close()

// 	// 1. Salin isi stream file ke dalam memory buffer
// 	buf := bytes.NewBuffer(nil)
// 	if _, err := io.Copy(buf, src); err != nil {
// 		return "", fmt.Errorf("gagal menyalin file ke memory: %v", err)
// 	}

// 	// 2. Kalkulasi atribut file secara pasti
// 	fileBytes := buf.Bytes()
// 	bodyReader := bytes.NewReader(fileBytes)
// 	contentLength := int64(len(fileBytes))
// 	contentType := http.DetectContentType(fileBytes) // Otomatis mendeteksi image/jpeg, image/png, dll.

// 	accessKey := os.Getenv("S3_KEY")
// 	secretKey := os.Getenv("S3_SECRET")
// 	endpoint := os.Getenv("S3_ENDPOINT")
// 	bucketName := os.Getenv("S3_BUCKET")

// 	if accessKey == "" || secretKey == "" || endpoint == "" || bucketName == "" {
// 		return "", fmt.Errorf("kredensial S3 belum lengkap di Vercel")
// 	}

// 	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
// 		return aws.Endpoint{
// 			URL:           "https://" + endpoint,
// 			SigningRegion: "us-east-1",
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

// 	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
// 		o.UsePathStyle = true
// 	})

// 	filename := uuid.New().String() + "-" + file.Filename

// 	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
// 		Bucket:        aws.String(bucketName),
// 		Key:           aws.String(filename),
// 		Body:          bodyReader,
// 		ACL:           "public-read",
// 		ContentLength: aws.Int64(contentLength), // Mencegah AWS SDK menggunakan aws-chunked
// 		ContentType:   aws.String(contentType),  // Menghindari penolakan format MIME dari server S3
// 	})
// 	if err != nil {
// 		return "", err
// 	}

// 	imageURL := fmt.Sprintf("https://%s/%s/%s", endpoint, bucketName, filename)
// 	return imageURL, nil
// }

package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types" // Tambahan package types
	"github.com/google/uuid"
)

func UploadToCleverCloud(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// 1. Salin ke memori
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, src); err != nil {
		return "", fmt.Errorf("gagal menyalin file ke memory: %v", err)
	}

	fileBytes := buf.Bytes()
	contentType := http.DetectContentType(fileBytes)

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

	// =====================================================================
	// SOLUSI: BYPASS FUNGSI UPLOAD AWS SDK MENGGUNAKAN PRESIGNED URL
	// =====================================================================
	
	// 2. Buat URL "Tiket Masuk" yang sah
	presignClient := s3.NewPresignClient(client)
	presignReq, err := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(filename),
		ContentType: aws.String(contentType),
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", fmt.Errorf("gagal membuat tiket URL S3: %v", err)
	}

	// 3. Eksekusi unggahan menggunakan HTTP Client bawaan Go yang kebal bug SHA256
	req, err := http.NewRequest("PUT", presignReq.URL, bytes.NewReader(fileBytes))
	if err != nil {
		return "", err
	}

	// Wajib menyertakan header yang sama persis dengan yang didaftarkan di tiket URL
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-amz-acl", string(types.ObjectCannedACLPublicRead))
	req.ContentLength = int64(len(fileBytes))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("gagal menembakkan HTTP PUT: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ditolak oleh Clever Cloud S3 (Status %d): %s", resp.StatusCode, string(respBody))
	}

	// 4. Kembalikan URL Gambar
	imageURL := fmt.Sprintf("https://%s/%s/%s", endpoint, bucketName, filename)
	return imageURL, nil
}