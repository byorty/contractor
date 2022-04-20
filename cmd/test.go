package main

import (
	"fmt"
	"github.com/byorty/contractor/common"
)

type L2 struct {
	common.List[int]
}

func (l L2) Print() {
	fmt.Println(l.Entries())
}

func New() L2 {
	return L2{common.NewList[int]()}
}

func main() {
	l := New()
	l.Add(1)
	l.Add(2)
	l.Print()
	//ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Second*10))
	//c := cl.NewFxGraylogClient(common.Arguments{
	//	Url: "graylog.setpartnerstv.ru",
	//})
	//sorting := "timestamp:asc"
	//resp, err := c.Saved.SearchRelative(&saved.SearchRelativeParams{
	//	//Fields:     nil,
	//	//Filter:     nil,
	//	//Limit:      nil,
	//	//Offset:     nil,
	//	Query:   "app:autopayment-processor",
	//	Range:   300,
	//	Sort:    &sorting,
	//	Context: ctx,
	//	//HTTPClient: nil,
	//}, client.BasicAuth("1geh46b37k3mib77lj8lt401mkah24d1alomosish521tsfqsuaf", "token"))
	//log.Print(resp.Payload.Messages[0].Message.(map[string]interface{})["correlation_id"])
	//log.Print(resp.Payload.Messages[0].Message.(map[string]interface{}))
	//log.Print(err)
}
