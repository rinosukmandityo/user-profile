package serializer

import (
	"github.com/rinosukmandityo/user-profile/helper"
	m "github.com/rinosukmandityo/user-profile/models"
)

type UserSerializer interface {
	Decode(input []byte) (*m.User, error)
	Encode(input *m.User) ([]byte, error)
	DecodeMap(input []byte) (map[string]interface{}, error)
	EncodeMap(input map[string]interface{}) ([]byte, error)
	DecodeResult(input []byte) (*helper.ResultInfo, error)
	EncodeResult(input *helper.ResultInfo) ([]byte, error)
}
