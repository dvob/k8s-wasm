use std::error::Error;
use k8s_openapi::api::authentication::v1::{TokenReview, TokenReviewStatus};

fn authenticate(input: String) -> Result<String, Box<dyn Error>> {
    // read input
    let token_review: TokenReview = serde_json::from_str(&input)?;

    // prepare result
    let mut status = TokenReviewStatus::default();
    if token_review.spec.token.unwrap_or_default() == "correct-token" {
        status.authenticated = Some(true)
    } else {
        status.authenticated = Some(false)
    }
    let mut response = TokenReview::default();
    response.status = Some(status);

    // prepare and return output
    let output = serde_json::to_string(&response)?;
    Ok(output)
}

#[no_mangle]
pub unsafe fn run(ptr: u32, len: u32) -> u64 {
    // obtain input
    let input = ptr_to_string(ptr, len);

    // prepare result
    let output = authenticate(input).expect("authentication failed");

    // return result
    let (ptr, len) = string_to_ptr(&output);
    std::mem::forget(output);

    return ((ptr as u64) << 32) | len as u64;
}

/// Allocate memory into the module's linear memory
/// and return the offset to the start of the block.
#[no_mangle]
pub fn alloc(len: usize) -> *mut u8 {
    // create a new mutable buffer with capacity `len`
    let mut buf = Vec::with_capacity(len);
    // take a mutable pointer to the buffer
    let ptr = buf.as_mut_ptr();
    // take ownership of the memory block and
    // ensure that its destructor is not
    // called when the object goes out of scope
    // at the end of the function
    std::mem::forget(buf);
    // return the pointer so the runtime
    // can write data at this offset
    return ptr;
}

#[no_mangle]
pub unsafe fn dealloc(ptr: *mut u8, len: usize) {
    // we only bring data into scope
    let data = Vec::from_raw_parts(ptr, len, len);
    // drop is actually no necessary but we write it down to make it clear what does happen here
    std::mem::drop(data);
}

pub unsafe fn ptr_to_string(ptr: u32, len: u32) -> String {
    //let mut data = Vec::from_raw_parts(ptr as *mut u8, len, len);
    let slice = std::slice::from_raw_parts_mut(ptr as *mut u8, len as usize);
    let utf8 = std::str::from_utf8_unchecked_mut(slice);

    return String::from(utf8);
}

pub unsafe fn string_to_ptr(s: &String) -> (u32, u32) {
    return (s.as_ptr() as u32, s.len() as u32);
}

