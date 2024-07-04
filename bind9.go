package main

import (
	"flag"
	"fmt"
	"log"

	"bind9-manager-service/internal/config"
	"bind9-manager-service/internal/handler"
	"bind9-manager-service/internal/model"
	"bind9-manager-service/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"

	_ "github.com/mattn/go-sqlite3"
)

var configFile = flag.String("f", "etc/bind9-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx, err := svc.NewServiceContext(c)
	if err != nil {
		log.Fatalf("failed to initialize service context: %v", err)
	}

	err = model.GenerateAllZoneFiles(ctx.DataSource, c.BindPath)
	if err != nil {
		log.Fatalf("failed to generate files: %v", err)
	}

	defer ctx.DataSource.Close()

	if err := svc.StartBind9(); err != nil {
		log.Fatalf("failed to start bind9: %v", err)
	}

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()

}
