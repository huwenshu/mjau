# Mjau

Mjau is a simple, fast, and flexible web font server; it efficiently delivers
CSS files containing `@font-face` rules via HTTP.

## Features

Mjau has the following features:

* Automatic generation of `@font-face` rules from web font and metadata files.
* Grouping of multiple font styles and weights in one response.
* Support for `EOT` and `WOFF` web font formats.
* Customizable CSS templates.
* Web font CSS embedding using `base64`-encoded data URIs.
* Customizable `Cache-Control: max-age` HTTP response header.
* Optional entity tags (`ETag`s) generation and validation.
* Optional HTTP response `gzip` compression.
* Optional Cross-Origin Resource Sharing (`CORS`).
* Whitelist-based HTTP referrer validation.
* Easy configuration through command-line flags.

## Drawbacks

Currently the web fonts are delivered by embedding them in a CSS file using
`base64`-encoded data URIs. Unfortunately this technique does not work at
all/correctly in Microsoft Internet Explorer prior to version 9.

Support for web fonts delivery through external file linking from the CSS file
is on the way.

## Installation

### Prerequisites

Mjau is written in the Go Programming Language and you will need a Go workspace
in order to build it. You can find instructions for downloading and installing
the Go compilers, tools, and libraries in the official Go [Getting Started][1]
guide.

Once you have a Go workspace, you must ensure that the `GOPATH` environment
variable is correctly set. For more details on how to set and use it please
read the [`GOPATH` environment variable][2] section from the official
documentation of the `go` command.

### Building

After having a Go workspace ready and correctly setting the `GOPATH`
environment variable, the compilation and installation process is fairly
straightforward.

Change directory to one of the `GOPATH` directories and issue the following
command:

	$ go get github.com/noll/mjau

The `go` command will automatically fetch the code and its dependencies, build
the program, and install it.

### Compatibility

Tests have been successfully ran using Go version `go1.0.2` on the following
operating systems:

* Arch Linux
* FreeBSD 9.0
* Mac OS X 10.8
* Microsoft Windows 7

## User Guide

### Prerequisites

Before proceeding, you must ensure that the `GOPATH/bin` directory, where the
Mjau executable resides, is present in the `PATH` environment variable.

Assuming you have only one Go workspace in your `GOPATH`, adding the binary
directory to your `PATH` can be achieved by issuing the following command:

	$ export PATH=$PATH:$GOPATH/bin

### Quickstart

First of all, in order to serve web fonts, Mjau requires a font library. Mjau
is distributed with a tiny font library, containing the Amaranth and Open Sans
font families, which resides in the `fonts` directory.

Another requirement is a set of template files, one template for each web font
format supported by the server. The default set of templates resides in the
`templates` directory.

The last requirement of the server is represented by a JSON-encoded whitelist
which enumerates all the domain names allowed to use the service. The default
whitelist resides in the root directory of the project and is called
`whitelist.json`. It enumerates only the `http://localhost/` domain name, so
this is the only domain name allowed to use the service.

#### Starting Mjau

At startup Mjau tries to find the font library, templates directory, and the
whitelist in the current directory. The fastest way to get it to work is to
start it from the source directory, where all the defaults reside:

	$ cd $GOPATH/src/github.com/noll/mjau
	$ mjau -b :8080

Mjau is now listening on all `IPv4` addresses available on the local machine,
on port `8080`.

#### Using Web Fonts

Once the server is running you are ready to start using web fonts in your HTML
pages. First, in the HTML document, you need to link to the CSS file containing
the `@font-face` rules:

	<link rel="stylesheet" href="http://localhost:8080/css/?family=Amaranth">

Second, in the stylesheet, you need to apply the requested `font-family` to the
appropriate selectors:

	p { font-family: "Amaranth", sans-serif }

If you refresh the page you should be able to see the web fonts in action.

### Configuration

Now that you've learned the basics, let's dive into more advanced Mjau
configuration options.

#### Font Library

The font library is a file system directory containing font families. Inside
the font library, each font family must be placed in its own directory. The
font family directory must contain all the subfamilies of the font family (in
one or more web font formats) and a JSON-encoded metadata file named
`metadata.json`.

The metadata file for a font family must enumerate all the font subfamilies and
web font formats which will be made available through Mjau. The unspecified
subfamilies and web font formats will remain private even if they are present
in the font family directory. Font families without metadata files are
completely ignored.

You can choose which font library to use with the `-l` command-line flag:

	$ mjau -l /path/to/font/library

An example font library is available in the `fonts` directory and may be used
as a starting point in building your own.

#### CSS Templates

Web fonts are served using CSS files containing `@font-face` rules. A set of
CSS templates, one template for each supported web font format, is used to
generate these files. Each template must be named after the web font format
for which the template is going to be used and must have the `.css.tmpl
extension.

You can choose which templates to use with the `-t` command-line flag:

	$ mjau -t /path/to/templates/directory

Example templates are provided in the `templates` directory and may be used as
a starting point in writing your own.

#### Whitelist

The whitelist is a JSON-encoded whitelist which enumerates all the domain names
allowed to use the service. When receiving a request the server verifies the
HTTP referrer provided in the request and delivers the requested web fonts only
to those domains which match one of the entries in the whitelist.

One trick is that you can allow any domain to use the service by specifying the
empty string in the domains list as in the following example:

	{ "domains": [ "" ] }

You can choose which whitelist to use with the `-w` command-line flag:

	$ mjau -w /path/to/whitelist.json

An example whitelist named `whitelist.json` is available in the root directory
of the project and may be used as a starting point in writing your own.

#### `Cache-Control` HTTP Response Headers

`Cache-Control` is a class of HTTP response headers designed to give web
publishers more control over their content, and to address the limitations of
the `Expires` HTTP response header.

The `Cache-Control: max-age=[seconds]` HTTP response header specifies the
maximum amount of time that a resource will be considered fresh. It informs the
client that it may cache the resource, but must revalidate with the server if
the `max-age` value is exceeded. While the `max-age` value is not exceeded, the
client may use the cached resource without revalidating it.

You can change the `max-age` value using the `-m` command-line flag:

	$ mjau -m 2592000

Web font files are not updated very often and by default Mjau sets the
`max-age` value to `2592000` (30 days) to avoid unnecessary requests.

#### Entity Tags

Entity tags (`ETag`s) are a mechanism used by HTTP servers and clients to
determine if a resource from the client's cache matches the one on the server.
When a resource is retrieved by a client, the server sends the resource and
specifies its `ETag` using the `ETag` HTTP response header. Later, if the
client has to validate a cached resource, it uses the `If-None-Match` HTTP
request header to send the `ETag` back to the server. If the current `ETag` of
the resource matches with the one sent by the client, the server responds with
a `304 Not Modified` HTTP status and the client uses the cached resource.

If you wish to enable `ETag`s you can use the `-e` command-line flag:

	$ mjau -e

`ETag`s are disabled by default.

#### Gzip Compression

Compression reduces response time by reducing the size of the HTTP response.

You can enable `gzip` compression using the `-g` command-line flag:

	$ mjau -g

`Gzip` compression is disabled by default.

#### Cross-Origin Resource Sharing (`CORS`)

`CORS` is a mechanism designed to enable client-side cross-origin requests.

You can enable `CORS` using the `-o` command-line flag:

	$ mjau -o

`CORS` is disabled by default.

### Request URL

Web fonts are delivered as CSS files containing one or more `@font-face`
rules. If you have followed the [Quickstart][3] guide to start Mjau, the base
URL of the CSS file is:

	http://localhost:8080/css/

In order to request a web font you have to add the `family=` URL parameter
containing the font family name of the requested web font:

	http://localhost:8080/css/?family=Amaranth

Any spaces in the font family name should be replaced with plus signs (`+`):

	http://localhost:8080/css/?family=Open+Sans

The web fonts are served in `WOFF` format by default. You can specify the web
font format by adding the `format=` URL parameter containing the name of the
requested web font format:

	http://localhost:8080/css/?family=Amaranth&format=eot

Multiple fonts may be grouped in one response by separating their names with a
pipe character (`|`):

	http://localhost:8080/css/?family=Amaranth|Open+Sans

When specifying only the font family name, the server delivers the regular
version of the requested web fonts by default. To request other weights or
styles append a colon (`:`) to the name of the font family, followed by a list
of numerical weights and styles separated by commas (`,`):

	http://localhost:8080/css/?family=Amaranth:400italic,700

When only the weight is specified the normal style of the font is delivered by
default: `700normal` is equivalent to `700`. You can't specify only styles,
you must always append the style to a numerical weight.

The font family names, styles, and weights are defined in the metadata files
from the font library.

### Supported Web Font Formats

Currently, only `EOT` and `WOFF` web font formats are supported.

### Command-line Flags

For a complete list of the available command-line flags use the `-h`
command-line flag:

	$ mjau -h

## License

Mjau is distributed under a BSD-style license. See [LICENSE][4] for details.

[1]: http://golang.org/doc/install
[2]: http://golang.org/cmd/go/#GOPATH_environment_variable
[3]: /noll/mjau#quickstart
[4]: /noll/mjau/blob/master/LICENSE
