package counter

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestSingle(t *testing.T) {
	c := New()

	for i := 1; i < 100; i++ {
		c.Increment("foo")
		assert.Equal(t, c.Values()["foo"], int64(i))
	}
}

func TestMultiple(t *testing.T) {
	c := New()

	for i := 1; i < 100; i++ {
		c.Increment("foo")
		c.Increment("bar")
		assert.Equal(t, c.Values()["foo"], int64(i))
		assert.Equal(t, c.Values()["bar"], int64(i))
	}
}

func TestSimple(t *testing.T) {
	c := New()

	c.Increment("foo")
	c.Increment("foo")
	c.Increment("foo")
	c.Increment("foo")
	c.Increment("bar")
	c.Increment("bar")
	c.Increment("bar")
	c.Increment("qaz")

	assert.Equal(t, c.Values()["foo"], int64(4))
	assert.Equal(t, c.Values()["bar"], int64(3))
	assert.Equal(t, c.Values()["qaz"], int64(1))

}
