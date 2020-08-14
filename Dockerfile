FROM gcr.io/jenkinsxio-labs-private/jx-cli-base:0.0.3

ENTRYPOINT ["jx-git-operator"]

COPY ./build/linux/jx-git-operator /usr/bin/jx-git-operator