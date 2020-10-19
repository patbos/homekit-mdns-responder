FROM golang:1 as build-env
# All these steps will be cached
RUN mkdir /build
WORKDIR /build
COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
# COPY the source code as the last step
COPY . .

# Build the binary
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o /go/bin/hap-mdns

FROM scratch
COPY --from=build-env /go/bin/hap-mdns /hap-mdns

COPY --from=build-env /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Europe/Stockholm

ENTRYPOINT ["/hap-mdns"]
