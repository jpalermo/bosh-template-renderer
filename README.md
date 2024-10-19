# Bosh Template Renderer

Bosh Template Renderer is intended to be the default template rendering library for BOSH releases.

The main goals are a stable grammar and fast rendering time.

# Development Usage

```bash
echo '{"properties": {"variable": {"complex": "thing"}}}' | go run . examples/example.btl
```

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

### So how do I do things like validate URLs?

A direction moving forward may be to not do input validate inside of template rendering. Input validation doesn't
really seem like the job of template rendering, it's just where we've always done it because we had the full
expressiveness of Ruby available within template rendering.

Possible ways forward:
- Push input validate down into the release job itself. So during job startup or pre-start the job could do the
  validation at that point. This has the upside of allowing release authors to do validation however they want, and 
  the downside that it happens much later in the deployment process. At best, the bootstrap VM is going to be down. If 
  it's a singleton VM, all of them.
- Create input validate rules within the job spec itself. Similar to the existing `default` value properties can
  have, there could be additional things added such as `required` or even more complex things. Upside of this is
  that it happens even earlier than template rendering, but the downside is you can only use what's already been built,
  and right now that's nothing.
