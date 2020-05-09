package main

import (
  "agrippanux/gobackup/backup"
  "github.com/joho/godotenv"
  "log"
  "os"
)

func main() {
  log.Println("gobackup started")

  godotenv.Load()

  source := os.Getenv("SOURCE")
  s3Bucket := os.Getenv("S3_BUCKET")
  s3Region := os.Getenv("S3_REGION")

  if source == "" ||
    s3Bucket == "" ||
    s3Region == "" {
    log.Fatalln("missing either SOURCE, S3_BUCKET, S3_REGION env variable")
  }

  backup, err := backup.NewBackup(source, s3Bucket, s3Region)
  if err != nil {
    log.Fatalln(err)
  }

  // start the backup process
  if err := backup.Compress(); err != nil {
    log.Fatalln(err)
  }

  // clean up after ourselves
  defer backup.Cleanup()

  // show to AWS
  if err := backup.Ship(); err != nil {
    log.Fatalln(err)
  }

}
