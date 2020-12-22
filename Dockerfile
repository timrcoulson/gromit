FROM golang

ENV GO111MODULE on

# Printing deps
RUN apt update && apt install --no-install-recommends -y enscript cups cups-bsd ca-certificates bash jq && rm -rf /var/lib/apt/lists/*
RUN cp /etc/cups/cupsd.conf /etc/cups/cupsd.conf.original
RUN chmod a-w /etc/cups/cupsd.conf.original

WORKDIR /src
COPY go.* ./
RUN go mod download
COPY . .
ENV PORT 8080
CMD ["make", "run"]
EXPOSE 8080
EXPOSE 631