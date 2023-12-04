FROM golang:alpine

# Install bash
RUN apk --no-cache add bash

# Set the SESSION_SECRET_KEY argument using uuidgen
RUN echo "export SESSION_SECRET_KEY=\$(uuidgen)" > /etc/profile.d/session_secret.sh && \
    chmod +x /etc/profile.d/session_secret.sh

RUN mkdir /app
WORKDIR /app
COPY . .

# Install dependencies and build the application
RUN go get github.com/gorilla/mux@v1.8.1 && \
    go get github.com/lib/pq && \
    go build -o main .

CMD ["./main"]