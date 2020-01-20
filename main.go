package main

import (
	"errors"
	"fmt"
)

func getMasterServerList() (masterservers []string) {
	for i := 1; i <= 4; i++ {
		masterserver := fmt.Sprintf("master%d.teeworlds.com:%d", i, 8283)
		masterservers = append(masterservers, masterserver)
	}
	return
}

func packControlMessageWithToken(tokenServer, tokenClient int) []byte {
	const netPacketFlagControl = 1
	const netControlMessageToken = 5
	const netTokenRequestDataSize = 512

	const size = 4 + 3 + netTokenRequestDataSize
	b := make([]byte, size, size)

	// Header
	b[0] = (netPacketFlagControl << 2) & 0b11111100
	b[3] = byte(tokenServer >> 24)
	b[4] = byte(tokenServer >> 16)
	b[5] = byte(tokenServer >> 8)
	b[6] = byte(tokenServer)
	// Data
	b[7] = netControlMessageToken
	b[8] = byte(tokenClient >> 24)
	b[9] = byte(tokenClient >> 16)
	b[10] = byte(tokenClient >> 8)
	b[11] = byte(tokenClient)

	return b
}

func unpackInt(b []byte) (result int, rest []byte, err error) {
	if len(b) < 4 {
		err = errors.New("cannot unpack, length of input insufficient")
	}
	list := b[:4]
	i := 0
	sign := int((list[i] >> 6) & 1)
	result = int(list[i] & 0b00111111)

	for {

		if list[i]&0b10000000 == 0 {
			break
		}
		i++
		result |= int((list[i] & 0b01111111)) << 6

		if list[i]&0b10000000 == 0 {
			break
		}
		i++
		result |= int((list[i] & 0b01111111)) << (6 + 7)

		if list[i]&0b10000000 == 0 {
			break
		}
		i++
		result |= int((list[i] & 0b01111111)) << (6 + 7 + 7)

		if list[i]&0b10000000 == 0 {
			break
		}
		i++
		result |= int((list[i] & 0b01111111)) << (6 + 7 + 7 + 7)
	}

	i++

	result ^= -sign
	rest = b[i:]
	return
}

func headerConnectionless(tokenServer, tokenClient int) []byte {
	const netPacketFlagConnectionless = 8
	const netPacketVersion = 1

	b := make([]byte, 9, 9)

	// Header
	b[0] = ((netPacketFlagConnectionless << 2) & 0b11111100) | (netPacketVersion & 0b00000011)
	b[1] = byte(tokenServer >> 24)
	b[2] = byte(tokenServer >> 16)
	b[3] = byte(tokenServer >> 8)
	b[4] = byte(tokenServer)
	// ResponseToken
	b[5] = byte(tokenClient >> 24)
	b[6] = byte(tokenClient >> 16)
	b[7] = byte(tokenClient >> 8)
	b[8] = byte(tokenClient)

	return b
}

func unpackControlMessageWithToken(message []byte) (tokenServer, tokenClient int, err error) {
	if len(message) < 12 {
		err = fmt.Errorf("control message is too small, %d byte, required 12 byte", len(message))
		return
	}
	tokenClient = (int(message[3]) << 24) + (int(message[4]) << 16) + (int(message[5]) << 8) + int(message[6])
	tokenServer = (int(message[8]) << 24) + (int(message[9]) << 16) + (int(message[10]) << 8) + int(message[11])
	return
}

func main() {
	ms, err := NewMasterServerFromAddress("master1.teeworlds.com:8283")
	if err != nil {
		fmt.Println("Failed to create a masterserver: ", err)
		return
	}
	defer ms.Close()

	ms.RefreshToken()
	_, err = ms.GetServerList()
	if err != nil {
		fmt.Println(err)
	}
}
