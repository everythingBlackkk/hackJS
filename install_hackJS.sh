#!/bin/bash

echo "Building hackJS..."
go build -o hackJS hackJS.go

echo "Moving hackJS to /usr/local/bin/"
sudo mv hackJS /usr/local/bin/

echo "Creating directory for wordlist..."
mkdir -p ~/bin

echo "Moving wordlist to ~/bin/"
mv WordList.txt ~/bin/

echo "Installation complete."
