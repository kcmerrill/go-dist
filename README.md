## go-dist
Go binary distribution made easy

## Installation
`go get -u github.com/go-dist`

## What is it
I got tired of making "releases" in github, especially for my tiny projects. go-dist will build any github repo's binaries via a website and allow you to download the binaries without you needing to build them yourself. Simply link to the binaries and go-dist will rebuild the binary for you and serve them up to your users.

## Installation via Docker
`docker run -d -P --restart=always --name go-dist kcmerrill/go-dist`

## Binaries
![Mac OSX](https://raw.githubusercontent.com/kcmerrill/go-dist/master/assets/apple_logo.png "Mac OSX") [386](http://go-dist.kcmerrill.com/kcmerrill/go-dist/mac/386) | [amd64](http://go-dist.kcmerrill.com/kcmerrill/go-dist/mac/amd64)

![Windows](https://raw.githubusercontent.com/kcmerrill/go-dist/master/assets/windows_logo.png "Windows") [386](http://go-dist.kcmerrill.com/kcmerrill/go-dist/windows/386) | [amd64](http://go-dist.kcmerrill.com/kcmerrill/go-dist/windows/amd64)

![Linux](https://raw.githubusercontent.com/kcmerrill/go-dist/master/assets/linux_logo.png "Linux") [386](http://go-dist.kcmerrill.com/kcmerrill/go-dist/linux/386) | [amd64](http://go-dist.kcmerrill.com/kcmerrill/go-dist/linux/amd64)

## Demo
The binary section is this project in action!

## Getting Started
Setup go-dist by seeing installation/binaries above, or check out this project running at: [https://go-dist.kcmerrill.com](https://go-dist.kcmerrill.com)
1. To get an easy copy/paste markdown `https://go-dist.kcmerrill.com/<github_username>/<project_name>`
  1. Example: [https://go-dist.kcmerrill.com/kcmerrill/go-dist](https://go-dist.kcmerrill.com/kcmerrill/go-dist)
  1. Example: [https://go-dist.kcmerrill.com/kcmerrill/alfred](https://go-dist.kcmerrill.com/kcmerrill/alfred)
2. Or ... you can manually create the links like so: `https://go-dist.kcmerrill.com/<github_username>/<project_name>/<OS:mac|linux|windows>/<arch_type:amd64|386|arm>`

Of course, I'm using `https://go-dist.kcmerrill.com` in these examples, but you can substitute wherever `go-dist` is running

Using the `--cache <string golang time.duration, 10s, 60m, 1h ... >` will invalidate the cache. This means that the binary will be rebuilt. You can use webhooks to invalidate the cache. Simply setup a webhok to `send everything` to `https://go-dist.kcmerrill.com/<github_username>/<github_project>`. Anytime anything is merged into master we'll invalidate the cache so the next user gets a fresh copy of your binary.

## Limitations & Known Issues
The first user gets to "warm" up the cache. This is intended.

Also, there are quite a few of [known limitations](https://github.com/golang/go/issues/6376) when it comes to cross compiling. If you are noticing issues with your binaries, chances are likely there are open/stale github issues in the golang issue tracker.

## How
When a user clicks on the link, if it's been over a half hour since the project was built, or if a binary/project doesn't exist, `go get -u <project>`. Then, using the great work over at [mitchellh/gox](https://github.com/mitchellh/gox), generate the binary on the fly. Until webhooks are integreated, the first person to get a non-cached version it will be a bit slower for.

## TODO
* More error checking
* Better user interface
* Better looking github markdown
* Currently only supports github public projects, enable private/non-github repos
