#!/bin/bash
fyne bundle -package assets -name Icon src/gui/assets/icon.png > src/gui/assets/bundled.go
fyne bundle -package assets -name Logo -append src/gui/assets/logo.png >> src/gui/assets/bundled.go
fyne package -name "Result Converter" -sourceDir src/gui -os windows -icon src/gui/assets/icon.png