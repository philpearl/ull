package ull

import "math"

// Constants for OptimalFGRAEstimator
const (
	eta0 = 4.663135422063788
	eta1 = 2.1378502137958524
	eta2 = 2.781144650979996
	eta3 = 0.9824082545153715
	tau  = 0.8194911375910897

	pow2Tau      = 1.7631258657688563 // 2.0^tau
	pow2MinusTau = 0.5670918786435586 // 2.0^(-tau)
	pow4MinusTau = 0.3216030842364037 // 4.0^(-tau)

	minusInvTau = -1.0 / tau
	etaX        = eta0 - eta1 - eta2 + eta3
	eta23X      = (eta2 - eta3) / etaX
	eta13X      = (eta1 - eta3) / etaX
	eta3012XX   = (eta3*eta0 - eta1*eta2) / (etaX * etaX)

	pow4MinusTauEta23 = pow4MinusTau * (eta2 - eta3)
	pow4MinusTauEta01 = pow4MinusTau * (eta0 - eta1)
	pow4MinusTauEta3  = pow4MinusTau * eta3
	pow4MinusTauEta1  = pow4MinusTau * eta1
	pow2MinusTauEtaX  = pow2MinusTau * etaX

	phi1     = eta0 / (pow2Tau * (2.0*pow2Tau - 1.0))
	pInitial = etaX * (pow4MinusTau / (2.0 - pow2MinusTau))
)

// Pre-computed estimation factors for precision 4 to 18
var estimationFactors = [15]float64{
	455.6358404615186,
	2159.476860400962,
	10149.51036338182,
	47499.52712820488,
	221818.76564766388,
	1034754.6840013304,
	4824374.384717942,
	2.2486750611989766e7,
	1.0479810199493326e8,
	4.8837185623048025e8,
	2.275794725435168e9,
	1.0604938814719946e10,
	4.9417362104242645e10,
	2.30276227770117e11,
	1.0730444972228585e12,
}

// Pre-computed register contributions
var registerContributions = [236]float64{
	0.8484061093359406, 0.38895829052007685, 0.5059986252327467, 0.17873835725405993,
	0.48074234060273024, 0.22040001471443574, 0.2867199572932749, 0.10128061935935387,
	0.2724086914332655, 0.12488785473931466, 0.16246750447680292, 0.057389829555353204,
	0.15435814343988866, 0.0707666752272979, 0.09206087452057209, 0.03251947467566813,
	0.08746577181824695, 0.0400993542020493, 0.05216553700867983, 0.018426892732996067,
	0.04956175987398336, 0.022721969094305374, 0.029559172293066274, 0.01044144713836362,
	0.02808376340530896, 0.012875216815740723, 0.01674946174724118, 0.005916560101748389,
	0.015913433441643893, 0.0072956356627506685, 0.009490944673308844, 0.0033525700962450116,
	0.009017216113341773, 0.004134011914931561, 0.0053779657012946284, 0.0018997062578498703,
	0.005109531310944485, 0.002342503834183061, 0.00304738001114257, 0.001076452918957914,
	0.0028952738727082267, 0.0013273605219527246, 0.0017267728074345586, 6.09963188753462e-4,
	0.0016405831157217021, 7.521379173550258e-4, 9.78461602292084e-4, 3.4563062172237723e-4,
	9.2962292270938e-4, 4.2619276177576713e-4, 5.544372155028133e-4, 1.958487477192352e-4,
	5.267631795945699e-4, 2.4149862146135835e-4, 3.141672858847145e-4, 1.1097608132071735e-4,
	2.9848602115777116e-4, 1.3684320663902123e-4, 1.7802030736817869e-4, 6.288368329501905e-5,
	1.6913464774658265e-4, 7.754107700464113e-5, 1.0087374230011362e-4, 3.563252169014952e-5,
	9.583875639268212e-5, 4.393801322487549e-5, 5.715927601779108e-5, 2.0190875207520577e-5,
	5.430624268457414e-5, 2.4897113642537945e-5, 3.2388833410757184e-5, 1.144099329232623e-5,
	3.0772185549154786e-5, 1.4107744575453657e-5, 1.8352865935237916e-5, 6.482944704957522e-6,
	1.7436805727319977e-5, 7.99403737572986e-6, 1.0399500462555932e-5, 3.67350727106242e-6,
	9.880422483694849e-6, 4.529755498675165e-6, 5.892791363067244e-6, 2.081562667074589e-6,
	5.5986600976661345e-6, 2.5667486794686803e-6, 3.339101736056405e-6, 1.1795003568090263e-6,
	3.1724346748254955e-6, 1.4544270182973653e-6, 1.8920745223756656e-6, 6.683541714686068e-7,
	1.7976340035771381e-6, 8.241391019206623e-7, 1.072128458850476e-6, 3.7871739159788393e-7,
	1.0186145159929963e-6, 4.6699164053601817e-7, 6.075127690181302e-7, 2.1459709360913574e-7,
	5.77189533646426e-7, 2.6461697039041317e-7, 3.442421115430427e-7, 1.2159967724530947e-7,
	3.27059699739513e-7, 1.4994302882644454e-7, 1.9506195985170504e-7, 6.890345650764188e-8,
	1.853256875916027e-7, 8.49639834530526e-8, 1.1053025444979778e-7, 3.904357664636507e-8,
	1.0501327589016596e-7, 4.814414208323267e-8, 6.263105916717392e-8, 2.2123721430020238e-8,
	5.9504908663745294e-8, 2.7280481949286693e-8, 3.548937430686624e-8, 1.2536224699555158e-8,
	3.371796684815404e-8, 1.545826061452554e-8, 2.0109761920695445e-8, 7.103548569567803e-9,
	1.910600846054063e-8, 8.759296176321385e-9, 1.139503111580109e-8, 4.0251673442004705e-9,
	1.082626247715867e-8, 4.963383100969499e-9, 6.456900615837058e-9, 2.28082795382416e-9,
	6.134612546958812e-9, 2.812460192131048e-9, 3.65874960227048e-9, 1.292412391857717e-9,
	3.476127720042246e-9, 1.5936574250689536e-9, 2.0732003554895977e-9, 7.323348470132607e-10,
	1.9697191686598677e-9, 9.030328662369446e-10, 1.1747619217600795e-9, 4.1497151491950363e-10,
	1.1161251587553774e-9, 5.116961428952198e-10, 6.656691762391315e-10, 2.351401942661752e-10,
	6.324431369849931e-10, 2.899484087937328e-10, 3.771959611450379e-10, 1.3324025619025952e-10,
	3.5836869940773545e-10, 1.6429687995368037e-10, 2.1373498756237659e-10, 7.549949478033437e-11,
	2.0306667462222755e-10, 9.309747508122088e-11, 1.2111117194789844e-10, 4.2781167456975155e-11,
	1.1506606020637118e-10, 5.275291818652914e-11, 6.86266490006118e-11, 2.424159650745726e-11,
	6.520123617549523e-11, 2.9892007004129765e-11, 3.888672595026375e-11, 1.3736301184893309e-11,
	3.6945743959497274e-11, 1.693805979747882e-11, 2.2034843273746723e-11, 7.783562034953282e-12,
	2.093500180037604e-11, 9.597812206565218e-12, 1.248586262365167e-11, 4.4104913787558985e-12,
	1.186264650299681e-11, 5.4385213096368525e-12, 7.075011313669894e-12, 2.499168647301308e-12,
	6.721871027139603e-12, 3.081693348317683e-12, 4.008996942969544e-12, 1.4161333491633975e-12,
	3.808892905481426e-12, 1.7462161775917615e-12, 2.271665129027518e-12, 8.024403094117999e-13,
	2.1582778227746425e-12, 9.89479027998621e-13, 1.2872203525845489e-12, 4.54696198313039e-13,
	1.2229703685228866e-12, 5.606801491206791e-13, 7.293928206826874e-13, 2.5764985922987735e-13,
	6.92986095905959e-13, 3.1770479284824887e-13, 4.1330443990824427e-13, 1.4599517261737423e-13,
	3.926748688923721e-13, 1.8002480658009348e-13, 2.3419555992885186e-13, 8.272696321778206e-14,
	2.225059832666067e-13, 1.0200957528418621e-13, 1.327049869160979e-13, 4.687655297461429e-14,
	1.2608118449008524e-13, 5.780288643182276e-14, 7.519618885068399e-14, 2.656221301145837e-14,
	7.144286571105751e-14, 3.2753529955811655e-14, 4.2609301647742677e-14, 1.5051259431302017e-14,
	4.0482511975524363e-14, 1.8559518231526075e-14, 2.4144210160882415e-14, 8.528672304925501e-15,
	2.293908229376684e-14, 1.0516598285774437e-14, 1.3681118012966618e-14, 4.832701981970378e-15,
	1.2998242223663023e-14, 5.959143881034847e-15, 7.752292944665042e-15, 2.7384108113817744e-15,
	7.365346997814574e-15, 3.376699844369893e-15, 4.392773006047039e-15, 1.5516979527951759e-15,
	4.173513269314059e-15, 1.9133791810691354e-15, 2.4891286772044455e-15, 8.792568765435867e-16,
}

func (u *UltraLogLog) FGRAEstimate() uint64 {
	m := len(u.registers)
	p := u.precision
	off := int(p<<2) + 4

	var sum float64
	var c0, c4, c8, c10 int
	var c4w0, c4w1, c4w2, c4w3 int

	// Process each register
	for _, reg := range u.registers {
		r := int(reg)
		r2 := r - off
		if r2 < 0 {
			if r2 < -8 {
				c0++
			}
			if r2 == -8 {
				c4++
			}
			if r2 == -4 {
				c8++
			}
			if r2 == -2 {
				c10++
			}
		} else if r < 252 {
			sum += registerContributions[r2]
		} else {
			switch r {
			case 252:
				c4w0++
			case 253:
				c4w1++
			case 254:
				c4w2++
			case 255:
				c4w3++
			}
		}
	}

	// Handle small range estimates if needed
	if c0 > 0 || c4 > 0 || c8 > 0 || c10 > 0 {
		z := smallRangeEstimate(c0, c4, c8, c10, m)
		if c0 > 0 {
			sum += calculateContribution0(c0, z)
		}
		if c4 > 0 {
			sum += calculateContribution4(c4, z)
		}
		if c8 > 0 {
			sum += calculateContribution8(c8, z)
		}
		if c10 > 0 {
			sum += calculateContribution10(c10, z)
		}
	}

	// Handle large range estimates if needed
	if c4w0 > 0 || c4w1 > 0 || c4w2 > 0 || c4w3 > 0 {
		sum += calculateLargeRangeContribution(c4w0, c4w1, c4w2, c4w3, m, 65-int(p))
	}

	// Return final estimate
	return uint64(estimationFactors[p-4] * math.Pow(sum, minusInvTau))
}

// smallRangeEstimate computes z for small range registers
func smallRangeEstimate(c0, c4, c8, c10, m int) float64 {
	alpha := m + 3*(c0+c4+c8+c10)
	beta := m - c0 - c4
	gamma := 4*c0 + 2*c4 + 3*c8 + c10

	alphaF := float64(alpha)
	betaF := float64(beta)
	gammaF := float64(gamma)

	quadRootZ := (math.Sqrt(betaF*betaF+4.0*alphaF*gammaF) - betaF) / (2.0 * alphaF)
	rootZ := quadRootZ * quadRootZ
	return rootZ * rootZ
}

// calculateContribution0 computes contribution for c0 registers
func calculateContribution0(c0 int, z float64) float64 {
	return float64(c0) * sigma(z)
}

// calculateContribution4 computes contribution for c4 registers
func calculateContribution4(c4 int, z float64) float64 {
	return float64(c4) * pow2MinusTauEtaX * psiPrime(z, z*z)
}

// calculateContribution8 computes contribution for c8 registers
func calculateContribution8(c8 int, z float64) float64 {
	return float64(c8) * (z*pow4MinusTauEta01 + pow4MinusTauEta1)
}

// calculateContribution10 computes contribution for c10 registers
func calculateContribution10(c10 int, z float64) float64 {
	return float64(c10) * (z*pow4MinusTauEta23 + pow4MinusTauEta3)
}

// psiPrime computes (z + eta23X) * (zSquare + eta13X) + eta3012XX
func psiPrime(z, zSquare float64) float64 {
	return (z+eta23X)*(zSquare+eta13X) + eta3012XX
}

// sigma computes the sigma function for contribution0
func sigma(z float64) float64 {
	if z <= 0.0 {
		return eta3
	}
	if z >= 1.0 {
		return math.Inf(1)
	}

	powZ := z
	nextPowZ := powZ * powZ
	s := 0.0
	powTau := etaX

	for {
		oldS := s
		nextNextPowZ := nextPowZ * nextPowZ
		s += powTau * (powZ - nextPowZ) * psiPrime(nextPowZ, nextNextPowZ)
		if !(s > oldS) {
			return s / z
		}
		powZ = nextPowZ
		nextPowZ = nextNextPowZ
		powTau *= pow2Tau
	}
}

// calculateLargeRangeContribution computes contribution for large range registers
func calculateLargeRangeContribution(c4w0, c4w1, c4w2, c4w3, m, w int) float64 {
	z := largeRangeEstimate(c4w0, c4w1, c4w2, c4w3, m)
	rootZ := math.Sqrt(z)
	s := phi(rootZ, z) * float64(c4w0+c4w1+c4w2+c4w3)
	s += z * (1.0 + rootZ) * (float64(c4w0)*eta0 + float64(c4w1)*eta1 + float64(c4w2)*eta2 + float64(c4w3)*eta3)
	s += rootZ * (float64(c4w0+c4w1)*(z*pow2MinusTau*(eta0-eta2)+pow2MinusTau*eta2) +
		float64(c4w2+c4w3)*(z*pow2MinusTau*(eta1-eta3)+pow2MinusTau*eta3))
	return s * math.Pow(pow2MinusTau, float64(w)) / ((1.0 + rootZ) * (1.0 + z))
}

// largeRangeEstimate computes z for large range registers
func largeRangeEstimate(c4w0, c4w1, c4w2, c4w3, m int) float64 {
	alpha := m + 3*(c4w0+c4w1+c4w2+c4w3)
	beta := c4w0 + c4w1 + 2*(c4w2+c4w3)
	gamma := m + 2*c4w0 + c4w2 - c4w3

	alphaF := float64(alpha)
	betaF := float64(beta)
	gammaF := float64(gamma)

	inner := (math.Sqrt(betaF*betaF+4.0*alphaF*gammaF) - betaF) / (2.0 * alphaF)
	return math.Sqrt(inner)
}

// phi computes the phi function for large range contribution
func phi(z, zSquare float64) float64 {
	if z <= 0.0 {
		return 0.0
	}
	if z >= 1.0 {
		return phi1
	}

	previousPowZ := zSquare
	powZ := z
	nextPowZ := math.Sqrt(z)
	p := pInitial / (1.0 + nextPowZ)
	ps := psiPrime(powZ, previousPowZ)
	s := nextPowZ * (ps + ps) * p

	for {
		previousPowZ = powZ
		powZ = nextPowZ
		oldS := s
		nextPowZ = math.Sqrt(powZ)
		nextPs := psiPrime(powZ, previousPowZ)
		p *= pow2MinusTau / (1.0 + nextPowZ)
		s += nextPowZ * ((nextPs + nextPs) - (powZ+nextPowZ)*ps) * p
		if !(s > oldS) {
			return s
		}
		ps = nextPs
	}
}
