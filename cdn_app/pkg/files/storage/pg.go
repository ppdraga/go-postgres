package storage

import (
	"fmt"
	"gorm.io/gorm"
)

type WithDB struct {
	db *gorm.DB
}

func New(db *gorm.DB) *WithDB {
	return &WithDB{
		db: db,
	}
}

func (wdb *WithDB) FindUserFiles(username string) (*[]string, error) {
	result := []string{}
	sql := `select file.name from userfile as uf left join file on uf.file_id = file.id 
		where uf.removed is null and uf.user_id in (select id from "user" where name = ?);`
	err := wdb.db.Raw(sql, username).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	fmt.Println(result)
	return &result, nil
}

func (wdb *WithDB) FindServerFiles(servername string) (*[]string, error) {
	result := []string{}
	sql := `SELECT file.name FROM serverfile AS sf left join file on sf.file_id = file.id 
		where sf.removed is null and sf.server_id in (select id from "server" s  where s.name = ?);`
	err := wdb.db.Raw(sql, servername).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	fmt.Println(result)
	return &result, nil
}

func (wdb *WithDB) FindAreaFiles(areaname string) (*[]string, error) {
	result := []string{}
	sql := `select distinct file.name FROM serverfile AS sf left join file on sf.file_id = file.id 
		where sf.removed is null and sf.server_id in (select id from "server" s  
		where s.area_id in (select id from area where name = ?));`
	err := wdb.db.Raw(sql, areaname).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	fmt.Println(result)
	return &result, nil
}