package aws

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/configuration"
)

type AWSService struct {
	config *aws.Config
	photo  configuration.Photo

	log infrastructure.Logger
}

func New(
	rep configuration.PhotoRepository,
	logger infrastructure.Logger,
) (*AWSService, error) {
	// check configuration repository given
	if rep == nil {
		return nil, errors.New(ErrNoConfiguration)
	}

	//overriding the nil value of Logger
	if logger == nil {
		logger = new(infrastructure.LoggerNil)
	}

	conf := rep.Get()
	var service = &AWSService{
		photo: conf,
		config: &aws.Config{
			Region: aws.String(conf.Region),
		},
		log: logger,
	}

	if service.photo.Endpoint != "" { // TODO что за костыль?
		service.config.Endpoint = aws.String(service.photo.Endpoint)
	}

	return service, nil
}

// MaxFileSize return the maximum size of the file
func (service *AWSService) MaxFileSize() int64 {
	return service.photo.MaxFileSize
}

// AllowedFileTypes return allowed file types
func (service *AWSService) AllowedFileTypes() []string {
	return service.photo.AllowedFileTypes
}

// GetDefaultAvatar return default avatar
func (service *AWSService) GetDefaultAvatar() string {
	return service.photo.DefaultAvatar
}

//GetImages get image from image storage and set it to every user
func (service *AWSService) GetImages(
	users ...*models.UserPublicInfo,
) {
	for _, user := range users {
		if user == nil {
			service.log.Println("image warning: user == nil")
			continue
		}
		if user.FileKey == "" {
			service.log.Println("image warning: FileKey == ''")
			continue
		}
		photoURL, err := service.GetImage(user.FileKey)
		if err != nil {
			service.log.Println("image error: ", err.Error())
			continue
		}
		user.PhotoURL = photoURL
	}
}

//SaveImage save image given by 'key' user
func (service *AWSService) SaveImage(
	key string,
	file multipart.File,
) error {
	if key == "" {
		return errors.New(ErrNoKey)
	}

	var buf bytes.Buffer
	_, err := io.Copy(&buf, file)
	// img, err := imaging.Decode(file)
	// if err != nil {
	// 	utils.Debug(false, "cant decode")
	// 	_, err = io.Copy(&buf, file)
	// } else {
	// 	err = imaging.Encode(&buf, img, imaging.JPEG)
	// }
	//utils.Debug(false, "buf:", string(buf.Bytes()))

	if err != nil {
		service.log.Println("cant encode")
		return err
	}

	fileType := http.DetectContentType(buf.Bytes())
	fileSize := buf.Len()
	params := &s3.PutObjectInput{
		Bucket: aws.String(
			service.photo.PlayersAvatarsStorage,
		),
		Key:  aws.String(key),
		Body: bytes.NewReader(buf.Bytes()),
		ACL:  aws.String(ACLDefault),

		ContentLength: aws.Int64(int64(fileSize)),
		ContentType:   aws.String(fileType),
	}

	sess := session.Must(session.NewSession(service.config))
	svc := s3.New(sess)
	resp, err := svc.PutObject(params)
	if err != nil {
		return err
	}

	service.log.Println("Done", resp)
	return nil
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

//GetImage get image by its key
func (service *AWSService) GetImage(key string) (string, error) {
	if key == "" {
		return "", errors.New(ErrNoKey)
	}

	sess, err := session.NewSession(service.config)
	if err != nil {
		return "", err
	}
	svc := s3.New(sess)
	if key == "1.png" { // TODO  убрать костыль
		key = "artyom/1.png"
	}
	req, _ := svc.GetObjectRequest(
		&s3.GetObjectInput{
			Bucket: aws.String(
				service.photo.PlayersAvatarsStorage,
			),
			Key: aws.String(key),
		},
	)
	return req.Presign(service.photo.Expire)
}

//DeleteImage delete image, which key is the same as given
func (service *AWSService) DeleteImage(key string) error {
	if key == "" {
		return errors.New(ErrNoKey)
	}

	sess, err := session.NewSession(service.config)
	if err != nil {
		return err
	}
	svc := s3.New(sess)

	_, err = svc.DeleteObject(
		&s3.DeleteObjectInput{
			Bucket: aws.String(
				service.photo.PlayersAvatarsStorage,
			),
			Key: aws.String(key),
		},
	)
	if err != nil {
		service.log.Println(
			"Unable to delete object", key,
			"Error message:", err,
		)
		return err
	}

	err = svc.WaitUntilObjectNotExists(
		&s3.HeadObjectInput{
			Bucket: aws.String(
				service.photo.PlayersAvatarsStorage,
			),
			Key: aws.String(key),
		},
	)
	if err != nil {
		service.log.Println(
			"Error occurred while waiting for object", key,
			"to be deleted. Error message:", err,
		)
		return err
	}

	return nil
}
