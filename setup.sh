#!/bin/bash
echo "Compiling..."

fyne bundle -package assets -name Icon src/gui/assets/icon.png > src/gui/assets/bundled.go
fyne bundle -package assets -name Logo -append src/gui/assets/logo.png >> src/gui/assets/bundled.go
fyne package -appID "Result.Converter" -appVersion 1.2.5 -icon src/gui/assets/icon.png -name "Result Converter" -os windows -release -sourceDir src/gui
