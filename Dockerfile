FROM gcr.io/jenkinsxio/jx-cli-base:0.0.32

ENTRYPOINT ["jx-git-operator"]

COPY ./build/linux/jx-git-operator /usr/bin/jx-git-operator