
FROM golang:1.18-alpine AS build


WORKDIR /app


COPY src/go.mod src/go.sum ./


RUN go mod download


COPY ./src .


RUN go build -o myapp ./cmd/main.go


FROM alpine:latest


EXPOSE 8080


COPY --from=build /app/myapp /usr/local/bin/myapp


ENTRYPOINT ["myapp"]
