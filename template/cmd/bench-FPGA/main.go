package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"testing"

	"github.com/ReconfigureIO/sdaccel/xcl"
)

func BenchmarkKernel(world xcl.World, b *testing.B) {
	// Get our program
	program := world.Import("kernel_test")
	defer program.Release()

	// Get our kernel
	krnl := program.GetKernel("reconfigure_io_sdaccel_builder_stub_0_1")
	defer krnl.Release()

	// We need to create an input the size of B.N, so that the kernel
	// iterates B.N times
	input := make([]uint32, b.N)

	// create some sample input data, as an example here we're just filling the input variable with incrementing uint32s
	for i, _ := range input {
		input[i] = uint32(i)
	}

	// Create input buffer
	inputBuff := world.Malloc(xcl.ReadOnly, uint(binary.Size(input)))
	defer inputBuff.Free()

	// Create variable and buffer for the result from the FPGA, in this template we're assuming the result is the same size as the input
	result := make([]byte, b.N)
	outputBuff := world.Malloc(xcl.ReadWrite, uint(binary.Size(result)))
	defer outputBuff.Free()

	// Write input buffer
	binary.Write(inputBuff.Writer(), binary.LittleEndian, &input)

	// Set arguments – input buffer, output buffer and data length
	krnl.SetMemoryArg(0, inputBuff)
	krnl.SetMemoryArg(1, outputBuff)
	krnl.SetArg(2, uint32(len(input)))

	// Reset the timer so that we only benchmark the runtime of the FPGA
	b.ResetTimer()
	krnl.Run(1, 1, 1)
}

func main() {
	// Create the world
	world := xcl.NewWorld()
	defer world.Release()

	// Create a function that the benchmarking machinery can call
	f := func(b *testing.B) {
		BenchmarkKernel(world, b)
	}

	// Benchmark it
	result := testing.Benchmark(f)

	log.Println()
	log.Println()
	log.Printf("Benchmark: ")
	log.Println()
	// Print the benchmark result
	fmt.Printf("%s\n", result.String())
}
