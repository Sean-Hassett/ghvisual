# GHVisual

![alt text](images/screen.png)

## Overview

View data about the activity of a GitHub User.

The centre circle represents the User with each of the circles branching out from the User representing a repository that they own.
The size of the circle scales with the size of the repo. The color of the circle is determined by how recently the repo was updated. Repos which haven't been updated in a long time will be a lighter shade while repos which have been updated in the past few days will be dark.

To the left are two bar charts.
The top chart shows the breakdown of all a User's commits based on the day of the week the commit was made.
The bottom chart shows a breakdown based on the time of day the commit was made.
The intent is to reveal patterns in a User's activity around the day of the week and the time of day. 

For time of day, the definitions are as follows:
* Morning: 6:00am - 11:59am
* Afternoon: 12:00pm - 5:59pm
* Evening: 6:00pm - 11:59pm
* Night: 12:00am - 5:59am

The data is fetched live so it will be up to date but only public repositories and commit history can be accessed.
It is possible to include the data from private activity if your access token has access to your private repositories but it's a little awkward.
You need to change line 61 of ghvisual/retrieve.go from:
```
repos, _, err := client.Repositories.List(ctx, config.Username, nil)
```
to:
```
repos, _, err := client.Repositories.List(ctx, "", nil)
```

This means it will default to the User that the access token is associated with and will include private activity if the token has access.
Changing the Username in the config to an empty string won't work since that value is used in other places and  setting the Username to the account that the token is associated with won't pull in private activity even if the token has access to it so for now the awkward solution is the only one.

## To Do

#### Timelapse

Originally I wanted to do a timelapse starting on the day the User created their GitHub account and proceeding to the present day.
The idea was that you could see your work on GitHub grow day by day, with the repos appearing, disappearing and changing size as time went by.

I think this would look really interesting and would also reveal patterns of activity.
For example if you were rewriting a lot of code you would see your commits going up but your repos not changing size very much.

The main problem I faced with this was the difficulty in getting the size of a repo at the time of each commit.
The only way I could see was to iterate through the file structure for each commmit counting up the size of each individual file and this would probably result in hitting the rate limit very fast.

#### Weights

I also had an idea to attempt to model how many lines of code it took to solve a problem in different languages and use this to weight the size of the repos depending on what language they were predominantly written in.

My plan for this was to take public repos where people solve problems from sites such as Hackerrank and Project Euler. With these repos, you would have examples of people solving the very same problems using different languages and you could compare to see the difference in size of the solutions.
This doesn't factor in different approaches to problem solving and different coding styles but I thought it would be interesting to see even a crude comparison.

I managed to locate some repos which would have served the purpose well but since it would have been fairly tedious and uninteresting to parse through the different naming conventions used to automatically gather the data, I ended up leaving it to one side and not coming back to it.

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
* Enter your access token and the username you want to look at.
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

* Clone the repo and make sure the cloned repo is in your GOPATH (which is usually set to $HOME/go).
* Open the file:
```
ghvisual/config/config.json
```
* Enter your access token and the username you want to look at.
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