mod mail;
mod request;
mod response;

use std::process;

use clap::{App, Arg};

#[tokio::main]
async fn main() {
    let app = App::new("payapp")
        .arg(
            Arg::new("cookie")
                .short('c')
                .long("cookie")
                .help("登录cookie")
                .takes_value(true)
                .required(true),
        )
        .arg(
            Arg::new("username")
                .short('u')
                .long("username")
                .help("邮件用户名")
                .takes_value(true)
                .required(true),
        )
        .arg(
            Arg::new("password")
                .short('p')
                .long("password")
                .help("邮件密码")
                .takes_value(true)
                .required(true),
        )
        .arg(
            Arg::new("server")
                .short('s')
                .long("server")
                .help("邮件服务器")
                .default_value("smtp.126.com"),
        )
        .arg(
            Arg::new("to")
                .long("to")
                .help("邮件接收方")
                .takes_value(true)
                .required(true),
        )
        .get_matches();

    let cookie = match app.value_of("cookie") {
        Some(cookie) => cookie,
        None => {
            println!("请输入 cookie");
            process::exit(-1);
        }
    };

    let username = match app.value_of("username") {
        Some(username) => username,
        None => {
            println!("请输入邮箱用户名");
            process::exit(-1);
        }
    };

    let password = match app.value_of("password") {
        Some(password) => password,
        None => {
            println!("请输入邮箱密码");
            process::exit(-1);
        }
    };

    let server = match app.value_of("server") {
        Some(server) => server,
        None => "smtp.126.com",
    };

    let to = match app.value_of("to") {
        Some(to) => to,
        None => {
            println!("请输入邮件接收方");
            process::exit(-1);
        }
    };

    request::handle(cookie, username, password, to, server).await;
}
