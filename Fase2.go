package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

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
			break
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
					AgregarUserTXT(aux2.part.PartStart, aux.path)
				} else {
					formatear(aux2.ebr.PartSize, aux.Nombre, aux2.ebr.PartStart, aux2.fecha, aux.path, aux2.nombreMontada)
					AgregarUserTXT(aux2.ebr.PartStart, aux.path)
				}

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
	//Se meten las estructuras del AVD
	files.Seek(SB.SbAptrStartAVD, 0)
	for i = 0; i < cantAVD; i++ {
		var b5 bytes.Buffer
		AVD := ArbolVirtualDirectorio{}
		binary.Write(&b5, binary.BigEndian, &AVD)
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
func GraficarSUPERBOOT(pathi string, nombre string) {

	dir := "/home/josselyn/Escritorio/SUPERBOOT_" + nombre + ".txt"
	var _, errr = os.Stat(dir)
	//Crea el archivo si no existe
	if os.IsNotExist(errr) {
		var file, errr = os.Create(dir)
		if existeError(errr) {
			return
		}
		defer file.Close()
	}

	start, path := BuscarParticionM(nombre)
	cadena := ""
	super, auxi := LeerSUPERBOOT(start, path)
	if start != -1 && auxi != false {
		cadena := ""
		cadena += "digraph G {\ngraph [pad=\"0.5\", nodesep=\"1\", ranksep=\"2\"];"
		cadena += "\nnode [shape=plain]\nrankdir=LR;\n"
		cadena += "Tabla[label=<\n<table border=\"0\" cellborder=\"1\" cellspacing=\"0\">\n"
		cadena += "<tr><td><i>Nombre</i></td>\n<td><i>Valor</i> </td>\n</tr>"

		cadena += "<tr><td>Sb_Nombre_HD</td><td>" + Nombres(super.SbNombreHD) + "</td></tr>\n"
		cadena += "<tr><td>Sb_Arbol_Virtual_Count</td><td>" + strconv.Itoa(int(super.SbAVDcount)) + "</td></tr>\n"
		cadena += "<tr><td>Sb__Detalle_Directorio_Count</td><td>" + strconv.Itoa(int(super.SbDetalleDirectorioCount)) + "</td></tr>\n"
		cadena += "<tr><td>Sb_INodos_Count</td><td>" + strconv.Itoa(int(super.SbINodoCount)) + "</td></tr>\n"
		cadena += "<tr><td>Sb_Bloques_Count</td><td>" + strconv.Itoa(int(super.SbBloquesCount)) + "</td></tr>\n"
		cadena += "<tr><td>Sb_Arbol_Virtual_Free</td><td>" + strconv.Itoa(int(super.SbAVDFree)) + "</td></tr>\n"
		cadena += "<tr><td>Sb_Detalle_Directorio_Free</td><td>" + strconv.Itoa(int(super.SbDetalleDirectorioFree)) + "</td></tr>\n"
		cadena += "<tr><td>Sb_INodo_Free</td><td>" + strconv.Itoa(int(super.SbINodoFree)) + "</td></tr>\n"
		cadena += "<tr><td>Sb_Bloques_Free</td><td>" + strconv.Itoa(int(super.SbBloqueFree)) + "</td></tr>\n"
		cadena += "<tr><td>Sb_Date_Creation</td><td>" + string(super.SbDateCreation[:]) + "</td></tr>\n"
		cadena += "<tr><td>Sb_Date_Montaje</td><td>" + string(super.SbDateMontaje[:]) + "</td></tr>\n"
		cadena += "<tr><td>Sb_MontajesCount</td><td>" + strconv.Itoa(int(super.SbMontajesCount)) + "</td></tr>\n"
		cadena += "<tr><td>Sb_Aptr_Start_Bipmap_AVD</td><td>" + strconv.Itoa(int(super.SbAptrStartBipmapAVD)) + "</td></tr>\n"
		cadena += "<tr><td>Sb_Aptr_Start_AVD</td><td>" + strconv.Itoa(int(super.SbAptrStartAVD)) + "</td></tr>\n"
		cadena += "<tr><td>Sb_Aptr_Start_Bipmap_DD</td><td>" + strconv.Itoa(int(super.SbAptrStartBipmapDD)) + "</td></tr>\n"
		cadena += "<tr><td>Sb_Aptr_Start_DD</td><td>" + strconv.Itoa(int(super.SbAptrStartDD)) + "</td></tr>\n"
		cadena += "<tr><td>Sb_Aptr_Start_Bipmap_INodo</td><td>" + strconv.Itoa(int(super.SbAptrStartBipmapINodo)) + "</td></tr>\n"
		cadena += "<tr><td>Sb_Aptr_Start_INodo</td><td>" + strconv.Itoa(int(super.SbAptrStartINodo)) + "</td></tr>\n"
		cadena += "<tr><td>Sb_Aptr_Start_Bipmap_Bloque</td><td>" + strconv.Itoa(int(super.SbAptrStartBipmapBloque)) + "</td></tr>\n"
		cadena += "<tr><td>Sb_Aptr_Start_Bloque</td><td>" + strconv.Itoa(int(super.SbAptrStartBloque)) + "</td></tr>\n"
		cadena += "<tr><td>Sb_Aptr_Start_Bitacora</td><td>" + strconv.Itoa(int(super.SbAptrStartLogBitacora)) + "</td></tr>\n"
		cadena += "<tr><td>Sb_Size_Struct_AVD</td><td>" + strconv.Itoa(int(super.SbSizeStructAVD)) + "</td></tr>\n"
		cadena += "<tr><td>Sb_Size_Struct_DD</td><td>" + strconv.Itoa(int(super.SbSizeStructDD)) + "</td></tr>\n"
		cadena += "<tr><td>Sb_Size_Struct_Inodo</td><td>" + strconv.Itoa(int(super.SbSizeStructINodo)) + "</td></tr>\n"
		cadena += "<tr><td>Sb_Size_Struct_Bloques</td><td>" + strconv.Itoa(int(super.SbSizeStructBloque)) + "</td></tr>\n"
		cadena += "<tr><td>Sb_First_Free_Bit_AVD</td><td>" + strconv.Itoa(int(super.SbFirstBitFreeAVD)) + "</td></tr>\n"
		cadena += "<tr><td>Sb_First_Free_Bit_DD</td><td>" + strconv.Itoa(int(super.SbFirstBitFreeDD)) + "</td></tr>\n"
		cadena += "<tr><td>Sb_First_Free_Bit_INodo</td><td>" + strconv.Itoa(int(super.SbFirstBitFreeINodo)) + "</td></tr>\n"
		cadena += "<tr><td>Sb_First_Free_Bit_Bloques</td><td>" + strconv.Itoa(int(super.SbFirstBitFreeBloque)) + "</td></tr>\n"
		cadena += "<tr><td>Sb_Magic_Num</td><td>" + strconv.Itoa(int(super.SbMagicNum)) + "</td></tr>\n"
	} else {
		fmt.Println("El ID de la partición no ha sido encontrado!")
	}
	cadena += "</table>>];}"
	errrr := ioutil.WriteFile(dir, []byte(cadena[:]), 0644)
	if errrr != nil {
		panic(errrr)
	}
	com1 := "dot"
	com2 := "-Tpng"
	com3 := dir
	com4 := "-o"
	com5 := path
	exec.Command(com1, com2, com3, com4, com5).Output()
	fmt.Println(colorGreen, "Success")
}

//REP comando para reportar todos los resultados
func REP() {

}

//AgregarUserTXT se agrega el archivo a la raiz
func AgregarUserTXT(start int64, path string) {
	cadena := "1,G,root\n"
	cadena += "1,U,root,root,201602676\n"
	super, valido := LeerSUPERBOOT(start, path)
	if valido {
		file, err := os.Open(path)
		defer file.Close()
		if err != nil {
			fmt.Println(colorRed, "No se encontró la ruta del archivo")

		}
		var ocupado byte = 1
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
			avd.AVDAptrDetalleDirectorio = 2
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
			//Actualiza el bipmap del bloque
			file.Seek(super.SbAptrStartBipmapBloque, 0)
			var b5 bytes.Buffer
			binary.Write(&b5, binary.BigEndian, &ocupado)
			escribirBytes(file, b5.Bytes())
			bloque1 := Bloque{}
			bloque2 := Bloque{}

			for i := 0; i < len(cadena); i++ {
				if i < 25 {
					bloque1.BDarray[i] = cadena[i]
				} else {
					bloque2.BDarray[i] = cadena[i]
				}
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

//Mkdir analiza los comandos que vienen después
func Mkdir(arreglo []string) {
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
			if path != "" && start != -1 {
				id = strings.ToLower(comandos[1])

			} else {
				fmt.Println(colorRed, "El id de la partición no fue encontrado")
				return
			}
			break
		case "-path":
			path := direccion(comandos[1])
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
		crearAC(id, path, cadena, int64(size), crear, start, dir)
	} else {
		fmt.Println(colorRed, "Faltan parámetros obligatorios.")
	}

}

//crearAC crear archivos y carpetas
func crearAC(id string, path string, contenido string, size int64, crear bool, start int64, dir string) {

	carpetas := strings.Split(path, "/")
	for i := 0; i < len(carpetas); i++ {

	}
}

func CrearCarpeta(carpeta padre, conjCarpetas []string) {

}
