package test

import (
	"context"
	v1 "gitee.com/moyusir/service-centre/api/serviceCenter/v1"
	"gitee.com/moyusir/service-centre/internal/conf"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"os"
	"testing"
	"time"
)

func newApp(logger log.Logger, hs *http.Server) *kratos.App {
	// go build -ldflags "-X main.Version=x.y.z"
	var (
		// Name is the name of the compiled software.
		Name string
		// Version is the version of the compiled software.
		Version string
		id, _   = os.Hostname()
	)
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			hs,
		),
	)
}
func StartServiceCenterServer(t *testing.T) v1.UserHTTPClient {
	bootstrap, err := conf.LoadConfig("../../configs/config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	app, cleanUp, err := initApp(bootstrap.Server, bootstrap.Data, log.NewStdLogger(os.Stdout))
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		err := app.Run()
		if err != nil {
			return
		}
	}()
	t.Cleanup(func() {
		app.Stop()
		<-done
		cleanUp()
	})

	for {
		select {
		case <-done:
			t.Fatal("failed to start app")
		default:
			client, err := http.NewClient(context.Background(),
				http.WithEndpoint("localhost:8000"),
				http.WithTimeout(time.Hour),
			)
			if err != nil {
				continue
			}
			t.Cleanup(func() {
				client.Close()
			})
			return v1.NewUserHTTPClient(client)
		}
	}
}
