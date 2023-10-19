FROM golang:alpine
RUN mkdir /app 
ADD . /app/
WORKDIR /app 
RUN go get github.com/lib/pq && go build -o main .
CMD ["./main"]