FROM golang:1.10

WORKDIR ${GOPATH}/src/github.com/hand-writing-authentication-team/credential-store

ADD . .

RUN go get -u github.com/kardianos/govendor
RUN govendor sync

CMD ["go", "run", "server.go"]