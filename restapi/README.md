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

## Configuration

### Synopsis

```
LoadPlugin restapi
<Plugin "restapi">
  Addr "::"
  Port "8443"
  CertFile "/path/to/cert_file.pem"
  KeyFile "/path/to/key_file.pem"
</Plugin>
```

### Options

*   **Addr** *Network address*

    Addredd to listen to. Defaults to `""` (any address).
*   **Port** *Port*

    Post to listen to. Defaults to `8080` (`8443` if **CertFile** is specified).
*   **CertFile** *Path*<br>
    **KeyFile** *Path*

    TLS certificate and key files. Refer to
    [`"net/http".ListenAndServeTLS`](https://golang.org/pkg/net/http/#ListenAndServeTLS)
    for more information on the TLS setup.
