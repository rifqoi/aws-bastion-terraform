FROM golang:alpine as builder

WORKDIR /app

COPY . . 

# manage dependencies
RUN go mod download

RUN  ["go", "build", "-o", "main", "."]

FROM alpine:latest as main

WORKDIR /app

EXPOSE 8000
COPY --from=builder /app/main ./main

CMD [ "/app/main" ]
