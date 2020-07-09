package middleware

import (
	"context"

	"golang.org/x/text/language"

	"google.golang.org/grpc"

	_ "github.com/caos/zitadel/internal/statik"
)

func TranslationHandler(defaultLanguage language.Tag) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	translator := newZitadelTranslator(defaultLanguage)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if loc, ok := resp.(localizers); ok && resp != nil {
			translateFields(ctx, loc, translator)
		}
		return resp, err
	}
}
