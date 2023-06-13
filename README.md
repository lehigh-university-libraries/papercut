# Papercut CLI

Command line utility to help fetch papers from various sources.

## Install

### Homebrew

You can install papercut using homebrew

```
brew tap lehigh-university-libraries/papercut
brew install papercut
```

### Download Binary

Instead of homebrew, you can download a binary for your system from [the latest release](https://github.com/lehigh-university-libraries/homebrew-papercut/releases/latest)

Then put the binary in a directory that is in your `$PATH`

## Usage

TODO

## Updating

### Homebrew

If homebrew was used, you can simply upgrade the homebrew formulae for papercut

```
brew update && brew upgrade papercut
```

### Download Binary

If the binary was downloaded and added to the `$PATH` updating papercut could look as follows. Requires [gh](https://cli.github.com/manual/installation) and `tar`

```
# update for your architecture
ARCH="papercut_Linux_x86_64.tar.gz"
TAG=$(gh release list --exclude-pre-releases --exclude-drafts --limit 1 --repo lehigh-university-libraries/homebrew-papercut | awk '{print $3}')
gh release download $TAG --repo lehigh-university-libraries/homebrew-papercut --pattern $ARCH
tar -zxvf $ARCH
mv papercut /directory/in/path/binary/was/placed
rm $ARCH
```
