#![feature(proc_macro_hygiene, decl_macro)]
#[macro_use] extern crate rocket;
#[macro_use] extern crate rocket_contrib;
#[macro_use] extern crate serde_derive;

extern crate reqwest;

mod blocking_client;
mod types;
mod server;

pub use crate::types::*;

fn main() {

    println!("Hello, world!!");
}


#[cfg(test)]
mod tests {
    pub use crate::server::*;
    use super::*;

    #[test]
    fn test_server() {
    }

}