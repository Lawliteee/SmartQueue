FROM alpine:latest
WORKDIR /app
COPY server .
COPY migrations/ migrations/
EXPOSE 8080
CMD ["./server"]
