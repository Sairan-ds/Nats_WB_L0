package database

import (
	"fmt"
	"log"
)

//

// Загрузка кэша из БД
func SetUp() *Database {
	dbStr := NewDatabase()
	orderRows, err := dbStr.Db.Query("SELECT * FROM orders")
	log.Println("Setting up 1")
	if err != nil {
		panic(err)
	}
	defer orderRows.Close()

	for orderRows.Next() {
		order := Order{}
		//order.OrderUid = "b563feb7b2b84b6test"
		orderRows.Scan(&order.OrderUid, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
			&order.CustomerId, &order.DeliveryService, &order.Shardkey, &order.SmId, &order.DateCreated, &order.OofShard)
		rowsPayment, err := dbStr.Db.Query("select transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee from payments where fk_payments_order=$1", &order.OrderUid)
		if err != nil {
			panic(err)
		}
		defer rowsPayment.Close()

		payment := Payment{}
		for rowsPayment.Next() {
			err = rowsPayment.Scan(&payment.Transaction, &payment.RequestId, &payment.Currency, &payment.Provider,
				&payment.Amount, &payment.PaymentDt, &payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee)
			if err != nil {
				panic(err)
			}
			order.Payment = payment
		}

		rowsDelivery, err := dbStr.Db.Query("select name, phone, zip, city, address, region, email from delivery where fk_delivery_order=$1", &order.OrderUid)
		if err != nil {
			panic(err)
		}
		defer rowsDelivery.Close()
		delivery := Delivery{}
		for rowsDelivery.Next() {
			err = rowsDelivery.Scan(&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address,
				&delivery.Region, &delivery.Email)
			if err != nil {
				panic(err)
			}
			order.Delivery = delivery
		}

		rowsItems, err := dbStr.Db.Query("select chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status from items where fk_items_order=$1", &order.OrderUid)
		if err != nil {
			panic(err)
		}
		defer rowsItems.Close()

		var items []Item
		for rowsItems.Next() {
			item := Item{}
			err := rowsItems.Scan(&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale,
				&item.Size, &item.TotalPrice, &item.NmId, &item.Brand, &item.Status)
			if err != nil {
				fmt.Println(err)
				continue
			}
			items = append(items, item)
			order.Items = items
			
		}
		dbStr.Cache[order.OrderUid] = order
	}

	return dbStr
}
