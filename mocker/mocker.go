package mocker

import (
	"bytes"
	"context"
	"fmt"
	"github.com/byorty/contractor/common"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"strings"
)

const (
	variableExpr = "${%s}"
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
	loggerFactory common.LoggerFactory,
) Mocker {
	return &mocker{
		args:              args,
		dataCrawler:       dataCrawler,
		mediaConverter:    mediaConverter,
		expressionFactory: expressionFactory,
		commonLogger:      loggerFactory.CreateCommonLogger(),
		successLogger:     loggerFactory.CreateSuccessLogger(),
		errorLogger:       loggerFactory.CreateErrorLogger(),
		router:            mux.NewRouter(),
	}
}

type mocker struct {
	args              common.Arguments
	dataCrawler       common.DataCrawler
	expressionFactory common.ExpressionFactory
	mediaConverter    common.MediaConverter
	commonLogger      common.Logger
	successLogger     common.Logger
	errorLogger       common.Logger
	router            *mux.Router
}

func (m *mocker) Configure(ctx context.Context, containers common.TemplateContainer) error {
	for name, template := range containers {
		for statusCode, exampleContainer := range template.ExpectedResponses {
			for mediaType, example := range exampleContainer {
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

				path := template.GetPath()
				headerParams := make(map[string]string)
				queryParams := make(map[string]string)
				for key, value := range m.args.Variables {
					path = strings.ReplaceAll(path, fmt.Sprintf(variableExpr, key), value)
					buf = bytes.ReplaceAll(buf, []byte(fmt.Sprintf(variableExpr, key)), []byte(value))

					for headerName, headerValue := range template.HeaderParams {
						headerParams[headerName] = strings.ReplaceAll(fmt.Sprint(headerValue), fmt.Sprintf(variableExpr, key), value)
					}

					for queryName, queryValue := range template.QueryParams {
						queryParams[queryName] = strings.ReplaceAll(fmt.Sprint(queryValue), fmt.Sprintf(variableExpr, key), value)
					}
				}

				route := m.router.Methods(template.Method)
				route.Path(path)

				template.HeaderParams[common.HeaderAccept] = mediaType
				for headerName, headerValue := range headerParams {
					route.Headers(headerName, fmt.Sprint(headerValue))
				}

				for queryName, queryValue := range queryParams {
					route.Queries(queryName, fmt.Sprint(queryValue))
				}

				m.commonLogger.PrintGroup("Request: %s", name)
				m.commonLogger.PrintParameter("Path", template.Path)
				m.commonLogger.PrintParameter("Method", strings.ToUpper(template.Method))
				m.commonLogger.PrintParameter("Status Code", statusCode)
				m.commonLogger.PrintParameters("Header Parameters", template.HeaderParams)
				m.commonLogger.PrintParameters("Path Parameters", template.PathParams)
				m.commonLogger.PrintParameters("Query Parameters", template.QueryParams)
				m.commonLogger.PrintParameter("Body", string(buf))

				route.HandlerFunc(m.createHandler(name, mediaType, statusCode, buf))
			}
		}
	}

	return nil
}

func (m *mocker) createHandler(name, mediaType string, statusCode int, buf []byte) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add(common.HeaderContentType, mediaType)
		writer.WriteHeader(statusCode)
		writer.Write(buf)
		m.successLogger.PrintSubGroupName("[SUCCESS] Example %s serve success", name)
	}
}

func (m *mocker) Run() error {
	baseUrl, err := url.Parse(m.args.Url)
	if err != nil {
		return err
	}

	m.commonLogger.PrintGroup("Serve requests...")
	m.router.NotFoundHandler = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		m.errorLogger.PrintSubGroupName("[FAILURE] Example not found for %s", request.URL.String())
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
