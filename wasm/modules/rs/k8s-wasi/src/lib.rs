use k8s_openapi::api::authentication::v1::{TokenReview, TokenReviewStatus};
use std::error::Error;
use std::io::{Read, Write};

fn authenticate(input: &[u8]) -> Result<Vec<u8>, Box<dyn Error>> {
    // read input
    let token_review: TokenReview = serde_json::from_slice(&input)?;

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
    let output = serde_json::to_vec(&response)?;
    Ok(output)
}

#[no_mangle]
fn run() {
    // read input from stdin into buffer
    let mut buf = Vec::new();
    std::io::stdin()
        .read_to_end(&mut buf)
        .expect("failed to read from stdin");

    // prepare result
    let output = authenticate(&buf).expect("failed to authenticate");

    // write result to stdout
    std::io::stdout()
        .write_all(&output)
        .expect("failed to write to stdout");
    std::io::stdout().flush().expect("flush of stdout failed");
}
