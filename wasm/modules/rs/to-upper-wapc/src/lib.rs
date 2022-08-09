fn run(input: &[u8]) -> wapc_guest::CallResult {
    // prepare result
    let output = input.to_ascii_uppercase();

    Ok(output)
}

#[no_mangle]
pub fn wapc_init() {
    wapc_guest::register_function("run", run);
}