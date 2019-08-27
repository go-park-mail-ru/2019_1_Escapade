package photo

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"mime/multipart"
	"time"
)

//GetImages get image from image storage and set it to every user
func GetImages(users ...*models.UserPublicInfo) {
	if !_AWS.public.set {
		utils.Debug(true, "package photo not initialized")
		return
	}

	for _, user := range users {
		if user == nil {
			utils.Debug(false, "image warning: user == nil")
			continue
		}
		if user.FileKey == "" {
			utils.Debug(false, "image warning: FileKey == ''")
			continue
		}
		var err error
		if user.PhotoURL, err = GetImageFromS3(user.FileKey); err != nil {
			utils.Debug(false, "image error: ", err.Error())
			continue
		}
	}
}

//SaveImageInS3 save image given by 'key' user
func SaveImageInS3(key string, file multipart.File) (err error) {
	if !_AWS.public.set {
		utils.Debug(true, "package photo not initialized")
		return
	}

	sess := session.Must(session.NewSession(_AWS.public.config))

	// Create S3 service client
	svc := s3.New(sess)

	//snippet-start:[s3.go.put_object.call]
	_, err = svc.PutObject((&s3.PutObjectInput{}).
		SetBucket(_AWS.public.PlayersAvatarsStorage).
		SetKey(key).
		SetBody(file),
	)

	return
}

//GetImageFromS3 get image by its key
func GetImageFromS3(key string) (url string, err error) {
	if !_AWS.public.set {
		utils.Debug(true, "package photo not initialized")
		return
	}

	sess, err := session.NewSession(_AWS.public.config)
	if err != nil {
		return
	}
	svc := s3.New(sess)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(_AWS.public.PlayersAvatarsStorage),
		Key:    aws.String(key),
	})
	url, err = req.Presign(24 * time.Hour)
	return
}

//DeleteImage delete image, which key is the same as given
func DeleteImageFromS3(key string) (err error) {
	if !_AWS.public.set {
		utils.Debug(true, "package photo not initialized")
		return
	}

	sess, err := session.NewSession(_AWS.public.config)
	if err != nil {
		return
	}
	svc := s3.New(sess)

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(_AWS.public.PlayersAvatarsStorage),
		Key:    aws.String(key)})
	if err != nil {
		utils.Debug(false, "Unable to delete object",
			key, "Error message:", err)
		return
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(_AWS.public.PlayersAvatarsStorage),
		Key:    aws.String(key),
	})
	if err != nil {
		utils.Debug(false, "Error occurred while waiting for object",
			key, "to be deleted. Error message:", err)
		return
	}

	return
}
