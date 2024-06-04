package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()

	res, err := s.db.Exec("INSERT INTO parcel (number, client, status, address, created_at) values (:number, :client, :status, :address, :created_at)",
		sql.Named("number", p.Number),
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	p := Parcel{}
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		return Parcel{}, err
	}
	defer db.Close()

	res := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = :number", sql.Named("number", number))
	err = res.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	var res []Parcel
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = :client",
		sql.Named("client", client))

	defer rows.Close()

	for rows.Next() {
		p := Parcel{}
		if err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt); err != nil {
		}

		res = append(res, p)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = s.db.Exec("UPDATE parcel SET address = :address WHERE status = :status",
		sql.Named("address", address),
		sql.Named("status", ParcelStatusRegistered))

	return nil
}

func (s ParcelStore) Delete(number int) error {
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = s.db.Exec("DELETE FROM parcel WHERE status = :status AND number = :number",
		sql.Named("status", ParcelStatusRegistered),
		sql.Named("number", number))

	return nil
}
