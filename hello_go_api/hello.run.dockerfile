# BUILD:
# docker build --force-rm=true -f hello.run.dockerfile -t hello .

# RUN:
# docker run -it -p 8000:8000 hello

FROM golang as builder

RUN mkdir /app
RUN mkdir /go/src/hello
ADD . /go/src/hello
WORKDIR /go/src/hello

# Go dep
RUN go get -d ./...

# Build a standalone binary
RUN ls && export GOPATH=$GOPATH:/go && set -ex && \
  CGO_ENABLED=0 go build \
        -tags netgo \
        -o /app/hello \
        -v -a \
        -ldflags '-extldflags "-static"' && \
  ls

# Create the second stage with a basic image.
# this will drop any previous
# stages (defined as `FROM <some_image> as <some_name>`)
# allowing us to start with a fat build image and end up with
# a very small runtime image.

FROM busybox

# add compiled binary
COPY --from=builder /app/hello /hello

# run
EXPOSE 8000
CMD ["/hello"]
