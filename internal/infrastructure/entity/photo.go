package entity

import "github.com/go-park-mail-ru/2019_1_Escapade/internal/domens"

// AwsPublicConfig public aws information as region and endpoint
//easyjson:json
type PhotoPublicConfig struct {
	Region                string          `json:"region"`
	Endpoint              string          `json:"endpoint"`
	PlayersAvatarsStorage string          `json:"playersAvatarsStorage"`
	DefaultAvatar         string          `json:"defaultAvatar"`
	MaxFileSize           int64           `json:"maxFileSize"`
	Expire                domens.Duration `json:"expire"`
	AllowedFileTypes      []string        `json:"allowedFileTypes"`
}

// AwsPrivateConfig private aws information. Need another json.
//easyjson:json
type PhotoPrivateConfig struct {
	AccessURL string `json:"accessUrl"`
	AccessKey string `json:"accessKey"`
	SecretURL string `json:"secretUrl"`
	SecretKey string `json:"secretKey"`
}

// TODO deleteme

// MarshalJSON supports json.Marshaler interface
func (v PhotoPublicConfig) MarshalJSON() ([]byte, error) {
	return []byte{}, nil
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PhotoPublicConfig) UnmarshalJSON(data []byte) error {
	return nil
}

// MarshalJSON supports json.Marshaler interface
func (v PhotoPrivateConfig) MarshalJSON() ([]byte, error) {
	return []byte{}, nil
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PhotoPrivateConfig) UnmarshalJSON(data []byte) error {
	return nil
}
