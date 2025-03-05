# DFA Vending Machine  

A **Deterministic Finite Automata (DFA) based Vending Machine** implemented in **C**, using a `.txt` configuration file to define the DFA states and transitions.

## 📜 Table of Contents  

- [📜 Table of Contents](#-table-of-contents)  
- [👥 Team Members](#-team-members)  
- [🛠 Installation](#-installation)  
- [🚀 Usage](#-usage)  
- [📂 Project Structure](#-project-structure)  
- [📜 DFA Configuration File](#-dfa-configuration-file)  
- [📝 Code Documentation](#-code-documentation)  
- [📃 License](#-license)  

---

## 👥 Team Members  
- **[Daffa Indra Wibowo (NIM)]** - DFA Config  
- **[M. Argya Vityasy (23/522547/PA/22475)]** - Code

---

## 🛠 Installation  

### Prerequisites  
Ensure you have the following installed:  
- **GCC Compiler** (`gcc`)  
- **Make** (for easy compilation)  

### Steps  
1. **Clone the repository**  
   ```sh
   git clone https://github.com/yourusername/DFA-vending-machine.git
   cd DFA-vending-machine
   ```
2. **Compile the code**  
   ```sh
   make
   ```
3. **Run the vending machine**  
   ```sh
   ./vending_machine
   ```

---

## 🚀 Usage  

1. Run the executable.  
2. Follow the on-screen instructions to insert money and select an item.  
3. The DFA logic processes inputs based on the predefined state transitions.  
4. If sufficient money is inserted, the product is dispensed.
5. If money is more than the product price, dispense the product and return the remaining money 
---

## 📂 Project Structure  

```
DFA-vending-machine/
│── src/
│   ├── main.c             # Main program logic
│   ├── dfa.c              # DFA processing functions
│── config/
│   ├── vending_config.txt # DFA state transitions
│── README.md              # Project documentation
│── Makefile               # Build automation
```

---

## 📜 DFA Configuration File  

The DFA configuration file (`vending_config.txt`) defines:  
- **States**
- **Alphabets**
- **Start state**
- **Transitions**  
- **Accepted inputs**  

Example format:  
```
States: S0, S1000, S2000, S5000, S100000
Alphabet: 1000, 2000, 5000, 10000
Start: S0
Accept: S3000, S4000, ...
Transitions:
S0 1000 S1000
S0 2000 S2000
...
```

---

## 📝 Code Documentation  
### Later will contain the screenshots of the program 

### `main.c`  
- Handles user input and calls DFA functions.  

### `dfa.c`  
- Reads the DFA configuration file.  
- Implements state transitions.  
- Checks if the final state is reached.  

### `dfa.h`  
- Header file with function declarations.  

---

## 📃 License  

This project is licensed under the **MIT License**.  

---
