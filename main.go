package main

import (
	"agrippanux/gobackup/backup"
	"flag"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	log.Println("gobackup started")

	if err := godotenv.Load(); err != nil {
		log.Println("did not find a .env file, proceeding")
	}

	if os.Getenv("AWS_ACCESS_KEY_ID") == "" ||
		os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		log.Fatalln("please set the AWS_SECRET_KEY_ID and AWS_SECRET_ACCESS_KEY env variables")
	}

	source := flag.String("source", "", "source to zip")
	s3Bucket := flag.String("s3bucket", "", "s3 bucket to use")
	s3Region := flag.String("s3region", "", "s3 region to use (example, us-west-2)")

	flag.Parse()

	if len(*source) == 0 ||
		len(*s3Bucket) == 0 ||
		len(*s3Region) == 0 {
		log.Fatalln("flags source, s3bucket, and s3region required")
	}

	b, err := backup.NewBackup(*source, *s3Bucket, *s3Region)
	if err != nil {
		log.Fatalln(err)
	}

	// start the backup process
	if err := b.Compress(); err != nil {
		log.Fatalln(err)
	}

	// clean up after ourselves
	defer func() {
		if err := b.Cleanup(); err != nil {
			log.Fatalln(err)
		}
	}()

	// show to AWS
	if err := b.Ship(); err != nil {
		log.Fatalln(err)
	}

}
