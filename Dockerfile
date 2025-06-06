# Build stage 
FROM golang:1.24.3-alpine3.21 AS builder 
# set working directory . as /app 
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz | tar xvz

# Run stage 
FROM alpine:3.21
# set working directory . as /app 
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
RUN chmod +x start.sh
RUN chmod +x wait-for.sh
COPY db/migration ./migration

# which port are intented to be published 
EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
