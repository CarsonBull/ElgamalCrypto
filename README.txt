Name: Alex Carson Bull
Email: Carson.Bull@gmail.com

Elgamal Crypto

Description: This project implements a 32 bit public key crypto system. It is build in the Go language. It has three different operating modes that can be
set in the main function of the program.

Key Generation: this provides you with two test files, pubkey.txt and prikey.txt, that contain the information required to encrypt and decrypt a message.

Encryption: Encrytipn requires that the user provide a file to encrypt and will create a text file called ctext.txt. It also needs a pubkey.txt file

Decryption: Decryption will take a ctext.txt file and output plain text to the file that is passed to the function. It requires a prikey.txt file

Note: This program should be able to handle ascii but may have issues with utf-8
      I did have an issue with windows defender thinking this was a virus because of file names. Issues can be solved by turning off windows defender. 

Files: (only needed to run the program)
publicKey.go    Program
message.txt     File that will be encrypted( you can set this to whatever you want)
README.txt      Tells you about the program
