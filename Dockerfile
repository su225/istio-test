FROM ubuntu:xenial

RUN mkdir -p /node
WORKDIR /node
COPY ./istio-test .
COPY ./init.sh .

ENV DOWNSTREAM=""
EXPOSE 8888
ENTRYPOINT ./init.sh
