package wfrsa

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"tunnel/internal/config"
	my_aes "tunnel/pkg/aes_gcm"
)

const (
	ENC_NONE = iota
	ENC_AES_GCM
	ENC_SM4
)

type RSAInfo struct {
	Prvkey        *rsa.PrivateKey
	Pubkey        rsa.PublicKey
	EncType       int
	ServerPublic  any
	ServerEncInfo struct {
		Key         []byte
		IV          []byte
		ADD         []byte
		AESCipher   cipher.Block
		AESGmc      cipher.AEAD
		MYAesCipher my_aes.AesCipher
	}
}

func NewRSA(c config.Config) (ras *RSAInfo, err error) {
	enc_type := c.GetString("enc_type", "")
	ras = &RSAInfo{}
	switch enc_type {
	case "aes":
		ras.EncType = ENC_AES_GCM
		break
	case "sm4":
		ras.EncType = ENC_SM4
		break
	case "none":
		ras.EncType = ENC_NONE
		break
	default:
		err = fmt.Errorf("enc_type is error")
		return
	}

	// Generates private key.
	ras.Prvkey, err = rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return
	}

	// Generates public key from private key.
	ras.Pubkey = ras.Prvkey.PublicKey
	ras.ServerEncInfo.Key = make([]byte, 0, 1024)
	ras.ServerEncInfo.IV = make([]byte, 0, 1024)
	ras.ServerEncInfo.ADD = make([]byte, 0, 1024)
	return
}

func (r RSAInfo) GetPublicKeyByte() (ret []byte, err error) {
	var derPkix []byte
	derPkix = x509.MarshalPKCS1PublicKey(&r.Prvkey.PublicKey)
	if err != nil {
		return nil, err
	}
	block := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: derPkix,
	}
	ret = pem.EncodeToMemory(block)
	return
}

func (r RSAInfo) RsaEncrypt(data []byte) (ret []byte, err error) {

	return
}

func (r RSAInfo) RsaDecrypt(cipherText []byte) (plainText []byte, err error) {
	return rsa.DecryptPKCS1v15(rand.Reader, r.Prvkey, cipherText)
}

func (r RSAInfo) SetSeverPublicKey(pk []byte) (err error) {
	block, _ := pem.Decode(pk)
	r.ServerPublic, err = x509.ParsePKIXPublicKey(block.Bytes)
	return
}

func (r *RSAInfo) GenAESCipher() (err error) {
	// cipher.newGCMWithNonceAndTagSize(cipher cipher.Block)
	r.ServerEncInfo.AESCipher, err = aes.NewCipher(r.ServerEncInfo.Key)

	if err != nil {
		return err
	}

	r.ServerEncInfo.AESGmc, err = cipher.NewGCM(r.ServerEncInfo.AESCipher)

	if err != nil {
		return err
	}
	// 原生的cipher中无法获取到生成签名的函数，导致无法校验，此处在初始化的时候，进行生成保存到自定义结构体中
	r.ServerEncInfo.MYAesCipher, err = my_aes.GenCipher(r.ServerEncInfo.Key, r.ServerEncInfo.AESGmc)
	return
}

func (r RSAInfo) Aes256GcmEnc(plaintext []byte) (ciphertext []byte) {
	ciphertext = r.ServerEncInfo.AESGmc.Seal(nil, r.ServerEncInfo.IV, plaintext, r.ServerEncInfo.ADD)

	return ciphertext[:len(plaintext)]
}

func (r RSAInfo) Aes256GcmDec(ciphertext []byte) (plaintext []byte, err error) {
	if len(ciphertext) <= r.ServerEncInfo.AESGmc.NonceSize() {
		err = fmt.Errorf("string: to short")
		return
	}

	ciphertext2 := make([]byte, 0, len(ciphertext)+r.ServerEncInfo.AESGmc.Overhead())

	ciphertext2 = append(ciphertext2, ciphertext...)
	for i := 0; i < r.ServerEncInfo.AESGmc.Overhead(); i++ {
		ciphertext2 = append(ciphertext2, 0)
	}

	// 使用原生库的Open 会有签名，而密文数据无法提供，所以通过函数进行更换
	// 正常使用  r.ServerEncInfo.AESGmc.Open(nil, r.ServerEncInfo.IV, ciphertext2, r.ServerEncInfo.ADD)
	plaintext, err = my_aes.Open(r.ServerEncInfo.AESGmc, r.ServerEncInfo.MYAesCipher, nil, r.ServerEncInfo.IV, ciphertext2, r.ServerEncInfo.ADD)
	return
}
