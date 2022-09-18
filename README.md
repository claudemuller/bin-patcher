<p align="center">
    <img src="./bin-patcher.png" alt="Bin Patcher" />
</p>

# Bin Patcher

A little Go application that is able to patch binaries given that the source and destination signatures are of the same length e.g. `JE` and `JNE` both being 2 bytes.

That said, ensure that a signature and patch of 5 bytes is passed in to ensure the correct locating of the instruction in question e.g. a signature of `ff08c07409` with a patch of `ff08c07509` where the last 2 bytes represent the instruction to be altered.

# Building

```bash
make build
```

# Cleaning the Project

```bash
make clean
```

# Running

## CLI

```bash
make cli in=<input_bin> out=<output_bin> sig=<signature_to_find> patch=<patch_signature>
```

## GUI

```bash
make gui
```
