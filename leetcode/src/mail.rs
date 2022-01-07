use lettre::transport::smtp::authentication::Credentials;
use lettre::{Message, SmtpTransport, Transport};
use std::error::Error;

pub fn send(
    text: &str,
    username: &str,
    password: &str,
    to: &str,
    server: &str,
) -> Result<(), Box<dyn Error>> {
    let email = Message::builder()
        .from(format!("Bot <{}>", username).parse()?)
        .to(format!("Edgar <{}>", to).parse()?)
        .subject("LeetCode打卡提醒")
        .body(text.to_string())?;

    let creds = Credentials::new(username.to_string(), password.to_string());

    // Open a remote connection to gmail
    let mailer = SmtpTransport::relay(server)?.credentials(creds).build();

    // Send the email
    mailer.send(&email)?;
    Ok(())
}
