# Beremiz (Go Edition)

**Beremiz-Go** is a stack-based toy programming language implemented in **Go**.  
It is a reimagining of the original [Beremiz (Lua version)](https://github.com/AdaiasMagdiel/Beremiz), inspired by [Porth](https://gitlab.com/tsoding/porth) by [Alexey Kutepov](https://twitch.tv/tsoding).

The goal of Beremiz is **playfulness and education** — it’s not meant for production use, but as a way to learn about interpreters, compilers, and stack-based language design.

The name **Beremiz** comes from [Beremiz Samir](https://en.wikipedia.org/wiki/Beremiz_Samir), _The Man Who Counted_, a character created by Brazilian writer [Malba Tahan](https://en.wikipedia.org/wiki/J%C3%BAlio_C%C3%A9sar_de_Mello_e_Souza).

---

## ✨ Features

- Concise **stack-based execution model**
- Basic types: integers, floats, strings, booleans, `nil`
- Arithmetic and stack operations
- Conditional execution (`if … elif … else … end`)
- Loops (`for … end`)
- Built-in words like `dup`, `eq`, `neq`, `write`, `writeln`
- Clear and minimal syntax
- 🔌 **VS Code support** — [syntax highlighting extension](https://marketplace.visualstudio.com/items?itemName=Adaias-Magdiel.beremiz)
- 🔌 **Sublime Text support** — [syntax highlighting file](./.syntax-highlight/sublime-text/Beremiz.sublime-syntax)

---

## 🚀 Getting Started

### Run a program

```bash
make
./dist/beremiz ./examples/hello_world.brz
```

### Hello World

```beremiz
"Hello, world!" writeln
```

---

## 🧭 Language Tour

### Numbers

```beremiz
42     writeln
3.14   writeln
0b1010 writeln
0xFF   writeln
0777   writeln
0o123  writeln
1_000  writeln
```

### Conditionals

```beremiz
5

if dup 5 eq do
    "Equal to 5"
end

writeln
```

```beremiz
10

if dup 7 eq do
    "Equal to 7"
else
    "Not equal to 7"
end

writeln
```

```beremiz
3

if dup 1 eq do
    "One"
elif dup 2 eq do
    "Two"
elif dup 3 eq do
    "Three"
else
    "Another number"
end

writeln
```

### Loops

```beremiz
1_000

for dup 0 neq do
    dup writeln
    1 -
end
```

### Operators

```beremiz
1 2 + writeln            # -> 3
10 5 - writeln           # -> 5
42 3.14 * writeln        # -> 131.88
10 5 / writeln           # -> 2
```

---

## 🔢 Numeric Literals

- **Decimal**: `42`, `-123`, `+123`, `1_000_000`
- **Floating point**: `3.14`, `.5`, `5.`, `-3.14`
- **Hexadecimal**: `0xFF`, `0XDEADBEEF`
- **Octal**: `0777`, `0o123`
- **Binary**: `0b1010`, `0B11110000`, `0b1100_0101`
- **64-bit extremes**: `9223372036854775807`, `-9223372036854775808`

---

## 🖥 VS Code Extension

For the best experience, install the official [Beremiz VS Code extension](https://marketplace.visualstudio.com/items?itemName=Adaias-Magdiel.beremiz):

- Syntax highlighting for `.brz` files
- Keyword, number, string, operator, and comment coloring
- Customizable scopes for built-ins (`dup`, `writeln`, `neq`, etc.)

---

## 📌 Status

🚧 **Work in Progress**

- ✅ Hello World
- ✅ Numbers (decimal, float, binary, hex, octal)
- ✅ Strings (with escapes)
- ✅ Booleans (`true`, `false`), `nil`
- ✅ Conditionals (`if / elif / else / end`)
- ✅ Loops (`for`)
- ✅ VS Code syntax highlighting
- 🚧 User-defined words / functions
- 🚧 Standard library

---

## 📄 License

This project is licensed under the [MIT License](LICENSE).
