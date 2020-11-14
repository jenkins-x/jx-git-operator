FROM alpine

ARG BUILD_DATE
ARG VERSION
ARG REVISION
ARG TARGETARCH
ARG TARGETOS

LABEL maintainer="jenkins-x"

RUN addgroup -S app \
    && adduser -S -g app app \
    && apk --no-cache add \
    ca-certificates curl git netcat-openbsd

# kubectl
ENV KUBECTL_VERSION 1.16.15

# see https://docs.docker.com/engine/reference/builder/#automatic-platform-args-in-the-global-scope
RUN echo using kubectl version ${KUBECTL_VERSION} and OS ${TARGETOS} arch ${TARGETARCH} && \
  curl -LO  https://storage.googleapis.com/kubernetes-release/release/v${KUBECTL_VERSION}/bin/${TARGETOS}/${TARGETARCH}/kubectl && \
  mv kubectl /usr/bin/kubectl && \
  chmod +x /usr/bin/kubectl

RUN echo using jx-git-operator version ${VERSION} and OS ${TARGETOS} arch ${TARGETARCH} && \
  mkdir -p /home/.jx3 && \
  curl -L https://github.com/jenkins-x/jx-git-operator/releases/download/v${VERSION}/jx-git-operator-${TARGETOS}-${TARGETARCH}.tar.gz | tar xzv && \
  mv jx-git-operator /usr/bin

WORKDIR /home/app

USER app

ENTRYPOINT ["jx-git-operator"]

