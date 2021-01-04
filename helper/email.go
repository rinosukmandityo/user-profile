package helper

import (
	"fmt"
	"log"
	"net/smtp"
	"net/url"
	"os"
)

var (
	SMTP_HOST      = os.Getenv("SMTP_HOST")
	SMTP_PORT      = os.Getenv("SMTP_PORT")
	SOURCE_EMAIL   = os.Getenv("SOURCE_EMAIL")
	EMAIL_PASSWORD = os.Getenv("EMAIL_PASSWORD")
)

func ResetPasswordMailContent(user map[string]string, tokenid, url string) error {
	fullname := user["fullname"]
	tokenid += "&d=" + user["id"]

	mailContent := `
	<html>
		<body>
			<table>
				<tbody>
					<tr>
						<td style="padding-bottom:20px">
							<h2 style="margin:0;color:#262626;font-weight:700;font-size:20px;line-height:1.2">Hi ` + fullname + `,</h2>
						</td>
					</tr>
					<tr>
						<td style="padding-bottom:20px"> 
							<p style="margin:0;color:#4c4c4c;font-weight:400;font-size:16px;line-height:1.25">
								Reset your password, and we'll get you on your way.
							</p>
						</td>
					</tr>
					<tr>
						<td style="padding-bottom:20px"> 
							<p style="margin:0;color:#4c4c4c;font-weight:400;font-size:16px;line-height:1.25">
								To change your password, click <a href="` + url + tokenid + `" style="color:#008cc9;display:inline-block;text-decoration:none" target="_blank">here</a> or paste the following link into your browser:
							</p>
						</td>
					</tr>
					<tr> 
						<td style="padding-bottom:20px">
							<p style="margin:0;color:#4c4c4c;font-weight:400;font-size:16px;line-height:1.25"><a href="` + url + tokenid + `">` + url + tokenid + `</a></p>
						</td> 
					</tr>
					<tr>
						<td style="padding-bottom:20px">
							<p style="margin:0;color:#4c4c4c;font-weight:400;font-size:16px;line-height:1.25">This link will expire in 15 minutes, so be sure to use it right away.</p>
						</td>
					</tr>
					<tr> 
						<td style="padding-bottom:20px"> 
							<p style="margin:0;color:#4c4c4c;font-weight:400;font-size:16px;line-height:1.25">Thank you for using TalentPro!</p> 
							<p style="margin:0;color:#4c4c4c;font-weight:400;font-size:16px;line-height:1.25">The TalentPro Team</p>
						</td>
					</tr>
				</tbody>
			</table>
		</body>
	</html>
	`
	mailSubject := fullname + ", here's the link to reset your password"
	if e := SendMail([]string{user["email"]}, mailSubject, mailContent); e != nil {
		return e
	}

	log.Println("Mail sent!")
	return nil
}

func SendMail(to []string, subject, message string) error {
	content := "Subject: " + subject + "\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		message

	auth := smtp.PlainAuth("", SOURCE_EMAIL, EMAIL_PASSWORD, SMTP_HOST)
	smtpAddr := fmt.Sprintf("%s:%s", SMTP_HOST, SMTP_PORT)

	if e := smtp.SendMail(smtpAddr, auth, SOURCE_EMAIL, to, []byte(content)); e != nil {
		return e
	}

	return nil
}

func ConstructEmailURL(rawurl string, config map[string]string) string {
	u, _ := url.Parse(rawurl)
	return u.Scheme + "://" + u.Host + config["suffix"]
}
