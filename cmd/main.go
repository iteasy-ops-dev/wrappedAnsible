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
	"iteasy.wrappedAnsible/internal/ansible"
	"iteasy.wrappedAnsible/internal/handlers"
	"iteasy.wrappedAnsible/internal/router"
)

var (
	port = ":8080"
)

func init() {
	var g errgroup.Group
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("⚙️ Wrapped Ansible Server Init.")
	initJsonData := `{
			"type": "init",
			"name": "서버 초기화 실행.",
		  "options": {}
	  }
	`

	// var payload []byte
	g.Go(func() error {
		var err error
		initAnsible := ansible.GetAnsibleFromFactory(ctx, []byte(initJsonData))
		// payload, err = ansible.Excuter(initAnsible)
		_, err = ansible.Excuter(initAnsible)
		return err
	})

	if err := g.Wait(); err != nil {
		panic("❌ 서버 초기화 실패.")
	}
}

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
