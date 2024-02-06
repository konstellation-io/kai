package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"regexp"
	"strings"

	"github.com/go-logr/logr"
	"github.com/labstack/echo/v4"
)

const _operationNameKey = "operationName"

func NewGraphQLOperationMiddleware(logger logr.Logger) echo.MiddlewareFunc {
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
				return next(c)
			}

			c.Set(_operationNameKey, operation)

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
