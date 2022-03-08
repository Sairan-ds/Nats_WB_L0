package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Database struct {
	Db *sql.DB
	Cache map[string]Order
}

// NewDatabase - returns a pointer to a new database connection
func NewDatabase() *Database {
	
	dbStruct := Database{}
	dbStruct.Cache = make(map[string]Order)
	log.Println("Setting up new database connection")

	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbTable := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", dbHost, dbPort, dbUsername, dbTable, dbPassword)

	dbObj, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println("err with open")
		panic(err)
	}


	err = dbObj.Ping()
	if err != nil {
		fmt.Println("err with ping")
		panic(err)
	}

	dbStruct.Db = dbObj
	return &dbStruct
}
// Add order to Database
func (db *Database) AddOrder(order *Order) {
	dbObj := db.Db


	dbObj.QueryRow(`insert into orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		order.OrderUid, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerId,
		order.DeliveryService, order.Shardkey, order.SmId, order.DateCreated, order.OofShard)
	dbObj.QueryRow(`insert into payments (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee, fk_payments_order) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		order.Payment.Transaction, order.Payment.RequestId, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount,
		order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal,
		order.Payment.CustomFee, order.OrderUid)
	dbObj.QueryRow(`insert into delivery (name, phone, zip, city, address, region, email, fk_delivery_order) values ($1,$2,$3,$4,$5,$6,$7,$8)`,
		order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address,
		order.Delivery.Region, order.Delivery.Email, order.OrderUid)
	for _, item := range order.Items {
		dbObj.QueryRow(`insert into items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status, fk_items_order) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
			item.ChrtId, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmId,
			item.Brand, item.Status, order.OrderUid)
	}
	db.Cache[order.OrderUid] = *order

}


