package mocker

import (
	"context"
	"fmt"
	"github.com/byorty/contractor/common"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type Mocker interface {
	Configure(ctx context.Context, containers common.TemplateContainer) error
	Run() error
}

func NewFxMocker(
	dataCrawler common.DataCrawler,
	expressionFactory common.ExpressionFactory,
	mediaConverter common.MediaConverter,
) Mocker {
	return &mocker{
		dataCrawler:       dataCrawler,
		mediaConverter:    mediaConverter,
		expressionFactory: expressionFactory,
		router:            mux.NewRouter(),
	}
}

type mocker struct {
	dataCrawler       common.DataCrawler
	expressionFactory common.ExpressionFactory
	mediaConverter    common.MediaConverter
	router            *mux.Router
}

func (m *mocker) Configure(ctx context.Context, containers common.TemplateContainer) error {
	for _, template := range containers {
		for statusCode, exampleContainer := range template.ExpectedResponses {
			for mediaType, example := range exampleContainer {
				route := m.router.Methods(template.Method)
				route.Path(template.GetPath())

				route.Headers(common.HeaderContentType, mediaType)
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

				route.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					writer.Header().Add(common.HeaderContentType, mediaType)
					writer.WriteHeader(statusCode)
					writer.Write(buf)
				})
			}
		}
	}

	return nil
}

func (m *mocker) Run() error {
	srv := &http.Server{
		Handler:      m.router,
		Addr:         ":8181",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv.ListenAndServe()
}
