package barbershop

import (
	"github.com/fatih/color"
	"time"
)

type BarberShop struct {
	ShopCapacity    int
	HairCutDuration time.Duration
	NumberOfBarbers int

	BarberDoneChan chan bool
	ClientChan     chan string
	Open           bool
}

func (shop *BarberShop) AddBarber(barber string) {
	shop.NumberOfBarbers++

	go func() {
		isSleeping := false
		color.Yellow("%s goes to the waiting room to check for clients.", barber)

		for {
			// if there are no clients, the barber goes to sleep
			if len(shop.ClientChan) == 0 {
				color.Yellow("There is nothing to do, so %s takes a nap.", barber)
				isSleeping = true
			}

			client, shopOpen := <-shop.ClientChan

			if shopOpen {
				if isSleeping {
					color.Yellow("%s wakes %s up.", client, barber)
					isSleeping = false
				}
				// cut hair
				shop.cutHair(barber, client)
			} else {
				// shop is closed, so send the barber home and close this goroutine
				shop.sendBarberHome(barber)
				return
			}
		}
	}()
}

func (shop *BarberShop) cutHair(barber, client string) {
	color.Green("%s is cutting %s's hair.", barber, client)
	time.Sleep(shop.HairCutDuration)
	color.Green("%s is finished cutting %s's hair.", barber, client)
}

func (shop *BarberShop) sendBarberHome(barber string) {
	color.Cyan("%s is going home.", barber)
	shop.BarberDoneChan <- true
}

func (shop *BarberShop) CloseShopForDay() {
	color.Cyan("Closing shop for the day.")

	close(shop.ClientChan)
	shop.Open = false

	for a := 1; a <= shop.NumberOfBarbers; a++ {
		<-shop.BarberDoneChan
	}

	close(shop.BarberDoneChan)
	color.Green("----------------------------------------------------")
	color.Green("Barber shop is closed for the day. everyone is home.")
}

func (shop *BarberShop) AddClient(client string) {
	color.Green("%s client enters the shop.", client)

	if shop.Open {
		select {
		case shop.ClientChan <- client:
			color.Yellow("%s takes a seat in the waiting room.", client)
		default: // shop is full
			color.Red("The waiting room is full. %s leaves.", client)
		}
	} else {
		color.Red("Shop is already closed. %s leaves.", client)
	}
}
