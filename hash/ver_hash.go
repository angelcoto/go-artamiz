package hash

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

// VerificaHash recorre la tabla de hashes y verifica si el hash
// de cada archivo corresponde al definido en la tabla
func VerificaHash(archivohash string) {
	f, err := os.Open(archivohash)
	if err != nil {
		fmt.Println("* Error: ", err)
	}
	defer f.Close()

	var lineapartida []string
	var hash []byte
	var archivo, hashsrc string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lineapartida = strings.Split(scanner.Text(), "*")
		hashsrc = lineapartida[0][0 : len(lineapartida[0])-1]
		_, err := hex.DecodeString(hashsrc)

		if err == nil {
			archivo = lineapartida[1]
			longitud := len(lineapartida[0]) - 1
			switch longitud {
			case 64:
				hash, err = SumArchivo(archivo, "sha256")
			case 40:
				hash, err = SumArchivo(archivo, "sha1")
			case 32:
				hash, err = SumArchivo(archivo, "md5")
			}
			hashstr := hex.EncodeToString(hash)
			if err == nil {
				if hashstr == hashsrc {
					fmt.Printf("%s\tOK - Los valores coinciden\n", archivo)
				} else {
					fmt.Printf("%s\tFALLA - Los valores son distintos\n", archivo)

				}
			} else {
				fmt.Println("* Error: ", err)
			}
		}
	}
	/*
		if err := scanner.Err(); err != nil {
			fmt.Println("* Error: ", err)
		}
	*/
}
