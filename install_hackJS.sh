#!/bin/bash

# الترجمة وبناء الأداة
echo "Building hackJS..."
go build -o hackJS hackJS.go

# نقل الأداة إلى مجلد ضمن PATH
echo "Moving hackJS to /usr/local/bin/"
sudo mv hackJS /usr/local/bin/

# إنشاء مجلد لتخزين ملف wordlist إذا لم يكن موجوداً
echo "Creating directory for wordlist..."
mkdir -p ~/bin

# نقل ملف wordlist إلى المجلد
echo "Moving wordlist to ~/bin/"
mv WordList.txt ~/bin/

echo "Installation complete."
