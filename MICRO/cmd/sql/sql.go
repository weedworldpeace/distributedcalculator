package sql

import (
	"context"
	"database/sql"
	"errors"
	"log"

	_ "github.com/mattn/go-sqlite3"
)


type DB struct {
	ctx context.Context
	db *sql.DB
}

var MyDB = DB{}

type User struct {
	Login string `json:"login"`
	Password string `json:"password"`
}

type Expression struct {
	Id int
	Status string
	Result float64
}

func (d *DB) InsertUser(u User) (int64, error) {
	var q = `INSERT INTO users (login, password) values ($1, $2)`

	res, err := d.db.ExecContext(d.ctx, q, u.Login, u.Password)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (d *DB) InsertExpression(user_id int64) (int64, error) {
	var q = `INSERT INTO expressions (status, result, user_id) values ($1, $2, $3)`
	res, err := d.db.ExecContext(d.ctx, q,"accepted", 0.0, user_id)
	if err != nil {
		log.Println(err)
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (d *DB) UpdateExpression(result float64, id int64) error {
	var q1 = `UPDATE expressions SET status = $1 WHERE id = $2`
	var q2 = `UPDATE expressions SET result = $1 WHERE id = $2`

	_, err := d.db.ExecContext(d.ctx, q1, "resolved", id)
	if err != nil {
		return err
	}

	_, err = d.db.ExecContext(d.ctx, q2, result, id)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) SelectPassword(login string) (string, error) {
	var q = `SELECT password FROM users WHERE login = $1`
	var pass = ""
	err := d.db.QueryRowContext(d.ctx, q, login).Scan(&pass)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return pass, nil
}

func (d *DB) SelectID(login string) (int64, error) {
	var q = `SELECT id FROM users WHERE login = $1`
	var id int64
	err := d.db.QueryRowContext(d.ctx, q, login).Scan(&id)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	return id, nil
}

func (d *DB) SelectExpressions(user_id int64) ([]Expression, error) {
	expArr := []Expression{}
	var q = `SELECT id, status, result FROM expressions WHERE user_id = $1`

	rows, err := d.db.QueryContext(d.ctx, q, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		exp := Expression{}
		err := rows.Scan(&exp.Id, &exp.Status, &exp.Result)
		if err != nil {
			return nil, err
		}
		expArr = append(expArr, exp)
	}

	return expArr, nil
}

func (d *DB) SelectExpressionById(expID int64, userID int64) (Expression, error) {
	var q1 = `SELECT user_id FROM expressions WHERE id = $1`
	var q2 = `SELECT id, status, result FROM expressions WHERE id = $1`

	var user_id int64
	exp := Expression{}

	err := d.db.QueryRowContext(d.ctx, q1, expID).Scan(&user_id)
	if err != nil {
		return exp, err
	}

	if userID != user_id {
		return exp, errors.New("no rights")
	}

	err = d.db.QueryRowContext(d.ctx, q2, expID).Scan(&exp.Id, &exp.Status, &exp.Result)
	if err != nil {
		return exp, err
	}

	return exp, nil
}

func (d *DB) LoginExists(login string) (bool, error) {
	var lg string
	var q = `SELECT id FROM users WHERE login = $1`
	err := d.db.QueryRowContext(d.ctx, q, login).Scan(&lg)
	if err != nil {
		if err != sql.ErrNoRows {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func (d *DB) Clear() (error) {
	q1 := `DELETE FROM users`
	q2 := `DELETE FROM expressions`
	_, err := d.db.ExecContext(d.ctx, q1)
	if err != nil {
		return err
	}
	_, err = d.db.ExecContext(d.ctx, q2)
	if err != nil {
		return err
	}
	return nil
}

func createTables(ctx context.Context, db *sql.DB) error {
	const (
		usersTable = `CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		login TEXT NOT NULL,
		password TEXT NOT NULL
		);`
		expressionsTable = `CREATE TABLE IF NOT EXISTS expressions(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		status TEXT NOT_NULL,
		result REAL,
		user_id INTEGER NOT NULL,
		
		FOREIGN KEY (user_id) REFERENCES expressions (id)
		);`
	)

	if _, err := db.ExecContext(ctx, usersTable); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, expressionsTable); err != nil {
		return err
	}
	return nil
}

func SqlUp() {
	ctx := context.TODO()
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		panic(err)
	}
	err = db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	err = createTables(ctx, db)
	if err != nil {
		panic(err)
	}

	MyDB.ctx = ctx
	MyDB.db = db
}

