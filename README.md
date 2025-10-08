# 🧮 Beremiz (Go Edition)

**Beremiz-Go** is a **stack-based toy programming language** implemented in **Go**.  
It is a reimagining of the original [Beremiz (Lua version)](https://github.com/AdaiasMagdiel/Beremiz), inspired by [Porth](https://gitlab.com/tsoding/porth) by [Alexey Kutepov](https://twitch.tv/tsoding).

The goal of Beremiz is **playfulness and education** — it’s not meant for production use,  
but rather a fun way to explore interpreters, compilers, and stack-based design.

The name **Beremiz** comes from [Beremiz Samir](https://en.wikipedia.org/wiki/Beremiz_Samir),  
_The Man Who Counted_, a character created by Brazilian writer  
[Júlio César de Mello e Souza](https://en.wikipedia.org/wiki/J%C3%BAlio_C%C3%A9sar_de_Mello_e_Souza), known as **Malba Tahan**.

---

## ✨ Features

- ⚙️ **Stack-based execution model**
- 🔢 Basic types: `int`, `float`, `string`, `bool`, `nil`
- ➕ Arithmetic and stack operations
- 🔁 Control flow: `if / elif / else / do / end`, `for / do / end`
- 🧩 `define` system for custom words and constants
- 🔗 String concatenation with `.`
- 🧰 Built-ins: `dup`, `swap`, `over`, `rot`, `pop`, `depth`, `clear`
- 🖨 Output: `write`, `writeln`
- 🧠 Type introspection: `type`
- 💡 REPL with `.help`, `.clear`, and `exit`
- 🧪 Buffer-optimized output (auto-flush in loops)
- 🎨 Syntax highlighting:
  - [VS Code extension](https://marketplace.visualstudio.com/items?itemName=Adaias-Magdiel.beremiz)
  - [Sublime Text syntax file](./.syntax-highlight/sublime-text/Beremiz.sublime-syntax)

---

## 🚀 Getting Started

### 🧩 Building

Requirements:

- **Go 1.22+**

```bash
git clone https://github.com/AdaiasMagdiel/beremiz-go.git
cd beremiz-go
go build -o beremiz ./cmd/beremiz/main.go  # or make
```

### ▶️ Running a File

```bash
./beremiz examples/hello_world.brz
```

### 💬 REPL Mode

```bash
./beremiz-go
```

Then type:

```
> "Hello" writeln
Hello
> 2 3 + writeln
5
```

Exit with:

```
> exit
```

---

## 🧠 Language Overview

### 🔢 Numbers

```beremiz
42                   writeln   # decimal
-123                 writeln   # negative decimal
+123                 writeln   # positive decimal
0                    writeln   # zero
1_000_000            writeln   # decimal with underscore
3.14                 writeln   # floating point
0.0                  writeln   # zero floating point
.5                   writeln   # floating point with leading dot
5.                   writeln   # floating point with trailing dot
-3.14                writeln   # negative floating point
0xFF                 writeln   # hexadecimal
0x1a2b3c             writeln   # hexadecimal
0XDEADBEEF           writeln   # hexadecimal (uppercase)
0777                 writeln   # octal
0o123                writeln   # modern octal
0b1010               writeln   # binary
0B11110000           writeln   # binary uppercase
0b1100_0101          writeln   # binary with underscore
9223372036854775807  writeln   # max 64-bit integer
-9223372036854775808 writeln   # min 64-bit integer
```

---

### ➕ Operators

```beremiz
1 2 + writeln            # -> 3
10 5 - writeln           # -> 5
42 3.14 * writeln        # -> 131.88
10 5 / writeln           # -> 2
0b1010 0b0011 + writeln  # -> 13
```

Also supported:

- `**` → exponentiation
- `%` → modulo
- `<`, `>`, `<=`, `>=`, `eq`, `neq`
- `. ` → concatenation (`true`, `false`, `nil` become `"true"`, `"false"`, `"nil"`)

---

### 🧩 Stack Operations

| Word    | Description             | Example                      |
| ------- | ----------------------- | ---------------------------- |
| `dup`   | Duplicate top           | `5 dup -> [5 5]`             |
| `swap`  | Swap top two            | `1 2 swap -> [2 1]`          |
| `over`  | Copy second to top      | `[a b] -> [a b a]`           |
| `rot`   | Rotate top three        | `[a b c] -> [b c a]`         |
| `pop`   | Drop top                | `[1 2 pop -> [1]]`           |
| `depth` | Push current stack size | `[1 2 3] depth -> [1 2 3 3]` |
| `clear` | Clear entire stack      | `[1 2 3] clear -> []`        |

---

### 🧭 Conditionals

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

---

### 🔁 Loops

```beremiz
1_000

for dup 0 neq do
    dup writeln
    1 -
end


10
for dup 0 neq do
    if dup 5 eq do
        "five!" writeln
    else
        "Number: " write
        dup writeln
    end
    1 -
end
```

---

### 🧮 Define — Constants and Functions

```beremiz
define PI
  3.14159
end

define square
  dup *
end

define circle_area
  dup * PI *
end

define add_two
  2 +
end

5 square writeln          # 25
10 circle_area writeln    # 314.159
8 add_two writeln         # 10

define a 5 end
define b a 10 * end
define c b 2 * end

c writeln                 # 100
```

---

### 🌀 Fibonacci Example

```beremiz
define fibonacci
over over +   # stack: 15 0 1 0 -> 15 0 1 0 1 -> 15 0 1 1
rot pop       # stack: 15 1 1 0 -> 15 1 1
end

15   # stack: 15
0 1  # stack: 15 0 1
for
    rot  # stack: 0 1 15
    dup  # stack: 0 1 15 15
    rot  # stack: 0 15 15 1
    dup  # stack: 0 15 15 1 1
    rot  # stack: 0 15 1 1 15
    <    # stack: 0 15 1 bool
do       # stack: 0 15 1
    rot swap   # stack: 15 1 0 -> 15 0 1
    fibonacci  # stack: 15 1 1
    swap       # stack: 15 1 1
    dup        # stack: 15 1 1 1
    writeln    # stack: 15 1 1
    swap       # stack: 15 1 1
end
```

---

### 🔔 FizzBuzz Example

```beremiz
1

define multiple_by_3
    dup 3 % 0 eq
end

define multiple_by_5
    dup 5 % 0 eq
end

for dup 100 <= do
    dup write ": " write

    if
        multiple_by_3
        swap
        multiple_by_5
        rot and
    do
        "FizzBuzz"
    elif multiple_by_3 do
        "Fizz"
    elif multiple_by_5 do
        "Buzz"
    else
        ""
    end

    writeln

    1 +
end
```

---

## 🖥 Editor Support

- 🧩 **VS Code** — [Official Extension](https://marketplace.visualstudio.com/items?itemName=Adaias-Magdiel.beremiz)
- 🎨 **Sublime Text** — [Syntax File](./.syntax-highlight/sublime-text/Beremiz.sublime-syntax)

---

## 🧰 Project Structure

| File           | Description                       |
| -------------- | --------------------------------- |
| `lexer/`       | Tokenization of source code       |
| `parser/`      | Stack-based interpreter           |
| `tokens.go`    | Token and keyword definitions     |
| `main.go`      | CLI & REPL entry point            |
| `pathutils.go` | Path resolution helpers           |
| `err.go`       | Error formatting                  |
| `examples/`    | Complete runnable `.brz` programs |

---

## 🧭 Roadmap

- [x] REPL
- [x] `define`, `if`, `for` blocks
- [x] Rich literals and concatenation
- [x] Buffered output with smart flush
- [ ] `import` for module support
- [ ] Standard library (`math`, `string`, etc.)

---

## 📜 License

This project is licensed under the **GNU General Public License v3.0** (GPLv3).
See the full text in [LICENSE](LICENSE).
