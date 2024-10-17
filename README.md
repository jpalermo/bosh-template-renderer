# Bosh Template Renderer

Bosh Template Renderer is intended to be the default template rendering library for BOSH releases.

The main goals are a stable grammar and fast rendering time.

# Questions

---
### Why not use ERB templates?

ERB templates are rendered with the current version of Ruby the BOSH Director is running with. As we've upgraded ruby
versions over the years, we've seen several BOSH releases fail to properly render templates with newer versions.
The goal of this library is to provide a stable template language that we can guarantee compatibility with going forward.

### What about template language X?

Using a drop in template library gives us a much more feature rich language for release authors, but comes with the
same problems as using ERB. We have no guarantees that the library authors won't make breaking changes or abandon
the project entirely. This library does have dependencies that could change, but we have the ability to replace them
or work around them because those dependencies aren't the provided interface.

### Will this help speed up BOSH template rendering?

Maybe... We have not yet run any benchmarking tests, but the lexing/parsing library we are using (participle), does seem
to have good benchmarking of it's own as well as some options for possibly pre-compiling the templates so re-rendering
them is faster.

### ERB does X and this does not!

That's not a question... but will always be true. The currently supported grammar is very small and since we have to
implement the lex-ing and parsing, it will never be huge. We are certainly open to adding features though.
