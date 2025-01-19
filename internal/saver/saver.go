package saver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/t88code/sh5-dublicator/internal/helper"
	"log"
	"log/slog"
	"net/http"
	"os"
)

type saver struct {
	logger  *slog.Logger
	postfix string
}

func New(logger *slog.Logger, postfix string) (Saver, error) {
	return &saver{
		logger:  logger,
		postfix: postfix,
	}, nil
}

//func SaveAllReqRep(ctx context.Context,
//	reqRepMapSrc map[sh5api.HeadName]getter.ReqRep,
//	reqRepMapDst map[sh5api.HeadName]getter.ReqRep,
//	headNamesForSync []sh5api.HeadName) error {
//
//	for _, headName := range headNamesForSync {
//
//		reqRepSrc, ok := reqRepMapSrc[headName]
//		if !ok {
//			return fmt.Errorf("справочник (%s) не найден в базе источнике", headName)
//		}
//
//		reqRepDst, ok := reqRepMapDst[headName]
//		if !ok {
//			return fmt.Errorf("справочник (%s) не найден в базе получателе", headName)
//		}
//
//		err := SaveReqRep(ctx, fmt.Sprintf("src_%s", headName), reqRepSrc.Req, reqRepSrc.Rep)
//		if err != nil {
//			return err
//		}
//		err = SaveReqRep(ctx, fmt.Sprintf("dst_%s", headName), reqRepDst.Req, reqRepDst.Rep)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}

//func SaveReqRep(ctx context.Context, filename string, req interface{}, rep interface{}) error {
//	err := SaveJsonFile(ctx, helper.GetPath(filename, "req"), req)
//	if err != nil {
//		return err
//	}
//	err = SaveJsonFile(ctx, helper.GetPath(filename, "rep"), rep)
//	if err != nil {
//		return err
//	}
//	return nil
//}

func (s *saver) OnQuery(req *http.Request, body []byte, reqName string) {
	postfix := "req"
	if s.postfix != "" {
		postfix = fmt.Sprintf("%s_%s", postfix, s.postfix)
	}

	err := s.SaveDataToFile(helper.GetPath(reqName, postfix), body)
	if err != nil {
		s.logger.Error(err.Error())
	}
}

func (s *saver) OnReply(res *http.Response, body []byte, reqName string) {
	postfix := "rep"
	if s.postfix != "" {
		postfix = fmt.Sprintf("%s_%s", postfix, s.postfix)
	}

	err := s.SaveDataToFile(helper.GetPath(reqName, postfix), body)
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

func SaveJsonFile(ctx context.Context, path string, dataJsonStruct interface{}) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	dataJsonBytes, err := json.MarshalIndent(dataJsonStruct, "", "    ")
	if err != nil {
		return err
	}

	n, err := file.Write(dataJsonBytes)
	if err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("Записано в (%s) байтов (%d)", path, n))
	return nil
}
