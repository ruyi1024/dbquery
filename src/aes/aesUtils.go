/*
Copyright 2014-2024 The Lepus Team Group, website: https://www.lepus.cc
Licensed under the GNU General Public License, Version 3.0 (the "GPLv3 License");
You may not use this file except in compliance with the License.
You may obtain a copy of the License at
    https://www.gnu.org/licenses/gpl-3.0.html
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
Special note:
Please do not use this source code for any commercial purpose,
or use it for commercial purposes after secondary development, otherwise you may bear legal risks.
*/

package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
)

// AesEncrypt AES加密
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// AesDecrypt 解密
func AesDecrypt(decryptStr, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("key 长度必须是16/24/32: %s", err)
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	origData := make([]byte, len(decryptStr))
	blockMode.CryptBlocks(origData, decryptStr)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func AesPassEncode(pwdStr, aesKey string) (string, error) {
	pwd := []byte(pwdStr)
	result, err := AesEncrypt(pwd, []byte(aesKey))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(result), nil
}

func AesPassDecode(pwd, aesKey string) (string, error) {
	temp, _ := hex.DecodeString(pwd)
	//执行AES解密
	res, err := AesDecrypt(temp, []byte(aesKey))
	if err != nil {
		return "", err
	}
	return string(res), nil
}
