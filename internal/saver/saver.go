package saver

import (
	"fmt"
	"github.com/t88code/sh5-dublicator/internal/helper"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type saver struct {
	logger     *slog.Logger
	postfix    string
	folderName string
}

func New(logger *slog.Logger, postfix string, folderName string) (Saver, error) {
	return &saver{
		logger:     logger,
		postfix:    postfix,
		folderName: folderName,
	}, nil
}

func (s *saver) OnQuery(req *http.Request, body []byte, reqName string) {
	postfix := "req"
	if s.postfix != "" {
		postfix = fmt.Sprintf("%s_%s", postfix, s.postfix)
	}
	postfix = fmt.Sprintf("%s_%d", postfix, time.Now().Unix())

	err := s.SaveDataToFile(helper.GetPath(s.folderName, reqName, postfix), body)
	if err != nil {
		s.logger.Error(err.Error())
	}
}

func (s *saver) OnReply(res *http.Response, body []byte, reqName string) {
	postfix := "rep"
	if s.postfix != "" {
		postfix = fmt.Sprintf("%s_%s", postfix, s.postfix)
	}
	postfix = fmt.Sprintf("%s_%d", postfix, time.Now().Unix())

	err := s.SaveDataToFile(helper.GetPath(s.folderName, reqName, postfix), body)
	if err != nil {
		s.logger.Error(err.Error())
	}
}

func (s *saver) SaveDataToFile(path string, data []byte) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			s.logger.Error(err.Error())
		}
	}(file)

	n, err := file.Write(data)
	if err != nil {
		return err
	}
	s.logger.Info(fmt.Sprintf("Записано в (%s) байтов (%d)", path, n))
	return nil
}
