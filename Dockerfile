FROM golang

ADD . /go/src/github.com/archi-chester/skolkovo-test

RUN go get github.com/gorilla/mux

RUN go get github.com/sirupsen/logrus

RUN go get github.com/twinj/uuid

RUN go install github.com/archi-chester/skolkovo-test

RUN pwd

#EXPOSE 9003

RUN ls -l src/github.com/archi-chester/skolkovo-test

RUN chmod a+x src/github.com/archi-chester/skolkovo-test/skolkovo-test

WORKDIR src/github.com/archi-chester/skolkovo-test

ENTRYPOINT ["skolkovo-test"]

