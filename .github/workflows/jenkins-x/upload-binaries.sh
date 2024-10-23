#!/bin/bash

echo "HOME is $HOME"
echo current git configuration

# See https://github.com/actions/checkout/issues/766
git config --global --add safe.directory "$GITHUB_WORKSPACE"

git config --global --get user.name
git config --global --get user.email

echo "setting git user"

git config --global user.name jenkins-x-bot-test
git config --global user.email "jenkins-x@googlegroups.com"

git add * || true
git commit -a -m "chore: release $VERSION" --allow-empty
git tag -fa v$VERSION -m "Release version $VERSION"
git push origin v$VERSION

export BRANCH=$(git rev-parse --abbrev-ref HEAD)
export BUILDDATE=$(date)
export REV=$(git rev-parse HEAD)
export GOVERSION="1.23.2"
export ROOTPACKAGE="github.com/$REPOSITORY"

# Install syft
curl -sSfL https://raw.githubusercontent.com/anchore/syft/main/install.sh | \
sh -s -- -b /usr/local/bin v0.54.0
chmod +x /usr/local/bin/syft


goreleaser release
