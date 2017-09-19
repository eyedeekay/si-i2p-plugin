package main

import (
        "fmt"
)

type samList struct{
        stackOfSams []samHttp
        samAddrString string
        samPortString string
}

func (samStack *samList) clientLoop(){
        //request := "i2p-projekt.i2p"
        if len(samStack.stackOfSams) != 0 {

        }else{
                samStack.stackOfSams[0].createClient(samStack.samAddrString, samStack.samPortString, "")
        }
}

func (samStack *samList) clientAdd(request string){
        //append(samStack.stackOfSams, newSamHttp(samList.samAddrString, samList.samPortString, request))
        //samStack.clientCheckExists(request)
}

func (samStack *samList) clientCheckExists(request string) (bool, int) {
        for samIndex, samInst := range samStack.stackOfSams {
                if samInst.hostCheck( samInst.hostSet(request) ) {
                        fmt.Println(request)
                        return true, samIndex
                }
        }
        return false, -1
}

func (samStack *samList) readInputPipe(){
        request := "http://i2p-projekt.i2p/en/docs./api/samv3"
        clientExistsAlready, clientIndex := samStack.clientCheckExists(request)
        if clientExistsAlready {
                samStack.stackOfSams[clientIndex].sendRequest(request)
        }else{
                samStack.clientAdd(request)
        }
}

func (samStack *samList) createClient(){

}

func createSamList(samAddr string, samPort string) samList{
        var samStack samList
        samStack.samAddrString = samAddr
        samStack.samPortString = samPort
        return samStack
}
