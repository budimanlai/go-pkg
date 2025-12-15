package main

import (
	"fmt"

	"github.com/budimanlai/go-pkg/storage"
)

func main() {
	fileStorage()
}

func fileStorage() {
	basePath := "./example"
	// localStorage := storage.NewLocalStorage(basePath+"/uploads", "http://localhost:8080/uploads")
	// fileStorage := storage.NewStorage(localStorage)

	config := storage.S3Config{
		Region:          "us-east-1",
		Bucket:          "public",
		AccessKeyID:     "admin",
		SecretAccessKey: "admin123",
		EndpointURL:     "http://localhost:8333",
		PublicURL:       "http://localhost:8888/buckets/public",
	}
	s3Storage := storage.NewS3Storage(config)
	fileStorage := storage.NewStorage(s3Storage)

	dest := "image1.png"
	source := basePath + "/data/" + dest

	err := fileStorage.Save(source, dest)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	fmt.Println("File uploaded")

	url, err := fileStorage.GetURL(dest)
	if err != nil {
		fmt.Printf("%s", err.Error())
	} else {
		fmt.Printf("File URL: %s\n", url)
	}

	urlSigned, err := fileStorage.GetSignedURL(dest, 60)
	if err != nil {
		fmt.Printf("%s", err.Error())
	} else {
		fmt.Printf("Signed File URL: %s\n", urlSigned)
	}

	exists, err := fileStorage.Exists(dest)
	if err != nil {
		fmt.Printf("%s", err.Error())
	} else {
		fmt.Printf("File exists: %t\n", exists)
	}

	// err = fileStorage.Delete("custom_bucket/tabel_1.png")
	// if err != nil {
	// 	fmt.Printf("%s", err.Error())
	// } else {
	// 	fmt.Println("File deleted")
	// }
}
