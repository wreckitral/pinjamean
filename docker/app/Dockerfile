FROM golang:1.27

RUN go install github.com/cespare/reflex@v0.3.2
COPY reflex.conf /

COPY start.sh /
COPY entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]
