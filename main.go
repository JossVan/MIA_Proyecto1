package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var Comandos map[int]string
var subcomando map[int]string

func main() {

	Analizador("exec -path->/home/josselyn/Documentos/Universidad/Archivos/prueba1.mia" + "$$")

}
func asignacionValores() {
	Comandos = make(map[int]string)
	Comandos[0] = "exec"
	Comandos[1] = "pause"
	Comandos[2] = "mkdisk"
	Comandos[3] = "rmdisk"
	Comandos[4] = "fdisk"
	Comandos[5] = "mount"
	Comandos[6] = "unmount"
	Comandos[7] = "mkfs"
	Comandos[8] = "login"
	Comandos[9] = "logout"
	Comandos[10] = "mkgrp"
	Comandos[11] = "rmgrp"
	Comandos[12] = "mkusr"
	Comandos[13] = "rmusr"
	Comandos[14] = "chmod"
	Comandos[15] = "mkfile"
	Comandos[16] = "cat"
	Comandos[17] = "rm"
	Comandos[18] = "edit"
	Comandos[19] = "ren"
	Comandos[20] = "mkdir"
	Comandos[21] = "cp"
	Comandos[22] = "mv"
	Comandos[23] = "find"
	Comandos[24] = "chown"
	Comandos[26] = "chgrp"
	Comandos[27] = "recovery"
	Comandos[28] = "loss"
	Comandos[29] = "rep"

}

func Analizador(cadena string) {
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
				fmt.Println(cadenita)
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
		direccion = "Comando incorrecto, se esperaba -PATH"
	}
	return direccion
}
func ValidarRuta(ruta string) bool {
	if _, err := os.Stat(ruta); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("La ruta o archivo no existe")
			return true
		} else {
			fmt.Println("Error al verificar ruta")
			return true
		}

	}
	return false
}
func AnalizarLineaComando(cadena string) {
	fmt.Println(cadena)
	arreglo := strings.Split(cadena, " ")
	if strings.ToLower(arreglo[0]) == "exec" {

		direccion := direccion(arreglo[1])
		if !ValidarRuta(direccion) {
			CargaMasiva(direccion)
		}

	}

	/*	for i := 0; i < len(arreglo); i++ {

		}*/
}
