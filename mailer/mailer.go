package mailer

import(
	"os"
	"github.com/kataras/go-mailer"
)

func SendMail(to []string, subject string, content string) map[string]interface{}{
	config := mailer.Config{
        Host:     os.Getenv("MAILER_HOST"),
        Username: os.Getenv("MAILER_USERNAME"),
        Password: os.Getenv("MAILER_PASSWORD"),
        FromAddr: "admin@paingha.me",
        Port:     587,
        // Enable UseCommand to support sendmail unix command,
        // if this field is true then Host, Username, Password and Port are not required,
        // because these info already exists in your local sendmail configuration.
        //
        // Defaults to false.
        UseCommand: false,
	}
	// initalize a new mail sender service.
    sender := mailer.New(config)
	// send the e-mail.
    err := sender.Send(subject, content, to...)

    if err != nil {
		println("error while sending the e-mail: " + err.Error())
		return map[string]interface{}{
			"message": "An error occured while sending email",
			"status": false,
			"error": err.Error(),
		}
	}
	return map[string]interface{}{
		"message": "Email sent successfully",
		"status": true,
		"error": nil,
	}
}

//Set func to accept html file path for email template