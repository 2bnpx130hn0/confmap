// Package encryptor provides AES-GCM encryption and decryption support
// for sensitive values stored in configuration maps.
//
// Encrypted values are stored as base64-encoded strings prefixed with "enc:",
// making them easy to identify and safe to store in config files.
//
// Example usage:
//
//	enc, err := encryptor.New([]byte("my-32-byte-secret-key-here!!!!!!"))
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	cipher, _ := enc.Encrypt("s3cr3t-password")
//	cfg := map[string]any{"db_password": cipher}
//
//	plain, _ := enc.DecryptConfig(cfg)
//	fmt.Println(plain["db_password"]) // s3cr3t-password
package encryptor
