package responders

import (
	codecservices "github.com/stretchrcom/codecs/services"
	"github.com/stretchrcom/goweb/context"
)

const (
	DefaultStandardFieldDataKey   string = "d"
	DefaultStandardFieldStatusKey string = "s"
	DefaultStandardFieldErrorsKey string = "e"
)

type GowebAPIResponder struct {
	httpResponder HTTPResponder
	codecService  codecservices.CodecService

	// field names

	// StandardFieldDataKey is the response object field name for the data.
	StandardFieldDataKey string

	// StandardFieldStatusKey is the response object field name for the status.
	StandardFieldStatusKey string

	// StandardFieldErrorsKey is the response object field name for the errors.
	StandardFieldErrorsKey string
}

func NewGowebAPIResponder(httpResponder HTTPResponder) *GowebAPIResponder {
	api := new(GowebAPIResponder)
	api.httpResponder = httpResponder
	api.StandardFieldDataKey = DefaultStandardFieldDataKey
	api.StandardFieldStatusKey = DefaultStandardFieldStatusKey
	api.StandardFieldErrorsKey = DefaultStandardFieldErrorsKey
	return api
}

// SetCodecService sets the codec service to use.
func (a *GowebAPIResponder) SetCodecService(service codecservices.CodecService) {
	a.codecService = service
}

// GetCodecService gets the codec service that will be used by this object.
func (a *GowebAPIResponder) GetCodecService() codecservices.CodecService {

	if a.codecService == nil {
		a.codecService = new(codecservices.WebCodecService)
	}

	return a.codecService
}

// WriteResponseObject writes the status code and response object to the HttpResponseWriter in
// the specified context, in the format best suited based on the request.
//
// Goweb uses the WebCodecService to decide which codec to use when responding
// see http://godoc.org/github.com/stretchrcom/codecs/services#WebCodecService for more information.
//
// This method should be used when the Goweb Standard Response Object does not satisfy the needs of
// the API, but other Respond* methods are recommended.
func (a *GowebAPIResponder) WriteResponseObject(ctx context.Context, status int, responseObject interface{}) error {

	service := a.GetCodecService()
	codec, codecError := service.GetCodec("application/json")

	if codecError != nil {
		return codecError
	}

	output, marshalErr := codec.Marshal(responseObject, nil)

	if marshalErr != nil {
		return marshalErr
	}

	// use the HTTP responder to respond
	a.httpResponder.With(ctx, status, output)

	return nil

}

// Responds to the Context with the specified status, data and errors.
func (a *GowebAPIResponder) Respond(ctx context.Context, status int, data interface{}, errors []string) error {

	// make the standard response object
	sro := map[string]interface{}{a.StandardFieldStatusKey: status}

	if data != nil {
		sro[a.StandardFieldDataKey] = data
	}
	if len(errors) > 0 {
		sro[a.StandardFieldErrorsKey] = errors
	}

	return a.WriteResponseObject(ctx, status, sro)

}
