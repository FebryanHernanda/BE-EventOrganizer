package templates

import "fmt"

func ActivationTemplate(name, email, date, link string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">

<body style="margin:0; padding:0; background:#f4f4f4; font-family:Arial, sans-serif;">
	<table width="100%%" cellpadding="0" cellspacing="0" style="background:#f4f4f4; padding:20px;">
		<tr>
			<td align="center">
				<table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff; border-radius:8px; padding:20px;">
					<tr>
						<td>
							<h2 style="color:#222; margin-bottom:10px;">Welcome, %s!</h2>
							<p style="color:#333; font-size:15px;">
								Thank you for registering. Below are your details:
							</p>

							<table style="margin:15px 0; font-size:14px;">
								<tr><td><strong>Email:</strong></td><td>%s</td></tr>
								<tr><td><strong>Registered At:</strong></td><td>%s</td></tr>
							</table>

							<p style="color:#333; font-size:15px; margin-bottom:20px;">
								Click the button below to activate your account:
							</p>

							<table cellpadding="0" cellspacing="0" style="margin:20px 0;">
								<tr>
									<td align="center">
										<a href="%s" 
											style="background:#007bff; color:#fff; padding:12px 20px; 
											text-decoration:none; border-radius:5px; font-size:16px;">
											Activate Account
										</a>
									</td>
								</tr>
							</table>

							<p style="font-size:13px; color:#666;">
								If the button doesn't work, use this link:<br>
								<a href="%s" style="color:#007bff;">%s</a>
							</p>

							<p style="font-size:12px; color:#999; margin-top:30px;">
								© %s Event Organizer. All rights reserved.
							</p>
						</td>
					</tr>
				</table>
			</td>
		</tr>
	</table>
</body>

</html>
	`, name, email, date, link, link, link, date[:4])
}
