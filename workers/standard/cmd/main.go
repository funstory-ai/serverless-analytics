package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	KAFKA_USER   string
	KAFKA_PASSWD string
	KAFKA_URL    string
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	KAFKA_PASSWD = os.Getenv("KAFKA_PASSWD")
	KAFKA_USER = os.Getenv("KAFKA_USER")
	KAFKA_URL = os.Getenv("KAFKA_URL")
}

type LogRequest struct {
	Topic     string   `json:"topic"`
	Logs      []string `json:"logs"`
	Timestamp int64    `json:"timestamp"`
}

func main() {
	router := gin.Default()
	router.POST("/collect", handleLog)
	router.Run(":8080")
}

func handleLog(c *gin.Context) {
	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	var logRequest LogRequest
	err = json.Unmarshal(b, &logRequest)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	// sendToKafka(logRequest.Topic, logRequest.Logs)
	if err := postToKafka(logRequest.Topic, logRequest.Logs); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.String(200, "ok")
}

func postToKafka(topic string, logs []string) error {

	url := KAFKA_URL + "/produce/" + topic
	fmt.Println(url)
	var payload []map[string]string
	for _, log := range logs {
		payload = append(payload, map[string]string{
			"value": log,
		})
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.SetBasicAuth(KAFKA_USER, KAFKA_PASSWD)
	fmt.Println(KAFKA_USER)
	fmt.Println(KAFKA_PASSWD)
	// req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// You might want to check the response status code here
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d %s", resp.StatusCode, string(b))
	}
	return nil
}
