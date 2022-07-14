/*
hash v base64 encoding:

 Base64 is decodable, SHA1 is not.

 They will both transform data into another format.

Encoding is reversible , hashing is not.
Endoding transforms data using a public algorithm so it can be easily reversed.

Hashing preserves the integrity of the data.


https://stackoverflow.com/questions/25058207/what-is-the-main-difference-between-base64-encode-hashing-and-sha1-md5
*/

//crypto hash

/*
SHA256 hashes are frequently used to compute short identities for binary or text blobs. For example, TLS/SSL certificates use SHA256 to compute a certificate’s signature. Here’s how to compute SHA256 hashes in Go.
*/
package main
//Go implements several hash functions in various crypto/* packages.
import (
    "crypto/sha256"
    "fmt"
)
func main() {
    s := "sha256 this string"
//Start with a new hash.
    h := sha256.New()
//Write expects bytes. If you have a string s, use []byte(s) to coerce it to bytes.
    h.Write([]byte(s))
//This gets the finalized hash result as a byte slice. The argument to Sum can be used to append to an existing byte slice: it usually isn’t needed.
    bs := h.Sum(nil)
    fmt.Println(s)
    fmt.Printf("%x\n", bs)
}
//Running the program computes the hash and prints it in a human-readable hex format.

//base64 encoding

//Go provides built-in support for base64 encoding/decoding.

package main
//This syntax imports the encoding/base64 package with the b64 name instead of the default base64. It’ll save us some space below.
import (
    b64 "encoding/base64"
    "fmt"
)
func main() {
//Here’s the string we’ll encode/decode.
    data := "abc123!?$*&()'-=@~"
//Go supports both standard and URL-compatible base64. Here’s how to encode using the standard encoder. The encoder requires a []byte so we convert our string to that type.
    sEnc := b64.StdEncoding.EncodeToString([]byte(data))
    fmt.Println(sEnc)
//Decoding may return an error, which you can check if you don’t already know the input to be well-formed.
    sDec, _ := b64.StdEncoding.DecodeString(sEnc)
    fmt.Println(string(sDec))
    fmt.Println()
//This encodes/decodes using a URL-compatible base64 format.
    uEnc := b64.URLEncoding.EncodeToString([]byte(data))
    fmt.Println(uEnc)
    uDec, _ := b64.URLEncoding.DecodeString(uEnc)
    fmt.Println(string(uDec))
}
//The string encodes to slightly different values with the standard and URL base64 encoders (trailing + vs -) but they both decode to the original string as desired.




//encryption
/*
Encryption is a way to hide data so that it is useless if it falls into the wrong hands. To encrypt in Go, we’ll use the Advanced Encryption Standard, provided by crypto/aes.

By adding random bytes, we can use them as an argument in the crypto/cipher module method, NewCFBEncrypter(). Then, before the Encode function, which encodes and returns the string to Base64, there is the MySecret constant that contains the secret for the encryption.

The Encrypt function, which takes two arguments, provides the text to encode and the secret to encode it. This then returns the Encode() function and passes the cipherText variable defined with the scope of Encrypt.

By running the file, the main function executes with the StringToEncrypt variable that contains the string to encrypt. The Encrypt() function also executes when the main function executes and now has two parameters: StringToEncrypt and MySecret.
*/

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

// This should be in an env file in production
const MySecret string = "abc&1*~#^2^#s0^=)^^7%b34"

func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
func Decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

// Encrypt method is to encrypt or hide any classified text
func Encrypt(text, MySecret string) (string, error) {
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return Encode(cipherText), nil
}

// Decrypt method is to extract back the encrypted text
func Decrypt(text, MySecret string) (string, error) {
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return "", err
	}
	cipherText := Decode(text)
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}
func main() {
	StringToEncrypt := "Encrypting this string"
	// To encrypt the StringToEncrypt
	encText, err := Encrypt(StringToEncrypt, MySecret)
	if err != nil {
		fmt.Println("error encrypting your classified text: ", err)
	}
	fmt.Println(encText)
	// To decrypt the original StringToEncrypt
	decText, err := Decrypt("Li5E8RFcV/EPZY/neyCXQYjrfa/atA==", MySecret)
	if err != nil {
		fmt.Println("error decrypting your encrypted text: ", err)
	}
	fmt.Println(decText)
}
