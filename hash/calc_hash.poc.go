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
	"sync"
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

// Definición del buffer para jobs.
type tjob struct {
	id   int
	path string
	algo string
}

// Definición del buffer para resultados
type tresultado struct {
	job  tjob
	hash []byte
	err  error
}

const totaljobs = 10
const totalworkers = 2
const totalresultados = 10

var jobs = make(chan tjob, totaljobs)
var resultados = make(chan tresultado, totalresultados)

// workerHash lee el buffer "jobs" para obtener el próximo job a ejecutar
// Cuando la tarea ha sido ejecutada el resultado se escribe en el buffer "resultado"
// El worker seguirá ejecutándose mientras el buffer no esté cerrado.
// Cuando el buffer se cierra el worker se declara como finalizado a través de
// wg.Done()
func workerHash(wg *sync.WaitGroup) {
	//Recorre el buffer de jobs
	for job := range jobs {
		hash, err := SumArchivo(job.path, job.algo)
		resultado := tresultado{job, hash, err}
		resultados <- resultado
	}
	wg.Done()
}

// imprimeSalida lee el buffer de resultados para imprimir cada resultado a terminal.
// La rutina se mantiene en ejecución hasta que el buffer es cerrado.
func imprimeSalida(done chan bool) {
	for resultado := range resultados {
		if resultado.err != nil {
			fmt.Printf("* Error: %s\n", resultado.err)
		} else {
			fmt.Printf("%x *%s\n", resultado.hash, resultado.job.path)
		}
	}
	done <- true // Informa a la función de llamado que el trabajo ha finalizado
}

// creaWorkerPool inicia los worker que estarán leyendo la cola de jobs
func creaWorkerPool(nWorkers int) {
	var wg sync.WaitGroup
	for i := 0; i < nWorkers; i++ {
		wg.Add(1)
		go workerHash(&wg)
	}
	wg.Wait()
	close(resultados) // Para indicarle al lector del buffer que no hay más valores a enviar
}

// SumRecursivo imprime el hash para los archivos de un directorio,
// incluyendo los subdirectorios.
//
// Esta versión implementa el patrón Worker Pools con la utilización de
// buffered channels y goroutines.  Con esta implementación se logra una
// ejecución asíncrona de las funciones del procesos, manteniendo la
// comunicación entre las funciones a través los buffers.
//
// IMPORTANTE: Esta versión existe como prueba de concepto y como un ejercicio
// para la utilización de concurrencia, pero las pruebas de desempeño arrojaron que
// el tiempo de ejecución de esta versión es mejor que la versión síncrona solamente
// cuando la cantidad de jobs es pequeña.  En casos de muchos jobs (decena a miles) esta
// versión mostro ofrecer un desempeño un poco inferior cuando se utiliza un solo worker,
// y significativamente inferior cuando se utiliza más de un worker.
// Posiblemente en escenarios de sistemas de almacenamiento con arreglos de discos la versión
// asíncrona ofrezca un mejor desempeño, pero en escenario de un solo disco el desempeño
// fue decepcionante.
func SumRecursivo(dir string, algo string) {

	fmt.Println("Con goroutine")

	done := make(chan bool)

	//
	go func() {

		i := 0
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			//fmt.Println(path, info.IsDir(), algo)
			if !info.IsDir() {
				job := tjob{i, path, algo}
				jobs <- job
				i++
			}
			return nil
		})
		close(jobs) // Para indicarle al lector del buffer que no hay más valores a enviar
	}()

	go imprimeSalida(done)
	creaWorkerPool(totalworkers)

	<-done // Genera un bloqueo hasta que ha finalizado la impresión de todos los resultados

}
