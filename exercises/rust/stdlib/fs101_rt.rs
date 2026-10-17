#![no_std]

use core::alloc::{GlobalAlloc, Layout};
use core::ffi::{c_int, c_void};
use core::fmt::{self, Write};
use core::panic::PanicInfo;
use core::ptr;

unsafe extern "C" {
    fn malloc(size: usize) -> *mut c_void;
    fn calloc(nmemb: usize, size: usize) -> *mut c_void;
    fn realloc(p: *mut c_void, size: usize) -> *mut c_void;
    fn free(p: *mut c_void);
    fn posix_memalign(memptr: *mut *mut c_void, alignment: usize, size: usize) -> c_int;
    fn write(fd: c_int, buf: *const c_void, count: usize) -> isize;
    fn abort() -> !;
}

const MALLOC_ALIGN: usize = 2 * size_of::<usize>();

struct LibcAllocator;

unsafe fn alloc_aligned(layout: Layout) -> *mut u8 {
    let align = layout.align().max(size_of::<*mut c_void>());
    let mut p = ptr::null_mut();
    if unsafe { posix_memalign(&mut p, align, layout.size()) } != 0 {
        return ptr::null_mut();
    }
    p.cast()
}

unsafe impl GlobalAlloc for LibcAllocator {
    unsafe fn alloc(&self, layout: Layout) -> *mut u8 {
        if layout.align() <= MALLOC_ALIGN {
            unsafe { malloc(layout.size()).cast() }
        } else {
            unsafe { alloc_aligned(layout) }
        }
    }

    unsafe fn alloc_zeroed(&self, layout: Layout) -> *mut u8 {
        if layout.align() <= MALLOC_ALIGN {
            unsafe { calloc(1, layout.size()).cast() }
        } else {
            let p = unsafe { alloc_aligned(layout) };
            if !p.is_null() {
                unsafe { p.write_bytes(0, layout.size()) };
            }
            p
        }
    }

    unsafe fn dealloc(&self, p: *mut u8, _layout: Layout) {
        unsafe { free(p.cast()) }
    }

    unsafe fn realloc(&self, p: *mut u8, layout: Layout, new_size: usize) -> *mut u8 {
        if layout.align() <= MALLOC_ALIGN {
            return unsafe { realloc(p.cast(), new_size).cast() };
        }
        let new_layout = unsafe { Layout::from_size_align_unchecked(new_size, layout.align()) };
        let q = unsafe { alloc_aligned(new_layout) };
        if !q.is_null() {
            unsafe {
                ptr::copy_nonoverlapping(p, q, layout.size().min(new_size));
                free(p.cast());
            }
        }
        q
    }
}

#[global_allocator]
static ALLOCATOR: LibcAllocator = LibcAllocator;

struct Stderr;

impl Write for Stderr {
    fn write_str(&mut self, s: &str) -> fmt::Result {
        let mut buf = s.as_bytes();
        while !buf.is_empty() {
            let n = unsafe { write(2, buf.as_ptr().cast(), buf.len()) };
            if n <= 0 {
                return Err(fmt::Error);
            }
            buf = &buf[n as usize..];
        }
        Ok(())
    }
}

struct VgMessage {
    buf: [u8; 1024],
    len: usize,
}

impl Write for VgMessage {
    fn write_str(&mut self, s: &str) -> fmt::Result {
        for &c in s.as_bytes() {
            let esc: &[u8] = if c == b'%' { b"%%" } else { core::slice::from_ref(&c) };
            if self.len + esc.len() >= self.buf.len() {
                return Err(fmt::Error);
            }
            self.buf[self.len..self.len + esc.len()].copy_from_slice(esc);
            self.len += esc.len();
        }
        Ok(())
    }
}

#[allow(unused_variables)]
unsafe fn valgrind_request(args: &[usize; 6]) -> usize {
    #[allow(unused_mut)]
    let mut result = 0;
    #[cfg(target_arch = "x86_64")]
    unsafe {
        core::arch::asm!(
            "rol rdi, 3", "rol rdi, 13", "rol rdi, 61", "rol rdi, 51",
            "xchg rbx, rbx",
            inout("rdx") result,
            in("rax") args.as_ptr(),
            out("rdi") _,
            options(nostack),
        );
    }
    #[cfg(target_arch = "aarch64")]
    unsafe {
        core::arch::asm!(
            "ror x12, x12, #3", "ror x12, x12, #13", "ror x12, x12, #51", "ror x12, x12, #61",
            "orr x10, x10, x10",
            inout("x3") result,
            in("x4") args.as_ptr(),
            out("x12") _,
            options(nostack),
        );
    }
    result
}

const VG_USERREQ_PRINTF_BACKTRACE_VALIST_BY_REF: usize = 0x1404;

#[panic_handler]
fn panic(info: &PanicInfo<'_>) -> ! {
    let _ = writeln!(Stderr, "{info}");

    let mut msg = VgMessage { buf: [0; 1024], len: 0 };
    let _ = writeln!(msg, "{info}");
    msg.buf[msg.len] = 0;
    let va_list = [0usize; 4];
    unsafe {
        valgrind_request(&[
            VG_USERREQ_PRINTF_BACKTRACE_VALIST_BY_REF,
            msg.buf.as_ptr() as usize,
            va_list.as_ptr() as usize,
            0,
            0,
            0,
        ])
    };

    unsafe { abort() }
}

#[unsafe(no_mangle)]
extern "C" fn rust_eh_personality() {}
