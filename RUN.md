docker run -d -p 127.0.0.1:5432:5432 -e "POSTGRES_PASSWORD=postgres" -e "POSTGRES_INITDB_ARGS=--lc-collate C" postgres:15

./bin/tests-runner serve --server.addr=:9091 --server.auth.disable --db.host=127.0.0.1 --db.database=postgres --db.dialect=postgres --db.user=postgres --db.password=postgres

# Rust solutions

Solutions written in Rust (see exercises/rust/README.md) are built on the grading host
by `make build`, like C solutions. The host needs, besides gcc and make:

- rustc 1.85 or newer on PATH
- curl, tar and sha256sum, and network access to static.crates.io.

The crates listed in exercises/rust/stdlib/crates.mk are downloaded and compiled on
demand into $XDG_CACHE_HOME/fs101-rust (default ~/.cache/fs101-rust of the user that
runs `tests-runner serve`; override with FS101_RUST_CACHE)
