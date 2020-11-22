# go-artamiz
Calculador de hash de archivos

## Instalación desde código fuente
Es requerido instalar Go (https://golang.org/doc/install)

Luego de asegurarse que Go está debidamente instalado, se ejecuta:

    cd $GOPATH
    go get github.com/angelcoto/go-artamiz

El binario generado en $GOBIN

## Uso
    go-artamiz -h
    Usage of go-artamiz:
    -a string
    	Archivo
    -d string
    	Directorio a recorrer (default <Directorio actual>)
    -m string
    	Algoritmo: md5, sha1, sha256 (default "sha256")
    -r	Recorrido recursivo
    -t string
    	Texto
