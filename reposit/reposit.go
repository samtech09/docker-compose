package reposit

import (
	"sam-learn/docker-compose/model"
	"time"

	"github.com/jackc/pgx"
)

var (
	_pdb *pgx.ConnPool
)

const (
	//prepared statements const
	getdata = "pstmtgetdata"
)

//InitDB initialize database connection and prepare statements
func InitDB(host, dbname, dbuser, dbpwd string, port int) error {
	pgxConfig := pgx.ConnConfig{
		Host:     host,
		Port:     uint16(port),
		Database: dbname,
		User:     dbuser,
		Password: dbpwd,
	}
	var err error
	pgxConnPoolConfig := pgx.ConnPoolConfig{pgxConfig, 8, nil, 5 * time.Second}
	_pdb, err = pgx.NewConnPool(pgxConnPoolConfig)
	if err != nil {
		//log.Fatalf("Psql Connection error %v\n", err)
		return err
	}

	// setup prepared sql statement for later use
	sql := "select id, name, count(*) over() as cnt from demotable;"
	_, err = _pdb.Prepare(getdata, sql)
	if err != nil {
		//log.Fatalf("Failed preparing stmt: %v\n", err)
		return err
	}

	return nil
}

//GetData returns list of records from demotable
func GetData() ([]model.Demodata, error) {
	//get list of usres not submitted given test
	rows, err := _pdb.Query(getdata)
	if err != nil {
		return nil, err
	}
	// make sure to always close rows
	defer rows.Close()

	var list []model.Demodata
	count := 0
	i := 0
	for rows.Next() {
		// scan rows to struct
		u := model.Demodata{}
		err := rows.Scan(&u.ID, &u.Name, &count)
		if err != nil {
			return nil, err
		}

		if i == 0 {
			// instead of append, make slice of fixed length to avoid reallocations
			list = make([]model.Demodata, count)
		}
		list[i] = u
		i++
	}

	return list, nil
}
