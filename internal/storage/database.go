package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

func NewStorage(dsn string) Storage {
	return Storage{
		DSN: dsn,
	}
}

func (s Storage) IncreaseBalance(u User) error {
	conn, err := pgx.Connect(context.Background(), s.DSN)
	if err != nil {
		log.Println("error conecting to database:", err)
		return err
	}
	defer conn.Close(context.Background())

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("err beginning transaction:", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		}
	}()

	sqlStr := `insert into users (id, balance)
		values ($2, $1)
		on conflict (id)
		do update set balance = users.balance + EXCLUDED.balance
		where EXCLUDED.id = $2;
	`

	_, err = tx.Exec(context.Background(), sqlStr, u.Balance, u.ID)
	if err != nil {
		log.Println("err increasing user balance", err)
		return err
	}

	sqlStr = `insert into history
		(user_id, operation_type, amount) 
		values ($1, $2, $3);
	`

	_, err = tx.Exec(context.Background(), sqlStr, u.ID, Arrival, u.Balance)
	if err != nil {
		log.Println("err inserting into history", err)
		return err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("err committing transaction:", err)
	}

	return nil
}

func (s Storage) TransferMoney(t Transaction) error {
	conn, err := pgx.Connect(context.Background(), s.DSN)
	if err != nil {
		log.Println("error conecting to database:", err)
		return err
	}
	defer conn.Close(context.Background())

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("err beginning transaction:", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		}
	}()

	sqlStr := `
		update users 
		set balance = balance - $1 
		where id = $2;
	`
	_, err = tx.Exec(context.Background(), sqlStr, t.Amount, t.SenderID)
	if err != nil {
		return fmt.Errorf("err updating sender balance:", err)
	}

	sqlStr = `
		update users 
		set balance = balance + $1 
		where id = $2;
	`
	_, err = tx.Exec(context.Background(), sqlStr, t.Amount, t.RecieverID)
	if err != nil {
		return fmt.Errorf("err updating receiver balance:", err)
	}

	sqlStr = `
	insert into transactions 
		(sender_id, money, reciever_id) 
	values ($1, $2, $3)
	`
	_, err = tx.Exec(context.Background(), sqlStr, t.SenderID, t.Amount, t.RecieverID)
	if err != nil {
		return fmt.Errorf("err inserting into transactions:", err)
	}

	sqlStr = `
	insert into history
		(user_id, operation_type, amount) 
	values ($1, $2, $3), ($4, $5, $6)
	`
	_, err = tx.Exec(context.Background(), sqlStr, t.SenderID, Transfer, t.Amount,
		t.RecieverID, Arrival, t.Amount)
	if err != nil {
		return fmt.Errorf("err inserting into history:", err)
	}

	if err := tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("err committing transaction:", err)
	}

	return nil

}

func (s Storage) GetHistory(userID int) ([]Operation, error) {
	conn, err := pgx.Connect(context.Background(), s.DSN)
	if err != nil {
		log.Println("error conecting to database:", err)
		return []Operation{}, err
	}
	defer conn.Close(context.Background())

	sqlStr := ` select id, user_id, operation_type, amount
				from history
				where user_id = $1
				order by time desc
				limit 10; 
	`

	rows, err := conn.Query(context.Background(), sqlStr, userID)
	if err != nil {
		log.Println("err getting user history:", err)
		return []Operation{}, err
	}

	fmt.Println(rows)

	var operations []Operation
	for rows.Next() {
		var o Operation
		err = rows.Scan(&o.ID, &o.UserID, &o.OperationType, &o.Amount)
		if err != nil {
			log.Println("err scanning a row:", err)
			return []Operation{}, err
		}

		operations = append(operations, o)
	}

	return operations, nil
}
