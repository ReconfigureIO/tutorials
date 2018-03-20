package main

import (
  	// import the entire framework (including bundled verilog)
  	_ "github.com/ReconfigureIO/sdaccel"

  	// Use the new SMI protocol package
  	"github.com/ReconfigureIO/sdaccel/smi"
)

func Top(
    // Specify inputs and outputs to and from the FPGA. Tell the FPGA where to find data in shared memory, what data type
    // to expect or pass single integers directly to the FPGA by sending them to the control register

    ...

    // Set up ports for interacting with the shared memory, here we have 2 SMI ports which can be used to read or write
    readReq chan<- smi.Flit64,
    readResp <-chan smi.Flit64,

    writeReq chan<- smi.Flit64,
    writeResp <-chan smi.Flit64) {

    // Do whatever needs doing with the data from the host

    ...

    // Write the result to the location in shared memory as requested by the host
    smi.WriteUInt32(
		  writeReq, writeResp, <results_pointer>, smi.DefaultOptions, 512, <results_data>)
}
