FROM golang:latest

COPY src/ /app
WORKDIR /app

RUN chmod +x /app/order && mkdir -p /app/outputs && mkdir -p /app/outlogs && go mod download

ENTRYPOINT ["bash","/app/order"]