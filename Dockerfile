FROM zhuzhengqian123/golang

LABEL maintainer="zhuzhengqian@coding.net"

COPY . /app
WORKDIR /app
ENV GOPROXY=https://goproxy.io
RUN go build -v -o code-watcher
RUN cp ./code-watcher /bin/code-watcher && chmod a+x /bin/code-watcher
COPY ./supervisord/code-watcher.conf /etc/supervisor/conf.d/