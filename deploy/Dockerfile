#FROM golang:1.17
FROM frolvlad/alpine-glibc:glibc-2.33

ARG docker_version=0.0.1

COPY ./build/release/nazari-broker_$docker_version /bin/broker

CMD ["/bin/broker"]