package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"testing"

	"github.com/ReconfigureIO/sdaccel/xcl"
)

func main() {
	state := NewState(1000)
	defer state.Release()

	log.Println()
	log.Println()
	log.Printf("Benchmark round:")

	result := testing.Benchmark(state.Run)
	log.Println()
	log.Println()
	log.Printf("Benchmark:")
	fmt.Println(result)
}

//func BenchmarkMultiply(b *testing.B) {
//  state := NewState(100)
//  defer state.Release()

//  b.Run("100_inputs", state.Run)
//}

type State struct {
	// Kernel, InputBuf, OutputBuf, input, result.
	world      xcl.World
	program    *xcl.Program
	krnl       *xcl.Kernel
	buff       *xcl.Memory
	outputBuff *xcl.Memory
	input      []uint32
	output     []uint32
}

func NewState(nInputs int) *State {
	w := xcl.NewWorld()          // variable for new World
	p := w.Import("kernel_test") // variable to import our kernel
	size := uint(nInputs) * 4    // number of bytes needed to hold the input and output data

	s := &State{
		world:      w,
		program:    p,
		krnl:       p.GetKernel("reconfigure_io_sdaccel_builder_stub_0_1"),
		buff:       w.Malloc(xcl.ReadOnly, size),  // constructed as a function of nInputs
		outputBuff: w.Malloc(xcl.ReadWrite, size), // output will be the same size as input
		input:      make([]uint32, nInputs),       //variable to store input data
		output:     make([]uint32, nInputs),       // variable to store result data
	}

	// Seed the input array with incrementing values
	//max := int(nInputs)

	for i, _ := range s.input {
		s.input[i] = uint32(i)
	}

	// To avoid measuring warmup cost of the first few calls (especially in sim)
	const warmup = 2
	for i := 0; i < warmup; i++ {
		log.Println()
		log.Println()
		log.Printf("Warm up round:")
		s.feedFPGA()
	}

	return s
}

func (s *State) Run(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s.feedFPGA()
	}
}

func (s *State) Release() {
	s.buff.Free()
	s.outputBuff.Free()
	s.program.Release()
	s.world.Release()
}

func (s *State) feedFPGA() {
	// write input to memory
	binary.Write(s.buff.Writer(), binary.LittleEndian, &s.input)
	// Zero out the space for the result
	//binary.Write(s.outputBuff.Writer(), binary.LittleEndian, &s.output)

	// Send the location of the input array as the first argument
	s.krnl.SetMemoryArg(0, s.buff)
	// Send the location the FPGA should put the result as the second argument
	s.krnl.SetMemoryArg(1, s.outputBuff)
	// Send the length of the input array as the third argument, so the FPGA knows what to expect
	s.krnl.SetArg(2, uint32(len(s.input)))

	// start the FPGA running
	s.krnl.Run(1, 1, 1)

	// Read the results into our output variable
	binary.Read(s.outputBuff.Reader(), binary.LittleEndian, &s.output)

	log.Printf("Input: %v ", s.input)
	log.Printf("Output: %v ", s.output)
}
