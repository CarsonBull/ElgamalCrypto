package main

import (
	//"bytes"
	"fmt"
	//"io"
	"math/big"
	"math/rand"
	"os"
	"strconv"
)

func main() {

	// Use these three functions to run the program in the way that you would like
	// Buildkey takes a seed
	// EncryptFile takes a file name and a seed
	// decryptFile takes a file name for output

	buildKeyFiles(2124)

	encryptFile("message.txt", 1233)

	decryptFile("output.txt")

}

func getPrime(randSeed *rand.Rand) uint32 {

	var prime uint32

	for {
		prime = randSeed.Uint32()
		if prime&0x1 == 0 || (prime>>30)&0x1 != 0x1 || (prime>>31)&0x1 == 0x1 {
			continue
		}
		if millerRabin(prime, 4) {
			return prime
		}
	}

}

func millerRabin(prime uint32, numChecks int) bool {

	var temp int = 2 // This variable keeps a running power that is used to figire out if k is to big
	var k int
	var m int
	prime = prime - 1

	for i := 1; ; i++ {
		if int(prime)%temp != 0 {
			k = i - 1
			m = int(prime) / (temp / 2)
			break
		} else {
			temp = temp * 2
		}
	}

	for i := 0; i < numChecks; i++ {

		//var a int = 2
		var a int = rand.Intn(int(prime)-4) + 2

		b := int(new(big.Int).Exp(big.NewInt(int64(a)), big.NewInt(int64(m)), big.NewInt(int64(prime+1))).Int64())

		if b == 1 || b == -1 {
			//fmt.Println("First check")
			continue
		}

		var count int = 0

		//fmt.Println("Here", tempB, tempM)

		for {
			//fmt.Println(b, m, "--")
			//fmt.Println(power(b, m))
			//b = power(b, 2) % (int(prime) + 1)
			b = int(new(big.Int).Exp(big.NewInt(int64(b)), big.NewInt(int64(2)), big.NewInt(int64(prime+1))).Int64())

			//fmt.Println(b)
			if b == int(prime) {
				break
			} else if b == 1 {
				return false
			} else if count > k {
				//fmt.Println("counted up: ", count)
				return false
			}

			count += 1

			//fmt.Println(b, m, "-------")
		}

	}
	return true

}

func power(base, power, mod int) int {
	return int(new(big.Int).Exp(big.NewInt(int64(base)), big.NewInt(int64(power)), big.NewInt(int64(mod))).Int64())
}

// func powerMulti(base, power, multi, mod int) int {

// 	e2PowK := new(big.Int).Exp(big.NewInt(int64(base)), big.NewInt(int64(power)), nil)

// 	a := new(big.Int).Mul(e2PowK, big.NewInt(int64(multi)))

// 	return int(new(big.Int).Exp(a, big.NewInt(int64(1)), big.NewInt(int64(mod))).Int64())
// }

func powerMulti(base, power, multi, mod int) int {

	modB := big.NewInt(int64(mod))
	one := big.NewInt(int64(1))

	e2PowK := new(big.Int).Exp(big.NewInt(int64(base)), big.NewInt(int64(power)), modB)

	m := new(big.Int).Exp(big.NewInt(int64(multi)), one, modB)

	a := new(big.Int).Mul(e2PowK, m)

	return int(new(big.Int).Exp(a, one, modB).Int64())
}

func safePrime(key int) uint32 {

	// Seed 231 will result in the prime 3978804419 after a long waiting period

	randSource := rand.NewSource(int64(key))
	randSeed := rand.New(randSource)

	var prime uint32 = getPrime(randSeed)
	var safePrime uint32

	// first finds a prime to build primes from
	for {
		prime = getPrime(randSeed)
		// fmt.Printf("%32b , %d\n", safePrime, safePrime)
		if prime%12 == 5 {
			// If the prime matches the condition then we check if we can double it and add one to get another prime. Thats our safe prime
			safePrime = (prime * 2) + 1
			//fmt.Printf("%32b , %d\n", safePrime, safePrime)
			if millerRabin(safePrime, 4) {
				return safePrime
			}
		}
	}

}

func encrypt(p, g, e2, k, message int) (int, int) {

	fmt.Println(p, e2, k)

	var c1 int = power(g, k, p)
	var c2 int = powerMulti(e2, k, message, p)

	return c1, c2

}

func decrypt(p, d, c1, c2 int) int {

	var message int = powerMulti(c1, p-1-d, c2, p)

	return message

}

// Slower and didn't really work great.
// func expBtSquaring(x, n big.Int) big.Int {

// 	zero := big.NewInt(0)
// 	one := big.NewInt(1)
// 	two := big.NewInt(2)

// 	if n.Cmp(zero) == -1 {
// 		return expBtSquaring(*x.Div(one, &x), *n.Neg(&n))
// 	} else if n.Cmp(zero) == -1 {
// 		return *one
// 	} else if x.And(&x, one).Cmp(zero) == 0 {
// 		return expBtSquaring(*x.Mul(&x, &x), *n.Div(&n, two))
// 	} else {
// 		return expBtSquaring(*x.Mul(&x, &x), *n.Sub(&n, one).Div(&n, two))
// 	}

// 	// if n < 0 {
// 	// 	return expBtSquaring(1/x, -n)
// 	// } else if n == 0 {
// 	// 	return 1
// 	// } else if n == 1 {
// 	// 	return x
// 	// } else if n&0x1 == 0 {
// 	// 	return expBtSquaring(x*x, n/2)
// 	// } else if n&0x1 == 1 {
// 	// 	return x * expBtSquaring(x*x, (n-1)/2)
// 	// } else {
// 	// 	return x * expBtSquaring(x*x, (n-1)/2)
// 	// }

// }

func buildKeyFiles(key int) {

	fmt.Println("Generating P")
	var p int = int(safePrime(key)) // in reality you would want a random number there that is based off the enviroment or something.
	var g int = 2
	var userD int
	var e2 int

	// for {
	// 	//reader := bufio.NewReader(os.Stdin)
	// 	fmt.Printf("Enter d that's less than %d: ", p)

	// 	numScanned, _ := fmt.Scanf("%d", &userD)

	// 	if numScanned != 1 {
	// 		fmt.Println("Error reading int. Please try again.")
	// 	} else {
	// 		break
	// 	}
	// }

	userD = rand.Intn(p - 1)

	// I should have used python
	e2 = powerMulti(g, userD, 1, p)

	//e2 = expBtSquaring(g, userD) % p

	pubFile, err := os.Create("pubkey.txt")

	if err != nil {
		fmt.Println("Error opening pubkey.txt file")
		return
	}

	priFile, err := os.Create("prikey.txt")

	if err != nil {
		fmt.Println("Error opening prikey.txt file")
		return
	}

	test, err := pubFile.WriteString(strconv.Itoa(p) + " " + strconv.Itoa(g) + " " + strconv.Itoa(e2))

	if err != nil || test == 0 {
		fmt.Println("Error")
	}

	test, err = priFile.WriteString(strconv.Itoa(p) + " " + strconv.Itoa(g) + " " + strconv.Itoa(userD))

	if err != nil || test == 0 {
		fmt.Println("Error")
	}

	fmt.Println("KEYS P: ", p, " g: ", g, " e2: ", e2, " d: ", userD)

}

func encryptFile(fileName string, key int) {

	randSource := rand.NewSource(int64(key))
	randSeed := rand.New(randSource)

	messageFile, err := os.Open(fileName)

	if err != nil {
		fmt.Println("Error opening ", fileName, " file")
		return
	}

	fileBytes := readBytes(messageFile)

	messageFile.Close()

	keyFile, err := os.Open("pubkey.txt")

	if err != nil {
		fmt.Println("Error opening pubkey.txt file")
		return
	}

	var p int
	var g int
	var e2 int

	_, err = fmt.Fscanf(keyFile, "%d %d %d", &p, &g, &e2)

	if err != nil {
		fmt.Println(err)
		return
	}

	keyFile.Close()

	file, err := os.Create("ctext.txt")

	if err != nil {
		fmt.Println("Error opening ctext.txt file")
		return
	}

	var k int = randSeed.Intn(p - 1)
	//var k int = 5

	fmt.Println("K: ", k)

	for i := 0; i < len(fileBytes); i++ {
		c1, c2 := encrypt(p, g, e2, k, int(fileBytes[i]))
		_, err := file.WriteString(strconv.Itoa(c1) + " " + strconv.Itoa(c2) + " ")
		if err != nil {
			fmt.Println("Error with output file", err)
			return
		}
		k = randSeed.Intn(p - 1)
	}

	file.Close()

}

func readBytes(file *os.File) []byte {

	var returnByte []byte
	var offset int64 = 0

	for {
		temp := make([]byte, 1)
		read, err := file.ReadAt(temp, offset)

		if err != nil || read < len(temp) {
			break
		} else {
			returnByte = append(returnByte, temp[0])
			offset += 1
		}
	}

	return returnByte

}

func decryptFile(outputFileName string) {

	cypherFile, err := os.Open("ctext.txt")

	if err != nil {
		fmt.Println("Error opening the cipher file")
		return
	}

	keyFile, err := os.Open("prikey.txt")

	if err != nil {
		fmt.Println("Error opening prikey.txt file")
		return
	}

	var p int
	var g int
	var d int

	_, err = fmt.Fscanf(keyFile, "%d %d %d", &p, &g, &d)

	fmt.Println(p, d)

	if err != nil {
		fmt.Println(err)
		return
	}

	keyFile.Close()

	file, err := os.Create(outputFileName)

	if err != nil {
		fmt.Println("Error opening output file")
		return
	}

	var c1, c2 int
	read, err := fmt.Fscanf(cypherFile, "%d %d", &c1, &c2)

	if err != nil {
		fmt.Println(err)
		return
	}

	for read == 2 {
		if err != nil {
			fmt.Println(err)
			return
		}
		test := decrypt(p, d, c1, c2)
		fmt.Printf("%c\n", test)
		fmt.Fprintf(file, "%c", test)
		read, err = fmt.Fscanf(cypherFile, "%d %d", &c1, &c2)
	}

	fmt.Println("EXITING")

	cypherFile.Close()
	file.Close()

}
