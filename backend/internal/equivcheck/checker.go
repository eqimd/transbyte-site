package equivcheck

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const (
	BinariesFolder   = "binaries/"
	CommandJava      = "java"
	CommandTransbyte = BinariesFolder + "transbyte.jar"
	CommandKissat    = BinariesFolder + "kissat"
	CommandAbc       = BinariesFolder + "abc"
	CommandAigToAig  = BinariesFolder + "aigtoaig"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const randomStringLength = 8

const InnerErrorText = "Inner error"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func createJavaFileWithContent(filename, content string) error {
	f, err := os.Create(filename + ".java")
	if err != nil {
		return err
	}

	_, err = f.WriteString(content)
	if err != nil {
		return err
	}

	f.Close()

	return nil
}

func execTransalg(filename string) (string, error) {
	cmd := exec.Command(
		CommandJava,
		"-jar",
		CommandTransbyte,
		filename+".java",
		"--array-sizes",
		"5",
		"--output",
		filename+".aag",
	)
	var errb bytes.Buffer
	cmd.Stderr = &errb

	err := cmd.Run()
	return errb.String(), err
}

func execAigToAig(filename string) (string, error) {
	cmd := exec.Command(CommandAigToAig, filename+".aag", filename+".aig")
	var errb bytes.Buffer
	cmd.Stderr = &errb

	err := cmd.Run()
	return errb.String(), err
}

func execAbc(firstFn, secondFn, dir string) (string, error) {
	miterCmd := fmt.Sprintf("miter %s %s ; fraig ; write_cnf %s", firstFn+".aig", secondFn+".aig", dir+"/cnf.cnf")
	cmd := exec.Command(
		CommandAbc,
		"-q",
		miterCmd,
	)
	var errb bytes.Buffer
	cmd.Stdout = &errb

	err := cmd.Run()
	return errb.String(), err
}

func execKissat(miterDir string) (string, error) {
	cmd := exec.Command(
		CommandKissat,
		"-q",
		"--unsat",
		miterDir+"/cnf.cnf",
	)
	var outb bytes.Buffer
	cmd.Stdout = &outb

	err := cmd.Run()
	return outb.String(), err
}

func CheckEquivalence(codeFirst, codeSecond string) (string, error) {
	randDir := randStringBytes(randomStringLength)
	_ = os.Mkdir(randDir, os.ModePerm)

	firstFn := randDir + "/" + "ClassOne"
	secondFn := randDir + "/" + "ClassTwo"

	defer func ()  {
		_ = os.RemoveAll(randDir)
	}()

	if err := createJavaFileWithContent(firstFn, codeFirst); err != nil {
		log.Println("Error creating first class:", err)
		return InnerErrorText, err
	}
	if err := createJavaFileWithContent(secondFn, codeSecond); err != nil {
		log.Println("Error creating second class:", err)
		return InnerErrorText, err
	}

	if outp, err := execTransalg(firstFn); err != nil {
		log.Println("Error executing transbyte on first class:", err, outp)
		return outp, err
	}

	if outp, err := execTransalg(secondFn); err != nil {
		log.Println("Error executing transbyte on second class:", err, outp)
		return outp, err
	}

	if outp, err := execAigToAig(firstFn); err != nil {
		log.Println("Error executing aigtoaig on first aag:", err, outp)
		return outp, err
	}

	if outp, err := execAigToAig(secondFn); err != nil {
		log.Println("Error executing aigtoaig on second aag:", err, outp)
		return outp, err
	}

	if outp, err := execAbc(firstFn, secondFn, randDir); err != nil {
		log.Println("Error executing abc:", err, outp)
		return InnerErrorText, err
	}

	outp, _ := execKissat(randDir)

	if strings.Contains(outp, "UNSAT") {
		return "Equivalent", nil
	} else if strings.Contains(outp, "SAT") {
		return "Not equivalent", nil
	} else {
		return InnerErrorText, nil
	}
}
