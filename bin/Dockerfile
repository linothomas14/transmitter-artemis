FROM alpine:3.18
WORKDIR /app
RUN apk add --no-cache tzdata
ENV TZ=Asia/Jakarta
COPY transmitter-artemis .
COPY config.yaml .
CMD ["/app/transmitter-artemis"]