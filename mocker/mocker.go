package mocker

import (
	"context"
	"fmt"
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/logger"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"strings"
)

type Mocker interface {
	Configure(ctx context.Context, containers common.TemplateContainer) error
	Run() error
}

func NewFxMocker(
	args common.Arguments,
	dataCrawler common.DataCrawler,
	expressionFactory common.ExpressionFactory,
	mediaConverter common.MediaConverter,
	logger logger.Logger,
) Mocker {
	return &mocker{
		args:              args,
		dataCrawler:       dataCrawler,
		mediaConverter:    mediaConverter,
		expressionFactory: expressionFactory,
		logger:            logger,
		router:            mux.NewRouter(),
	}
}

type mocker struct {
	args              common.Arguments
	dataCrawler       common.DataCrawler
	expressionFactory common.ExpressionFactory
	mediaConverter    common.MediaConverter
	logger            logger.Logger
	router            *mux.Router
}

func (m *mocker) Configure(ctx context.Context, containers common.TemplateContainer) error {
	for name, template := range containers {
		for statusCode, exampleContainer := range template.ExpectedResponses {
			for mediaType, example := range exampleContainer {
				route := m.router.Methods(template.Method)
				route.Path(template.GetPath())

				template.HeaderParams[common.HeaderContentType] = mediaType
				for headerName, headerValue := range template.HeaderParams {
					route.Headers(headerName, fmt.Sprint(headerValue))
				}

				for queryName, queryValue := range template.QueryParams {
					route.Queries(queryName, fmt.Sprint(queryValue))
				}

				m.dataCrawler.Walk(example, func(_ string, data interface{}) {
					switch d := data.(type) {
					case map[string]interface{}:
						for key, value := range d {
							switch v := value.(type) {
							case map[string]interface{}, []interface{}:
								continue
							default:
								generator, err := m.expressionFactory.Create(common.ExpressionTypeGenerator, v)
								if err != nil {
									continue
								}

								d[key] = generator.(Generator).Generate()
							}
						}
					case []interface{}:
						for i, item := range d {
							switch item.(type) {
							case map[string]interface{}, []interface{}:
								continue
							default:
								generator, err := m.expressionFactory.Create(common.ExpressionTypeGenerator, item)
								if err != nil {
									continue
								}

								d[i] = generator.(Generator).Generate()
							}
						}
					}
				})

				buf, err := m.mediaConverter.Marshal(common.MediaType(mediaType), example)
				if err != nil {
					return err
				}

				m.logger.PrintGroup("Request: %s", name)
				m.logger.PrintSubGroup("Path: %s", template.Path)
				m.logger.PrintSubGroup("Method: %s", strings.ToUpper(template.Method))
				m.logger.PrintSubGroup("Status Code: %d", statusCode)
				m.logger.PrintParameters("Header Parameters", template.HeaderParams)
				m.logger.PrintParameters("Path Parameters", template.PathParams)
				m.logger.PrintParameters("Query Parameters", template.QueryParams)
				m.logger.PrintSubGroup("Body: %s", string(buf))

				route.HandlerFunc(m.createHandler(name, mediaType, statusCode, buf))
			}
		}
	}

	return nil
}

func (m *mocker) createHandler(name, mediaType string, statusCode int, buf []byte) func(http.ResponseWriter, *http.Request) {
	successLogger := m.logger.ToSuccessColors()
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add(common.HeaderContentType, mediaType)
		writer.WriteHeader(statusCode)
		writer.Write(buf)
		successLogger.PrintSubGroup("[SUCCESS] Example %s serve success", name)
	}
}

func (m *mocker) Run() error {
	baseUrl, err := url.Parse(m.args.Url)
	if err != nil {
		return err
	}

	m.logger.PrintGroup("Serve requests...")

	errorLogger := m.logger.ToErrorColors()
	m.router.NotFoundHandler = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		errorLogger.PrintSubGroup("[FAILURE] Example not found for %s", request.URL.String())
	})
	srv := &http.Server{
		Handler: m.router,
		Addr:    baseUrl.Host,
	}

	if len(m.args.CertFilename) > 0 && len(m.args.KeyFilename) > 0 {
		return srv.ListenAndServeTLS(m.args.CertFilename, m.args.KeyFilename)
	} else {
		return srv.ListenAndServe()
	}
}
