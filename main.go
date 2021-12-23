package main

import (
	"context"
	"github.com/byorty/contractor/common"
	"github.com/byorty/contractor/converter"
	"github.com/byorty/contractor/logger"
	"github.com/byorty/contractor/tester"
	"log"
)

func main() {
	ctx := context.Background()
	//loader := &openapi3.Loader{Context: ctx}
	//doc, err := loader.LoadFromFile("open_api_v3.yml")
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//err = doc.Validate(ctx)
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//router, err := legacyrouter.NewRouter(doc)
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//httpReq, err := http.NewRequest(http.MethodGet, "https://st-widget.ogon.ru/v1/news/11876", nil)
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//// Find route
	//route, pathParams, err := router.FindRoute(httpReq)
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	////doc.Paths[""].Patch.Responses[""].
	//
	//// Validate request
	//requestValidationInput := &openapi3filter.RequestValidationInput{
	//	request:    httpReq,
	//	PathParams: pathParams,
	//	Route:      route,
	//}
	//if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
	//	log.Panic(err)
	//}
	//
	//var (
	//	respStatus      = 200
	//	respContentType = "application/json"
	//	respBody        = bytes.NewBufferString(`{}`)
	//)
	//
	//log.Println("Response:", respStatus)
	//responseValidationInput := &openapi3filter.ResponseValidationInput{
	//	RequestValidationInput: requestValidationInput,
	//	Status:                 respStatus,
	//	Header:                 http.Header{"Content-Type": []string{respContentType}},
	//}
	//if respBody != nil {
	//	data, _ := json.Marshal(respBody)
	//	responseValidationInput.SetBodyBytes(data)
	//}
	//
	//if err := openapi3filter.ValidateResponse(ctx, responseValidationInput); err != nil {
	//	log.Printf("%T\n", err.(*openapi3filter.ResponseError).Unwrap().(*openapi3.SchemaError).Value)
	//	log.Println(err.(*openapi3filter.ResponseError).Unwrap().(*openapi3.SchemaError).JSONPointer())
	//	//log.Panic(err)
	//}

	l := &logger.Logger{}

	w := tester.NewFxWorker(
		converter.NewFxOa3Converter(l),
		tester.NewFxTester(),
	)

	err := w.Configure(ctx, common.Arguments{
		SpecFilename: "open_api_v3.yml",
	})
	if err != nil {
		log.Panic(err)
	}

	err = w.Run()
	if err != nil {
		log.Panic(err)
	}
}
