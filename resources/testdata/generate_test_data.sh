#!/bin/bash

###########################################################
#
# Copyright (c) 2017 codeliveroil. All rights reserved.
#
# This work is licensed under the terms of the MIT license.
# For a copy, see <https://opensource.org/licenses/MIT>.
#
###########################################################

set -e

if [ "$(basename $(pwd))" != "testdata" ]; then
  echo "Run the script from the folder in which it is contained."
  exit
fi

cd ../..
go build
cd -
../../img -w 80 -o color_matrix.sh color_matrix.png
../../img -o disposalBackground.sh disposalBackground.gif
../../img -o disposalNone.sh disposalNone.gif
../../img -o disposalNoneTransparency.sh disposalNoneTransparency.gif
../../img -o disposalUnspecified.sh disposalUnspecified.gif
../../img -l 3 -s 2 -w 60 -o all.sh disposalNone.gif

echo ""
echo "Test data generated."
read -p "Now, visually verify the test data. Press any key to continue." x
echo $0
for f in $(ls *.sh); do
  [[ $0 = *"$f"* ]] && continue #ignore generation script.
  clear
  ./$f
  read -p "Hit any key to proceed to the next image." x
done
