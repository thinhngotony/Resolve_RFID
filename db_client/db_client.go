package db_client

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", Username, Password, Hostname, Dbname)
}

func DbConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn(""))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return nil, err
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+Dbname)
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return nil, err
	}
	no, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when fetching rows", err)
		return nil, err
	}
	log.Printf("rows affected %d\n", no)
	db.Close()

	db, err = sql.Open("mysql", dsn(Dbname))
	if err != nil {
		log.Printf("Error %s when opening database", err)
		return nil, err
	}

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging database", err)
		return nil, err
	}
	log.Printf("Verified connection from %s database with Ping\n", Dbname)
	return db, nil
}

func ConvertFromRFID(db *sql.DB, rfid string) (string, string, string, bool, error) {
	log.Printf("Getting JAN code")
	query := `select drgm_jan, drgm_jan2, status from Covert_RFID_JANCODE where drgm_rfid_cd = ?;`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return "", "", "", false, err
	}
	defer stmt.Close()
	var jan_code_1, jan_code_2, status string
	row := stmt.QueryRowContext(ctx, rfid)
	if err := row.Scan(&jan_code_1, &jan_code_2, &status); err != nil {

		return "", "", "", false, err
	}
	return jan_code_1, jan_code_2, status, true, nil

}

func ConvertFromJan1(db *sql.DB, jancode_1 string) ([]string, bool, error) {
	log.Printf("Getting list of RFID from Jan code 1")
	query := `select drgm_rfid_cd from Covert_RFID_JANCODE where drgm_jan = ?;`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return nil, false, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, jancode_1)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	var rfid_list []string
	for rows.Next() {
		var prd string
		if err := rows.Scan(&prd); err != nil {
			return nil, false, err
		}
		rfid_list = append(rfid_list, prd)
	}

	if err := rows.Err(); err != nil {
		return nil, false, err
	}

	if len(rfid_list) == 0 {
		return nil, false, err
	}

	return rfid_list, true, nil
}

func ConvertFromJan2(db *sql.DB, jancode_2 string) ([]string, bool, error) {
	log.Printf("Getting list of RFID from Jan code 2")
	query := `select drgm_rfid_cd from Covert_RFID_JANCODE where drgm_jan2 = ?;`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return nil, false, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, jancode_2)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	var rfid_list []string
	for rows.Next() {
		var items string
		if err := rows.Scan(&items); err != nil {
			return nil, false, err
		}
		rfid_list = append(rfid_list, items)
	}

	if err := rows.Err(); err != nil {
		return nil, false, err
	}

	if len(rfid_list) == 0 {
		return nil, false, err
	}

	return rfid_list, true, nil
}

func UpdateLog(db *sql.DB, create_date string, store_code string, rfid string, mode string) error {
	log.Printf("Importing data to drfid_taglog table")
	query := `INSERT INTO drfid_taglog(dt_create_date, dt_store_code, dt_rfid, dt_mode ) VALUES (?, ?, ?, ?);`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, create_date, store_code, rfid, mode)
	if err != nil {
		log.Printf("Error %s when inserting row into log table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d rows created ", rows)
	return nil
}
