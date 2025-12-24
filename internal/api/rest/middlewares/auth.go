package middlewares

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sater-151/todo-list/internal/models"
	"github.com/sater-151/todo-list/internal/pkg/errorspkg"
	"github.com/sater-151/todo-list/internal/pkg/validate"
)

type MiddlewaresDependencies struct {
	Password string `validate:"required"`
}

type Middlewares struct {
	password string
}

func NewMiddlewares(d *MiddlewaresDependencies) (*Middlewares, error) {
	if err := validate.Struct(d); err != nil {
		return nil, errorspkg.NewValidationError("rest.NewTodoTaskHandlers", d, err)
	}

	return &Middlewares{
		password: d.Password,
	}, nil
}

func (m *Middlewares) Auth(next http.Handler) http.Handler {
	fn := func(res http.ResponseWriter, req *http.Request) {
		pass := m.password
		if pass != "" {
			cookie, err := req.Cookie("token")
			if err != nil {
				slog.Error(err.Error())
				ErrorHandler(res, err, http.StatusUnauthorized)
				return
			}
			jwtCookie := cookie.Value
			jwtToken, err := jwt.Parse(jwtCookie, func(_ *jwt.Token) (interface{}, error) {
				return []byte(pass), nil
			})
			if err != nil {
				slog.Error(err.Error())
				ErrorHandler(res, err, http.StatusUnauthorized)
				return
			}
			if !jwtToken.Valid {
				slog.Error("token is invalid")
				ErrorHandler(res, fmt.Errorf("token is invalid"), http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(res, req)
	}

	return http.HandlerFunc(fn)

}

func ErrorHandler(res http.ResponseWriter, err error, status int) {
	var errJS models.Error
	errJS.Err = err.Error()
	res.WriteHeader(status)
	if err := json.NewEncoder(res).Encode(errJS); err != nil {
		slog.Error(err.Error())
		ErrorHandler(res, err, http.StatusInternalServerError)
		return
	}
}
