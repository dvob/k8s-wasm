use std::io::Write;
use k8s_openapi::api::authentication::v1::*;
use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize)]
struct Request {
    request: TokenReview,
}

#[derive(Serialize, Deserialize)]
struct Response {
    response: TokenReview,
    error: Option<String>,
}

// why does this signature produce a functino which takes a parameter?
//fn auth() -> Result<(), Box<dyn std::error::Error>> {
#[no_mangle]
fn authn() {
    let req: Request = serde_json::from_reader(std::io::stdin()).unwrap();
    let token_review = req.request;
    let token = token_review.spec.token.expect("token missing");

    let mut response = TokenReview::default();
    let mut status = TokenReviewStatus::default();

    if token == "my-test-token" {
        status.authenticated = Some(true);
        status.user = Some(UserInfo{
            username: Some("my-user".to_string()),
            uid: Some("1337".to_string()),
            groups: Some(vec![
                "system:masters".to_string(),
            ]),
            extra: None,
        });
    } else {
        status.authenticated = Some(false);
        status.error = Some("invalid token".to_string())
    }

    response.status = Some(status);

    let resp = Response{
        response,
        error: None,
    };

    serde_json::to_writer(std::io::stdout(), &resp).unwrap();
    std::io::stdout().flush().unwrap();
}