package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/spanner"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"github.com/kelseyhightower/envconfig"
	"github.com/sinmetal/gcpmetadata"
	"go.opencensus.io/trace"
)

type EnvConfig struct {
	SpannerDatabase string `required:"true"`
	Goroutine       int    `default:"3"`
	RowCount        int    `default:"1"`
}

func main() {
	var env EnvConfig
	if err := envconfig.Process("ssrhr", &env); err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("ENV_CONFIG %+v\n", env)

	project, err := gcpmetadata.GetProjectID()
	if err != nil {
		panic(err)
	}

	{
		exporter, err := stackdriver.NewExporter(stackdriver.Options{
			ProjectID: project,
		})
		if err != nil {
			panic(err)
		}
		trace.RegisterExporter(exporter)
		// trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	}

	ctx := context.Background()

	endCh := make(chan error, 10)

	sc, err := spanner.NewClient(ctx, env.SpannerDatabase)
	if err != nil {
		panic(err)
	}
	sss := NewSmallSizeStore(sc)

	for i := 1; i <= env.RowCount; i++ {
		goGetSmallSize(sss, env.Goroutine, i, endCh)
	}

	fmt.Println("Start Wait EndCh...")
	err = <-endCh
	fmt.Printf("BOMB %+v", err)
}
