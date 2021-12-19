package handler

import "fmt"

type CoffeeMaker interface {
	Make(intensity int) string
	PrintPrice()
}

type ExpresoCoffeeMaker struct {}
type LatteCoffeeMaker struct {
	version int
}

func (c *ExpresoCoffeeMaker) Make(intensity int) string {
	return "Making an awesome expreso coffee"
}

func (c *LatteCoffeeMaker) Make(intensity int) string {
	return "Making an awesome latte coffee"
}

func (c *ExpresoCoffeeMaker) PrintPrice(){
	fmt.Println("Expressos cost a shit ton of money")
}

func (c *LatteCoffeeMaker) PrintPrice(){
	fmt.Println("2 bucks")
}

func New() *LatteCoffeeMaker {
	return &LatteCoffeeMaker{}
}

func main() {
	latteMaker := New()
	fmt.Println(latteMaker.version)
	fmt.Println(latteMaker.Make(10))
	latteMaker.PrintPrice()
}