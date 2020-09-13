package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
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
var ebr EBR
var contador = 0

// ListDiscos inicio de la lista
var ListDiscos ListaDisco

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
	f := "mount -path->/home/josselyn/Escritorio/archivoBinario/disco.dsk -name->particion8\n"
	f += "rep -id->vda1 -Path->/home/Prueba/reporteDisk1.png -name->disk"
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
	colorcitos()
	estado := 0
	cadenita := ""
	lineaComando := ""
	escape := false
	comilla := false
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
				if cadena[i+1] >= 48 && cadena[i+1] <= 57 {
					estado = 3
					cadenita += caracter
				} else {
					estado = 8
					lineaComando += caracter
				}
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
				comilla = true
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
				lineaComando += cadenita + " "
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
			} else if cadena[i] == 32 && comilla == true {
				cadenita += "@"
			} else if cadena[i] == 32 && comilla == false {
				estado = 0
				lineaComando += cadenita + " "
				cadenita = ""
			} else if cadena[i] == 34 {
				estado = 0
				lineaComando += cadenita + " "
				cadenita = ""
				comilla = false
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

			if caracter != "\n" && (i+1) != len(cadena) {
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
				estado = 8
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
			} else if caracter == " " {
				cadenita += " "
				estado = 0
				lineaComando += cadenita
				cadenita = ""
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
	fmt.Println(colorBlue, "Leyendo archivo de entrada ubicado en la dirección: "+ruta)
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
func duracion() {
	duration := time.Duration(1) * time.Second
	time.Sleep(duration)
}

//AnalizarLineaComando esta verifica cada linea de comando enviada por analizador
func AnalizarLineaComando(cadena string) {
	fmt.Println(colorCyan, "*****Comando detectado*****")
	arreglo := strings.Split(cadena, " ")
	switch strings.ToLower(arreglo[0]) {
	case "exec":
		direccion := direccion(arreglo[1])
		if !ValidarRuta(direccion) {
			fmt.Println(colorBlue, "analizando ruta...")
			duracion()
			fmt.Println(colorCyan, cadena)
			fmt.Println(colorCyan, "***************************")
			duracion()
			CargaMasiva(direccion)
		}
		break
	case "mkdisk":
		fmt.Println(colorCyan, cadena)
		fmt.Println(colorCyan, "***************************")
		duracion()
		fmt.Println(colorBlue, "Creando disco...")
		MKDISK(arreglo)
		break
	case "pause":
		fmt.Println(colorCyan, cadena)
		fmt.Println(colorCyan, "***************************")
		duracion()
		fmt.Println(colorBlue, "Presione una tecla para continuar...")
		fmt.Scanln()
		break
	case "rmdisk":
		fmt.Println(colorCyan, cadena)
		fmt.Println(colorCyan, "***************************")
		RMDISK(arreglo[1])

		break
	case "fdisk":
		fmt.Println(colorCyan, cadena)
		fmt.Println(colorCyan, "***************************")
		duracion()
		FDISK(arreglo)
		break
	case "mount":
		fmt.Println(colorCyan, cadena)
		fmt.Println(colorCyan, "***************************")
		fmt.Println(len(arreglo))
		if len(arreglo) <= 2 {
			listarMontadas()
		} else {
			Mount(arreglo)
		}
		break
	case "unmount":
		fmt.Println(colorCyan, cadena)
		fmt.Println(colorCyan, "***************************")
		Unmount(arreglo)
		break
	case "mkfs":
		fmt.Println(colorCyan, cadena)
		fmt.Println(colorCyan, "***************************")
		fmt.Println(colorBlue, "Verificando requisitos de formateo...")
		fmt.Println(colorRed, "¿Seguro que desea formatear el disco? s/n")
		reader := bufio.NewReader(os.Stdin)
		entrada, _ := reader.ReadString('\n')
		eleccion := strings.TrimRight(entrada, "\r\n")

		if eleccion == "s" {
			MKFS(arreglo)
		} else {
			fmt.Println(colorCyan, "Formateo cancelado.")
			return
		}
		break
	case "mkdir":
		fmt.Println(colorCyan, cadena)
		fmt.Println(colorCyan, "***************************")
		MkDir(arreglo)
		break
	case "mkfile":
		fmt.Println(colorCyan, cadena)
		fmt.Println(colorCyan, "***************************")
		MkFile(arreglo)
		break
	case "rep":
		fmt.Println(colorCyan, cadena)
		fmt.Println(colorCyan, "***************************")
		REPORTE(arreglo)
		break
	default:
		fmt.Println(colorYellow, "Comandos no reconocidos...")
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
		} else if com[0] != "" {
			fmt.Println(colorRed, "Comando no permitido!")
			return
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
	if len(extension) > 0 {
		if strings.ToLower(extension[1]) != "dsk" {
			return false
		}
	} else {
		return false
	}
	return true
}
func verificarNombreParticion(nombre string) bool {
	for i := 0; i < len(nombre); i++ {
		if !(nombre[i] >= 48 && nombre[i] <= 57 || nombre[i] >= 65 && nombre[i] <= 90 ||
			nombre[i] >= 97 && nombre[i] <= 122 || nombre[i] == 95 || nombre[i] == 46) {
			return false
		}
	}
	return true
}

//UNIT funcion que verifica si la unidad es correcta
func UNIT(unidad string) byte {
	unidad = strings.ToLower(unidad)
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

		file.Seek(int64(size), 0) // 0 inicio del archivo, pos 0, 1->donde se quedo, 2->al final del archivo

		var binario2 bytes.Buffer
		binary.Write(&binario2, binary.BigEndian, &cero)
		escribirBytes(file, binario2.Bytes())
		file.Seek(0, 0)
		CrearMBR(int64(size)+1, file)
		duracion()
		fmt.Println(colorGreen, "*****Información del disco creado*****")
		fmt.Println(colorGreen, "Nombre del disco: "+nombre)
		AbrirArchivo()
	}
}
func readNextBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		fmt.Println(colorRed, "No hay bytes que leer")
		return bytes
	}

	return bytes
}

func escribirBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
		return
	}
}

//MBR lleva todos los datos que requiere el mbr
type MBR struct {
	MbrTam           int64
	MbrFechaCreacion [19]byte
	MbrDiskID        uint8
	MbrRecorrido     int64
	Particiones      [4]particion
	MbrActivas       byte
}

//particion información de cada partición en el archivo
type particion struct {
	PartStatus    byte
	PartType      byte
	PartFit       byte
	PartStart     int64
	PartSize      int64
	PartPartition bool
	PartName      [16]byte
	PartDelete    bool
	PartUnida     bool
}

//EBR contenido del EBR
type EBR struct {
	PartStatus   byte
	PartFit      byte
	PartStart    int64
	PartSize     int64
	PartNext     int64
	PartName     [16]byte
	PartPrevious int64
	PartDelete   bool
}

//CrearMBR aquí escribe el mbr en el archivo binario
func CrearMBR(mbrTam int64, file *os.File) {
	mbr = MBR{}
	mbr.MbrTam = mbrTam
	var n uint8
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	mbr.MbrDiskID = n
	for i := 0; i < 4; i++ {
		mbr.Particiones[i] = particion{PartStatus: 73}
	}
	t := time.Now()
	fecha := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	copy(mbr.MbrFechaCreacion[:], fecha)
	tamMBR := int64(unsafe.Sizeof(mbr))

	mbr.Particiones[0].PartSize = mbr.MbrTam - tamMBR
	mbr.Particiones[0].PartStart = tamMBR
	//agrega el mbr al disco
	var b2 bytes.Buffer
	binary.Write(&b2, binary.BigEndian, &mbr)
	escribirBytes(file, b2.Bytes())

}

//AbrirArchivo Se abre el disco para leer el MBR
func AbrirArchivo() {
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
	fmt.Println(colorGreen)
	fmt.Printf("%s%d", " Tamaño del disco: ", mbr2.MbrTam)
	fmt.Println()
	fmt.Println(" Fecha de creación: " + BytesToString(mbr2.MbrFechaCreacion))
	fmt.Printf("%s%d", " ID de disco: ", mbr2.MbrDiskID)
	fmt.Println()
	duracion()
}

//BytesToString convierte un array de bytes a cadena
func BytesToString(datos [19]byte) string {
	cadena := ""
	for i := 0; i < len(datos); i++ {
		cadena += string(rune(datos[i]))
	}
	return cadena
}

//ActualizarMBR en este metodo se va actualizando la información de las particones

//RMDISK ES PARA ELIMINAR EL ARCHIVO
func RMDISK(direc string) {
	dir := direccion(direc)

	if dir != "" {
		if !ValidarRuta(dir) {
			fmt.Println(colorBlue, "Verificando requisitos para eliminación...")
			fmt.Println(colorRed, "¿Seguro que desea eliminar el disco? s/n")
			reader := bufio.NewReader(os.Stdin)
			entrada, _ := reader.ReadString('\n')
			eleccion := strings.TrimRight(entrada, "\r\n")
			if eleccion == "s" {
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
			} else {
				fmt.Println(colorCyan, "Eliminación cancelada.")
				return
			}
		}
	}
}

//FDISK administra las particiones del disco
func FDISK(subcomandos []string) {
	aux := 0
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
		analiza := strings.ToLower(subcadena[0])
		switch analiza {

		case "-size":
			tamanio = int(size(subcadena[1]))
			if tamanio != -1 {
				aux++
			} else {
				fmt.Println(colorRed, "Imposible crear una partición del tamaño solicitado")
				return
			}
			break
		case "-unit":
			tam = UNITFDISK(subcadena[1])
			if tam != -1 {
				aux++
			} else {
				fmt.Println(colorYellow, "Verificar el parámetro de unit")
				return
			}
			break
		case "-path":
			dir = direccion(subcomandos[i])
			if ExisteArchivo(dir) {
				aux++
			} else {
				return
			}
			break
		case "-type":
			if TYPE(subcadena[1]) {
				aux++
				tipo = subcadena[1]
			} else {
				fmt.Println(colorRed, "Parámetro del comando -type incorrecto")
				return
			}
			break
		case "-fit":
			if FIT(subcadena[1]) {
				fit = subcadena[1]
				aux++
			} else {
				fmt.Println(colorRed, "El parámetro del comando fit es incorrecto")
				return
			}
			break
		case "-delete":
			if DELETE(subcadena[1]) {
				delete = subcadena[1]
				aux++
			} else {
				fmt.Println(colorRed, "El parámetro del comando delete es incorrecto")
				return
			}
			break
		case "-name":
			if verificarNombreParticion(subcadena[1]) {
				aux++
				name = subcadena[1]
			} else {
				fmt.Println(colorYellow, "El nombre de la partición no tiene el formato correcto!")
				return
			}
			break
		case "-add":
			numero, correcto := VerificarNumero(subcadena[1])
			if correcto == true {
				aux++
				add = int(numero)
			} else {
				return
			}
			break
		case "":
			break
		default:
			fmt.Println(colorYellow, "Parámetro no reconocido!")
			return
		}

	}
	if aux >= 3 {
		if delete != "" && dir != "" && name != "" {
			fmt.Println(colorBlue, "Verificando requisitos para eliminación...")
			fmt.Println(colorRed, "¿Seguro que desea eliminar la partición? s/n")
			reader := bufio.NewReader(os.Stdin)
			entrada, _ := reader.ReadString('\n')
			eleccion := strings.TrimRight(entrada, "\r\n")
			if eleccion == "s" {
				EliminarParticion(dir, name, delete)
			} else {
				fmt.Println(colorCyan, "Eliminación cancelada.")
				return
			}
		} else if add != 0 && dir != "" && name != "" {
			AgregarOQuitar(dir, int64(add), name, int64(tam))
		} else if dir != "" && name != "" && tamanio != 0 {
			CrearParticionNueva(int64(tamanio), int64(tam), dir, tipo, fit, name)
			//	graphic(dir)
		}
	} else {
		fmt.Println(colorYellow, "Faltan parámetros requeridos!")
	}
}

//AgregarOQuitar este metodo agrega o quita espacio de una particion
func AgregarOQuitar(path string, add int64, name string, unidades int64) {
	b := false
	add = add * unidades
	mbr, b = LeerMBR(path)
	if b == true {
		for i := 0; i < len(mbr.Particiones); i++ {
			nombreParticion := Nombres(mbr.Particiones[i].PartName)
			if nombreParticion == name {
				if add < 0 {
					if mbr.Particiones[i].PartType == byte('e') {
						ebr = ExtraerEBR(path, mbr.Particiones[i].PartStart)
						for ebr.PartNext != -1 {
							ebr = ExtraerEBR(path, ebr.PartNext)
						}
						if ebr.PartStatus == 73 {
							if ebr.PartSize > (-1 * add) {
								ebr.PartSize += add
								mbr.Particiones[i].PartSize += add
								EscribirMBR(path)
								EscribirEBR(ebr.PartStart, path)
								mensajeCreado(path, nombreParticion, mbr.Particiones[i].PartSize-add, mbr.Particiones[i].PartSize)
								return
							}
						} else {
							espacio := (mbr.Particiones[i].PartStart + mbr.Particiones[i].PartSize) - (ebr.PartStart + ebr.PartSize)
							if espacio >= (-1 * add) {
								mbr.Particiones[i].PartSize += add
								EscribirMBR(path)
								mensajeCreado(path, nombreParticion, mbr.Particiones[i].PartSize-add, mbr.Particiones[i].PartSize)
								return
							}
						}
						fmt.Println(colorYellow, "******************************************************")
						fmt.Println(colorYellow, "No hay suficiente espacio para reducir la partición")
						fmt.Println(colorYellow, "******************************************************")
						return

					}
					if mbr.Particiones[i].PartSize > (-1 * add) {
						mbr.Particiones[i].PartSize += add
						EscribirMBR(path)
						mensajeCreado(path, nombreParticion, mbr.Particiones[i].PartSize-add, mbr.Particiones[i].PartSize)
						return
					}

				}
				termina := mbr.Particiones[i].PartStart + mbr.Particiones[i].PartSize
				if i < 3 {
					if termina != mbr.Particiones[i+1].PartStart {
						libre := mbr.Particiones[i+1].PartStart - termina
						if add <= libre {
							if mbr.Particiones[i].PartType == byte('e') {
								ebr = ExtraerEBR(path, mbr.Particiones[i].PartStart)
								for ebr.PartNext != -1 {
									ebr = ExtraerEBR(path, ebr.PartNext)
								}
								if ebr.PartStatus == 73 {
									ebr.PartSize += add
								}
								EscribirEBR(ebr.PartStart, path)
							}
							mbr.Particiones[i].PartSize += add
							EscribirMBR(path)
							mensajeCreado(path, nombreParticion, mbr.Particiones[i].PartSize-add, mbr.Particiones[i].PartSize)
							return
						}
					}
				} else {
					libre := mbr.MbrTam - mbr.Particiones[i].PartStart - mbr.Particiones[i].PartSize
					if add <= libre {
						mbr.Particiones[i].PartSize += add
						if mbr.Particiones[i].PartType == byte('e') {
							ebr = ExtraerEBR(path, mbr.Particiones[i].PartStart)
							for ebr.PartNext != -1 {
								ebr = ExtraerEBR(path, ebr.PartNext)
							}
							if ebr.PartStatus == 73 {
								ebr.PartSize += add
							}
							EscribirEBR(ebr.PartStart, path)
						}
						EscribirMBR(path)
						mensajeCreado(path, nombreParticion, mbr.Particiones[i].PartSize-add, mbr.Particiones[i].PartSize)
						return
					}
				}

				fmt.Println(colorYellow, "******************************************************")
				fmt.Println(colorYellow, "No hay suficiente espacio para añadir a la partición")
				fmt.Println(colorYellow, "******************************************************")
				return
			}
		}
		///////////////////////////////////////////////////////////////
		b := false
		mbr, b = LeerMBR(path)
		if b == true {
			start, size := BuscarExtendida()
			if start != -1 {
				ebr = ExtraerEBR(path, start)
				nombreParticion := Nombres(ebr.PartName)
				AgregarOQuitarLogicas(path, add, name, unidades, nombreParticion, size, start)
				return
			}
			fmt.Println(colorRed, "No existe una partición extendida para agregar a una partición lógica!")
			return
		}

	}
}

//AgregarOQuitarLogicas este metodo añade o quita espacio a las lógicas
func AgregarOQuitarLogicas(path string, add int64, name string, unidades int64, nombreParticion string, tamExtendida int64, comienzo int64) {

	if nombreParticion == name {
		if add < 0 {
			if ebr.PartSize > add {
				ebr.PartSize += add
				EscribirEBR(ebr.PartStart, path)
				mensajeCreado(path, nombreParticion, ebr.PartSize-add, ebr.PartSize)
				return
			}
			fmt.Println(colorYellow, "************************INFORMACIÓN************************")
			fmt.Println("La partición no es lo suficientemente grande para reducirla")
			return
		}
		tamEBR := int64(unsafe.Sizeof(ebr)) - 1
		if ebr.PartNext != -1 {
			libre := ebr.PartNext - (ebr.PartSize + ebr.PartStart + tamEBR)
			if libre >= add {
				ebr.PartSize += add
				EscribirEBR(ebr.PartStart, path)
				mensajeCreado(path, nombreParticion, ebr.PartSize-add, ebr.PartSize)
				return
			}
			fmt.Println(colorYellow, "*********************************INFORMACIÓN*******************************")
			fmt.Println("No hay suficiente espacio libre para aumentar la partición con el tamaño solicitado")
			return
		}
		libre := (tamExtendida + comienzo) - (ebr.PartSize + ebr.PartStart + tamEBR)
		if libre >= add {
			ebr.PartSize += add
			EscribirEBR(ebr.PartStart, path)
			mensajeCreado(path, nombreParticion, ebr.PartSize-add, ebr.PartSize)
			return
		}
		fmt.Println(colorYellow, "*********************************INFORMACIÓN*******************************")
		fmt.Println("No hay suficiente espacio libre para aumentar la partición con el tamaño solicitado")
		return
	}

	if ebr.PartNext != -1 {
		ebr = ExtraerEBR(path, ebr.PartNext)
		nombreParticion = Nombres(ebr.PartName)
		AgregarOQuitarLogicas(path, add, name, unidades, nombreParticion, tamExtendida, comienzo)
		return
	}
	fmt.Println("**********************MENSAJE**************************")
	fmt.Println("No se ha encontrado el nombre de la partición a reducir")

}

func mensajeCreado(path string, nombreParticion string, antes int64, despues int64) {
	EscribirMBR(path)
	fmt.Println(colorYellow, "******************************************************")
	fmt.Println(colorYellow, " Se ha implementado el comando add a la partición")
	fmt.Println(colorYellow, "******************************************************")
	fmt.Println(" Nombre de la partición: " + nombreParticion)
	fmt.Printf("%s%d%s", " Tamaño anterior de la partición: ", antes, " bytes\n")
	fmt.Printf("%s%d%s", " Tamaño actual de la partición: ", despues, " bytes\n")
	fmt.Println(colorYellow, "******************************************************")
}

//EliminarParticion este metodo realiza la eliminación de una partición
func EliminarParticion(path string, name string, tipo string) {
	b := false
	mbr, b = LeerMBR(path)
	if b == true {
		for i := 0; i < len(mbr.Particiones); i++ {
			nombreParticion := Nombres(mbr.Particiones[i].PartName)
			// Verifica si está en la partición
			if name == nombreParticion {
				var nuevoNombre [16]byte
				tt := mbr.Particiones[i].PartType
				mbr.Particiones[i].PartName = nuevoNombre
				mbr.Particiones[i].PartStatus = 73
				mbr.Particiones[i].PartType = 0
				mbr.Particiones[i].PartFit = 0
				mbr.Particiones[i].PartPartition = false
				mbr.Particiones[i].PartDelete = true
				tamm := mbr.Particiones[i].PartSize
				st := mbr.Particiones[i].PartStart
				mbr.Particiones[i].PartSize = 0
				mbr.Particiones[i].PartStart = 0
				mbr.MbrActivas--
				if strings.ToLower(tipo) == "fast" {
					EscribirMBR(path)
					mensajeEliminar(tamm, name, "Parcial", string(rune(tt)))
				} else {
					EscribirMBR(path)
					EliminacionFULLP(st, path, tamm)
					mensajeEliminar(tamm, name, "Total", string(rune(tt)))
				}
				return
			}

		}
		for i := 0; i < len(mbr.Particiones); i++ {
			if mbr.Particiones[i].PartType == byte('e') {
				ebr = ExtraerEBR(path, mbr.Particiones[i].PartStart)
				BuscarEliminarLogica(name, path, tipo, mbr.Particiones[i].PartStart)
				return
			}
		}
	} else {
		return
	}

	fmt.Println(colorYellow, "No existe el nombre de la partición, imposible eliminarla.")
}
func mensajeEliminar(ss int64, name string, tipo string, tipo2 string) {
	fmt.Println(colorRed, "***Información de partición eliminada***")
	fmt.Println(" Nombre de la partición: " + name)
	fmt.Printf("%s%d%s", " Tamaño de la partición: ", ss, "\n")
	fmt.Println(" Tipo de partición: " + tipo2)
	fmt.Println(" Tipo de eliminación: " + tipo)
	fmt.Println(colorRed, "****************************************")
}

//BuscarEliminarLogica este metodo busca la partición que se desea eliminar, si está la elimina
func BuscarEliminarLogica(name string, path string, tipo string, empieza int64) {

	nombre := Nombres(ebr.PartName)

	if nombre == name {
		if ebr.PartStart != empieza {
			ebrAnterior := EBR{}
			ebrAnterior = ExtraerEBR(path, ebr.PartPrevious)
			ebrAnterior.PartNext = ebr.PartNext
			actual := EBR{}
			actual = ebr
			ebr = ebrAnterior
			EscribirEBR(ebr.PartStart, path)
			siguiente := EBR{}
			siguiente = ExtraerEBR(path, actual.PartNext)
			siguiente.PartPrevious = ebr.PartStart
			ebr = siguiente
			EscribirEBR(ebr.PartStart, path)
			ebr = actual
			ebr.PartStatus = 73
			EscribirEBR(ebr.PartStart, path)
			mensajeEliminar(actual.PartSize, nombre, tipo, "Logica")
			if strings.ToLower(tipo) == "full" {
				EliminacionFULLP(ebr.PartStart, path, ebr.PartSize)
			}
			return
		}
		mensajeEliminar(ebr.PartSize, nombre, tipo, "Logica")
		ebr.PartStatus = 73
		ebr.PartFit = 0
		ebr.PartName = [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		EscribirEBR(ebr.PartStart, path)
		return

	}
	if ebr.PartNext != -1 {
		ebr = ExtraerEBR(path, ebr.PartNext)
		BuscarEliminarLogica(name, path, tipo, empieza)
		return
	}

	fmt.Println(colorYellow, "************************Mensaje**************************")
	fmt.Println(colorYellow, "No existe el nombre de la partición, imposible eliminarla.")
	fmt.Println(colorYellow, "**********************************************************")

}

//EscribirMBR modifica la información del mbr
func EscribirMBR(path string) {
	files, err := os.OpenFile(path, os.O_RDWR, 0644)
	defer files.Close()
	if err != nil {
		log.Fatal(err)
	}
	files.Seek(0, 0)
	var b3 bytes.Buffer
	binary.Write(&b3, binary.BigEndian, &mbr)
	escribirBytes(files, b3.Bytes())
}

//EliminacionFULLP hace una eliminacion completa de la particion
func EliminacionFULLP(start int64, path string, size int64) {
	files, err := os.OpenFile(path, os.O_RDWR, 0644)
	defer files.Close()
	if err != nil {
		log.Fatal(err)
	}
	var FULL uint8 = 0
	files.Seek(start-1, 0)
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, &FULL)
	escribirBytes(files, binario.Bytes())

	files.Seek(size, 1)
	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, &FULL)
	escribirBytes(files, binario2.Bytes())
}

//ExisteNombreParticion busca en el mbr si hay una particion con el mismo nombre
func ExisteNombreParticion(nom string, path string) bool {
	for i := 0; i < 4; i++ {
		nombreAnalizar := Nombres(mbr.Particiones[i].PartName)
		if nom == nombreAnalizar {
			return true
		}
		if mbr.Particiones[i].PartType == byte('e') {
			ebr = ExtraerEBR(path, mbr.Particiones[i].PartStart)
			nombreAnalizar := Nombres(ebr.PartName)
			if nombreAnalizar == nom {
				return true
			}
			for ebr.PartNext != -1 {
				ebr = ExtraerEBR(path, ebr.PartNext)
				nombreAnalizar := Nombres(ebr.PartName)
				if nom == nombreAnalizar {
					return true
				}
			}
		}
	}
	return false
}

//VerificarExistenciaExtendida este metodo verifica que no haya más de una particion extendida
func VerificarExistenciaExtendida() bool {
	for i := 0; i < 4; i++ {
		if mbr.Particiones[i].PartType == 101 {
			return true
		}
	}
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
	size = size * unidad
	var s int64
	var part int
	b := false
	mbr, b = LeerMBR(path)
	if b == true {
		if strings.ToLower(tipo) == "e" && VerificarExistenciaExtendida() {
			fmt.Println(colorYellow, "Ya existe una partición extendida")
			return
		}
		if ExisteNombreParticion(name, path) {
			fmt.Println(colorYellow, "El nombre de la partición ya existe!")
			return
		}
		if strings.ToLower(tipo) == "l" {
			st, ss := BuscarExtendida()
			if st != -1 {
				ebr = ExtraerEBR(path, st)
				nuevofit := ' '
				if strings.ToLower(fit) == "bf" {
					nuevofit = 'b'
				} else if strings.ToLower(fit) == "ff" {
					nuevofit = 'f'
				} else if strings.ToLower(fit) == "wf" {
					nuevofit = 'w'
				}
				CrearLogica(path, size, name, byte(nuevofit), ss)
				return
			}
			fmt.Println(colorYellow, "No se puede crear una partición lógica si no existe una partición extendida.")
			return

		}
		s, part = PrimerAjuste(size)
		if s != 0 && part != -1 {

			InformacionParticion(name, tipo, fit, size, s, part, path)
			CrearParticion(path, name, tipo, part)

		}

	}
}

//CrearEBR crea el ebr y lo situa en el archivo
func CrearEBR(start int64, size int64, previous int64) {
	ebr = EBR{PartStatus: 73, PartStart: start}
	ebr.PartSize = size
	ebr.PartNext = -1
	ebr.PartPrevious = previous
}

//EscribirEBR escribe los EBR que se van formando en la partición extendida
func EscribirEBR(start int64, path string) {
	files, err := os.OpenFile(path, os.O_RDWR, 0644)
	defer files.Close()
	if err != nil {
		log.Fatal(err)
	}
	files.Seek(start, 0)
	tamEBR := int64(unsafe.Sizeof(ebr))
	files.Seek(tamEBR, 1)
	var b3 bytes.Buffer
	binary.Write(&b3, binary.BigEndian, &ebr)
	escribirBytes(files, b3.Bytes())
	//ExtraerEBR(path, start)
}

//InformacionParticion en este metodo se agrega toda la información al struct particion
func InformacionParticion(name string, tipo string, fit string, size int64, start int64, numero int, path string) {
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
	} else if strings.ToLower(tipo) == "l" {
		nuevoTipo = 'l'
	}
	if nuevoTipo == 'e' {
		//tamebr := int64(unsafe.Sizeof(ebr)) - 1
		CrearEBR(start, size, -1)
		EscribirEBR(mbr.Particiones[numero].PartStart, path)
	}
	copy(mbr.Particiones[numero].PartName[:], name)
	mbr.Particiones[numero].PartSize = size
	mbr.Particiones[numero].PartFit = byte(nuevofit)
	mbr.Particiones[numero].PartType = byte(nuevoTipo)
	mbr.Particiones[numero].PartStatus = 65
	mbr.Particiones[numero].PartPartition = true
}

//CrearParticion crea la particion
func CrearParticion(path string, name string, tipo string, numero int) {
	files, err := os.OpenFile(path, os.O_RDWR, 0644)
	defer files.Close()
	if err != nil {
		log.Fatal(err)
	}
	files.Seek(0, 0)
	var b3 bytes.Buffer
	binary.Write(&b3, binary.BigEndian, &mbr)
	escribirBytes(files, b3.Bytes())

	ParticionesActivas(path)
}

//ParticionesActivas muestra todas las particiones en el disco
func ParticionesActivas(path string) {
	fmt.Println(colorGreen, "**********Particiones actuales**********")
	OrdenarArregloParticion()
	for i := 0; i < 4; i++ {
		if mbr.Particiones[i].PartStatus == 65 {
			fmt.Println(colorGreen, "Nombre de la partición: "+Nombres(mbr.Particiones[i].PartName))
			fmt.Println(colorGreen, "Tamaño de la partición: "+strconv.Itoa(int(mbr.Particiones[i].PartSize)))
			fmt.Println(colorGreen, "Tipo de ajuste: "+string(rune(mbr.Particiones[i].PartFit)))
			var tipos byte = mbr.Particiones[i].PartType
			switch string(rune(tipos)) {
			case "p":
				fmt.Println(colorGreen, "Tipo de partición: Primaria")
				break
			case "e":
				fmt.Println(colorGreen, "Tipo de partición: Extendida")
				ebr = ExtraerEBR(path, mbr.Particiones[i].PartStart)
				if ebr.PartStatus == 65 || ebr.PartNext != -1 {
					fmt.Println("__________Particiones lógicas__________")
				}
				if ebr.PartStatus == 65 {
					fmt.Println(colorGreen, "Nombre de la partición: "+Nombres(ebr.PartName))
					fmt.Println("Tamaño de partición: " + strconv.Itoa(int(ebr.PartSize)))
					fmt.Println("Ajuste de partición: " + string(rune(ebr.PartFit)))
				}
				for ebr.PartNext != -1 {
					ebr = ExtraerEBR(path, ebr.PartNext)
					fmt.Println(colorGreen, "Nombre de la partición: "+Nombres(ebr.PartName))
					fmt.Println("Tamaño de partición: " + strconv.Itoa(int(ebr.PartSize)))
					fmt.Println("Ajuste de partición: " + string(rune(ebr.PartFit)))
				}
				fmt.Println("_______________________________________")
				break

			}
			fmt.Println("****************************************")
		}
	}

}

//PrimerAjuste este metodo devuelve la posicion inicial del primer espacio que encuentre
func PrimerAjuste(tam int64) (int64, int) {

	if mbr.MbrActivas == 4 {
		fmt.Println(colorYellow, "Ya existen 4 particiones en el disco actual")
		return 0, -1
	}
	for i := 0; i < 4; i++ {
		TAM := mbr.Particiones[i].PartSize
		if !mbr.Particiones[i].PartPartition {
			if TAM >= tam {
				mbr.Particiones[i].PartSize = tam
				mbr.MbrRecorrido += mbr.Particiones[i].PartSize

				if i < 3 && mbr.Particiones[i].PartDelete == false {
					mbr.Particiones[i+1].PartStart = mbr.Particiones[i].PartStart + tam
					mbr.Particiones[i+1].PartSize = mbr.MbrTam - int64(unsafe.Sizeof(mbr)) - mbr.MbrRecorrido
				}
				mbr.MbrActivas++
				return mbr.Particiones[i].PartStart, i
			}
		}
	}
	OrdenarArregloParticion()
	return Ajustar(tam, int64(unsafe.Sizeof(mbr)), 0)
}

//Ajustar pruebas
func Ajustar(tam int64, Inicio int64, i int) (int64, int) {

	if mbr.Particiones[i].PartStart != 0 {
		if mbr.Particiones[i].PartStart > Inicio {
			libre := mbr.Particiones[i].PartStart - Inicio
			if libre >= tam {
				for a := 0; a < 4; a++ {
					if mbr.Particiones[a].PartStart == 0 {
						mbr.Particiones[a].PartStart = Inicio
						mbr.Particiones[a].PartSize = tam
						mbr.MbrActivas++
						return mbr.Particiones[a].PartStart, a
					}
				}
			}
			Inicio = mbr.Particiones[i].PartStart + mbr.Particiones[i].PartSize
			if (i + 1) != 3 {
				return Ajustar(tam, Inicio, i+1)
			}
			libre = mbr.MbrTam - (mbr.Particiones[i+1].PartStart + mbr.Particiones[i+1].PartSize)
			Inicio = mbr.Particiones[i+1].PartStart + mbr.Particiones[i+1].PartSize
			if libre >= tam {
				for a := 0; a < 4; a++ {
					if mbr.Particiones[a].PartStart == 0 {
						mbr.Particiones[a].PartStart = Inicio
						mbr.Particiones[a].PartSize = tam
						mbr.MbrActivas++
						return mbr.Particiones[a].PartStart, a
					}
				}
			}

		}
	} else {
		return Ajustar(tam, Inicio, i+1)
	}
	fmt.Println(colorYellow, "Ya existen 4 particiones en el disco actual")
	return 0, -1
}

//OrdenarArregloParticion ordena el arreglo de particiones
func OrdenarArregloParticion() {
	for i := 0; i < len(mbr.Particiones); i++ {
		for j := 0; j < len(mbr.Particiones)-1; j++ {
			if mbr.Particiones[j].PartStart > mbr.Particiones[j+1].PartStart {
				temp := mbr.Particiones[j]
				mbr.Particiones[j] = mbr.Particiones[j+1]
				mbr.Particiones[j+1] = temp
			}
		}
	}
}

//BuscarExtendida este metodo busca la partición extendida para extraer su ebr
func BuscarExtendida() (int64, int64) {
	for i := 0; i < 4; i++ {
		if mbr.Particiones[i].PartType == 101 {
			return mbr.Particiones[i].PartStart, mbr.Particiones[i].PartSize
		}
	}
	return -1, 0
}

//ExtraerEBR este método extrae el struct del primer ebr de la partición extendida
func ExtraerEBR(path string, start int64) EBR {
	files, err := os.OpenFile(path, os.O_RDWR, 0644)
	defer files.Close()
	if err != nil {
		log.Fatal(err)
	}
	files.Seek(0, 0)
	files.Seek(start, 0)
	ebr2 := EBR{}
	var size int = int(unsafe.Sizeof(ebr2))

	files.Seek(int64(size), 1)

	data := readNextBytes(files, size)
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &ebr2)
	if err != nil {
		panic(err)
	}
	return ebr2
}

//Libres dd
func Libres(path string, start int64) EBR {
	ebrC := EBR{}
	ebrC = ExtraerEBR(path, start)
	if ebrC.PartStatus == 65 {
		return ebrC
	}
	if ebrC.PartNext != -1 {
		return Libres(path, ebrC.PartNext)
	}
	return ebrC
}

//InicioVacio verifica el primer ebr
func InicioVacio() bool {
	if ebr.PartStatus == 73 {
		return true
	}
	return false
}

//CrearLogica Verifica si se puede crear una logica, si hay espacio la crea
func CrearLogica(path string, size int64, name string, fit byte, tamExtendida int64) {

	inicioExtendida := ebr.PartStart
	if InicioVacio() {
		tami := int64(unsafe.Sizeof(ebr)) - 1
		if ebr.PartNext == -1 && (ebr.PartSize-tami) >= size {
			ebr.PartSize = size
			ebr.PartFit = fit
			ebr.PartStatus = 65
			copy(ebr.PartName[:], name)
			EscribirEBR(ebr.PartStart, path)
			MensajeConfirmacion()
			return
		}
		auxi := EBR{}
		auxi = ebr
		auxi = ExtraerEBR(path, ebr.PartNext)
		sobra := auxi.PartStart - ebr.PartStart
		if sobra >= size {
			ebr.PartSize = size
			ebr.PartFit = fit
			ebr.PartStatus = 65
			copy(ebr.PartName[:], name)
			ebr.PartNext = auxi.PartStart
			auxi.PartPrevious = ebr.PartStart
			EscribirEBR(ebr.PartStart, path)
			ebr = auxi
			EscribirEBR(ebr.PartStart, path)
			MensajeConfirmacion()
			return
		}

	}
	TamEBR := int64(unsafe.Sizeof(ebr)) - 1
	for ebr.PartNext != -1 {
		sobra := ebr.PartNext - (ebr.PartStart + ebr.PartSize + TamEBR)
		if (sobra) >= (size + TamEBR) {
			inicio := ebr.PartStart + ebr.PartSize + TamEBR
			nuevo := EBR{}
			nuevo.PartNext = ebr.PartNext
			actual := EBR{}
			actual = ebr
			ebr = ExtraerEBR(path, actual.PartNext)
			ebr.PartPrevious = inicio
			EscribirEBR(ebr.PartStart, path)
			ebr = actual
			nuevo.PartPrevious = ebr.PartStart
			nuevo.PartStart = inicio
			nuevo.PartFit = fit
			nuevo.PartSize = size
			copy(nuevo.PartName[:], name)
			nuevo.PartStatus = 65
			ebr.PartNext = nuevo.PartStart
			EscribirEBR(ebr.PartStart, path)
			ebr = nuevo
			EscribirEBR(ebr.PartStart, path)
			MensajeConfirmacion()
			return
		}
		ebr = ExtraerEBR(path, ebr.PartNext)
	}

	libre := (inicioExtendida + tamExtendida) - (ebr.PartSize + ebr.PartStart + TamEBR)

	if libre >= (size + TamEBR) {
		siguiente := ebr.PartSize + ebr.PartStart + TamEBR
		anterior := ebr.PartStart
		ebr.PartNext = siguiente
		EscribirEBR(ebr.PartStart, path)
		ebr = EBR{}
		CrearEBR(siguiente, size, anterior)
		ebr.PartFit = fit
		copy(ebr.PartName[:], name)
		ebr.PartStatus = 65
		ebr.PartNext = -1
		EscribirEBR(ebr.PartStart, path)
		MensajeConfirmacion()
		return
	}

	fmt.Println("**********************************************")
	fmt.Println("*No hay más espacio en la partición extendida*")
	fmt.Println("**********************************************")

}

//MensajeConfirmacion este metodo imprime un mensaje
func MensajeConfirmacion() {

	fmt.Println(colorGreen, "****Información de la partición lógica****")
	fmt.Println(" Nombre de la partición: " + Nombres(ebr.PartName))
	fmt.Printf("%s%d", "Tamaño de la partición: ", ebr.PartSize)
	fmt.Println("\n*******************************************")

}

//LeerMBR este metodo devuelve el mbr actual del disco
func LeerMBR(path string) (MBR, bool) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		fmt.Println(colorRed, "No se encontró la ruta del archivo")
		return mbr, false
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

	return mbr2, true
}
func graficarMBR(path string, ubicacion string) {

	dir := ""
	rutas := strings.Split(ubicacion, "/")
	for i := 0; i < len(rutas)-1; i++ {
		dir += rutas[i] + "/"
	}
	nombre := rutas[len(rutas)-1]
	extension := strings.Split(nombre, ".")

	if AnalizarRuta(dir) {

		dd := ubicacion
		dirdoc := dir + extension[0] + ".txt"
		var _, errr = os.Stat(dirdoc)
		//Crea el archivo si no existe
		if os.IsNotExist(errr) {
			var file, errr = os.Create(dirdoc)
			if existeError(errr) {
				return
			}
			defer file.Close()
		}

		cadena := ""
		cadena += "digraph G {\ngraph [pad=\"0.5\", nodesep=\"1\", ranksep=\"2\"];"
		cadena += "\nnode [shape=plain]\nrankdir=LR;\n"
		cadena += "Tabla[label=<\n<table border=\"0\" cellborder=\"1\" cellspacing=\"0\">\n"
		cadena += "<tr><td><i>Nombre</i></td>\n<td><i>Valor</i> </td>\n</tr>"

		LeerMBR(path)

		cadena += "<tr><td>Mbr_sizeDisk</td><td>" + strconv.Itoa(int(mbr.MbrTam)) + "</td></tr>\n"
		cadena += "<tr><td>Mbr_FechaCreacion</td><td>" + string(mbr.MbrFechaCreacion[:]) + "</td></tr>\n"
		cadena += "<tr><td>Mbr_DiskSignature</td><td>" + strconv.Itoa(int(mbr.MbrDiskID)) + "</td></tr>\n"

		for i := 0; i < len(mbr.Particiones); i++ {
			nombre := ""
			for j := 0; j < len(mbr.Particiones[i].PartName); j++ {
				if mbr.Particiones[i].PartName[j] != 0 {
					nombre += string(rune(mbr.Particiones[i].PartName[j]))
				} else {
					break
				}
			}
			if nombre == "" {
				nombre = "---"
			}
			cadena += "<tr><td>Part" + strconv.Itoa((i + 1)) + "_Name</td><td>" + nombre + "</td></tr>\n"
			cadena += "<tr><td>Part" + strconv.Itoa((i + 1)) + "_Size</td><td>" + strconv.Itoa(int(mbr.Particiones[i].PartSize)) + "</td></tr>\n"
			cadena += "<tr><td>Part" + strconv.Itoa((i + 1)) + "_Start</td><td>" + strconv.Itoa(int(mbr.Particiones[i].PartStart)) + "</td></tr>\n"
			cadena += "<tr><td>Part" + strconv.Itoa((i + 1)) + "_Status</td><td>" + string(rune(mbr.Particiones[i].PartStatus)) + "</td></tr>\n"
			cadena += "<tr><td>Part" + strconv.Itoa((i + 1)) + "_Fit</td><td>" + string(rune(mbr.Particiones[i].PartFit)) + "</td></tr>\n"
			cadena += "<tr><td>Part" + strconv.Itoa((i + 1)) + "_Type</td><td>" + string(rune(mbr.Particiones[i].PartType)) + "</td></tr>\n"
		}

		cadena += "</table>>];}"
		errrr := ioutil.WriteFile(dirdoc, []byte(cadena[:]), 0644)
		if errrr != nil {
			panic(errrr)
		}
		com1 := "dot"
		com2 := "-T" + strings.ToLower(extension[1])
		com3 := dirdoc
		com4 := "-o"
		com5 := dd
		exec.Command(com1, com2, com3, com4, com5).Output()
		fmt.Println(colorGreen, "Success")
	}
}
func existeError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}
	return (err != nil)
}

//Nombres este metodo devuelve el nombre en string
func Nombres(n [16]byte) string {
	nombre := ""
	for j := 0; j < len(n); j++ {
		if n[j] != 0 {
			nombre += string(rune(n[j]))
		} else {
			break
		}
	}
	return nombre
}

//GraficarDisco crea el txt del graphviz para graficar
func GraficarDisco(path string, ubicacion string) {

	dir := ""
	rutas := strings.Split(ubicacion, "/")
	for i := 0; i < len(rutas)-1; i++ {
		dir += rutas[i] + "/"
	}
	nombre := rutas[len(rutas)-1]
	extension := strings.Split(nombre, ".")

	if AnalizarRuta(dir) {
		dd := dir
		dir = dir + extension[0] + ".txt"
		var _, errr = os.Stat(dir)
		//Crea el archivo si no existe
		if os.IsNotExist(errr) {
			var file, errr = os.Create(dir)
			if existeError(errr) {
				return
			}
			defer file.Close()
		}
		d := false
		mbr, d = LeerMBR(path)
		if d == true {
			cadena := "digraph structs {\n"
			cadena += "node [shape=record];\n"
			cadena += "disco[label=\"MBR&#92;nSize: " + strconv.Itoa(int(mbr.MbrTam))

			OrdenarArregloParticion()

			Inicio := int64(unsafe.Sizeof(mbr))

			for i := 0; i < 4; i++ {
				if mbr.Particiones[i].PartSize != 0 {
					if mbr.Particiones[i].PartStart > Inicio {
						cadena += "|"
						cadena += "Libre: "
						disponible := mbr.Particiones[i].PartStart - Inicio
						cadena += strconv.Itoa(int(disponible))
						Inicio = mbr.Particiones[i].PartStart + mbr.Particiones[i].PartSize
						i--
					} else {
						cadena += "|"
						nombre := Nombres(mbr.Particiones[i].PartName)
						if mbr.Particiones[i].PartType == byte('e') {
							cadena += Grafextendida(path, i, nombre)
						} else {
							cadena += "Nombre: " + nombre + "&#92;n"
							cadena += "Tipo: " + "Primaria" + "&#92;n"
							cadena += "Size: " + strconv.Itoa(int(mbr.Particiones[i].PartSize))
						}

						Inicio = mbr.Particiones[i].PartStart + mbr.Particiones[i].PartSize
					}
				}
			}
			if mbr.Particiones[3].PartStart != 0 {
				libre := mbr.MbrTam - mbr.Particiones[3].PartStart - mbr.Particiones[3].PartSize
				cadena += "|"
				cadena += "Libre: "
				cadena += strconv.Itoa(int(libre))
			}

			cadena += "\"];}"
			errrr := ioutil.WriteFile(dir, []byte(cadena[:]), 0644)
			if errrr != nil {
				panic(errrr)
			}
			com1 := "dot"
			com2 := "-T" + strings.ToLower(extension[1])
			com3 := dir
			com4 := "-o"
			com5 := dd + nombre
			exec.Command(com1, com2, com3, com4, com5).Output()
			fmt.Println(colorGreen, "Success")
		}
	}
}

//Grafextendida devuelve una cadena con el codigo para graficar particiones logicas
func Grafextendida(path string, i int, nombre string) string {
	cadena := ""
	cadena += "{"
	cadena += "Nombre: " + nombre + "&#92;n"
	cadena += "Tipo: " + "Extendida&#92;n"
	cadena += "Size: " + strconv.Itoa(int(mbr.Particiones[i].PartSize)) + " bytes|{"
	com, ss := BuscarExtendida()

	ebr = ExtraerEBR(path, com)
	a := false
	if ebr.PartStatus == 73 {
		if ebr.PartNext != -1 {
			libre := ebr.PartNext - (ebr.PartStart)
			cadena += "Libre: " + strconv.Itoa(int(libre))
			a = true
		} else {
			tamebr := int64(unsafe.Sizeof(ebr)) - 1
			cadena += "Libre: " + strconv.Itoa(int(ss-tamebr))
		}
	}
	if a == true || (ebr.PartStatus == 65) {
		if a == false {
			cadena += "EBR&#92;n"
			nombre := Nombres(ebr.PartName)
			cadena += "Nombre: " + nombre + "&#92;n"
			cadena += "|"
			cadena += "Logica&#92;n"
			cadena += "Size: " + strconv.Itoa(int(ebr.PartSize))
		} else {
			ebr = ExtraerEBR(path, ebr.PartNext)
			if ebr.PartStatus == 65 {
				cadena += "|"
				cadena += "EBR&#92;n"
				nombre := Nombres(ebr.PartName)
				cadena += "Nombre: " + nombre + "&#92;n"
				cadena += "|"
				cadena += "Logica&#92;n"
				cadena += "Size: " + strconv.Itoa(int(ebr.PartSize))
			}
		}
		tamebr := int64(unsafe.Sizeof(ebr)) - 1
		for ebr.PartNext != -1 {
			libre := ebr.PartNext - (tamebr + ebr.PartStart + ebr.PartSize)

			if libre == 0 {

				cadena += "|"

				if ebr.PartNext != -1 {
					auxi := EBR{}
					auxi = ExtraerEBR(path, ebr.PartNext)
					nombre := Nombres(auxi.PartName)
					cadena += "EBR&#92;n"
					cadena += "Nombre: " + nombre + "&#92;n"
					cadena += "|"
					cadena += "Logica&#92;n"
					cadena += "Size: " + strconv.Itoa(int(auxi.PartSize))

				}
			} else {
				cadena += "|"
				cadena += "Libre: " + strconv.Itoa(int(libre))
				if ebr.PartStatus == 65 {
					cadena += "|"
					auxi := EBR{}
					auxi = ExtraerEBR(path, ebr.PartNext)
					nombre := Nombres(auxi.PartName)
					cadena += "EBR&#92;n"
					cadena += "Nombre: " + nombre + "&#92;n"
					cadena += "|"
					cadena += "Logica&#92;n"
					cadena += "Size: " + strconv.Itoa(int(auxi.PartSize))

				}
			}
			ebr = ExtraerEBR(path, ebr.PartNext)
		}
		libre := (ss + com) - (ebr.PartStart + tamebr + ebr.PartSize)
		if libre != 0 {
			cadena += "|"
			cadena += "Libre: " + strconv.Itoa(int(libre))
		}
	}
	cadena += "}}"
	return cadena
}
func imprimirComienzo(path string) {
	for ebr.PartNext != -1 {
		fmt.Println("partición: " + Nombres(ebr.PartName) + ", inicia: " + strconv.Itoa(int(ebr.PartStart)))
		fmt.Println("siguiente: " + strconv.Itoa(int(ebr.PartNext)))
		fmt.Println("anterior: " + strconv.Itoa(int(ebr.PartPrevious)))
		ebr = ExtraerEBR(path, ebr.PartNext)
	}
}

//ExisteArchivo verifica si el disco existe
func ExisteArchivo(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

//Mount monta una partición en la RAM
func Mount(arreglo []string) {
	dir := ""
	aux := 0
	nombreParticion := ""
	for i := 1; i < len(arreglo); i++ {
		subcomando := strings.Split(arreglo[i], "->")
		switch strings.ToLower(subcomando[0]) {
		case "-path":
			dir = direccion(arreglo[i])
			if ExisteArchivo(dir) {
				aux++
			} else {
				fmt.Println("La ruta no ha sido encontrada!")
				return
			}
			break
		case "-name":
			if verificarNombreParticion(subcomando[1]) {
				aux++
				nombreParticion = subcomando[1]
			} else {
				return
			}
			break
		case "":
			break
		default:
			fmt.Println(colorRed, "Comando no reconocido!")
			return
		}
	}
	if aux == 2 {
		montar(nombreParticion, dir)
	} else {
		fmt.Println(colorRed, "Faltan parámetros solicitados!")
	}
}

func montar(nombreParticion string, path string) {
	b := false
	mbr, b = LeerMBR(path)
	if b {
		AgregarDisco(path, nombreParticion)
	}
}

//NodoParticion este struct contiene los datos que va a tener la lista de particiones montadas
type NodoParticion struct {
	name          [16]byte
	nombreMontada string
	numero        int32
	ebr           EBR
	part          particion
	tipo          string
	siguiente     *NodoParticion
	fecha         [19]byte
}

//NodoDisco el nodo contendrá la lista de discos
type NodoDisco struct {
	path             string
	Nombre           string
	Letra            byte
	listaParticiones ListaParticion
	siguiente        *NodoDisco
}

//ListaDisco este struct guarda los atributos de la lista disco
type ListaDisco struct {
	inicio *NodoDisco
}

//ListaParticion este struct guarda los atributos de la lista
type ListaParticion struct {
	inicio *NodoParticion
}

//ListaDiscoVacia devuelve verdadero si la lista está vacía
func ListaDiscoVacia() bool {
	if ListDiscos.inicio == nil {
		return true
	}
	return false
}

//AgregarDisco este metodo mete el disco a la lista
func AgregarDisco(path string, nombreParticion string) {
	p := false
	mbr, p = LeerMBR(path)
	if p {
		log := false
		PE, i := BuscarParticionPE(nombreParticion, path)
		if !PE && i == -1 {
			log = BuscarParticionExtendida(nombreParticion, path)
		}
		if ListaDiscoVacia() && (PE || log) && ExisteArchivo(path) {
			var ini NodoDisco = NodoDisco{}
			ListDiscos.inicio = &ini
			ListDiscos.inicio.Letra = 97
			array := strings.Split(path, "/")
			nombre := array[len(array)-1]
			ListDiscos.inicio.Nombre = nombre
			ListDiscos.inicio.path = path
			//Llenar lista de particion
			var listParticion ListaParticion

			var ini2 NodoParticion = NodoParticion{}
			listParticion.inicio = &ini2
			listParticion.inicio.numero = 1
			if PE {
				listParticion.inicio.part = mbr.Particiones[i]
				if mbr.Particiones[i].PartType == 'p' {
					listParticion.inicio.tipo = "primaria"
				} else {
					listParticion.inicio.tipo = "extendida"
				}
			} else if log {
				listParticion.inicio.ebr = ebr
				listParticion.inicio.tipo = "logica"
			}
			copy(listParticion.inicio.name[:], nombreParticion)
			t := time.Now()
			fecha := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", t.Year(), t.Month(), t.Day(),
				t.Hour(), t.Minute(), t.Second())
			copy(listParticion.inicio.fecha[:], fecha)
			listParticion.inicio.nombreMontada = "vda1"
			ListDiscos.inicio.listaParticiones = listParticion
			ListDiscos.inicio.siguiente = nil
			fmt.Println(colorGreen, "**************Información**************")
			fmt.Println(colorGreen, "Se ha montado la partición exitosamente")
		} else if (PE || log) && ExisteArchivo(path) {
			array := strings.Split(path, "/")
			nombre := array[len(array)-1]

			var auxiliar *NodoDisco
			auxiliar = ListDiscos.inicio
			a1 := false
			for auxiliar != nil {
				if auxiliar.Nombre == nombre {
					PosListaParticion(auxiliar.listaParticiones, nombreParticion, string(rune(auxiliar.Letra)), PE, log, i)
					a1 = true
					break
				}
			}

			if !a1 {
				var auxiliar2 *NodoDisco
				auxiliar2 = ListDiscos.inicio
				for auxiliar2.siguiente != nil {
					auxiliar2 = auxiliar2.siguiente
				}
				auxiliar2.siguiente.Letra = auxiliar2.Letra + 1
				array := strings.Split(path, "/")
				nombre := array[len(array)-1]
				auxiliar2.siguiente.Nombre = nombre
				auxiliar2.siguiente.path = path
				//Llenar lista de particion
				var listParticion ListaParticion
				listParticion.inicio.numero = 1
				if PE {
					listParticion.inicio.part = mbr.Particiones[i]
					if mbr.Particiones[i].PartType == 'p' {
						listParticion.inicio.tipo = "primaria"
					}
				} else if log {
					listParticion.inicio.ebr = ebr
					listParticion.inicio.tipo = "logica"
				}
				t := time.Now()
				fecha := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", t.Year(), t.Month(), t.Day(),
					t.Hour(), t.Minute(), t.Second())
				copy(listParticion.inicio.fecha[:], fecha)
				copy(listParticion.inicio.name[:], nombreParticion)
				listParticion.inicio.nombreMontada = "vd" + string(rune(auxiliar2.siguiente.Letra+1)) + "1"
				ListDiscos.inicio.listaParticiones = listParticion
				ListDiscos.inicio.siguiente = nil
				fmt.Println(colorGreen, "**************Información**************")
				fmt.Println(colorGreen, "Se ha montado la partición exitosamente")
			}
		} else {
			fmt.Println(colorRed, "Verifique el nombre y tipo de partición")
		}
	}
}

//BuscarParticionPE retorna informacion del tipo de particion
func BuscarParticionPE(nombre string, path string) (bool, int) {
	for i := 0; i < len(mbr.Particiones); i++ {
		nom := Nombres(mbr.Particiones[i].PartName)
		if nom == nombre {
			if mbr.Particiones[i].PartType != byte('e') {
				return true, i
			}
			return false, i

		}
	}
	return false, -1
}

//BuscarParticionExtendida si encuentra el nombre, devuelve un boleano
func BuscarParticionExtendida(nombre string, path string) bool {
	for i := 0; i < 4; i++ {
		if mbr.Particiones[i].PartType == byte('e') {
			ebr = ExtraerEBR(path, mbr.Particiones[i].PartStart)
			nom := Nombres(ebr.PartName)
			if nom == nombre {
				return true
			}
			for ebr.PartNext != -1 {
				ebr = ExtraerEBR(path, ebr.PartNext)
				nom := Nombres(ebr.PartName)
				if nom == nombre {
					return true
				}
			}
		}
	}
	return false
}

//PosListaParticion agrega un nuevo elemento a la lista particion
func PosListaParticion(Lista ListaParticion, nombreparticion string, letra string, PE bool, log bool, i int) {
	var auxiliar *NodoParticion
	auxiliar = Lista.inicio
	for auxiliar.siguiente != nil {
		auxiliar = auxiliar.siguiente
	}
	ini := NodoParticion{}
	auxiliar.siguiente = &ini

	if PE {
		auxiliar.siguiente.part = mbr.Particiones[i]
		if mbr.Particiones[i].PartType == 'p' {
			auxiliar.siguiente.tipo = "primaria"
		}
	} else if log {
		auxiliar.siguiente.ebr = ebr
		auxiliar.siguiente.tipo = "logica"
	}
	t := time.Now()
	fecha := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	copy(auxiliar.siguiente.fecha[:], fecha)
	copy(auxiliar.siguiente.name[:], nombreparticion)
	auxiliar.siguiente.numero = auxiliar.numero + 1
	auxiliar.siguiente.nombreMontada = "vd" + letra + strconv.Itoa(int(auxiliar.siguiente.numero))
	fmt.Println(colorGreen, "**************Información**************")
	fmt.Println(colorGreen, "Se ha montado la partición exitosamente")
}
func listarMontadas() {
	auxi := 1
	aa := ListDiscos.inicio
	if aa == nil {
		fmt.Println(colorYellow, "No hay ninguna partición montada")
		return
	}
	fmt.Println(colorGreen, "*****************Lista de particiones montadas*****************")
	for aa != nil {
		fmt.Println("Nombre del disco: " + aa.Nombre)
		lista := aa.listaParticiones.inicio
		for lista != nil {
			fmt.Println(colorGreen, strconv.Itoa(auxi)+".")
			fmt.Println("Nombre de la partición: " + Nombres(lista.name))
			fmt.Println("ID asignado: " + lista.nombreMontada)
			fmt.Println("Tipo de partición: " + lista.tipo)
			lista = lista.siguiente
			auxi++
		}
		aa = aa.siguiente
		fmt.Println(colorGreen, "___________________________________")
	}
}

//Unmount verificación de parametros para el comando unmount
func Unmount(arreglo []string) {

	for i := 1; i < len(arreglo); i++ {
		ids := strings.Split(arreglo[i], "->")
		numeroID := int64(0)
		esNum := false
		//Este for recorre todos los ids que encuentre
		if len(ids) == 1 {
			return
		}
		for j := 0; j < len(ids); j++ {
			ID := ids[j]
			montada := ""
			if (j + 1) < len(arreglo) {
				montada = ids[j+1]
				j++
			} else {
				return
			}
			numeroID, esNum = VerificarNumero(string(rune(ID[3])))
			if !esNum {
				fmt.Println("El id es incorrecto")
				return
			}
			if montada != "" {
				if strings.Contains(montada, "vd") {
					letra := montada[2]
					num, a := VerificarNumero(string(rune(montada[3])))
					if a {
						Desmontar(numeroID, letra, int32(num))
					}
				} else {
					fmt.Println("El nombre de la partición a desmontar no es correcto")
					return
				}
			} else {
				fmt.Println("Falta el id de la partición montada!")
				return
			}
		}
	}
}

//Desmontar esta función busca la partición a desmontar y si la encuentra pues la desmonta xd
func Desmontar(numID int64, letra byte, numP int32) {
	auxi := ListDiscos.inicio
	contador := int64(1)
	for auxi != nil {
		if contador == numID {
			if letra == auxi.Letra {
				buscarParticionMontadaDesmontar(numP, &auxi.listaParticiones, auxi.path)
			} else {
				fmt.Println(colorYellow, "La letra ingresada no coincide con la asignada")
				return
			}
		}
		contador++
		auxi = auxi.siguiente
	}
}

func buscarParticionMontadaDesmontar(numP int32, lista *ListaParticion, path string) {
	var auxi *NodoParticion
	auxi = lista.inicio

	for auxi.siguiente != nil {

		if auxi.numero == numP && auxi == lista.inicio {
			lista.inicio = lista.inicio.siguiente
			fmt.Println("Se ha desmontado la partición")
			return
		} else if auxi.siguiente.numero == numP {
			auxi.siguiente = auxi.siguiente.siguiente
			fmt.Println("Se ha desmontado la partición")
			return
		}

		auxi = auxi.siguiente
	}
}

var super SUPERBOOT

//SUPERBOOT atributos del superboot
type SUPERBOOT struct {
	SbNombreHD [16]byte
	//Cantidad de estructuras en la partición
	SbAVDcount               int64
	SbDetalleDirectorioCount int64
	SbINodoCount             int64
	SbBloquesCount           int64
	//Cantidad de estructuras libres
	SbAVDFree               int64
	SbDetalleDirectorioFree int64
	SbINodoFree             int64
	SbBloqueFree            int64
	//Fechas
	SbDateCreation [19]byte
	SbDateMontaje  [19]byte
	//Cantidad de montajes
	SbMontajesCount int64
	//Apuntador de inicio de bipmaps y estructuras
	SbAptrStartBipmapAVD    int64
	SbAptrStartAVD          int64
	SbAptrStartBipmapDD     int64
	SbAptrStartDD           int64
	SbAptrStartBipmapINodo  int64
	SbAptrStartINodo        int64
	SbAptrStartBipmapBloque int64
	SbAptrStartBloque       int64
	SbAptrStartLogBitacora  int64
	//Tamaño de una estructura
	SbSizeStructAVD    int64
	SbSizeStructDD     int64
	SbSizeStructINodo  int64
	SbSizeStructBloque int64
	//Primer bite libre en los bipmaps
	SbFirstBitFreeAVD    int64
	SbFirstBitFreeDD     int64
	SbFirstBitFreeINodo  int64
	SbFirstBitFreeBloque int64
	SbMagicNum           int64
}

//ArbolVirtualDirectorio información del AVD
type ArbolVirtualDirectorio struct {
	AVDFechaCreacion           [19]byte
	AVDNameDirectoy            [16]byte
	AVDAptrArraySubdirectorios [6]int64
	AVDAptrDetalleDirectorio   int64
	AVDAptrInd                 int64
	AVDProper                  int64
	AVDGid                     int64
	AVDPerm                    uint8
}

//DetalleDirectorio información del DD
type DetalleDirectorio struct {
	DDArrayAptrINodo [5]InfoArchivo
	DDAptrIndirecto  int64
}

//InfoArchivo info del archivo en el DD
type InfoArchivo struct {
	DDfileName             [20]byte
	DDAptrINodo            int64
	DDFileDateCreation     [19]byte
	DDFileDateModificacion [19]byte
}

//INodo info del archivo
type INodo struct {
	INodoNumero           int64
	INodoFileSize         int64
	INodoBloquesAsignados int64
	INodoAptrDeBloque     [4]int64
	INodoAptrInd          int64
	INodoIDPropier        int64
	INodoGid              int64
	INodoPerm             uint8
}

//Bloque contenido del bloque
type Bloque struct {
	BDarray [25]byte
}

//LogBitacora contenido del log
type LogBitacora struct {
	logTipoOperacion byte
	logTipo          byte
	logNombre        [20]byte
	//logContenido
	logFechaTransicion [19]byte
}

//MKFS evalua los subcomandos dle mkfs
func MKFS(arreglo []string) {
	id := ""
	tipo := ""
	aux := 0
	for i := 1; i < len(arreglo); i++ {
		comandos := strings.Split(arreglo[i], "->")
		switch strings.ToLower(comandos[0]) {
		case "-id":
			id = comandos[1]
			aux++
			break
		case "-type":
			if DELETE(comandos[1]) {
				aux++
			} else {
				fmt.Println("Tipo de formateo inválido")
				return
			}
			break
		case "":
			break
		default:
			fmt.Println(colorRed, "Comando no reconocido!")
			return

		}
	}
	if aux == 2 {
		BuscarID(tipo, id)
	} else if tipo == "" && id != "" {
		BuscarID("full", id)
	}

}

//BuscarID busca el id de la particion montada
func BuscarID(tipo string, id string) {
	var aux *NodoDisco
	aux = ListDiscos.inicio
	for aux != nil {
		var lista ListaParticion
		lista = aux.listaParticiones
		aux2 := lista.inicio
		for aux2 != nil {
			if aux2.nombreMontada == id {
				if aux2.tipo == "primaria" {
					formatear(aux2.part.PartSize, aux.Nombre, aux2.part.PartStart, aux2.fecha, aux.path, aux2.nombreMontada)
					///	GraficarSUPERBOOT("/home/josselyn/Escritorio/SUPER.png", id)
					AgregarUserTXT(aux2.part.PartStart, aux.path)

					return
				}
				formatear(aux2.ebr.PartSize, aux.Nombre, aux2.ebr.PartStart, aux2.fecha, aux.path, aux2.nombreMontada)
				AgregarUserTXT(aux2.ebr.PartStart, aux.path)
				return

			}
			aux2 = aux2.siguiente
		}
		aux = aux.siguiente
	}

	fmt.Println("No se ha encontrado el nombre de la partición que se desea formatear.")
}

//SetValores calcula los tamaños de las estructuras
func SetValores() (tamSB int64, tamAVD int64, tamDD int64, tamINodo int64, tamBloque int64, tamBitacora int64) {
	SB := SUPERBOOT{}
	tamSB = int64(unsafe.Sizeof(SB))
	AVD := ArbolVirtualDirectorio{}
	tamAVD = int64(unsafe.Sizeof(AVD))
	DD := DetalleDirectorio{}
	tamDD = int64(unsafe.Sizeof(DD))
	INODO := INodo{}
	tamINodo = int64(unsafe.Sizeof(INODO))
	BLOQUE := Bloque{}
	tamBloque = int64(unsafe.Sizeof(BLOQUE))
	BITACORA := LogBitacora{}
	tamBitacora = int64(unsafe.Sizeof(BITACORA))
	return
}
func formatear(size int64, NombreDisco string, inicioParticion int64, fechaM [19]byte, path string, nombre string) {
	tamSB, tamAVD, tamDD, tamINodo, tamBloque, tamBitacora := SetValores()
	tamEstructuras := (size - (2 * tamSB)) / (27 + tamAVD + tamDD + (5*tamINodo + (20 * tamBloque) + tamBitacora))
	cantAVD := tamEstructuras
	cantDD := tamEstructuras
	cantInodo := 5 * tamEstructuras
	cantBloques := 4 * cantInodo
	cantBitacora := tamEstructuras
	SB := SUPERBOOT{}
	copy(SB.SbNombreHD[:], NombreDisco)
	SB.SbAVDcount = cantAVD
	SB.SbDetalleDirectorioCount = cantDD
	SB.SbINodoCount = cantInodo
	SB.SbBloquesCount = cantBloques
	SB.SbAVDFree = cantAVD
	SB.SbDetalleDirectorioFree = cantDD
	SB.SbINodoFree = cantInodo
	SB.SbBloqueFree = cantBloques
	t := time.Now()
	fecha := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	copy(SB.SbDateCreation[:], fecha)
	SB.SbDateMontaje = fechaM

	SB.SbAptrStartBipmapAVD = inicioParticion + tamSB
	SB.SbAptrStartAVD = SB.SbAptrStartBipmapAVD + cantAVD
	SB.SbAptrStartBipmapDD = SB.SbAptrStartAVD + (tamAVD * cantAVD)
	SB.SbAptrStartDD = SB.SbAptrStartBipmapDD + cantDD
	SB.SbAptrStartBipmapINodo = SB.SbAptrStartDD + (tamDD * cantDD)
	SB.SbAptrStartINodo = SB.SbAptrStartBipmapINodo + cantInodo
	SB.SbAptrStartBipmapBloque = SB.SbAptrStartINodo + (tamINodo * cantInodo)
	SB.SbAptrStartBloque = SB.SbAptrStartBipmapBloque + cantBloques
	SB.SbAptrStartLogBitacora = SB.SbAptrStartBloque + (tamBloque + cantBloques)

	SB.SbSizeStructAVD = tamAVD
	SB.SbSizeStructDD = tamDD
	SB.SbSizeStructINodo = tamINodo
	SB.SbSizeStructBloque = tamBloque

	SB.SbFirstBitFreeAVD = SB.SbAptrStartBipmapAVD
	SB.SbFirstBitFreeDD = SB.SbAptrStartBipmapDD
	SB.SbFirstBitFreeINodo = SB.SbAptrStartBipmapINodo
	SB.SbFirstBitFreeBloque = SB.SbAptrStartBipmapBloque

	SB.SbMontajesCount++
	SB.SbMagicNum = 201602676

	files, err := os.OpenFile(path, os.O_RDWR, 0644)
	defer files.Close()
	if err != nil {
		log.Fatal(err)
	}
	//Se escribe el superboot
	files.Seek(inicioParticion, 0)
	var b3 bytes.Buffer
	binary.Write(&b3, binary.BigEndian, &SB)
	escribirBytes(files, b3.Bytes())
	//Se reserva espacio para el bipmap del AVD
	files.Seek(SB.SbAptrStartBipmapAVD, 0)
	var cero byte = 0
	i := int64(0)
	for i = 0; i < cantAVD; i++ {
		var b4 bytes.Buffer
		binary.Write(&b4, binary.BigEndian, &cero)
		escribirBytes(files, b4.Bytes())
	}

	start := SB.SbAptrStartAVD
	files.Seek(start, 0)
	for i = 0; i < cantAVD; i++ {
		var b5 bytes.Buffer
		avd := ArbolVirtualDirectorio{}
		binary.Write(&b5, binary.BigEndian, &avd)
		escribirBytes(files, b5.Bytes())
	}
	//Se reserva espacio para los bipmaps del DD
	files.Seek(SB.SbAptrStartBipmapDD, 0)
	for i = 0; i < cantDD; i++ {
		var b6 bytes.Buffer
		binary.Write(&b6, binary.BigEndian, &cero)
		escribirBytes(files, b6.Bytes())
	}
	//Se agregan las estructuras DD a la particion
	files.Seek(SB.SbAptrStartDD, 0)
	for i = 0; i < cantDD; i++ {
		var b7 bytes.Buffer
		DD := DetalleDirectorio{}
		binary.Write(&b7, binary.BigEndian, &DD)
		escribirBytes(files, b7.Bytes())
	}
	//Se asigna el espacio para el bipmap de los inodos
	files.Seek(SB.SbAptrStartBipmapINodo, 0)
	for i = 0; i < cantInodo; i++ {
		var b8 bytes.Buffer
		binary.Write(&b8, binary.BigEndian, &cero)
		escribirBytes(files, b8.Bytes())
	}
	//Se agregan las estructuras de Inodo a la particion
	files.Seek(SB.SbAptrStartINodo, 0)
	for i = 0; i < cantInodo; i++ {
		var b8 bytes.Buffer
		inodo := INodo{}
		binary.Write(&b8, binary.BigEndian, &inodo)
		escribirBytes(files, b8.Bytes())
	}
	//Se asigna espacio para el bipmap de los bloques
	files.Seek(SB.SbAptrStartBipmapBloque, 0)
	for i = 0; i < cantBloques; i++ {
		var b8 bytes.Buffer
		binary.Write(&b8, binary.BigEndian, &cero)
		escribirBytes(files, b8.Bytes())
	}
	files.Seek(SB.SbAptrStartBloque, 0)
	for i = 0; i < cantBloques; i++ {
		var b8 bytes.Buffer
		bloque := Bloque{}
		binary.Write(&b8, binary.BigEndian, &bloque)
		escribirBytes(files, b8.Bytes())
	}
	files.Seek(SB.SbAptrStartLogBitacora, 0)
	for i = 0; i < cantBitacora; i++ {
		var b8 bytes.Buffer
		bitacora := LogBitacora{}
		binary.Write(&b8, binary.BigEndian, &bitacora)
		escribirBytes(files, b8.Bytes())
	}
	fmt.Println(colorGreen, "************INFORMACIÓN**************")
	fmt.Println("Se ha formateado la partición " + nombre + " correctamente.")
}

//EscribirSUPERBOOT escribe el superboot al inicio de la particion
func EscribirSUPERBOOT(start int64, path string, super SUPERBOOT) {
	files, err := os.OpenFile(path, os.O_RDWR, 0644)
	defer files.Close()
	if err != nil {
		log.Fatal(err)
	}
	files.Seek(start, 0)
	var b3 bytes.Buffer
	binary.Write(&b3, binary.BigEndian, &super)
	escribirBytes(files, b3.Bytes())

}

//EscribirAVD reescribe la estructura avd
func EscribirAVD(start int64, path string, avd ArbolVirtualDirectorio) {
	files, err := os.OpenFile(path, os.O_RDWR, 0644)
	defer files.Close()
	if err != nil {
		log.Fatal(err)
	}
	files.Seek(start, 0)

	var b3 bytes.Buffer
	binary.Write(&b3, binary.BigEndian, &avd)
	escribirBytes(files, b3.Bytes())
}

//EscribirDD actualiza la estructura del dd
func EscribirDD(start int64, path string, dd DetalleDirectorio) {
	files, err := os.OpenFile(path, os.O_RDWR, 0644)
	defer files.Close()
	if err != nil {
		log.Fatal(err)
	}
	files.Seek(start, 0)

	var b3 bytes.Buffer
	binary.Write(&b3, binary.BigEndian, &dd)
	escribirBytes(files, b3.Bytes())
}

//EscribirINodo actualiza la información del inodo
func EscribirINodo(start int64, path string, inodo INodo) {
	files, err := os.OpenFile(path, os.O_RDWR, 0644)
	defer files.Close()
	if err != nil {
		log.Fatal(err)
	}
	files.Seek(start, 0)

	var b3 bytes.Buffer
	binary.Write(&b3, binary.BigEndian, &inodo)
	escribirBytes(files, b3.Bytes())
}

//EscribirBloque actualiza el bloque en la posicion solicitada
func EscribirBloque(start int64, path string, bloque Bloque) {
	files, err := os.OpenFile(path, os.O_RDWR, 0644)
	defer files.Close()
	if err != nil {
		log.Fatal(err)
	}
	files.Seek(start, 0)

	var b3 bytes.Buffer
	binary.Write(&b3, binary.BigEndian, &bloque)
	escribirBytes(files, b3.Bytes())
}

//LeerSUPERBOOT extra la estructura del superboot
func LeerSUPERBOOT(start int64, path string) (SUPERBOOT, bool) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		fmt.Println(colorRed, "No se encontró la ruta del archivo")
		return SUPERBOOT{}, false
	}
	super := SUPERBOOT{}
	tamSB := int(unsafe.Sizeof(super))
	file.Seek(start, 0)
	data := readNextBytes(file, tamSB)
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &super)
	if err != nil {
		panic(err)
	}

	return super, true

}

//LeerAVD este metodo extrae el AVD requerido de la partición
func LeerAVD(start int64, path string, tamAVD int64) (ArbolVirtualDirectorio, bool) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		fmt.Println(colorRed, "No se encontró la ruta del archivo")
		return ArbolVirtualDirectorio{}, false
	}
	AVD := ArbolVirtualDirectorio{}
	file.Seek(start, 0)
	data := readNextBytes(file, int(tamAVD))
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &AVD)
	if err != nil {
		panic(err)
	}

	return AVD, true
}

//LeerDD extrae el detalle de directorio en la posicion solicitada
func LeerDD(start int64, path string, tamDD int64) DetalleDirectorio {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		fmt.Println(colorRed, "No se encontró la ruta del archivo")
		return DetalleDirectorio{}
	}
	DD := DetalleDirectorio{}
	file.Seek(start, 0)
	data := readNextBytes(file, int(tamDD))
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &DD)
	if err != nil {
		panic(err)
	}
	return DD
}

//LeerINodo extrae el inodo en la posicion solicitada
func LeerINodo(start int64, path string, tamINodo int64) INodo {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		fmt.Println(colorRed, "No se encontró la ruta del archivo")
		return INodo{}
	}
	INodo := INodo{}
	file.Seek(start, 0)
	data := readNextBytes(file, int(tamINodo))
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &INodo)
	if err != nil {
		panic(err)
	}
	return INodo
}

//LeerBloque extrae el bloque en la posición solicitada
func LeerBloque(start int64, path string, tamBloque int64) Bloque {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		fmt.Println(colorRed, "No se encontró la ruta del archivo")
		return Bloque{}
	}
	Bloque := Bloque{}
	file.Seek(start, 0)
	data := readNextBytes(file, int(tamBloque))
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &Bloque)
	if err != nil {
		panic(err)
	}
	return Bloque
}

//BuscarParticionM busca el id de la particion montada
func BuscarParticionM(id string) (int64, string) {
	id = strings.ToLower(id)
	if strings.Contains(id, "vd") {
		letraAbuscar := id[2]
		var nodoAux *NodoDisco
		nodoAux = ListDiscos.inicio
		for nodoAux != nil {
			if nodoAux.Letra == letraAbuscar {
				var NodoAuxL *NodoParticion
				NodoAuxL = nodoAux.listaParticiones.inicio
				for NodoAuxL != nil {
					if id == NodoAuxL.nombreMontada {
						if NodoAuxL.tipo == "primaria" {
							return NodoAuxL.part.PartStart, nodoAux.path
						}
						return NodoAuxL.ebr.PartStart, nodoAux.path

					}
					NodoAuxL = NodoAuxL.siguiente
				}
			}
			nodoAux = nodoAux.siguiente
		}
	}
	return -1, ""
}

//GraficarSUPERBOOT Grafica en una tabla la información del sistema de archivos
func GraficarSUPERBOOT(ubicacion string, start int64, path string) {
	dir := ""
	rutas := strings.Split(ubicacion, "/")
	for i := 0; i < len(rutas)-1; i++ {
		dir += rutas[i] + "/"
	}
	nombre := rutas[len(rutas)-1]
	extension := strings.Split(nombre, ".")

	if AnalizarRuta(dir) {

		dir = dir + extension[0] + ".txt"
		var _, errr = os.Stat(dir)
		//Crea el archivo si no existe
		if os.IsNotExist(errr) {
			var file, errr = os.Create(dir)
			if existeError(errr) {
				return
			}
			defer file.Close()
		}

		super, auxi := LeerSUPERBOOT(start, path)
		if start != -1 && auxi != false {
			cadena := ""
			cadena += "digraph G {\ngraph [pad=\"0.5\", nodesep=\"1\", ranksep=\"2\"];"
			cadena += "\nnode [shape=plain]\nrankdir=LR;\n"
			cadena += "Tabla[label=<\n<table border=\"1\" cellborder=\"1\" cellspacing=\"1\">\n"
			cadena += "<tr><td bgcolor=\"#76D7C4\"><i>Nombre</i></td>\n<td bgcolor=\"#76D7C4\"><i>Valor</i> </td>\n</tr>"

			cadena += "<tr><td bgcolor=\"#A3E4D7\">Sb_Nombre_HD</td><td>" + Nombres(super.SbNombreHD) + "</td></tr>\n"
			cadena += "<tr><td bgcolor=\"#A3E4D7\">Sb_Arbol_Virtual_Count</td><td>" + strconv.Itoa(int(super.SbAVDcount)) + "</td></tr>\n"
			cadena += "<tr><td bgcolor=\"#A3E4D7\">Sb__Detalle_Directorio_Count</td><td>" + strconv.Itoa(int(super.SbDetalleDirectorioCount)) + "</td></tr>\n"
			cadena += "<tr><td bgcolor=\"#A3E4D7\">Sb_INodos_Count</td><td>" + strconv.Itoa(int(super.SbINodoCount)) + "</td></tr>\n"
			cadena += "<tr><td bgcolor=\"#A3E4D7\">Sb_Bloques_Count</td><td>" + strconv.Itoa(int(super.SbBloquesCount)) + "</td></tr>\n"
			cadena += "<tr><td bgcolor=\"#A3E4D7\">Sb_Arbol_Virtual_Free</td><td>" + strconv.Itoa(int(super.SbAVDFree)) + "</td></tr>\n"
			cadena += "<tr><td bgcolor=\"#A3E4D7\">Sb_Detalle_Directorio_Free</td><td>" + strconv.Itoa(int(super.SbDetalleDirectorioFree)) + "</td></tr>\n"
			cadena += "<tr><td bgcolor=\"#A3E4D7\">Sb_INodo_Free</td><td>" + strconv.Itoa(int(super.SbINodoFree)) + "</td></tr>\n"
			cadena += "<tr><td bgcolor=\"#A3E4D7\">Sb_Bloques_Free</td><td>" + strconv.Itoa(int(super.SbBloqueFree)) + "</td></tr>\n"
			cadena += "<tr><td bgcolor=\"#A3E4D7\">Sb_Date_Creation</td><td>" + string(super.SbDateCreation[:]) + "</td></tr>\n"
			cadena += "<tr><td bgcolor=\"#A3E4D7\">Sb_Date_Montaje</td><td>" + string(super.SbDateMontaje[:]) + "</td></tr>\n"
			cadena += "<tr><td bgcolor=\"#A3E4D7\">Sb_MontajesCount</td><td>" + strconv.Itoa(int(super.SbMontajesCount)) + "</td></tr>\n"
			cadena += "<tr><td bgcolor=\"#A3E4D7\">Sb_Aptr_Start_Bipmap_AVD</td><td>" + strconv.Itoa(int(super.SbAptrStartBipmapAVD)) + "</td></tr>\n"
			cadena += "<tr><td bgcolor=\"#A3E4D7\">Sb_Aptr_Start_AVD</td><td>" + strconv.Itoa(int(super.SbAptrStartAVD)) + "</td></tr>\n"
			cadena += "<tr><td bgcolor=\"#A3E4D7\">Sb_Aptr_Start_Bipmap_DD</td><td>" + strconv.Itoa(int(super.SbAptrStartBipmapDD)) + "</td></tr>\n"
			cadena += "<tr><td bgcolor=\"#A3E4D7\">Sb_Aptr_Start_DD</td><td>" + strconv.Itoa(int(super.SbAptrStartDD)) + "</td></tr>\n"
			cadena += "<tr><td bgcolor=\"#A3E4D7\">Sb_Aptr_Start_Bipmap_INodo</td><td>" + strconv.Itoa(int(super.SbAptrStartBipmapINodo)) + "</td></tr>\n"
			cadena += "<tr><td bgcolor=\"#A3E4D7\">Sb_Aptr_Start_INodo</td><td>" + strconv.Itoa(int(super.SbAptrStartINodo)) + "</td></tr>\n"
			cadena += "<tr><td bgcolor=\"#A3E4D7\">Sb_Aptr_Start_Bipmap_Bloque</td><td>" + strconv.Itoa(int(super.SbAptrStartBipmapBloque)) + "</td></tr>\n"
			cadena += "<tr><td bgcolor=\"#A3E4D7\">Sb_Aptr_Start_Bloque</td><td>" + strconv.Itoa(int(super.SbAptrStartBloque)) + "</td></tr>\n"
			cadena += "<tr><td bgcolor=\"#A3E4D7\">Sb_Aptr_Start_Bitacora</td><td>" + strconv.Itoa(int(super.SbAptrStartLogBitacora)) + "</td></tr>\n"
			cadena += "<tr><td bgcolor=\"#A3E4D7\">Sb_Size_Struct_AVD</td><td>" + strconv.Itoa(int(super.SbSizeStructAVD)) + "</td></tr>\n"
			cadena += "<tr><td bgcolor=\"#A3E4D7\">Sb_Size_Struct_DD</td><td>" + strconv.Itoa(int(super.SbSizeStructDD)) + "</td></tr>\n"
			cadena += "<tr><td bgcolor=\"#A3E4D7\">Sb_Size_Struct_Inodo</td><td>" + strconv.Itoa(int(super.SbSizeStructINodo)) + "</td></tr>\n"
			cadena += "<tr><td bgcolor=\"#A3E4D7\">Sb_Size_Struct_Bloques</td><td>" + strconv.Itoa(int(super.SbSizeStructBloque)) + "</td></tr>\n"
			cadena += "<tr><td bgcolor=\"#A3E4D7\">Sb_First_Free_Bit_AVD</td><td>" + strconv.Itoa(int(super.SbFirstBitFreeAVD)) + "</td></tr>\n"
			cadena += "<tr><td bgcolor=\"#A3E4D7\">Sb_First_Free_Bit_DD</td><td>" + strconv.Itoa(int(super.SbFirstBitFreeDD)) + "</td></tr>\n"
			cadena += "<tr><td bgcolor=\"#A3E4D7\">Sb_First_Free_Bit_INodo</td><td>" + strconv.Itoa(int(super.SbFirstBitFreeINodo)) + "</td></tr>\n"
			cadena += "<tr><td bgcolor=\"#A3E4D7\">Sb_First_Free_Bit_Bloques</td><td>" + strconv.Itoa(int(super.SbFirstBitFreeBloque)) + "</td></tr>\n"
			cadena += "<tr><td bgcolor=\"#A3E4D7\">Sb_Magic_Num</td><td>" + strconv.Itoa(int(super.SbMagicNum)) + "</td></tr>\n"
			cadena += "</table>>];}"
			errrr := ioutil.WriteFile(dir, []byte(cadena[:]), 0644)
			if errrr != nil {
				panic(errrr)
			}
			com1 := "dot"
			com2 := "-T" + extension[1]
			com3 := dir
			com4 := "-o"
			com5 := ubicacion
			exec.Command(com1, com2, com3, com4, com5).Output()
			fmt.Println(colorGreen, "Success")
		} else {
			fmt.Println("El ID de la partición no ha sido encontrado!")
		}
	}
}

//REPORTE comando para reportar todos los resultados
func REPORTE(arreglo []string) {
	name := ""
	path := ""
	id := ""
	var start int64
	var direccion string
	ruta := ""
	for i := 1; i < len(arreglo); i++ {
		comandos := strings.Split(arreglo[i], "->")
		com := strings.ToLower(comandos[0])
		switch com {
		case "-name":
			n := strings.ToLower(comandos[1])
			if NombreReporte(n) {
				name = n
			} else {
				fmt.Println("No existe un reporte con ese nombre!")
				return
			}
			break
		case "-path":
			path = comandos[1]
			if strings.Contains(comandos[1], "@") {
				path = strings.ReplaceAll(comandos[1], "@", " ")
			}
			if path == "" {
				fmt.Println("Ruta no válida")
				return
			}
			break
		case "-id":
			p := strings.ToLower(comandos[1])
			start, direccion = BuscarParticionM(p)
			if start != -1 && direccion != "" {
				id = p
			} else {
				fmt.Println(colorRed, "No se ha encontrado el id de partición")
				return
			}
			break
		case "-ruta":
			if strings.Contains(comandos[1], "@") {
				ruta = strings.ReplaceAll(comandos[1], "@", " ")
			} else {
				ruta = comandos[1]
			}
			break
		case "":
			break
		default:
			fmt.Println(colorRed, "Comando inválido!")
			break
		}
	}

	if name != "" && id != "" && path != "" {
		SeleccionarReporte(name, start, direccion, path, ruta)
	} else {
		fmt.Println(colorRed, "Faltan parámetros para completar la instrucción")
	}
}

//NombreReporte verifica si el nombre del directorio es correcto
func NombreReporte(nombre string) bool {
	switch nombre {
	case "mbr":
		return true
	case "disk":
		return true
	case "sb":
		return true
	case "bm_arbdir":
		return true
	case "bm_detdir":
		return true
	case "bm_inode":
		return true
	case "bm_block":
		return true
	case "bitacora":
		return true
	case "directorio":
		return true
	case "tree_file":
		return true
	case "tree_directorio":
		return true
	case "tree_complete":
		return true
	case "ls":
		return true
	default:
		return false
	}
}

//SeleccionarReporte aqui se hace la seleccion del reporte que se quiere mostrar
func SeleccionarReporte(name string, start int64, dirDisco string, ubicacionReporte string, dirBuscar string) {
	switch name {
	case "mbr":
		graficarMBR(dirDisco, ubicacionReporte)
		break
	case "disk":
		GraficarDisco(dirDisco, ubicacionReporte)
		break
	case "sb":
		GraficarSUPERBOOT(ubicacionReporte, start, dirDisco)
		break
	case "bm_arbdir":
		dd := false
		super, dd = LeerSUPERBOOT(start, dirDisco)
		if dd {
			GraficarBitMap(ubicacionReporte, super.SbAptrStartBipmapAVD, super.SbSizeStructAVD, dirDisco)
		}
		break
	case "bm_detdir":
		dd := false
		super, dd = LeerSUPERBOOT(start, dirDisco)
		if dd {
			GraficarBitMap(ubicacionReporte, super.SbAptrStartBipmapDD, super.SbSizeStructDD, dirDisco)
		}
		break
	case "bm_inode":
		dd := false
		super, dd = LeerSUPERBOOT(start, dirDisco)
		if dd {
			GraficarBitMap(ubicacionReporte, super.SbAptrStartBipmapINodo, super.SbSizeStructINodo, dirDisco)
		}
		break
	case "bm_block":
		dd := false
		super, dd = LeerSUPERBOOT(start, dirDisco)
		if dd {
			GraficarBitMap(ubicacionReporte, super.SbAptrStartBipmapBloque, super.SbSizeStructBloque, dirDisco)
		}
		break
	case "directorio":
		GraficarDirectorio(ubicacionReporte, start, dirDisco)
		break
	}
}

//AgregarUserTXT se agrega el archivo a la raiz
func AgregarUserTXT(start int64, path string) {
	cadena := "1,G,root\n"
	cadena += "1,U,root,root,201602676\n"
	super, valido := LeerSUPERBOOT(start, path)
	if valido {
		file, err := os.OpenFile(path, os.O_RDWR, 0644)
		defer file.Close()
		if err != nil {
			fmt.Println(colorRed, "No se encontró la ruta del archivo")

		}
		var ocupado byte = '1'
		//Actualiza el bipmap del avd
		file.Seek(super.SbAptrStartBipmapAVD, 0)
		var b3 bytes.Buffer
		binary.Write(&b3, binary.BigEndian, &ocupado)
		escribirBytes(file, b3.Bytes())

		avd, aux := LeerAVD(super.SbAptrStartAVD, path, super.SbSizeStructAVD)
		if aux {
			super.SbAVDFree--
			copy(avd.AVDNameDirectoy[:], "/")
			t := time.Now()
			fecha := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", t.Year(), t.Month(), t.Day(),
				t.Hour(), t.Minute(), t.Second())
			copy(avd.AVDFechaCreacion[:], fecha)
			avd.AVDAptrDetalleDirectorio = 1
			avd.AVDProper = 1
			EscribirAVD(super.SbAptrStartAVD, path, avd)
			super.SbDetalleDirectorioFree--
			//Actualiza el bipmap del DD
			file.Seek(super.SbAptrStartBipmapDD, 0)
			var b3 bytes.Buffer
			binary.Write(&b3, binary.BigEndian, &ocupado)
			escribirBytes(file, b3.Bytes())
			dd := LeerDD(super.SbAptrStartDD, path, super.SbSizeStructDD)
			dd.DDArrayAptrINodo[0].DDAptrINodo = 1
			copy(dd.DDArrayAptrINodo[0].DDFileDateCreation[:], fecha)
			copy(dd.DDArrayAptrINodo[0].DDFileDateModificacion[:], fecha)
			copy(dd.DDArrayAptrINodo[0].DDfileName[:], "user.txt")
			EscribirDD(super.SbAptrStartDD, path, dd)
			super.SbINodoFree--
			//Actualiza el bipmap del inodo
			file.Seek(super.SbAptrStartBipmapINodo, 0)
			var b4 bytes.Buffer
			binary.Write(&b4, binary.BigEndian, &ocupado)
			escribirBytes(file, b4.Bytes())
			inodo := INodo{}
			inodo.INodoNumero = 1
			inodo.INodoFileSize = 40
			inodo.INodoBloquesAsignados = 2
			inodo.INodoAptrDeBloque[0] = 1
			inodo.INodoAptrDeBloque[1] = 2
			inodo.INodoIDPropier = 1
			inodo.INodoGid = 1
			EscribirINodo(super.SbAptrStartINodo, path, inodo)
			super.SbBloqueFree = super.SbBloqueFree - 2
			//Actualiza el bipmap del bloque
			file.Seek(super.SbAptrStartBipmapBloque, 0)
			var b5 bytes.Buffer
			binary.Write(&b5, binary.BigEndian, &ocupado)
			escribirBytes(file, b5.Bytes())
			bloque1 := Bloque{}
			bloque2 := Bloque{}
			resto := 0
			tam := len(cadena)
			for i := 0; i < len(cadena); i++ {
				if i < 25 {
					bloque1.BDarray[i] = cadena[i]
					resto = i
				}
			}
			resto--
			cad := tam - resto
			for i := 0; i < cad; i++ {
				bloque2.BDarray[i] = cadena[resto]
				resto++
			}
			EscribirBloque(super.SbAptrStartBloque, path, bloque1)
			EscribirBloque(super.SbAptrStartBloque+super.SbSizeStructBloque, path, bloque2)
			super.SbFirstBitFreeAVD = super.SbAptrStartBipmapAVD + 1
			super.SbFirstBitFreeDD = super.SbAptrStartBipmapDD + 1
			super.SbFirstBitFreeINodo = super.SbAptrStartBipmapINodo + 1
			super.SbFirstBitFreeBloque = super.SbAptrStartBipmapBloque + 2
			EscribirSUPERBOOT(start, path, super)
		}

	}
}

//MkFile analiza los comandos que vienen después
func MkFile(arreglo []string) {
	id := ""
	path := ""
	crear := false
	size := 0
	cadena := ""
	start := int64(0)
	dir := ""
	for i := 1; i < len(arreglo); i++ {
		comandos := strings.Split(arreglo[i], "->")
		com := strings.ToLower(comandos[0])
		switch com {
		case "-id":
			start, dir = BuscarParticionM(strings.ToLower(comandos[1]))
			if dir != "" && start != -1 {
				id = strings.ToLower(comandos[1])

			} else {
				fmt.Println(colorRed, "El id de la partición no fue encontrado")
				return
			}
			break
		case "-path":
			path = direccion(comandos[1])
			if path != "" {

			} else {
				fmt.Println(colorRed, "Dirección inválida")
			}
			break
		case "-p":
			crear = true
			break
		case "-size":
			n, err := strconv.Atoi(comandos[i])
			if err != nil {
				fmt.Println(colorRed, "Verifique el valor del size.")
				return
			}
			size = n
			break
		case "-cont":
			cadena = strings.ReplaceAll(comandos[i], "\"", "")

			break
		case "":
			break
		}
	}

	if id != "" && path != "" {
		CrearArchivos(id, path, cadena, int64(size), crear, start, dir)
	} else {
		fmt.Println(colorRed, "Faltan parámetros obligatorios.")
	}

}

//CrearArchivos verifica el archivo a crear
func CrearArchivos(id string, path string, cont string, size int64, crear bool, start int64, dir string) {
	b := false
	super, b = LeerSUPERBOOT(start, dir)
	if b {
		direccion := strings.Split(path, "/")
		archivo := direccion[len(direccion)-1]
		var carpetas []string
		for i := 0; i < (len(direccion) - 1); i++ {
			carpetas[i] = direccion[i]
		}
		if direccion[0] == "" {
			direccion[0] = "/"
		} else {
			fmt.Println(colorRed, "Dirección incorrecta")
			return
		}
		avd, b := LeerAVD(super.SbAptrStartAVD, dir, super.SbSizeStructAVD)
		if b {
			var inicio int64
			if crear {
				inicio = CrearCarpeta(direccion[0], carpetas, avd, dir, 1, crear, start, super.SbSizeStructAVD)
			} else {
				inicio = buscarCarpeta(carpetas, direccion[0], avd, dir, 1, start, super.SbSizeStructAVD)
			}
			if inicio != -1 {
				EntrarAVD(inicio, dir, archivo, cont, size)
			} else {
				fmt.Println(colorYellow, "No se ha encontrado el directorio solicitado")
				return
			}
		}
	}
}

//EntrarAVD verifica si el AVD ya tiene un arbol de directorio
func EntrarAVD(AVDStart int64, path string, archivo string, cont string, size int64) {
	avd, p := LeerAVD(AVDStart, path, super.SbSizeStructAVD)
	if p {
		if avd.AVDAptrDetalleDirectorio != 0 {
			CrearDD(path, archivo, cont, size, avd.AVDAptrDetalleDirectorio)
		} else {
			//Se crea un nuevo apuntador en la siguiente posición libre de DD
			op := ((super.SbAptrStartBipmapDD + super.SbDetalleDirectorioCount) - super.SbFirstBitFreeDD)
			apuntador := (super.SbDetalleDirectorioCount - op) + 1
			avd.AVDAptrDetalleDirectorio = apuntador
			//Se modifica la información del AVD asignandole un valor a la celda de DD
			EscribirAVD(AVDStart, path, avd)
			CrearDDInodo(path, archivo, cont, size, 0, apuntador)
			return
		}
	}
}

//CrearDD hace los Detalle de directorio
func CrearDD(path string, archivo string, cont string, size int64, aptr int64) {
	inicioDD := super.SbAptrStartDD + (super.SbSizeStructDD * (aptr - 1))
	dd := LeerDD(inicioDD, path, super.SbSizeStructDD)
	for i := 0; i < len(dd.DDArrayAptrINodo); i++ {
		var vacio [20]byte
		if dd.DDArrayAptrINodo[i].DDfileName == vacio {
			CrearDDInodo(path, archivo, cont, size, i, aptr)
			return
		}
	}
	if dd.DDAptrIndirecto != 0 {
		CrearDD(path, archivo, cont, size, dd.DDAptrIndirecto)
	} else {
		op := ((super.SbAptrStartBipmapDD + super.SbDetalleDirectorioCount) - super.SbFirstBitFreeDD)
		apuntador := (super.SbDetalleDirectorioCount - op) + 1
		dd.DDAptrIndirecto = apuntador
		CrearDDInodo(path, archivo, cont, size, 0, apuntador)
	}

}

//CrearDDInodo crea los archivos
func CrearDDInodo(path string, archivo string, cont string, size int64, posicion int, apuntador int64) {

	//Se calcula el start de la estructura del detalle de directorio
	inicioDD := super.SbAptrStartDD + (super.SbSizeStructDD * (apuntador - 1))
	//Se crea el nuevoDD que se va a meter en el archivo con todos sus datos
	NuevoDD := DetalleDirectorio{}
	t := time.Now()
	fecha := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	copy(NuevoDD.DDArrayAptrINodo[posicion].DDFileDateCreation[:], fecha)
	copy(NuevoDD.DDArrayAptrINodo[posicion].DDFileDateModificacion[:], fecha)
	copy(NuevoDD.DDArrayAptrINodo[posicion].DDfileName[:], archivo)
	//Se crea el apuntador nuevo del INODO
	apuntadorInodo := (super.SbAptrStartBipmapINodo + super.SbINodoCount) - super.SbFirstBitFreeINodo
	aInodo := super.SbINodoCount - apuntadorInodo + 1
	//Se le asigna el apuntador
	NuevoDD.DDArrayAptrINodo[posicion].DDAptrINodo = aInodo
	//Se guarda el detalle de directorio en el archivo
	EscribirDD(inicioDD, path, NuevoDD)
	//Se actualiza el bipmap del DD
	super.SbFirstBitFreeDD = actualizarBipmap(path, super.SbFirstBitFreeDD, (super.SbAptrStartBipmapDD + super.SbSizeStructDD))
	//Se actualiza el bipmap del inodo
	super.SbFirstBitFreeINodo = actualizarBipmap(path, super.SbFirstBitFreeINodo, (super.SbAptrStartBipmapINodo + super.SbSizeStructINodo))
	comienzoInodo := super.SbAptrStartINodo + (super.SbSizeStructINodo * (aInodo - 1))
	NuevoInodo := INodo{}
	NuevoInodo.INodoNumero = aInodo
	tamCad := 0
	if cont != "" {
		tamCad = len(cont)
	}
	if int64(tamCad) > size {
		size = int64(tamCad)
	}
	redondear := float64(size / 25)
	a := int64(Roundf(redondear))
	NuevoInodo.INodoBloquesAsignados = a
	NuevoInodo.INodoFileSize = size
	auxi := a
	for i := 0; i < 4; i++ {
		var cad [25]byte
		cad, cont = dividirCadena(cont)
		NuevoInodo = CrearBloque(path, NuevoInodo, i, cad)
		a--
	}
	if a != 0 {
		apuntadorInodoind := (super.SbAptrStartBipmapINodo + super.SbINodoCount) - super.SbFirstBitFreeINodo
		aInodoind := super.SbINodoCount - apuntadorInodoind + 1
		NuevoInodo.INodoAptrInd = aInodoind
		EscribirINodo(comienzoInodo, path, NuevoInodo)
		CrearINodo(path, cont, size, NuevoInodo.INodoNumero, aInodoind, auxi)
	}
}

//CrearINodo crea todos los inodos que se necesiten
func CrearINodo(path string, cont string, size int64, numero int64, aptrind int64, a int64) {
	comienzoInodo := super.SbAptrStartINodo + (super.SbSizeStructINodo * (aptrind - 1))
	NuevoInodo := INodo{}
	NuevoInodo.INodoNumero = numero
	NuevoInodo.INodoBloquesAsignados = a
	NuevoInodo.INodoFileSize = size
	auxi := a
	for i := 0; i < 4; i++ {
		var cad [25]byte
		cad, cont = dividirCadena(cont)
		NuevoInodo = CrearBloque(path, NuevoInodo, i, cad)
		a--
	}
	if a != 0 {
		apuntadorInodoind := (super.SbAptrStartBipmapINodo + super.SbINodoCount) - super.SbFirstBitFreeINodo
		aInodoind := super.SbINodoCount - apuntadorInodoind + 1
		NuevoInodo.INodoAptrInd = aInodoind
		super.SbFirstBitFreeINodo = actualizarBipmap(path, super.SbFirstBitFreeINodo, (super.SbAptrStartBipmapINodo + super.SbSizeStructINodo))
		EscribirINodo(comienzoInodo, path, NuevoInodo)
		CrearINodo(path, cont, size, NuevoInodo.INodoNumero, aInodoind, auxi)
	}
}
func dividirCadena(cadena string) (cad [25]byte, resto string) {
	noSeUsa := ""
	for i := 0; i < len(cadena); i++ {
		noSeUsa += string(rune(cadena[i]))
		cad[i] = cadena[i]
	}
	resto = strings.Replace(cadena, noSeUsa, "", 1)
	return cad, resto
}

//RetornarArreglo divide el contenido
func RetornarArreglo(size int64) [25]byte {
	var aux [25]byte

	var letra byte = 65
	var i int64
	for i = 0; i < size; i++ {
		aux[i] = letra
		letra++
		if letra == 91 {
			letra = 65
		}
	}
	return aux
}

//CrearBloque va creando los bloques depende del tamaño del archivo
func CrearBloque(path string, NuevoInodo INodo, i int, arreglo [25]byte) INodo {
	apuntadorBloque := (super.SbAptrStartBipmapBloque + super.SbBloquesCount) - super.SbFirstBitFreeBloque
	aBloque := super.SbBloquesCount - apuntadorBloque + 1
	inicioBloque := super.SbAptrStartBloque + (super.SbSizeStructBloque * (aBloque - 1))
	super.SbFirstBitFreeBloque = actualizarBipmap(path, super.SbFirstBitFreeBloque, (super.SbAptrStartBipmapBloque + super.SbSizeStructBloque))
	nuevoBloque := Bloque{}
	NuevoInodo.INodoAptrDeBloque[i] = aBloque
	//Añade el bloque
	nuevoBloque.BDarray = arreglo
	EscribirBloque(inicioBloque, path, nuevoBloque)
	return NuevoInodo
}

//Roundf dd
func Roundf(x float64) float64 {
	t := math.Trunc(x)
	if math.Abs(x-t) >= 0 {
		return t + math.Copysign(1, x)
	}
	return t
}
func actualizarBipmap(path string, primero int64, fin int64) int64 {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	defer file.Close()
	if err != nil {
		fmt.Println(colorRed, "No se encontró la ruta del archivo")
	} else if primero <= fin {
		var ocupado byte = '1'
		//Actualiza el bipmap
		file.Seek(primero, 0)
		var b3 bytes.Buffer
		binary.Write(&b3, binary.BigEndian, &ocupado)
		escribirBytes(file, b3.Bytes())

		//Busca el siguiente bipmap desocupado

		act := false
		var cont int64 = 0
		for !act {
			var oc byte
			file.Seek((primero + cont), 0)
			data := readNextBytes(file, int(unsafe.Sizeof(oc)))
			buffer := bytes.NewBuffer(data)
			err = binary.Read(buffer, binary.BigEndian, &oc)
			if err != nil {
				panic(err)
			}
			if oc == 0 {
				primero = primero + cont
				act = true
			}
			cont++
		}
	}
	return primero
}

func buscarCarpeta(carpetas []string, padre string, avd ArbolVirtualDirectorio, path string, siguiente int, ss int64, inicioAVD int64) int64 {

	for i := 0; i < len(avd.AVDAptrArraySubdirectorios); i++ {
		if avd.AVDAptrArraySubdirectorios[i] != 0 {
			avdAux, p := LeerAVD(super.SbAptrStartAVD+((avd.AVDAptrArraySubdirectorios[i]-1)*super.SbSizeStructAVD), path, super.SbSizeStructAVD)
			if p {
				if Nombres(avdAux.AVDNameDirectoy) == carpetas[siguiente] {
					fmt.Println(avd)
					padre = carpetas[siguiente]
					if siguiente < len(carpetas)-1 {
						siguiente++
						inicioAVD = super.SbAptrStartAVD + ((avd.AVDAptrArraySubdirectorios[i] - 1) * super.SbSizeStructAVD)
						return buscarCarpeta(carpetas, padre, avdAux, path, siguiente, ss, inicioAVD)
					}
					return inicioAVD

				}
			}
		}

	}
	if avd.AVDAptrInd != 0 {
		inicioAVD = super.SbAptrStartAVD + ((avd.AVDAptrInd - 1) * super.SbSizeStructAVD)
		avdAux, p := LeerAVD(inicioAVD, path, super.SbSizeStructAVD)
		if p {
			buscarCarpeta(carpetas, padre, avdAux, path, siguiente, ss, inicioAVD)
		}
	}
	return -1
}

//MkDir crear carpetas
func MkDir(arreglo []string) {
	id := ""
	path := ""
	crear := false
	start := int64(0)
	dir := ""
	for i := 1; i < len(arreglo); i++ {
		comandos := strings.Split(arreglo[i], "->")
		com := strings.ToLower(comandos[0])
		switch com {
		case "-id":
			start, dir = BuscarParticionM(strings.ToLower(comandos[1]))
			if dir != "" && start != -1 {
				id = strings.ToLower(comandos[1])

			} else {
				fmt.Println(colorRed, "El id de la partición no fue encontrado")
				return
			}
			break
		case "-path":
			path = direccion(arreglo[i])
			if path == "" {
				fmt.Println(colorRed, "Dirección inválida")
			}
			break
		case "-p":
			crear = true
			break

		case "":
			break
		}
	}

	if id != "" && path != "" {
		crearAC(id, path, crear, start, dir)
	} else {
		fmt.Println(colorRed, "Faltan parámetros obligatorios.")
	}
}

//crearAC crear archivos y carpetas
func crearAC(id string, path string, crear bool, start int64, dir string) {
	carpetas := strings.Split(path, "/")
	var a = false
	super, a = LeerSUPERBOOT(start, dir)
	if a {
		avds := ArbolVirtualDirectorio{}
		avds, aux := LeerAVD(super.SbAptrStartAVD, dir, super.SbSizeStructAVD)
		if aux {
			if carpetas[0] == "" {
				carpetas[0] = "/"
			}
			CrearCarpeta(carpetas[0], carpetas, avds, dir, 1, crear, start, super.SbAptrStartAVD)
			fmt.Println(colorGreen, "**************Información***************")
			fmt.Println(colorGreen, " Se han creado las o la carpeta(s)")
			fmt.Println(colorGreen, "****************************************")
		}
	}

}

//CrearCarpeta se crean carpeta
func CrearCarpeta(carpetaPadre string, conjCarpetas []string, avd ArbolVirtualDirectorio, path string, siguiente int, crear bool, ss int64, inicioAVD int64) int64 {
	aux := 0

	for i := 0; i < len(avd.AVDAptrArraySubdirectorios); i++ {
		if avd.AVDAptrArraySubdirectorios[i] != 0 {
			aux++
			avdAux, p := LeerAVD(super.SbAptrStartAVD+((avd.AVDAptrArraySubdirectorios[i]-1)*super.SbSizeStructAVD), path, super.SbSizeStructAVD)
			if p {
				if Nombres(avdAux.AVDNameDirectoy) == conjCarpetas[siguiente] {
					carpetaPadre = conjCarpetas[siguiente]
					if siguiente < len(conjCarpetas)-1 {
						siguiente++
						inicioAVD = super.SbAptrStartAVD + ((avd.AVDAptrArraySubdirectorios[i] - 1) * super.SbSizeStructAVD)
						return CrearCarpeta(carpetaPadre, conjCarpetas, avdAux, path, siguiente, crear, ss, inicioAVD)
					}
					return inicioAVD

				}
			}
		}
	}
	if avd.AVDAptrInd != 0 {
		aux = 6
	}
	if aux == 6 {
		if avd.AVDAptrInd != 0 {
			inicioAVD = super.SbAptrStartAVD + ((avd.AVDAptrInd - 1) * super.SbSizeStructAVD)
			avdAux, p := LeerAVD(inicioAVD, path, super.SbSizeStructAVD)
			if p {
				return CrearCarpeta(carpetaPadre, conjCarpetas, avdAux, path, siguiente, crear, ss, inicioAVD)
			}
		} else {
			op := ((super.SbAptrStartBipmapAVD + super.SbAVDcount) - super.SbFirstBitFreeAVD)
			apuntador := (super.SbAVDcount - op) + 1
			//Se abre el archivo para actualizar bipmap
			super.SbFirstBitFreeAVD = actualizarBipmap(path, super.SbFirstBitFreeAVD, (super.SbAptrStartBipmapAVD + super.SbSizeStructAVD))
			avd.AVDAptrInd = apuntador
			EscribirAVD(inicioAVD, path, avd)
			nuevoAVD := ArbolVirtualDirectorio{}
			t := time.Now()
			fecha := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", t.Year(), t.Month(), t.Day(),
				t.Hour(), t.Minute(), t.Second())
			copy(nuevoAVD.AVDFechaCreacion[:], fecha)
			copy(nuevoAVD.AVDNameDirectoy[:], carpetaPadre)
			inicioAVD = super.SbAptrStartAVD + (super.SbSizeStructAVD * (apuntador - 1))
			EscribirAVD(inicioAVD, path, nuevoAVD)
			super.SbAVDFree--
			EscribirSUPERBOOT(ss, path, super)
			return CrearCarpeta(carpetaPadre, conjCarpetas, nuevoAVD, path, siguiente, crear, ss, inicioAVD)

		}
	} else if crear == true && super.SbAVDFree != 0 {
		for i := 0; i < len(avd.AVDAptrArraySubdirectorios); i++ {
			if avd.AVDAptrArraySubdirectorios[i] == 0 {
				op := ((super.SbAptrStartBipmapAVD + super.SbAVDcount) - super.SbFirstBitFreeAVD)
				apuntador := (super.SbAVDcount - op) + 1
				//Se abre el archivo para actualizar bipmap
				super.SbFirstBitFreeAVD = actualizarBipmap(path, super.SbFirstBitFreeAVD, (super.SbAptrStartBipmapAVD + super.SbSizeStructAVD))
				avd.AVDAptrArraySubdirectorios[i] = apuntador
				EscribirAVD(inicioAVD, path, avd)
				nuevoAVD := ArbolVirtualDirectorio{}
				t := time.Now()
				fecha := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", t.Year(), t.Month(), t.Day(),
					t.Hour(), t.Minute(), t.Second())
				copy(nuevoAVD.AVDFechaCreacion[:], fecha)
				copy(nuevoAVD.AVDNameDirectoy[:], conjCarpetas[siguiente])
				inicioAVD = super.SbAptrStartAVD + (super.SbSizeStructAVD * (apuntador - 1))
				EscribirAVD(inicioAVD, path, nuevoAVD)
				super.SbAVDFree--
				EscribirSUPERBOOT(ss, path, super)
				if siguiente < len(conjCarpetas)-1 {
					carpetaPadre = conjCarpetas[siguiente]
					siguiente++
					return CrearCarpeta(carpetaPadre, conjCarpetas, nuevoAVD, path, siguiente, crear, ss, inicioAVD)

				}
				return inicioAVD

			}
		}
	} else {
		fmt.Println(colorRed, "No existe el directorio: "+conjCarpetas[siguiente])
		return -1
	}
	return -1
}

//GraficarDirectorio hace la gráfica de los avd
func GraficarDirectorio(ubicacion string, start int64, path string) {

	dir := ""
	rutas := strings.Split(ubicacion, "/")
	for i := 0; i < len(rutas)-1; i++ {
		dir += rutas[i] + "/"
	}
	nombre := rutas[len(rutas)-1]
	extension := strings.Split(nombre, ".")

	if AnalizarRuta(dir) {

		dir = dir + extension[0] + ".txt"
		var _, errr = os.Stat(dir)
		//Crea el archivo si no existe
		if os.IsNotExist(errr) {
			var file, errr = os.Create(dir)
			if existeError(errr) {
				return
			}
			defer file.Close()
		}
		auxi := false
		super = SUPERBOOT{}
		super, auxi = LeerSUPERBOOT(start, path)
		if start != -1 && auxi {
			cadena := ""
			cadena += "digraph G {\ngraph [pad=\"0.5\", nodesep=\"1\", ranksep=\"2\"];"
			cadena += "\nnode [shape=plain]\n rankdir=LR\n"
			avd, act := LeerAVD(super.SbAptrStartAVD, path, super.SbSizeStructAVD)
			if act {
				cadena += GraficarDIR(avd, "AVD", 1, path)
				cadena += UniverDir(path, 1, avd)
				cadena += "}"
				errrr := ioutil.WriteFile(dir, []byte(cadena[:]), 0644)
				if errrr != nil {
					panic(errrr)
				}
				com1 := "dot"
				com2 := "-T" + strings.ToLower(extension[1])
				com3 := dir
				com4 := "-o"
				com5 := ubicacion
				exec.Command(com1, com2, com3, com4, com5).Output()
				fmt.Println(colorGreen, "Success")
			}
		}
	}
}

//GraficarDIR reportes de las carpetas
func GraficarDIR(avd ArbolVirtualDirectorio, nodo string, cont int, path string) string {
	cadena := nodo + strconv.Itoa(cont)
	cadena += "[label=<\n<table border=\"0\" cellborder=\"1\" cellspacing=\"0\">\n"
	cadena += "<tr><td colspan=\"8\" bgcolor=\"#AF7AC5\"><i>" + "AVD" + strconv.Itoa(cont) + " \"" + string(avd.AVDFechaCreacion[:]) + "\"" + "</i></td></tr>\n"
	cadena += "<tr><td colspan=\"8\" bgcolor=\"#AF7AC5\"><i>" + Nombres(avd.AVDNameDirectoy) + "</i></td></tr>\n"
	cadena += "<tr><td bgcolor=\"#D7BDE2\">Aptr1</td><td bgcolor=\"#D7BDE2\">Aptr2</td><td bgcolor=\"#D7BDE2\">Aptr3</td><td bgcolor=\"#D7BDE2\">Aptr4</td><td bgcolor=\"#D7BDE2\">Aptr5</td><td bgcolor=\"#D7BDE2\">Aptr6</td><td bgcolor=\"#D7BDE2\">AptrDD</td><td bgcolor=\"#D7BDE2\">AptrInd</td></tr>\n"
	cadena += "<tr>\n"
	for i := 0; i < len(avd.AVDAptrArraySubdirectorios); i++ {
		if avd.AVDAptrArraySubdirectorios[i] == 0 {
			cadena += "<td> </td>"
		} else {
			cadena += "<td>" + strconv.Itoa(int(avd.AVDAptrArraySubdirectorios[i])) + "</td>\n"
		}
	}
	if avd.AVDAptrDetalleDirectorio != 0 {
		cadena += "<td>" + strconv.Itoa(int(avd.AVDAptrDetalleDirectorio)) + "</td>\n"
	} else {
		cadena += "<td> </td>\n"
	}
	if avd.AVDAptrInd != 0 {
		cadena += "<td bgcolor=\"#D2B4DE\">" + strconv.Itoa(int(avd.AVDAptrInd)) + "</td>\n"
	} else {
		cadena += "<td bgcolor=\"#D2B4DE\"> </td>\n"
	}
	cadena += "</tr>\n"
	cadena += "</table>>];\n"

	for i := 0; i < len(avd.AVDAptrArraySubdirectorios); i++ {
		if avd.AVDAptrArraySubdirectorios[i] != 0 {
			pos := super.SbAptrStartAVD + (super.SbSizeStructAVD * (avd.AVDAptrArraySubdirectorios[i] - 1))
			avdaux, ac := LeerAVD(pos, path, super.SbSizeStructAVD)
			if ac {
				cont = int(avd.AVDAptrArraySubdirectorios[i])
				cadena += GraficarDIR(avdaux, nodo, cont, path)

			}
		}
	}
	if avd.AVDAptrInd != 0 {
		pos := super.SbAptrStartAVD + (super.SbSizeStructAVD * (avd.AVDAptrInd - 1))
		avdaux, ac := LeerAVD(pos, path, super.SbSizeStructAVD)
		if ac {
			cont = int(avd.AVDAptrInd)
			cadena += GraficarDIR(avdaux, nodo, cont, path)
		}
	}
	return cadena
}

//UniverDir hace la union entre directorios
func UniverDir(path string, c int, avd ArbolVirtualDirectorio) string {
	cad := ""
	//	avd, aux := LeerAVD(super.SbAptrStartAVD, path, super.SbSizeStructAVD)

	for i := 0; i < len(avd.AVDAptrArraySubdirectorios); i++ {
		if avd.AVDAptrArraySubdirectorios[i] != 0 {
			cad += "AVD" + strconv.Itoa(c) + "->"
			cad += "AVD" + strconv.Itoa(int(avd.AVDAptrArraySubdirectorios[i])) + ";\n"
		}
	}
	if avd.AVDAptrInd != 0 {
		cad += "AVD" + strconv.Itoa(c) + "->"
		cad += "AVD" + strconv.Itoa(int(avd.AVDAptrInd)) + ";\n"
	}
	for i := 0; i < len(avd.AVDAptrArraySubdirectorios); i++ {
		if avd.AVDAptrArraySubdirectorios[i] != 0 {
			pos := super.SbAptrStartAVD + (super.SbSizeStructAVD * (avd.AVDAptrArraySubdirectorios[i] - 1))
			avdaux, ac := LeerAVD(pos, path, super.SbSizeStructAVD)
			if ac {
				cad += UniverDir(path, int(avd.AVDAptrArraySubdirectorios[i]), avdaux)
			}
		}

	}
	if avd.AVDAptrInd != 0 {
		pos := super.SbAptrStartAVD + (super.SbSizeStructAVD * (avd.AVDAptrInd - 1))
		avdaux, ac := LeerAVD(pos, path, super.SbSizeStructAVD)
		if ac {
			cad += UniverDir(path, int(avd.AVDAptrInd), avdaux)
		}
	}
	return cad
}

//GraficarBitMap grafica los bitmaps de la partición
func GraficarBitMap(rutaUbicacion string, comienzo int64, tam int64, path string) {
	dir := ""
	rutas := strings.Split(rutaUbicacion, "/")
	for i := 0; i < len(rutas)-1; i++ {
		dir += rutas[i] + "/"
	}
	if AnalizarRuta(dir) {
		var _, err = os.Stat(rutaUbicacion)
		//Crea el archivo si no existe
		if os.IsNotExist(err) {
			var file, err = os.Create(rutaUbicacion)
			if existeError(err) {
				return
			}
			defer file.Close()
		}
	}

	// Abre archivo usando permisos READ & WRITE
	var reporte, err = os.OpenFile(rutaUbicacion, os.O_RDWR, 0644)
	if existeError(err) {
		return
	}
	defer reporte.Close()

	contador := -1
	var i int64

	files, err := os.OpenFile(path, os.O_RDWR, 0644)
	defer files.Close()
	if err != nil {
		fmt.Println(colorRed, "No se encontró la ruta del archivo")
	}
	files.Seek(comienzo, 0)
	cad := ""
	var bit byte
	for i = 0; i < tam; i++ {
		contador++
		data := readNextBytes(files, int(unsafe.Sizeof(bit)))
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &bit)
		if err != nil {
			panic(err)
		}

		if contador < 20 {
			if bit == 0 {
				cad += "0 "
			} else {
				cad += string(rune(bit)) + " "
			}
		} else {
			contador = 0
			cad += "\n"
			if bit == 0 {
				cad += "0 "
			} else {
				cad += string(rune(bit)) + " "
			}
		}
	}
	_, err = reporte.WriteString(cad)
	if existeError(err) {
		return
	}
	// Salva los cambios
	err = reporte.Sync()
	if existeError(err) {
		return
	}
	fmt.Println("Se ha creado el reporte de bitmap existosamente.")
}
