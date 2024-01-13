package logger

import (
	"context"
	"os"
	"regexp"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

const (
	bufferSize     = 2048
	delayByTimeout = time.Millisecond * 100

	testMessage       = "Test message"
	logMessagePattern = `\[\d{4}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}\.\d{3}\]\s\[(\w+)\]\s(.*)`
)

var regExp *regexp.Regexp

type MessageData struct {
	buffer []byte
	err    error
}

func init() {
	regExp = regexp.MustCompile(logMessagePattern)
}

// Прочитать сообщение из лога.
func readMessage(pipe *os.File) MessageData {
	var (
		messageCh = make(chan MessageData)
		message   MessageData
		wg        sync.WaitGroup
	)

	ctx, cancel := context.WithTimeout(context.Background(), delayByTimeout)
	defer cancel()

	wg.Add(1)

	go func() {
		defer close(messageCh)

		buffer := make([]byte, bufferSize)
		read, err := pipe.Read(buffer)

		messageCh <- MessageData{
			buffer: buffer[:read],
			err:    err,
		}
	}()

	go func() {
		defer wg.Done()

		select {
		case <-ctx.Done():
			cancel()

		case message = <-messageCh:
		}
	}()
	wg.Wait()

	return message
}

type loggerSuite struct {
	suite.Suite
	oldStdout *os.File
	reader    *os.File
	writer    *os.File
	err       error
}

func (s *loggerSuite) SetupTest() {
	s.oldStdout = os.Stdout
	s.reader, s.writer, s.err = os.Pipe()

	s.Require().NoError(s.err)
	os.Stdout = s.writer
}

func (s *loggerSuite) TeardownTest() {
	err := s.reader.Close()
	s.Require().NoError(err)

	err = s.writer.Close()
	s.Require().NoError(err)

	os.Stdout = s.oldStdout
}

// Проверка что лог записан без ошибок.
func (s *loggerSuite) checkLogMessage() {
	msg := readMessage(s.reader)
	s.Require().NoError(msg.err)
	s.Require().NotNil(msg.buffer)
	s.Require().True(regExp.Match(msg.buffer))
}

// Проверка что лог отсутствует.
func (s *loggerSuite) checkEmptyLogMessage() {
	msg := readMessage(s.reader)
	s.Require().NoError(msg.err)
	s.Require().Nil(msg.buffer)
}

func (s *loggerSuite) TestErrorLevel() {
	logg := New(errorTitle)

	logg.Error(testMessage)
	s.checkLogMessage()

	logg.Warning(testMessage)
	s.checkEmptyLogMessage()

	logg.Info(testMessage)
	s.checkEmptyLogMessage()

	logg.Debug(testMessage)
	s.checkEmptyLogMessage()
}

func (s *loggerSuite) TestWarningLevel() {
	logg := New(warningTitle)

	logg.Error(testMessage)
	s.checkLogMessage()

	logg.Warning(testMessage)
	s.checkLogMessage()

	logg.Info(testMessage)
	s.checkEmptyLogMessage()

	logg.Debug(testMessage)
	s.checkEmptyLogMessage()
}

func (s *loggerSuite) TestInfoLevel() {
	logg := New(infoTitle)

	logg.Error(testMessage)
	s.checkLogMessage()

	logg.Warning(testMessage)
	s.checkLogMessage()

	logg.Info(testMessage)
	s.checkLogMessage()

	logg.Debug(testMessage)
	s.checkEmptyLogMessage()
}

func (s *loggerSuite) TestDebugLevel() {
	logg := New(debugTitle)

	logg.Error(testMessage)
	s.checkLogMessage()

	logg.Warning(testMessage)
	s.checkLogMessage()

	logg.Info(testMessage)
	s.checkLogMessage()

	logg.Debug(testMessage)
	s.checkLogMessage()
}

func TestLoggerSuite(t *testing.T) {
	suite.Run(t, new(loggerSuite))
}
