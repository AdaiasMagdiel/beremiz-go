# Beremiz (Go Edition)

**Beremiz-Go** is a stack-based toy programming language implemented in **Go**.
It is a reimagining of the original [Beremiz](https://github.com/AdaiasMagdiel/Beremiz) (Lua version), inspired by [Porth](https://gitlab.com/tsoding/porth) by [Alexey Kutepov](https://twitch.tv/tsoding).

The goal of Beremiz is **playfulness and education** — it’s not intended for production, but to learn about interpreters, compilers, and stack-based language design.

The name **Beremiz** comes from [Beremiz Samir](https://en.wikipedia.org/wiki/Beremiz_Samir), _The Man Who Counted_, a character created by Brazilian writer [Malba Tahan](https://en.wikipedia.org/wiki/J%C3%BAlio_C%C3%A9sar_de_Mello_e_Souza).

---

## Features

- Concise **stack-based execution model**
- Basic types: integers, floats, strings, booleans, nil
- Arithmetic and stack operations
- Conditional execution (`if … else … endif`)
- Keywords and simple built-in words
- Clear and minimal syntax

---

## Getting Started

### Run a program

```bash
go run ./cmd/beremiz ./examples/hello_world.brz
```

### Hello World

```beremiz
"Hello, world!" writeln
```

---

## Language Tour

### Numbers

```beremiz
42 writeln
3.14 writeln
0b1010 writeln
0xFF writeln
```

### Conditionals

```beremiz
true if
    "It's true"
else
    "It's false"
endif

writeln
```

```beremiz
0 if
    "Zero is truthy"
else
    "Zero is falsy"
endif

writeln
```

### Keywords

```beremiz
true writeln
false writeln
nil writeln
```

---

## Status

🚧 **Work in Progress**

- ✅ Hello world
- ✅ Numbers (int, float, binary, hex, octal)
- ✅ Strings (with escapes)
- ✅ Booleans (`true`, `false`), `nil`
- ✅ Conditionals (`if / else / endif`)
- 🚧 Loops
- 🚧 User-defined words / functions
- 🚧 Standard library

---

## License

This project is licensed under the terms of the [GNU General Public License v3](LICENSE).
