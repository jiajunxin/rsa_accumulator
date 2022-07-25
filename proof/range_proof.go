// Package proof range proof
// Variant of Lipmaa’s Compact Argument for Positivity proposed by Geoffroy Couteau et al. for range proof
// To prove an integer x lies in the range [a, b], we can show that x - a and b - x are positive by decomposing
// them as sum of four squares
// Paper: Removing the Strong RSA Assumption from Arguments over the Integers
// Link: https://eprint.iacr.org/2016/128
package proof

import (
	"bytes"
	"crypto"
	"crypto/sha256"
	"math/big"
)

const (
	rpChallengeStatement = "c = (g^x)(h^r), x is non-negative"
	sha256ResultLen      = 32
	commitLen            = sha256ResultLen * 5
)

var rpB = big.NewInt(4096) // bound B

// RangeProof is the proof for range proof
type RangeProof struct {
	// c = (g^x)(h^r)
	c *big.Int
	// commitment of x,
	// containing c1, c2, c3, c4, ci = (g^xi)(h^ri),
	// x = x1^2 + x2^2 + x3^2 + x4^2
	commitX FourNum
	// the commitment delta
	commitment rpCommitment
	// the response to the challenge
	response *rpResponse
}

// NewRangeProof generates a new proof for range proof
func NewRangeProof(c *big.Int, commitX FourNum, commitment rpCommitment, response *rpResponse) *RangeProof {
	return &RangeProof{
		c:          c,
		commitX:    commitX,
		commitment: commitment,
		response:   response,
	}
}

// rpCommitment is the range proof commitment generated by the prover
type rpCommitment [commitLen]byte

// rpChallenge is the challenge for range proof
type rpChallenge struct {
	statement string   // the statement for the challenge
	g, h, n   *big.Int // public parameters: G, H, N
	c4        FourNum  // commitment of x containing c1, c2, c3, c4
}

// newRPChallenge generates a new challenge for range proof
func newRPChallenge(pp *PublicParameters, c4 FourNum) *rpChallenge {
	return &rpChallenge{
		statement: rpChallengeStatement,
		g:         pp.G,
		h:         pp.H,
		n:         pp.N,
		c4:        c4,
	}
}

// Serialize generates the serialized data for range proof challenge in byte format
func (r *rpChallenge) serialize() []byte {
	var buf bytes.Buffer
	buf.WriteString(r.statement)
	buf.WriteString(r.g.String())
	buf.WriteString(r.h.String())
	buf.WriteString(r.n.String())
	for _, c := range r.c4 {
		buf.WriteString(c.String())
	}
	return buf.Bytes()
}

// sha256 generates the SHA256 hash of the range proof challenge
func (r *rpChallenge) sha256() []byte {
	hashF := crypto.SHA256.New()
	hashF.Write(r.serialize())
	hashResult := hashF.Sum(nil)
	return hashResult
}

// bigInt serializes the range proof challenge to bytes, generates the SHA256 hash of the byte data,
// and convert the hash to big integer
func (r *rpChallenge) bigInt() *big.Int {
	hashVal := r.sha256()
	return new(big.Int).SetBytes(hashVal)
}

// rpResponse is the response sent by the prover after receiving verifier's challenge
type rpResponse struct {
	Z4 FourNum
	T4 FourNum
	T  *big.Int
}

// newRPCommitment generates a new commitment for range proof
func newRPCommitment(d4 FourNum, d *big.Int) rpCommitment {
	var dByteList [4][]byte
	for i := 0; i < 4; i++ {
		dByteList[i] = d4[i].Bytes()
	}
	dBytes := d.Bytes()
	hashF := crypto.SHA256.New()
	var sha256List [4][]byte
	for i, dByte := range dByteList {
		hashF.Write(dByte)
		sha256List[i] = hashF.Sum(nil)
		hashF.Reset()
	}
	var commitment rpCommitment
	for idx, s := range sha256List {
		copy(commitment[idx*sha256ResultLen:(idx+1)*sha256ResultLen], s)
	}
	hashF.Write(dBytes)
	copy(commitment[commitLen-sha256ResultLen:], hashF.Sum(nil))
	return commitment
}

// RPProver refers to the Prover in zero-knowledge integer range proof
type RPProver struct {
	pp          *PublicParameters // public parameters
	x           *big.Int          // x, non-negative integer
	r           *big.Int          // r
	sp          *big.Int          // security parameter, kappa
	C           *big.Int          // c = (g^x)(h^r)
	fourSquareX FourNum           // Lagrange four square of x: x = x1^2 + x2^2 + x3^2 + x4^2
	commitFSX   FourNum           // commitment of four square of x: c1, c2, c3, c4, ci = (g^xi)(h^ri)
	randM4      FourNum           // random coins: m1, m2, m3, m4, mi is in [0, 2^(B/2 + 2kappa)]
	randR4      FourNum           // random coins: r1, r2, r3, r4, ri is in [0, n]
	randS4      FourNum           // random coins: s1, s2, s3, s4, si is in [0, 2^(2kappa)*n]
	randS       *big.Int          // random coin s in [0, 2^(B/2 + 2kappa)*n]
}

// NewRPProver generates a new range proof prover
func NewRPProver(pp *PublicParameters, r, x *big.Int) *RPProver {
	prover := &RPProver{
		pp: pp,
		x:  x,
		r:  r,
		sp: big.NewInt(securityParam),
	}
	prover.calC()
	return prover
}

// calculate parameter c, c = (g^x)(h^r)
func (r *RPProver) calC() *big.Int {
	r.C = new(big.Int).Exp(r.pp.G, r.x, r.pp.N)
	opt := iPool.Get().(*big.Int)
	defer iPool.Put(opt)
	r.C.Mul(r.C, opt.Exp(r.pp.H, r.r, r.pp.N))
	r.C.Mod(r.C, r.pp.N)
	return r.C
}

// Prove generates the proof for range proof
func (r *RPProver) Prove() (*RangeProof, error) {
	cx, err := r.commitForX()
	if err != nil {
		return nil, err
	}
	commitment, err := r.commit()
	if err != nil {
		return nil, err
	}
	response, err := r.response()
	if err != nil {
		return nil, err
	}
	return NewRangeProof(r.C, cx, commitment, response), nil
}

// commitForX generates the commitment for x
func (r *RPProver) commitForX() (FourNum, error) {
	// calculate lagrange four squares for x
	fs, err := LagrangeFourSquares(r.x)
	if err != nil {
		return FourNum{}, err
	}
	r.fourSquareX = fs
	// calculate commitment for x
	rc, err := newFourRandCoins(r.pp.N)
	if err != nil {
		return FourNum{}, err
	}
	r.randR4 = rc
	c4 := newRPCommitFromFS(r.pp, rc, fs)
	r.commitFSX = c4
	return c4, nil
}

// newRPCommitFromFS generates a range proof commitment for a given integer
func newRPCommitFromFS(pp *PublicParameters, coins FourNum, fs FourNum) (cList FourNum) {
	for i := 0; i < 4; i++ {
		cList[i] = new(big.Int).Exp(pp.G, fs[i], pp.N)
		cList[i].Mul(cList[i], new(big.Int).Exp(pp.H, coins[i], pp.N))
	}
	return
}

// commit composes the commitment for range proof
func (r *RPProver) commit() (rpCommitment, error) {
	// pick m1, m2, m3, m4, mi is in [0, 2^(B/2 + 2kappa)]
	mLmt := big.NewInt(2)
	//powMLmt := new(big.Int).Set(rpB)
	powMLmt := iPool.Get().(*big.Int).Set(rpB)
	defer iPool.Put(powMLmt)
	powMLmt.Rsh(powMLmt, 1)
	//powMLmtPart := new(big.Int).Set(r.sp)
	powMLmtPart := iPool.Get().(*big.Int).Set(r.sp)
	defer iPool.Put(powMLmtPart)
	powMLmtPart.Mul(powMLmtPart, big2)
	powMLmt.Add(powMLmt, powMLmtPart)
	mLmt.Exp(mLmt, powMLmt, nil)
	m4, err := newFourRandCoins(mLmt)
	if err != nil {
		return rpCommitment{}, err
	}
	r.randM4 = m4
	// pick s1, s2, s3, s4, si is in [0, 2^(B/2 + 2kappa)*n]
	sLmt := big.NewInt(2)
	//powSLmt := new(big.Int).Mul(r.sp, big2)
	powSLmt := iPool.Get().(*big.Int).Mul(r.sp, big2)
	defer iPool.Put(powSLmt)
	sLmt.Exp(sLmt, powSLmt, nil)
	sLmt.Mul(sLmt, r.pp.N)
	s4, err := newFourRandCoins(sLmt)
	if err != nil {
		return rpCommitment{}, err
	}
	r.randS4 = s4
	// pick s in [0, 2^(B/2 + 2kappa)*n]
	sLmt.Set(mLmt)
	sLmt.Mul(sLmt, r.pp.N)
	s, err := freshRandCoin(sLmt)
	if err != nil {
		return rpCommitment{}, err
	}
	r.randS = s
	// calculate commitment
	d4 := calD4(r.pp, m4, s4)
	d := calD(s, r.pp.H, r.pp.N, r.commitFSX, m4)
	c := newRPCommitment(d4, d)
	return c, nil
}

// calD4 calculates d1, d2, d3, d4, di = (g^mi)(h^si) mod n
func calD4(pp *PublicParameters, m, s FourNum) FourNum {
	var d4 FourNum
	for i := 0; i < 4; i++ {
		d4[i] = calDi(pp.G, pp.H, m[i], s[i], pp.N)
	}
	return d4
}

// calDi calculates di = (g^mi)(h^si) mod n
func calDi(g, h, mi, si, n *big.Int) *big.Int {
	res := new(big.Int).Set(g)
	res.Exp(res, mi, n)
	opt := iPool.Get().(*big.Int)
	defer iPool.Put(opt)
	res.Mul(res, opt.Exp(h, si, n))
	res.Mod(res, n)
	return res
}

// calD calculates d = product of (ci^mi)(h^s) mod n
func calD(s, h, n *big.Int, c FourNum, m FourNum) *big.Int {
	// h^s
	//hPowS := new(big.Int).Exp(h, s, n)
	hPowS := iPool.Get().(*big.Int).Exp(h, s, n)
	defer iPool.Put(hPowS)
	// ci^mi
	var cPowM4 FourNum
	for i := 0; i < 4; i++ {
		cPowM4[i] = new(big.Int).Exp(c[i], m[i], n)
	}
	// product of ci^mi
	d := big.NewInt(1)
	for i := 0; i < 4; i++ {
		d.Mul(d, cPowM4[i])
		d.Mod(d, n)
	}
	d.Mul(d, hPowS)
	d.Mod(d, n)
	return d
}

// calChallengeBigInt calculates the challenge for range proof in big integer format
func (r *RPProver) calChallengeBigInt() *big.Int {
	challenge := newRPChallenge(r.pp, r.commitFSX)
	return challenge.bigInt()
}

// response generates the response for verifier's challenge
func (r *RPProver) response() (*rpResponse, error) {
	c := r.calChallengeBigInt()
	var z4 FourNum
	for i := 0; i < 4; i++ {
		z4[i] = new(big.Int).Mul(c, r.fourSquareX[i])
		z4[i].Add(z4[i], r.randM4[i])
	}
	var t4 FourNum
	for i := 0; i < 4; i++ {
		t4[i] = new(big.Int).Mul(c, r.randR4[i])
		t4[i].Add(t4[i], r.randS4[i])
	}

	sumXR := big.NewInt(0)
	opt := iPool.Get().(*big.Int)
	defer iPool.Put(opt)
	for i := 0; i < 4; i++ {
		sumXR.Add(sumXR, opt.Mul(r.fourSquareX[i], r.randR4[i]))
	}
	t := new(big.Int).Sub(r.r, sumXR)
	t.Mul(t, c)
	t.Add(t, r.randS)
	response := &rpResponse{
		Z4: z4,
		T4: t4,
		T:  t,
	}
	return response, nil
}

// RPVerifier refers to the Verifier in zero-knowledge integer range proof
type RPVerifier struct {
	pp         *PublicParameters // public parameters
	sp         *big.Int          // security parameters
	C          *big.Int          // C, (g^x)(h^r)
	commitment rpCommitment      // commitment, delta = H(d1, d2, d3, d4, d)
	commitFSX  FourNum
}

// NewRPVerifier generates a new range proof verifier
func NewRPVerifier(pp *PublicParameters) *RPVerifier {
	verifier := &RPVerifier{
		pp: pp,
		sp: big.NewInt(securityParam),
	}
	return verifier
}

// Verify verifies the range proof
func (r *RPVerifier) Verify(proof *RangeProof) bool {
	r.SetC(proof.c)
	r.setCommitForX(proof.commitX)
	r.setCommitment(proof.commitment)
	return r.VerifyResponse(proof.response)
}

// SetC sets C to the verifier
func (r *RPVerifier) SetC(c *big.Int) {
	r.C = c
}

// setCommitment sets the commitment to the verifier
func (r *RPVerifier) setCommitment(c rpCommitment) {
	r.commitment = c
}

// setCommitForX sets the commitment of x to the verifier
// Commitment of x: c1, c2, c3, c4, ci = (g^x1=i)(h^ri)
func (r *RPVerifier) setCommitForX(c4 FourNum) {
	r.commitFSX = c4
}

// challenge generates a challenge for prover's commitment
func (r *RPVerifier) challenge() *big.Int {
	challenge := newRPChallenge(r.pp, r.commitFSX)
	return challenge.bigInt()
}

// VerifyResponse verifies the response, if accepts, return true; otherwise, return false
func (r *RPVerifier) VerifyResponse(response *rpResponse) bool {
	c := r.challenge()
	// the first 4 parameters: (g^zi)(h^ti)(ci^(-e)) mod n
	var firstFourParams FourNum
	//negC := new(big.Int).Neg(c)
	negC := iPool.Get().(*big.Int).Neg(c)
	defer iPool.Put(negC)
	opt := iPool.Get().(*big.Int)
	defer iPool.Put(opt)
	for i := 0; i < 4; i++ {
		firstFourParams[i] = new(big.Int).Exp(r.pp.G, response.Z4[i], r.pp.N)
		firstFourParams[i].Mul(
			firstFourParams[i],
			opt.Exp(r.pp.H, response.T4[i], r.pp.N),
		)
		firstFourParams[i].Mul(
			firstFourParams[i],
			opt.Exp(r.commitFSX[i], negC, r.pp.N),
		)
		firstFourParams[i].Mod(firstFourParams[i], r.pp.N)
	}

	//cPowNegE := new(big.Int).Exp(r.C, negC, r.pp.N)       // c^(-e)
	cPowNegE := iPool.Get().(*big.Int).Exp(r.C, negC, r.pp.N) // c^(-e)
	defer iPool.Put(cPowNegE)
	//hPowT := new(big.Int).Exp(r.pp.H, response.T, r.pp.N) // h^t
	hPowT := iPool.Get().(*big.Int).Exp(r.pp.H, response.T, r.pp.N) // h^t
	defer iPool.Put(hPowT)
	//product of (ci^zi)(h^t)(c^(-e)) mod n
	prodParam := big.NewInt(1)
	for i := 0; i < 4; i++ {
		prodParam.Mul(
			prodParam,
			opt.Exp(r.commitFSX[i], response.Z4[i], r.pp.N),
		)
		prodParam.Mod(prodParam, r.pp.N)
	}
	prodParam.Mul(prodParam, hPowT)
	prodParam.Mod(prodParam, r.pp.N)
	prodParam.Mul(prodParam, cPowNegE)
	prodParam.Mod(prodParam, r.pp.N)

	hashF := sha256.New()
	var sha256List [4][]byte
	for i := 0; i < 4; i++ {
		hashF.Write(firstFourParams[i].Bytes())
		sha256List[i] = hashF.Sum(nil)
		hashF.Reset()
	}
	hashF.Write(prodParam.Bytes())
	h := hashF.Sum(nil)
	var commitment rpCommitment
	for i := 0; i < 4; i++ {
		copy(commitment[i*sha256ResultLen:(i+1)*sha256ResultLen], sha256List[i])
	}
	copy(commitment[commitLen-sha256ResultLen:], h)
	return commitment == r.commitment
}
