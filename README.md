# GHVisual

## Install & Run

#### 1. Docker

##### Prerequisites:

* [Docker](https://www.docker.com/get-started)
* [Github Personal Access Token](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/)

##### Steps

* Clone the repo.
* Open the file:
```
ghvisual/config/config.json
```
* Change the placeholder values to your own details.
* To build and run the Docker image enter the following commands with sudo privileges:

```
bash docker_build.sh
bash docker_run.sh
``` 
* With the container running, navigate to localhost:8080 in your browser to see the program.

#### 2. No Docker

##### Prerequisites:

* [Golang](https://golang.org/doc/install)
* [Govendor](https://github.com/kardianos/govendor)
* [Github Personal Access Token](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/)

##### Steps

* Clone the repo and make sure the cloned repo is in your GOPATH.
* Open the file:
```
ghvisual/config/config.json
```
* Change the placeholder values to your own details.
* From the root directory run:
```
govendor sync
```
* This will install the dependencies. You will have to ensure your GOPATH and PATH environment variables are set up correctly. If govendor isn't working for you, enter the following commands before trying govendor again:
```
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```
* Execute the run script from the root directory:
```
bash run.sh
```
* Navigate to localhost:8080 in your browser to see the program.