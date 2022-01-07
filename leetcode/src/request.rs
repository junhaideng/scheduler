use crate::mail;
use crate::response::Response;
use reqwest;

pub async fn get(cookie: &str) -> Result<Response, String> {
    // http://172.18.96.1:8888
    // let proxy = reqwest::Proxy::http("http://172.18.96.1:8888").map_err(|e| e.to_string())?;

    let client = reqwest::Client::builder()
        // .proxy(proxy)
        .build()
        .map_err(|e| e.to_string())?;

    let resp = client
        .post("https://leetcode-cn.com/graphql/")
        .body(r#"{"query":"\n    query questionOfToday {\n  todayRecord {\n    date\n    userStatus\n    question {\n      questionId\n      frontendQuestionId: questionFrontendId\n      difficulty\n      title\n      titleCn: translatedTitle\n      titleSlug\n      paidOnly: isPaidOnly\n      freqBar\n      isFavor\n      acRate\n      status\n      solutionNum\n      hasVideoSolution\n      topicTags {\n        name\n        nameTranslated: translatedName\n        id\n      }\n      extra {\n        topCompanyTags {\n          imgUrl\n          slug\n          numSubscribed\n        }\n      }\n    }\n    lastSubmission {\n      id\n    }\n  }\n}\n    ","variables":{}}"#)
        // .json(r#"{"query":"\n    query dailyQuestionRecords($year: Int!, $month: Int!) {\n  dailyQuestionRecords(year: $year, month: $month) {\n    date\n    userStatus\n    question {\n      questionFrontendId\n      title\n      titleSlug\n      translatedTitle\n    }\n  }\n}\n    ","variables":{"year":2022,"month":1}}"#)
        .header("referer", "https://leetcode-cn.com/problemset/all/")
        .header("cookie", cookie)
        .header("content-type", "application/json")
        .send()
        .await
        .map_err(|e| e.to_string())?;

    // println!("{}", resp.text().await.unwrap());
    // Err("".into())

    let res: Response = resp.json().await.map_err(|e| e.to_string())?;
    Ok(res)
}

pub async fn handle(cookie: &str, username: &str, password:&str, to: &str, server: &str) {
    let res = get(cookie).await;

    match res {
        Ok(res) => {
            for record in res.data.today_record.iter() {
                if record.user_status.eq("FINISH") {
                    println!("今日已经打卡");
                    return;
                }
            }

            match mail::send("LeetCode 今天还没有打卡哦", username, password, to, server){
              Ok(_) => println!("发送邮件成功"),
              Err(e)=> println!("发送邮件失败: {}", e.to_string())
            }
        }
        Err(e) => {
            println!("get response failed: {}", e);
        }
    }
}
