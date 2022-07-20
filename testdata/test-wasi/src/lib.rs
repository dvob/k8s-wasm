#[no_mangle]
fn my_function() {
    println!("stdout output");
}

#[no_mangle]
fn my_panic() {
    panic!("panic output");
}

#[no_mangle]
fn my_error() {
    eprintln!("stderr output");
    std::process::exit(1);
}
