package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dimalewshin98-glitch/ShortyURL/internal/logger"
	models "github.com/dimalewshin98-glitch/ShortyURL/internal/model"
	"go.uber.org/zap"
)

type AuditorProducer struct {
	file   *os.File
	writer *bufio.Writer
}

func NewProducer(filename string) (*AuditorProducer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &AuditorProducer{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func (p *AuditorProducer) WriteData(writeData models.AuditData) error {
	data, err := json.Marshal(&writeData)
	if err != nil {
		return err
	}
	if _, err := p.writer.Write(data); err != nil {
		return err
	}
	if err := p.writer.WriteByte('\n'); err != nil {
		return err
	}
	return p.writer.Flush()
}

func (p *AuditorProducer) Close() error {
	return p.file.Close()
}

type Auditor interface {
	OnEvent(enevtType string, userID int, URL string)
}

type InfileAuditor struct {
	auditFile string
	producer  AuditorProducer
}

func NewInfileAuditor(auditFile string) (*InfileAuditor, error) {
	requestAuditor := InfileAuditor{}
	requestAuditor.auditFile = auditFile
	producer, err := NewProducer(auditFile)
	if err != nil {
		return nil, err
	}
	requestAuditor.producer = *producer
	return &requestAuditor, nil
}

func (a *InfileAuditor) OnEvent(enevtType string, userID int, URL string) {
	auditData := models.AuditData{
		TS:     123,
		Action: enevtType,
		UserID: strconv.Itoa(userID),
		URL:    URL,
	}
	a.producer.WriteData(auditData)
}

type RemoteAuditor struct {
	auditURL string
}

func NewRemoteAuditor(auditURL string) *RemoteAuditor {
	return &RemoteAuditor{auditURL: auditURL}
}

func (a *RemoteAuditor) OnEvent(enevtType string, userID int, URL string) {
	auditData := models.AuditData{
		TS:     123,
		Action: enevtType,
		UserID: strconv.Itoa(userID),
		URL:    URL,
	}
	reqData, err := json.Marshal(auditData)
	if err != nil {
		logger.Log.Error("Error send remote audit", zap.String("remote URL", a.auditURL), zap.String("error", err.Error()))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		a.auditURL,
		bytes.NewBuffer(reqData),
	)
	if err != nil {
		logger.Log.Error("Failed to create HTTP request to remote audit", zap.String("remote URL", a.auditURL), zap.String("error", err.Error()))
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Log.Error("Error on response from remote audit", zap.String("remote URL", a.auditURL), zap.String("error", err.Error()))
		return
	}
	if resp.StatusCode != http.StatusOK {
		logger.Log.Error("Unsuccess status code from remote audit",
			zap.String("remote URL", a.auditURL),
			zap.String("statusCode", resp.Status))
		return
	}
}
