package renderer_test

import (
	"github.com/Jeffail/gabs/v2"
	"github.com/cloudfoundry/bosh-template-renderer/renderer"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"strings"
)

var _ = Describe("Rendering", func() {
	Describe("plain text", func() {
		It("returns the same string", func() {
			text := `line1
line2 
line3 `

			template, _ := renderer.Parse(strings.NewReader(text))
			Expect(template.Render(nil)).To(Equal(text))
		})
	})

	Describe("interpolating", func() {
		Context("properties", func() {
			It("string values", func() {
				text := `{{p.key}}`
				data, _ := gabs.ParseJSON([]byte(`{"properties": {"key": "value"}}`))
				template, _ := renderer.Parse(strings.NewReader(text))
				Expect(template.Render(data)).To(Equal("value"))
			})

			It("number values", func() {
				text := `{{p.integer}} {{p.floating}}`
				data, _ := gabs.ParseJSON([]byte(`{"properties": {"integer": 5, "floating": 3.14}}`))
				template, _ := renderer.Parse(strings.NewReader(text))
				Expect(template.Render(data)).To(Equal("5 3.14"))
			})

			It("array values", func() {
				text := `{{p.data}}`
				data, _ := gabs.ParseJSON([]byte(`{"properties": {"data": [1, "2"]}}`))
				template, _ := renderer.Parse(strings.NewReader(text))
				Expect(template.Render(data)).To(Equal(`[1,"2"]`))
			})

			It("object values", func() {
				text := `{{p.data}}`
				data, _ := gabs.ParseJSON([]byte(`{"properties": {"data": {"nested": [1, "2"]}}}`))
				template, _ := renderer.Parse(strings.NewReader(text))
				Expect(template.Render(data)).To(Equal(`{"nested":[1,"2"]}`))
			})

			It("nested searching", func() {
				text := `{{p.one.two.three}}`
				data, _ := gabs.ParseJSON([]byte(`{"properties": {"one": {"two": {"three": 3}}}}`))
				template, err := renderer.Parse(strings.NewReader(text))
				Expect(err).ToNot(HaveOccurred())
				Expect(template.Render(data)).To(Equal("3"))
			})

			It("errors when the property does not exist", func() {
				text := `{{p.nope.not-there}}`
				data, _ := gabs.ParseJSON([]byte(`{"properties": {"one": 1}}`))
				template, err := renderer.Parse(strings.NewReader(text))
				Expect(err).ToNot(HaveOccurred())
				_, err = template.Render(data)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("p.nope.not-there did not match any provided properties"))
			})

			It("does not interpolate single braces", func() {
				text := `{p.one}}`
				data, _ := gabs.ParseJSON([]byte(`{"properties": {"one": 1}}`))
				template, err := renderer.Parse(strings.NewReader(text))
				Expect(err).ToNot(HaveOccurred())
				Expect(template.Render(data)).To(Equal("{p.one}}"))
			})

			It("braces can be escaped", func() {
				text := `\{{p.one}}`
				data, _ := gabs.ParseJSON([]byte(`{"properties": {"one": 1}}`))
				template, err := renderer.Parse(strings.NewReader(text))
				Expect(err).ToNot(HaveOccurred())
				Expect(template.Render(data)).To(Equal("{{p.one}}"))
			})
		})

		Context("specs", func() {
			It("spec values", func() {
				text := `{{spec.thing_one}} {{spec.thing_two}}`
				data, _ := gabs.ParseJSON([]byte(`{"spec": {"thing_one": 1, "thing_two": 2}}`))
				template, err := renderer.Parse(strings.NewReader(text))
				Expect(err).ToNot(HaveOccurred())
				Expect(template.Render(data)).To(Equal("1 2"))
			})

			It("errors when the property does not exist", func() {
				text := `{{spec.not_a_thing}}`
				data, _ := gabs.ParseJSON([]byte(`{"spec": {}}`))
				template, err := renderer.Parse(strings.NewReader(text))
				Expect(err).ToNot(HaveOccurred())
				_, err = template.Render(data)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("spec.not_a_thing did not match any provided properties"))
			})

		})
	})
})
