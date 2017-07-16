# p2ptfg

Requiere instalar go desde la página https://golang.org/

Para ejecutar el código se tienen que indicar la dirección hexadecimal interna (GUID) que se quiere dar interna al nodo, comenzando por "0x", seguido de las direcciones IP:puerto a emplear por el nodo, y la dirección IP:puerto a la que conectarse.
Por ejemplo, para ejecutar un nodo de dirección GUID 0x10, ejecutado con la dirección local y puerto 9000 conectado a sí mismo:

go run nodo.go 0x10 127.0.0.1:9000 127.0.0.1:9000

Si se quiere ejecutar un nodo de GUID 0x50 conectado al anterior:

go run nodo.go 0x50 127.0.0.1:9001 127.0.0.1:9000



Para modificar el número de nodos que caben en la red, se puede modificar la variable "LongDirec".
Los ejemplos anteriores están escritos para LongDirec = 2, con un valor de 6 las direcciones serían 0x100000 y 0x500000.

Con la variable "leaf" se pueden modificar la cantidad de vecinos hacia cada dirección. Con un valor leaf = 2 cada nodo almacena 2 vecinos superiores y 2 vecinos inferiores.



Los comandos que se pueden ejecutar en el nodo son:

id :Imprime por pantalla la dirección GUID que posee el nodo en el que se ejecuta.

t :Imprime las tablas de direccionamiento Pastry que contiene el nodo.

ltablas :imprime las tablas hash de los vecinos del nodo, que tiene almacenadas en caso de fallo de estos.

hashtable :imprime los valores de la tabla hash interna que contiene el nodo

l :Imprime una tabla con las direcciones GUID de los nodos vecinos.

add :Se indica que se quiere añadir un par valor-clave a la tabla hash global del sistema distribuido, se solicitarán posteriormente el valor de la clave, que debe ser hexadecimal (comenzando por 0x), para después pedir el valor que se quiere asociar.

addh :Al igual que add añade un valor, pero pudiendo establecer un valor de tipo string como clave.

get :Orden para obtener el valor asociado a una clave de la tabla hash global, tras ejecutar get se pedirá el valor hexadecimal de la clave a buscar.

geth :Al igual que get obtiene un valor, pero pudiendo establecer una clave de tipo string.

remove :Orden para eliminar un par valor-clave de la tabla hash global del sistema.

r :Pide al nodo que recompruebe su tabla hash interna por si alguno de sus valores debería estar asociado al vecino (No tiene utilidad práctica puesto que esta acción se hará de forma automática, este comando solo tiene utilidad para pruebas).



############################################
El archivo nodo128 contiene las mismas funciones que el anterior, pero permitiendo un mayor número en la variable de longitud de direcciones.

- Para los parámetros de las funciones a ejecutar, como add y get, no se debe añadir el prefijo 0x a la hora de tratar con valores hexadecimales.
