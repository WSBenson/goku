package kafku

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/WSBenson/goku/internal"
	"github.com/tidwall/gjson"
)

// PublishFighterMessage publishes an event message when a fighter has been added to the ES index
func PublishFighterMessage(p sarama.AsyncProducer, pCount int, topic string, zfighter string) {
	internal.Logger.Info().Msg("Writing this message to the producer: " + zfighter)
	p.Input() <- &sarama.ProducerMessage{Topic: topic, Value: sarama.ByteEncoder(zfighter)}
}

// GenerateUniqueRandomMessage is a test-only function for creating a unique mock GM data event message
func GenerateUniqueRandomMessage() (message string, err error) {
	files, err := ioutil.ReadDir("/etc/static")
	if err != nil {
		internal.Logger.Fatal().Err(err).Msg("failed to open utils/static dirs")
	}
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	// Randomly grab one of these GM Data events as a string
	eventIndex := r1.Intn(len(files))
	selectedFileName := files[eventIndex].Name()

	fileContent, err := ioutil.ReadFile("/etc/static/" + selectedFileName)
	// Append on a timestamp so every message is guaranteed to be unique
	newID := "\"id\": \"" + time.Now().String() + "\""
	oldID := "\"id\": \"\""
	editedFileContent := strings.ReplaceAll(string(fileContent), oldID, newID)
	if err != nil {
		return "", err
	}
	return editedFileContent, err
}

// GenerateUniqueRandomMessages is a test-only function for creating a unique mock GM data event message
func GenerateUniqueRandomMessages(numMessages int) (messages []string, err error) {
	for i := 0; i < numMessages; i++ {
		msg, err := GenerateUniqueRandomMessage()
		messages = append(messages, msg)
		if err != nil {
			return messages, err
		}
	}
	return
}

// PublishRandomMessages is a test-only function for publishing mock GM data event messages
func PublishRandomMessages(p sarama.AsyncProducer, pCount int, topic string) (err error) {
	messages, err := GenerateUniqueRandomMessages(pCount)
	if err != nil {
		return err
	}
	for i := 0; i < len(messages); i++ {
		p.Input() <- &sarama.ProducerMessage{Topic: topic, Value: sarama.ByteEncoder(messages[i])}
	}
	return
}

// RandomKafkaTopics creates random strings to be used for kafka topic names
func RandomKafkaTopics(numString int, stringLength int) (names []string) {
	names = make([]string, numString)
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i := 0; i < numString; i++ {
		b := make([]byte, stringLength)
		for j := range b {
			b[j] = letterBytes[r1.Intn(len(letterBytes))]
		}
		names[i] = string(b)
	}
	return names
}

// RandomIndexName creates random string
func RandomIndexName(stringLength int) (name string) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	letterBytes := "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, stringLength)
	for j := range b {
		b[j] = letterBytes[r1.Intn(len(letterBytes))]
	}
	name = string(b)
	return
}

// AcmExists checks is the acm is not nil or blank
func AcmExists(msg *sarama.ConsumerMessage) (acmExists bool) {
	acm := gjson.Parse(string(msg.Value)).Get("originalobjectpolicy").String()
	acmExists = strings.TrimSpace(acm) != ""
	return
}

// CleanupText cleans out bad chars
func CleanupText(text string) (cleanedText []byte, err error) {
	textBytes := []byte(text)
	var character rune
	// Space's ascii character decimal value is 32
	space := rune(32)
	// 0 to 31 (inclusive), 92, 127 to 255 (inclusive)
	// https://www.asciitable.com/
	for i := 0; i <= 31; i++ {
		character = rune(i)
		textBytes = bytes.ReplaceAll(textBytes, []byte(string(character)), []byte(string(space)))
	}
	// Spacy does not like backslash (which has a ascii decimal value of 92)
	character = rune(92)
	textBytes = bytes.ReplaceAll(textBytes, []byte(string(character)), []byte(string(space)))

	// Also eliminate these ascii character
	for i := 127; i <= 255; i++ {
		character = rune(i)
		textBytes = bytes.ReplaceAll(textBytes, []byte(string(character)), []byte(string(space)))
	}
	return textBytes, err
}
