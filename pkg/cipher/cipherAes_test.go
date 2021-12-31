package cipher

import (
	"encoding/base64"
	"encoding/hex"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAesEncryptECB(t *testing.T) {
	assert := require.New(t)

	origData := []byte("123123")      // 待加密的数据
	key := []byte("ABCDEFGHIJKLMNOP") // 加密的密钥
	log.Println("原文：", string(origData))

	log.Println("------------------ ECB模式 --------------------")
	encrypted := AesEncryptECB(origData, key)
	log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	decrypted := AesDecryptECB(encrypted, key)
	log.Println("解密结果：", string(decrypted))

	assert.Equal(string(origData), string(decrypted))
}

func TestAesEncryptCBC(t *testing.T) {
	assert := require.New(t)

	origData := []byte("Hello World") // 待加密的数据
	key := []byte("ABCDEFGHIJKLMNOP") // 加密的密钥
	log.Println("原文：", string(origData))

	log.Println("------------------ CBC模式 --------------------")
	encrypted := AesEncryptCBC(origData, key)
	log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	decrypted := AesDecryptCBC(encrypted, key)
	log.Println("解密结果：", string(decrypted))
	assert.Equal(string(origData), string(decrypted))
}

func TestAesEncryptCFB(t *testing.T) {
	assert := require.New(t)

	origData := []byte("Hello World") // 待加密的数据
	key := []byte("ABCDEFGHIJKLMNOP") // 加密的密钥
	log.Println("原文：", string(origData))

	log.Println("------------------ CFB模式 --------------------")
	encrypted := AesEncryptCFB(origData, key)
	log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	decrypted := AesDecryptCFB(encrypted, key)
	log.Println("解密结果：", string(decrypted))

	assert.Equal(string(origData), string(decrypted))
}
