CRATES := libc bitflags memchr arrayvec byteorder bytemuck spin errno linux-raw-sys rustix equivalent foldhash allocator-api2 hashbrown
PUBLIC_CRATES := libc bitflags memchr arrayvec byteorder bytemuck spin rustix hashbrown

libc_VERSION := 0.2.190
libc_SHA256 := ce5d3ddc6d3fa000eb1536d85e147bfe31aacaba692ed6a876f95cb7c855be78
libc_EDITION := 2021
libc_CFG_linux := linux_time_bits64

bitflags_VERSION := 2.13.2
bitflags_SHA256 := 3ded4057c258ba199e2d26386d3af3780957ecaee6c4ef4041c6b4b8b97c0b06
bitflags_EDITION := 2021

memchr_VERSION := 2.8.3
memchr_SHA256 := cf8baf1c55e62ffcace7a9f06f4bd9cd3f0c4beb022d3b367256b91b87513d98
memchr_EDITION := 2021
memchr_FEATURES := alloc

arrayvec_VERSION := 0.7.8
arrayvec_SHA256 := d3fb67a6e08acf24fdeccbac2cb6ac4305825bd1f117462e0e6f2f193345ad56
arrayvec_EDITION := 2018

byteorder_VERSION := 1.5.0
byteorder_SHA256 := 1fd0f2584146f6f2ef48085050886acf353beff7305ebd1ae69500e27c67f64b
byteorder_EDITION := 2021

bytemuck_VERSION := 1.25.2
bytemuck_SHA256 := 95832e849adfb21180ccb6826a99da14e5d266ae5c2e668e1602cf234f153797
bytemuck_EDITION := 2018
bytemuck_FEATURES := extern_crate_alloc aarch64_simd avx512_simd align_offset alloc_uninit const_zeroed impl_core_error min_const_generics must_cast must_cast_extra pod_saturating track_caller transparentwrapper_extra zeroable_atomics zeroable_maybe_uninit zeroable_unwind_fn

spin_VERSION := 0.12.3
spin_SHA256 := 0134f9043ed38b087ac4f7d4af44c79e2c9e5094421fe3164f435ce585953b10
spin_EDITION := 2021
spin_FEATURES := mutex spin_mutex rwlock once lazylock barrier

errno_VERSION := 0.3.14
errno_SHA256 := 39cab71617ae0d63f51a36d69f866391735b51691dbda63cf6f96d042b63efeb
errno_EDITION := 2018
errno_DEPS := libc=libc

linux-raw-sys_VERSION := 0.12.1
linux-raw-sys_SHA256 := 32a66949e030da00e8c7d4434b251670a91556f4144941d37452769c25d58a53
linux-raw-sys_EDITION := 2021
linux-raw-sys_FEATURES := general ioctl no_std prctl system

rustix_VERSION := 1.1.5
rustix_SHA256 := 891efababe418670775f199f0d233d84843c227a0949a883ce15b37c78d6629d
rustix_EDITION := 2021
rustix_FEATURES := use-libc libc libc_errno alloc event fs mm param pipe process rand stdio system thread time
rustix_DEPS := bitflags=bitflags libc=libc libc_errno=errno linux_raw_sys=linux-raw-sys
rustix_CFG := lower_upper_exp_for_non_zero rustc_diagnostics
rustix_CFG_linux := libc linux_raw_dep linux_like linux_kernel
rustix_CFG_macos := libc apple bsd

equivalent_VERSION := 1.0.3
equivalent_SHA256 := 00d174d5400e5e8fd687ad1049e2f578285fa914201b1af7e8b112a4546bd826
equivalent_EDITION := 2015

foldhash_VERSION := 0.2.0
foldhash_SHA256 := 77ce24cb58228fbb8aa041425bb1050850ac19177686ea6e0f41a70416f56fdb
foldhash_EDITION := 2021

allocator-api2_VERSION := 0.2.21
allocator-api2_SHA256 := 683d7910e743518b0e34f1186f92494becacb047c7b6bf616c96772180fef923
allocator-api2_EDITION := 2018
allocator-api2_FEATURES := alloc

hashbrown_VERSION := 0.17.1
hashbrown_SHA256 := ed5909b6e89a2db4456e54cd5f673791d7eca6732202bbf2a9cc504fe2f9b84a
hashbrown_EDITION := 2024
hashbrown_FEATURES := default default-hasher inline-more allocator-api2 equivalent raw-entry
hashbrown_DEPS := allocator_api2=allocator-api2 equivalent=equivalent foldhash=foldhash
