use std::env;

use anyhow::{Ok, Result};
use request::get_bugs;
use reqwest::Client;

use crate::request::collect_bug;
mod request;
mod response;

#[tokio::main]
async fn main() -> Result<()> {
    let cookie = match env::args().skip(1).next() {
        Some(c) => c,
        None => {
            println!("no cookie");
            return Ok(());
        }
    };
    let client = Client::new();
    let bugs = get_bugs(&client, &cookie).await?;

    for bug in bugs.iter() {
        collect_bug(&client, &cookie, bug).await?;
    }

    println!("{:?}", bugs);

    Ok(())
}
