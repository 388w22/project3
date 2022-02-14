FROM eecs388/project3-runner

RUN apk add --no-cache tcpdump
COPY network/ network/
RUN go build -o ./network/client/client ./network/client
RUN go build -o ./network/dns/dns_server ./network/dns
RUN go build -o ./network/http/http_server ./network/http

COPY mitm/ mitm/
RUN go build -o ./mitm/mitm ./mitm
