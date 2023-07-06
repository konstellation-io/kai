package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/delivery/http/httperrors"
	"github.com/konstellation-io/kai/engine/admin-api/delivery/http/token"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logging"
	"github.com/labstack/echo"
)

func extractTokenFromAuthHeader(authHeader string) string {
	if len(strings.Split(authHeader, " ")) == 2 {
		return strings.Split(authHeader, " ")[1]
	}

	return ""
}

func NewJwtAuthMiddleware(_ *config.Config, logger logging.Logger, tokenParser *token.Parser) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			b, err := io.ReadAll(c.Request().Body)
			if err != nil {
				return err
			}

			bodyBuffer := bytes.NewBuffer(b)
			c.Request().Body = io.NopCloser(bodyBuffer)

			operation, err := getOperation(b)
			if err != nil {
				logger.Info("Unable to get graphql operation")
			}

			authHeader := c.Request().Header.Get("Authorization")
			plainToken := extractTokenFromAuthHeader(authHeader)

			user, err := tokenParser.GetUser(plainToken)
			if err != nil {
				logger.Warn("No token found in context")

				if operation == "info" {
					fmt.Println("here")
					return next(c)
				}

				return httperrors.HTTPErrUnauthorized
			}

			c.Set("user", user)

			return next(c)
		}
	}
}

type graphqlReqBody struct {
	Query string `json:"query"`
}

func getOperation(body []byte) (string, error) {
	var graphqlBody graphqlReqBody
	err := json.Unmarshal(body, &graphqlBody)
	if err != nil {
		return "", err
	}

	operation := regexp.MustCompile("{.*{").FindString(graphqlBody.Query)
	operation = strings.ReplaceAll(operation, `\n`, "")
	operation = regexp.MustCompile(`[a-zA-Z]+`).FindString(operation)

	return operation, nil
}
