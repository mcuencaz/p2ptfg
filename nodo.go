package main
import (
 "fmt"
 "math"
 "encoding/gob"
 "net"
 "strconv"
 // "time"
 "os"
)

// Longitud de las direcciones en hexadecimal (cantidad de cifras)
const LongDirec int = 2

// Número de vecinos hacia cada dirección a tener en cuenta en un nodo.
const leaf int = 2



// Para ejecutar el nodo, incluir en el primer argumento la dirección GUID del nodo, en el segundo argumento la dirección
// IP del nodo, y en el tercer argumento la dirección IP de un nodo ya existente en la red, para solicitar la unión.
//ej: go run nodo.go 0x1234 127.0.0.1:9000 127.0.0.1:9002


func main() {

	var inp string



	// fmt.Println("Introducir dirIp")

	if (len(os.Args) < 4){
		fmt.Println("Pocos argumentos, se ha de incluir: DirecciónGUID   DirecciónIP:puerto(Nodo actual)   DirecciónIP:puerto(Nodo al que conectarse)")

		return
	}


	// fmt.Scanln(&inp)

	var nodo Nodo


	// Direccion, err := strconv.Atoi(inp)

	// if(err!=nil){
	// 	fmt.Println(err)
	// }else{

	// 	fmt.Println("Introducir GUID para el nodo")


	// 	var guid string
	// 	fmt.Scanln(&guid)

	// 	dirtabla, err := strconv.ParseInt(guid,0,64)	

	// 	if(err!=nil){
	// 		fmt.Println(err)
	// 	}else{


	// 		nodo = newNodo(int(dirtabla), Direccion)
	// 	}
	// 	}

	argGUI, err := strconv.ParseInt(os.Args[1],0,64)
	// argDirec, err := strconv.ParseInt(os.Args[2],0,64)
	argDirec := os.Args[2]

	nodo = newNodo(int(argGUI), argDirec)

	if(err != nil){
		fmt.Println(err)
	}

	fmt.Printf("Nodo con dirección %X activo \n", argGUI)

	go server(&nodo)



	// var inp string
	// fmt.Scanln(&inp)

	// conectarse, err := strconv.Atoi(inp)

	// if(err!=nil){
	// 	fmt.Println(err)
	// }else{

	// 		// nodo2 := newNodo(0x000000, Direccion)

	// 		solicitarUnion(nodo, conectarse)
	// 	}

	// argPort, err := strconv.ParseInt(os.Args[3],0,64)
	argPort := os.Args[3]

	//#DEBUG
	// fmt.Println("IPdir del nodo al que conectarse ")
	// fmt.Println(argPort)


	solicitarUnion(nodo, argPort)






	i := 1
	for(i>0){
		
		

		fmt.Scanln(&inp)

		if(inp == "id"){
			fmt.Printf("Soy el nodo con GUID %X \n", nodo.Direccion.GUIdir)
		}

		if(inp == "t"){
			 nodo.imprimeTablas()
		}

		if(inp == "ltablas"){
			 nodo.imprimeHashVecinos()
		}

		if(inp == "hashtable"){
			 nodo.imprimeHash()
		}

		if(inp == "l"){
			nodo.imprimeVecinos()
		}

		if(inp == "r"){
			nodo.reajustarHashmap()
		}

		if(inp == "remove"){
			fmt.Println("Introducir key")
			fmt.Scanln(&inp)
			keyint, _ := strconv.ParseInt(inp,0,64)

			deleteValorHashtable(int(keyint), nodo.Direccion.IPdir)
		}

		if(inp == "add"){

			var strkey string
			var strval0 string
			var strval1 string
			var strval2 string
			fmt.Println("Introducir key")
			fmt.Scanln(&strkey)
			fmt.Println("Introducir valor")
			fmt.Scanln(&strval0, &strval1, &strval2)

			var valorhash = strval0+" "+strval1+" "+strval2
			// for _,s := range strval{
			// 	valorhash += s
			// }

			keyint, _ := strconv.ParseInt(strkey,0,64)

			addValorHashtable(int(keyint), valorhash, nodo.Direccion.IPdir)

		}



		if(inp == "get"){

			fmt.Println("Introducir key")
			fmt.Scanln(&inp)
			keyint, _ := strconv.ParseInt(inp,0,64)

			getValorHashtable(int(keyint), nodo.Direccion.IPdir, nodo.Direccion.IPdir)

		}




	}

}


type ParDireccion struct{

	GUIdir int
	IPdir string


}

type Nodo struct {
	
	Direccion ParDireccion
	//su parte del hash table
	Hashtable map[int]string
	HashVecinos [leaf*2]map[int]string

	//tablas de direcciones
	TablaDirec [LongDirec] []ParDireccion
	TablaVecinosI [leaf] ParDireccion
	TablaVecinosS [leaf] ParDireccion


}

func newParDireccion(gui int, dirIp string) ParDireccion{

	return ParDireccion{GUIdir: gui, IPdir: dirIp}
}


func newNodo(dir int, dirIp string) Nodo {

	var tablaDirec [LongDirec] []ParDireccion
	var tablaVecinosI [leaf]ParDireccion
	var tablaVecinosS [leaf]ParDireccion
	var hashVecinos [leaf*2]map[int]string
	nuevoPar := newParDireccion(dir, dirIp)
	parNulo := newParDireccion(-1,"-1")


	for i := 0; i<LongDirec; i++ {

		tablaDirec[i] = []ParDireccion{nuevoPar, nuevoPar, nuevoPar, nuevoPar, nuevoPar, nuevoPar, nuevoPar, nuevoPar, nuevoPar, nuevoPar, nuevoPar, nuevoPar, nuevoPar, nuevoPar, nuevoPar, nuevoPar} 
	}


	for k := 0; k<leaf; k++{
		tablaVecinosI[k] = parNulo
		tablaVecinosS[k] = parNulo

	}

	for k := 0; k<leaf*2; k++{
		hashVecinos[k] = make(map[int] string)
	}

	
	return Nodo{Direccion: nuevoPar, Hashtable: make(map[int] string), TablaDirec: tablaDirec, TablaVecinosI: tablaVecinosI, TablaVecinosS: tablaVecinosS, HashVecinos: hashVecinos }
}






//Funciones de los nodos

// Agrega un valor a la Hashtable interna y actualiza ese valor en las copias almacenadas en sus vecinos.
func (n *Nodo) addValor(key int , valor string ) {

	n.Hashtable[key] = valor
	n.actualizarCopiasVecinos()

}


// Obtiene el valor buscado del Hashtable, en caso de no existir se  devuelve "".
// Si no aparece el dato en la Hashtable interna, se busca en las copias de los vecinos.
func (n *Nodo) getValor(key int) string {

	res := n.Hashtable[key]

	if(res == ""){
		for _,mapa := range n.HashVecinos {
			res = mapa[key]
			if(res != ""){
				break
			}
		}
	}

	// if(res == ""){
	// 	fmt.Println("A buscar la key en otro nodo")
	// 	n.buscarDireccion(key, 1)
	// }	


	return res

}

// Se elimina el valor del Hashtable
func (n *Nodo) deleteValor(key int) {

	delete(n.Hashtable, key)
	n.actualizarCopiasVecinos()

}


//Devuelve el par GUI y dirección que debería tener acceso a la dirección indicada
func (n *Nodo) buscarDireccion(key int, num int) ParDireccion{

	if(num > LongDirec){
		fmt.Println("Debería estar en este nodo pero no está, no existe el dato")
		return newParDireccion(-1, "-1")
	}

	
	//La columna de la tabla en la que debería estar la dirección, si es 0x3422, fila 0: será 3. Si es fila 1 será 4.
	keyDiv := (key/int(math.Pow(16, float64(LongDirec-num))))
	//#DEBUG
	// fmt.Printf("Division por %X con %d \n", int(math.Pow(16, float64(LongDirec-num))), num)
	resta := ((key/int(math.Pow(16, float64(LongDirec+1-num))))*16)
	//#DEBUG
	// fmt.Printf("keyDiv vale %X menos resta %X que sale %X \n", keyDiv, resta, keyDiv-resta)
	keyDiv = keyDiv-resta
	


	Direccion := (n.Direccion.GUIdir/int(math.Pow(16, float64(LongDirec-num))))
	resta = ((n.Direccion.GUIdir/int(math.Pow(16, float64(LongDirec+1-num))))*16)
	//#DEBUG
	// fmt.Printf("Direccion vale %X menos resta %X que sale %X \n", Direccion, resta, Direccion-resta)
	Direccion = Direccion-resta

	//#DEBUG
	// fmt.Println()
	// fmt.Print("Busco : ")
	// fmt.Print(key)
	// fmt.Print("   con division: ")
	// fmt.Print(keyDiv)
	// fmt.Print("         Denom : ")
	// fmt.Print(((0x1000000)/int(math.Pow(16, float64(num)))))
	// fmt.Println("")


	if((keyDiv > 0xF) || (keyDiv < 0)){
		fmt.Println("Dirección fuera de rango")
		return newParDireccion(-1, "-1")
	}

	if(num == LongDirec){
		Direccion = -1
	}

	if(keyDiv == Direccion){
		//#DEBUG
		// fmt.Println("La tiene este nodo, busco en siguiente tabla")

		//Si es el nodo, su direcciçon estará entre la posición anterior comprobada, y la actual
		return n.buscarDireccion(key, num+1)

	}else{
		//#DEBUG
		// fmt.Printf("Se envía al nodo cuya dirección empieza con %X de la fila numero %d, direccion completa %X \n", keyDiv, num, n.TablaDirec[num-1][keyDiv])
		// fmt.Println(keyDiv)

		// if(num == LongDirec){

			return n.TablaDirec[num-1][keyDiv]

		// }else{
		// 	return newParDireccion(-2, "-2")
		// }

	}

	return newParDireccion(-3, "-3")
}



// Función que modifica las tablas de vecinos, para comprobar si debe añadirse un nuevo nodo.
func (n *Nodo) editarVecinos(vecinoDir ParDireccion, avisar bool){


	gui := vecinoDir.GUIdir
	direcIp := vecinoDir.IPdir

	guiNodoN := n.Direccion.GUIdir

	//variable empleada para introducir un valor mayor que la dirección del nodo, en la tabla inferior, en caso de que sea el más próximo (circular).
	noPonerMenor := true
	var j int


	// Closure para introducir nodo en la lista de vecinos superiores
	ponerNodoEnSuperior := func(k int) {
		//Poner el nodo
					// #DEBUG
					// fmt.Printf("Se pone %X como vecino \n", gui)

					//si hay un nodo al final de la lista recalculo si debería estar al acabar la función.
					if(n.TablaVecinosS[leaf-1].GUIdir>=0){
						ultNodo := n.TablaVecinosS[leaf-1]

						//recalculo el antiguo vecino, por si se pasa a la otra tabla
						defer n.editarVecinos(ultNodo, true)
					}

					j = leaf-1
						for j-1 >= k {
							n.TablaVecinosS[j] = n.TablaVecinosS[j-1]
							j--
						}

						n.TablaVecinosS[k] = newParDireccion(gui, direcIp)
						pedirHashvecino(newParDireccion(gui, direcIp), n.Direccion.IPdir)
						n.reajustarHashmap()
						return
	}

	// Closure para introducir nodo en la lista de vecinos inferiores
	ponerNodoEnInferior := func(i int){

		//Se sustituye el nodo
					// #DEBUG
					// fmt.Printf("Se pone %X como vecino \n", gui)

					//si hay un nodo al "final" (principio) de la lista y recalculo si debería estar al acabar la función.
					if(n.TablaVecinosI[0].GUIdir>=0){
						ultNodo := n.TablaVecinosI[0]
						// recalculo el antiguo vecino, por si se pasa a la otra tabla
						defer n.editarVecinos(ultNodo, true)
					}

					j = 0
						for j+1 <= i {
							n.TablaVecinosI[j] = n.TablaVecinosI[j+1]
							j++
						}

						n.TablaVecinosI[i] = newParDireccion(gui, direcIp)
						pedirHashvecino(newParDireccion(gui, direcIp), n.Direccion.IPdir)
						n.reajustarHashmap()
						return

	}

	// Si el vecino que se mira tiene la misma dirección que el nodo, se deja de buscar (no puede ser vecino de sí mismo).
	if(gui == guiNodoN){
		return
	}

	
	//Si el nuevo es menor que el nodo:
	if(gui < guiNodoN){
			//comprobar si iria en la lista menor, si puede ir poner noPonerMenor a false: casos de ser mayor que alguno de la lista menor
			// o que alguno de la lista menor sea mayor que el nodo
			for i:=leaf-1; i>=0; i-- {

				//si ya está en la lista acaba
					if(n.TablaVecinosI[i].GUIdir == gui){ return }

				// Si el nuevo vecino es mayor que el existente (más próximo por la izquierda) o hay algun nodo superior que sea más próximo por
				// distancia circular
				if(gui>n.TablaVecinosI[i].GUIdir || n.TablaVecinosI[i].GUIdir>guiNodoN){
					
					noPonerMenor = false
				}
			}
		}




	//Comprobar lista mayor
	if noPonerMenor{
		for k:=0; k<leaf; k++ {
			if(n.TablaVecinosS[k].GUIdir == gui){ return }
		//Si el nuevo es menor que el nodo:
			if(gui < guiNodoN){
				//Si el nuevo es mayor que el nodo antiguo, que también será menor que el nodo actual (Si es nodo = 6, en la lista superior está el 4, y pongo un 3)
				// n.TablaVecinosS[k].GUIdir<0 implica que está vacío
				if(n.TablaVecinosS[k].GUIdir < guiNodoN){
					if((gui < n.TablaVecinosS[k].GUIdir || n.TablaVecinosS[k].GUIdir<0)){
						//Poner el nodo
						ponerNodoEnSuperior(k)
						n.reajustarHashmap()
						return
					}
				}

			}else if( gui > guiNodoN) {
				//Si el nodo es mayor que el nodo actual, y es menor que el nodo antiguo, o el nodo antiguo es menor que el nodo actual
				if ((n.TablaVecinosS[k].GUIdir < guiNodoN) || (gui < n.TablaVecinosS[k].GUIdir)){
					//Poner el nodo
					ponerNodoEnSuperior(k)
					return

				}


			}
		}
	}


	//Comprobar lista menor.
	for i:=leaf-1; i>=0; i-- {
		if(n.TablaVecinosI[i].GUIdir == gui){ return }
			// Si el nodo nuevo es menor que el nodo actual, se comprueba solo si el nodo nuevo es mayor que el antiguo (está más cerca del nodo n)
			
			//Si el nodo antiguo es mayor que el nodo, (en el círculo, está más cerca por la izquierda del 0), se sustituye si 
			// el nodo nuevo es mayor que el antiguo, o si el nodo nuevo es menor que la dirección del nodo
				//<0 indica si es -1, si está vacío
			if(n.TablaVecinosI[i].GUIdir > guiNodoN || n.TablaVecinosI[i].GUIdir<0 ){
				if( (gui > n.TablaVecinosI[i].GUIdir) || (gui < guiNodoN)){
					//Se pone el nodo
					ponerNodoEnInferior(i)
					return

				}
			}else if(gui < guiNodoN){
				if(gui>n.TablaVecinosI[i].GUIdir){
					//Se pone el nodo
					ponerNodoEnInferior(i)
					return

				}
			}
		}



}



//Imprime las tablas de dirección
func (n *Nodo) imprimeTablas() {


	for index, tabla := range n.TablaDirec {
		fmt.Printf("Tabla %d : \n", index)
		for ind2,pardir := range tabla {

			fmt.Printf("%X: %X %s| ", ind2, pardir.GUIdir, pardir.IPdir)
		}
		fmt.Printf("\n")
		
	}

}



func (n *Nodo) imprimeHash() {

	fmt.Println(n.Hashtable)
	
}

func (n *Nodo) imprimeHashVecinos() {

	for i, hashvecino := range n.HashVecinos{
		if(i<leaf){
			fmt.Printf("Tabla de %X ", n.TablaVecinosI[i].GUIdir)
			fmt.Println(hashvecino)
		}else if(i>=leaf){
			fmt.Printf("Tabla de %X ", n.TablaVecinosS[i-leaf].GUIdir)
			fmt.Println(hashvecino)
		}

	}
	
}


func (n *Nodo) imprimeVecinos() {
	for _,vecino := range n.TablaVecinosI{
		fmt.Printf("%X ", vecino.GUIdir)
	}

	fmt.Printf("  ")

	for _,vecino := range n.TablaVecinosS{
		fmt.Printf("%X ", vecino.GUIdir)
	}

	fmt.Printf("\n")
}


// Se envían la tabla Hashstable a los vecinos, para que guarden una copia
func (n *Nodo) actualizarCopiasVecinos(){

	for _,vecino := range n.TablaVecinosI{
		if(vecino.GUIdir>=0){
			enviarHashvecino(n.Hashtable, n.Direccion, vecino.IPdir)
		}
	}

	for _,vecino := range n.TablaVecinosS{
		if(vecino.GUIdir>=0){
			enviarHashvecino(n.Hashtable, n.Direccion, vecino.IPdir)
		}
	}

}

// AL incorporarse nuevos vecinos, se comprueba si algún valor del nodo debería estar almacenado por el nuevo vecino 
// en lugar de por el nodo.
func (n *Nodo) reajustarHashmap(){

	for key, valor := range n.Hashtable {

		direccionCorrecta := n.buscarDireccion(key, 1)
		//#DEBUG
		// fmt.Printf("Direccion correcta es: %X \n", direccionCorrecta.GUIdir)
		if( direccionCorrecta.GUIdir != n.Direccion.GUIdir){

			addValorHashtable(key, valor, direccionCorrecta.IPdir)

			//Para que fuera seguro, habría que borrar el valor al recibir confirmación de que el otro nodo ha añadido el valor
			n.deleteValor(key)
		}


	}
	
	n.actualizarCopiasVecinos()

}




// Realiza (o no) las modificaciones necesarias a la tabla Pastry "TablaDirec" para añadir a un nodo nuevo.
func (n *Nodo) editarTablaR(nuevaDir ParDireccion) {
	tabla := 0
	for tabla < LongDirec {
		var direcTotal = int(math.Pow(0x10, float64((LongDirec-1)-tabla)))
		var comienzo = 0
		var comienzoNuevo = 0
		var noCambiado = true


		


		//comprueba si empieza por el número que debe para editar las tablas, si no no hace nada
		if(tabla != 0){
			// fmt.Printf("Resto %X con %X para ver si empiezan igual, debería ser menor que %X \n", n.Direccion.GUIdir, nuevaDir.GUIdir, int(math.Pow(0x10, float64(LongDirec-tabla))) )
			// fmt.Println(math.Abs(float64((n.Direccion.GUIdir - nuevaDir.GUIdir))) >= math.Pow(0x10, float64(LongDirec-tabla)))
			// fmt.Printf("Resultado %X \n", int(math.Abs(float64((n.Direccion.GUIdir - nuevaDir.GUIdir)))))
			// if( math.Abs(float64((n.Direccion.GUIdir - nuevaDir.GUIdir))) >= math.Pow(0x10, float64(LongDirec-tabla)) ){
			// 	return
			// }
			
			// //Se eliminan los números no necesarios para la comprobación de la tabla 
			// comienzo = n.Direccion.GUIdir/int(math.Pow(0x10, float64(LongDirec-tabla)))
			// comienzoNuevo := nuevaDir.GUIdir/int(math.Pow(0x10, float64(LongDirec-tabla)))

			//Se cogen los prefijos de las direcciones para comprobación de comienzos.
			comienzo = n.Direccion.GUIdir/int(math.Pow(0x10, float64(LongDirec-tabla)))
			comienzoNuevo = nuevaDir.GUIdir/int(math.Pow(0x10, float64(LongDirec-tabla)))

			//#DEBUG
			// fmt.Printf("La comparación de comienzos es entre %X y %X \n", comienzo, comienzoNuevo)


			//Se eliminan los números no necesarios para la comprobación de la tabla, si el comienzo de las direcciones no es el mismo 
			// no se debe incluir en las filas que no son la primera, pues para ello deberían compartir el prefijo.
			
			if( comienzo!=comienzoNuevo){
				return
			}

			comienzo = comienzo * int(math.Pow(0x10, float64(LongDirec-tabla)))


		}

		for i,columna := range n.TablaDirec[tabla] {


			// Ahora mismo lo cambia si es igual o mayor a la Direccion inicial del rango 
			dif1 := columna.GUIdir - ((direcTotal*i)+comienzo)//+0x0FFFFF)
			dif2 := nuevaDir.GUIdir - ((direcTotal*i)+comienzo)//+0x0FFFFF)

			//#DEBUG
			// fmt.Printf("Comparo %X con %X para la dirección %X \n", columna.GUIdir, nuevaDir.GUIdir, direcTotal*i+comienzo )

			//SI lo cambio tengo que avisar a mi lista (Al final no)


			//Se debe poner en la tabla 0 con más preferencia a un nodo que comience por esa dirección que uno que no.
			// Para la primera fila, se coloca el nuevo nodo si su dirección comienza por lo mismo que la columna, aunque esté más lejano al comienzo del rango
			if(tabla < LongDirec-1){

				comienzoNuevo = nuevaDir.GUIdir/int(math.Pow(0x10, float64(LongDirec-1-tabla)))
				comienzoYaEnColumna := n.TablaDirec[tabla][i].GUIdir/int(math.Pow(0x10, float64(LongDirec-1-tabla)))

				comienzoPref := n.Direccion.GUIdir/int(math.Pow(0x10, float64(LongDirec-tabla)))

				// Si se cumple esta condicion significa que el nuevo nodo comienza por la direccion correspondiente al indice de la columna
				// y el que había no lo hacía
				prefijoMirado := comienzoPref*0x10 + i


				// fmt.Printf("Comparo %X con %X para la dirección %X y sale \n", comienzoNuevo, comienzoYaEnColumna, prefijoMirado )

				if(comienzoNuevo == prefijoMirado && comienzoYaEnColumna != prefijoMirado){
					n.TablaDirec[tabla][i].GUIdir = nuevaDir.GUIdir
					n.TablaDirec[tabla][i].IPdir = nuevaDir.IPdir
					noCambiado = false

				}

				// fmt.Println(!noCambiado)
			}



			// Si está más cerca de la dirección lo cambio, y en caso de estar a la misma distancia, pongo la direccion mas pequeña
			// Si es la primera tabla y ya se ha cambiado, no se hace nada
			if(noCambiado && (math.Abs(float64(dif2))<=math.Abs(float64(dif1)))){
				if((math.Abs(float64(dif2))==math.Abs(float64(dif1)))){
						if(nuevaDir.GUIdir<columna.GUIdir){
							n.TablaDirec[tabla][i].GUIdir = nuevaDir.GUIdir
							n.TablaDirec[tabla][i].IPdir = nuevaDir.IPdir
						}
					}else{
								n.TablaDirec[tabla][i].GUIdir = nuevaDir.GUIdir
								n.TablaDirec[tabla][i].IPdir = nuevaDir.IPdir
							}
			}

			noCambiado = true
		}

		tabla++
	}

}


// FUNCIONES DE RED SERVIDOR

type Transmision struct {

	Orden string
	Nodo Nodo
	NumTabla int

	//funciones hashtable
	Hashkey int
	Hashvalor string

	Hashtable map[int]string

}


func server(nodo *Nodo) {

	// port := strconv.Itoa(nodo.IPdir)

	 // listen on a port
	 ln, err := net.Listen("tcp", /*":"+port*/ nodo.Direccion.IPdir )
	 if err != nil {
		 fmt.Println(err)
		 return
	 }


	for {
		 // accept a connection
		 c, err := ln.Accept()
		 if err != nil {
			 fmt.Println(err)
		 continue
	 }
	 	// handle the connection
		 go handleServerConnection(c, nodo)
	 }
}


func handleServerConnection(c net.Conn, nodo *Nodo) {
	 // receive the message
	 var trans Transmision
	 err := gob.NewDecoder(c).Decode(&trans)
	 if err != nil {
		 fmt.Println(err)
	 } else {
	 	//#DEBUG
	 	// fmt.Println("Orden a ejecutar: ", trans.Orden)
	 }


	 // Ejecutado en el nodo nuevo, cuando el nodo que debe añadirle a la lista le envía sus tablas.
	 // Se actualizan las tablas propias, y se envia mensaje de unir 
	 if(trans.Orden == "editarTabla"){

	 	var anunciados map[int]int

	 	for i := 0; i<LongDirec; i++ {
	 		anunciados = make(map[int]int)
	 		// nodo.editarTablaR(trans.Nodo.Direccion, i)



	 		for _,direccionTabla := range trans.Nodo.TablaDirec[i] {
	 			//Edita su tabla para incluir todas las direcciones de la tabla que ha recibido, en caso de que debiera incluirlas,
	 			//y envia a esas direcciones su direccion, para que estas le tengan en cuenta en caso de tener que hacerlo

		   		if((direccionTabla.GUIdir != nodo.Direccion.GUIdir) && (anunciados[direccionTabla.GUIdir]!=1)){

		   			//#DEBUG
		   			// fmt.Printf("Soy %X, y voy a agregar a %X \n",nodo.Direccion.GUIdir, direccionTabla.GUIdir)

		   			nodo.editarTablaR(direccionTabla)
		   			nodo2 := *nodo
		   			anunciados[direccionTabla.GUIdir] = 1
		   			agregarNodo(nodo2, direccionTabla.IPdir, i)
		   		}
		   }


			for _,vecino := range trans.Nodo.TablaVecinosI {
	 			if(vecino.GUIdir>=0){
		 			enviarVecino(nodo.Direccion, vecino.IPdir)
		 			nodo.editarVecinos(vecino, false)
		 		}
	 		}

	 		for _,vecino := range trans.Nodo.TablaVecinosS {
	 			if(vecino.GUIdir>=0){
		 			enviarVecino(nodo.Direccion, vecino.IPdir)
		 			nodo.editarVecinos(vecino, false)
		 		}
	 		}

	 		//pone como vecino a quien le ha enviado las tablas
	 		nodo.editarVecinos(trans.Nodo.Direccion, false)
	 		enviarVecino(nodo.Direccion, trans.Nodo.Direccion.IPdir)



	 	}

	}


	// Recibido en nodos ya existentes de la red, por el nodo que quiere unirse, o como 
	// retransmisión de un nodo que ha recibido lo mismo, porque este está más cerca de la dirección que debe.
	if(trans.Orden == "solicitarUnion"){
		dirAnunciar := nodo.buscarDireccion(trans.Nodo.Direccion.GUIdir, 1)


		if(dirAnunciar.GUIdir == nodo.Direccion.GUIdir){

			nodo2 := *nodo
			enviarTablaR(nodo2, trans.Nodo.Direccion.IPdir, 0)
		}else{
			solicitarUnion(trans.Nodo, dirAnunciar.IPdir)
		}



	}

	// Recibido por un nodo que aparece en las tablas de un nodo recién registrado, 
	if(trans.Orden == "agregarNodo"){

		//#DEBUG
		fmt.Printf("Agregando nodo: soy %X y añado a %X  \n", nodo.Direccion.GUIdir, trans.Nodo.Direccion.GUIdir)
		nodo.editarTablaR(trans.Nodo.Direccion)
		// nodo.editarVecinos(trans.Nodo.Direccion, true)

	}



	// Funciones hash

	// Recibe orden de añadir un valor de la hashtable, y busca a quien debe añadirlo. Si es el nodo lo añade.
	if(trans.Orden == "addValorHashtable"){

		dirValor := nodo.buscarDireccion(trans.Hashkey, 1)

		if(dirValor.GUIdir == nodo.Direccion.GUIdir){
			nodo.addValor(trans.Hashkey, trans.Hashvalor)
		}else{
			addValorHashtable(trans.Hashkey, trans.Hashvalor, dirValor.IPdir)
		}

	}


	// Recibe orden de eliminar un valor de la hashtable, y busca quien debe eliminarlo. Si es el nodo lo elimina.
	if(trans.Orden == "deleteValorHashtable"){

		dirValor := nodo.buscarDireccion(trans.Hashkey, 1)

		if(dirValor.GUIdir == nodo.Direccion.GUIdir){
			nodo.deleteValor(trans.Hashkey)
		}else{
			deleteValorHashtable(trans.Hashkey, dirValor.IPdir)
		}

	}

	// Recibe orden de obtener un valor de la hashtable, y busca a quien lo tiene. Si es el nodo que debería tenerlo pero no lo tiene,
 	// imprime que no existe, si lo tiene lo envía al nodo solicitante

	if(trans.Orden == "getValorHashtable"){

		

		res := nodo.getValor(trans.Hashkey)

		if(res == ""){

			//#DEBUG
			// fmt.Println("A buscar la key en otro nodo")

			dirValor := nodo.buscarDireccion(trans.Hashkey, 1)

			if(dirValor.GUIdir == nodo.Direccion.GUIdir){
				fmt.Println("No existe ese dato")
				enviarValor("[none]", trans.Nodo.Direccion.IPdir)
			}else{
				getValorHashtable(trans.Hashkey, dirValor.IPdir, trans.Nodo.Direccion.IPdir)
			}


		}else{
			enviarValor(res, trans.Nodo.Direccion.IPdir)
		}
		
	}


	// Recibido por un nodo que ha solicitado un dato, imprime el valor del dato.
	if(trans.Orden == "enviarValor"){

		if trans.Hashvalor == "[none]" {
			fmt.Println("No hay ningún valor asociado a esa clave")
		}else{ 
			fmt.Println("El valor asociado a la clave es:")
			fmt.Println(trans.Hashvalor)
		}

	}


	//Recibido por un nodo que quiere que compruebe si debo añadirle como vecino
	if(trans.Orden == "enviarVecino"){

		//#DEBUG
		// fmt.Println("Se envia un vecino")
		nodo.editarVecinos(trans.Nodo.Direccion, true)

		nodo.reajustarHashmap()
		nodo.actualizarCopiasVecinos()

		fmt.Println(trans.Hashvalor)


	}


	//Un vecino ha pedido el hashtable
	if(trans.Orden == "pedirHashvecino"){

		//#DEBUG
		// fmt.Println("Se solicita el hashmap por un vecino")
		enviarHashvecino(nodo.Hashtable, nodo.Direccion, trans.Nodo.Direccion.IPdir)

	}


	//Se recibe un hashtable de un vecino
	if(trans.Orden == "enviarHashvecino"){

		guiOrigen := trans.Nodo.Direccion.GUIdir

		//#DEBUG
		// fmt.Printf("Se recibe el hashmap de %X \n", guiOrigen)

		var i = 0
		for _,vecino := range nodo.TablaVecinosI{

			if(guiOrigen == vecino.GUIdir){
				nodo.HashVecinos[i] = trans.Hashtable
				break
			}
			i++

		}

		for _,vecino := range nodo.TablaVecinosS{

			if(guiOrigen == vecino.GUIdir){
				nodo.HashVecinos[i] = trans.Hashtable
				break
			}
			i++
		}



	}



	 c.Close()
}


// Se envian las tablas a un nodo que ha solicitado unirse
func enviarTablaR(nodo Nodo, destino string, numTabla int){


	// port := strconv.Itoa(destino)
	//#DEBUG
	// fmt.Println("La Direccion a la que se quiere enviar es "+destino)

	c, err := net.Dial("tcp", /*"127.0.0.1:"+port*/destino)
	 if err != nil {
	 	fmt.Println(err)
	 	return
	 }


	 transmision := newTransmision("editarTabla", nodo)
	 transmision.NumTabla = numTabla


	 // fmt.Println("Recibida orden de: ",transmision.Orden)
	 err = gob.NewEncoder(c).Encode(transmision)
	 if err != nil {
	 	fmt.Println(err)
	 }
	 c.Close()

}


// Se pide una tabla de direcciones
// No se usa
func pedirTablaR(nodo Nodo, destino string, numTabla int){


	// port := strconv.Itoa(destino)


	c, err := net.Dial("tcp", /*"127.0.0.1:"+port*/destino)
	 if err != nil {
	 	fmt.Println(err)
	 	return
	 }

	 transmision := newTransmision("pedirTabla", nodo)
	 transmision.NumTabla = numTabla

	 //#DEBUG
	 // fmt.Println("Solicitando tabla a ", destino)
	 // fmt.Println("Recibida orden de: ",transmision.Orden)
	 err = gob.NewEncoder(c).Encode(transmision)
	 if err != nil {
	 	fmt.Println(err)
	 }
	 c.Close()

}

// Constructor de estructura transmision
func newTransmision(orden string, nodo Nodo) Transmision {

	return Transmision{Orden: orden, Nodo: nodo, NumTabla: 0 }
}


// Nodo nuevo lo ejecuta para solicitar una unión a la red, 
// o es reenviado por otro nodo hacia el que debería añadir el nuevo.
func solicitarUnion(nodo Nodo, destino string) {

	// port := strconv.Itoa(destino)



	c, err := net.Dial("tcp", /*"127.0.0.1:"+port*/destino)
	 if err != nil {
	 	fmt.Println(err)
	 	return
	 }

	 transmision := newTransmision("solicitarUnion", nodo)


	 err = gob.NewEncoder(c).Encode(transmision)
	 if err != nil {
	 	fmt.Println(err)
	 }
	 c.Close()


}

// Enviado por los nodos recién llegados, para que se les tenga en cuenta en las tablas.
func agregarNodo(nodo Nodo, destino string, numTabla int) {

	// port := strconv.Itoa(destino)


	c, err := net.Dial("tcp", /*"127.0.0.1:"+port*/destino)
	 if err != nil {
	 	fmt.Println(err)
	 	return
	 }
	 // send the message

	 transmision := newTransmision("agregarNodo", nodo)
	 transmision.NumTabla = numTabla


	 err = gob.NewEncoder(c).Encode(transmision)
	 if err != nil {
	 	fmt.Println(err)
	 }
	 c.Close()


}




// Funciones Hashtable

// Utilizado en un nodo, para mandar el mensaje a todos los demás de que se añada un nuevo valor a la tabla.
func addValorHashtable(key int, valor string, destino string){

	// port := strconv.Itoa(destino)


	c, err := net.Dial("tcp", /*"127.0.0.1:"+port*/destino)
	 if err != nil {
	 	fmt.Println(err)
	 	return
	 }
	 // send the message

	 nodo := newNodo(0x000000, destino)
	 transmision := newTransmision("addValorHashtable", nodo)
	 transmision.Hashkey = key
	 transmision.Hashvalor = valor


	 err = gob.NewEncoder(c).Encode(transmision)
	 if err != nil {
	 	fmt.Println(err)
	 }
	 c.Close()


}


// Se anula un valor de la tabla, poniéndolo a valor vacío
func deleteValorHashtable(key int, destino string){

	// port := strconv.Itoa(destino)


	c, err := net.Dial("tcp", destino)
	 if err != nil {
	 	fmt.Println(err)
	 	return
	 }
	 // send the message

	 nodo := newNodo(0x000000, destino)
	 transmision := newTransmision("deleteValorHashtable", nodo)
	 transmision.Hashkey = key


	 err = gob.NewEncoder(c).Encode(transmision)
	 if err != nil {
	 	fmt.Println(err)
	 }
	 c.Close()


}



// Ejecutado en un nodo, para mandar el mensaje al que deba de que se pide un valor de la tabla.
func getValorHashtable(key int, destino string, origen string){

	// port := strconv.Itoa(destino)


	c, err := net.Dial("tcp", /*"127.0.0.1:"+port*/destino)
	 if err != nil {
	 	fmt.Println(err)
	 	return
	 }
	 // send the message
	
	nodo := newNodo(0x000000, origen)


	 transmision := newTransmision("getValorHashtable", nodo)
	 transmision.Hashkey = key


	 err = gob.NewEncoder(c).Encode(transmision)
	 if err != nil {
	 	fmt.Println(err)
	 }
	 c.Close()


}


// Se envía un mensaje a una dirección.
func enviarValor(valor string, destino string){

	// port := strconv.Itoa(destino)

	c, err := net.Dial("tcp", /*"127.0.0.1:"+port*/destino)
	 if err != nil {
	 	fmt.Println(err)
	 	return
	 }

	nodo := newNodo(0x000000, destino)

	transmision := newTransmision("enviarValor", nodo)
	transmision.Hashvalor = valor

	err = gob.NewEncoder(c).Encode(transmision)
	 if err != nil {
	 	fmt.Println(err)
	 }
	 c.Close()

}


// Avisa a un nodo para que compruebe si debe añadirle como vecino
func enviarVecino(vecino ParDireccion, destino string){

	// port := strconv.Itoa(destino)

	c, err := net.Dial("tcp", destino)
	 if err != nil {
	 	fmt.Println(err)
	 	return
	 }

	nodo := newNodo(vecino.GUIdir, vecino.IPdir)

	transmision := newTransmision("enviarVecino", nodo)

	err = gob.NewEncoder(c).Encode(transmision)
	 if err != nil {
	 	fmt.Println(err)
	 }
	 c.Close()



}


func pedirHashvecino(vecino ParDireccion, origen string){

	c, err := net.Dial("tcp", vecino.IPdir)
	 if err != nil {
	 	fmt.Println(err)
	 	return
	 }

	nodo := newNodo(0x0, origen)



	transmision := newTransmision("pedirHashvecino", nodo)

	err = gob.NewEncoder(c).Encode(transmision)
	 if err != nil {
	 	fmt.Println(err)
	 }
	 c.Close()

}


func enviarHashvecino(hashmap map[int]string, origen ParDireccion, destino string){

	c, err := net.Dial("tcp", destino)
	 if err != nil {
	 	fmt.Println(err)
	 	return
	 }

	nodo := newNodo(origen.GUIdir, origen.IPdir)



	transmision := newTransmision("enviarHashvecino", nodo)
	transmision.Hashtable = hashmap

	err = gob.NewEncoder(c).Encode(transmision)
	 if err != nil {
	 	fmt.Println(err)
	 }
	 c.Close()


}


