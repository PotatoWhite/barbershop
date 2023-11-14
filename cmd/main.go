package main

import (
	"barber/pkg/barbershop"
	"fmt"
	"github.com/fatih/color"
	"math/rand"
	"time"
)

// variables
var (
	seatingCapacity = 10
	arrivalRate     = 100
	cutDuration     = 1000 * time.Millisecond
	timeOpen        = 10 * time.Second
)

func main() {
	// print welcome message
	color.Yellow("Welcome to the barber shop!")
	color.Yellow("---------------------------")

	// create channels if we need any
	clientChan := make(chan string, seatingCapacity)
	doneChan := make(chan bool)

	// create the barber shop
	shop := &barbershop.BarberShop{
		ShopCapacity:    seatingCapacity,
		HairCutDuration: cutDuration,
		NumberOfBarbers: 0,
		BarberDoneChan:  doneChan,
		ClientChan:      clientChan,
		Open:            true,
	}
	color.Green("Barber shop is open for business!")

	// add barbers
	shop.AddBarber("Frank")
	shop.AddBarber("Gerald")
	shop.AddBarber("Hank")
	shop.AddBarber("Ivan")
	shop.AddBarber("Karl")

	// start the barber shop as a goroutine
	shopClosing := make(chan bool)
	closed := make(chan bool)

	go func() {
		<-time.After(timeOpen) // at 10 seconds, close the shop
		shopClosing <- true
		shop.CloseShopForDay()
		closed <- true
	}()

	// add clients
	i := 1
	go func() {
		for {
			randomMillisecond := rand.Int() % (2 * arrivalRate) // random number between 0 and 2 * arrivalRate
			select {
			case <-shopClosing:
				return
			case <-time.After(time.Duration(randomMillisecond) * time.Millisecond):
				shop.AddClient(fmt.Sprintf("Client %d", i))
				i++
			}
		}
	}()

	<-closed
}
