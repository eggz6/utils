package checker

import "fmt"

type checkerChain struct {
	chain []func() error
}

type OK func() bool

func (c *checkerChain) Invalid(name string, ok OK) Checker {
	c.d(func() error {
		if !ok() {
			return fmt.Errorf("%v is invalid", name)
		}

		return nil
	})

	return c
}

func (c *checkerChain) Yes() error {
	for _, f := range c.chain {
		e := f()
		if e != nil {
			return e
		}
	}

	return nil
}

func (c *checkerChain) d(handle func() error) {
	c.chain = append(c.chain, handle)
}

func (c *checkerChain) NoEmptyString(name string, val string) Checker {
	c.d(func() error {
		if len(val) == 0 {
			return fmt.Errorf("empty string %s", name)
		}

		return nil
	})

	return c
}

func (c *checkerChain) NoZero(name string, val int) Checker {
	c.d(func() error {
		if val == 0 {
			return fmt.Errorf("%s is zero", name)
		}

		return nil
	})

	return c
}

type Checker interface {
	Yes() error
	NoEmptyString(name string, val string) Checker
	NoZero(name string, val int) Checker
	Invalid(name string, handle OK) Checker
}

func Check() Checker {
	return &checkerChain{chain: []func() error{}}
}
