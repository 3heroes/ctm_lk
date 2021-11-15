package main

import (
	"context"
	"net/http"
	"sync"

	"ctm_lk/internal/config"
	"ctm_lk/internal/database"
	"ctm_lk/internal/server"
	"ctm_lk/pkg/logger"
	"ctm_lk/pkg/ossignal"
	"ctm_lk/pkg/workers"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello, World</h1>"))
}

func main() {
	//makeMigrations()
	var wg sync.WaitGroup
	logger.NewLogs()
	defer logger.Close()
	logger.Info("Старт сервера")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	config.NewConfig()

	sDB := database.OpenDBConnect()
	defer sDB.Close()

	wg.Add(1)
	go func() {
		ossignal.HandleQuit(cancel)
		wg.Done()
	}()

	w := workers.NewWorkersPool(10)
	defer w.Close()

	s := new(server.Server)
	s.ServerDB = sDB
	s.Start(ctx)
	wg.Wait()
	logger.Info("Сервер остановлен")

}
