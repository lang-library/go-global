#! /usr/bin/env bash
set -uvx
set -e
cwd=`pwd`
version=0.1.4
git add .
git commit -m"Release v$version"
git tag -a v$version -mv$version
git push origin v$version
git push
git remote -v
gh auth login --with-token < ~/settings/github-all-tokne.txt
#gh release create v$version --repo lang-library/go-global --generate-notes --target main
gh release create v$version --generate-notes --target main
