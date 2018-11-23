# GHVisual

## Install & Run

#### 1. Docker

##### Prerequisites:

* [Docker](https://www.docker.com/get-started)
* [Github Personal Access Token](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/)

##### Steps

* Clone the repo
* Open the file:
```
ghvisual/config/config.json
```
* Change the placeholder values to your own details
* Run the following commands:

```
docker build -t ghvisual .
docker run -p 8080:8080 ghvisual
```
* Navigate to localhost:8080 in your browser to see the program

#### 2. No Docker

##### Prerequisites:

* [Golang](https://golang.org/doc/install)
* [Govendor](https://github.com/kardianos/govendor)
* [Github Personal Access Token](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/)

##### Steps

* Clone the repo and make sure the cloned repo is in your GOPATH
* From the root directory run:
```
govendor sync
```
* This will install the dependencies
* Run the following commands:
```
cd ghvisual
go build -o main . && cd ..
```
* To run the program, from the root directory run:
```
./ghvisual/main
```
* Navigate to localhost:8080 in your browser to see the program