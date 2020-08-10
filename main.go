package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var Comandos map[int]string
var subcomando map[int]string

func main() {

	fmt.Print("Introduzca un comando----:: ")
	reader := bufio.NewReader(os.Stdin)
	entrada, _ := reader.ReadString('\n')
	eleccion := strings.TrimRight(entrada, "\r\n")
	Analizador(eleccion + "$$")

	for eleccion != "exit" {
		fmt.Print("Introduzca un comando----:: ")
		reader = bufio.NewReader(os.Stdin)
		entrada, _ = reader.ReadString('\n')
		eleccion = strings.TrimRight(entrada, "\r\n")
		Analizador(eleccion + "$$")
	}

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

func asignacionSubcomandos() {
	subcomando = make(map[int]string)
	subcomando[0] = "path"
	subcomando[1] = "size"
	subcomando[2] = "name"
	subcomando[3] = "unit"
	subcomando[4] = "type"
	subcomando[5] = "fit"
	subcomando[6] = "delete"
	subcomando[7] = "add"
	subcomando[8] = "id"
	subcomando[9] = "usr"
	subcomando[10] = "pwd"
	subcomando[11] = "grp"
	subcomando[12] = "ugo"
	subcomando[13] = "r"
	subcomando[14] = "p"
	subcomando[15] = "cont"
	subcomando[16] = "file"
	subcomando[17] = "rf"
	subcomando[18] = "dest"
	subcomando[19] = "route"
}

func Analizador(cadena string) {
	asignacionValores()
	asignacionSubcomandos()
	estado := 0
	cadenita := ""
	lineaComando := ""
	escape := false
	ruta := false
	extension := ""
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
				estado = 5
			} else if cadena[i] == 45 {
				estado = 8
				lineaComando += caracter
			} else if cadena[i] == 92 {
				estado = 4
			} else if cadena[i] == 34 {
				estado = 5
				ruta = true
			} else if rune(cadena[i]) == '\n' || escape == false {
				AnalizarLineaComando(lineaComando)
				lineaComando = ""
			} else if rune(cadena[i]) == '\n' || escape == true {
				estado = 0
			} else if cadena[i] == 35 {
				estado = 7
				cadenita += caracter
			} else if caracter == "$" {
				AnalizarLineaComando(lineaComando)
				lineaComando = ""
			}

			break
		case 1:
			if cadena[i] >= 65 && cadena[i] <= 90 || cadena[i] >= 97 && cadena[i] <= 122 {
				cadenita += caracter
				estado = 1
			} else if cadena[i] == 47 {
				cadenita += caracter
				estado = 5
			} else if cadena[i] == 32 {
				for j := range Comandos {
					if Comandos[j] == strings.ToLower(cadenita) {
						lineaComando += cadenita + " "
						cadenita = ""
						estado = 0
						break
					}
				}
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
				lineaComando += cadenita
				cadenita = ""
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
			} else if cadena[i] == 32 || cadena[i] == '\t' {
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
			}
			break
		case 5:
			if cadena[i] == 47 {
				estado = 5
				cadenita += caracter
			} else if cadena[i] == 32 {
				cadenita += "@"

			} else if cadena[i] == 46 {
				estado = 6
				cadenita += caracter
			} else if rune(cadena[i]) != '\n' || cadena[i] != 46 {
				estado = 5
				cadenita += caracter
			}
			break
		case 6:

			if cadena[i] >= 65 && cadena[i] <= 90 || cadena[i] >= 97 && cadena[i] <= 122 {
				estado = 6
				extension += caracter
				println(extension)
			} else if cadena[i] == 34 && ruta == true {
				if strings.ToLower(extension) == "mia" {
					cadenita += extension
					estado = 0
					lineaComando += cadenita
					cadenita = ""
					extension = ""
				}
			} else if cadena[i] != 34 && ruta == false {
				if strings.ToLower(extension) == "mia" {
					println("paso por extension")
					cadenita += extension
					estado = 0
					lineaComando += cadenita
					cadenita = ""
					extension = ""
				}
			} else {
				estado = 0
			}
			break
		case 7:

			if rune(cadena[i]) != '\n' {
				cadenita = caracter
				estado = 7
			} else {
				fmt.Println(cadenita)
				cadenita = ""
			}
			break
		case 8:
			if cadena[i] >= 65 && cadena[i] <= 90 || cadena[i] >= 97 && cadena[i] <= 122 {
				cadenita += caracter
				estado = 8
			} else if cadena[i] >= 48 || cadena[i] <= 57 {
				cadenita += caracter
				estado = 3
			} else if cadena[i] == 45 {
				for j := range subcomando {
					if subcomando[j] == strings.ToLower(cadenita) {
						cadenita += string(rune(cadena[i]))
						i++
						if cadena[i] == 62 {
							cadenita += string(rune(cadena[i]))
							lineaComando += cadenita
							cadenita = ""
							estado = 0
						}
						break
					}
				}

			}
			break
		}
	}
}
func AnalizarLineaComando(cadena string) {
	fmt.Println(cadena)
}
