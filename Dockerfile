FROM golang

MAINTAINER dreamans@163.com

RUN  mkdir -p /usr/local/src && \
     mkdir /usr/local/syncd

COPY . /usr/local/src

WORKDIR /usr/local/src

RUN make && \
    mv output/* /usr/local/syncd

EXPOSE 8878

CMD ["/bin/bash"]
