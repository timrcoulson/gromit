FROM golang

ENV GO111MODULE on

# Printing deps
COPY provision.sh .
RUN ./provision.sh

WORKDIR /src
COPY go.* ./
RUN go mod download

# Add a new user "john" with user id 8877
RUN useradd -u 8877 john
# Change to non-root privilege
USER john

COPY . .
ENV PORT 8080
CMD ["make", "run"]
EXPOSE 8080
EXPOSE 631

