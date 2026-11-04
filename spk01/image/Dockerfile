FROM golang:1.27.3-alpine

WORKDIR /workspace
COPY trusted/ .

# Warm up the build cache so runtime test execution is fast
RUN go test -run TestAddSuccess ./...

CMD ["go", "test", "-json", "-run", "^TestAddSuccess$", "./..."]