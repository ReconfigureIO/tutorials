package main

import (
	// Import the entire framework for interracting with SDAccel from Go (including bundled verilog)
	_ "github.com/ReconfigureIO/sdaccel"

	// Use the SMI protocol package
	"github.com/ReconfigureIO/sdaccel/smi"
)

// function to multiply two uint32s
func Multiply(a uint32) uint32 {
	return a * 2
}

func Top(
	// Pass two operands to the FPGA, the integer to be multiplied and a pointer to the
	// space in shared memory where it should store the result.

	a uint32,
	addr uintptr,

	// Set up port for interacting with the shared memory
	writeReq chan<- smi.Flit64,
	writeResp <-chan smi.Flit64) {

	// Multiply incoming data by 2 using Multiply function
	val := Multiply(a)

	// Write the result to the location in shared memory as requested by the host
	smi.WriteUInt32(
		writeReq, writeResp, addr, smi.DefaultOptions, val)
}
