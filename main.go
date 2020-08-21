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
var rutita string
var colorYellow string
var mbr MBR

func colorcitos() {
	colorRed = "\033[31m"
	colorGreen = "\033[32m"
	colorYellow = "\033[33m"
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
		fmt.Print(colorBlanco, "\nIntroduzca un comando----:: ")
		reader = bufio.NewReader(os.Stdin)
		entrada, _ = reader.ReadString('\n')
		eleccion = strings.TrimRight(entrada, "\r\n")
		Analizador(eleccion + "$$")
	}
}

//Analizador funcion que analiza todo el texto
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
			} else if cadena[i] == 92 {
				lineaComando += cadenita
				cadenita = ""
				estado = 0
				i--
			} else {
				estado = 0
				i--
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
				i--
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
					lineaComando += cadenita + " "
					cadenita = ""
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
					lineaComando += " "
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
			} else if caracter == "\n" {
				lineaComando += cadenita
				cadenita = ""
				estado = 0
				i--
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

//CargaMasiva función para cargar datos
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
		}
		return cad[1]
	}
	fmt.Println(colorRed, "Comando incorrecto, se esperaba -PATH")
	return ""
}

//ValidarRuta valida si la dirección es correcta
func ValidarRuta(ruta string) bool {
	fmt.Println(colorYellow, ruta)
	if _, err := os.Stat(ruta); err != nil {
		if os.IsNotExist(err) {
			fmt.Println(colorRed, "La ruta o archivo no existe")
			return true
		}
		fmt.Println(colorRed, "Error al verificar ruta")
		return true
	}
	return false
}

//AnalizarLineaComando esta verifica cada linea de comando enviada por analizador
func AnalizarLineaComando(cadena string) {
	fmt.Println(cadena)
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
		fmt.Println(colorBlue, "Creando disco...")
		MKDISK(arreglo)
		break
	case "pause":
		fmt.Println("Precione Enter para continuar")
		var input string
		fmt.Scanln(&input)
		break
	case "rmdisk":
		fmt.Println(colorBlue, "Verificando requisitos para eliminación...")
		RMDISK(arreglo[1])
		break
	case "fdisk":
		FDISK(arreglo)
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

//MKDISK SE USA PARA COMPROBAR LOS COMANDOS
func MKDISK(cadena []string) {
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
			direccion := com[1]
			if strings.Contains(direccion, "@") {
				direccion = strings.ReplaceAll(direccion, "@", " ")
			}
			if AnalizarRuta(direccion) {
				aux++
				disk.ext = direccion
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
			if aux == 3 {
				disk.unit = 'm'
			}
			CrearDisco(disk.name, disk.tam, disk.unit, disk.ext)
		} else {
			fmt.Println(colorRed, "Falta un subcomando requerido")
		}
	}
}

//AnalizarRuta sirve para comprobar que la ruta existe
func AnalizarRuta(direccion string) bool {
	carpetas := strings.Split(direccion, "/")
	directorio := ""
	for i := 0; i < len(carpetas); i++ {
		directorio += "/" + carpetas[i]
		_, error := os.Stat(directorio)
		if os.IsNotExist(error) {
			error = os.MkdirAll(direccion, 0777)
			if error != nil {
				fmt.Println(colorRed, "Se ha producido un error al intentar acceder a la ruta")
				return false
			}
		}
	}
	return true
}

//VerificacionNombre sirve para comprobar que el nombre es correcto
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

//UNIT funcion que verifica si la unidad es correcta
func UNIT(unidad string) byte {
	if unidad == "m" {
		return 'm'
	} else if unidad == "k" {
		return 'k'
	} else {
		return 'E'
	}

}

// CrearDisco crea el archivo binario verificando cada uno de sus atributos
func CrearDisco(nombre string, tam int64, unidad byte, ruta string) {
	rutita = ruta + "/" + nombre
	file, err := os.Create(rutita)
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
		var cero int8 = 0

		size = size - 1
		var binario bytes.Buffer
		binary.Write(&binario, binary.BigEndian, &cero)
		escribirBytes(file, binario.Bytes())

		file.Seek(int64(size), 0)

		var binario2 bytes.Buffer
		binary.Write(&binario2, binary.BigEndian, &cero)
		escribirBytes(file, binario2.Bytes())

		file.Seek(0, 0)
		CrearMBR(int64(size), file)
		fmt.Println(colorGreen, "Nombre del disco: "+nombre)
		AbrirArchivo()
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

//MBR lleva todos los datos que requiere el mbr
type MBR struct {
	MbrTam           int64
	MbrFechaCreacion [19]byte
	MbrDiskID        uint8
	Part1            particion
	Part2            particion
	Part3            particion
	Part4            particion
}

//particion información de cada partición en el archivo
type particion struct {
	PartStatus byte
	PartType   byte
	PartFit    byte
	PartStart  int64
	PartSize   int64
	PartName   [16]byte
	PartDelete bool
}

//CrearMBR aquí escribe el mbr en el archivo binario
func CrearMBR(mbrTam int64, file *os.File) {
	mbr = MBR{}
	mbr.MbrTam = mbrTam
	mbr.MbrDiskID = 1
	mbr.Part1 = particion{PartStatus: 73}
	mbr.Part2 = particion{PartStatus: 73}
	mbr.Part3 = particion{PartStatus: 73}
	mbr.Part4 = particion{PartStatus: 73}
	t := time.Now()
	fecha := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	copy(mbr.MbrFechaCreacion[:], fecha)
	/*fmt.Print(colorGreen, " Fecha de creación: ")
	for i := 0; i < 10; i++ {
		fmt.Print(colorGreen, string(mbr.MbrFechaCreacion[i]))
	}
	fmt.Println(colorGreen, " ")*/
	/*tami := fmt.Sprintf("%d", mbr.MbrTam)
	fmt.Println(colorGreen, "Tamaño: "+tami+" bytes")*/
	tamMBR := int64(unsafe.Sizeof(mbr))
	mbr.Part1.PartSize = mbr.MbrTam - tamMBR
	mbr.Part1.PartStart = tamMBR + 1
	//agrega el mbr al disco
	var b2 bytes.Buffer
	binary.Write(&b2, binary.BigEndian, &mbr)
	escribirBytes(file, b2.Bytes())

}

//AbrirArchivo Se abre el disco para leer el MBR
func AbrirArchivo() {
	file, err := os.Open(rutita)
	defer file.Close()
	if err != nil { //validar que no sea nulo.
		log.Fatal(err)
	}
	mbr2 := MBR{}
	var size int = int(unsafe.Sizeof(mbr2))
	file.Seek(0, 0)
	data := readNextBytes(file, size)
	buffer := bytes.NewBuffer(data)

	err = binary.Read(buffer, binary.BigEndian, &mbr2)
	if err != nil {
		panic(err)
	}

	//fmt.Println(mbr2)
	fmt.Printf("%s%d%s%s%s%d", "Tamaño del disco: ", mbr2.MbrTam, "Fecha de creación: ", mbr2.MbrFechaCreacion[:], "ID de disco: ", mbr2.MbrDiskID)
	fmt.Println()
}

//ActualizarMBR en este metodo se va actualizando la información de las particones

//RMDISK ES PARA ELIMINAR EL ARCHIVO
func RMDISK(direc string) {
	dir := direccion(direc)

	if dir != "" {
		ext := strings.Split(dir, ".")
		if ext[1] == "dsk" {
			err := os.Remove(dir)
			if err != nil {
				fmt.Println(colorRed, "Error al intentar eliminar archivo")
			}
			fmt.Println(colorGreen, "Success, archivo eliminado")
		} else {
			fmt.Println(colorRed, "El disco a eliminar debe ser .dsk")
		}

	}
}

//FDISK administra las particiones del disco
func FDISK(subcomandos []string) {
	aux := 0
	err := false
	tam := 1024
	tamanio := 0
	dir := ""
	tipo := "p"
	fit := "wf"
	delete := ""
	name := ""
	add := 0
	for i := 1; i < len(subcomandos); i++ {
		subcadena := strings.Split(subcomandos[i], "->")
		switch strings.ToLower(subcadena[0]) {

		case "-size":
			tamanio = int(size(subcadena[1]))
			if tamanio != -1 {
				aux++
			} else {
				fmt.Println(colorRed, "Imposible crear una partición del tamaño solicitado")
				err = true
			}
			break
		case "-unit":
			tam = UNITFDISK(subcadena[1])
			if tam != -1 {
				aux++
			} else {
				fmt.Println(colorYellow, "Verificar el parámetro de unit")
				err = true
			}
			break
		case "-path":
			dir = direccion(subcomandos[i])
			if dir != "" {
				aux++
			} else {
				err = true
			}
			break
		case "-type":
			if TYPE(subcadena[1]) {
				aux++
				tipo = subcadena[1]
			} else {
				fmt.Println(colorRed, "Parámetro del comando -type incorrecto")
				err = true
			}
			break
		case "-fit":
			if FIT(subcadena[1]) {
				fit = subcadena[1]
				aux++
			} else {
				fmt.Println(colorRed, "El parámetro del comando fit es incorrecto")
				err = true
			}
			break
		case "-delete":
			if DELETE(subcadena[1]) {
				aux++
				delete = subcadena[1]
			} else {
				fmt.Println(colorRed, "El parámetro del comando delete es incorrecto")
				err = true
			}
			break
		case "-name":
			if !ExisteNombreParticion(subcadena[1]) {
				aux++
				name = subcadena[1]
			} else {
				fmt.Println(colorYellow, "El nombre de la partición ya existe!")
				err = true
			}
			break
		case "-add":
			numero, correcto := VerificarNumero(subcadena[1])
			if correcto == true {
				aux++
				add = int(numero)
			} else {
				err = true
			}
			break
		}

	}
	if err == false {
		if aux >= 3 {
			if delete != "" {

			} else if add != 0 {

			} else if dir != "" && name != "" && tamanio != 0 {

				CrearParticionNueva(int64(tamanio), int64(tam), dir, tipo, fit, name)
			}
		} else {
			fmt.Println(colorYellow, "Faltan parámetros requeridos!")
		}
	}
}

//ExisteNombreParticion busca en el mbr si hay una particion con el mismo nombre
func ExisteNombreParticion(nombre string) bool {

	return false
}

//VerificarNumero identifica un número negativo o positivo
func VerificarNumero(num string) (int64, bool) {
	numero, err := strconv.Atoi(num)
	if err != nil {
		fmt.Println(colorRed, "Tamaño incorrecto:", err)
	} else {
		//FALTA VERIFICAR SI HAY ESPACIO
		return int64(numero), true
	}
	return 0, false
}

//UNITFDISK verifica el tamaño del fdisk
func UNITFDISK(unidad string) int {
	unidad = strings.ToLower(unidad)
	switch unidad {
	case "k":
		return 1024
	case "b":
		return 1

	case "m":
		return 1024 * 1024
	}
	return -1
}

//TYPE verifica si el parametro es correcto
func TYPE(tipo string) bool {
	tipo = strings.ToLower(tipo)
	if tipo == "p" || tipo == "e" || tipo == "l" {
		return true
	}
	return false
}

//FIT verifica los parametros para las particiones
func FIT(fit string) bool {
	fit = strings.ToLower(fit)
	if fit == "bf" || fit == "ff" || fit == "wf" {
		return true
	}
	return false
}

//DELETE verifica que los comandos de DELETE sean correctos
func DELETE(delete string) bool {
	delete = strings.ToLower(delete)
	if delete == "fast" || delete == "full" {
		return true
	}
	return false
}

//CrearParticionNueva crea una particion nueva en el disco
func CrearParticionNueva(size int64, unidad int64, path string, tipo string, fit string, name string) {
	mbr = LeerMBR()
	size = size * unidad
	var s int64
	var part int
	s, part = PrimerAjuste(size, strings.ToLower(tipo))
	if s != 0 && part != 0 {
		InformacionParticion(name, tipo, fit, size, s, part)
		CrearParticion(path)
	}

}

//InformacionParticion en este metodo se agrega toda la información al struct particion
func InformacionParticion(name string, tipo string, fit string, size int64, start int64, numero int) {
	nuevofit := ' '
	if strings.ToLower(fit) == "bf" {
		nuevofit = 'b'
	} else if strings.ToLower(fit) == "ff" {
		nuevofit = 'f'
	} else if strings.ToLower(fit) == "wf" {
		nuevofit = 'w'
	}
	nuevoTipo := ' '
	if strings.ToLower(tipo) == "p" {
		nuevoTipo = 'p'
	} else if strings.ToLower(tipo) == "e" {
		nuevoTipo = 'e'
	}
	if numero == 1 {
		copy(mbr.Part1.PartName[:], name)
		mbr.Part1.PartSize = size
		mbr.Part1.PartStart = start
		mbr.Part1.PartFit = byte(nuevofit)
		mbr.Part1.PartType = byte(nuevoTipo)
		mbr.Part1.PartStatus = 65

	} else if numero == 2 {
		copy(mbr.Part2.PartName[:], name)
		mbr.Part2.PartSize = size
		mbr.Part2.PartStart = start
		mbr.Part2.PartFit = byte(nuevofit)
		mbr.Part2.PartType = byte(nuevoTipo)
		mbr.Part2.PartStatus = 65
	} else if numero == 3 {
		copy(mbr.Part3.PartName[:], name)
		mbr.Part3.PartSize = size
		mbr.Part3.PartStart = start
		mbr.Part3.PartFit = byte(nuevofit)
		mbr.Part3.PartType = byte(nuevoTipo)
		mbr.Part3.PartStatus = 65
	} else if numero == 4 {
		copy(mbr.Part4.PartName[:], name)
		mbr.Part4.PartSize = size
		mbr.Part4.PartStart = start
		mbr.Part4.PartFit = byte(nuevofit)
		mbr.Part4.PartType = byte(nuevoTipo)
		mbr.Part4.PartStatus = 65
	}

}

//CrearParticion crea la particion
func CrearParticion(path string) {
	files, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, os.ModeAppend)
	defer files.Close()
	if err != nil {
		log.Fatal(err)
	}
	files.Seek(0, 0)
	var b3 bytes.Buffer
	binary.Write(&b3, binary.BigEndian, &mbr)
	escribirBytes(files, b3.Bytes())
}

//PrimerAjuste este metodo devuelve la posicion inicial del primer espacio que encuentre
func PrimerAjuste(tam int64, tipo string) (int64, int) {

	if mbr.Part1.PartStatus == 73 && mbr.Part1.PartSize >= tam {
		mbr.Part1.PartSize = tam
		if mbr.Part1.PartDelete == false {
			mbr.Part2.PartStart = int64(unsafe.Sizeof(mbr)) + tam
			mbr.Part2.PartSize = mbr.MbrTam - int64(unsafe.Sizeof(mbr)) - tam
		}
		return mbr.Part1.PartStart, 1
	} else if mbr.Part2.PartStatus == 73 && mbr.Part2.PartSize >= tam {
		if mbr.Part1.PartType == 101 && tipo == "e" {
			fmt.Println(colorRed, "Ya existe una partición extendida")
			return 0, 0
		}
		mbr.Part2.PartSize = tam
		if mbr.Part1.PartDelete == false {
			mbr.Part3.PartStart = (int64(unsafe.Sizeof(mbr))) + mbr.Part1.PartSize + mbr.Part2.PartSize
			mbr.Part3.PartSize = mbr.MbrTam - mbr.Part1.PartSize - mbr.Part2.PartSize
		}
		return mbr.Part2.PartStart, 2
	} else if mbr.Part3.PartStatus == 73 && mbr.Part3.PartSize >= tam {
		if (mbr.Part1.PartType == 101 || mbr.Part2.PartType == 101) && tipo == "e" {
			fmt.Println(colorRed, "Ya existe una partición extendida")
			return 0, 0
		}
		mbr.Part3.PartSize = tam
		if mbr.Part1.PartDelete == false {
			mbr.Part4.PartStart = (int64(unsafe.Sizeof(mbr))) + mbr.Part3.PartSize + mbr.Part2.PartSize + mbr.Part1.PartSize
			mbr.Part4.PartSize = mbr.MbrTam - mbr.Part1.PartSize - mbr.Part2.PartSize - mbr.Part3.PartSize
		}
		return mbr.Part3.PartStart, 3
	} else if mbr.Part4.PartStatus == 73 && mbr.Part4.PartSize >= tam {
		if (mbr.Part1.PartType == 101 || mbr.Part2.PartType == 101 || mbr.Part3.PartType == 101) && tipo == "e" {
			fmt.Println(colorRed, "Ya existe una partición extendida")
			return 0, 0
		}
		mbr.Part4.PartSize = tam
		return mbr.Part4.PartStart, 4
	}
	fmt.Println(colorYellow, "No hay espacio en la partición")
	return 0, 0
}

//LeerMBR este metodo devuelve el mbr actual del disco
func LeerMBR() MBR {
	file, err := os.Open(rutita)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	mbr2 := MBR{}
	var size int = int(unsafe.Sizeof(mbr2))
	file.Seek(0, 0)
	data := readNextBytes(file, size)
	buffer := bytes.NewBuffer(data)

	err = binary.Read(buffer, binary.BigEndian, &mbr2)
	if err != nil {
		panic(err)
	}
	fmt.Println(mbr2)
	return mbr2
}
