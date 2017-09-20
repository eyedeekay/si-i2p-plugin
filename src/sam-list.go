package main

import (
        "bufio"
        "path/filepath"
        "fmt"
        "log"
        "os"
        "syscall"
)

type samList struct{
        stackOfSams []samHttp
        err error

        samAddrString string
        samPortString string

        sendPath string
        sendPipe *os.File
        sendBuff bufio.Reader

        recvPath string
        recvPipe *os.File
}

func (samStack * samList) initPipes(){
        pathConnectionExists, err := exists(filepath.Join(connectionDirectory, "parent"))
        samStack.checkErr(err)
        if ! pathConnectionExists {
                fmt.Println("Creating a connection:", "parent")
                os.Mkdir(filepath.Join(connectionDirectory, "parent"), 0755)
        }

        samStack.sendPath = filepath.Join(connectionDirectory, "parent", "send")
        pathSendExists, sendErr := exists(samStack.sendPath)
        samStack.checkErr(sendErr)
        if ! pathSendExists {
                samStack.err = syscall.Mkfifo(samStack.sendPath, 0755)
                fmt.Println("Preparing to create Pipe:", samStack.sendPath)
                samStack.checkErr(samStack.err)
                fmt.Println("checking for problems...")
                samStack.sendPipe, samStack.err = os.OpenFile(samStack.sendPath , os.O_RDWR|os.O_CREATE, 0755)
                fmt.Println("Opening the Named Pipe as a File...")
                samStack.sendBuff = *bufio.NewReader(samStack.sendPipe)
                fmt.Println("Opening the Named Pipe as a Buffer...")
                fmt.Println("Created a named Pipe for sending requests:", samStack.sendPath)
        }

        samStack.recvPath = filepath.Join(connectionDirectory, "parent", "recv")
        pathRecvExists, recvErr := exists(samStack.recvPath)
        samStack.checkErr(recvErr)
        if ! pathRecvExists {
                samStack.err = syscall.Mkfifo(samStack.recvPath, 0755)
                fmt.Println("Preparing to create Pipe:", samStack.recvPath)
                samStack.checkErr(samStack.err)
                fmt.Println("checking for problems...")
                samStack.recvPipe, samStack.err = os.OpenFile(samStack.recvPath , os.O_RDWR|os.O_CREATE, 0755)
                fmt.Println("Created a named Pipe for recieving responses:", samStack.recvPath)
        }
}

func (samStack *samList) createClient(samAddrString string, samPortString string){

}

func (samStack *samList) clientLoop(){
        //request := "i2p-projekt.i2p"
        if len(samStack.stackOfSams) != 0 {

        }else{
                samStack.stackOfSams[0].createClient(samStack.samAddrString, samStack.samPortString, "")
        }
}

func (samStack *samList) cleanupClient(){
        samStack.sendPipe.Close()
        samStack.recvPipe.Close()
        //samStack.sam.Close()
        //os.RemoveAll(filepath.Join(connectionDirectory, samStack.host))
}

func (samStack *samList) checkErr(err error) {
	if err != nil {
                samStack.cleanupClient()
		log.Fatal(err)
	}
}

func createSamList(samAddr string, samPort string) samList{
        var samStack samList
        samStack.samAddrString = samAddr
        samStack.samPortString = samPort
        return samStack
}
