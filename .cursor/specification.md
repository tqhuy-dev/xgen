# Project Specification

## 1. Overview
**Project Name:** XGen

**Purpose:**  
This project is a **Golang library** providing reusable functions and components that can be integrated into multiple other Go applications.  

**Example:**  
A library that includes array function, string function, data algorithms, and more. 

---

## 2. Architecture

### 2.1 Tech Stack
- **Language:** Go 1.25+
- **Dependencies:**
    - `github.com/stretchr/testify` – testing framework
    - `go.uber.org/zap` – structured logging
    - `github.com/golang-jwt/jwt/v5` – token management
    - `github.com/spf13/viper` – configuration loader
- **Testing:** Standard Go `testing` package with `testify` assertions.

---

## 3. Folder Structure
/lo -> lo folder contains all the logic related to the array functions.
/strings -> stings folder contains all the logic related to the string functions.
/data -> data folder contains all the logic related to the data algorithms.

