package api

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"

	"fmt"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (h *Handler) Setfiles(users ...*models.UserPublicInfo) (err error) {

	for _, user := range users {
		if user == nil {
			continue
		}
		if user.FileKey == "" {
			continue
			//return re.ErrorAvatarNotFound()
		}
		if user.PhotoURL, err = h.getURLToAvatar(user.FileKey); err != nil {
			continue
			//return re.ErrorAvatarNotFound()
		}
	}
	return nil
}

func (h *Handler) saveFile(key string, file multipart.File) (err error) {

	sess := session.Must(session.NewSession(h.AWS.AwsConfig))

	// Create S3 service client
	svc := s3.New(sess)

	//snippet-start:[s3.go.put_object.call]
	_, err = svc.PutObject((&s3.PutObjectInput{}).
		SetBucket(h.Storage.PlayersAvatarsStorage).
		SetKey(key).
		SetBody(file),
	)

	return
}

func (h *Handler) getURLToAvatar(key string) (url string, err error) {

	sess, err := session.NewSession(h.AWS.AwsConfig)
	if err != nil {
		return
	}
	svc := s3.New(sess)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(h.Storage.PlayersAvatarsStorage),
		Key:    aws.String(key),
	})
	url, err = req.Presign(24 * time.Hour)
	return
}

func (h *Handler) deleteAvatar(key string) (err error) {
	sess, err := session.NewSession(h.AWS.AwsConfig)
	if err != nil {
		return
	}
	svc := s3.New(sess)

	// Delete the item
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(h.Storage.PlayersAvatarsStorage),
		Key:    aws.String(key)})
	if err != nil {
		fmt.Printf("Unable to delete object %q from bucket %q, %v\n", key, h.Storage.PlayersAvatarsStorage, err)
		return
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(h.Storage.PlayersAvatarsStorage),
		Key:    aws.String(key),
	})
	if err != nil {
		fmt.Printf("Error occurred while waiting for object %q to be deleted\n", key)
		return
	}

	fmt.Printf("Object %q successfully deleted\n", key)
	return
}
