package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

var subcomando map[int]string
var mkdiskcomands map[int]string
var colorPurple string
var colorRed string
var colorCyan string
var colorBlanco string
var colorGreen string
var colorBlue string
var disk comMKDISK

func colorcitos() {
	colorRed = "\033[31m"
	colorGreen = "\033[32m"
	//	colorYellow := "\033[33m"
	colorBlue = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan = "\033[36m"
	colorBlanco = "\033[37m"
}
func main() {
	fmt.Print(colorBlanco, "Introduzca un comando----:: ")
	reader := bufio.NewReader(os.Stdin)
	entrada, _ := reader.ReadString('\n')
	eleccion := strings.TrimRight(entrada, "\r\n")
	Analizador(eleccion + "$$")

	for eleccion != "exit" {
		fmt.Print(colorBlanco, "Introduzca un comando----:: ")
		reader = bufio.NewReader(os.Stdin)
		entrada, _ = reader.ReadString('\n')
		eleccion = strings.TrimRight(entrada, "\r\n")
		Analizador(eleccion + "$$")
	}
}

func Analizador(cadena string) {
	colorcitos()
	estado := 0
	cadenita := ""
	lineaComando := ""
	escape := false
	for i := 0; i < len(cadena); i++ {
		caracter := string(rune(cadena[i]))

		switch estado {
		case 0:
			if cadena[i] == 32 || caracter == "\t" {
				estado = 0
			} else if cadena[i] >= 65 && cadena[i] <= 90 || cadena[i] >= 97 && cadena[i] <= 122 {
				cadenita += caracter
				estado = 1
			} else if cadena[i] >= 48 && cadena[i] <= 57 {
				estado = 2
				cadenita += caracter
			} else if cadena[i] == 47 {
				cadenita += caracter
				estado = 9
			} else if cadena[i] == 45 {
				estado = 8
				lineaComando += caracter
			} else if cadena[i] == 46 {
				estado = 0
				lineaComando += caracter
			} else if cadena[i] == 58 {
				estado = 0
				lineaComando += caracter
			} else if cadena[i] == 92 {
				estado = 4
			} else if cadena[i] == 34 {
				estado = 5
			} else if cadena[i] == 35 {
				estado = 7
				cadenita += caracter
				if lineaComando != "" {
					AnalizarLineaComando(lineaComando)
					lineaComando = ""
				}
			} else if caracter == "$" {
				if lineaComando != "" {
					AnalizarLineaComando(lineaComando)
					lineaComando = ""
				}

			} else if caracter == "\n" || escape == false {
				if lineaComando != "" {
					AnalizarLineaComando(lineaComando)
					lineaComando = ""
				}
			} else if caracter == "\n" || escape == true {
				estado = 0
			} else {
				fmt.Println(colorRed, "Caracter no reconocido "+caracter)
			}

			break
		case 1:
			if cadena[i] >= 65 && cadena[i] <= 90 || cadena[i] >= 97 && cadena[i] <= 122 ||
				cadena[i] >= 48 && cadena[i] <= 57 || cadena[i] == 95 || cadena[i] == 46 {
				cadenita += caracter
				estado = 1
			} else if cadena[i] == 47 {
				cadenita += caracter
				estado = 5
			} else if cadena[i] == 32 {
				lineaComando += cadenita + " "
				cadenita = ""
				estado = 0
			} else if len(cadena) == (i + 2) {
				lineaComando += cadenita + " "
				cadenita = ""
				estado = 0
			} else if caracter == "\n" {
				lineaComando += cadenita
				cadenita = ""
				estado = 0
				AnalizarLineaComando(lineaComando)
				lineaComando = ""
			} else {
				estado = 0
			}
			break
		case 2:
			if cadena[i] >= 48 && cadena[i] <= 57 {
				estado = 2
				cadenita += caracter
			} else if cadena[i] == 46 {
				estado = 3
				cadenita += caracter
			} else if cadena[i] == 47 {
				cadenita += caracter
				estado = 5
			} else if cadena[i] == 32 || cadena[i] == '\t' {
				estado = 0
				lineaComando += cadenita + " "
				cadenita = ""
			} else if len(cadena) == (i + 2) {
				lineaComando += cadenita
				cadenita = ""
				estado = 0
			} else if caracter == "\n" {
				lineaComando += cadenita
				cadenita = ""
				estado = 0
			} else {
				estado = 0
			}

			break
		case 3:
			if cadena[i] >= 48 && cadena[i] <= 57 {
				estado = 2
				cadenita += caracter
			} else if cadena[i] == 47 {
				cadenita += caracter
				estado = 5
			} else if cadena[i] == 32 || cadena[i] == '\t' || caracter == "\n" {
				estado = 0
				lineaComando += cadenita
				cadenita = ""
			} else {
				estado = 0
			}

			break
		case 4:
			if cadena[i] == 42 {
				escape = true
				estado = 0
				if string(rune(cadena[i+1])) == "\n" {
					i++
				}
			}
			break
		case 5:
			if cadena[i] == 47 {
				estado = 5
				cadenita += caracter
			} else if cadena[i] == 32 {
				cadenita += "@"

			} else if cadena[i] == 34 {
				estado = 0
				lineaComando += cadenita + " "
				cadenita = ""
			} else if caracter != "\n" && cadena[i] != 92 && (len(cadena) != (i + 2)) {
				estado = 5
				cadenita += caracter
			} else if cadena[i] == 92 {
				i++
				if cadena[i] == 42 {
					i++
				}
			}
			break
		case 7:

			if caracter != "\n" && (len(cadena) != (i + 2)) {
				cadenita += caracter
				estado = 7
			} else {
				fmt.Println(string(colorPurple), cadenita)
				cadenita = ""
				estado = 0
			}
			break
		case 8:
			if cadena[i] >= 65 && cadena[i] <= 90 || cadena[i] >= 97 && cadena[i] <= 122 {
				cadenita += caracter
				estado = 8
			} else if cadena[i] >= 48 && cadena[i] <= 57 {
				cadenita += caracter
				estado = 3
			} else if cadena[i] == 92 {
				lineaComando += cadenita
				cadenita = ""
				estado = 0
				i--
			} else if cadena[i] == 45 {
				cadenita += string(rune(cadena[i]))
				i++
				if cadena[i] == 62 {
					cadenita += string(rune(cadena[i]))
					lineaComando += cadenita
					cadenita = ""
					estado = 0
				}
			}
			break
		case 9:
			if cadena[i] == 92 {
				i--
				estado = 0
			} else if cadena[i] != 32 && (len(cadena) != (i + 2)) {
				cadenita += caracter
			} else {
				lineaComando += cadenita + " "
				cadenita = ""
				estado = 0
			}
			break
		}
	}
}
func CargaMasiva(direccion string) {
	file, err := os.Open(direccion)
	if err != nil {
		log.Fatal(err)

	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	texto := ""
	for scanner.Scan() {
		texto += scanner.Text() + "\n"
		//	fmt.Println(scanner.Bytes())
	}
	Analizador(texto)
}
func direccion(cadena string) string {
	cad := strings.Split(cadena, "->")
	direccion := ""
	if cad[0] == "-path" {
		if strings.Contains(cad[1], "@") {
			for h := 0; h < len(cad[1]); h++ {
				if cad[1][h] == 64 {
					direccion += " "
				} else {
					direccion += string(rune(cad[1][h]))
				}
			}
			return direccion
		} else {
			return cad[1]
		}
	} else {
		fmt.Println(colorRed, "Comando incorrecto, se esperaba -PATH")
	}
	return ""
}
func ValidarRuta(ruta string) bool {
	if _, err := os.Stat(ruta); err != nil {
		if os.IsNotExist(err) {
			fmt.Println(colorRed, "La ruta o archivo no existe")
			return true
		} else {
			fmt.Println(colorRed, "Error al verificar ruta")
			return true
		}

	}
	return false
}
func AnalizarLineaComando(cadena string) {
	arreglo := strings.Split(cadena, " ")
	switch strings.ToLower(arreglo[0]) {
	case "exec":
		fmt.Println(colorBlue, "analizando ruta...")
		direccion := direccion(arreglo[1])
		if !ValidarRuta(direccion) {
			CargaMasiva(direccion)
		}
		break
	case "mkdisk":
		MKDSIK(arreglo)
		break
	case "pause":
		fmt.Println("Precione Enter para continuar")
		var input string
		fmt.Scanln(&input)
		break
	}
}

func size(num string) int64 {
	numero, err := strconv.Atoi(num)
	if err != nil {
		fmt.Println(colorRed, "Tamaño incorrecto:", err)
	} else if numero >= 0 {
		return int64(numero)
	}
	return -1
}

type comMKDISK struct {
	name string
	tam  int64
	unit byte
	ext  string
}

func MKDSIK(cadena []string) {
	aux := 0
	err := false
	for i := 1; i < len(cadena); i++ {
		com := strings.Split(cadena[i], "->")
		if strings.ToLower(com[0]) == "-size" {
			disk.tam = size(com[1])
			if disk.tam != -1 {
				aux++
			} else {
				err = true
			}
		} else if strings.ToLower(com[0]) == "-path" {
			if strings.Contains(com[1], "@") {
				strings.ReplaceAll(com[1], "@", " ")
			}
			if AnalizarRuta(com[1]) {
				aux++
				disk.ext = com[1]
			} else {
				err = true
			}
		} else if strings.ToLower(com[0]) == "-name" {
			if VerificacionNombre(com[1]) {
				aux++
				disk.name = com[1]
			} else {
				err = true
			}

		} else if strings.ToLower(com[0]) == "-unit" {
			disk.unit = UNIT(com[1])
			if disk.unit != 'E' {
				aux++
			} else {
				err = true
			}
		}
	}
	if err == true {
		fmt.Println(colorRed, "Error en las características del disco")
	} else {
		if aux >= 3 {
			CrearDisco(disk.name, disk.tam, disk.unit, disk.ext)
		} else {
			fmt.Println(colorRed, "Falta un subcomando requerido")
		}
	}
}
func AnalizarRuta(direccion string) bool {

	_, error := os.Stat(direccion)
	if os.IsNotExist(error) {
		error = os.Mkdir(direccion, 0777)
		if error != nil {
			fmt.Println(colorRed, "Se ha producido un error al intentar acceder a la ruta")
			return false
		}
	}
	return true
}

func VerificacionNombre(nombre string) bool {
	for i := 0; i < len(nombre); i++ {
		if !(nombre[i] >= 48 && nombre[i] <= 57 || nombre[i] >= 65 && nombre[i] <= 90 ||
			nombre[i] >= 97 && nombre[i] <= 122 || nombre[i] == 95 || nombre[i] == 46) {
			return false
		}
	}
	extension := strings.Split(nombre, ".")
	if strings.ToLower(extension[1]) != "dsk" {
		return false
	}
	return true
}
func UNIT(unidad string) byte {
	if unidad == "m" {
		return 'm'
	} else if unidad == "k" {
		return 'k'
	} else {
		return 'E'
	}

}

func CrearDisco(nombre string, tam int64, unidad byte, ruta string) {

	file, err := os.Create(ruta + "/" + nombre)
	defer file.Close()
	if err != nil {
		fmt.Println(colorRed, err)
		return
	}

	size := 0
	if unidad == 'k' {
		size = 1024 * int(tam)
	} else {
		size = 1024 * 1024 * int(tam)
	}
	if size != 0 {
		var otro int8 = 0

		s := &otro
		size = size - 1
		fmt.Println(unsafe.Sizeof(otro))

		var binario bytes.Buffer
		binary.Write(&binario, binary.BigEndian, s)
		escribirBytes(file, binario.Bytes())

		file.Seek(int64(size), 0)
		var binario2 bytes.Buffer
		binary.Write(&binario2, binary.BigEndian, s)
		escribirBytes(file, binario2.Bytes())

		file.Seek(0, 0)
		CrearMBR(int64(size), file)
	}
}
func readNextBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

func escribirBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}
}

type MBR struct {
	mbr_tam            int64
	mbr_fecha_creacion [10]byte
	mbr_disk_id        int8
	part               particion
}

type particion struct {
}

func CrearMBR(mbr_tam int64, file *os.File) {
	mbr := MBR{}
	mbr.mbr_tam = mbr_tam
	mbr.mbr_disk_id = 1
	t := time.Now()
	fecha := string(t.Day()) + "-" + string(t.Month()) + "-" + string(t.Year())
	copy(mbr.mbr_fecha_creacion[:], fecha)
	var b2 bytes.Buffer
	binary.Write(&b2, binary.BigEndian, &mbr)
	escribirBytes(file, b2.Bytes())
}
