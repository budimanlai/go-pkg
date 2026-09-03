package main

import (
	"fmt"

	"github.com/budimanlai/go-pkg/v3/storage"
)

func main() {
	fileStorage()
}

func fileStorage() {
	basePath := "./example"

	config := storage.S3Config{
		Region:          "your_bucket_region",
		Bucket:          "your_bucket_name",
		AccessKeyID:     "your_access_key_id",
		SecretAccessKey: "your_secret_access_key",
		EndpointURL:     "your_endpoint_url",
	}
	s3Storage := storage.NewS3Storage(config)
	fileStorage := storage.NewStorage(s3Storage)

	dest := "public/avatar/image1.png"
	source := basePath + "/data/image1.png"

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
}
