<p align="center">
    <img src="./src/gui/assets/logo.png" alt="Result Converter Logo"/>
</p>

<p align="center">
    <img src="https://img.shields.io/github/go-mod/go-version/luciferchase/cbse-result-converter" alt="GitHub go.mod Go version"/>
    <img src="https://img.shields.io/github/license/luciferchase/cbse-result-converter" alt="GitHub"/>
    <a href="https://github.com/luciferchase/CBSE-Result-Converter/releases" title="Total Downloads" rel="nofollow"><img src="https://img.shields.io/github/downloads/luciferchase/cbse-result-converter/v0.4.0/total" alt="Total Downloads"/></a>
    <img src="https://img.shields.io/badge/app%20size-10%20MB-brightgreen" alt="App Size"/>
    <img src="https://img.shields.io/github/issues/luciferchase/cbse-result-converter" alt="GitHub issues"/>
</p>

# Introduction

Every year CBSE releases report of all the students in a school. However, in their infinite wisdom, they are still stuck with the old Newspaper Gazette Style Report in a plain .txt file. Analysis in a plain .txt file is not feasible, hence every year IT dept. of the School spends countless manhours to convert that file to a .csv file.      

This simple app converts that .txt file to .csv file in seconds.

<p align="center">
    <img src="./src/gui/assets/main-screen.jpg" alt="Main Screen"/>
</p>

The .csv file converted through this allows easy analysis and comparision. The output file has these fields:

```python
"ROLL", "GENDER", "NAME",
"SUBJECT CODE 1", "SUBJECT NAME", "MARKS OBTAINED", "GRADE",
"SUBJECT CODE 2", "SUBJECT NAME", "MARKS OBTAINED", "GRADE",
"SUBJECT CODE 3", "SUBJECT NAME", "MARKS OBTAINED", "GRADE",
"SUBJECT CODE 4", "SUBJECT NAME", "MARKS OBTAINED", "GRADE",
"SUBJECT CODE 5", "SUBJECT NAME", "MARKS OBTAINED", "GRADE",
"SUBJECT CODE 6", "SUBJECT NAME", "MARKS OBTAINED", "GRADE",
"TOTAL", "PERCENTAGE", "RESULT",
```

# Analysis

The app provides basic analysis as well, but for detailed anaylsis, use the .csv file.

<p align="center">
    <img src="./src/gui/assets/analyse-tab.jpg" alt="Analysis Tab"/>
</p>

# Download

Binaries can be easily download from here [Release](https://github.com/luciferchase/CBSE-Result-Converter/releases).

# Build

Build it from source by cloning this repo and run `setup.sh`. It will build `gui.exe` in the `src/gui` folder (this is a known issue [#2395](https://github.com/fyne-io/fyne/issues/2395)).