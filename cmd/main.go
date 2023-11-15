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
	seatingCapacity = 5
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
	shop.AddBarber("Potato")

	// start the barber shop as a goroutine
	shopClosing := make(chan bool)
	closed := make(chan bool)

	// close the shop after 10 seconds
	go func() {
		<-time.After(timeOpen) // at 10 seconds, execute the following code
		shopClosing <- true
		shop.CloseShopForDay()
		closed <- true
	}()

	// add clients
	i := 1
	go func() {
		for {
			randomMillisecond := rand.Intn(100)
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
