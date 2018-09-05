package httpr_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/go-playground/form"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phogolabs/http/httpr"
)

var _ = Describe("Binder", func() {
	var (
		r    *http.Request
		body *bytes.Buffer
	)

	BeforeEach(func() {
		body = &bytes.Buffer{}
		r = httptest.NewRequest("GET", "http://example.com", body)
		r.Header.Add("Content-Type", "application/json")
	})

	It("descodes the request successfully", func() {
		t := T{Name: "Jack"}
		Expect(json.NewEncoder(body).Encode(&t)).To(Succeed())

		t2 := T{}
		Expect(httpr.Bind(r, &t2)).To(Succeed())
		Expect(t2).To(Equal(t))
	})

	Context("when a form post is made", func() {
		BeforeEach(func() {
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			r.Method = "POST"
		})

		It("descodes the request successfully", func() {
			t := T{Name: "Jack"}

			values, err := form.NewEncoder().Encode(&t)
			Expect(err).To(Succeed())

			body.WriteString(values.Encode())

			t2 := T{}
			Expect(httpr.Bind(r, &t2)).To(Succeed())
			Expect(t2).To(Equal(t))
		})
	})

	Context("when the binder fails", func() {
		It("returns the error", func() {
			t := T{Name: "Jack", Err: "Oh no"}
			Expect(json.NewEncoder(body).Encode(&t)).To(Succeed())

			t2 := T{}
			Expect(httpr.Bind(r, &t2)).To(MatchError("Oh no"))
		})
	})

	Context("when the validation fails", func() {
		It("returns the error", func() {
			t := T{}
			Expect(json.NewEncoder(body).Encode(&t)).To(Succeed())

			t2 := T{}
			Expect(httpr.Bind(r, &t2)).To(MatchError("Key: 'T.name' Error:Field validation for 'name' failed on the 'required' tag"))
		})
	})
})
