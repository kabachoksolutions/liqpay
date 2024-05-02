package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/kabachoksolutions/liqpay"
)

func main() {
	cfg := liqpay.NewConfig("sandbox_", "sandbox_", true)
	c := liqpay.NewClient(cfg, http.DefaultClient)

	orderID := uuid.New().String()

	r, err := c.CreateInvoice(&liqpay.InvoiceRequest{
		Amount:        100,
		Currency:      liqpay.CurrencyUAH,
		Description:   "Test",
		Email:         "test@gmail.com",
		OrderID:       orderID,
		Phone:         "380969696969",
		ActionPayment: "pay",
		ExpiredDate:   time.Now().Add(time.Hour * 5).Format(time.DateTime),
		Goods: []liqpay.InvoiceItem{
			{
				Amount: 100,
				Count:  2,
				Unit:   "pcs.",
				Name:   "Test",
			},
		},
		Language: liqpay.LanguageUK,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%#+v", r)

	resp, err := c.CancelInvoice(orderID)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%#+v", resp)

	slink, err := c.CreateSubscription(&liqpay.SubscriptionRequest{
		OrderID:            orderID,
		Amount:             100,
		Currency:           liqpay.CurrencyUAH,
		Description:        "test1",
		Phone:              "380969696969",
		SubscribeDateStart: time.Now().Format(time.DateTime),
		SubscribePeriod:    liqpay.SubscribePeriodMonthly,
		ServerURL:          "https://2844-193-56-13-203.ngrok-free.app/callback",
	})
	if err != nil {
		log.Println(err)
	}

	log.Printf("subscription link: %s", slink)

	clink, err := c.CreateCheckout(&liqpay.CheckoutRequest{
		OrderID:     orderID,
		Amount:      100,
		Currency:    liqpay.CurrencyUAH,
		Description: "test1",
		ServerURL:   "https://2844-193-56-13-203.ngrok-free.app/callback",
		ResultURL:   "https://2844-193-56-13-203.ngrok-free.app/test",
	})
	if err != nil {
		log.Println(err)
	}

	log.Printf("checkout link: %s", clink)

	status, err := c.Status(orderID)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("status: %#+v", status)

	app := fiber.New()

	app.Post("/callback", func(ctx *fiber.Ctx) error {
		var (
			data      = ctx.FormValue("data")
			signature = ctx.FormValue("signature")
		)

		if err := c.ValidateCallback(data, signature); err != nil {
			return err
		}

		var callback liqpay.Callback
		if err := ctx.BodyParser(&callback); err != nil {
			return err
		}

		log.Printf("callback: %#+v", callback)

		return ctx.SendStatus(fiber.StatusOK)
	})

	app.Get("/test", func(ctx *fiber.Ctx) error {
		return ctx.Status(200).SendString("hey")
	})

	log.Fatal(app.Listen(":3000"))
}
