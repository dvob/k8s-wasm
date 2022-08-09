use std::io::{Read, Write};

#[no_mangle]
fn run() {
    // read input from stdin into buffer
    let mut buf = Vec::new();
    std::io::stdin()
        .read_to_end(&mut buf)
        .expect("failed to read from stdin");

    // prepare result
    let output = buf.to_ascii_uppercase();

    // write result to stdout
    std::io::stdout()
        .write_all(&output)
        .expect("failed to write to stdout");
    std::io::stdout().flush().expect("flush of stdout failed");
}

