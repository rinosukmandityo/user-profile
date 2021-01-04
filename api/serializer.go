package api

import (
	slz "github.com/rinosukmandityo/user-profile/api/serializer"
	js "github.com/rinosukmandityo/user-profile/api/serializer/json"
)

var (
	ContentTypeJson    = "application/json"
	ContentTypeMsgPack = "application/x-msgpack"
)

func GetSerializer(contentType string) slz.UserSerializer {
	if contentType == ContentTypeMsgPack {
		// return message pack object here
	}
	return &js.User{}
}
