package cmn

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

// 使用公钥进行RSA加密后按Base64编码字符串
func EncodeRsa(str string, pubKey string) (string, error) {
	bt, err := EncodeRsaBytes(StringToBytes(str), StringToBytes(pubKey))
	if err != nil {
		return "", err
	}
	return Base64(bt), nil
}

// 按Base64解码字符串后使用私钥进行RSA解密
func DecodeRsa(str string, priKey string) (string, error) {
	by, err := Base64Decode(str)
	if err != nil {
		return "", err
	}
	bt, err := DecodeRsaBytes(by, StringToBytes(priKey))
	if err != nil {
		return "", err
	}
	return BytesToString(bt), nil
}

// 使用公钥文件进行RSA加密后按Base64编码字符串
func EncodeRsaByPubFile(str string, pubKeyFileName string) (string, error) {
	bt, err := EncodeRsaBytesByPubFile(StringToBytes(str), pubKeyFileName)
	if err != nil {
		return "", err
	}
	return Base64(bt), nil
}

// 按Base64解码字符串后使用私钥文件进行RSA解密
func DecodeRsaByPriFile(str string, pubKeyFileName string) (string, error) {
	by, err := Base64Decode(str)
	if err != nil {
		return "", err
	}
	bt, err := DecodeRsaBytesByPriFile(by, pubKeyFileName)
	if err != nil {
		return "", err
	}
	return BytesToString(bt), nil
}

// 当前目录下创建2048位的秘钥文件"rsa_private.pem、rsa_public.pem"
func GenerateRsaKey() error {
	return GenerateRsaKeyFile(2048, "rsa_private.pem", "rsa_public.pem")
}

// 使用公钥文件进行RSA加密
func EncodeRsaBytesByPubFile(data []byte, pubKeyFileName string) ([]byte, error) {
	file, err := os.Open(pubKeyFileName)
	if err != nil {
		return nil, err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buf := make([]byte, fileInfo.Size())
	file.Read(buf)

	return EncodeRsaBytes(data, buf)
}

// 使用私钥文件进行RSA解密
func DecodeRsaBytesByPriFile(data []byte, privateKeyFileName string) ([]byte, error) {
	file, err := os.Open(privateKeyFileName)
	if err != nil {
		return nil, err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	buf := make([]byte, fileInfo.Size())
	defer file.Close()
	file.Read(buf)

	return DecodeRsaBytes(data, buf)
}

// 创建秘钥文件（keySize通常是1024、2048、4096）
func GenerateRsaKeyFile(keySize int, priKeyFile, pubKeyFile string) error {
	pri, pub, err := GenerateRSAKey(keySize)
	if err != nil {
		return err
	}

	WriteFileBytes(priKeyFile, pri)
	WriteFileBytes(pubKeyFile, pub)
	return nil
}

// 创建秘钥
func GenerateRSAKey(keySize int) (privateKey []byte, publicKey []byte, err error) {
	prvKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return
	}

	pkixb, err := x509.MarshalPKIXPublicKey(&prvKey.PublicKey)
	if err != nil {
		return
	}

	privateKey = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(prvKey),
	})
	publicKey = pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pkixb,
	})
	return
}

// 使用公钥加密
func EncodeRsaBytes(data []byte, publicKey []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("无效的公钥")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	key, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("无效的公钥")
	}

	return rsa.EncryptPKCS1v15(rand.Reader, key, data)
}

// 使用秘钥解密
func DecodeRsaBytes(cipherText, privateKey []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("无效的秘钥")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return rsa.DecryptPKCS1v15(rand.Reader, key, cipherText)
}
