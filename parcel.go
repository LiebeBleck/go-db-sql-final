package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {

	query := `INSERT INTO parcel (client_id, status, address, created_at) 
		VALUES (?, ?, ?, ?)`

	result, err := s.db.Exec(query, p.Client, p.Status, p.Address, p.CreatedAt)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	query := `SELECT id, client_id, status, address, created_at FROM parcel WHERE id = ?`
	row := s.db.QueryRow(query, number)

	var p Parcel
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return p, fmt.Errorf("посылка с номером %d не найдена", number)
		}
		return p, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	query := `SELECT id, client_id, status, address, created_at FROM parcel WHERE client_id = ?`
	rows, err := s.db.Query(query, client)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parcels []Parcel
	for rows.Next() {
		var p Parcel
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		parcels = append(parcels, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return parcels, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {

	query := `UPDATE parcel SET status = ? WHERE id = ?`
	_, err := s.db.Exec(query, status, number)
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {

	parcel, err := s.Get(number)
	if err != nil {
		return err
	}

	if parcel.Status != ParcelStatusRegistered {
		return fmt.Errorf("нельзя изменить адрес для посылки с таким статусом")
	}

	query := `UPDATE parcel SET address = ? WHERE id = ?`
	_, err = s.db.Exec(query, address, number)
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {

	parcel, err := s.Get(number)
	if err != nil {
		return err
	}

	if parcel.Status != ParcelStatusRegistered {
		return fmt.Errorf("нельзя удалить посылку с таким статусом")
	}

	query := `DELETE FROM parcel WHERE id = ?`
	_, err = s.db.Exec(query, number)
	if err != nil {
		return err
	}

	return nil
}
