FROM golang:1.27.3-alpine

WORKDIR /workspace
COPY trusted/ .

# Precompile the test binary to eliminate runtime compilation overhead
# and allow execution on read-only root filesystems without GOCACHE.
RUN CGO_ENABLED=0 go test -c -o /usr/local/bin/trusted.test .

# Compile the standard test2json converter tool from Go toolchain
RUN go build -o /usr/local/bin/test2json cmd/test2json

# SC-07: output limiter that caps test output and keeps the verdict
COPY outcap/ /src/outcap/
RUN cd /src/outcap && CGO_ENABLED=0 go build -o /usr/local/bin/outcap .

CMD ["/usr/local/bin/test2json", "-t", "/usr/local/bin/trusted.test", "-test.v=test2json", "-test.run", "^TestAddSuccess$"]