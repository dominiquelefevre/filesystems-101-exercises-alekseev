# Solving exercises in Rust

You can write solutions to these exercises in Rust instead of C:

| exercise | the function lib.rs must define |
|---|---|
| 00-ps | `pub fn ps()` |
| 01-lsof | `pub fn lsof()` |
| 02-fuse-helloworld | `pub fn helloworld(mntp: &CStr) -> c_int` |
| 13-realpath | `pub fn abspath(path: &CStr)` |

Rust solutions live in `rust/<exercise>/lib.rs`, which starts as a template
with the function to implement (like `solution.h` in C). The directory also
has a `Makefile`, `fs101.rs` and the C test harness (`main.c`,
`callbacks.c`, ...), which calls your function. The tests are the same for
both languages.

The grader checks `c/<exercise>` unless you select Rust: create an empty file
`rust/<exercise>/.selected` and commit it. Delete it to switch back to C.
`rust/<exercise>` must not contain `.c` files of your own.

## Rules

- `lib.rs` must start with `#![no_std]`. `std` is not available; use `core`,
  and `alloc` for `Vec`, `String`, `Box`, `format!` (add `extern crate alloc;`).
- The edition is 2024. Warnings are errors, as with `-Werror` in C. Debug
  assertions and integer overflow checks are on.
- Define the function from the table as a plain `pub fn` in `lib.rs`. Don't
  add `#[no_mangle]` or `extern "C"`; that glue is generated for you. A
  wrong signature fails to compile with a type mismatch.
- Call the harness callbacks through the `fs101` crate (`fs101.rs` in
  `rust/<exercise>`): for example `fs101::report_file(path)`. The raw C
  declarations are in `fs101::sys`. Don't edit `fs101.rs`; the grader
  replaces it.
- Make system calls with the `rustix` crate (safe wrappers, no `unsafe`
  needed) or with `libc`. Both go through the C library, so the same
  restrictions as in C apply, because the checks run in the same binary: `strcpy()`, `sprintf()`,
  `lseek()` and friends from `rust/stdlib/forbidden_funcs.c` exit with an error
  (use `pread()`/`pwrite()` instead of `lseek()`), and 13-realpath rejects
  `realpath()` and `chdir()`.
- Tests run under Valgrind, so memory you leak (e.g. `mem::forget`,
  `Box::leak`) and file descriptors you don't close fail the tests. The
  global allocator is libc `malloc()`/`free()`.
- A panic prints its message to stderr and to the Valgrind log, then aborts.
  Don't define your own `#[panic_handler]` or `#[global_allocator]`.
- Split code into modules as usual: `mod util;` in `lib.rs` loads `util.rs`.

These crates are available with `use` (no `Cargo.toml` needed):

| crate | version | notes |
|---|---|---|
| rustix | 1.1.5 | libc backend; `fs`, `process`, `pipe`, `mm`, `event`, `thread`, `time`, `stdio`, `param`, `rand`, `system`, `alloc` |
| libc | 0.2.190 | |
| bitflags | 2.13.2 | |
| memchr | 2.8.3 | feature `alloc` |
| arrayvec | 0.7.8 | no `std` |
| byteorder | 1.5.0 | no `std` |
| bytemuck | 1.25.2 | no `derive` |
| spin | 0.12.3 | no `lock_api` |
| hashbrown | 0.17.1 | default features |

## Example: 01-lsof

```rust
#![no_std]
#![forbid(unsafe_code)]

extern crate alloc;

use alloc::ffi::CString;
use alloc::format;
use alloc::vec::Vec;

use rustix::fs::{CWD, Dir, Mode, OFlags, openat, readlinkat};

pub fn lsof() {
    let dir = openat(CWD, c"/proc/self/fd", OFlags::RDONLY | OFlags::DIRECTORY, Mode::empty())
        .and_then(Dir::new);
    let dir = match dir {
        Ok(dir) => dir,
        Err(e) => return fs101::report_error(c"/proc/self/fd", e.raw_os_error()),
    };
    for entry in dir {
        let Ok(entry) = entry else { break };
        let name = entry.file_name().to_str().unwrap();
        let link = CString::new(format!("/proc/self/fd/{name}")).unwrap();
        if let Ok(target) = readlinkat(CWD, link.as_c_str(), Vec::new()) {
            fs101::report_file(&target);
        }
    }
}
```

(This lists only the solution's own files, so it doesn't pass the tests.)

With `rustix` and the `fs101` wrappers, 00-ps, 01-lsof and 13-realpath can
be solved without `unsafe`; `#![forbid(unsafe_code)]` makes the compiler check
that. 02-fuse-helloworld still needs `unsafe`: libfuse calls your callbacks
with raw pointers.

## Building

`make` in `rust/<exercise>` builds and runs `a.out`, as in `c/<exercise>`. It
needs rustc 1.85 or newer on `PATH` (or set `RUSTC`), plus `curl`, `tar` and
`sha256sum`. The first build downloads the crates above and compiles them into
`~/.cache/fs101-rust` (set `FS101_RUST_CACHE` to use another location); later
builds reuse them until rustc changes. Build artifacts go to `.rust/`.

The grader may run a different rustc version than yours. Since warnings are
errors, a lint that only one of the versions knows can make the build fail on
one side only.

For rust-analyzer, run `make rust-project.json` in `rust/<exercise>` and open
that directory in your editor.
