# Dependencies: wget, unzip, python3, java
#
# Requires PlaydateOS.pdos to be in this directory

mkdir PlaydateOS
cd PlaydateOS || exit
unzip ../PlaydateOS.pdos
cd System/Settings.pdx || exit

wget https://raw.githubusercontent.com/cranksters/playdate-reverse-engineering/refs/heads/main/tools/pdz.py
wget https://github.com/scratchminer/unluac/releases/download/v2023.03.22/unluac.jar
wget https://github.com/scratchminer/unluac/releases/download/v2023.03.22/unluac.sh

python3 pdz.py -o . -i main.pdz
python3 pdz.py -o . -i crankButton.pdz
rm main.pdz
rm crankButton.pdz

chmod +x unluac.sh
./unluac.sh -r . .

rm unluac.sh
rm unluac.jar
rm pdz.py

rm -r CoreLibs

cp -r . ../../../src
cd ../../..
rm -r PlaydateOS
