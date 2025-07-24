package ses

import (
	"context"
	"errors"
	"os"

	ur "blog-api/repositories/user"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

var ErrDaemonAddrNotFound = errors.New("failed to lookup daemon address")
var ErrRegionNotFound = errors.New("failed to lookup aws region")

func SendEmail(input *sesv2.SendEmailInput) error {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		return ErrRegionNotFound
	}

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return err
	}

	client := sesv2.NewFromConfig(cfg)

	_, err = client.SendEmail(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func PrepareContactEmail(params *ur.UserSendEmailPost) (*sesv2.SendEmailInput, error) {
	daemon, ok := os.LookupEnv("DAEMON_ADDRESS")
	if !ok {
		return nil, ErrDaemonAddrNotFound
	}

	plainText := "This message was delivered on behalf of: " + params.From + "\n\n" +
		"EMAIL IS AS FOLLOWS:\n\n" +
		"--------------------\n\n" +
		params.Message + "\n\n" +
		"Respond to: " + params.From

	html := "<!DOCTYPE html>" +
		"<html>" +
		"<body>" +
		"<div><h3>This message was delivered on behalf of: " + params.From + "</h3></div>" +
		"<div><p>EMAIL IS AS FOLLOWS:</p></div>" +
		"<div>--------------------</div>" +
		"<div><p>" + params.Message + "</p></div>" +
		"</br>" +
		"<div><strong>" + "Respond to: " + params.From + "</strong></div>" +
		"</body>" +
		"</html>"

	return &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(daemon),
		Destination: &types.Destination{
			ToAddresses: []string{params.To},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data: aws.String(params.Subject),
				},
				Body: &types.Body{
					Html: &types.Content{
						Data: aws.String(html),
					},
					Text: &types.Content{
						Data: aws.String(plainText),
					},
				},
			},
		},
	}, nil
}

func PreparePasswordResetEmail(token, to string) (*sesv2.SendEmailInput, error) {
	url := "https://jonahbutler.dev/password-reset?resetToken=" + token

	daemon, ok := os.LookupEnv("DAEMON_ADDRESS")
	if !ok {
		return nil, ErrDaemonAddrNotFound
	}

	plainText := "A request to reset your password was submitted.\n\n" +
		"If you did not make this request, ignore this email as someone may have accidentally typed your email address.\n\n" +
		"To update your password please visit: \n\n" +
		url + "\n\n"

	html := "<div><h3>A request to reset your password was submitted.</h3></div>" +
		"<div><strong>If you did not make this request, ignore this email as someone may have accidentally typed your email address.</strong></div>" +
		"<div><p>To update your password please visit:</p></div>" +
		"<div><a href='" + url + "'" + ">" + url + "</a></div>"

	return &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(daemon),
		Destination: &types.Destination{
			ToAddresses: []string{to},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data: aws.String("Password Reset"),
				},
				Body: &types.Body{
					Html: &types.Content{
						Data: aws.String(html),
					},
					Text: &types.Content{
						Data: aws.String(plainText),
					},
				},
			},
		},
	}, nil
}
