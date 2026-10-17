#![no_std]

#[path = "../stdlib/fs101_fuse.rs"]
pub mod fuse;

#[doc(hidden)]
#[macro_export]
macro_rules! export_solution {
    ($solution:ident) => {
        #[unsafe(no_mangle)]
        unsafe extern "C" fn helloworld(mntp: *const ::core::ffi::c_char) -> ::core::ffi::c_int {
            let f: fn(&::core::ffi::CStr) -> ::core::ffi::c_int = $solution::helloworld;
            f(unsafe { ::core::ffi::CStr::from_ptr(mntp) })
        }
    };
}
