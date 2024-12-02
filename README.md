# lf-install

## lf-install is not a package manager

This library is intended for use within Go programs or automated environments (such as CIs)
which have some business downloading or otherwise locating Linux Foundation binaries.

The included command-line utility, `lf-install`, is a convenient way of using
the library in ad-hoc or CI shell scripting outside of Go.

`lf-install` does **not**:

- Determine suitable installation path based on target system. e.g. in `/usr/bin` or `/usr/local/bin` on Unix based system.
- Deal with execution of installed binaries (via service files or otherwise).
- Upgrade existing binaries on your system.
- Add nor link downloaded binaries to your `$PATH`.

## API

The `Installer` offers a few high-level methods:

- `Ensure(context.Context, []src.Source)` to find, install, or build a product version
- `Install(context.Context, []src.Installable)` to install a product version

### Sources

The `Installer` methods accept number of different `Source` types.
Each comes with different trade-offs described below.

- `fs.{AnyVersion,ExactVersion,Version}` - Finds a binary in `$PATH` (or additional paths)
  - **Pros:**
    - This is most convenient when you already have the product installed on your system
      which you already manage.
  - **Cons:**
    - Only relies on a single version, expects _you_ to manage the installation
    - _Not recommended_ for any environment where product installation is not controlled or managed by you (e.g. default GitHub Actions image managed by GitHub)
- `checkpoint.LatestVersion` - Downloads, verifies & installs any known product available in HashiCorp Checkpoint
  - **Pros:**
    - Checkpoint typically contains only product versions considered stable
  - **Cons:**
    - Installation may consume some bandwidth, disk space and a little time
    - Currently doesn't allow installation of old versions or enterprise versions (see `releases` above)
- `build.GitRevision` - Clones raw source code and builds the product from it
  - **Pros:**
    - Useful for catching bugs and incompatibilities as early as possible (prior to product release).
  - **Cons:**
    - Building from scratch can consume significant amount of time & resources (CPU, memory, bandwidth, disk space)
    - There are no guarantees that build instructions will always be up-to-date
    - There's increased likelihood of build containing bugs prior to release
    - Any CI builds relying on this are likely to be fragile

## Example Usage

See examples at <https://pkg.go.dev/github.com/chushi-io/lf-install#example-Installer>.

## CLI

In addition to the Go library, which is the intended primary use case of `lf-install`, we also distribute CLI.

The CLI comes with some trade-offs:

- more limited interface compared to the flexible Go API (installs specific versions of products via `releases.ExactVersion`)
- minimal environment pre-requisites (no need to compile Go code)
- see ["lf-install is not a package manager"](https://github.com/chushi-io/lf-install#lf-install-is-not-a-package-manager)

### Usage

```text
Usage: lf-install install [options] -version <version> <product>

  This command installs a Linux Foundation product.
  Options:
    -version  [REQUIRED] Version of product to install.
    -path     Path to directory where the product will be installed.
              Defaults to current working directory.
    -log-file Path to file where logs will be written. /dev/stdout
              or /dev/stderr can be used to log to STDOUT/STDERR.
```

```sh
lf-install install -version 1.3.7 tofu
```

```sh
lf-install: will install tofu@1.3.7
installed tofu@1.3.7 to /current/working/dir/tofu
```
