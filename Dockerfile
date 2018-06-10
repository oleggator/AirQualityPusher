FROM golang:latest AS build

RUN  mkdir -p /go/src \
  && mkdir -p /go/bin \
  && mkdir -p /go/pkg

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN mkdir -p $GOPATH/src/pusher
ADD . $GOPATH/src/pusher

WORKDIR $GOPATH/src/pusher
RUN go build -a -installsuffix cgo -o /bin/pusher .


FROM scratch

COPY --from=build /bin/pusher /bin/pusher
COPY --from=build /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip

ENTRYPOINT ["/bin/pusher"]
