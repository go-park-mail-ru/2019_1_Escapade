package photo

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"io/ioutil"
	"mime/multipart"
	"os"
	"time"
)

// AwsPublicConfig public aws information as region and endpoint
//easyjson:json
type AwsPublicConfig struct {
	set                   bool
	config                *aws.Config
	Region                string `json:"region"`
	Endpoint              string `json:"endpoint"`
	PlayersAvatarsStorage string `json:"playersAvatarsStorage"`
	DefaultAvatar         string `json:"defaultAvatar"`
}

// AwsPrivateConfig private aws information. Need another json.
//easyjson:json
type AwsPrivateConfig struct {
	set       bool
	AccessURL string `json:"accessUrl"`
	AccessKey string `json:"accessKey"`
	SecretURL string `json:"secretUrl"`
	SecretKey string `json:"secretKey"`
}

// _AWS singleton
var _AWS struct {
	public  AwsPublicConfig
	private AwsPrivateConfig
}

//GetImages get image from image storage and set it to every user
func GetImages(users ...*models.UserPublicInfo) {
	if !_AWS.public.set {
		utils.Debug(true, "package photo not initialized")
		return
	}

	for _, user := range users {
		if user == nil {
			continue
		}
		if user.FileKey == "" {
			continue
		}
		var err error
		if user.PhotoURL, err = GetImage(user.FileKey); err != nil {
			continue
		}
	}
}

//SaveImage save image given by 'key' user
func SaveImage(key string, file multipart.File) (err error) {
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

//GetImage get image by its key
func GetImage(key string) (url string, err error) {
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
func DeleteImage(key string) (err error) {
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

// Init load photo configuration
// publicConfigPath - path to json, initialize public AWS information
// privateConfigPath - path to json, initialize private AWS information
func Init(publicConfigPath string, privateConfigPath string) error {

	if err := initPublic(publicConfigPath); err != nil {
		utils.Debug(false, "Cant load AwsPublicConfig. Error message:",
			err.Error())
		return err
	}

	if err := initPrivate(privateConfigPath); err != nil {
		utils.Debug(false, "Cant load AwsPrivateConfig. Error message:",
			err.Error())
		return err
	}
	return nil
}

// initPublic initialize public AWS information
func initPublic(publicConfigPath string) error {
	var (
		data []byte
		err  error
	)

	if data, err = ioutil.ReadFile(publicConfigPath); err != nil {
		return err
	}
	var publicAWS = &AwsPublicConfig{}
	if err = publicAWS.UnmarshalJSON(data); err != nil {
		return err
	}
	publicAWS.config = &aws.Config{
		Region:   aws.String(publicAWS.Region),
		Endpoint: aws.String(publicAWS.Endpoint)}
	publicAWS.set = true

	_AWS.public = *publicAWS
	return err
}

// initPrivate initialize private AWS information
func initPrivate(path string) error {
	var (
		data []byte
		err  error
	)

	if data, err = ioutil.ReadFile(path); err != nil {
		return err
	}
	var privateAWS = &AwsPrivateConfig{}
	if err = privateAWS.UnmarshalJSON(data); err != nil {
		return err
	}

	privateAWS.set = true

	_AWS.private = *privateAWS

	os.Setenv(privateAWS.AccessURL, privateAWS.AccessKey)
	os.Setenv(privateAWS.SecretURL, privateAWS.SecretKey)
	return err
}

func GetDefaultAvatar() string {
	return _AWS.public.DefaultAvatar
}
