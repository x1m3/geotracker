package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"github.com/x1m3/geotracker/HTTPServer"
	"github.com/x1m3/geotracker/command"
	"github.com/x1m3/geotracker/repo"
	"log"
	"runtime"
	"sync"
)

type config struct {
	Address    *string
	Database   *string
	DBUser     *string
	DBPassword *string
	DBPort     *int
	HttpPort   *int
}

func (c *config) Guard() error {
	if c.Database == nil {
		return errors.New("Database param cannot be empty.")
	}
	if c.DBUser == nil {
		return errors.New("DB user param cannot be empty.")
	}
	return nil
}

func (c *config) GetDBsourceName() string {
	return fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true", *c.DBUser, *c.DBPassword, *c.Address, *c.DBPort, *c.Database)
}

// Launches a command in another goroutine and increments the waitgroup counter
// when starting, and decrements it when the process finish.
func launchInBackground(function func(), wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		function()
		wg.Done()
	}()
}

func openDBConnection(dsn string) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Error testing the connection <%s>", err)
	}
	return db
}

func createHttpServer(port int) *HTTPServer.Server {
	router := HTTPServer.NewRouter()
	protocolAdapter := HTTPServer.NewJSONAdapter()
	return HTTPServer.New(router, protocolAdapter, "", port)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	config := &config{}
	config.Address = flag.String("address", "127.0.0.1", "Address of the database server, like 127.0.0.1 or db.whatever.com")
	config.Database = flag.String("database", "", "The name of the MySql database to use.")
	config.DBUser = flag.String("user", "", "The username to use in the MySql connection.")
	config.DBPassword = flag.String("password", "", "The password to use in the MySql connection.")
	config.DBPort = flag.Int("port", 3306, "The port where MySql is listening.")
	config.HttpPort = flag.Int("http_port", 80, "The port where the server will accept incoming http connections.")

	flag.Parse()
	if err := config.Guard(); err != nil {
		log.Fatal(err.Error() + " Use -help param to get more info.")
	}

	db := openDBConnection(config.GetDBsourceName())
	defer db.Close()
	log.Println("Connected to " + config.GetDBsourceName())

	httpServer := createHttpServer(*config.HttpPort)

	httpServer.RegisterEndpoint("/ping", command.NewPing(), "GET")
	httpServer.RegisterEndpoint("/track/store", command.NewSaveTrack(repo.NewTrackRepoAsync(repo.NewTrackRepoMYSQL(db), 10000, 20)), "POST")

	// Server will run in his own goroutine. We need to wait for it to finish
	wg := &sync.WaitGroup{}
	launchInBackground(httpServer.Run, wg)
	wg.Wait()
}
