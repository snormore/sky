FROM snormore/dependencies

MAINTAINER SkyDB skydb.io 

ENV GOPATH /go
ENV GOBIN /go/bin
ENV SKY_OWNER_PATH /go/src/github.com/snormore

RUN mkdir -p $GOBIN

RUN mkdir -p $SKY_OWNER_PATH
ADD . $SKY_OWNER_PATH/sky

RUN cd $SKY_OWNER_PATH/sky && go get && go build -a -o /usr/local/bin/skyd

CMD ["-port 8589"]

ENTRYPOINT /usr/local/bin/skyd

