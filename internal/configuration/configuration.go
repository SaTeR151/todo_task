package configuration

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/sater-151/todo-list/internal/pkg/errorspkg"
	"github.com/sater-151/todo-list/internal/pkg/validate"
	"github.com/spf13/viper"
)

const (
	_defaultConfigurationsPath = "configuration.yaml"
)

type (
	Configurations struct {
		Logger     *Logger     `mapstructure:"Logger" validate:"required"`
		HTTPServer *HTTPServer `mapstructure:"HttpClient" validate:"required"`
		Version    string      `validate:"-"`
	}

	HTTPServer struct {
		Port              string        `mapstructure:"Port" validate:"required,min=1"`
		ReadTimeout       time.Duration `mapstructure:"ReadTimeout" validate:"required"`
		WriteTimeout      time.Duration `mapstructure:"WriteTimeout" validate:"required"`
		ShutdownTimeout   time.Duration `mapstructure:"ShutdownTimeout" validate:"required"`
		ReadHeaderTimeout time.Duration `mapstructure:"ReadHeaderTimeout" validate:"required"`
		IdleTimeout       time.Duration `mapstructure:"IdleTimeout" validate:"required"`
		MaxHeaderBytes    int           `mapstructure:"MaxHeaderBytes" validate:"gt=0"`
		TestPass          string        `mapstructure:"TestPass" validate:"required"`
	}

	Logger struct {
		Level slog.Level `mapstructure:"Level" validate:"min=-4,max=8"`
	}
)

var password string

func NewConfig() (*Configurations, error) {
	vp := viper.New()

	vp.SetConfigFile(_defaultConfigurationsPath)
	if err := vp.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading configuration: %w", err)
	}

	var configurations Configurations
	if err := vp.Unmarshal(&configurations); err != nil {
		return nil, fmt.Errorf("unmarshaling configuration: %w", err)
	}

	if err := validate.Struct(configurations); err != nil {
		return nil, errorspkg.NewValidationError("configuration.NewConfig", configurations, err)
	}

	return &configurations, nil
}

func GetPass() (pass string) {
	pass = password
	return pass
}
