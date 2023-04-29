package main

import (
	"context"
	"fmt"
	"github.com/rueian/rueidis"
	"github.com/snowmerak/rmq/gen/bean"
	"github.com/snowmerak/rmq/lib/client/subset"
	"google.golang.org/grpc"
	"time"
)

func main() {
	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{"localhost:6379"},
	})
	if err != nil {
		panic(err)
	}

	beanBuilder := bean.New()
	beanBuilder.AddSubset(subset.New(client))
	beanContainer := beanBuilder.Build()

	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		fmt.Println(beanContainer.Subset().GetMembers(context.Background(), "test"))
	}

	grpc.Dial()
}
