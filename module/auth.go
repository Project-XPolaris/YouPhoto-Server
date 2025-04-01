package module

import (
	"bytes"
	"encoding/gob"

	"github.com/allentom/harukap"
	"github.com/allentom/harukap/commons"
	"github.com/allentom/harukap/module/auth"
	"github.com/projectxpolaris/youphoto/config"
	"github.com/projectxpolaris/youphoto/database"
)

var Auth = &auth.AuthModule{
	Plugins: []harukap.AuthPlugin{},
}

func CreateAuthModule() {
	Auth.ConfigProvider = config.DefaultConfigProvider
	// Auth.AddCacheStore(&UserSerializer{})
	Auth.InitModule()
}

type UserSerializer struct {
}

func (s *UserSerializer) Serialize(data interface{}) ([]byte, error) {
	user := data.(*database.User)
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(user)
	if err != nil {
		return nil, err
	}
	return result.Bytes(), nil
}

func (s *UserSerializer) Deserialize(raw []byte) (commons.AuthUser, error) {
	var user database.User
	decoder := gob.NewDecoder(bytes.NewReader(raw))
	err := decoder.Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
