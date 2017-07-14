package service

import (
	"golang.org/x/net/context"

	"github.com/go-kit/kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(PartnerService) PartnerService

// LoggingMiddleware takes a logger as a dependency
// and returns a ServiceMiddleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next PartnerService) PartnerService {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   PartnerService
}

func (mw loggingMiddleware) GetPartnerDataByKeyValue(ctx context.Context, key string, value string, group string) (partnerId int32, partnerCode string, attributes map[string]string, err error) {
	defer func() {
		mw.logger.Log("method", "KeyValue", "id", partnerId, "code", partnerCode, "attributes", attributes, "err", err)
	}()
	return mw.next.GetPartnerDataByKeyValue(ctx, key, value, group)
}

func (mw loggingMiddleware) GetDataById(ctx context.Context, id int32, code string, group string) (partnerId int32, partnerCode string, attributes map[string]string, err error) {
	defer func() {
		mw.logger.Log("method", "ById", "id", partnerId, "code", partnerCode, "attributes", attributes, "err", err)
	}()
 	return mw.next.GetDataById(ctx, id, code, group)
}
