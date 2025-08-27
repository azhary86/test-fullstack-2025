// f(n) = ((n!) / 2^n)

package main

import 	(
	"fmt"
	"math"
	"log"
)

func faktorial (n int) int {
	if n <= 1{
		return 1
	}

	return n * faktorial(n-1)
}

func pangkat (n int) int {
	return 1 << uint(n)
}

func kalkulasi (n int) int{
	if n < 0{
		return 0
	}

	fakt := faktorial(n)
	pow := pangkat(n)

	result := float64(fakt) / float64(pow)
	return int(math.Ceil(result))
}

func main(){
	fmt.Println("Enter number: ")
	var number int
	_, err := fmt.Scanf("%d",&number)
	if err != nil {	
		log.Fatalf("Error reading integer: %v", err)
	}
	
	fmt.Printf("f(%d)= %d",number, kalkulasi(number))
}