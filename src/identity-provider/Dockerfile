# syntax=docker/dockerfile:1

FROM golang:1.17

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p logs temp

RUN go build -o ./app

EXPOSE 8040
ENV PORT=8040

ARG OKTA_CLIENT_SECRET
ENV OKTA_CLIENT_SECRET=$OKTA_CLIENT_SECRET

ARG OKTA_SSWS_TOKEN
ENV OKTA_SSWS_TOKEN=$OKTA_SSWS_TOKEN

ENTRYPOINT [ "./app" ]
