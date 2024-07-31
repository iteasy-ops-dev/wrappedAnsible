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
	config "iteasy.wrappedAnsible/configs"
	"iteasy.wrappedAnsible/internal/ansible"
	"iteasy.wrappedAnsible/internal/handlers"
	"iteasy.wrappedAnsible/internal/model"
	"iteasy.wrappedAnsible/internal/router"
	"iteasy.wrappedAnsible/pkg/utils"
)

var (
	port = ":8080"
)

func _mongo() {
	client := model.GetMongoInstance()
	if client != nil {
		model.PingMongoDB(client)
	}

	adminEmail := config.GlobalConfig.Default.Admin
	adminName := "admin"
	password := config.GlobalConfig.Default.Password

	hashedPassword, _ := utils.HashingPassword(password)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a := model.NewAuth(ctx)
	a.SetEmail(adminEmail)
	a.SetName(adminName)
	a.SetPassword(string(hashedPassword))

	if err := a.SignUp(); err != nil {

		switch err.(type) {
		case *model.AlreadyExistsError:
			return
		default:
			panic(err.Error())
		}
	}

}

func _ansible() {
	fmt.Println("⚙️ Wrapped Ansible Server Init.")
	ctx, cancel := context.WithCancel(context.Background())
	// g, ctx := errgroup.WithContext(ctx)
	defer cancel()

	initJsonData := `{
			"type": "init",
			"name": "서버 초기화 실행.",
		  "options": {}
	  }
	`
	init := ansible.GennerateInitType{
		Ctx:      ctx,
		JsonData: []byte(initJsonData),
	}

	initAnsible, _ := ansible.GetAnsibleFromFactory(init)
	// if err != nil {
	// 	panic("❌ 서버 초기화 실패.")
	// }
	r, err := ansible.Excuter(initAnsible)
	if !r.Status || err != nil {
		panic("❌ 서버 초기화 실패.")
	}
}

func init() {
	_mongo()
	_ansible()
}

func main() {
	var g errgroup.Group
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := &http.Server{
		Addr:    port,
		Handler: handlers.CorsMiddleware(router.NewRouter()),
		// TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12},
	}

	g.Go(func() error {
		fmt.Printf("✅ Welcome Wrapped Ansible Server. PORT %s\n", port)
		return server.ListenAndServe()
		// return server.ListenAndServeTLS("server.crt", "server.key")
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
