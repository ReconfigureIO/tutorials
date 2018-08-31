package main

import (
	"encoding/binary"
	"log"

	"github.com/ReconfigureIO/sdaccel/xcl"
)

func main() {
	// Allocate a 'world' for interacting with and FPGA
	world := xcl.NewWorld()
	defer world.Release()

	// Import the compiled code that will be loaded onto the FPGA (referred to here as a kernel)
	// Right now these two identifiers are hard coded as an output from the build process
	krnl := world.Import("kernel_test").GetKernel("reconfigure_io_sdaccel_builder_stub_0_1")
	defer krnl.Release()

	// Make an array to send to the FPGA for processing
	input := make([]uint32, 10)

	// Seed the array with incrementing values
	for i, _ := range input {
		input[i] = uint32(i)
	}

	// Create space in shared memory for our input array
	buff := world.Malloc(xcl.ReadOnly, uint(binary.Size(input)))
	defer buff.Free()

	// Create a variable to hold the output from the FPGA
	var output [10]uint32

	// Create space in the shared memory for the output from the FPGA
	outputBuff := world.Malloc(xcl.ReadWrite, uint(binary.Size(output)))
	defer outputBuff.Free()

	// Write our input array to the shared memory at the location we specified previously
	binary.Write(buff.Writer(), binary.LittleEndian, &input)

	// Zero out the space for the result
	binary.Write(outputBuff.Writer(), binary.LittleEndian, &output)

	// Send the location of the input array as the first argument
	krnl.SetMemoryArg(0, buff)
	// Send the location the FPGA should put the result as the second argument
	krnl.SetMemoryArg(1, outputBuff)
	// Send the length of the input array as the third argument, so the FPGA knows what to expect
	krnl.SetArg(2, uint32(len(input)))

	// Run the FPGA with the supplied arguments. This is the same for all projects.
	// The arguments ``(1, 1, 1)`` relate to x, y, z co-ordinates and correspond to our current
	// underlying technology.
	krnl.Run(1, 1, 1)

	// Read the results into our output array and then print them out
	binary.Read(outputBuff.Reader(), binary.LittleEndian, &output)

	log.Println()
	log.Println()
	log.Printf("This is our input to the FPGA: \n")
	log.Printf("Input: %v \n", input)
	log.Printf("We programmed the FPGA to multiply each integer by '2' and pass back the results as a new array. This is the result we got: \n")
	log.Printf("Output: %v \n", output)
	log.Println()
}
