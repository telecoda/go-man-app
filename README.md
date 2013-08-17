#go-man

Pacman type game built as a RESTful API using go-lang

This project is my first venture into using go.

- A. I normally spend my day developing API's in java / python.

- B. I normally spend my commute time trying to throw together some game code for fun.

- A + B = C. Here is my playground project to learn about the go language and have a little fun at the same time.

Inspired by the idea of https://speakerdeck.com/christkv/mongoman-a-nodejs-powered-pacman-clone I saw demoed at

2012 MongoDB conference in London.

##Running all the tests

The project contains a number of unit and functional test which can be run with following command

    go test ./...

##Running all the server

To compile the go-man server type:

    go build

To run the go-man server type:

    ./go-man-app

There is a javascript client located at:-

http://github.com/telecoda/go-man-javascript-client.git

This can be used to play the game locally.

Or there is a demo version of the game hosted on heroku and google app engine.

Client: http://go-man-client.heroku.com

API: http://go-man-app.appspot.com