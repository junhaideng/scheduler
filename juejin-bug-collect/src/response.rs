use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug)]
pub struct Response {
    #[serde(rename = "err_no")]
    pub err_no: i64,

    #[serde(rename = "err_msg")]
    pub err_msg: String,

    #[serde(rename = "data")]
    pub data: Option<Vec<Data>>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct Data {
    #[serde(rename = "bug_type")]
    pub bug_type: i64,

    #[serde(rename = "bug_time")]
    pub bug_time: i64,

    #[serde(rename = "bug_show_type", skip_serializing)]
    pub bug_show_type: i64,

    #[serde(rename = "is_first", skip_serializing)]
    pub is_first: bool,
}
