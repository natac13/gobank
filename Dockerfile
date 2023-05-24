FROM golang:1.16-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/bin/ ./cmd/...

# Path: Dockerfile

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/ /app/bin/

CMD ["/app/bin/"]

## 3. Build the image

# ```bash
# docker build -t go-docker .
# ```

## 4. Run the image

# ```bash
# docker run -it --rm go-docker
# ```
