package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
	"iteasy.wrappedAnsible/internal/handlers"
	"iteasy.wrappedAnsible/internal/router"
)

var (
	port = ":8080"
)

// func init() {
// 	client := model.GetMongoInstance()
// 	if client != nil {
// 		model.PingMongoDB(client)
// 	}

// 	fmt.Println("⚙️ Wrapped Ansible Server Init.")
// 	ctx, cancel := context.WithCancel(context.Background())
// 	g, ctx := errgroup.WithContext(ctx)
// 	defer cancel()

// 	initJsonData := `{
// 			"type": "init",
// 			"name": "서버 초기화 실행.",
// 		  "options": {}
// 	  }
// 	`
// 	init := ansible.GennerateInitType{
// 		Ctx:      ctx,
// 		JsonData: []byte(initJsonData),
// 	}

// 	g.Go(func() error {
// 		var err error
// 		initAnsible, err := ansible.GetAnsibleFromFactory(init)
// 		if err != nil {
// 			return fmt.Errorf("failed to get Ansible from factory: %w", err)
// 		}
// 		r, err := ansible.Excuter(initAnsible)
// 		if !r.Status || err != nil {
// 			return fmt.Errorf("failed to execute Ansible: %w", err)
// 		}
// 		return nil
// 	})

// 	if err := g.Wait(); err != nil {
// 		panic("❌ 서버 초기화 실패.")
// 	}
// }

func main() {
	var g errgroup.Group
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := &http.Server{
		Addr:    port,
		Handler: handlers.CorsMiddleware(router.NewRouter()),
	}

	g.Go(func() error {
		fmt.Printf("✅ Welcome Wrapped Ansible Server. PORT %s\n", port)
		return server.ListenAndServe()
	})

	// Signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	g.Go(func() error {
		sig := <-sigCh
		fmt.Printf("👋 Received signal: %v, initiating shutdown...\n", sig)
		cancel()
		return nil
	})

	g.Go(func() error {
		<-ctx.Done()
		fmt.Println("👋 Shutting down server...")
		return server.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		log.Fatalf("❌ Server error: %v", err)
	}
}
