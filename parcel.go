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
	query := `INSERT INTO parcel (client_id, status, address, created_at) VALUES (?, ?, ?, ?)`
	result, err := s.db.Exec(query, p.Client, p.Status, p.Address, p.CreatedAt)
	if err != nil {
		return 0, fmt.Errorf("ошибка при добавлении посылки: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ошибка получения ID новой посылки: %w", err)
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
			return Parcel{}, fmt.Errorf("посылка с номером %d не найдена: %w", number, err)
		}
		return Parcel{}, fmt.Errorf("ошибка при получении посылки: %w", err)
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	query := `SELECT id, client_id, status, address, created_at FROM parcel WHERE client_id = ?`
	rows, err := s.db.Query(query, client)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении посылок клиента %d: %w", client, err)
	}
	defer rows.Close()

	var parcels []Parcel
	for rows.Next() {
		var p Parcel
		if err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("ошибка при чтении данных посылки: %w", err)
		}
		parcels = append(parcels, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка обработки результатов запроса: %w", err)
	}

	return parcels, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	query := `UPDATE parcel SET status = ? WHERE id = ?`
	result, err := s.db.Exec(query, status, number)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении статуса посылки %d: %w", number, err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("посылка с номером %d не найдена", number)
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	query := `UPDATE parcel SET address = ? WHERE id = ? AND status = ?`
	result, err := s.db.Exec(query, address, number, ParcelStatusRegistered)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении адреса посылки %d: %w", number, err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("нельзя изменить адрес: либо посылка не найдена, либо статус не '%s'", ParcelStatusRegistered)
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	query := `DELETE FROM parcel WHERE id = ? AND status = ?`
	result, err := s.db.Exec(query, number, ParcelStatusRegistered)
	if err != nil {
		return fmt.Errorf("ошибка при удалении посылки %d: %w", number, err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("нельзя удалить посылку %d: либо она не найдена, либо статус не '%s'", number, ParcelStatusRegistered)
	}

	return nil
}
