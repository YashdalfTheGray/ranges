FROM golang:buster as builder

WORKDIR /app
COPY . .
RUN go build

# FROM scratch

# COPY --from=builder /app/ranges /bin/
EXPOSE 8080-8085
ENTRYPOINT [ "./ranges" ]