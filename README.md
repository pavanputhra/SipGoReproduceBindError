# UDP Bind Error Reproduction

This project is set up to reproduce and investigate UDP binding issues.


## Getting Started

1. Make sure you have Go installed on your system
2. Clone this repository
3. Run the application:
   ```bash
   go run main.go
   ```
4. Execute following command to see the error
  ```bash
  sed 's/$/\r/' invite.txt | nc -u 127.0.0.1 5060
  ```

## License

This project is open source and available under the MIT License. 