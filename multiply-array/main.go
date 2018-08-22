package main

import (
	// Import the entire framework for interracting with SDAccel from Go (including bundled verilog)
	_ "github.com/ReconfigureIO/sdaccel"
	"github.com/ReconfigureIO/sdaccel/smi"
)

// function to multiply two uint32s
func Multiply(a uint32) uint32 {
	return a * 2
}

func Top(
	// Pass a pointer to shared memory to tell the FPGA where to find the data to be operated on,
	// and a pointer to the space in shared memory where the result should be stored. Also tell the FPGA
	// the length that the data will be.

	inputData uintptr,
	outputData uintptr,
	length uint32,

	// Set up ports for interacting with the shared memory
	readReq chan<- smi.Flit64,
	readResp <-chan smi.Flit64,

	writeReq chan<- smi.Flit64,
	writeResp <-chan smi.Flit64) {

	// Read all the input data into a channel
	inputChan := make(chan uint32)
	go smi.ReadBurstUInt32(
		readReq, readResp, inputData, smi.DefaultOptions, length, inputChan)

	// Create a channel for the result of the calculation
	transformedChan := make(chan uint32)
	// multiply each element of the input channel by 2 and send to the channel we just made to hold the result
	go func() {
		// no need to stop here, which will save us some clocks checking
		for {
			sample := <-inputChan
			val := Multiply(sample)
			transformedChan <- val
		}
	}()

	// Write transformed results back to memory
	smi.WriteBurstUInt32(
		writeReq, writeResp, outputData, smi.DefaultOptions, length, transformedChan)
}
