# A luciferchase original
# https://gist.github.com/luciferchase/ede7a276f61a338d5c371eca945164fe

import csv
import re

roll_regex = re.compile("^\\d{8}$")
name_regex = re.compile("([A-Z])\\w+")

# Edit this to the path to your result file
path = ""

with open(path, "r") as f:
    raw_data = f.read()
    raw_data = raw_data.split()

subject_codes = {
    "027": "HISTORY",
    "028": "POLITICAL SCIENCE",
    "029": "GEOGRAPHY",
    "030": "ECONOMICS",
    "041": "MATHEMATICS",
    "042": "PHYSICS",
    "043": "CHEMISTRY",
    "044": "BIOLOGY",
    "048": "PHYSICAL EDUCATION",
    "049": "PAINTING",
    "054": "BUSINESS STUDIES",
    "055": "ACCOUNTANCY",
    "083": "COMPUTER SCIENCE",
    "301": "ENGLISH CORE",
    "302": "HINDI CORE",
    "806": "TOURISM",
    "809": "FOOD PRODUCTION"
}

parsed_data = [
    [
        "ROLL", "GENDER", "NAME", 
        "SUBJECT CODE 1", "SUBJECT NAME", "MARKS OBTAINED", "GRADE", 
        "SUBJECT CODE 2", "SUBJECT NAME", "MARKS OBTAINED", "GRADE", 
        "SUBJECT CODE 3", "SUBJECT NAME", "MARKS OBTAINED", "GRADE", 
        "SUBJECT CODE 4", "SUBJECT NAME", "MARKS OBTAINED", "GRADE", 
        "SUBJECT CODE 5", "SUBJECT NAME", "MARKS OBTAINED", "GRADE", 
        "SUBJECT CODE 6", "SUBJECT NAME", "MARKS OBTAINED", "GRADE",
        "TOTAL", "PERCENTAGE", "RESULT"
    ]
]

for i in range(len(raw_data)):
    # Filter anything which looks like a roll number
    if roll_regex.search(raw_data[i]):
        # Parse individual student's data

        roll, gender = raw_data[i: i + 2]

        # Extremely rare but if a student has five word name
        if (name_regex.search(raw_data[i + 6])):
            name = " ".join(raw_data[i + 2: i + 7])
            i += 3
        # Rare but if a student has four word name
        elif name_regex.search(raw_data[i + 5]):
            name = " ".join(raw_data[i + 2: i + 6])
            i += 2
        # If a student has three word name
        elif name_regex.search(raw_data[i + 4]):
            name = " ".join(raw_data[i + 2: i + 5])
            i += 1     
        # If a student has only one word name
        elif not name_regex.search(raw_data[i + 3]):
            name = raw_data[i + 2]
            i -= 1
        # Default two word names
        else:
            name = " ".join(raw_data[i + 2: i + 4])
        
        result = raw_data[i + 13]
        if (result not in ["PASS", "COMP"]):
            continue

        student_subjects = [[j, subject_codes[j]] for j in raw_data[i + 4: i + 10]]
        marks = [[raw_data[i + j], raw_data[i + j + 1]] for j in range(14, 25, 2)]
        int_marks = [int(i[0]) for i in marks]

        total = str(sum(int_marks))
        percentage = str((int(total) - int(raw_data[i + 24])) / 5)[:4]
        # If percentage is of the form XX.X add 0 to the last
        if (len(percentage) == 3):
            percentage += "0"

        student_data = [roll, gender, name]
        for j in range(6):
            student_data.extend([student_subjects[j][0], student_subjects[j][1], marks[j][0], marks[j][1]])
        student_data.extend([total, percentage, result])
        parsed_data.append(student_data)

with open("class_12th_result.csv", "w", newline = "\n") as f:
    writer = csv.writer(f)
    writer.writerows(parsed_data)