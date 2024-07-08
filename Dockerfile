FROM golang:1.12.5-stretch
LABEL maintainer="Nader Carun <comebacknader@gmail.com>"

WORKDIR ~

ENV WINK_DB_U wink-db-username 
ENV WINK_DB_P wink-db-password 
ENV WINK_DB_HOST wink-db-container 
ENV WINK_DB_NAME wink-db-name 
ENV WINK_PATH /go/src/github.com/comebacknader/wink/ 
ENV WINK_ENVIRON local 

RUN go get github.com/lib/pq \
	 && go get github.com/gorilla/websocket \
	 && go get github.com/julienschmidt/httprouter \
	 && go get github.com/microcosm-cc/bluemonday \
	 && go get github.com/badoux/checkmail \
	 && go get golang.org/x/crypto/bcrypt \
	 && go get github.com/satori/go.uuid \
	 && go get github.com/pilu/fresh

RUN mkdir -p /go/src/github.com/comebacknader/wink/ 

WORKDIR /root 

RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.34.0/install.sh | bash \
	&& source .profile \ 
	&& nvm install 12.4.0 \
	&& npm install -g sass@1.21.0 

WORKDIR /go/src/github.com/comebacknader/wink/ 

COPY . .

EXPOSE 80 8080

CMD /go/bin/fresh 
