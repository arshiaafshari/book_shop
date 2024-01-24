package initializers

import (
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func ConnectToMinio() {

	useSSL := false

	//initialize minio client
	var err error
	MinioClient, err = minio.New(os.Getenv("endpoint_MINIO"), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("accessKey_MINIO"), os.Getenv("secretKey_MINIO"), ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

}
