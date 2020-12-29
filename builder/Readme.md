Builder
=======
A simple HTML templating engine written in Go. Works in similar manner to C macros. 

To run the builder:

- Navigate to the root directory of the builder
- Run `go get` to fetch dependencies
- Run `go build`
- Run the binary while in the root directory of the website.

All pages and components of the website must be located in a directory called `src`. The output will be in a directory (that will be created if necessary) called `out` (**WARNING:** each execution of the builder will clear the contents of the `out` directory).

### HTML
There are two types of html files:

- **pages**: A normal `.html` file that represents a full website page
- **components**: A building block that used within other pages or components

All pages are built and outputted with the same directory nesting as the source file had.

A component is defined by creating a HTML file that is named matching the pattern `<component name>.component.html`. A component named `test.component.html` is named `test`. No other component in a project can have the same name.

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
All other file types except `.html` and `.css` files will be linked into the output directory with the same directory nesting as the source directory.