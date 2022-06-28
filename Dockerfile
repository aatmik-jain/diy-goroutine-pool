ARG SERVICE_NAME=diy_goroutine_pool

FROM golang:latest as builder

ARG SERVICE_NAME
ENV SERVICE_NAME=$SERVICE_NAME

WORKDIR /Documents/DIY/${SERVICE_NAME}

COPY . .
CMD ["go", "test", "."]