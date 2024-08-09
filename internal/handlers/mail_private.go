package handlers

import (
	"fmt"

	config "iteasy.wrappedAnsible/configs"
	"iteasy.wrappedAnsible/pkg/utils"
)

// mail 메일 인증
func _sendVerificationEmail(to, token string) error {
	subject := "Email Verification"
	verificationLink := fmt.Sprintf("%s/verify?token=%s", config.GlobalConfig.Default.Host, token)
	mailBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<body>
			<p>Please verify your email using the following link:</p>
			<p><a href="%s">Verify Email</a></p>
		</body>
		</html>`, verificationLink)

	if err := utils.SendEmail(to, subject, mailBody); err != nil {
		return err
	}
	return nil
}

func _sendResetPasswordEmail(to, tempPassword string) error {
	subject := "Password Reset"
	mailBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<body>
			<p>Your temporary password is: <b>%s</b></p>
		</body>
		</html>`, tempPassword)

	if err := utils.SendEmail(to, subject, mailBody); err != nil {
		return err
	}
	return nil
}
