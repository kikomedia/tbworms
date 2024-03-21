package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"tbworms/server"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type ApplicationConfig struct {
	Port         int
	CookieStore  *sessions.CookieStore
	ShutdownWait *time.Duration
}

type ApplicationContext struct {
	Config     ApplicationConfig
	GameServer *server.GameServer
	HttpServer http.Server
}

var applicationContext = &ApplicationContext{}

func prepareApplicationConfig(appCtx *ApplicationContext) {
	flag.IntVar(&appCtx.Config.Port, "p", 8080, "server port")
	appCtx.Config.ShutdownWait = flag.Duration("shutdown-timeout", time.Second*15, "Duration Shutdown in Seconds")

	flag.Parse()

	appCtx.Config.CookieStore = sessions.NewCookieStore([]byte("someprettyprivatekey")) // :D
}

func prepareGameServer(appCtx *ApplicationContext) {
	appCtx.GameServer = server.NewGameServer()
	go appCtx.GameServer.Run()
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "dist/index.html")
}

func cleanup(appCtx *ApplicationContext) {
	log.Println("Shutting Down...")

	appCtx.GameServer.Shutdown()
	ctx, cancel := context.WithTimeout(context.Background(), *appCtx.Config.ShutdownWait)
	defer cancel()

	appCtx.HttpServer.Shutdown(ctx)
}

func prepareShutdown(appCtx *ApplicationContext) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup(appCtx)
		log.Println("...server closed.")
		os.Exit(1)
	}()
}

func runApplicationServer(appCtx *ApplicationContext) {
	// Prepare routes
	route := mux.NewRouter()
	route.HandleFunc("/", serveHome)
	route.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.ServeWs(appCtx.GameServer, w, r)
	})

	// Prepare server config
	server := &http.Server{
		Addr:              ":" + strconv.Itoa(appCtx.Config.Port),
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           route,
	}

	// Run server
	log.Println("Starting at http://localhost" + server.Addr)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error: ", err)
		os.Exit(1)
	}
}

func main() {
	applicationContext = &ApplicationContext{}
	prepareShutdown(applicationContext)
	prepareApplicationConfig(applicationContext)
	prepareGameServer(applicationContext)
	runApplicationServer(applicationContext)
}
