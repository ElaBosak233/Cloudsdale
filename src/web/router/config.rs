use axum::{routing::get, Router};

use crate::web::controller;

pub fn router() -> Router {
    return Router::new()
        .route("/", get(controller::config::find))
        .route("/favicon", get(controller::config::get_favicon));
}