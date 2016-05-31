// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package lib

import (
	"errors"
	"strings"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

)

const (
	QUEUE_SOURCE_METHOD_API = "api"
	QUEUE_SOURCE_METHOD_FILE = "file"
	QUEUE_SOURCE_METHOD_MYSQL = "mysql"
)

type QueueSource struct {
	config *QueueSourceConfig
}

//queue data source fetch, SQL or restful API
type QueueSourceConfig struct {
	Method   string
	//sql dsn config
	MysqlDsn string
	//ApiPrefix+value config
	ApiPrefix string
	//Value for specific method
	Value string
}

// Construct a new QueueSource, need no pointer
func NewQueueSource(queue string, config QueueSourceConfig) (*QueueSource, error) {
	conNew:=&(config)
	conNew.Value=queue

	return NewQueueSourceByConfig(conNew)
}

func NewQueueSourceByConfig(config *QueueSourceConfig) (*QueueSource, error) {
	if config.Method==QUEUE_SOURCE_METHOD_API ||
	config.Method==QUEUE_SOURCE_METHOD_MYSQL ||
	config.Method==QUEUE_SOURCE_METHOD_FILE {

		if strings.Trim(config.Value, " \t")=="" {
			return nil, errors.New("QueueSourceConfig Vaule field empty.")
		}
	}else{
		return nil, errors.New("Unsupport QueueSource method.")
	}

	return &QueueSource{config:config}, nil
}


//use cache first, update when needed
func (qs *QueueSource) GetData()(list []string, err error) {
	if qs.config.Method==QUEUE_SOURCE_METHOD_API{
		list, err=qs.geneApiSouce()
	}else if qs.config.Method==QUEUE_SOURCE_METHOD_MYSQL{
		list, err=qs.geneMysqlSouce()
	}else if qs.config.Method==QUEUE_SOURCE_METHOD_FILE{
		list, err=qs.geneFileSouce()
	}else{
		return nil, errors.New("Unsupport QueueSource method.")
	}

	return
}

func (qs *QueueSource) geneMysqlSouce()(list []string, err error) {
	db, err := sql.Open("mysql", qs.config.MysqlDsn)
	if err != nil {
		return nil, errors.New("Error when sql.Open(): "+err.Error())
	}
	defer db.Close()

	var PushID string
	var PushList []string
	rows, err := db.Query(qs.config.Value)
	if err != nil {
		return nil, errors.New("Error when db.Query: "+err.Error())
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&PushID)
		if err != nil {
			return nil, errors.New("Error when rows.Scan(&PushID): "+err.Error())
		}
		PushList=append(PushList, PushID)
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.New("Error when rows.Err(): "+err.Error())
	}

	return PushList, nil
}

func (qs *QueueSource) geneApiSouce()(list []string, err error) {


	return
}

func (qs *QueueSource) geneFileSouce()(list []string, err error) {
	return
}

func (qs *QueueSource) Cache() (bool, error){
	return false, nil
}

func (qs *QueueSource) Update() (bool, error) {
	return false, nil
}





