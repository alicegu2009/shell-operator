package conversion

import (
	"testing"

	. "github.com/onsi/gomega"
)

func Test_ConversionRule_FromString(t *testing.T) {
	g := NewWithT(t)

	var r string
	var c ConversionRule

	r = "asd->qwe"
	c = RuleFromString(r)
	g.Expect(c.FromVersion).Should(Equal("asd"))
	g.Expect(c.ToVersion).Should(Equal("qwe"))

	r = "unstable.example.com/asd->stable.example.com/qwe"
	c = RuleFromString(r)
	g.Expect(c.FromVersion).Should(Equal("unstable.example.com/asd"))
	g.Expect(c.ToVersion).Should(Equal("stable.example.com/qwe"))

	r = "stable.example.com/asd->v1"
	c = RuleFromString(r)
	g.Expect(c.FromVersion).Should(Equal("stable.example.com/asd"))
	g.Expect(c.ToVersion).Should(Equal("v1"))
}
