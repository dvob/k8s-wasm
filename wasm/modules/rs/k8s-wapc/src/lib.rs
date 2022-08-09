use k8s_openapi::api::authentication::v1::{TokenReview, TokenReviewStatus};
use std::error::Error;

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

fn run(input: &[u8]) -> wapc_guest::CallResult {
    // prepare result
    let output = authenticate(input).expect("failed to authenticate");

    Ok(output)
}

#[no_mangle]
pub fn wapc_init() {
    wapc_guest::register_function("run", run);
}
