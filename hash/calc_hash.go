package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// SumTexto devuelve el hash para un texto
func SumTexto(s string, algo string) []byte {

	var sum []byte

	switch algo {
	case "md5":
		h := md5.New()
		h.Write([]byte(s))
		sum = h.Sum(nil)
	case "sha256":
		h := sha256.New()
		h.Write([]byte(s))
		sum = h.Sum(nil)
	default:
		h := sha1.New()
		h.Write([]byte(s))
		sum = h.Sum(nil)
	}

	return sum
}

// SumArchivo devuelve el hash de un archivo
// utilizando el algoritmo especificado
func SumArchivo(a string, algo string) ([]byte, error) {

	var sum []byte

	f, err := os.Open(a)
	if err != nil {
		return sum, err
	}
	defer f.Close()

	switch algo {
	case "sha256":
		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			return sum, err
		}
		sum = h.Sum(nil)
	case "sha1":
		h := sha1.New()
		if _, err := io.Copy(h, f); err != nil {
			return sum, err
		}
		sum = h.Sum(nil)
	case "md5":
		h := md5.New()
		if _, err := io.Copy(h, f); err != nil {
			return sum, err
		}
		sum = h.Sum(nil)
	}

	return sum, err

}

// SumDirectorio imprime el hash para los archivos
// de un directorio, sin incluir los subdirectorios
func SumDirectorio(dir string, algo string) {

	archivos, err := ioutil.ReadDir(dir)

	if err != nil {
		fmt.Printf("* Error: %s\n", err)
	}

	for _, f := range archivos {
		if !f.IsDir() {

			// Se usa la ruta completa para poder localizar el archivo
			// al momento de calcular el hash
			archivo := filepath.Join(dir, f.Name())
			hash, err := SumArchivo(archivo, algo)
			if err != nil {
				fmt.Printf("* Error: %s\n", err)
			} else {
				fmt.Printf("%x *%s\n", hash, f.Name())
			}

		}
	}
}

// SumRecursivo imprime el hash para los archivos de
// un directorio, incluyendo los subdirectorios
func SumRecursivo(dir string, algo string) {

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		//fmt.Println(path, info.IsDir(), algo)
		if !info.IsDir() {
			hash, err := SumArchivo(path, algo)
			if err != nil {
				fmt.Printf("* Error: %s\n", err)
			} else {
				fmt.Printf("%x *%s\n", hash, path)
			}
		}
		return nil

	})
}
