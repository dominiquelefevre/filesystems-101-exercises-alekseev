#![no_std]

use core::ffi::{CStr, c_int};

pub mod sys {
    use core::ffi::{c_char, c_int};

    unsafe extern "C" {
        pub fn report_path(path: *const c_char);
        pub fn report_error(parent: *const c_char, child: *const c_char, errno_code: c_int);
    }
}

pub fn report_path(path: &CStr) {
    unsafe { sys::report_path(path.as_ptr()) }
}

pub fn report_error(parent: &CStr, child: &CStr, errno_code: c_int) {
    unsafe { sys::report_error(parent.as_ptr(), child.as_ptr(), errno_code) }
}

#[doc(hidden)]
#[macro_export]
macro_rules! export_solution {
    ($solution:ident) => {
        #[unsafe(no_mangle)]
        unsafe extern "C" fn abspath(path: *const ::core::ffi::c_char) {
            let f: fn(&::core::ffi::CStr) = $solution::abspath;
            f(unsafe { ::core::ffi::CStr::from_ptr(path) })
        }
    };
}
