# Usage
#   bash scripts/tag.sh v0.3.2

set -x

REMOTE=https://github.com/suzuki-shunsuke/candy

BRANCH=`git branch | grep "^\* " | sed -e "s/^\* \(.*\)/\1/"`
if [ "$BRANCH" != "master" ]; then
  read -p "The current branch isn't master but $BRANCH. Are you ok? (y/n)" YN
  if [ "${YN}" != "y" ]; then
    echo "cancel to release"
    exit 0
  fi
fi

TAG=$1
VERSION=${TAG#v}

if [ "$TAG" = "$VERSION" ]; then
  echo "the tag must start with 'v'" >&2
  exit 1
fi

cd `dirname $0`/..

VERSION_FILE=pkg/domain/version.go

echo "create $VERSION_FILE"
cat << EOS > $VERSION_FILE || exit 1
package domain

// Don't edit this file.
// This file is generated by the release command.

// Version is the candy's version.
const Version = "$VERSION"
EOS

git add $VERSION_FILE || exit 1
git commit -m "build: update version to $TAG" || exit 1
git tag $TAG || exit 1
git push $REMOTE $BRANCH || exit 1
git push $REMOTE $TAG
