package repository

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"os"

	models "github.com/dimalewshin98-glitch/ShortyURL/internal/model"
	"github.com/google/uuid"
)

type Producer struct {
	file   *os.File
	writer *bufio.Writer
}

func NewProducer(filename string) (*Producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func (p *Producer) WriteElement(element *URLelement) error {
	data, err := json.Marshal(&element)
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

func (p *Producer) Close() error {
	return p.file.Close()
}

type Consumer struct {
	file *os.File
}

func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file: file,
	}, nil
}

func (c *Consumer) ReadElement(URLId string) (*URLelement, error) {
	c.file.Seek(0, 0)
	scanner := bufio.NewScanner(c.file)
	for scanner.Scan() {
		var element URLelement
		line := scanner.Bytes()
		err := json.Unmarshal(line, &element)
		if err != nil {
			return nil, err
		}
		if element.ShortUrl == URLId {
			return &element, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return nil, errors.New("URL not found in file repository")
}

func (c *Consumer) ReadElements() ([]URLelement, error) {
	var result []URLelement
	c.file.Seek(0, 0)
	scanner := bufio.NewScanner(c.file)
	for scanner.Scan() {
		var element URLelement
		line := scanner.Bytes()
		err := json.Unmarshal(line, &element)
		if err != nil {
			return nil, err
		}
		result = append(result, element)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Consumer) Close() error {
	return c.file.Close()
}

type URLelement struct {
	UUID        string `json:"uuid"`
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
	UserID      int    `json:"user_id"`
}

type InfileRepository struct {
	fileStoragePath string
	producer        Producer
	consumer        Consumer
}

func NewfileRepository(fileStoragePath string) (*InfileRepository, error) {
	producer, err := NewProducer(fileStoragePath)
	if err != nil {
		return nil, err
	}
	consumer, err := NewConsumer(fileStoragePath)
	if err != nil {
		return nil, err
	}
	return &InfileRepository{
		fileStoragePath: fileStoragePath,
		producer:        *producer,
		consumer:        *consumer,
	}, err
}

func (r *InfileRepository) Ping(ctx context.Context) error {
	return nil
}

func (r *InfileRepository) Store(ctx context.Context, userID int, urlID string, URL string) (string, error) {
	element := URLelement{
		UUID:        uuid.New().String(),
		ShortUrl:    urlID,
		OriginalUrl: URL,
		UserID:      userID,
	}
	r.producer.WriteElement(&element)
	return urlID, nil
}

func (r *InfileRepository) Get(ctx context.Context, urlID string) (string, error) {
	element, err := r.consumer.ReadElement(urlID)
	if err != nil {
		return "", err
	}
	URL := element.OriginalUrl
	return URL, nil
}

func (r *InfileRepository) GetUsersID(ctx context.Context) ([]int, error) {
	elements, err := r.consumer.ReadElements()
	if err != nil {
		return nil, err
	}
	uniqueIDs := make(map[int]bool)
	var result []int
	for _, urlInfo := range elements {
		if !uniqueIDs[urlInfo.UserID] {
			uniqueIDs[urlInfo.UserID] = true
			result = append(result, urlInfo.UserID)
		}
	}
	return result, nil
}

func (r *InfileRepository) GetUserUrls(ctx context.Context, userID int) (models.ApiUserUrlsRes, error) {
	elements, err := r.consumer.ReadElements()
	if err != nil {
		return nil, err
	}
	var userURLs models.ApiUserUrlsRes
	for _, urlInfo := range elements {
		var userURL models.UserUrlRes
		if urlInfo.UserID == userID {
			userURL.OriginalURL = urlInfo.OriginalUrl
			userURL.ShortURL = urlInfo.ShortUrl
			userURLs = append(userURLs, userURL)
		}
	}
	return userURLs, nil
}
