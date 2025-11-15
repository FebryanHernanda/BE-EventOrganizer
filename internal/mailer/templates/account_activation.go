package templates

import "fmt"

func ActivationTemplate(name, email, date, link string) string {
	return fmt.Sprintf(`
	<div style="font-family:Arial, sans-serif; background:#f4f4f4; padding:20px;">
		<div style="max-width:600px; margin:0 auto; background:#fff; border-radius:8px; padding:20px;">
			<h2 style="color:#222;">Welcome, %s!</h2>
			<p>Thank you for registering with us. Below are your registration details:</p>

			<table style="margin:15px 0;">
				<tr><td><strong>Email:</strong></td><td>%s</td></tr>
				<tr><td><strong>Registered At:</strong></td><td>%s</td></tr>
			</table>

			<p>Click the button below to activate your account:</p>

			<a href="%s" style="background:#007bff; color:#fff; padding:12px 20px; text-decoration:none; border-radius:5px;">
				Activate Account
			</a>

			<p style="font-size:14px; margin-top:20px; color:#666;">
				If the button doesn't work, use this link:<br>
				<a href="%s">%s</a>
			</p>
		</div>
	</div>
	`, name, email, date, link, link, link)
}
