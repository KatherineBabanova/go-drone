FROM golang as build-stage

WORKDIR /build
COPY main.go ./
COPY go.* ./
RUN go build -o server

FROM ubuntu as deploy-stage
WORKDIR /app
COPY --from=build-stage /build/server /app/
#EXPOSE 8090

CMD ["./server"]