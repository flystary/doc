package my_aes

import (
	"crypto/cipher"
	"fmt"
)

// A cipher is an instance of AES encryption using a particular key.
type AesCipher struct {
	enc []uint32
	dec []uint32

	// ks is the key schedule, the length of which depends on the size of
	// the AES key.
	ks []uint32
	// productTable contains pre-computed multiples of the binary-field
	// element used in GHASH.
	productTable [256]byte
	// nonceSize contains the expected size of the nonce, in bytes.
	nonceSize int
	// tagSize contains the size of the tag, in bytes.
	tagSize int
}

func GenCipher(key []byte, g cipher.AEAD) (ret AesCipher, err error) {
	n := len(key) + 28
	ret.dec = make([]uint32, n)
	ret.enc = make([]uint32, n)
	ret.nonceSize = g.NonceSize()
	ret.tagSize = g.Overhead()

	var rounds int
	switch len(key) {
	case 128 / 8:
		rounds = 10
	case 192 / 8:
		rounds = 12
	case 256 / 8:
		rounds = 14
	default:
		err = fmt.Errorf("key error")
		return
	}
	expandKeyAsm(rounds, &key[0], &ret.enc[0], &ret.dec[0])
	gcmAesInit(&ret.productTable, ret.enc)
	ret.ks = ret.enc
	return
}

// Open authenticates and decrypts ciphertext. See the [cipher.AEAD] interface
// for details.
func Open(g cipher.AEAD, gcmAsm AesCipher, dst, nonce, ciphertext, data []byte) ([]byte, error) {
	// var gcmAsm gcmAsm
	// g := wf.ServerEncInfo.AESGmc
	if len(nonce) != g.NonceSize() {
		panic("crypto/cipher: incorrect nonce length given to GCM")
	}
	// Sanity check to prevent the authentication from always succeeding if an implementation
	// leaves tagSize uninitialized, for example.
	if g.NonceSize() < gcmMinimumTagSize {
		panic("crypto/cipher: incorrect GCM tag size")
	}

	if len(ciphertext) < g.Overhead() {
		return nil, errOpen
	}
	if uint64(len(ciphertext)) > ((1<<32)-2)*uint64(BlockSize)+uint64(g.Overhead()) {
		return nil, errOpen
	}

	// tag := ciphertext[len(ciphertext)-g.Overhead():]
	ciphertext = ciphertext[:len(ciphertext)-g.Overhead()]

	// See GCM spec, section 7.1.
	var counter, tagMask [gcmBlockSize]byte

	if len(nonce) == gcmStandardNonceSize {
		// Init counter to nonce||1
		copy(counter[:], nonce)
		counter[gcmBlockSize-1] = 1
	} else {
		// Otherwise counter = GHASH(nonce)
		gcmAesData(&gcmAsm.productTable, nonce, &counter)
		gcmAesFinish(&gcmAsm.productTable, &tagMask, &counter, uint64(len(nonce)), uint64(0))
	}

	encryptBlockAsm(len(gcmAsm.ks)/4-1, &gcmAsm.ks[0], &tagMask[0], &counter[0])

	var expectedTag [gcmTagSize]byte
	gcmAesData(&gcmAsm.productTable, data, &expectedTag)

	ret, out := sliceForAppend(dst, len(ciphertext))
	if InexactOverlap(out, ciphertext) {
		panic("crypto/cipher: invalid buffer overlap")
	}
	if len(ciphertext) > 0 {
		gcmAesDec(&gcmAsm.productTable, out, ciphertext, &counter, &expectedTag, gcmAsm.ks)
	}
	gcmAesFinish(&gcmAsm.productTable, &tagMask, &expectedTag, uint64(len(ciphertext)), uint64(len(data)))

	// if subtle.ConstantTimeCompare(expectedTag[:g.gcmAsm.tagSize], tag) != 1 {
	// 	// for i := range out {
	// 	// 	out[i] = 0
	// 	// }
	// 	return nil, errOpen
	// }

	return ret, nil
}
