# go-artamiz
Calculador de hash de archivos

## Instalación desde fuente
Decargar el código fuente

cd $GOPATH
git get github.com/angelcoto/go-artamiz

El binario queda generado en $GOBIN

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
