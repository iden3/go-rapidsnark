package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/iden3/go-rapidsnark"
)

var zkeyFName = flag.String("zkey", "", "circuit zkey file")
var wtnsFName = flag.String("witness", "", "witness file")
var proofFName = flag.String("proof", "proof.json", "proof file")
var publicInputsFName = flag.String("public", "public.json",
	"public inputs file")

func main() {
	flag.Parse()

	if *zkeyFName == "" {
		_, _ = fmt.Fprintf(os.Stderr, "zkey file is required\n")
		os.Exit(1)
	}
	if *wtnsFName == "" {
		_, _ = fmt.Fprintf(os.Stderr, "witness file is required\n")
		os.Exit(1)
	}
	if *proofFName == "" {
		_, _ = fmt.Fprintf(os.Stderr, "proof file is required\n")
		os.Exit(1)
	}
	if *publicInputsFName == "" {
		_, _ = fmt.Fprintf(os.Stderr, "public inputs file is required\n")
		os.Exit(1)
	}

	zkeyBytes, err := os.ReadFile(*zkeyFName)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to read zkey file: %v\n", err)
		os.Exit(1)
	}

	wtnsBytes, err := os.ReadFile(*wtnsFName)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to read witness file: %v\n", err)
		os.Exit(1)
	}

	proof, publicInputs, err := rapidsnark.Groth16Prover(zkeyBytes, wtnsBytes)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(*proofFName, []byte(proof), 0644)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(*publicInputsFName, []byte(publicInputs), 0644)
	if err != nil {
		panic(err)
	}
}
