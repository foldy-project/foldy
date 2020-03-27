#[macro_use]
extern crate log;
#[macro_use]
extern crate serde_derive;
#[macro_use]
extern crate k8s_openapi;
#[macro_use]
extern crate kube_derive;
#[macro_use]
extern crate kube;
#[macro_use] extern crate async_trait;
extern crate reqwest;
#[macro_use] extern crate futures;
use async_trait::async_trait;

use actix_web::{middleware, web, App, HttpRequest, HttpResponse, HttpServer};
use std::cell::Cell;
use std::io;
use std::sync::atomic::{AtomicUsize, Ordering};
use std::sync::Mutex;

use thiserror::Error;


mod types;
mod server_actixweb;

pub use crate::types::*;

fn main() {

    println!("Hello, world!!");
}


#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_server() {
    }

}