FROM golang:1.21 AS builder

WORKDIR /home/builder
COPY . .
ENV CGO_ENABLED=0
RUN go build -o planner ./cmd/server


FROM alpine:latest AS runner

WORKDIR /home/planner
COPY --from=builder /home/builder/planner .

ENTRYPOINT ["/home/planner/planner"]
