package accumulator

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
	"strconv"
)

const (
	securityPara      = 128
	RSABitLength      = 2048
	randomizerSetSize = 256
	// Note that the securityParaHashToPrime is running securityParaHashToPrime rounds of Miller-Robin test
	// together with one time Baillie-PSW test. Totally heuristic value for now.
	securityParaHashToPrime = 10

	// G, g, h are generated by TrustedSetupForQRN. They can only be used for testing purposes. DO NOT use in production.
	N2048String = "22582513446883649683242153375773765418277977026848618150278436227443969113525388360965414596382292671632010154272027792498289390464326093128963474525925743125404187090638221587455285089494562751793489098182761320953828657439130044252338283109583198301789045090284695934345711523245381620643226632165168827411546661236460973389982263385406789443858985073091473529732325356098830825299275985202060852102775942940039443155227986748457261585440368528834910182851433705587223040610934954417065434756145769875043620201897615075786323297141320586481340831246603933018654794846594742280842668198512719618188992528830140149361"
	G2048String = "3734320578166922768976307305081280303658237303482921793243310032002132951325426885895423150554487167609218974062079302792001919827304933109188668552532361245089029380294384169787606911401094856511916709999954764232948323779503820860893459514928713744983707360078264267038900798843893405664990521531326919997106338139056096176409033756102908667173913246197068450150318832809948977367751025873698025220766782003611956130604742644746610708520581969538416206455665972248047959779079118036299417601968576259426648158714614452861031491553305187113545916330322686053758561416773919173504690956803771722726889946697788319929"
	H2048String = "1582433196042535773898642856814926874501199844772808209798545765882857391073717631360065816613373509202691737458490830509979879771883168398785856056110736083435040549860024938378796318753064835110482441115760897524667343221753799849207723195729358565521753697076761550453675996906942484179834968386568757636433579938945322152073309477120701766107272148535093122238519340372766971216124175473667780382425281013570558875523373504108433319932127851859684947025440123382599601611460274335280822834972913253420025827402904805226163959418839188054187383250553791823431534564282919675786841775533806609995586228017407921459"
	// HashToPrimeFromSha256 is a prime number generated from Sha256
	HashToPrimeFromSha256 = iota
	// DIHashFromPoseidon is a division intractable Hash output
	DIHashFromPoseidon
)

var (
	big1  = big.NewInt(1)
	big2  = big.NewInt(2)
	big3  = big.NewInt(3)
	big5  = big.NewInt(5)
	big7  = big.NewInt(7)
	big11 = big.NewInt(11)
	big13 = big.NewInt(13)
	big17 = big.NewInt(17)
	big19 = big.NewInt(19)
	big23 = big.NewInt(23)
	big29 = big.NewInt(29)
	big31 = big.NewInt(31)
	big37 = big.NewInt(37)
	// Min2048 is set to a 2048 bits number with most significant bit 1 and other bits 0
	// This can speed up the calculation
	Min2048 = big.NewInt(0)
)

// Setup is a basic struct for a hidden order group
type Setup struct {
	N *big.Int
	G *big.Int //default generator in Z*_N
}

// Element should be able to be accumulated into RSA accumulator
type Element []byte

// EncodeType is the type of generating Element, should be consistent all the time
type EncodeType int

// GenerateG generates a generator for a hidden order group randomly
func GenerateG() {
	buffer := make([]big.Int, 8)
	buffer[0].Set(SHA256ToInt([]byte(N2048String))) //g1 should be 256 bit.
	for i := 1; i < 8; i++ {
		buffer[i].Set(SHA256ToInt(buffer[i-1].Bytes()))
	}
	prod := SetProduct(buffer)
	var N big.Int
	N.SetString(N2048String, 10)
	prod.Mod(prod, &N)
	fmt.Println("prod = ", prod.String())
	var gcd big.Int
	gcd.GCD(nil, nil, &N, prod)
	if gcd.Cmp(big1) != 0 {
		// gcd != 1
		//this condition should never happen
		fmt.Println("g and N not co-prime! We win the RSA-2048 challenge!")
	}
}

// SetProduct calculates the products of the input set
func SetProduct(inputSet []big.Int) *big.Int {
	var ret big.Int
	setSize := len(inputSet)
	ret.Set(big1)
	// ret is set to 1
	for i := 0; i < setSize; i++ {
		ret.Mul(&ret, &inputSet[i])
	}
	return &ret
}

// GetPseudoRandomElement returns the pseudo random element from the input integer, for test use only
func GetPseudoRandomElement(input int) *Element {
	var ret Element
	temp := strconv.Itoa(input)
	ret = []byte(temp[:])
	return &ret
}

// flipCoin outputs 1/0 with equal probability
func flipCoin() bool {
	//Int returns a uniform random value in [0, max)
	result, err := crand.Int(crand.Reader, big.NewInt(100))
	if err != nil {
		panic(err)
	}
	if result.Int64() < 50 {
		return true
	}
	return false
}

// TrustedSetupForQRN outputs a hidden order group
func TrustedSetupForQRN() {
	var p, q, N, g, h big.Int
	p = *getSafePrime()
	q = *getSafePrime()
	fmt.Println("Bit length of p = ", p.BitLen())
	fmt.Println("Bit length of q = ", q.BitLen())
	N.Mul(&p, &q)

	g = *getRanQR(&p, &q)
	// get a uniform random value randomNum in the QR_N, where the order of the group is p'q'
	randomNum, err := crand.Prime(crand.Reader, RSABitLength)
	if err != nil {
		panic(err)
	}
	order := getOrder(&p, &q)
	randomNum.Mod(randomNum, order)
	h.Exp(&g, randomNum, &N)
	fmt.Println("N = ", N.String())
	fmt.Println("g = ", g.String())
	fmt.Println("h = ", h.String())
}

func getOrder(p, q *big.Int) *big.Int {
	var pPrime, qPrime, phiN big.Int
	pPrime.Sub(p, big1)
	pPrime.Div(&pPrime, big2)
	qPrime.Sub(q, big1)
	qPrime.Div(&qPrime, big2)
	phiN.Mul(&pPrime, &qPrime)
	return &phiN
}

func testRemainder(input, modulo *big.Int) bool {
	var remainder, cmp big.Int
	cmp.Sub(modulo, big1)
	cmp.Div(&cmp, big2)
	remainder.Mod(input, modulo)
	if remainder.Cmp(&cmp) != 0 {
		return true
	}
	return true
}

// The following function implements the method described in "Safe Prime Generation with a Combined Sieve"
// for small prime r, if the input == (r-1)/2 mod r, then return false. How many r should be tested is purely experimental.
func safePrimeSieve(input *big.Int) bool {
	if !testRemainder(input, big3) {
		return false
	}
	if !testRemainder(input, big5) {
		return false
	}
	if !testRemainder(input, big7) {
		return false
	}
	if !testRemainder(input, big11) {
		return false
	}
	if !testRemainder(input, big13) {
		return false
	}
	if !testRemainder(input, big17) {
		return false
	}
	if !testRemainder(input, big19) {
		return false
	}
	if !testRemainder(input, big23) {
		return false
	}
	if !testRemainder(input, big29) {
		return false
	}
	if !testRemainder(input, big31) {
		return false
	}
	if !testRemainder(input, big37) {
		return false
	}

	return true
}

func getSuitablePrime() *big.Int {
	flag := false
	for !flag {
		ranNum, err := crand.Prime(crand.Reader, RSABitLength/2-1)
		if err != nil {
			panic(err)
		}
		flag = safePrimeSieve(ranNum)
		if !flag {
			continue
		}
		flag = ranNum.ProbablyPrime(securityPara / 2)
		if !flag {
			continue
		}
		return ranNum
	}
	return nil
}

// a safe prime p = 2p' +1 where p' is also a prime number
func getSafePrime() *big.Int {
	flag := false
	for !flag {
		ranNum := getSuitablePrime()
		//fmt.Println("get a prime = ", ranNum.String())
		ranNum.Mul(ranNum, big2)
		ranNum.Add(ranNum, big1)
		flag = ranNum.ProbablyPrime(securityPara / 2)
		if !flag {
			continue
		} else {
			fmt.Println("Found one safe prime = ", ranNum.String())
			return ranNum
		}
	}
	return nil
}

func getRanQR(p, q *big.Int) *big.Int {
	var N big.Int
	N.Mul(p, q)

	flag := false
	for !flag {
		ranNum, err := crand.Int(crand.Reader, Min2048)
		if err != nil {
			panic(err)
		}
		flag = isQR(ranNum, p, q)
		if !flag {
			continue
		} else {
			return ranNum
		}
	}
	return nil
}

func isQR(input, p, q *big.Int) bool {
	if big.Jacobi(input, p) == 1 && big.Jacobi(input, q) == 1 {
		return true
	}
	return false
}
