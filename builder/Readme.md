Builder
=======
A simple HTML templating engine written in Go. Works in similar manner to C macros. 

To run the builder, first build the application using `go build` and run the binary in the root directory of the website. All pages and components of the website must be located in a directory called `src`.

### HTML
There are two types of html files:

- **pages**: A HTML file that represents a full website page
- **components**: An HTML template that can be imported and reused by other HTML files

All pages are built and outputted with the same directory nesting as the source file had.

A component is created by creating an HTML file and a file with the extension `.component.json`. The syntax of this file is:

```json
{
    "name": "<name of component>",
    "templatePath": "<path to the HMTL file relative to this file>",
}
```

Once a component is defined it can be used in any other HTML file. This is done by using the component include directive:

```html
<!--@<name of component>{<parameters>}-->
```

Think of the parameters section as a JSON object (this is how it is parsed). However, every value must be a string.

A component can use parameters by prepending a '`#`' character before the parameter name almost anywhere within the component's HTML file. A parameter cannot be used in a component include directive except for the values of the parameters JSON object.

For example, `<!--@message{"content": "#message"}-->` is allowed, but `<!--@message{"#message": "hello!"-->` will not result in any parameter substitution (though it is still valid syntax).

### CSS
For simplicity and caching efficiency, **all** `*.css` files in the `src` directory are concatenated into the `out/styles.css` file. This means that you will need to avoid using conflicting selectors between pages and components.

### Other Files
For other file types such as JavaScript and images you can include the resource in the ouput by referencing within a `resources` array in the `*.component.json` file:

```json
{
    "name": "Example",
    "templatePath": "example.html",
    "resources": [
        "<path to the file relative to this file>"
    ]
}
```

The resource will be copied to a matching directory within the output.

### TODO
- Minify CSS and JavaScript in the output
- Automatically link the stylesheet on pages
- Allow including resources in pages