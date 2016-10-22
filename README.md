## go-dist
Go binary distribution made easy

## Installation
`go get -u github.com/go-dist`

## Installation via Docker
`docker run -d -P --name go-dist kcmerrill/go-dist`

## Demo
The binary section is this project in action!

## Getting Started
Setup go-dist by seeing installation/binaries above, or check out this project running at: [https://go-dist.kcmerrill.com](https://go-dist.kcmerrill.com)
1. To get an easy copy/paste markdown `https://go-dist.kcmerrill.com/<github_username>/<project_name>`
  1. Example: [https://go-dist.kcmerrill.com/kcmerrill/go-dist](https://go-dist.kcmerrill.com/kcmerrill/go-dist)
  1. Example: [https://go-dist.kcmerrill.com/kcmerrill/alfred](https://go-dist.kcmerrill.com/kcmerrill/alfred)
2. Or ... you can manually create the links like so: `https://go-dist.kcmerrill.com/<github_username>/<project_name>/<OS:mac|linux|windows>/<arch_type:amd64|386|arm>`

Of course, I'm using `https://go-dist.kcmerrill.com` in these examples, but you can substitute wherever `go-dist` is running

## Limitations & Known Issues
Until some of the features outlined below in my todo list, there are of course some limitations. I wouldn't recommend using this for anything that gets a ton of downloads a day, but for the hobbyist it should be just fine. Once webhooks are integrated, I'd then recommend it for larger projects, but not until then.

Also, there are quite a few of [known limitations](https://github.com/golang/go/issues/6376) when it comes to cross compiling. If you are noticing issues with your binaries, chances are likely there are open/stale github issues in the golang issue tracker.

## How
Every half hour, or if a binary/project doesn't exist, `go get -u <project>`. Then, using the great work over at [mitchellh/gox](https://github.com/mitchellh/gox), generate the binary on the fly.

## TODO
* Add webhook capabilities when merged into master, invalidate cache and automagically build everything(preferred)
* More error checking
* Better user interface
* Better looking github markdown
* Currently only supports github public projects, enable private/non-github repos
