package main

import (
	"fmt"
	"strings"
	"time"
	"unsafe"
)

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
					//CalculosPrimaria(aux2.part)
				} else {

				}
			}
			aux2 = aux2.siguiente
		}
		aux = aux.siguiente
	}
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

func formatear(size int64, NombreDisco string, inicioParticion int64) {
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
	SB.SbMagicNum = 201602676
}
