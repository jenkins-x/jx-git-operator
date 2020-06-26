FROM gcr.io/jenkinsxio-labs-private/jxl-base:0.0.53

ENTRYPOINT ["jx-git-operator"]

COPY ./build/linux/jx-git-operator /usr/bin/jx-git-operator