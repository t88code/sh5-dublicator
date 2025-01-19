package main

import (
	"context"
	"github.com/t88code/sh5-dublicator/domain"
	"github.com/t88code/sh5-dublicator/internal/comparer"
	"github.com/t88code/sh5-dublicator/internal/getter"
	"github.com/t88code/sh5-dublicator/internal/saver"
	"github.com/t88code/sh5-dublicator/pkg/config"
	"github.com/t88code/sh5-dublicator/pkg/sh5api"
	"log/slog"
	"os"
)

func checkError(err error) {
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	var err error
	ctx := context.Background()
	cfg, err := config.NewConfig()
	checkError(err)

	if _, err = os.Stat("./json"); os.IsNotExist(err) {
		err = os.Mkdir("json", 0777)
		checkError(err)
	}

	// создаем saver src
	svSrc, err := saver.New(logger, "src")
	if err != nil {
		return
	}

	// создаем saver dst
	svDst, err := saver.New(logger, "dst")
	if err != nil {
		return
	}

	// создаем клиентов для подключения к базам
	sh5ConfigSrc := &sh5api.Config{
		HttpClient: nil,
		BaseURL:    cfg.SH5SRC.BaseURL,
		Username:   cfg.SH5SRC.Username,
		Password:   cfg.SH5SRC.Password,
		DebugLog:   cfg.SH5SRC.DebugLog,
		OnQuery:    svSrc.OnQuery,
		OnReply:    svSrc.OnReply,
		Logger:     logger,
	}

	sh5ClientSrc, err := sh5api.New(sh5ConfigSrc)
	checkError(err)

	sh5ConfigDst := &sh5api.Config{
		HttpClient: nil,
		BaseURL:    cfg.SH5DST.BaseURL,
		Username:   cfg.SH5DST.Username,
		Password:   cfg.SH5DST.Password,
		DebugLog:   cfg.SH5DST.DebugLog,
		OnQuery:    svDst.OnQuery,
		OnReply:    svDst.OnReply,
		Logger:     logger,
	}

	sh5ClientDst, err := sh5api.New(sh5ConfigDst)
	checkError(err)

	// создаем геттеры
	gtSrc, err := getter.New(sh5ClientSrc)
	checkError(err)

	gtDst, err := getter.New(sh5ClientDst)
	checkError(err)

	// создаем comparer
	cmr, err := comparer.New(logger, gtSrc, gtDst)
	checkError(err)

	///////////////////
	procs := []*domain.ProcSync{
		{
			Name: sh5api.GGroupsTree,
			Head: sh5api.HeadCodeGGROUP,
			Original: []sh5api.FieldPath{
				sh5api.FIELD_4_GUID, // 4 GUID Товарнoй группы
				sh5api.FIELD_209_1,  // 209#1\1 Товарная группа - предок RID предка
				sh5api.FIELD_3,      // 3 Наименование
				//sh5api.FIELD_6,      // 6 Атрибуты
			},
		},
		//{
		//	Name:     "",
		//	Head:     "",
		//	Original: nil,
		//},
		//sh5api.GGroups,
		//sh5api.GGroupsTree,
	}

	compareResult, err := cmr.CompareDictionary(ctx, procs)
	checkError(err)

	if len(compareResult) == 0 {
		slog.Info("Обновление справочников не требуется")
		return
	}

	/////////////////////////

	////////////////////////////////////////

	////////////////////////////////////////
	//dupl, err := sh5dublicator.New(
	//	sh5ClientSrc,
	//	sh5ClientDst,
	//	procNamesForSync)
	//if err != nil {
	//	slog.Error(err.Error())
	//	return
	//}
	//
	//err = dupl.CopyDictionary(ctx)
	//if err != nil {
	//	slog.Error(err.Error())
	//	return
	//}

	////////////////////////////////////////

	//err = updater.Update(ctx, sh5ClientSrc, sh5ClientDst, compareResults, procNamesForSync)
	//if err != nil {
	//	slog.Error(err.Error())
	//	return
	//}

	//for _, ref := range refs {
	//	var req, rep interface{}
	//
	//	match, err := comparer.Compare(ctx, string(ref), rep)
	//	if err != nil {
	//		slog.Error(err.Error())
	//		return
	//	}
	//
	//	if !match {
	//		slog.Info(fmt.Sprintf("Необходимо обновление (%s)", ref))
	//		///////////////////////
	//		insGGroupReq := new(sh5api.InsGGroupReq)
	//		insGGroupRep, err := sh5Client.InsertGGroup(ctx, insGGroupReq)
	//		if err != nil {
	//			slog.Error(err.Error())
	//			return
	//		}
	//
	//		bytes, err := json.Marshal(insGGroupRep)
	//		if err != nil {
	//			return
	//		}
	//		fmt.Printf("%s", pretty.Pretty(bytes))
	//
	//		///////////////////////
	//	}
	//
	//	err = saver.SaveReqRep(context.Background(), string(ref), req, rep)
	//	if err != nil {
	//		slog.Error(err.Error())
	//		return
	//	}
	//
	//}

}
