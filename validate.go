package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/bufbuild/connect-go"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

type validator interface {
	ValidateAll() error
}

type validationError interface {
	Field() string
	Reason() string
	Cause() error
	Key() bool
	ErrorName() string
	Error() string
}

type validateMultiError interface {
	Error() string
	AllErrors() []error
}

func validate(req any) error {
	switch v := req.(type) {
	case validator:
		if err := v.ValidateAll(); err != nil {
			fieldInfo := &errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{},
			}

			switch vme := err.(type) {

			case validateMultiError:
				errs := vme.AllErrors()

				for i := 0; i < len(errs); i++ {
					switch ve := errs[i].(type) {
					case validationError:

						fieldInfo.FieldViolations = append(
							fieldInfo.FieldViolations,
							&errdetails.BadRequest_FieldViolation{
								Field: fmt.Sprintf("%s.%s",
									strings.TrimSuffix(ve.ErrorName(), "ValidationError"),
									ve.Field()),
								Description: ve.Reason(),
							},
						)
					}
				}

			default:
				log.Print(err)
			}

			err := connect.NewError(connect.CodeInvalidArgument, errors.New("Client specified an invalid argument. Check error details for more information."))

			if detail, detailErr := connect.NewErrorDetail(fieldInfo); detailErr == nil {
				err.AddDetail(detail)
			}
			return err
		}
	}
	return nil
}

func ValidatorInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if err := validate(req.Any()); err != nil {
				return nil, err
			}
			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
