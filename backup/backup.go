package backup

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/mholt/archiver/v3"
	"log"
	"os"
	"strings"
	"time"
)

type Backup struct {
	Source   string
	Dest     string
	File     string
	S3Bucket string
	S3Region string
}

func NewBackup(source, s3bucket, s3region string) (*Backup, error) {
	b := new(Backup)
	b.S3Region = s3region
	b.S3Bucket = s3bucket
	b.Source = source

	now := time.Now()
	timeStr := now.Format(time.RFC3339)
	if _, err := os.Stat(b.Source); err != nil {
		return nil, err
	}
	f := strings.Split(b.Source, "/")
	b.File = f[len(f)-1] + "_" + timeStr + ".zip"
	b.Dest = "/tmp/" + b.File

	return b, nil
}

// compresses source
func (b *Backup) Compress() error {
	log.Println(fmt.Sprintf("source: %s", b.Source))
	log.Println(fmt.Sprintf("destination: %s", b.Dest))

	if err := archiver.Archive([]string{b.Source}, b.Dest); err != nil {
		return err
	}

	return nil
}

// send the file to our aws bucket
func (b *Backup) Ship() error {

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(b.S3Region)},
	))

	uploader := s3manager.NewUploader(sess)

	uploadFile, err := os.Open(b.Dest)
	if err != nil {
		return err
	}

	// Upload input parameters
	upParams := &s3manager.UploadInput{
		Bucket: &b.S3Bucket,
		Key:    &b.File,
		Body:   uploadFile,
	}

	// Perform an upload.
	result, err := uploader.Upload(upParams)
	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("stored: %s", result.Location))

	return nil
}

// cleans up and deletes the file
func (b *Backup) Cleanup() error {
	log.Println(fmt.Sprintf("removing: %s", b.Dest))
	if _, err := os.Stat(b.Dest); err == nil {
		if err = os.Remove(b.Dest); err != nil {
			return err
		}
	}
	return nil
}
