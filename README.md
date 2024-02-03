#Snow
---
![](https://github.com/MattRattigan/Snow/blob/master/snow.webp)

##Snow is is a Go-based software designed for file encryption and decryption, utilizing SQLite for data storage and management. 
---

###Key Features

    File Encryption and Decryption: Securely encrypts and decrypts files using AES-GCM algorithm.
    User Management: Manages user data securely, leveraging SQLite for storage.
    File Extension Handling: Dynamically handles file extensions during encryption/decryption processes.
    Windows Registry Integration: Registers custom file extensions for encrypted files in the Windows Registry.

###Command-Line Flags

    -username: Username for login.
    -filepath: Path to the file for encryption/decryption.
    -dirpath: Path to the directory.
    -e: Encrypt file.
    -d: Decrypt file.
    -ext: Name of the file extension.

###Notes
    Run the application with administrative privileges on Windows to modify the registry.
    Tested on Windows and Ubuntu Linux
