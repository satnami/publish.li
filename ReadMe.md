# publish.li #

"Publish Your Articles Quickly and Easily."

A simple publishing site much like [Medium](http://medium.com/) or [telegra.ph](http://telegra.ph/). No need to sign up
or log in - just type your content and hit publish!

To edit, you will need the Unique ID shown at the top of the page.

Please see [alternatives to publish.li](https://alternativeto.net/software/publish-li/).

## Deploying to your Own Server ##

This project uses the excellent [gb](https://getgb.io/) as it's build tool. Here's a quickstart:

1. clone the repo, with `git clone git@github.com:appsattic/publish.li.git`
2. run `gb build` in the project root
3. run `./bin/publish` in the project root

When you run it, you must provide two environment variables:

* `PORT` - the local port you want to listen on, e.g. `8000`
* `BASE_URL` - how your server looks from the outside world (e.g. `https://publish.li` or `http://localhost:8000`)

Run the `./bin/publish` executable from the project root, so that the program can load up the templates and serve the
static pages. It outputs to both STDIN and STDERR, so it's up to you to redirect those where appropriate.

You might use a command like `BASE_URL=http://localhost:8080 PORT=8080 ./bin/server` in development.

## The DataStore ##

Since publish.li uses the BoltDB embedded datastore, this project won't run on PaaS solutions like Heroku or
OpenShift. There are very few operations which use the datastore (essentially get and put) so adding another backend
such as MongoDB or Postgres should be pretty easy. The filesystem would also be easy for self-hosted on your own server
but would have the same problems as the BoltDB backend.

Once anyone shows interest in another backend, we'll do a small refactor to use an interface in the application instead
of a concrete type. Until then though, let's just leave it as-is.

## Author ##

[Andrew Chilton](https://chilts.org), [@andychilton](https://twitter.com/andychilton).

For [AppsAttic](https://appsattic.com), [@AppsAttic](https://twitter.com/AppsAttic).

## License ##

This project is free software and can be forked, downloaded, used, and shared.

[AGPLv3](https://www.gnu.org/licenses/agpl-3.0.txt).

(Ends)
