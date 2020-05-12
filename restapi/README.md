# restapi

Plugin **restapi** provides a REST API to communicate with the *collectd*
daemon.

## Description

The *restapi plugin* starts a webserver and waits for incoming REST API
requests.

## Building

To build this plugin, the collectd header files are required.

On Debian and Ubuntu, the collectd headers are available from the
`collectd-dev` package. Once installed, add the import paths to the
`CGI_CPPFLAGS`:

```bash
export CGO_CPPFLAGS="-I/usr/include/collectd/core/daemon \
-I/usr/include/collectd/core -I/usr/include/collectd"
```

Alternatively, you can grab the collectd sources, run the `configure` script,
and reference the header files from there:

```bash
TOP_SRCDIR="${HOME}/collectd"
export CGO_CPPFLAGS="-I${TOP_SRCDIR}/src -I${TOP_SRCDIR}/src/daemon"
```

Then build the plugin with the "c-shared" buildmode:

```bash
go build -buildmode=c-shared -o restapi.so
```
