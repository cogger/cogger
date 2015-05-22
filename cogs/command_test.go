package cogs_test

import (
	"bytes"
	"os/exec"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
	"gopkg.in/cogger/cogger.v1"
	. "gopkg.in/cogger/cogger.v1/cogs"
)

var coggerInterface = reflect.TypeOf((*cogger.Cog)(nil)).Elem()

var _ = Describe("Command", func() {
	It("should execute a command", func() {
		cog := Command(context.Background(), "echo")
		Expect(reflect.TypeOf(cog).Implements(coggerInterface)).To(BeTrue())
		Expect(<-cog.Do(context.Background())).ToNot(HaveOccurred())
	})

	It("should return an error when no command", func() {
		cmdName := "randomunknowncommand"
		cog := Command(context.Background(), cmdName)
		Expect(reflect.TypeOf(cog).Implements(coggerInterface)).To(BeTrue())
		err := <-cog.Do(context.Background())
		execErr, ok := err.(*exec.Error)
		Expect(ok).To(BeTrue())
		Expect(execErr.Name).To(Equal(cmdName))
		Expect(execErr.Err).To(Equal(exec.ErrNotFound))
	})
})

var _ = Describe("CommandWithOutput", func() {
	It("should execute a command", func() {
		buf := &bytes.Buffer{}
		cog := CommandWithOutput(context.Background(), buf, "echo", "bob")
		Expect(reflect.TypeOf(cog).Implements(coggerInterface)).To(BeTrue())
		Expect(<-cog.Do(context.Background())).ToNot(HaveOccurred())
		Expect(buf.String()).To(Equal("bob\n"))
	})

	It("should return an error when no command", func() {
		cmdName := "randomunknowncommand"
		buf := &bytes.Buffer{}

		cog := CommandWithOutput(context.Background(), buf, cmdName)
		Expect(reflect.TypeOf(cog).Implements(coggerInterface)).To(BeTrue())
		err := <-cog.Do(context.Background())
		execErr, ok := err.(*exec.Error)
		Expect(ok).To(BeTrue())
		Expect(execErr.Name).To(Equal(cmdName))
		Expect(execErr.Err).To(Equal(exec.ErrNotFound))
		Expect(buf.String()).To(Equal(""))
	})
})

var _ = Describe("ExecuteCommand", func() {
	It("should execute a command", func() {
		buf := &bytes.Buffer{}
		cmd := exec.Command("echo", "bob")
		cmd.Stdout = buf
		cog := ExecuteCommand(context.Background(), cmd)

		Expect(reflect.TypeOf(cog).Implements(coggerInterface)).To(BeTrue())
		Expect(<-cog.Do(context.Background())).ToNot(HaveOccurred())
		Expect(buf.String()).To(Equal("bob\n"))
	})

	It("should return an error when no command", func() {
		cmdName := "randomunknowncommand"
		buf := &bytes.Buffer{}
		cmd := exec.Command(cmdName)
		cmd.Stdout = buf
		cog := ExecuteCommand(context.Background(), cmd)

		Expect(reflect.TypeOf(cog).Implements(coggerInterface)).To(BeTrue())
		err := <-cog.Do(context.Background())
		execErr, ok := err.(*exec.Error)
		Expect(ok).To(BeTrue())
		Expect(execErr.Name).To(Equal(cmdName))
		Expect(execErr.Err).To(Equal(exec.ErrNotFound))
		Expect(buf.String()).To(Equal(""))
	})
})
