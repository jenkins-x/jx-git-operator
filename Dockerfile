FROM golang:1.15 as builder

ARG REVISION

RUN mkdir -p /jx-git-operator/

WORKDIR /jx-git-operator

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -ldflags "-s -w \
    -X github.com/jenkins-x/jx-git-operator/pkg/version.REVISION=${REVISION}" \
    -a -o bin/jx-git-operator main.go

FROM alpine

ARG BUILD_DATE
ARG VERSION
ARG REVISION

LABEL maintainer="jenkins-x"

# kubectl
ENV KUBECTL_VERSION 1.16.15
RUN curl -LO  https://storage.googleapis.com/kubernetes-release/release/v${KUBECTL_VERSION}/bin/linux/amd64/kubectl && \
  mv kubectl /usr/bin/kubectl && \
  chmod +x /usr/bin/kubectl

RUN addgroup -S app \
    && adduser -S -g app app \
    && apk --no-cache add \
    ca-certificates curl git netcat-openbsd

WORKDIR /home/app

COPY --from=builder /jx-git-operator/bin/jx-git-operator /usr/bin/jx-git-operator

USER app

ENTRYPOINT ["jx-git-operator"]

