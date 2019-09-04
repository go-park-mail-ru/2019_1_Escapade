package photo

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"github.com/aws/aws-sdk-go/aws"
	_ "github.com/aws/aws-sdk-go/aws/session"

	"io/ioutil"
	"os"
)

// AwsPublicConfig public aws information as region and endpoint
//easyjson:json
type AwsPublicConfig struct {
	set                   bool
	config                *aws.Config
	Region                string   `json:"region"`
	Endpoint              string   `json:"endpoint"`
	PlayersAvatarsStorage string   `json:"playersAvatarsStorage"`
	DefaultAvatar         string   `json:"defaultAvatar"`
	MaxFileSize           int64    `json:"maxFileSize"`
	AllowedFileTypes      []string `json:"allowedFileTypes"`
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

// MaxFileSize return the maximum size of the file
func MaxFileSize() int64 {
	if !_AWS.public.set {
		utils.Debug(true, "package photo not initialized")
		return 0
	}
	return _AWS.public.MaxFileSize
}

// AllowedFileTypes return allowed file types
func AllowedFileTypes() []string {
	if !_AWS.public.set {
		utils.Debug(true, "package photo not initialized")
		return nil
	}
	return _AWS.public.AllowedFileTypes
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
		Region: aws.String(publicAWS.Region)}
	if publicAWS.Endpoint != "" {
		publicAWS.config.Endpoint = aws.String(publicAWS.Endpoint)
	}

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

	//_AWS.public.config.Credentials = credentials.NewStaticCredentials(privateAWS.AccessKey, privateAWS.SecretKey, "TOKEN")

	os.Setenv(privateAWS.AccessURL, privateAWS.AccessKey)
	os.Setenv(privateAWS.SecretURL, privateAWS.SecretKey)
	return err
}

// GetDefaultAvatar return default avatar
func GetDefaultAvatar() string {
	return _AWS.public.DefaultAvatar
}
