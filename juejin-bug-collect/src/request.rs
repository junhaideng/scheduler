use anyhow::{anyhow, Ok, Result};
use reqwest::header::COOKIE;
use serde::{Deserialize, Serialize};

use crate::response::{Data, Response};

const NOT_COLLECT_BUG_URL: &str = "https://api.juejin.cn/user_api/v1/bugfix/not_collect";
const COLLECT_BUG_URL: &str = "https://api.juejin.cn/user_api/v1/bugfix/collect";

#[derive(Serialize, Deserialize, Debug)]
pub struct Request {
    #[serde(rename = "bug_type")]
    bug_type: i64,

    #[serde(rename = "bug_time")]
    bug_time: i64,
}

pub async fn get_bugs(client: &reqwest::Client, cookie: &str) -> Result<Vec<Data>> {
    let resp = client
        .post(NOT_COLLECT_BUG_URL)
        .header(COOKIE, cookie)
        .send()
        .await?;

    if !resp.status().is_success() {
        return Err(anyhow!("status is {}", resp.status().as_u16()));
    }

    let resp = resp.json::<Response>().await?;
    if resp.err_no != 0 {
        return Err(anyhow!(resp.err_msg));
    }

    match resp.data {
        Some(data) => Ok(data),
        None => return Err(anyhow!("data is null")),
    }
}

pub async fn collect_bug(client: &reqwest::Client, cookie: &str, data: &Data) -> Result<bool> {
    let resp = client
        .post(COLLECT_BUG_URL)
        .header(COOKIE, cookie)
        .json(data)
        .send()
        .await?;

    if !resp.status().is_success() {
        return Err(anyhow!("status is {}", resp.status().as_u16()));
    }

    let resp: Response = resp.json().await?;
    if resp.err_no != 0 {
        return Err(anyhow!(resp.err_msg));
    }
    println!("collect bug: {:?}", resp); 

    Ok(true)
}
