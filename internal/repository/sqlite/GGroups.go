package sqlite

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/t88code/sh5-dublicator/domain"
	"log/slog"
)

type GGroupsRepository struct {
	DB     *sqlx.DB
	logger *slog.Logger
}

// NewGGroupsRepository will create an object that represent the article.Repository interface
func NewGGroupsRepository(conn *sqlx.DB, l *slog.Logger) *GGroupsRepository {
	return &GGroupsRepository{conn, l}
}

func (c *GGroupsRepository) insertOrUpdateByName(ctx context.Context, query string, ggroups *domain.GGroups) (err error) {
	result, err := c.DB.NamedExecContext(ctx, query, ggroups)
	if err != nil {
		return err
	}
	affect, err := result.RowsAffected()
	if err != nil {
	}
	if affect != 1 {
		err = fmt.Errorf("weird  Behavior. Total Affected: %d", affect)
		return
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		return
	}
	ggroups.ID = uint32(lastID)
	return
}

func (c *GGroupsRepository) GetByRid(ctx context.Context, tableName string, rid uint32) (ggroups []domain.GGroups, err error) {
	var query string
	if tableName == domain.TABLE_ggroups_src {
		query = `SELECT * FROM ggroups_src WHERE Rid = $1`
	} else {
		query = `SELECT * FROM ggroups_dst WHERE Rid = $1`
	}
	err = c.DB.SelectContext(ctx, &ggroups, query, rid)
	return
}

func (c *GGroupsRepository) GetByAttr(ctx context.Context, tableName string, attr uint32) (ggroups []domain.GGroups, err error) {
	var query string
	if tableName == domain.TABLE_ggroups_src {
		query = `SELECT * FROM ggroups_src WHERE Attr = $1`
	} else {
		query = `SELECT * FROM ggroups_dst WHERE Attr = $1`
	}
	err = c.DB.SelectContext(ctx, &ggroups, query, attr)
	return
}

func (c *GGroupsRepository) InsertAll(ctx context.Context, ggroup *domain.GGroups, tableName string) (err error) {
	var query string
	if tableName == domain.TABLE_ggroups_src {
		query =
			`INSERT INTO ggroups_src (
				Rid,
				Guid,
				FIELD_42,
				Name,
				Attr,
				Rid_PARENT,
				Guid_PARENT,
				FIELD_42_PARENT,
				Name_PARENT,
				Attr_PARENT,
				UserGroup,
				Rid_106,
				Name_106,
                Action         
			) VALUES (
				:Rid,
				:Guid,
				:FIELD_42,
				:Name,
				:Attr,
				:Rid_PARENT,
				:Guid_PARENT,
				:FIELD_42_PARENT,
				:Name_PARENT,
				:Attr_PARENT,
				:UserGroup,
				:Rid_106,
				:Name_106,
			    ''
		);`
	} else {
		query =
			`INSERT INTO ggroups_dst (
				Rid,
				Guid,
				FIELD_42,
				Name,
				Attr,
				Rid_PARENT,
				Guid_PARENT,
				FIELD_42_PARENT,
				Name_PARENT,
				Attr_PARENT,
				UserGroup,
				Rid_106,
				Name_106,
                Action
			) VALUES (
				:Rid,
				:Guid,
				:FIELD_42,
				:Name,
				:Attr,
				:Rid_PARENT,
				:Guid_PARENT,
				:FIELD_42_PARENT,
				:Name_PARENT,
				:Attr_PARENT,
				:UserGroup,
				:Rid_106,
				:Name_106,
			    ''
		);`
	}
	return c.insertOrUpdateByName(ctx, query, ggroup)
}

func (c *GGroupsRepository) Delete(ctx context.Context, ggroup *domain.GGroups, tableName string) (err error) {
	var query string
	if tableName == domain.TABLE_ggroups_src {
		query = `DELETE FROM ggroups_src WHERE Rid = :Rid`
	} else {
		query = `DELETE FROM ggroups_dst WHERE Rid = :Rid`
	}
	result, err := c.DB.NamedExecContext(ctx, query, ggroup)
	if err != nil {
		return
	}
	rowsAfected, err := result.RowsAffected()
	if err != nil {
		return
	}
	if rowsAfected != 1 {
		err = fmt.Errorf("weird  Behavior. Total Affected: %d", rowsAfected)
		return
	}
	return
}

func (c *GGroupsRepository) GetAll(ctx context.Context, tableName string) (ggroups []domain.GGroups, err error) {
	var query string
	if tableName == domain.TABLE_ggroups_src {
		query = `SELECT * FROM ggroups_src`
	} else {
		query = `SELECT * FROM ggroups_dst`
	}
	err = c.DB.SelectContext(ctx, &ggroups, query)
	return
}

func (c *GGroupsRepository) GetByGuid(ctx context.Context, tableName string, Guid string) (ggroups []domain.GGroups, err error) {
	var query string
	if tableName == domain.TABLE_ggroups_src {
		query = `SELECT * FROM ggroups_src WHERE Guid=$1`
	} else {
		query = `SELECT * FROM ggroups_dst WHERE Guid=$1`
	}
	err = c.DB.SelectContext(
		ctx,
		&ggroups,
		query,
		Guid)
	return
}

func (c *GGroupsRepository) UpdateAction(ctx context.Context, tableName string, ggroup *domain.GGroups) (err error) {
	var query string
	if tableName == domain.TABLE_ggroups_src {
		query = `UPDATE ggroups_src set Action=:Action WHERE Rid=:Rid`
	} else {
		query = `UPDATE ggroups_dst set Action=:Action WHERE Rid=:Rid`
	}

	return c.insertOrUpdateByName(ctx, query, ggroup)
}
