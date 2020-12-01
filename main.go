package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"time"

	"github.com/angelcoto/go-artamiz/hash"
)

const linea = "----------------------------------------------------------------"

func encabezado(t time.Time) {
	usuario, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Tabla hash generada por:", usuario.Username)
	fmt.Println("Inicio:", t)
	fmt.Println(linea)
}

func pie(t, inicio time.Time) {
	fmt.Println(linea)
	fmt.Println("Fin:", t)
	fmt.Println("Duración: ", time.Since(inicio))
}

func main() {

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	inicio := time.Now()

	filePtr := flag.String("a", "", "Archivo")
	txtPtr := flag.String("t", "", "Texto")
	dirPtr := flag.String("d", dir, "Directorio a recorrer")
	algoPtr := flag.String("m", "sha256", "Algoritmo: md5, sha1, sha256")
	recPtr := flag.Bool("r", false, "Recorrido recursivo")

	flag.Parse()

	archivo := *filePtr
	texto := *txtPtr
	directorio := *dirPtr
	algo := *algoPtr
	recursivo := *recPtr

	if archivo != "" {
		hash, err := hash.SumArchivo(archivo, algo)

		if err != nil {
			fmt.Printf("* Error: %s\n", err)
		} else {
			fmt.Printf("%x *%s\n", hash, archivo)
		}

	} else if texto != "" {
		fmt.Printf("%x\n", hash.SumTexto(texto, algo))

	} else {

		encabezado(time.Now())
		defer pie(time.Now(), inicio)

		if !recursivo {
			hash.SumDirectorio(directorio, algo)
		} else {
			hash.SumRecursivo(directorio, algo)
		}

	}

}
