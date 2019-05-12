FROM golang:1.11-alpine AS builder

# Install some dependencies needed to build the project
RUN apk add bash git

# create a working directory
WORKDIR /ffx/github.com/pamelag/blue

# Force the go compiler to use modules
ENV GO111MODULE=on

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .

RUN go mod download


# add source code
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o blue .

FROM scratch
WORKDIR /ffx
COPY --from=builder /ffx/github.com/pamelag/blue/ .
EXPOSE 8080
ENTRYPOINT [ "./blue" ]
