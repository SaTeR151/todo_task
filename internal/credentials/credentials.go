package credentials

import (
	"fmt"

	"github.com/sater-151/todo-list/internal/pkg/errorspkg"
	"github.com/sater-151/todo-list/internal/pkg/validate"
	"github.com/spf13/viper"
)

const (
	_defaultCredentialPath = "credentials.yaml"
)

type (
	Credentials struct {
		Postgres *Postgres `mapstructure:"Postgres" validate:"required"`
	}

	Postgres struct {
		ConnString string `mapstructure:"ConnString" validate:"required"`
	}
)

func NewCredentials() (*Credentials, error) {
	vp := viper.New()

	vp.SetConfigFile(_defaultCredentialPath)
	if err := vp.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading credentials: %w", err)
	}

	var credentials Credentials
	if err := vp.Unmarshal(&credentials); err != nil {
		return nil, fmt.Errorf("unmarshal credentials: %w", err)
	}

	if err := validate.Struct(credentials); err != nil {
		return nil, errorspkg.NewValidationError("credentials.NewCredentials", credentials, err)
	}

	return &credentials, nil

}
