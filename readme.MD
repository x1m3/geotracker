A sample code to demonstrate an HTTP server that accepts json requests
with a {GPS coordinate, tracker id} and stores asynchronously in a mysql database.

## How to install

1. Ensure that your GOPATH variable is set. If it is not defined, run "export GOPATH=${HOME}/mygo"

2. Check your go version with "go version" command. This project has been developed with go 1.9

3. Run "go get -u github.com/x1m3/geotracker"

4. cd $GOPATH/src/github.com/x1m3/geotracker

5. Run "go build" and you'll get an executable file called geotracker

## Configure MySQL database

Create a Mysql database, and add a table like this.

```
CREATE TABLE `track` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `driver_id` int(11) NOT NULL,
  `lat` double NOT NULL,
  `lon` double NOT NULL,
  `created_on` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`id`),
  KEY `track_by_driver_idx` (`driver_id`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8;
```

## Run the server

Execute "geotracker help" to get some information. You need to pass in command line the database parameters and the port to listen.

Database host, database port and http port are optional and defaults to 127.0.0.1, 3306 and 8080.

```
xime@superfrontserver:~/go/src/github.com/x1m3/geotracker$ ./geotracker -help
Usage of ./geotracker:
  -address string
        Address of the database server, like 127.0.0.1 or db.whatever.com (default "127.0.0.1")
  -database string
        The name of the MySql database to use.
  -http_port int
        The port where the server will accept incoming http connections. (default 80)
  -password string
        The password to use in the MySql connection.
  -port int
        The port where MySql is listening. (default 3306)
  -user string
        The username to use in the MySql connection.
````

## Testing the server

Open a browser, etc.. and call the endpoint /ping. Assuming that you are running on localhost with the default port 8080, visit http://127.0.0.1/ping.

You will see a "pong" response.

## Sending data to be stored in the database

Send a POST request. (CAUTION: GET method will fail) to the endpoint /track/store with a json body like this

```
{ "latitude" : 22, "longitude" : 3.34324, "driver_id": 2 }
```


You will get a 200 response with a body message "OK". Your data should be on the database ;-)

## How it works

Code is more or less self explanatory. The interesting point is that in stores tracks asynchronously. The client gets and immediate response, the track is enqueued and is stored later via a pool of workers in the database.

The most interesting part of all of this is the file repo/trackRepoAsync.go. Looking at the constructor we see that it accepts 3 params: An interface to a repo, the length of the queue and the number of workers.

```
func NewTrackRepoAsync(r Track, queueSize int, nWorkers int) *TrackRepoAsync {
```

The queue is implemented with a buffered channel. Writing to this channel will not block until the channel is full.

The third parameter is the number of workers. The repo will launch some goroutines that will be blocked, listening to any order from the channel. This adds pressure to the database, allowing to send simultaneous data via multiple connections. Take care to do not set this number very high, because it can kill your database. Right now, it is configured on main.go.

```
httpServer.RegisterEndpoint("/track/store", command.NewSaveTrack(repo.NewTrackRepoAsync(repo.NewTrackRepoMYSQL(db), 10000, 20)), "POST")
```

If, at some point, the queue is full, then any call to the repo will start to block, and the system will increase the response time to the client. This system is only capable of asking more work to the database and enqueue work that could be done at some point later.

The other interesting point is the concept of metarepo. Thanks to the use of interfaces, we can use another repo to store data in other places, and we can use it directly or combine with the AsyncRepo.

## Things to be done.

- Authorize.
- Document endpoints.
- Graceful shutdown. (Current code could lose data when shuting down)
- Better loging.
- Adding metrics.
- Etc...

