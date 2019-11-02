package photo

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
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

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	// img, err := imaging.Decode(file)
	// if err != nil {
	// 	utils.Debug(false, "cant decode")
	// 	_, err = io.Copy(&buf, file)
	// } else {
	// 	err = imaging.Encode(&buf, img, imaging.JPEG)
	// }
	//utils.Debug(false, "buf:", string(buf.Bytes()))

	if err != nil {
		utils.Debug(false, "cant encode")
		return err
	}

	fileType := http.DetectContentType(buf.Bytes())
	fileSize := buf.Len()
	params := &s3.PutObjectInput{
		Bucket: aws.String(_AWS.public.PlayersAvatarsStorage),
		Key:    aws.String(key),
		Body:   bytes.NewReader(buf.Bytes()),
		ACL:    aws.String("public-read"),

		ContentLength: aws.Int64(int64(fileSize)),
		ContentType:   aws.String(fileType),
	}

	sess := session.Must(session.NewSession(_AWS.public.config))
	svc := s3.New(sess)
	resp, err := svc.PutObject(params)
	if err != nil {
		return err
	}

	fmt.Println("Done", resp)
	return
}

/*
func SaveImageInS3RAW(key string, buf *bufio.Reader) (err error) {
	if !_AWS.public.set {
		utils.Debug(true, "package photo not initialized")
		return
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	// img, err := imaging.Decode(file)
	// if err != nil {
	// 	utils.Debug(false, "cant decode")
	// 	_, err = io.Copy(&buf, file)
	// } else {
	// 	err = imaging.Encode(&buf, img, imaging.JPEG)
	// }
	//utils.Debug(false, "buf:", string(buf.Bytes()))

	if err != nil {
		utils.Debug(false, "cant encode")
		return err
	}

	fileType := http.DetectContentType(buf.Bytes())
	fileSize := buf.Len()
	params := &s3.PutObjectInput{
		Bucket: aws.String(_AWS.public.PlayersAvatarsStorage),
		Key:    aws.String(key),
		Body:   bytes.NewReader(buf.Bytes()),
		ACL:    aws.String("public-read"),

		ContentLength: aws.Int64(int64(fileSize)),
		ContentType:   aws.String(fileType),
	}

	sess := session.Must(session.NewSession(_AWS.public.config))
	svc := s3.New(sess)
	resp, err := svc.PutObject(params)
	if err != nil {
		return err
	}

	fmt.Println("Done", resp)
	return
}*/

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
	if key == "1.png" {
		key = "artyom/1.png"
	}
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
