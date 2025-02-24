package main

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/t88code/sh5-dublicator/dbsheme"
	"github.com/t88code/sh5-dublicator/domain"
	"github.com/t88code/sh5-dublicator/internal/comparer"
	"github.com/t88code/sh5-dublicator/internal/getter"
	"github.com/t88code/sh5-dublicator/internal/repository/sqlite"
	"github.com/t88code/sh5-dublicator/internal/saver"
	"github.com/t88code/sh5-dublicator/internal/syncer"
	"github.com/t88code/sh5-dublicator/internal/utils"
	"github.com/t88code/sh5-dublicator/pkg/config"
	"github.com/t88code/sh5-dublicator/pkg/sh5api"
	"log/slog"
	"os"
	"time"
)

func checkError(err error) {
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

const DbSqliteFileName = "db.db"

func main() {
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelDebug)

	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level:     lvl,
				AddSource: false},
		))

	var err error
	cfg, err := config.NewConfig()
	checkError(err)

	databaseExists := utils.Exists(DbSqliteFileName)
	if databaseExists {
		err = os.Remove(DbSqliteFileName)
		checkError(err)
	}

	// init DB sqlite connection
	db, err := sqlx.Connect("sqlite3", DbSqliteFileName)
	checkError(err)
	defer db.Close()

	if databaseExists {
		db.MustExec(dbsheme.GGroupsSrc)
		db.MustExec(dbsheme.GGroupsDst)
	}

	if _, err = os.Stat("./json"); os.IsNotExist(err) {
		err = os.Mkdir("json", 0777)
		checkError(err)
	}

	folderNameForHttpRequest := time.Now().Format("2006-01-02 15-04-05")
	if _, err = os.Stat(fmt.Sprintf("./json/%s", folderNameForHttpRequest)); os.IsNotExist(err) {
		err = os.Mkdir(fmt.Sprintf("./json/%s", folderNameForHttpRequest), 0777)
		checkError(err)
	}

	// создаем saver src
	svSrc, err := saver.New(
		logger,
		"src",
		folderNameForHttpRequest)
	checkError(err)

	// создаем saver dst
	svDst, err := saver.New(
		logger,
		"dst",
		folderNameForHttpRequest)
	checkError(err)

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

	// вызываемые процедуры
	procs := []*domain.ProcSync{
		{
			Name: sh5api.GGroupsTree,
			Head: sh5api.HeadCodeGGROUP,
			OriginalForSaveToSql: []sh5api.FieldPath{
				sh5api.FIELD_1_RID,
				sh5api.FIELD_4_GUID,
				sh5api.FIELD_42,
				sh5api.FIELD_3_NAME,
				sh5api.FIELD_6_RidDst,
				sh5api.FIELD_209_1_RID_PARENT,
				sh5api.FIELD_209_4_GUID_PARENT,
				sh5api.FIELD_209_42,
				sh5api.FIELD_209_3_NAME_PARENT,
				sh5api.FIELD_209_6_RidDst_PARENT,
				sh5api.FIELD_239,
				sh5api.FIELD_106_1,
				sh5api.FIELD_106_3,
				sh5api.FIELD_6,
			},
			OriginalForInsertToDst: []sh5api.FieldPath{
				sh5api.FIELD_4_GUID,
				sh5api.FIELD_3_NAME,
				sh5api.FIELD_209_1_RID_PARENT,
			},
			OriginalForUpdateAttrInSrc: []sh5api.FieldPath{
				sh5api.FIELD_1_RID,
				sh5api.FIELD_4_GUID,
				sh5api.FIELD_209_1_RID_PARENT,
				sh5api.FIELD_3_NAME,
				sh5api.FIELD_6_RidDst,
			},
			OriginalForUpdateInDst: []sh5api.FieldPath{
				sh5api.FIELD_1_RID,
				sh5api.FIELD_4_GUID,
				sh5api.FIELD_209_1_RID_PARENT,
				sh5api.FIELD_3_NAME,
				sh5api.FIELD_6_RidDst,
			},
			OriginalForDeleteInDst: []sh5api.FieldPath{
				sh5api.FIELD_1_RID,
				sh5api.FIELD_4_GUID,
			},
			IsActionDoIt: true,
		},
	}

	// создаем репозиторий для сохранения справочников в DB SQLite
	gGroupsRepository := sqlite.NewGGroupsRepository(db, logger)

	// создаем comparer
	cmr, err := comparer.New(logger, gtSrc, gtDst, gGroupsRepository)
	checkError(err)

	ctx := context.Background()

	// создаем syncer
	scr, err := syncer.New(
		logger,
		cmr,
		sh5ClientSrc,
		sh5ClientDst,
		gGroupsRepository)
	checkError(err)

	// копируем справочники
	err = scr.SyncDictionary(ctx, procs)
	checkError(err)
}
