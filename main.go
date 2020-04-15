package main

import (
	"github.com/A1311981684/raspberryPi/DigitalLEDDisplay/LED3461BS"
	"log"
	"os"
	"strconv"
)

func main() {
	var content, duration string
	if len(os.Args) != 3 {
		content, duration = "13.14", "3"
		log.Println("actually, you can assign content and" +
			" display duration by giving me these 2 cmd args: ./DigitalLEDDisplay 13.14 3")
	}else{
		content, duration = os.Args[1], os.Args[2]
	}
	led := LED3461BS.Led3461BS{}
	err := led.Init()
	if err != nil {
		panic(err)
	}

	du, err := strconv.Atoi(duration)
	if err != nil {
		panic(err)
	}
	log.Println(content, du)

	//uncomment this to try flexible control!
	//led.FlexibleControl(10,6,6,6,[4]bool{true, false, false, false}, 10)


	err = led.Execute(content, du)
	if err != nil {
		panic(err)
	}
}
