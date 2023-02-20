package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/labstack/echo/v4"
)

func main() {
	// Create a new Echo instance
	e := echo.New()

	// Create a new SNS client
	snsClient := sns.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-2"),
	})))

	// Handler for push notifications
	e.POST("/push", func(c echo.Context) error {
		// Get the user's device token from the request
		// token := c.FormValue("token")

		// Parse the JSON message from the request body
		message := map[string]interface{}{
			"default":      "You got a message",
			"APNS_SANDBOX": `{"aps":{"alert":"You have a new message"}}`,
			"APNS":         `{"aps":{"alert":"You have a new message"}}`,
		}

		jsonMessage, err := json.Marshal(message)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to serialize message")
		}

		output, err := snsClient.Publish(&sns.PublishInput{
			TargetArn:        aws.String("arn:aws:sns:ap-northeast-2:117371106642:notification"),
			Message:          aws.String(string(jsonMessage)),
			MessageStructure: aws.String("json"),
			MessageAttributes: map[string]*sns.MessageAttributeValue{
				"MY.SNS.MOBILE.APNS.PUSH_TYPE": {
					DataType:    aws.String("String"),
					StringValue: aws.String("alert"),
				},
				"MY.SNS.MOBILE.APNS.TTL": {
					DataType:    aws.String("Number"),
					StringValue: aws.String("3600"),
				},
				"MY.SNS.MOBILE.APNS.PRIORITY": {
					DataType:    aws.String("String"),
					StringValue: aws.String("high"),
				},

				"MY.SNS.MOBILE.APNS_SANDBOX.PUSH_TYPE": {
					DataType:    aws.String("String"),
					StringValue: aws.String("alert"),
				},
				"MY.SNS.MOBILE.APNS_SANDBOX.TTL": {
					DataType:    aws.String("Number"),
					StringValue: aws.String("3600"),
				},
				"MY.SNS.MOBILE.APNS_SANDBOX.PRIORITY": {
					DataType:    aws.String("String"),
					StringValue: aws.String("high"),
				},
			},
		})

		if err != nil {
			println(err.Error())
			return c.String(500, "Failed to send push notification")
		}
		println(output.String())
		return c.String(200, "Push notification sent successfully")
	})

	e.GET("/subscribe", func(c echo.Context) error {
		body, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			// Handle error
			println(err.Error())
			return c.NoContent(http.StatusInternalServerError)
		}

		// Extract SubscribeURL from the request body
		var message map[string]interface{}
		err = json.Unmarshal(body, &message)
		if err != nil {
			// Handle error
			println(err.Error())
			return c.NoContent(http.StatusInternalServerError)
		}
		subscribeURL := message["SubscribeURL"].(string)

		// Confirm the subscription by sending an HTTP GET request to the SubscribeURL
		resp, err := http.Get(subscribeURL)
		if err != nil {
			// Handle error
			println(err.Error())
			return c.NoContent(http.StatusInternalServerError)
		}
		defer resp.Body.Close()

		// Return a 200 OK response to SNS to confirm the subscription
		return c.NoContent(http.StatusOK)
	})

	// Start the server
	e.Start(":3000")
}
