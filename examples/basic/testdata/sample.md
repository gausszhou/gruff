# Gruff Markdown Renderer

A lightweight, high-performance **markdown** renderer for the terminal.

## Features

- *Italic* and **bold** text
- `Inline code` support
- Unordered and ordered lists
- Headings (H1 through H6)

## Usage Example

```go
package main

import (
    "fmt"
    "log"

    "github.com/gausszhou/gruff"
)

func main() {
    out, err := gruff.Render("# Hello World")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Print(out)
}
```

## Formatting Showcase

This paragraph demonstrates **bold**, *italic*, `inline code`, and ***bold italic*** styles working together seamlessly.

1. First ordered item with **bold text**
2. Second ordered item with *italic text*
3. Third ordered item with `inline code`

### Nested Emphasis

Here we have **bold with *italic inside*** and *italic with **bold inside***.

## Code Elements

Inline code like `var x = 42` should stand out with a distinct background color.

## Lists with Mixed Content

- Item with **bold**
- Item with *italic*
- Item with `code`
- Item with ***both***

## Thematic Break

---

Above is a horizontal rule.

A casual introduction. 你好世界!

## Let’s talk about artichokes

The _artichoke_ is mentioned as a garden plant in the 8th century BC by Homer
**and** Hesiod. The naturally occurring variant of the artichoke, the cardoon,
which is native to the Mediterranean area, also has records of use as a food
among the ancient Greeks and Romans. Pliny the Elder mentioned growing of
_carduus_ in Carthage and Cordoba.

> He holds him with a skinny hand,
> ‘There was a ship,’ quoth he.
> ‘Hold off! unhand me, grey-beard loon!’
> An artichoke, dropt he.

--Samuel Taylor Coleridge, [The Rime of the Ancient Mariner][rime]

[rime]: https://poetryfoundation.org/poems/43997/

## Other foods worth mentioning

1. Carrots
1. Celery
1. Tacos
    * Soft
    * Hard
1. Cucumber

## Things to eat today

* [x] Carrots
* [x] Ramen
* [ ] Currywurst

### Power levels of the aforementioned foods

| Name       | Power | Comment          |
| ---        | ---   | ---              |
| Carrots    | 9001  | It’s over 9000?! |
| Ramen      | 9002  | Also over 9000?! |
| Currywurst | 10000 | What?!           |

## Currying Artichokes

Here’s a bit of code in [Haskell](https://haskell.org), because we are fancy.
Remember that to compile Haskell you’ll need `ghc`.

```haskell
module Main where

import Data.List (intercalate)

hello :: String -> String
hello s = "Hello, " <> s <> "."

main :: IO ()
main = putStrLn
     $ intercalate "\n"
     $ hello <$> [ "artichoke", "alcachofa" ]
```

***

_Alcachofa_, if you were wondering, is artichoke in Spanish.

## Final Section

A plain paragraph to close the document with some **formatting** to make sure everything works correctly in the *final* output.

Visit [Gruff on GitHub](https://github.com/gausszhou/gruff) for more information.

## Table

| Feature   | Status | Priority |
|:----------|:------:|:--------:|
| Headings  | ✅     | High     |
| Bold      | ✅     | High     |
| Italic    | ✅     | High     |
| Inline Code | ✅   | High     |
| Lists     | ✅     | High     |
| Links     | ✅     | Medium   |
| Tables    | ✅     | Medium   |
